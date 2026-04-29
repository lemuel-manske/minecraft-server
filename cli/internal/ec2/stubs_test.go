package ec2_test

import (
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/ec2"
)

func TestStubImplementsClient(t *testing.T) {
	var _ ec2.Client = (*ec2.Stub)(nil)
}
