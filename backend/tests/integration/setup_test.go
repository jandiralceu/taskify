//go:build integration

package integration

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/jandiralceu/taskify/internal/config"
	"github.com/jandiralceu/taskify/internal/database"
	"github.com/jandiralceu/taskify/internal/handlers"
	"github.com/jandiralceu/taskify/internal/pkg"
	"github.com/jandiralceu/taskify/internal/repository"
	"github.com/jandiralceu/taskify/internal/routes"
	"github.com/jandiralceu/taskify/internal/service"
	"github.com/jandiralceu/taskify/tests/integration/testhelpers"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// generateIntegrationRSAKeys creates a temporary RSA key pair for JWT signing during tests.
func generateIntegrationRSAKeys() (privateKeyPEM string, publicKeyPEM string, err error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", err
	}
	privBlock := &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}
	privPEM := pem.EncodeToMemory(privBlock)

	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}
	pubPEM := pem.EncodeToMemory(pubBlock)

	return string(privPEM), string(pubPEM), nil
}

// setupApp boots the entire application stack using real Postgres + Redis containers.
func setupApp(t *testing.T) (*httptest.Server, *gorm.DB, func()) {
	return setupAppCustom(t, nil)
}

func setupAppCustom(t *testing.T, modifyConfig func(*config.Config)) (*httptest.Server, *gorm.DB, func()) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Resolve migrations path relative to the tests/integration dir.
	migrationsDir := "../../migrations"

	// Start Postgres container with migrations.
	pgContainer, err := testhelpers.SetupPostgres(ctx, migrationsDir)
	require.NoError(t, err, "Failed to start Postgres container")

	// Start Redis container.
	redisContainer, err := testhelpers.SetupRedis(ctx)
	require.NoError(t, err, "Failed to start Redis container")

	// Build config from container info.
	cfg := &config.Config{
		DBHost:          pgContainer.Host,
		DBPort:          pgContainer.Port,
		DBUser:          pgContainer.User,
		DBPassword:      pgContainer.Password,
		DBName:          pgContainer.DBName,
		AppPort:         "0",
		Env:             "test",
		AppName:         "inventory-api-test",
		RedisHost:       redisContainer.Host,
		RedisPort:       redisContainer.Port,
		OTLPEndpoint:    "localhost:4317",
		RateLimitGlobal: "1000-M",
		RateLimitAuth:   "1000-M",
	}

	if modifyConfig != nil {
		modifyConfig(cfg)
	}

	// Initialize database (run migrations if needed, but testhelpers already does it)
	db, err := database.Init(ctx, cfg)
	require.NoError(t, err, "Failed to initialize database")

	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Cache.
	cacheManager := pkg.NewRedisCacheManager(cfg)

	// Generate RSA keys for JWT.
	privPEM, pubPEM, err := generateIntegrationRSAKeys()
	require.NoError(t, err, "Failed to generate RSA keys")

	jwtManager, err := pkg.NewJWTManager(privPEM, pubPEM)
	require.NoError(t, err, "Failed to initialize JWT manager")

	hasher := pkg.NewArgon2PasswordHasher()

	// Initialize Casbin Enforcer with GORM adapter for persistent RBAC during tests.
	adapter, err := gormadapter.NewAdapterByDB(db)
	require.NoError(t, err, "Failed to initialize Casbin adapter")

	enforcer, err := casbin.NewEnforcer("../../model.conf", adapter)
	require.NoError(t, err, "Failed to initialize Casbin enforcer")

	// Load policies from the database (seeded by migrations in testhelpers.SetupPostgres)
	require.NoError(t, enforcer.LoadPolicy(), "Failed to load Casbin policies")

	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	userService := service.NewUserService(userRepository, roleRepository, hasher, cfg.UploadPath, cacheManager)
	authHandler := handlers.NewAuthHandler(userService, jwtManager, cacheManager, hasher)
	userHandler := handlers.NewUserHandler(userService, enforcer)

	routeConfig := &routes.RouteConfig{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}

	// Suppress Gin debug output during tests.
	gin.SetMode(gin.TestMode)
	r := routes.Setup(routeConfig, cfg, jwtManager, enforcer, cacheManager)

	// Health check (same as main.go).
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	ts := httptest.NewServer(r)

	cleanup := func() {
		ts.Close()
		_ = cacheManager.Close()
		_ = sqlDB.Close()
		// Termination might take time, use a separate context
		termCtx := context.Background()
		_ = pgContainer.Terminate(termCtx)
		_ = redisContainer.Terminate(termCtx)
	}

	return ts, db, cleanup
}

// signUpUser registers a user using the API.
func signUpUser(t *testing.T, baseURL, firstName, lastName, email, password, role string) {
	t.Helper()

	body := fmt.Sprintf(`{"firstName":"%s","lastName":"%s","email":"%s","password":"%s","role":"%s"}`, firstName, lastName, email, password, role)
	resp, err := http.Post(baseURL+"/api/v1/auth/register", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	require.Equal(t, http.StatusNoContent, resp.StatusCode, "Register should return 204")
}

// signInUser authenticates a user and returns markers.
func signInUser(t *testing.T, baseURL, email, password string) (string, string) {
	t.Helper()

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	resp, err := http.Post(baseURL+"/api/v1/auth/signin", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		t.Fatalf("SignIn failed with status %d: %v", resp.StatusCode, errResp)
	}

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	accessToken, ok1 := result["accessToken"].(string)
	refreshToken, ok2 := result["refreshToken"].(string)
	require.True(t, ok1 && ok2, "Response must contain tokens")

	return accessToken, refreshToken
}

// authedRequest sends an HTTP request with an Authorization Bearer token.
func authedRequest(t *testing.T, method, url, token string, bodyData any) *http.Response {
	t.Helper()

	var reqBody *strings.Reader
	if bodyData != nil {
		jsonBytes, _ := json.Marshal(bodyData)
		reqBody = strings.NewReader(string(jsonBytes))
	} else {
		reqBody = strings.NewReader("")
	}

	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

