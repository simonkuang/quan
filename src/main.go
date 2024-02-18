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
	// set log output
	log.SetOutput(os.Stdout)

	r := gin.Default()

	config.FlagsInit(r)

	//var database *db.DBModel
	// database = &db.DBModel{}
	// database.Connect(db.GetDbFileName())

	// load data from backup bucket
	database := backup.Restore(config.BackupFileName, config.DBVersionStep)

	// FIXME: set database to context every single request
	r.Use(func(ctx *gin.Context) {
		ctx.Set("database", database)
		ctx.Next()
	})
	controllers.RegisterControllers(r)

	// tick-tock task
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(config.BackupScanDirtyInterval))
		for {
			select {
			case <-ticker.C:
				// if data is not dirty, skip this one
				if !config.BackupDirtyFlag {
					continue
				}
				// if another backup is running, skip this one
				if config.BackupRunningFlag {
					continue
				}
				// backup operation
				config.BackupRunningFlag = true
				flag := backup.Backup(database.GetLevelDB(), config.BackupFileName, config.BackupSize)
				config.BackupRunningFlag = false
				if flag {
					config.BackupDirtyFlag = false
				}
			}
		}
	}()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(fmt.Sprintf(":%d", config.Port))
}
