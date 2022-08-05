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

	entity.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (s *status) FindByID(ctx context.Context, id int64) (*object.Status, error) {
	entity := new(object.Status)
	err := s.db.QueryRowxContext(
		ctx,
		`SELECT id, account_id, content, create_at
		FROM status
		WHERE id = ?`,
		id,
	).StructScan(entity)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}

func (s *status) GetPublicTimelines(ctx context.Context, q object.GetTimelineQuery) ([]object.Status, error) {
	rows, err := s.db.QueryxContext(
		ctx,
		`SELECT id, account_id, content, create_at
		FROM status WHERE id > ? AND id < ?
		LIMIT ?`,
		q.SinceID, q.MaxID, q.Limit,
	)

	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
	}

	timeline := make([]object.Status, 0)

	for rows.Next() {
		var status object.Status
		rows.StructScan(&status)
		timeline = append(timeline, status)
	}

	return timeline, nil
}

func (s *status) GetHomeTimelines(ctx context.Context, userID int64, q object.GetTimelineQuery, following []int64) ([]object.Status, error) {
	if len(following) == 0 {
		return nil, nil
	}

	query := `
	SELECT id, account_id, content, create_at
	FROM status
	WHERE id > ? AND id < ? AND account_id IN (?)
	LIMIT ?
	`
	query, params, err := sqlx.In(query, q.SinceID, q.MaxID, following, q.Limit)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryxContext(
		ctx,
		query,
		params...,
	)
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
	}

	timeline := make([]object.Status, 0)

	for rows.Next() {
		var status object.Status
		rows.StructScan(&status)
		timeline = append(timeline, status)
	}

	return timeline, nil
}

func (s *status) DeleteStatus(ctx context.Context, statusID int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM status WHERE id = ?", statusID)
	if err != nil {
		return err
	}
	return nil
}
