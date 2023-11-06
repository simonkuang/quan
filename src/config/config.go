package config

var (
	LevelDBCacheMemory int    = 32 // in megabytes
	LevelDBConnections int    = 32 // concurreny level
	LevelDBNamespace   string = "quan."

	LevelDBOpenFilesCacheCapacity int = 1

	SecretPrefix string = "k8EtkYvyDuzQSU9N"
)
