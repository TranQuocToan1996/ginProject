package gcloud

import (
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/TranQuocToan1996/ginProject/models"
)

var (
	sysCfg   *models.Config
	fbApp    *firebase.App
	fsClient *firestore.Client
	rtClient *db.Client
	// storageClient *storage.Client
)

func getKey() string {
	key := os.Getenv("key")
	if len(key) == 0 {
		key = "firestore.json"
	}
	return key
}
