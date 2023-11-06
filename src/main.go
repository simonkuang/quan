package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/simonkuang/quan/src/config"
	"github.com/simonkuang/quan/src/controllers"
	"github.com/simonkuang/quan/src/db"
)

func main() {
	r := gin.Default()

	config.FlagsInit(r)

	database := &db.DBModel{}
	database.Connect()
	// FIXME: set database to context every single request
	r.Use(func(ctx *gin.Context) {
		ctx.Set("database", database)
		ctx.Next()
	})
	controllers.RegisterControllers(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(fmt.Sprintf(":%d", config.Port))
}
