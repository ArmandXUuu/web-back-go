package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"web-back-go/pkg/consts"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type App struct {
	router *gin.Engine
	db     *sqlx.DB
	config Config
}

type Config struct {
	DBUsername string `envconfig:"DB_USERNAME" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DBHostAddr string `envconfig:"DB_HOSTADDR" required:"true"`
}

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Couldnot get env var : %v", err)
	}

	var app App

	app.init()

	log.Fatal(app.router.Run(":9234"))
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "OK",
		"pi":          consts.CONST_TEST,
		"timestamp":   time.Now().Unix(),
		"currentPath": c.Request.URL.Path,
	})
}

func (app *App) init() {
	// CONFIG
	var err error
	var config Config

	err = envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	app.config = config

	configb, err := json.Marshal(app.config)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(configb))

	// DATABASE
	app.db = sqlx.MustConnect("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/ziyixu-tech?parseTime=true&tls=true",
			app.config.DBUsername, app.config.DBPassword, app.config.DBHostAddr))

	// GIN Router
	app.router = gin.Default()
	api := app.router.Group("/api", health)
	{
		api.GET("/healthz", health)
	}

	app.router.GET("/health", health)
}
