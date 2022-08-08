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
	//フォロー
	Follow(ctx context.Context, myID, targetID int64) error
	//フォロー解除
	Unfollow(ctx context.Context, myID, targetID int64) error
	//ユーザ間の関係を取得
	GetRelation(ctx context.Context, myID, targetID int64) (*object.Relation, error)
	//フォローしているユーザを取得
	GetFollowing(ctx context.Context, ID int64, limit int) ([]object.Account, error)
	//フォロワーを取得
	GetFollowers(ctx context.Context, ID int64, query object.FollowersQuery) ([]object.Account, error)
}
