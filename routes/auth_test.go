package routes_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"system.buon18.com/m/config"
	"system.buon18.com/m/database"
	"system.buon18.com/m/routes"
	"system.buon18.com/m/services"
	"system.buon18.com/m/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestAuthRoutes(t *testing.T) {
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
		),
		postgres.BasicWaitStrategies(),
	)
	assert.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	DB := database.InitSQL(connectionString)

	router := gin.Default()
	routes.Auth(router, DB)

	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       "admin@buon18.com",
		Role:        "bot",
		Permissions: []string{"FULL_ACCESS"},
	})
	assert.NoError(t, err)

	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: "admin@buon18.com",
	})
	assert.NoError(t, err)

	t.Run("SuccessGetProfile", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/auth/me", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("AddPasswordToUserWithoutPassword", func(t *testing.T) {
		w := httptest.NewRecorder()

		updatePassword := services.UpdatePasswordRequest{
			OldPassword: "no-password", // -- No password test
			NewPassword: "new-password",
		}

		updatePasswordJson, err := json.Marshal(updatePassword)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/auth/update-password", strings.NewReader(string(updatePasswordJson)))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("WrongOldPasswordUpdatePassword", func(t *testing.T) {
		w := httptest.NewRecorder()

		updatePassword := services.UpdatePasswordRequest{
			OldPassword: "wrong-password",
			NewPassword: "new-password",
		}

		updatePasswordJson, err := json.Marshal(updatePassword)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/auth/update-password", strings.NewReader(string(updatePasswordJson)))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("WrongPasswordLogin", func(t *testing.T) {
		w := httptest.NewRecorder()

		login := services.LoginRequest{
			Email:    "admin@buon18.com",
			Password: "wrong-password",
		}

		loginJson, err := json.Marshal(login)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(string(loginJson)))
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("CorrectPasswordLogin", func(t *testing.T) {
		w := httptest.NewRecorder()

		login := services.LoginRequest{
			Email:    "admin@buon18.com",
			Password: "new-password",
		}

		loginJson, err := json.Marshal(login)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(string(loginJson)))
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("SuccessRefreshToken", func(t *testing.T) {
		w := httptest.NewRecorder()

		refreshTokenRequest := services.RefreshTokenRequest{
			RefreshToken: refreshToken,
		}
		refreshTokenJson, err := json.Marshal(refreshTokenRequest)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/auth/refresh-token", strings.NewReader(string(refreshTokenJson)))
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("SuccessUpdateProfile", func(t *testing.T) {
		w := httptest.NewRecorder()

		updateName := "Admin"
		updateProfile := &services.UpdateProfileRequest{
			Name: &updateName,
		}

		updateProfileJson, err := json.Marshal(*updateProfile)
		assert.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/api/auth/me", strings.NewReader(string(updateProfileJson)))
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		// Check if name is updated
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/api/auth/me", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		expectedBodyJSON := `{"code":200,"message":"","data":{"user":{"id":2,"name":"Admin","email":"admin@buon18.com","type":"user","role":{"id":1,"name":"bot","description":"BOT","permissions":[{"id":1,"name":"FULL_ACCESS"}]}}}}`

		assert.JSONEq(t, expectedBodyJSON, w.Body.String())
	})

}
