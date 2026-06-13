package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type App struct {
	AppName                string
	AppEnv                 string
	AppKey                 string
	AppDebug               bool
	AppURL                 string
	AppPort                string
	AllowedOrigin          string
	AllowedHosts           []string
	JWTAccessSecret        string
	JWTRefreshSecret       string
	JWTAccessExpiryMinutes int
	JWTRefreshExpiryDays   int
	SessionSecret          string
}

func LoadApp() (App, error) {
	debug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	accessExpiry, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	refreshExpiry, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY_DAYS"))

	app := App{
		AppName:                os.Getenv("APP_NAME"),
		AppEnv:                 os.Getenv("APP_ENV"),
		AppKey:                 os.Getenv("APP_KEY"),
		AppDebug:               debug,
		AppURL:                 os.Getenv("APP_URL"),
		AppPort:                os.Getenv("APP_PORT"),
		AllowedOrigin:          os.Getenv("ALLOWED_ORIGIN"),
		AllowedHosts:           splitCSV(os.Getenv("ALLOWED_HOSTS")),
		JWTAccessSecret:        os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret:       os.Getenv("JWT_REFRESH_SECRET"),
		JWTAccessExpiryMinutes: accessExpiry,
		JWTRefreshExpiryDays:   refreshExpiry,
		SessionSecret:          os.Getenv("SESSION_SECRET"),
	}

	if err := app.validate(); err != nil {
		return App{}, err
	}
	return app, nil
}

func (a App) IsProd() bool { return strings.EqualFold(a.AppEnv, "production") }

func (a App) validate() error {
	var missing []string
	need := func(name, val string) {
		if val == "" {
			missing = append(missing, name)
		}
	}
	need("APP_KEY", a.AppKey)
	need("JWT_ACCESS_SECRET", a.JWTAccessSecret)
	need("JWT_REFRESH_SECRET", a.JWTRefreshSecret)
	need("SESSION_SECRET", a.SessionSecret)
	need("ALLOWED_ORIGIN", a.AllowedOrigin)
	need("APP_PORT", a.AppPort)

	if len(missing) > 0 {
		if a.AppEnv == "" || strings.EqualFold(a.AppEnv, "local") {
		} else {
			return fmt.Errorf("config: missing required env vars: %s", strings.Join(missing, ", "))
		}
	}

	if len(a.SessionSecret) > 0 && len(a.SessionSecret) < 32 {
		return errors.New("config: SESSION_SECRET must be at least 32 bytes; run `go run ./cmd/keygen` to generate one")
	}
	if a.JWTAccessSecret != "" && len(a.JWTAccessSecret) < 32 {
		return errors.New("config: JWT_ACCESS_SECRET must be at least 32 bytes")
	}
	if a.JWTRefreshSecret != "" && len(a.JWTRefreshSecret) < 32 {
		return errors.New("config: JWT_REFRESH_SECRET must be at least 32 bytes")
	}
	return nil
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
