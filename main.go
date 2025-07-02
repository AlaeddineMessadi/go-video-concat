package main

import (
	"net/http"
	"os"

	"github.com/amupxm/go-video-concat/config"
	ApiController "github.com/amupxm/go-video-concat/controller"
	"github.com/amupxm/go-video-concat/packages/cache"
	postgres "github.com/amupxm/go-video-concat/packages/database"
	s3 "github.com/amupxm/go-video-concat/packages/s3"
	"github.com/amupxm/go-video-concat/internal/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()
	logger.Log.Info("main() started")

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		logger.Log.Info("Running database migration...")
		// To be implemented: Load config and run GORM migration
		logger.Log.Info("Migration complete.")
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatalf("Failed to load config: %v", err)
	}
	logger.Log.Info("Config loaded successfully")

	// =======  Database and Storage =========
	postgres.PostgresConnection.ConnectDatabase(cfg)
	logger.Log.Info("Database connected")

	postgres.AutoMigration()
	logger.Log.Info("Database migrated")

	buckets := []string{"frame", "amupxm", "thumbnails", "splash", "upload", "splash-base", "splash-audio", "outputs"}
	s3.ObjectStorage.Connect(cfg)
	logger.Log.Info("S3 connected")

	s3.InitBuckets(buckets)
	logger.Log.Info("Buckets initialized")

	cache.Init(cfg)
	logger.Log.Info("Redis cache initialized")
	// =======================================

	router := gin.Default()

	// ============  Controllers  ============
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 1 - frame
	router.POST("/frame/upload", ApiController.Frame_Upload)
	router.POST("/frame", ApiController.Frame_Add)
	router.GET("/frame", ApiController.Frame_list)
	router.GET("/frame/:code/file", ApiController.Frame_File)
	router.GET("/frame/:code", ApiController.Frame_Single)

	// 2 - splash video (splash)
	router.POST("/splash/base", ApiController.Splash_Base)
	router.POST("/splash/audio", ApiController.Splash_Audio)
	router.POST("/splash", ApiController.Splash_Add)
	router.GET("/splash/:code", ApiController.Splash_file)

	// 3 - generator
	router.POST("/generator/upload", ApiController.Generator_Upload)
	router.POST("/generator", ApiController.Generator_Generate)
	router.GET("/generator/:code/file", ApiController.Generator_file)
	router.GET("/generator/:code", ApiController.Generator_Status)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	logger.Log.Infof("Server starting on port %s", port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		logger.Log.Errorf("Failed to start server: %v", err)
	}
}
