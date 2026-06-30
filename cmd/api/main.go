package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/controllers/health"
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	"github.com/akshit_tyagi/postgresql_project/internal/routes"
	"github.com/danielkov/gin-helmet/ginhelmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/akshit_tyagi/postgresql_project/docs"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found; relying on process env")
	}

	appConfig, err := config.LoadApp()
	if err != nil {
		slog.Error("config validation failed", "error", err)
		os.Exit(1)
	}
	if appConfig.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := slogLogger(appConfig.AppEnv)
	slog.SetDefault(logger)

	if err := config.InitializeDatabase(); err != nil {
		slog.Error("database init failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := config.Close(); err != nil {
			slog.Error("database close failed", "error", err)
		}
	}()

	slog.Info("startup",
		"app", appConfig.AppName,
		"env", appConfig.AppEnv,
		"port", appConfig.AppPort,
	)

	setupGinLogger()
	r := gin.New()
	r.Use(ginLoggerWithRequestID())
	r.Use(gin.CustomRecovery(recoveryHandler))
	r.Use(middlewares.RateLimiter())
	r.HandleMethodNotAllowed = true
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":     false,
			"statusCode": http.StatusNotFound,
			"message":    "Route not found",
		})
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":     false,
			"statusCode": http.StatusMethodNotAllowed,
			"message":    "Method not allowed",
		})
	})
	r.Use(ginhelmet.Default())
	r.Use(hostAllowlist(appConfig.AllowedHosts))
	r.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{appConfig.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	if appConfig.SessionSecret == "" {
		slog.Error("SESSION_SECRET is required to use the sessions middleware")
		os.Exit(1)
	}
	store := cookie.NewStore([]byte(appConfig.SessionSecret))
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "1.0.0"})
	})
	r.GET("/healthz", health.Healthz)
	r.GET("/readyz", health.Readyz)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	v1 := r.Group("/api/v1")
	{
		adminGroup := v1.Group("/admin")
		routes.AdminRoutes(adminGroup)
	}

	srv := &http.Server{
		Addr:         ":" + appConfig.AppPort,
		Handler:      r.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
	slog.Info("server exited")
}

func slogLogger(env string) *slog.Logger {
	level := slog.LevelInfo
	if strings.EqualFold(env, "local") {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{Level: level}
	if strings.EqualFold(env, "production") {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func setupGinLogger() {
	logPath := os.Getenv("GIN_LOG_PATH")
	if logPath == "" {
		logPath = "gin.log"
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Warn("gin log file open failed; using stdout", "error", err)
		gin.DefaultWriter = os.Stdout
		return
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(f, os.Stdout)
}

func ginLoggerWithRequestID() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[req_id=%s] %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func recoveryHandler(c *gin.Context, recovered any) {
	slog.Error("panic recovered",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"panic", fmt.Sprint(recovered),
	)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status":     false,
		"statusCode": http.StatusInternalServerError,
		"message":    "Something went wrong. Please try again.",
	})
}

func hostAllowlist(allowed []string) gin.HandlerFunc {
	if len(allowed) == 0 {
		return func(c *gin.Context) { c.Next() }
	}
	set := make(map[string]struct{}, len(allowed))
	for _, h := range allowed {
		if normalized, err := normalizeHost(h); err == nil {
			set[normalized] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		reqHost, err := normalizeHost(c.Request.Host)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		if _, ok := set[reqHost]; !ok {
			slog.Warn("rejected request from disallowed host",

				"host", c.Request.Host,
			)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Next()
	}
}
func normalizeHost(raw string) (string, error) {
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	parsed, err := url.Parse("//" + raw)
	if err != nil {
		return "", err
	}
	host := parsed.Hostname()
	if host == "" {
		return strings.ToLower(raw), nil
	}
	port := parsed.Port()
	if port != "" {
		return strings.ToLower(host + ":" + port), nil
	}
	return strings.ToLower(host), nil
}
