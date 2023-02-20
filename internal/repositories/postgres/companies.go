package postgres

import (
	"context"
	"errors"

	"github.com/itzurabhi/companies-micro/internal/models"
	"github.com/itzurabhi/companies-micro/internal/repositories"
	"github.com/jackc/pgx/v5/pgconn"
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

		var pgErr *pgconn.PgError

		// error with a mapped database code
		if errors.As(tx.Error, &pgErr) && pgErr.Code == "23505" {
			return data, repositories.ErrorRecordAlreadyExist
		}

		return data, tx.Error
	}
	return data, nil
}

// Delete implements repositories.Companies
func (repo pgcompanies) Delete(ctx context.Context, id string) error {
	err := repo.conn.WithContext(ctx).Model(&models.Company{}).Delete(models.Company{
		ID: id,
	}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repositories.ErrorRecordNotFound
	}

	return err
}

// Get implements repositories.Companies
func (repo pgcompanies) Get(ctx context.Context, id string) (models.Company, error) {
	comp := models.Company{}
	err := repo.conn.WithContext(ctx).Model(&models.Company{}).Find(&comp, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return comp, repositories.ErrorRecordNotFound
	}
	return comp, err
}

// Patch implements repositories.Companies
func (repo pgcompanies) Patch(ctx context.Context, id string, req models.Company) (models.Company, error) {
	err := repo.conn.WithContext(ctx).Save(&req).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return req, repositories.ErrorRecordNotFound
	}
	return req, err
}
