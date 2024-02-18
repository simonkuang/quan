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
			// if data is not dirty, skip this one
			if !config.BackupDirtyFlag {
				continue
			}
			// if another backup is running, skip this one
			if config.BackupRunningFlag {
				continue
			}
			config.BackupRunningFlag = true
			flag := backup.Backup(database.GetLevelDB(), config.BackupFileName, config.BackupSize)
			config.BackupRunningFlag = false
			if flag {
				config.BackupDirtyFlag = false
			}
		}
	}()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(fmt.Sprintf(":%d", config.Port))
}
