package tf_test

import (
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/tf"
)

func TestStubImplementsClient(t *testing.T) {
	var _ tf.Client = (*tf.Stub)(nil)
}
