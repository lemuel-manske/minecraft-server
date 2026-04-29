package git_test

import (
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/git"
)

func TestStubImplementsClient(t *testing.T) {
	var _ git.Client = (*git.Stub)(nil)
}
