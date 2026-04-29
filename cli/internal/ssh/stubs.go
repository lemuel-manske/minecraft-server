package ssh

import (
	"context"
	"io"
)

type Stub struct {
	WaitSSHErr     error
	RunErr         error
	OutputErr      error
	StreamErr      error
	UploadFileErr  error
	UploadTarErr   error
	DownloadTarErr error

	OutputResults map[string]string

	RunCmds          []string
	OutputCmds       []string
	StreamCmds       []string
	UploadFilePaths  [][2]string
	UploadTarPaths   [][2]string
	DownloadTarPaths [][2]string
	WaitForSSHCalled bool
}

func (s *Stub) WaitForSSH(ctx context.Context) error {
	s.WaitForSSHCalled = true
	return s.WaitSSHErr
}

func (s *Stub) Run(cmd string) error {
	s.RunCmds = append(s.RunCmds, cmd)
	return s.RunErr
}

func (s *Stub) Output(cmd string) (string, error) {
	s.OutputCmds = append(s.OutputCmds, cmd)
	if s.OutputErr != nil {
		return "", s.OutputErr
	}
	return s.OutputResults[cmd], nil
}

func (s *Stub) Stream(cmd string, w io.Writer) error {
	s.StreamCmds = append(s.StreamCmds, cmd)
	return s.StreamErr
}

func (s *Stub) UploadFile(localPath, remotePath string) error {
	s.UploadFilePaths = append(s.UploadFilePaths, [2]string{localPath, remotePath})
	return s.UploadFileErr
}

func (s *Stub) UploadTar(localDir, remotePath string) error {
	s.UploadTarPaths = append(s.UploadTarPaths, [2]string{localDir, remotePath})
	return s.UploadTarErr
}

func (s *Stub) DownloadTar(remotePath, localPath string) error {
	s.DownloadTarPaths = append(s.DownloadTarPaths, [2]string{remotePath, localPath})
	return s.DownloadTarErr
}
