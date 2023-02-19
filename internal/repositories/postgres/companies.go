package postgres

import (
	"context"

	"github.com/itzurabhi/companies-micro/internal/models"
	"github.com/itzurabhi/companies-micro/internal/repositories"
	"gorm.io/gorm"
)

type pgcompanies struct {
}

func CreateCompaniesRepo(conn *gorm.DB) repositories.Companies {
	return pgcompanies{}
}

// Create implements repositories.Companies
func (pgcompanies) Create(ctx context.Context, data models.Company) (models.Company, error) {
	panic("unimplemented")
}

// Delete implements repositories.Companies
func (pgcompanies) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// Get implements repositories.Companies
func (pgcompanies) Get(ctx context.Context, id string) (models.Company, error) {
	panic("unimplemented")
}

// Patch implements repositories.Companies
func (pgcompanies) Patch(ctx context.Context, id string, req models.Company) (models.Company, error) {
	panic("unimplemented")
}
