package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simonkuang/quan/src/backup"
	"github.com/simonkuang/quan/src/config"
	"github.com/simonkuang/quan/src/controllers"
)

func main() {
	log.SetOutput(os.Stdout)

	r := gin.Default()
	config.FlagsInit(r)

	database := backup.Restore(config.BackupFileName, config.DBVersionStep)

	r.Use(func(ctx *gin.Context) {
		ctx.Set("database", database)
		ctx.Next()
	})

	controllers.RegisterControllers(r)

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(config.BackupInterval))
		for range ticker.C {
			backup.Backup(database.GetLevelDB(), config.BackupFileName, config.BackupSize)
		}
	}()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(fmt.Sprintf(":%d", config.Port))
}
