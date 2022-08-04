package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Account interface {
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)
	//新規アカウントの作成
	CreateNewAccount(ctx context.Context, account object.Account) error
	//IDでアカウントを検索
	FindByID(ctx context.Context, ID int64) (*object.Account, error)
}
