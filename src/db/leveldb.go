package db

import (
	"errors"
	"fmt"

	"github.com/simonkuang/quan/src/config"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type DB interface {
	Connect()
	Get(key string) (string, error)
	Put(key string, value string) error
	Close() error
}

type DBModel struct {
	ldb *leveldb.DB
}

func GetDbFileName() string {
	config.DBVersionStep++
	return fmt.Sprintf("%s.%04d", config.DbFile, config.DBVersionStep)
}

// connect to leveldb
func (m *DBModel) Connect(dbFile string) {
	levelDBOpt := &opt.Options{
		OpenFilesCacheCapacity: config.LevelDBOpenFilesCacheCapacity,
		BlockCacheCapacity:     config.LevelDBCacheMemory / 2,
		Filter:                 nil,
		ReadOnly:               false,
	}

	ldb, err := leveldb.OpenFile(dbFile, levelDBOpt)
	if err != nil {
		panic(err)
	}
	m.ldb = ldb
}

// close db
func (m *DBModel) Close() error {
	if m.ldb == nil {
		return nil
	}
	return m.ldb.Close()
}

// retrieve value from db
func (m *DBModel) Get(key string) (string, error) {
	if m.ldb == nil {
		return "", errors.New("DB Not Connected")
	}
	exists, err := m.ldb.Has([]byte(key), nil)
	if err != nil {
		return "", err
	}
	if !exists { // not found, return empty string without error
		return "", nil
	}
	val, err := m.ldb.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

// put value
func (m *DBModel) Put(key string, value string) error {
	if m.ldb == nil {
		return errors.New("DB Not Connected")
	}
	return m.ldb.Put([]byte(key), []byte(value), &opt.WriteOptions{
		Sync: true,
	})
}

// get *leveldb.DB
func (m *DBModel) GetLevelDB() *leveldb.DB {
	return m.ldb
}

func (m *DBModel) Delete(id string) error {
	return m.ldb.Delete([]byte(id), &opt.WriteOptions{
		Sync: true,
	})
}
