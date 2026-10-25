package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseVerifier struct {
	client *firebaseauth.Client
}

func NewFirebaseVerifier(ctx context.Context, projectID string)(*FirebaseVerifier, error){
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, option.WithoutAuthentication())
	if err != nil{
		return nil, fmt.Errorf("Firebaseの初期化に失敗しました: %w", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("Firebase Authの初期化に失敗しました: %w", err)
	}
	return &FirebaseVerifier{client: client}, nil
}

// VerifyIDToken は、トークンが本物なら Firebase の UID を返す。
func (v *FirebaseVerifier) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}