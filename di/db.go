package di

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func InitDB() *gorm.DB {
	type Config struct {
		Driver string `json:"driver"`
		DSN    string `yaml:"dsn"`
	}

	var c Config
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(err)
	}

	var db *gorm.DB
	if c.Driver == "mysql" {
		db, err = gorm.Open(mysql.Open(c.DSN), &gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\n", log.LstdFlags), logger.Config{LogLevel: logger.Warn, Colorful: true, ParameterizedQueries: false}),
		})
	}

	if err != nil {
		panic(err)
	}

	return db
}
