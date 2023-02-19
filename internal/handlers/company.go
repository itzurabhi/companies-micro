package handlers

import (
	"github.com/gofiber/fiber"
	"github.com/itzurabhi/companies-micro/internal/logic"
	"github.com/itzurabhi/companies-micro/internal/models"

	"github.com/go-playground/validator/v10"

	log "github.com/sirupsen/logrus"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("validateCompanyTypeEnum", validateCompanyTypeEnum)
}

func validateCompanyTypeEnum(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	default:
		return false
	case models.CompanyCooperative:
		fallthrough
	case models.CompanyNonProfit:
		fallthrough
	case models.CompanyCorporations:
		fallthrough
	case models.CompanySoleProprietorship:
	}

	return true
}

type CompanyHandler struct {
	companyLogic *logic.CompanyLogic
}

type CompanyRequest struct {
	Name              string `validate:"required,min=1,max=15"`
	Description       string `validate:"max=3000"`
	AmountOfEmployees int    `validate:"required"`
	Registered        bool   `validate:"required"`
	Type              string `validate:"validateCompanyTypeEnum"`
}

func CreateCompanyHandler(companyLogic *logic.CompanyLogic) *CompanyHandler {
	return &CompanyHandler{
		companyLogic: companyLogic,
	}
}

func (handler *CompanyHandler) Create(c *fiber.Ctx) {

	req := new(CompanyRequest)

	if err := c.BodyParser(req); err != nil {
		log.Error("Company:Create", "could not decode body", err)
		_ = writeError(c, err)
	}

	if err := validate.Struct(req); err != nil {
		resp := createInvalidBodyError(err.Error())
		_ = writeError(c, resp)
	}

	cType, ok := models.CompanyTypeNameMap[req.Type]

	if !ok {
		log.Error("Company:Create", "company type mapping error, from :", req.Type)
		resp := createInternalError("could not create company")
		_ = writeError(c, resp)
	}

	data := models.Company{
		Name:        req.Name,
		Description: req.Description,
		Registered:  req.Registered,
		Type:        cType,
	}

	created, err := handler.companyLogic.Create(c.Context(), data)

	if err != nil {
		_ = writeError(c, err)
	}

	_ = writeSuccessJSON(c, &created)
}

func (handler *CompanyHandler) Get(c *fiber.Ctx) {
	_ = c.Params("id")
	panic("not implemented")
}

func (handler *CompanyHandler) Patch(c *fiber.Ctx) {
	_ = c.Params("id")
	panic("not implemented")
}

func (handler *CompanyHandler) Delete(c *fiber.Ctx) {
	_ = c.Params("id")
	panic("not implemented")
}
