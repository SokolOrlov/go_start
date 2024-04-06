package http

import "context"

type Adapter interface {
	Serve() error
	Shutdown(ctx context.Context)
}
