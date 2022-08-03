package dao

import (
	"context"
	"fmt"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	status struct {
		db *sqlx.DB
	}
)

func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

func debugPostTable(ctx context.Context, s *status) {
	rows, _ := s.db.QueryContext(ctx, "SELECT id, account_id, content FROM STATUS")
	defer rows.Close()

	var a object.Status

	for rows.Next() {
		rows.Scan(&a.ID, &a.AccountID, &a.Content)
		fmt.Println(a.ID, a.AccountID, a.Content)
	}
}

func (s *status) PostStatus(ctx context.Context, entity object.Status) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO STATUS (account_id, content) VALUES (?, ?)",
		entity.AccountID, entity.Content,
	)

	if err != nil {
		return err
	}

	debugPostTable(ctx, s)

	return nil
}
