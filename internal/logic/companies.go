package logic

import (
	"context"

	"github.com/itzurabhi/companies-micro/internal/models"
	"github.com/itzurabhi/companies-micro/internal/repositories"
	"github.com/sirupsen/logrus"
)

type CompanyLogic struct {
	companyRepo   repositories.Companies
	companyEvents repositories.EventBus
}

func CreateCompanyLogic(repo repositories.Companies, evbus repositories.EventBus) *CompanyLogic {
	return &CompanyLogic{
		companyRepo:   repo,
		companyEvents: evbus,
	}
}

func (cl *CompanyLogic) Create(ctx context.Context, data models.Company) (models.Company, error) {
	data, err := cl.companyRepo.Create(ctx, data)

	if err != nil {
		return data, err
	}

	if err := cl.companyEvents.PostEvent(models.CompanyEvent{
		Type: models.CompanyCreatedEvent,
		Uuid: data.ID,
	}); err != nil {
		logrus.Error("posting create event failed:", err)
	}

	return data, err
}

func (cl *CompanyLogic) Patch(ctx context.Context, id string, data models.Company) (models.Company, error) {

	data, err := cl.companyRepo.Patch(ctx, id, data)

	if err != nil {
		return data, err
	}

	if err := cl.companyEvents.PostEvent(models.CompanyEvent{
		Type: models.CompanyPatchedEvent,
		Uuid: data.ID,
	}); err != nil {
		logrus.Error("posting patch event failed:", id, err)
	}

	return data, err

}

func (cl *CompanyLogic) Delete(ctx context.Context, id string) error {

	err := cl.companyRepo.Delete(ctx, id)

	if err != nil {
		return err
	}

	if err := cl.companyEvents.PostEvent(models.CompanyEvent{
		Type: models.CompanyDeletedEvent,
		Uuid: id,
	}); err != nil {
		logrus.Error("posting patch event failed:", id, err)
	}
	return err
}

func (cl *CompanyLogic) Get(ctx context.Context, id string) (models.Company, error) {
	return cl.companyRepo.Get(ctx, id)
}
