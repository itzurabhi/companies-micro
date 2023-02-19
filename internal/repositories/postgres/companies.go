package postgres

import (
	"context"

	"github.com/itzurabhi/companies-micro/internal/models"
	"github.com/itzurabhi/companies-micro/internal/repositories"
	"gorm.io/gorm"
)

type pgcompanies struct {
	conn *gorm.DB
}

func CreateCompaniesRepo(conn *gorm.DB) repositories.Companies {
	return pgcompanies{
		conn: conn,
	}
}

// Create implements repositories.Companies
func (repo pgcompanies) Create(ctx context.Context, data models.Company) (models.Company, error) {
	tx := repo.conn.WithContext(ctx).Model(&models.Company{}).Create(&data)
	if tx.Error != nil {
		return data, tx.Error
	}

	return data, nil
}

// Delete implements repositories.Companies
func (repo pgcompanies) Delete(ctx context.Context, id string) error {
	return repo.conn.WithContext(ctx).Model(&models.Company{}).Delete(models.Company{
		ID: id,
	}).Error
}

// Get implements repositories.Companies
func (repo pgcompanies) Get(ctx context.Context, id string) (models.Company, error) {
	comp := models.Company{}
	err := repo.conn.WithContext(ctx).Model(&models.Company{}).Find(&comp, id).Error
	return comp, err
}

// Patch implements repositories.Companies
func (repo pgcompanies) Patch(ctx context.Context, id string, req models.Company) (models.Company, error) {
	err := repo.conn.WithContext(ctx).Save(&req).Error
	return req, err
}
