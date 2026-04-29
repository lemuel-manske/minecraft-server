package ssh

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

var _ Client = (*SSHImpl)(nil)

type SSHImpl struct {
	host    string
	keyPath string
	client  *gossh.Client
}

func New(host, keyPath string) *SSHImpl {
	return &SSHImpl{host: host, keyPath: keyPath}
}

func (s *SSHImpl) dial() (*gossh.Client, error) {
	key, err := os.ReadFile(s.keyPath)
	if err != nil {
		return nil, fmt.Errorf("read key: %w", err)
	}

	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse key: %w", err)
	}

	cfg := &gossh.ClientConfig{
		User:            "ubuntu",
		Auth:            []gossh.AuthMethod{gossh.PublicKeys(signer)},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	return gossh.Dial("tcp", s.host+":22", cfg)
}

func (s *SSHImpl) ensureConnected() error {
	if s.client != nil {
		return nil
	}

	c, err := s.dial()
	if err != nil {
		return err
	}

	s.client = c
	return nil
}

func (s *SSHImpl) WaitForSSH(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for SSH after 3 minutes")
		default:
		}

		c, err := s.dial()
		if err == nil {
			s.client = c
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for SSH after 3 minutes")
		case <-time.After(5 * time.Second):
		}
	}
}

func (s *SSHImpl) Run(cmd string) error {
	if err := s.ensureConnected(); err != nil {
		return err
	}

	sess, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	return sess.Run(cmd)
}

func (s *SSHImpl) Output(cmd string) (string, error) {
	if err := s.ensureConnected(); err != nil {
		return "", err
	}

	sess, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	out, err := sess.Output(cmd)
	return string(out), err
}

func (s *SSHImpl) Stream(cmd string, w io.Writer) error {
	if err := s.ensureConnected(); err != nil {
		return err
	}

	sess, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	sess.Stdout = w
	return sess.Run(cmd)
}

func (s *SSHImpl) UploadFile(localPath, remotePath string) error {
	if err := s.ensureConnected(); err != nil {
		return err
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	sess, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	stdin, err := sess.StdinPipe()
	if err != nil {
		return err
	}

	if err := sess.Start("cat > " + remotePath); err != nil {
		return err
	}

	if _, err := stdin.Write(data); err != nil {
		return err
	}

	stdin.Close()
	return sess.Wait()
}

func (s *SSHImpl) UploadTar(localDir, remotePath string) error {
	if err := s.ensureConnected(); err != nil {
		return err
	}

	sess, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	stdin, err := sess.StdinPipe()
	if err != nil {
		return err
	}

	if err := sess.Start(fmt.Sprintf("tar xzf - -C %s", remotePath)); err != nil {
		return err
	}

	if err := writeTar(stdin, localDir); err != nil {
		return err
	}

	stdin.Close()
	return sess.Wait()
}

func (s *SSHImpl) DownloadTar(remotePath, localPath string) error {
	if err := s.ensureConnected(); err != nil {
		return err
	}

	sess, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	stdout, err := sess.StdoutPipe()
	if err != nil {
		return err
	}

	if err := sess.Start(fmt.Sprintf("tar czf - -C %s .", remotePath)); err != nil {
		return err
	}

	if err := extractTar(stdout, localPath); err != nil {
		return err
	}

	return sess.Wait()
}

func writeTar(w io.Writer, dir string) error {
	gz := gzip.NewWriter(w)
	tw := tar.NewWriter(gz)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = rel

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})
	if err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return gz.Close()
}

func extractTar(r io.Reader, dir string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dir, hdr.Name)

		if hdr.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		f, err := os.Create(target)
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, tr); err != nil {
			f.Close()
			return err
		}

		f.Close()
	}

	return nil
}
