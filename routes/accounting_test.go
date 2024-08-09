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

func TestAccountingRoutes(t *testing.T) {
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
			filepath.Join("..", "database", "dev_scripts", "seed-payment-term.sh"),
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
	routes.Accounting(router, DB)

	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       "admin@buon18.com",
		Role:        "bot",
		Permissions: []string{"FULL_ACCESS"},
	})
	assert.NoError(t, err)

	t.Run("SuccessGetListOfQuotations", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/accounting/payment-terms", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"payment_terms":[{"id":1,"name":"Net 30","description":"Net 30","lines":[{"id":1,"sequence":1,"value_amount_percent":100,"number_of_days":30}]},{"id":2,"name":"Net 60","description":"Net 60","lines":[{"id":2,"sequence":1,"value_amount_percent":100,"number_of_days":60}]},{"id":3,"name":"Net 90","description":"Net 90","lines":[{"id":3,"sequence":1,"value_amount_percent":100,"number_of_days":90}]},{"id":4,"name":"30% Now, Balance 60 Days","description":"Pay 30% now, balance due in 60 days","lines":[{"id":4,"sequence":1,"value_amount_percent":30,"number_of_days":0},{"id":5,"sequence":2,"value_amount_percent":70,"number_of_days":60}]}]}}`
		expectedXTotalCountHeader := "4"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FilterSuccessGetListOfQuotations", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/accounting/payment-terms?name-like=30", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"payment_terms":[{"id":1,"name":"Net 30","description":"Net 30","lines":[{"id":1,"sequence":1,"value_amount_percent":100,"number_of_days":30}]},{"id":4,"name":"30% Now, Balance 60 Days","description":"Pay 30% now, balance due in 60 days","lines":[{"id":4,"sequence":1,"value_amount_percent":30,"number_of_days":0},{"id":5,"sequence":2,"value_amount_percent":70,"number_of_days":60}]}]}}`
		expectedXTotalCountHeader := "2"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessGetQuotationById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/accounting/payment-terms/1", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"payment_term":{"id":1,"name":"Net 30","description":"Net 30","lines":[{"id":1,"sequence":1,"value_amount_percent":100,"number_of_days":30}]}}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("NotFoundGetQuotationById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/accounting/payment-terms/0", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"payment term not found","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})
}
