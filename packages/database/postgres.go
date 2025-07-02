package postgres

import (
	"fmt"
	"github.com/amupxm/go-video-concat/internal/logger"

	"github.com/amupxm/go-video-concat/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DBCli *gorm.DB
}

var PostgresConnection = Postgres{}

func (d *Postgres) ConnectDatabase(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal(err)
	}

	d.DBCli = db
}
