package config

import (
	"flag"
	"net/url"

	"github.com/gin-gonic/gin"
)

type URLValue struct {
	URL *url.URL
}

func (v URLValue) String() string {
	if v.URL != nil {
		return v.URL.String()
	}
	return ""
}

func (v URLValue) Set(s string) error {
	if u, err := url.Parse(s); err != nil {
		return err
	} else {
		*v.URL = *u
	}
	return nil
}

var (
	Length             int      // short url hash length. Default: 5
	DbFile             string   // database file path. Default: ./quan.db
	LogFile            string   // log file path. Default: ./quan.log
	BaseUrl            URLValue // base url. Default: http://localhost:8080
	CharRange          int      // character range. 62 for [0-9a-zA-Z], 36 for [0-9a-z]. Default: 62
	DefaultRedirectUrl string   // redirect url. Default: /
	ListSize           int
	Port               int
)

func FlagsInit(ctx *gin.Engine) {
	// fl = flag.NewFlagSet("quan", flag.ExitOnError)

	flag.Var(&BaseUrl, "base-url", "base url for short url. Default: http://localhost:8080")
	flag.IntVar(&Length, "length", 5, "short url hash length. Default: 5")
	flag.StringVar(&DbFile, "db-file", "./quan.db", "database file path. Default: ./quan.db")
	flag.StringVar(&LogFile, "log-file", "./quan.log", "log file path. Default: ./quan.log")
	flag.IntVar(&CharRange, "char-range", 62, "character range. 62 for [0-9a-zA-Z], 36 for [0-9a-z]. Default: 62")
	flag.StringVar(&DefaultRedirectUrl, "default-redirect-url", "/", "redirect url. Default: /")
	flag.IntVar(&ListSize, "list-size", 50, "list size. Default: 50")
	flag.IntVar(&Port, "port", 8080, "port. Default: 8080")

	// set default value for BaseUrl
	if BaseUrl.String() == "" {
		if u, err := url.Parse("http://localhost:8080"); err == nil {
			BaseUrl.URL = u
		}
	}

	flag.Parse()
}
