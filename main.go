package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amupxm/go-video-concat/config"
	ApiController "github.com/amupxm/go-video-concat/controller"
	"github.com/amupxm/go-video-concat/packages/cache"
	postgres "github.comcom/amupxm/go-video-concat/packages/database"
	s3 "github.com/amupxm/go-video-concat/packages/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		fmt.Println("Running database migration...")
		// To be implemented: Load config and run GORM migration
		fmt.Println("Migration complete.")
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// =======  Database and Storage =========
	postgres.PostgresConnection.ConnectDatabase(cfg)
	buckets := []string{"frame", "amupxm", "thumbnails", "splash", "upload", "splash-base", "splash-audio", "outputs"}
	s3.ObjectStorage.Connect(cfg)
	s3.InitBuckets(buckets)
	cache.Init(cfg)
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
	router.Run("0.0.0.0:" + port)
}
