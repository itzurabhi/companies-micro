package repositories

type EventBus interface {
	PostEvent(items ...interface{}) error
}
