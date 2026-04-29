package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/infra/db"
	"backend/internal/infra/db/dao"
	"backend/internal/infra/web"
	"backend/internal/infra/web/gin"
	"backend/pkg"
	"backend/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	appEnv := pkg.GetEnvDefault("APP_ENV", "development")
	envFile := ".env"
	if appEnv == "development" {
		envFile = ".env.development"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		logger.Warn("Error loading environment file: " + envFile)
	}

	db, err := db.NewDatabaseSQLFactory(db.InstancePostgres)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := db.AutoMigrate(dao.NewDAOs()...); err != nil {
		logger.Fatal("Failed to migrate database: " + err.Error())
	}

	config := web.NewConfigWeb()
	server, err := gin.NewGinServer(config.Host, config.Port, config.CorsAllowOrigins, db)
	if err != nil {
		logger.Fatal(err.Error())
	}
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server Shutdown: %s", err.Error()))
	}
	<-ctx.Done()
}
