package repositories

import (
	"context"

	"github.com/itzurabhi/companies-micro/internal/models"
)

type Companies interface {
	Create(ctx context.Context, data models.Company) (models.Company, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Company, error)
	Patch(ctx context.Context, id string, req models.Company) (models.Company, error)
}
