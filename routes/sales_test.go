package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"server/config"
	"server/database"
	"server/middlewares"
	"server/models"
	"server/models/sales"
	"server/routes"
	"server/utils"
	"testing"
	"time"

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
			filepath.Join("..", "database", "dev_scripts", "001_create-schema.sh"),
			filepath.Join("..", "database", "dev_scripts", "002_seed.sh"),
			filepath.Join("..", "database", "dev_scripts", "100_seed-customer.sh"),
			filepath.Join("..", "database", "dev_scripts", "101_seed-quotation.sh"),
			filepath.Join("..", "database", "dev_scripts", "102_seed-payment-term.sh"),
			filepath.Join("..", "database", "dev_scripts", "103_seed-order.sh"),
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
	routes.Sales(router, &database.Connection{
		DB:     DB,
		Valkey: nil,
	})

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

		expectedBodyJSON := `{"code":200,"message":"","data":{"quotations":[{"id":1,"name":"Quotation 1","creation_date":"2021-01-01T00:00:00Z","validity_date":"2021-01-31T00:00:00Z","discount":50,"amount_delivery":100,"status":"quotation","total_amount":350,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":1,"name":"Item 1","description":"Item 1 description","price":100,"discount":0,"amount_total":100},{"id":2,"name":"Item 2","description":"Item 2 description","price":200,"discount":0,"amount_total":200}]},{"id":2,"name":"Quotation 2","creation_date":"2021-02-01T00:00:00Z","validity_date":"2021-02-28T00:00:00Z","discount":100,"amount_delivery":200,"status":"quotation_sent","total_amount":550,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":3,"name":"Item 3","description":"Item 3 description","price":500,"discount":50,"amount_total":450}]},{"id":3,"name":"Quotation 3","creation_date":"2021-03-01T00:00:00Z","validity_date":"2021-03-31T00:00:00Z","discount":150,"amount_delivery":300,"status":"quotation_sent","total_amount":1050,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":4,"name":"Item 4","description":"Item 4 description","price":1000,"discount":100,"amount_total":900}]},{"id":4,"name":"Quotation 4","creation_date":"2021-04-01T00:00:00Z","validity_date":"2021-04-30T00:00:00Z","discount":200,"amount_delivery":400,"status":"sales_order","total_amount":2050,"customer":{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},"items":[{"id":5,"name":"Item 5","description":"Item 5 description","price":2000,"discount":150,"amount_total":1850}]},{"id":5,"name":"Quotation 5","creation_date":"2021-05-01T00:00:00Z","validity_date":"2021-05-31T00:00:00Z","discount":250,"amount_delivery":500,"status":"sales_order","total_amount":3050,"customer":{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},"items":[{"id":6,"name":"Item 6","description":"Item 6 description","price":3000,"discount":200,"amount_total":2800}]},{"id":6,"name":"Quotation 6","creation_date":"2021-06-01T00:00:00Z","validity_date":"2021-06-30T00:00:00Z","discount":300,"amount_delivery":600,"status":"cancelled","total_amount":4050,"customer":{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},"items":[{"id":7,"name":"Item 7","description":"Item 7 description","price":4000,"discount":250,"amount_total":3750}]}]}}`
		expectedXTotalCountHeader := "6"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FilterSuccessGetListOfQuotations", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/quotations?status:eq=quotation_sent", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"quotations":[{"id":2,"name":"Quotation 2","creation_date":"2021-02-01T00:00:00Z","validity_date":"2021-02-28T00:00:00Z","discount":100,"amount_delivery":200,"status":"quotation_sent","total_amount":550,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":3,"name":"Item 3","description":"Item 3 description","price":500,"discount":50,"amount_total":450}]},{"id":3,"name":"Quotation 3","creation_date":"2021-03-01T00:00:00Z","validity_date":"2021-03-31T00:00:00Z","discount":150,"amount_delivery":300,"status":"quotation_sent","total_amount":1050,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":4,"name":"Item 4","description":"Item 4 description","price":1000,"discount":100,"amount_total":900}]}]}}`
		expectedXTotalCountHeader := "2"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessGetQuotationById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/quotations/1", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"quotation":{"id":1,"name":"Quotation 1","creation_date":"2021-01-01T00:00:00Z","validity_date":"2021-01-31T00:00:00Z","discount":50,"amount_delivery":100,"status":"quotation","total_amount":350,"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},"items":[{"id":1,"name":"Item 1","description":"Item 1 description","price":100,"discount":0,"amount_total":100},{"id":2,"name":"Item 2","description":"Item 2 description","price":200,"discount":0,"amount_total":200}]}}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("NotFoundGetQuotationById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/quotations/0", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"quotation not found","data":null}`
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessGetListOfOrders", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/orders", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"orders":[{"id":1,"name":"Order 1","commitment_date":"2021-04-05T00:00:00Z","note":"","quotation":{"id":4,"name":"Quotation 4","creation_date":"2021-04-01T00:00:00Z","validity_date":"2021-04-30T00:00:00Z","discount":200,"amount_delivery":400,"status":"sales_order","total_amount":2050,"customer":{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},"items":[{"id":5,"name":"Item 5","description":"Item 5 description","price":2000,"discount":150,"amount_total":1850}]},"payment_term":{"id":4,"name":"30% Now, Balance 60 Days","description":"Pay 30% now, balance due in 60 days","lines":[{"id":4,"sequence":1,"value_amount_percent":30,"number_of_days":0},{"id":5,"sequence":2,"value_amount_percent":70,"number_of_days":60}]}},{"id":2,"name":"Order 2","commitment_date":"2021-05-05T00:00:00Z","note":"","quotation":{"id":5,"name":"Quotation 5","creation_date":"2021-05-01T00:00:00Z","validity_date":"2021-05-31T00:00:00Z","discount":250,"amount_delivery":500,"status":"sales_order","total_amount":3050,"customer":{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},"items":[{"id":6,"name":"Item 6","description":"Item 6 description","price":3000,"discount":200,"amount_total":2800}]},"payment_term":{"id":2,"name":"Net 60","description":"Net 60","lines":[{"id":2,"sequence":1,"value_amount_percent":100,"number_of_days":60}]}},{"id":3,"name":"Order 3","commitment_date":"2021-06-05T00:00:00Z","note":"","quotation":{"id":6,"name":"Quotation 6","creation_date":"2021-06-01T00:00:00Z","validity_date":"2021-06-30T00:00:00Z","discount":300,"amount_delivery":600,"status":"cancelled","total_amount":4050,"customer":{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},"items":[{"id":7,"name":"Item 7","description":"Item 7 description","price":4000,"discount":250,"amount_total":3750}]},"payment_term":{"id":1,"name":"Net 30","description":"Net 30","lines":[{"id":1,"sequence":1,"value_amount_percent":100,"number_of_days":30}]}}]}}`
		expectedXTotalCountHeader := "3"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FilterSuccessGetListOfOrders", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/orders?commitment_date:eq=2021-04-05", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"orders":[{"id":1,"name":"Order 1","commitment_date":"2021-04-05T00:00:00Z","note":"","quotation":{"id":4,"name":"Quotation 4","creation_date":"2021-04-01T00:00:00Z","validity_date":"2021-04-30T00:00:00Z","discount":200,"amount_delivery":400,"status":"sales_order","total_amount":2050,"customer":{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},"items":[{"id":5,"name":"Item 5","description":"Item 5 description","price":2000,"discount":150,"amount_total":1850}]},"payment_term":{"id":4,"name":"30% Now, Balance 60 Days","description":"Pay 30% now, balance due in 60 days","lines":[{"id":4,"sequence":1,"value_amount_percent":30,"number_of_days":0},{"id":5,"sequence":2,"value_amount_percent":70,"number_of_days":60}]}}]}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessGetOrderById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/orders/1", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"order":{"id":1,"name":"Order 1","commitment_date":"2021-04-05T00:00:00Z","note":"","quotation":{"id":4,"name":"Quotation 4","creation_date":"2021-04-01T00:00:00Z","validity_date":"2021-04-30T00:00:00Z","discount":200,"amount_delivery":400,"status":"sales_order","total_amount":2050,"customer":{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},"items":[{"id":5,"name":"Item 5","description":"Item 5 description","price":2000,"discount":150,"amount_total":1850}]},"payment_term":{"id":4,"name":"30% Now, Balance 60 Days","description":"Pay 30% now, balance due in 60 days","lines":[{"id":4,"sequence":1,"value_amount_percent":30,"number_of_days":0},{"id":5,"sequence":2,"value_amount_percent":70,"number_of_days":60}]}}}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("NotFoundGetOrderById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/sales/orders/0", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"order not found","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessCreateQuotation", func(t *testing.T) {
		w := httptest.NewRecorder()

		testTime, err := time.Parse(time.RFC3339, "2021-07-01T00:00:00Z")
		assert.NoError(t, err)

		request := sales.SalesQuotationCreateRequest{
			Name:           "Quotation 7",
			CustomerId:     500,
			CreationDate:   testTime,
			ValidityDate:   testTime.AddDate(0, 0, 30),
			Discount:       0,
			AmountDelivery: 0,
			Status:         models.SalesQuotationStatusQuotation,
			SalesOrderItems: []sales.SalesOrderItemCreateRequest{
				{
					Name:        "Item 7",
					Description: "Item 7 description",
					Price:       4000,
					Discount:    0,
				},
			},
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/sales/quotations", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":201,"message":"quotation created successfully","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FailedCreateQuotation", func(t *testing.T) {
		w := httptest.NewRecorder()

		testTime, err := time.Parse(time.RFC3339, "2021-07-01T00:00:00Z")
		assert.NoError(t, err)

		request := sales.SalesQuotationCreateRequest{
			Name:           "Quotation 1",
			CustomerId:     500,
			CreationDate:   testTime,
			ValidityDate:   testTime.AddDate(0, 0, 30),
			Discount:       0,
			AmountDelivery: 0,
			Status:         models.SalesQuotationStatusQuotation,
			SalesOrderItems: []sales.SalesOrderItemCreateRequest{
				{
					Name:        "Item 7",
					Description: "Item 7 description",
					Price:       4000,
					Discount:    0,
				},
			},
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/sales/quotations", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":409,"message":"quotation name already exists","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})
}
