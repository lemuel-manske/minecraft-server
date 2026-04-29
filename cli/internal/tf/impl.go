package tf

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

var _ Client = (*TFImpl)(nil)

type TFImpl struct {
	dir     string
	verbose bool
}

func New(dir string, verbose bool) *TFImpl {
	return &TFImpl{dir: dir, verbose: verbose}
}

func (c *TFImpl) SetVerbose(v bool) {
	c.verbose = v
}

func (c *TFImpl) Apply(ctx context.Context, vars []string) error {
	args := []string{"apply", "-auto-approve"}
	for _, v := range vars {
		args = append(args, "-var="+v)
	}
	return c.run(ctx, args...)
}

func (c *TFImpl) Destroy(ctx context.Context) error {
	return c.run(ctx, "destroy", "-auto-approve")
}

func (c *TFImpl) Output(ctx context.Context, key string) (string, error) {
	cmd := exec.CommandContext(ctx, "terraform", "output", "-raw", key)
	cmd.Dir = c.dir

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("terraform output %s: %w", key, err)
	}

	return string(out), nil
}

func (c *TFImpl) run(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = c.dir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if c.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		if !c.verbose {
			return fmt.Errorf("terraform %s: %s", args[0], stderr.String())
		}
		return fmt.Errorf("terraform %s failed", args[0])
	}

	return nil
}
