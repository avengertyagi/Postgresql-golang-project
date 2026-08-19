package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/bootstrap"
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/database/migrations"
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/health"
	"github.com/akshit_tyagi/postgresql_project/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found; relying on process env")
	}

	logWriter, closeLog := setupLogging()
	defer closeLog()

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
	if err := migrations.Migrate(); err != nil {
		slog.Error("database migration failed", "error", err)
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

	//setupGinLogger()
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		slog.Error("failed to set trusted proxies", "error", err)
		os.Exit(1)
	}
	r.Use(gin.LoggerWithWriter(logWriter))
	r.Use(requestid.New())
	r.Use(gin.CustomRecovery(recoveryHandler))
	r.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(appConfig.AllowedOrigin, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middlewares.RateLimiter())
	r.HandleMethodNotAllowed = true
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":     false,
			"statusCode": http.StatusNotFound,
			"message":    "Route not found",
			"data":       gin.H{},
		})
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":     false,
			"statusCode": http.StatusMethodNotAllowed,
			"message":    "Method not allowed",
			"data":       gin.H{},
		})
	})
	r.Use(secure.New(secure.Config{
		AllowedHosts:          strings.Split(appConfig.AllowedHosts, ","),
		SSLRedirect:           appConfig.IsProd(),
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	}))
	r.Use(middlewares.RequestSizeLimiter(10 * 1024 * 1024))
	r.Use(middlewares.TimeoutMiddleware(5 * time.Second))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "1.0.0", "message": "Welcome to the Go lang."})
	})
	r.GET("/healthz", health.Healthz)
	r.GET("/readyz", health.Readyz)

	v1 := r.Group("/api/v1")
	{
		container := bootstrap.NewContainer(config.DB)
		routes.Routes(v1, container)

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

func setupLogging() (io.Writer, func()) {
	logPath := os.Getenv("GIN_LOG_PATH")
	if logPath == "" {
		logPath = "gin.log"
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Warn("gin log file open failed; logging to stdout only", "error", err)
		return os.Stdout, func() {}
	}
	w := io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = w
	gin.DefaultErrorWriter = w
	return w, func() { _ = f.Close() }
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
