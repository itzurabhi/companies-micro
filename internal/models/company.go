package models

// Type of a company
type CompanyType string

// All type of valid companies
const (
	CompanyCorporations       = "Corporations"
	CompanyNonProfit          = "NonProfit"
	CompanyCooperative        = "Cooperative"
	CompanySoleProprietorship = "Sole Proprietorship"
)

var CompanyTypeNameMap = map[string]CompanyType{
	CompanyCorporations:       CompanyCorporations,
	CompanyNonProfit:          CompanyNonProfit,
	CompanyCooperative:        CompanyCooperative,
	CompanySoleProprietorship: CompanySoleProprietorship,
}

// • ID (uuid) required
// • Name (15 characters) required - unique
// • Description (3000 characters) optional
// • Amount of Employees (int) required
// • Registered (boolean) required
// • Type (Corporations | NonProfit | Cooperative | Sole Proprietorship) required
type Company struct {
	ID                string      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name              string      `gorm:"unique;not null"`
	Description       string      `gorm:"not null"`
	AmountOfEmployees int         `gorm:"not null"`
	Registered        bool        `gorm:"not null"`
	Type              CompanyType `gorm:"not null"`
}
