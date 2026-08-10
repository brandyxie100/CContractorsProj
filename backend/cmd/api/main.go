package main

import (
	"log"

	"github.com/clementscontractors/equipment/internal/config"
	"github.com/clementscontractors/equipment/internal/handler"
	"github.com/clementscontractors/equipment/internal/middleware"
	"github.com/clementscontractors/equipment/internal/models"
	"github.com/clementscontractors/equipment/internal/repository"
	"github.com/clementscontractors/equipment/internal/seed"
	"github.com/clementscontractors/equipment/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Location{},
		&models.Asset{},
		&models.Project{},
		&models.Assignment{},
		&models.MaintenanceRecord{},
		&models.RefreshToken{},
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	if err := seed.Run(db, cfg); err != nil {
		log.Fatalf("seed: %v", err)
	}

	repos := repository.New(db)
	svc := service.New(repos, cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	api := &handler.API{S: svc, R: repos}

	r := gin.Default()
	r.Use(middleware.CORS(cfg.CORSOrigin))
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api.Register(r)

	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
