package repositories

import "context"

type EventBus interface {
	PostEvent(ctx context.Context, items ...interface{}) error
}
