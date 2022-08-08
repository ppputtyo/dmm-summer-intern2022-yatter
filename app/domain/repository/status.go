package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	//IDでStatusを取得
	FindByID(ctx context.Context, userID int64) (*object.Status, error)
	//Statusを新規作成
	PostStatus(ctx context.Context, status *object.Status) error
	//公開タイムラインを取得
	GetPublicTimelines(ctx context.Context, q object.GetTimelineQuery) ([]object.Status, error)
	//ホームタイムラインを取得
	GetHomeTimelines(ctx context.Context, userID int64, q object.GetTimelineQuery, following []int64) ([]object.Status, error)
	//Statusを削除
	DeleteStatus(ctx context.Context, statusID int64) error
}
