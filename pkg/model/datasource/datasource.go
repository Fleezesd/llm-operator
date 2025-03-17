package datasource

import (
	"context"
	"io"
)

type DataSource interface {
	Stat(ctx context.Context, info any) error
	Remove(ctx context.Context, info any) error
	ReadFile(ctx context.Context, info any) (io.ReadCloser, error)
	StatFile(ctx context.Context, info any) (any, error)
	GetTags(ctx context.Context, info any) (map[string]string, error)
	ListObjects(ctx context.Context, source string, info any) (any, error)
}
