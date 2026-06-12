//go:build integration

// Package services_test exercises the admin auth flow against a
// throwaway Postgres (dockertest). Run with:
//
//   go test -tags=integration ./src/services/...
//
// Requires Docker. The image used is postgres:16-alpine; the test is
// skipped automatically if Docker isn't reachable.
package services_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/akshit_tyagi/postgresql_project/src/config"
	"github.com/akshit_tyagi/postgresql_project/src/constants"
	dbmig "github.com/akshit_tyagi/postgresql_project/src/database/migrations"
	authmodel "github.com/akshit_tyagi/postgresql_project/src/models"
	"github.com/akshit_tyagi/postgresql_project/src/services"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestAdminAuthFlow(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Skipf("docker unavailable, skipping integration test: %v", err)
	}
	pool.MaxWait = 2 * time.Minute

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=postgres",
		},
	}, func(c *docker.HostConfig) {
		c.AutoRemove = true
		c.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Skipf("could not start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = pool.Purge(resource) })

	hostPort := resource.GetHostPort("5432/tcp")
	host, port, splitErr := splitHostPortLocal(hostPort)
	if splitErr != nil {
		t.Fatalf("invalid hostport %q: %v", hostPort, splitErr)
	}
	t.Setenv("DB_HOST", host)
	t.Setenv("DB_PORT", port)
	t.Setenv("DB_USERNAME", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("DB_DATABASE", "postgres")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("APP_ENV", "test")
	t.Setenv("JWT_REFRESH_EXPIRY_DAYS", strconv.Itoa(30))
	t.Setenv("JWT_ACCESS_EXPIRY_MINUTES", strconv.Itoa(60))
	// Real secrets that satisfy the >= 32 byte length rule.
	t.Setenv("JWT_ACCESS_SECRET", "0123456789abcdef0123456789abcdef0123456789abcdef")
	t.Setenv("JWT_REFRESH_SECRET", "fedcba9876543210fedcba9876543210fedcba9876543210")
	t.Setenv("SESSION_SECRET", "0123456789abcdef0123456789abcdef0123456789abcdef")

	if err := config.InitializeDatabase(); err != nil {
		t.Fatalf("db init: %v", err)
	}
	t.Cleanup(func() { _ = config.Close() })

	// Apply schema by running the same migrations cmd/migrate runs.
	if err := dbmig.Run("up"); err != nil {
		t.Fatalf("migrations: %v", err)
	}

	// Seed a super admin.
	admin := &authmodel.User{
		Name:      "Test Admin",
		Email:     "test-admin@example.com",
		Role:      constants.SUPER_ADMIN_ROLE,
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := admin.HashPassword("correct-horse-battery-staple"); err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if err := config.DB.Create(admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}

	// 1. Wrong password is rejected.
	if _, err := services.Login(authmodel.AdminLoginRequest{
		Email:    "test-admin@example.com",
		Password: "wrong",
	}); err == nil {
		t.Fatal("expected error on wrong password, got nil")
	}

	// 2. Correct password returns tokens.
	resp, err := services.Login(authmodel.AdminLoginRequest{
		Email:    "test-admin@example.com",
		Password: "correct-horse-battery-staple",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("login returned empty tokens")
	}

	// 3. Refresh with the issued token issues a new access token.
	refreshed, err := services.RefreshToken(resp.RefreshToken)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if refreshed.AccessToken == "" {
		t.Fatal("refresh returned empty access token")
	}

	// 4. Logout revokes the refresh token.
	if err := services.Logout(resp.RefreshToken); err != nil {
		t.Fatalf("logout: %v", err)
	}

	// 5. A revoked token can no longer be used to refresh.
	if _, err := services.RefreshToken(resp.RefreshToken); err == nil {
		t.Fatal("expected error refreshing revoked token, got nil")
	}

	// 6. Logging out a second time reports the session as already revoked.
	if err := services.Logout(resp.RefreshToken); err == nil {
		t.Fatal("expected error logging out twice")
	}
}

// splitHostPortLocal splits a "host:port" string. dockertest used to
// return this as two values; the latest version returns a single
// string, so we parse it ourselves to stay portable.
func splitHostPortLocal(s string) (string, string, error) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' {
			return s[:i], s[i+1:], nil
		}
	}
	return "", "", errBadHostPort
}

var errBadHostPort = &hostPortError{}

type hostPortError struct{}

func (*hostPortError) Error() string { return "missing port in address" }
