package tf

import "context"

type Stub struct {
	ApplyErr   error
	DestroyErr error
	OutputErr  error
	Outputs    map[string]string

	ApplyCalled   bool
	DestroyCalled bool
	ApplyVars     []string
	OutputKeys    []string
}

func (s *Stub) SetVerbose(_ bool) {}

func (s *Stub) Apply(ctx context.Context, vars []string) error {
	s.ApplyCalled = true
	s.ApplyVars = vars
	return s.ApplyErr
}

func (s *Stub) Destroy(ctx context.Context) error {
	s.DestroyCalled = true
	return s.DestroyErr
}

func (s *Stub) Output(ctx context.Context, key string) (string, error) {
	s.OutputKeys = append(s.OutputKeys, key)
	if s.OutputErr != nil {
		return "", s.OutputErr
	}
	return s.Outputs[key], nil
}
