package ec2

import "context"

type Stub struct {
	Inst            *Instance
	Err             error
	GetInstanceName string

	StartCalled bool
	StopCalled  bool
	StartID     string
	StopID      string
	WaitIP      string
	WaitErr     error
}

func (s *Stub) GetInstance(ctx context.Context, name string) (*Instance, error) {
	s.GetInstanceName = name
	return s.Inst, s.Err
}

func (s *Stub) StartInstance(ctx context.Context, id string) error {
	s.StartCalled = true
	s.StartID = id
	return s.Err
}

func (s *Stub) StopInstance(ctx context.Context, id string) error {
	s.StopCalled = true
	s.StopID = id
	return s.Err
}

func (s *Stub) WaitRunning(ctx context.Context, id string) (string, error) {
	return s.WaitIP, s.WaitErr
}
