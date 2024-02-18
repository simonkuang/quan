package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/simonkuang/quan/src/codec"
	"github.com/simonkuang/quan/src/config"
	"github.com/simonkuang/quan/src/db"
	"github.com/syndtr/goleveldb/leveldb"
)

var stor *BackupStorage

func Backup(db *leveldb.DB, backupName string, backupSize int) bool {
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	var count int = 0
	var fileCount = 1
	var fileList []string
	list := make(map[string]codec.ShortUrlEntity)
	ctx := context.Background()
	stor := getStorage(&ctx)
	for iter.Next() {
		count++
		var entity codec.ShortUrlEntity
		var hashStr string
		err := json.Unmarshal(iter.Value(), &entity)
		if err != nil {
			log.Printf("[ERROR] json.Unmarshal() key (%s) while backup: %s\n", iter.Key(), iter.Value())
		}
		if entity.Hash == "" {
			pos := strings.LastIndex(entity.ShortUrl, "/")
			hashStr = entity.ShortUrl[pos+1:]
		} else {
			hashStr = entity.Hash
		}
		list[hashStr] = entity

		if count > backupSize {
			if b, err := json.Marshal(list); err != nil {
				log.Printf("[ERROR] failed json.Marshal() while encoding list, err: %v", err)
			} else {
				fileName := fmt.Sprintf("%s-%03d.json", backupName, fileCount)
				stor.Upload(b, fileName)
				fileList = append(fileList, fileName)
				fileCount++
			}

			count = 0
			list = make(map[string]codec.ShortUrlEntity)
		}
	}
	if len(list) > 0 {
		if b, err := json.Marshal(list); err != nil {
			log.Printf("[ERROR] failed json.Marshal() while encoding list, err: %v", err)
		} else {
			fileName := fmt.Sprintf("%s-%03d.json", backupName, fileCount)
			fileList = append(fileList, fileName)
			stor.Upload(b, fileName)
		}
	}

	// index file
	if len(fileList) > 0 {
		fileName := fmt.Sprintf("%s.index.json", backupName)
		if b, err := json.Marshal(fileList); err != nil {
			log.Printf("[ERROR] failed json.Marshal() while encoding fileList, err: %v", err)
		} else {
			stor.Upload(b, fileName)
		}
	}

	stor.Close()
	iter.Release()
	return true
}

func Restore(backupName string, stepVersion int) *db.DBModel {
	currentDB := &db.DBModel{}
	currentDB.Connect(db.GetDbFileName())

	ctx := context.Background()
	stor := getStorage(&ctx)
	// read from index file
	indexFile := fmt.Sprintf("%s.index.json", backupName)
	indexContent := stor.Download(indexFile)
	if indexContent == nil {
		log.Printf("[ERROR] failed to download index file: %s", indexFile)
		return nil
	}
	var fileList []string
	indexContentStr := strings.Trim(string(indexContent), "\x00")
	if err := json.Unmarshal([]byte(indexContentStr), &fileList); err != nil {
		log.Printf("[ERROR] failed to json.Unmarshal() error: %v, indexContent: %s", err, indexContent)
		return nil
	}

	// restore data
	for _, fileName := range fileList {
		content := stor.Download(fileName)
		if content == nil {
			log.Printf("[ERROR] failed to download file: %s", fileName)
			continue
		}
		contentStr := strings.Trim(string(content), "\x00")
		var list map[string]codec.ShortUrlEntity
		if err := json.Unmarshal([]byte(contentStr), &list); err != nil {
			log.Printf("[ERROR] failed to json.Unmarshal() backup file: %s", fileName)
			continue
		}
		// fmt.Printf(" >> [DEBUG] backup data is: %s\n", contentStr)
		for hashStr, entity := range list {
			// fmt.Printf(" >> [DEBUG] hashStr is: %s\n", hashStr)
			// fmt.Printf(" >> [DEBUG] entity is: %v\n", entity)
			if b, err := json.Marshal(entity); err != nil {
				log.Printf("[ERROR] failed json.Marshal() entity while restore data: %v, %s", entity, err)
			} else {
				if err = currentDB.Put(hashStr, string(b)); err != nil {
					log.Printf("[ERROR] failed to db.Put() key: %s, value: %s", hashStr, entity.Url)
					return nil
				}
			}
		}
	}

	stor.Close()
	return currentDB
}

func getStorage(ctx *context.Context) *BackupStorage {
	if stor != nil {
		return stor
	}
	stor := NewBackupStorage()
	var credentialContent string
	if config.CredentialContentGoogle != "" {
		credentialContent = config.CredentialContentGoogle
	} else if config.CredentialFileGoogle != "" {
		b, err := os.ReadFile(config.CredentialFileGoogle)
		if err != nil {
			log.Printf("[ERROR] failed to read credential file: %s", config.CredentialFileGoogle)
		}
		credentialContent = string(b)
	} else {
		log.Printf("[ERROR] no credential file or content. Backup failed.")
		return nil
	}
	return stor.Connect(credentialContent, ctx).Setup(config.BackupBucket, config.BackupSize)
}
