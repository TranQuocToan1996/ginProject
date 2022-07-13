package gcloud

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/TranQuocToan1996/ginProject/models"
	"google.golang.org/api/option"
)

func InitFbApp(ctx context.Context) error {
	key := getKey()
	otps := option.WithCredentialsFile(key)
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var err error
	fbApp, err = firebase.NewApp(ctx, nil, otps)
	if err != nil {
		log.Fatalln(err)
	}

	fsClient, err = fbApp.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// defer fsClient.Close()

	rtClient, err = fbApp.Database(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func Close() {
	fsClient.Close()
	// storageClient.Close()
}

func RunFS(ctx context.Context, collName string, f func(client *firestore.Client) error) error {
	if fsClient == nil {
		InitFbApp(ctx)
	}

	err := f(fsClient)
	return err
}

func ExecuteRT(ctx context.Context, f func(app *firebase.App) error) error {

	if fbApp == nil {
		err := InitFbApp(context.Background())
		if err != nil {
			return err
		}
	}

	return f(fbApp)
}

func GetConfig(ctx context.Context) (*models.Config, error) {
	if sysCfg != nil {
		return sysCfg, nil
	}

	sysCfg = &models.Config{}

	//TODO: get config in the firestore

	return sysCfg, nil
}
