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

func TestSettingRoutes(t *testing.T) {
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
	routes.Setting(router, DB)

	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       "admin@buon18.com",
		Role:        "bot",
		Permissions: []string{"FULL_ACCESS"},
	})
	assert.NoError(t, err)

	t.Run("SuccessGetListOfUsers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"users":[{"id":1,"name":"bot","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}},{"id":2,"name":"admin","email":"admin@buon18.com","type":"user","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}}]}}`, w.Body.String())
	})

	t.Run("FilterSuccessGetListOfUsers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users?name-like=bot", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"users":[{"id":1,"name":"bot","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}}]}}`, w.Body.String())
	})

	t.Run("SortSuccessGetListOfUsers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users?sort-email=asc", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"users":[{"id":2,"name":"admin","email":"admin@buon18.com","type":"user","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}},{"id":1,"name":"bot","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}}]}}`, w.Body.String())
	})

	t.Run("LimitAndOffsetSuccessGetListOfUsers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users?offset=1&limit=1", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"users":[{"id":2,"name":"admin","email":"admin@buon18.com","type":"user","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}}]}}`, w.Body.String())
	})

	t.Run("SuccessGetUserById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users/1", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"user":{"id":1,"name":"bot","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}}}}`, w.Body.String())
	})

	t.Run("NotFoundGetUserById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users/999", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":404,"message":"user not found","data":null}`, w.Body.String())
	})

	t.Run("SucessGetListOfCustomers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/customers", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"customers":[{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},{"id":502,"full_name":"John Foo","gender":"u","email":"jf@dummy-data.com","phone":"012789123","additional_information":{"note":"This is a dummy data from john foo"}}]}}`, w.Body.String())
	})

	t.Run("FilterSucessGetListOfCustomers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/customers?fullname-like=Jane", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"customers":[{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}}]}}`, w.Body.String())
	})

	t.Run("NotFoundGetCustomerById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/customers/500", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}}}}`, w.Body.String())
	})

	t.Run("SuccessGetCustomerById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/customers/999", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":404,"message":"customer not found","data":null}`, w.Body.String())
	})

	t.Run("SuccessGetListOfRoles", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/roles", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"roles":[{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]},{"id":2,"name":"user","description":"User","Permissions":[{"id":6,"name":"VIEW_PROFILE"},{"id":7,"name":"UPDATE_PROFILE"}]}]}}`, w.Body.String())
	})

	t.Run("FilterSuccessGetListOfRoles", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/roles?name-like=user", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"roles":[{"id":2,"name":"user","description":"User","Permissions":[{"id":6,"name":"VIEW_PROFILE"},{"id":7,"name":"UPDATE_PROFILE"}]}]}}`, w.Body.String())
	})

	t.Run("SortSuccessGetListOfRoles", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/roles?sort-name=desc", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"roles":[{"id":2,"name":"user","description":"User","Permissions":[{"id":6,"name":"VIEW_PROFILE"},{"id":7,"name":"UPDATE_PROFILE"}]},{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}]}}`, w.Body.String())
	})
}
