package tf

import "context"

type Client interface {
	Apply(ctx context.Context, vars []string) error
	Destroy(ctx context.Context) error
	Output(ctx context.Context, key string) (string, error)
	SetVerbose(v bool)
}
