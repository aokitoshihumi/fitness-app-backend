package auth

import (
	"context"
	"errors"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseVerifier struct {
	client *firebaseauth.Client
}

func NewFirebaseVerifier(ctx context.Context, projectID string) (*FirebaseVerifier, error) {
	// 本番での事故を防ぐためにエミュレーターの環境変数が設定されいたらエラー出す
	// エミュレーターという開発環境用の、Googleの秘密鍵で検証を行うステップをかっ飛ばす機能があるらしい
	if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") != "" {
		return nil, errors.New("FIREBASE_AUTH_EMULATOR_HOST が設定されています。署名検証が無効になるため起動を中止します")
	}
	// FirebaseSDKを使うために、Appを作る
	// 検証に必要なのは Google の公開鍵だけなので、認証情報なし(WithoutAuthentication)でOK
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, option.WithoutAuthentication())
	if err != nil {
		return nil, fmt.Errorf("Firebaseの初期化に失敗しました: %w", err)
	}
	// authクライアントを作る
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
