package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/itzurabhi/companies-micro/internal/logic"
	"github.com/itzurabhi/companies-micro/internal/models"
	"github.com/itzurabhi/companies-micro/internal/repositories"
	mockRepos "github.com/itzurabhi/companies-micro/mocks/internal_/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupCompanyWithRepos(comRepo repositories.Companies, evRepo repositories.EventBus) (*fiber.App, *CompanyHandler) {
	app := fiber.New(fiber.Config{Prefork: false})
	compLogic := logic.CreateCompanyLogic(comRepo, evRepo)
	compHandler := CreateCompanyHandler(compLogic)
	compHandler.AddRoutes(app)
	return app, compHandler
}

func TestCompanyHandler_Create(t *testing.T) {

	// invalid input cases, mock should uncover unexpected calls to repos
	{
		companyNoOpRepo := mockRepos.NewCompanies(t)
		eventNoOpRepo := mockRepos.NewEventBus(t)
		app, _ := setupCompanyWithRepos(companyNoOpRepo, eventNoOpRepo)

		tests := []struct {
			description  string                 // description of the test case
			route        string                 // route path to test
			method       string                 // http method name
			body         map[string]interface{} // body of the request
			expectedCode int                    // expected HTTP status code
		}{
			{"create with empty body", "/companies", http.MethodPost, map[string]interface{}{}, http.StatusBadRequest},
			{"create with missing Name", "/companies", http.MethodPost, map[string]interface{}{
				"Description":       "Valid description",
				"AmountOfEmployees": 10,
				"Registered":        false,
				"Type":              "Corporations",
			}, http.StatusBadRequest},
			{"create with missing AmountOfEmployees", "/companies", http.MethodPost, map[string]interface{}{
				"Description": "Valid description",
				"Name":        "ValidName",
				"Registered":  false,
				"Type":        "Corporations",
			}, http.StatusBadRequest},
			{"create with missing Registered", "/companies", http.MethodPost, map[string]interface{}{
				"Description":       "Valid description",
				"AmountOfEmployees": 10,
				"Name":              "ValidName",
				"Type":              "Corporations",
			}, http.StatusBadRequest},
			{"create with missing Type", "/companies", http.MethodPost, map[string]interface{}{
				"Description":       "Valid description",
				"AmountOfEmployees": 10,
				"Name":              "ValidName",
				"Registered":        false,
			}, http.StatusBadRequest},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				bodyBuffer := new(bytes.Buffer)
				_ = json.NewEncoder(bodyBuffer).Encode(test.body)
				req := httptest.NewRequest(test.method, test.route, bodyBuffer)
				req.Header.Add("Content-Type", "application/json")
				resp, _ := app.Test(req, -1)
				assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
			})
		}
	}

	// valid input cases, mock should verify calls to repo methods
	{
		companySuccessRepo := mockRepos.NewCompanies(t)
		companySuccessRepo.On("Create", mock.Anything, mock.Anything).Return(models.Company{}, nil)
		eventSuccessRepo := mockRepos.NewEventBus(t)
		eventSuccessRepo.On("PostEvent", mock.Anything, mock.Anything).Return(nil)
		app, _ := setupCompanyWithRepos(companySuccessRepo, eventSuccessRepo)

		tests := []struct {
			description  string                 // description of the test case
			route        string                 // route path to test
			method       string                 // http method name
			body         map[string]interface{} // body of the request
			expectedCode int                    // expected HTTP status code
		}{
			{"create with all proper fields", "/companies", http.MethodPost, map[string]interface{}{
				"Description":       "Valid description",
				"Name":              "ValidName",
				"AmountOfEmployees": 10,
				"Registered":        false,
				"Type":              "Corporations",
			}, http.StatusCreated},
			{"create with all proper fields", "/companies", http.MethodPost, map[string]interface{}{
				"Description":       "Valid description",
				"Name":              "ValidName",
				"AmountOfEmployees": 10,
				"Registered":        false,
				"Type":              "Corporations",
			}, http.StatusCreated}}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				bodyBuffer := new(bytes.Buffer)
				_ = json.NewEncoder(bodyBuffer).Encode(test.body)
				req := httptest.NewRequest(test.method, test.route, bodyBuffer)
				req.Header.Add("Content-Type", "application/json")
				resp, _ := app.Test(req, -1)
				assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
			})
		}
	}
}
