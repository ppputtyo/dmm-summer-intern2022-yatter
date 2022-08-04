package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	FindByPostID(ctx context.Context, userID int64) (*object.Status, error)
	PostStatus(ctx context.Context, status *object.Status) error
	GetPublicTimelines(ctx context.Context, q object.Query) ([]object.Status, error)
}
