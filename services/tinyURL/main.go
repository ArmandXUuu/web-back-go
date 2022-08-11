package main

import (
	"fmt"
	"net/http"
	"time"
	"web-back-go/pkg/consts"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Hello world")
	router := gin.New()
	router.GET("/health", health)

	log.Fatal(router.Run(":9234"))
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"pi":        consts.CONST_TEST,
		"timestamp": time.Now().Unix(),
	})
}
