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

func (s *status) PostStatus(ctx context.Context, entity *object.Status) error {
	res, err := s.db.ExecContext(
		ctx,
		"INSERT INTO STATUS (account_id, content) VALUES (?, ?)",
		entity.AccountID, entity.Content,
	)

	if err != nil {
		return err
	}

	entity.ID, _ = res.LastInsertId()
	// debugPostTable(ctx, s)

	return nil
}

func (s *status) FindByPostID(ctx context.Context, id int64) (*object.Status, error) {
	var a object.Status
	err := s.db.QueryRow("SELECT id, account_id, content FROM status WHERE id = ?", id).Scan(&a.ID, &a.AccountID, &a.Content)

	if err != nil {
		return nil, err
	}
	fmt.Println(a.ID, a.AccountID, a.Content)

	return &a, nil
}
