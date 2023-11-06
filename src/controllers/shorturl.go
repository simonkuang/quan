package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simonkuang/quan/src/codec"
	"github.com/simonkuang/quan/src/config"
	"github.com/simonkuang/quan/src/db"
)

type EntityList struct {
	Shorten string
	Url     string
	Time    *time.Time
}

// print configs and flags
func StatusController(ctx *gin.Context) {
	t := time.Now()
	ctx.JSON(http.StatusOK, gin.H{
		"config": gin.H{
			"BaseUrl":   config.BaseUrl.String(),
			"Length":    config.Length,
			"DbFile":    config.DbFile,
			"LogFile":   config.LogFile,
			"CharRange": config.CharRange,
		},

		"time":      t,
		"timestamp": t.UnixMilli(),
	})
}

func SetController(ctx *gin.Context) {
	originUrl := ctx.PostForm("url")
	forceHash := ctx.PostForm("id")

	if originUrl == "" {
		Err(ctx, 501, "param 'url' should not be empty")
		return
	}

	database, exists := GetDB(ctx)
	if !exists {
		Err(ctx, 502, "Database not initialized at beginning.")
		return
	}

	startPostion := 0
	var (
		hashStr string
		err     error
	)
	for {
		if forceHash != "" {
			hashStr = forceHash
			forceHash = ""
		} else {
			hashStr, err = codec.Encode(originUrl, startPostion)
			startPostion += config.Length
			if err != nil {
				Err(ctx, 503, "Failed on Encoding Original URL: "+err.Error())
				return
			}
		}
		if checkConflict(database, hashStr, originUrl) {
			continue
		}

		t := time.Now()
		entity := &codec.ShortUrlEntity{
			Url:      originUrl,
			ShortUrl: HashToShortenUrl(hashStr),
			Time:     &t,
		}
		jsonStr, err := json.Marshal(entity)
		if err != nil {
			continue
		}

		err = database.Put(hashStr, string(jsonStr))
		if err != nil {
			Err(ctx, 504, "Failed Putting Data: "+err.Error())
			return
		}

		OK(ctx, gin.H{
			"hash": hashStr,
			"url":  HashToShortenUrl(hashStr),
		})
		return
	}
}

func checkConflict(database *db.DBModel, hash string, url string) bool {
	storedUrl, err := database.Get(hash)
	if err != nil {
		return false
	}
	if storedUrl == "" {
		return false
	}
	shortUrlEntity := &codec.ShortUrlEntity{}
	err = json.Unmarshal([]byte(storedUrl), shortUrlEntity)
	if err != nil {
		return false
	}
	if shortUrlEntity.Url == url {
		return false
	}
	return true
}

func RedirectController(ctx *gin.Context) {
	database, exists := GetDB(ctx)
	if !exists {
		Err(ctx, 502, "Database not initialized at beginning.")
		return
	}

	hash := ctx.Param("hash")
	if hash == "" {
		Redirect(ctx, http.StatusTemporaryRedirect, config.DefaultRedirectUrl)
		return
	}

	entityStr, err := database.Get(hash)
	if err != nil {
		Redirect(ctx, http.StatusTemporaryRedirect, config.DefaultRedirectUrl)
		return
	}
	entity := &codec.ShortUrlEntity{}
	err = json.Unmarshal([]byte(entityStr), entity)
	if err != nil {
		Redirect(ctx, http.StatusTemporaryRedirect, config.DefaultRedirectUrl)
		return
	}
	Redirect(ctx, http.StatusTemporaryRedirect, entity.Url)
}

func ListController(ctx *gin.Context) {
	lastKey := ctx.Query("last")
	sizeStr := ctx.Query("size")
	var size int
	if sizeStr == "" {
		size = config.ListSize
	} else {
		if size1, err := strconv.Atoi(sizeStr); err != nil {
			size = config.ListSize
		} else {
			size = size1
		}
	}
	if size > 10000 || size < 1 {
		size = config.ListSize
	}

	database, exists := GetDB(ctx)
	if !exists {
		Err(ctx, 502, "Database not initialized at beginning.")
		return
	}
	iter := database.GetLevelDB().NewIterator(nil, nil)
	defer iter.Release()
	if lastKey == "" {
		iter.Seek([]byte(lastKey))
	}
	count := 0
	var EntityListArr []EntityList
	for iter.Next() {
		count++
		entity := &codec.ShortUrlEntity{}
		err := json.Unmarshal([]byte(iter.Value()), entity)
		if err != nil {
			continue
		}
		EntityListArr = append(EntityListArr, EntityList{
			Shorten: entity.ShortUrl,
			Url:     entity.Url,
			Time:    entity.Time,
		})
		if count > size {
			break
		}
	}
	OK(ctx, gin.H{
		"list": EntityListArr,
	})
}

func DeleteController(ctx *gin.Context) {
	ids := ctx.PostFormArray("ids")
	if len(ids) == 0 {
		OK(ctx, gin.H{
			"deleted": []string{},
		})
	}

	database, exists := GetDB(ctx)
	if !exists {
		Err(ctx, 502, "Database temporarily unavailable.")
		return
	}

	done := []string{}
	undone := []string{}
	for i := range ids {
		if err := database.Delete(ids[i]); err != nil {
			undone = append(undone, ids[i])
		} else {
			done = append(done, ids[i])
		}
	}
	OK(ctx, gin.H{
		"deleted": done,
		"failed":  undone,
	})
}
