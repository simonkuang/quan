package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/simonkuang/quan/src/config"
	"github.com/simonkuang/quan/src/db"
)

func RegisterControllers(r *gin.Engine) {
	r.GET("/status", StatusController)
	r.GET("/i/:hash", RedirectController)

	username := os.Getenv("QUAN_ADMIN_USERNAME")
	password := os.Getenv("QUAN_ADMIN_PASSWORD")
	if username == "" || password == "" {
		panic("QUAN_ADMIN_USERNAME or QUAN_ADMIN_PASSWORD not set.")
	}

	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	authorized.POST("/set", SetController)
	authorized.GET("/list", ListController)
	authorized.POST("/delete", DeleteController)
}

func OK(ctx *gin.Context, data gin.H) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"head": gin.H{
			"code": 0,
			"msg":  "ok",
		},
		"body": data,
	})
}

func Err(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
		"head": gin.H{
			"code": code,
			"msg":  msg,
		},
		"body": gin.H{},
	})
}

func GetDB(ctx *gin.Context) (*db.DBModel, bool) {
	database_, exists := ctx.Get("database")
	if !exists {
		return nil, false
	}
	database, ok := database_.(*db.DBModel)
	if !ok {
		return nil, false
	}
	return database, true
}

func HashToShortenUrl(hash string) string {
	u1 := *(config.BaseUrl.URL)
	u := &u1
	u.Path = "/i/" + hash
	return u.String()
}

func Redirect(ctx *gin.Context, code int, location string) {
	ctx.Redirect(code, location)
}
