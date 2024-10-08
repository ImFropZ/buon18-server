package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"system.buon18.com/m/config"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/routes"
	"system.buon18.com/m/utils"

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
			filepath.Join("..", "database", "dev_scripts", "001_create-schema.sh"),
			filepath.Join("..", "database", "dev_scripts", "002_seed.sh"),
			filepath.Join("..", "database", "dev_scripts", "100_seed-customer.sh"),
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
	routes.Setting(router, &database.Connection{
		DB:     DB,
		Valkey: nil,
	})

	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       "admin@buon18.com",
		Role:        "bot",
		Permissions: []string{"FULL_ACCESS"},
	})
	assert.NoError(t, err)

	settingToken, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       "setting@buon18.com",
		Role:        "bot",
		Permissions: []string{"FULL_ACCESS"},
	})
	assert.NoError(t, err)

	t.Run("SuccessGetListOfUsers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"users":[{"id":1,"name":"bot","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","permissions":[{"id":1,"name":"FULL_ACCESS"}]}},{"id":2,"name":"admin","email":"admin@buon18.com","type":"user","role":{"id":1,"name":"bot","description":"BOT","permissions":[{"id":1,"name":"FULL_ACCESS"}]}},{"id":3,"name":"Setting Admin","email":"setting@buon18.com","type":"user","role":{"id":3,"name":"Setting Administrator","description":"Full access to all settings","permissions":[{"id":3,"name":"FULL_SETTING"}]}},{"id":4,"name":"Sales Admin","email":"sales@buon18.com","type":"user","role":{"id":4,"name":"Sales Administrator","description":"Full access to all sales","permissions":[{"id":4,"name":"FULL_SALES"}]}},{"id":5,"name":"Accounting Admin","email":"accounting@buon18.com","type":"user","role":{"id":5,"name":"Accounting Administrator","description":"Full access to all accounting","permissions":[{"id":5,"name":"FULL_ACCOUNTING"}]}}]}}`
		expectedXTotalCountHeader := "5"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FilterSuccessGetListOfUsers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users?name:like=bot", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"users":[{"id":1,"name":"bot","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","permissions":[{"id":1,"name":"FULL_ACCESS"}]}}]}}`
		expectedXTotalCountHeader := "1"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessGetUserById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users/1", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"user":{"id":1,"name":"bot","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","permissions":[{"id":1,"name":"FULL_ACCESS"}]}}}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("NotFoundGetUserById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users/0", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"user not found","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SucessGetListOfCustomers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/customers", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"customers":[{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}},{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}},{"id":502,"full_name":"John Foo","gender":"u","email":"jf@dummy-data.com","phone":"012789123","additional_information":{"note":"This is a dummy data from john foo"}}]}}`
		expectedXTotalCountHeader := "3"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FilterSucessGetListOfCustomers", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/customers?fullname:like=Jane", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"customers":[{"id":501,"full_name":"Jane Doe","gender":"f","email":"jad@dummy-data.com","phone":"064456789","additional_information":{"note":"This is a dummy data from jane doe"}}]}}`
		expectedXTotalCountHeader := "1"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessGetCustomerById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/customers/500", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"customer":{"id":500,"full_name":"John Doe","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}}}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("NotFoundGetCustomerById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/customers/0", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"customer not found","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessGetListOfRoles", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/roles", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"roles":[{"id":1,"name":"bot","description":"BOT","permissions":[{"id":1,"name":"FULL_ACCESS"}]},{"id":2,"name":"user","description":"User","permissions":[{"id":6,"name":"VIEW_PROFILE"},{"id":7,"name":"UPDATE_PROFILE"}]},{"id":3,"name":"Setting Administrator","description":"Full access to all settings","permissions":[{"id":3,"name":"FULL_SETTING"}]},{"id":4,"name":"Sales Administrator","description":"Full access to all sales","permissions":[{"id":4,"name":"FULL_SALES"}]},{"id":5,"name":"Accounting Administrator","description":"Full access to all accounting","permissions":[{"id":5,"name":"FULL_ACCOUNTING"}]}]}}`
		expectedXTotalCountHeader := "5"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FilterSuccessGetListOfRoles", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/roles?name:like=user", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"roles":[{"id":2,"name":"user","description":"User","permissions":[{"id":6,"name":"VIEW_PROFILE"},{"id":7,"name":"UPDATE_PROFILE"}]}]}}`
		expectedXTotalCountHeader := "1"

		assert.Equal(t, expectedXTotalCountHeader, w.Header().Get("X-Total-Count"))
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessGetRoleById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/roles/2", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"role":{"id":2,"name":"user","description":"User","permissions":[{"id":6,"name":"VIEW_PROFILE"},{"id":7,"name":"UPDATE_PROFILE"}]}}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("NotFoundGetRoleById", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/roles/0", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"role not found","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessCreateUser", func(t *testing.T) {
		w := httptest.NewRecorder()

		request := setting.SettingUserCreateRequest{
			Name:   "test",
			Email:  "test@buon18.com",
			RoleId: 2,
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/setting/users", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":201,"message":"user created successfully","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FailedCreateUser", func(t *testing.T) {
		w := httptest.NewRecorder()

		request := setting.SettingUserCreateRequest{
			Name:   "test",
			Email:  "admin@buon18.com",
			RoleId: 2,
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/setting/users", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":409,"message":"user email already exists","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessCreateCustomer", func(t *testing.T) {
		w := httptest.NewRecorder()

		request := setting.SettingCustomerCreateRequest{
			FullName:              "success test",
			Gender:                "m",
			Email:                 "success-test@buon18.com",
			Phone:                 "+85512123123",
			AdditionalInformation: `{"note":"This is a dummy data from test"}`,
		}

		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/setting/customers", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":201,"message":"customer created successfully","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FailedCreateCustomer", func(t *testing.T) {
		w := httptest.NewRecorder()

		request := setting.SettingCustomerCreateRequest{
			FullName:              "failed test",
			Gender:                "m",
			Email:                 "failed-test@buon18.com",
			Phone:                 "+85512123123",
			AdditionalInformation: `{}`,
		}

		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/setting/customers", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		w = httptest.NewRecorder()

		req = httptest.NewRequest("POST", "/api/setting/customers", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":409,"message":"customer email already exists","data":null}`
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessCreateRole", func(t *testing.T) {
		w := httptest.NewRecorder()

		request := setting.SettingRoleCreateRequest{
			Name:          "test",
			Description:   "test",
			PermissionIds: []uint{1},
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/setting/roles", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":201,"message":"role created successfully","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FailedCreateRole", func(t *testing.T) {
		w := httptest.NewRecorder()

		request := setting.SettingRoleCreateRequest{
			Name:          "test",
			Description:   "test",
			PermissionIds: []uint{1},
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/setting/roles", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+settingToken)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":403,"message":"unable to create role with full permission","data":null}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessUpdateCustomer", func(t *testing.T) {
		w := httptest.NewRecorder()

		fullName := "success"
		request := setting.SettingCustomerUpdateRequest{
			FullName: &fullName,
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/setting/customers/500", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"customer updated successfully","data":null}`
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())

		w = httptest.NewRecorder()

		req = httptest.NewRequest("GET", "/api/setting/customers/500", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON = `{"code":200,"message":"","data":{"customer":{"id":500,"full_name":"success","gender":"m","email":"jd@dummy-data.com","phone":"096123456","additional_information":{"note":"This is a dummy data from john doe"}}}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("FailedUpdateCustomer", func(t *testing.T) {
		w := httptest.NewRecorder()

		fullName := "success"
		request := setting.SettingCustomerUpdateRequest{
			FullName: &fullName,
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/setting/customers/0", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"customer not found","data":null}`
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessUpdateUser", func(t *testing.T) {
		w := httptest.NewRecorder()

		name := "success"
		request := setting.SettingUserUpdateRequest{
			Name: &name,
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/setting/users/1", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"user updated successfully","data":null}`
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())

		w = httptest.NewRecorder()

		req = httptest.NewRequest("GET", "/api/setting/users/1", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON = `{"code":200,"message":"","data":{"user":{"id":1,"name":"success","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","permissions":[{"id":1,"name":"FULL_ACCESS"}]}}}}`
		assert.Contains(t, w.Body.String(), expectedBodyJSON)
	})

	t.Run("FailedUpdateUser", func(t *testing.T) {
		w := httptest.NewRecorder()

		name := "success"
		request := setting.SettingUserUpdateRequest{
			Name: &name,
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/setting/users/0", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"user not found","data":null}`
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

	t.Run("SuccessUpdateRole", func(t *testing.T) {
		w := httptest.NewRecorder()

		name := "success"
		request := setting.SettingRoleUpdateRequest{
			Name: &name,
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/setting/roles/2", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"role updated successfully","data":null}`
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())

		w = httptest.NewRecorder()

		req = httptest.NewRequest("GET", "/api/setting/roles/2", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON = `{"code":200,"message":"","data":{"role":{"id":2,"name":"success","description":"User","permissions":[{"id":6,"name":"VIEW_PROFILE"},{"id":7,"name":"UPDATE_PROFILE"}]}}}`
		assert.JSONEq(t, w.Body.String(), expectedBodyJSON)
	})

	t.Run("FailedUpdateRole", func(t *testing.T) {
		w := httptest.NewRecorder()

		name := "success"
		request := setting.SettingRoleUpdateRequest{
			Name: &name,
		}
		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/setting/roles/0", bytes.NewReader(jsonData))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":404,"message":"role not found","data":null}`
		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})
}
