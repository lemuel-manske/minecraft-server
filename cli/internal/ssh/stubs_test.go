package ssh_test

import (
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/ssh"
)

func TestStubImplementsClient(t *testing.T) {
	var _ ssh.Client = (*ssh.Stub)(nil)
}
