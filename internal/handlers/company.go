package handlers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/itzurabhi/companies-micro/internal/logic"
	"github.com/itzurabhi/companies-micro/internal/models"
	"github.com/itzurabhi/companies-micro/internal/repositories"

	"github.com/go-playground/validator/v10"

	"github.com/sirupsen/logrus"
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
	Registered        *bool  `validate:"required"`
	Type              string `validate:"validateCompanyTypeEnum"`
}

func CreateCompanyHandler(companyLogic *logic.CompanyLogic) *CompanyHandler {
	return &CompanyHandler{
		companyLogic: companyLogic,
	}
}

func (handler *CompanyHandler) AddRoutes(app *fiber.App, middleWares ...fiber.Handler) {
	companiesRoute := app.Group("companies", middleWares...)
	companiesRoute.Get("/:id", handler.Get)
	companiesRoute.Post("/", handler.Create)
	companiesRoute.Patch("/:id", handler.Patch)
	companiesRoute.Delete("/:id", handler.Delete)
}

func (handler *CompanyHandler) Create(c *fiber.Ctx) error {

	req := new(CompanyRequest)

	if err := c.BodyParser(req); err != nil {
		logrus.Error("Company:Create", "could not decode body", err)
		return writeError(c, err)
	}

	if err := validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		resp := createInvalidBodyError(errs[0].Error())
		return writeError(c, resp)
	}

	cType, ok := models.CompanyTypeNameMap[req.Type]

	if !ok {
		logrus.Error("Company:Create", "company type mapping error, from :", req.Type)
		resp := createInternalError("could not create company")
		return writeError(c, resp)
	}

	data := models.Company{
		Name:              req.Name,
		Description:       req.Description,
		Registered:        *req.Registered,
		AmountOfEmployees: req.AmountOfEmployees,
		Type:              cType,
	}

	created, err := handler.companyLogic.Create(c.Context(), data)

	if err != nil {

		if errors.Is(err, repositories.ErrorRecordAlreadyExist) {
			return writeError(c, err, http.StatusBadRequest)
		}

		return writeError(c, err)
	}

	return writeSuccessJSON(c, &created, http.StatusCreated)
}

func (handler *CompanyHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := handler.companyLogic.Get(c.Context(), id)

	if err != nil {
		return writeError(c, err)
	}

	return writeSuccessJSON(c, &data)
}

func (handler *CompanyHandler) Patch(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(CompanyRequest)

	if err := c.BodyParser(req); err != nil {
		logrus.Error("Company:Create", "could not decode body", err)
		return writeError(c, err)
	}

	if err := validate.Struct(req); err != nil {
		resp := createInvalidBodyError(err.Error())
		return writeError(c, resp)
	}

	cType, ok := models.CompanyTypeNameMap[req.Type]

	if !ok {
		logrus.Error("Company:Create", "company type mapping error, from :", req.Type)
		resp := createInternalError("could not create company")
		return writeError(c, resp)
	}

	data := models.Company{
		ID:                id,
		Name:              req.Name,
		Description:       req.Description,
		Registered:        *req.Registered,
		AmountOfEmployees: req.AmountOfEmployees,
		Type:              cType,
	}

	created, err := handler.companyLogic.Patch(c.Context(), id, data)

	if err != nil {
		return writeError(c, err)
	}

	return writeSuccessJSON(c, &created)
}

func (handler *CompanyHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := handler.companyLogic.Delete(c.Context(), id); err != nil {
		logrus.Error("Company:Delete", "failed for ", id, err)
	}
	return c.Status(http.StatusAccepted).SendString("")
}
