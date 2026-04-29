package ec2

import "context"

type Instance struct {
	ID       string
	PublicIP string
	State    string
	Profile  string
}

type Client interface {
	GetInstance(ctx context.Context, name string) (*Instance, error)
	StartInstance(ctx context.Context, id string) error
	StopInstance(ctx context.Context, id string) error
	WaitRunning(ctx context.Context, id string) (ip string, err error)
}
