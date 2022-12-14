package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"web-back-go/pkg/consts"
	tinyURL "web-back-go/pkg/tinyURL"

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
	tlsOption := "&tls=true"
	if strings.HasPrefix(app.config.DBHostAddr, "localhost") || strings.HasPrefix(app.config.DBHostAddr, "mysql") {
		tlsOption = ""
	}
	app.db = sqlx.MustConnect("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/ziyixu-tech?parseTime=true%s",
			app.config.DBUsername, app.config.DBPassword, app.config.DBHostAddr, tlsOption))

	// GIN Router
	app.router = gin.Default()
	api := app.router.Group("/api")
	{
		api.GET("/healthz", health)
		api.GET("/testMD5", app.testMD5)

		api.GET("/tinyURL/short", app.getShortCode)
		api.GET("/tinyURL/base", app.getBaseURL)
		api.GET("/tinyURL/list", app.listAllTinyURL)

		api.GET("/todo", app.getAllTodo)
		api.GET("/todo/:id", app.getTodo)
		api.POST("/todo", app.createTodo)
		api.POST("/todo/:id", app.toggleTodo)
		api.PATCH("/todo/:id", app.updateTodo)
		api.DELETE("/todo/:id", app.deleteTodo)
	}

	app.router.GET("/health", health)
}

func (app *App) testMD5(c *gin.Context) {
	inputString := c.DefaultQuery("input", "test_string")
	log.Debugf("input string is %s", inputString)

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"MD5":    tinyURL.GeneragteMD5Value(inputString),
	})
}
