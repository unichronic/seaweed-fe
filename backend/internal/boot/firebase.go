package boot

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func NewFirebaseAuth() (*auth.Client, error) {
	if os.Getenv("DUMMY_AUTH") == "true" && os.Getenv("STAGE") != "prod" {
		return nil, nil
	}
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	if serviceAccountPath == "" {
		app, err := firebase.NewApp(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		return app.Auth(context.Background())
	}

	opt := option.WithServiceAccountFile(serviceAccountPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return app.Auth(context.Background())
}
