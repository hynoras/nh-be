package main

// @title NoHeir API
// @version 1.0
// @description API server for NoHeir application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@noheir.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey SessionAuth
// @in cookie
// @name auth_session

import (
	"log/slog"
	"os"

	"nh-be/internal/app"
	"nh-be/internal/config"
	"nh-be/internal/features/auth"
	"nh-be/internal/features/experiment"
	"nh-be/internal/features/experiment/result"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/internal/router"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)


var migrateSchemas = []any{
	&auth.VerificationToken{},
	&user.User{},
	&permission.Permission{},
	&permission.PermissionGroup{},
	&user.UserPermission{},
	&experiment.Experiment{},
	&result.ExperimentResult{},
}

func main() {
	app.SetupLogging()

	if err := godotenv.Load(); err != nil {
		slog.Warn("Warning: No .env file found")
	}

	cfg := config.LoadConfig()

	slog.Info("starting app", "env", cfg.AppEnv)

	service, err := app.InitializeServices(cfg)
	if err != nil {
		slog.Error("failed to initialize services", "error", err)
		os.Exit(1)
	}

	autoMigrate(cfg, service.DB)

	tracerShutdown := app.SetupObservability(service)
	defer service.ShutdownServices(tracerShutdown)

	r := app.NewRouter(cfg)

	deps := service.NewSharedDeps(cfg)
	router.SetupRoutes(r, deps)

	app.ListenAndServe(cfg, r, service)
}

// autoMigrate runs GORM AutoMigrate for all domain models when the
// application is running in development mode (AppEnv == "dev"). This is
// a no-op in production. It terminates the process if migration fails.
//
// This function lives in main rather than internal/app to avoid an import
// cycle: feature model packages (auth, user, etc.) import app.SharedDeps
// via their router.go files, so app cannot import them back.
func autoMigrate(cfg *config.Config, db *gorm.DB) {
	if cfg.AppEnv != "dev" {
		return
	}

	if err := db.AutoMigrate(migrateSchemas...); err != nil {
		slog.Error("failed to run GORM AutoMigrate", "error", err)
		os.Exit(1)
	}
	slog.Info("Running AutoMigrate in dev mode")
}
