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
		postgres.WithInitScripts(filepath.Join("..", "database", "dev_scripts", "create-schema.sh"), filepath.Join("..", "database", "dev_scripts", "seed.sh")),
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

	t.Run("GET /api/setting/users", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, `{"code":200,"message":"","data":{"users":[{"id":1,"name":"bot","email":"bot@buon18.com","type":"bot","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}},{"id":2,"name":"admin","email":"admin@buon18.com","type":"user","role":{"id":1,"name":"bot","description":"BOT","Permissions":[{"id":1,"name":"FULL_ACCESS"}]}}]}}`, w.Body.String())
	})

	t.Run("GET /api/setting/users?name-like=bot", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/setting/users?name-like=bot", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		t.Logf("Response: %s", w.Body.String())

		assert.Equal(t, 200, w.Code)
	})
}
