package ssh

import (
	"context"
	"io"
)

type Client interface {
	WaitForSSH(ctx context.Context) error
	Run(cmd string) error
	Output(cmd string) (string, error)
	Stream(cmd string, w io.Writer) error
	UploadFile(localPath, remotePath string) error
	UploadTar(localDir, remotePath string) error
	DownloadTar(remotePath, localPath string) error
}
