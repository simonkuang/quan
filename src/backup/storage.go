package backup

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type BackupStorage struct {
	initialized bool
	bucketName  string
	bucketSize  int
	client      *storage.Client
	ctx         *context.Context
}

// for default values
func NewBackupStorage() *BackupStorage {
	stor := &BackupStorage{initialized: true}
	return stor
}

func (m *BackupStorage) Connect(credentialContent string, ctx *context.Context) *BackupStorage {
	config, err := google.JWTConfigFromJSON([]byte(credentialContent), storage.ScopeReadWrite)
	if err != nil {
		log.Printf("[ERROR] get backup storage client failed: %v\n", err)
		return m
	}

	client, err := storage.NewClient(*ctx, option.WithTokenSource(config.TokenSource(*ctx)))
	if err != nil {
		log.Printf("[ERROR] get backup storage client failed: %v\n", err)
		return m
	}
	m.client = client
	m.ctx = ctx
	m.initialized = true
	log.Printf("[INFO] backup storage client connected")
	return m
}

func (m *BackupStorage) Close() {
	if m.client == nil {
		return
	}
	m.client.Close()
	m.client = nil
	m.initialized = false
	log.Printf("[INFO] backup storage client closed")
}

func (m *BackupStorage) Setup(bucket string, size int) *BackupStorage {
	if !m.initialized {
		log.Printf("[WARN] backup storage client is not initialized")
	}
	if m.client == nil {
		log.Printf("[ERROR] backup storage client is nil")
		return m
	}
	m.bucketName = bucket
	m.bucketSize = size
	return m
}

func (m *BackupStorage) Upload(content []byte, objectName string) *BackupStorage {
	if !m.initialized {
		log.Printf("[WARN] backup storage client is not initialized")
	}
	if m.client == nil {
		log.Printf("[ERROR] backup storage client is nil")
		return m
	}
	wc := m.client.Bucket(m.bucketName).Object(objectName).NewWriter(*m.ctx)
	if _, err := wc.Write(content); err != nil {
		log.Printf("[ERROR] Failed to write to bucket %s, object %s: %v", m.bucketName, objectName, err)
	}
	if err := wc.Close(); err != nil {
		log.Printf("[ERROR] Failed to close bucket %s, object %s: %v", m.bucketName, objectName, err)
	}
	return m
}

func (m *BackupStorage) Download(objectName string) []byte {
	if !m.initialized {
		log.Printf("[WARN] backup storage client is not initialized")
	}
	if m.client == nil {
		log.Printf("[ERROR] backup storage client is nil")
		return nil
	}
	rc, err := m.client.Bucket(m.bucketName).Object(objectName).NewReader(*m.ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to read from bucket %s, object %s: %v", m.bucketName, objectName, err)
		return nil
	}
	defer rc.Close()
	var data = make([]byte, m.bucketSize*2048)
	_, err = rc.Read(data)
	if err != nil {
		log.Printf("[ERROR] Failed to read from bucket %s, object %s: %v", m.bucketName, objectName, err)
		return nil
	}
	rc.Close()
	return data
}
