package kafka

import "github.com/itzurabhi/companies-micro/internal/repositories"

type companiesEventBus struct {
}

// PostEvent implements repositories.EventBus
func (companiesEventBus) PostEvent(items ...interface{}) error {
	panic("unimplemented")
}

func CreateCompaniesEventBus() repositories.EventBus {
	return companiesEventBus{}
}
