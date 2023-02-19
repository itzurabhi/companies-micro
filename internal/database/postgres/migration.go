package postgres

import (
	"github.com/itzurabhi/companies-micro/internal/models"
	"gorm.io/gorm"
)

func Migrate(conn *gorm.DB) error {

	return conn.AutoMigrate(&models.Company{})
}
