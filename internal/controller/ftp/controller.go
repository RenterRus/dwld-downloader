package ftp

import "context"

type Sender interface {
	Loader(ctx context.Context)
}
