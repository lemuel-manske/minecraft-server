package git

import (
	"fmt"
	"os/exec"
	"strings"
)

var _ Client = (*GitImpl)(nil)

type GitImpl struct {
	dir string
}

func New(dir string) *GitImpl {
	return &GitImpl{dir: dir}
}

func (g *GitImpl) CommitBackup(tag string) error {
	cmds := [][]string{
		{"git", "-C", g.dir, "add", "."},
		{"git", "-C", g.dir, "commit", "-m", "backup: " + tag},
		{"git", "-C", g.dir, "tag", tag},
	}

	for _, c := range cmds {
		out, err := exec.Command(c[0], c[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %s", c, out)
		}
	}

	return nil
}

func (g *GitImpl) ListTags() ([]string, error) {
	out, err := exec.Command("git", "-C", g.dir, "tag", "--sort=-creatordate").Output()
	if err != nil {
		return nil, fmt.Errorf("git tag list: %w", err)
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil, nil
	}

	return strings.Split(raw, "\n"), nil
}

func (g *GitImpl) Prune(keep int) error {
	tags, err := g.ListTags()
	if err != nil {
		return err
	}

	if len(tags) <= keep {
		return nil
	}

	for _, tag := range tags[keep:] {
		out, err := exec.Command("git", "-C", g.dir, "tag", "-d", tag).CombinedOutput()
		if err != nil {
			return fmt.Errorf("delete tag %s: %s", tag, out)
		}
	}

	return nil
}
