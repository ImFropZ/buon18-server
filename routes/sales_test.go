package routes_test

import (
	"context"
	"net/http/httptest"
	"path/filepath"
	"server/config"
	"server/database"
	"server/middlewares"
	"server/routes"
	"server/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestSalesRoutes(t *testing.T) {
	config.GetConfigInstance()
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(
			filepath.Join("..", "database", "dev_scripts", "create-schema.sh"),
			filepath.Join("..", "database", "dev_scripts", "seed.sh"),
			filepath.Join("..", "database", "dev_scripts", "seed-customer.sh"),
			filepath.Join("..", "database", "dev_scripts", "seed-quotation.sh"),
		),
		postgres.BasicWaitStrategies(),
	)
	assert.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	DB := database.InitSQL(connectionString)

	router := gin.Default()
	router.Use(middlewares.Authenticate(DB))
	routes.Sales(router, DB)

	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       "admin@buon18.com",
		Role:        "bot",
		Permissions: []string{"FULL_ACCESS"},
	})
	assert.NoError(t, err)

	t.Run("SuccessGetListOfQuotations", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/quotations", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"quotations":[{"id":1,"name":"Quotation 1","creation_date":"2021-01-01T00:00:00Z","validity_date":"2021-01-31T00:00:00Z","discount":50,"amount_delivery":100,"status":"quotation","total_amount":350,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":1,"name":"Item 1","description":"Item 1 description","price":100,"discount":0,"amount_total":100},{"id":2,"name":"Item 2","description":"Item 2 description","price":200,"discount":0,"amount_total":200}]},{"id":2,"name":"Quotation 2","creation_date":"2021-02-01T00:00:00Z","validity_date":"2021-02-28T00:00:00Z","discount":100,"amount_delivery":200,"status":"quotation_sent","total_amount":550,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":3,"name":"Item 3","description":"Item 3 description","price":500,"discount":50,"amount_total":450}]},{"id":3,"name":"Quotation 3","creation_date":"2021-03-01T00:00:00Z","validity_date":"2021-03-31T00:00:00Z","discount":150,"amount_delivery":300,"status":"quotation_sent","total_amount":1050,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":4,"name":"Item 4","description":"Item 4 description","price":1000,"discount":100,"amount_total":900}]}]}}`
		expectedXTotalCountHeader := "3"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FilterSuccessGetListOfQuotations", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/quotations?status-eq=quotation_sent", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"quotations":[{"id":2,"name":"Quotation 2","creation_date":"2021-02-01T00:00:00Z","validity_date":"2021-02-28T00:00:00Z","discount":100,"amount_delivery":200,"status":"quotation_sent","total_amount":550,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":3,"name":"Item 3","description":"Item 3 description","price":500,"discount":50,"amount_total":450}]},{"id":3,"name":"Quotation 3","creation_date":"2021-03-01T00:00:00Z","validity_date":"2021-03-31T00:00:00Z","discount":150,"amount_delivery":300,"status":"quotation_sent","total_amount":1050,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":4,"name":"Item 4","description":"Item 4 description","price":1000,"discount":100,"amount_total":900}]}]}}`
		expectedXTotalCountHeader := "2"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SortSuccessGetListOfQuotations", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/quotations?sort-name=desc", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"quotations":[{"id":3,"name":"Quotation 3","creation_date":"2021-03-01T00:00:00Z","validity_date":"2021-03-31T00:00:00Z","discount":150,"amount_delivery":300,"status":"quotation_sent","total_amount":1050,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":4,"name":"Item 4","description":"Item 4 description","price":1000,"discount":100,"amount_total":900}]},{"id":2,"name":"Quotation 2","creation_date":"2021-02-01T00:00:00Z","validity_date":"2021-02-28T00:00:00Z","discount":100,"amount_delivery":200,"status":"quotation_sent","total_amount":550,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":3,"name":"Item 3","description":"Item 3 description","price":500,"discount":50,"amount_total":450}]},{"id":1,"name":"Quotation 1","creation_date":"2021-01-01T00:00:00Z","validity_date":"2021-01-31T00:00:00Z","discount":50,"amount_delivery":100,"status":"quotation","total_amount":350,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":1,"name":"Item 1","description":"Item 1 description","price":100,"discount":0,"amount_total":100},{"id":2,"name":"Item 2","description":"Item 2 description","price":200,"discount":0,"amount_total":200}]}]}}`
		expectedXTotalCountHeader := "3"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})
}
