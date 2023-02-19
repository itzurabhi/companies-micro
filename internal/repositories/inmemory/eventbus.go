package inmemory

import (
	"context"

	"github.com/itzurabhi/companies-micro/internal/repositories"
	"github.com/sirupsen/logrus"
)

type inMemoryEventBus struct{}

// PostEvent implements repositories.EventBus
func (inMemoryEventBus) PostEvent(_ context.Context, events ...interface{}) error {
	logrus.Println("InMemoryEventBus:", events)
	return nil
}

func CreateInMemoryEventBus() repositories.EventBus {
	return inMemoryEventBus{}
}
