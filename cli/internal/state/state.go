package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const filename = ".mc-state"

type State struct {
	Profile    string `json:"profile"`
	InstanceID string `json:"instance_id"`
	InstanceIP string `json:"instance_ip"`
}

func Load(dir string) (*State, error) {
	data, err := os.ReadFile(filepath.Join(dir, filename))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var s State
	return &s, json.Unmarshal(data, &s)
}

func Save(dir string, s *State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, filename), data, 0644)
}

func Clear(dir string) error {
	err := os.Remove(filepath.Join(dir, filename))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
