package dao

import (
	"context"
	"fmt"
	"strconv"

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

func (s *status) FindByID(ctx context.Context, id int64) (*object.Status, error) {
	entity := new(object.Status)
	err := s.db.QueryRowxContext(ctx, "SELECT id, account_id, content, create_at FROM status WHERE id = ?", id).StructScan(entity)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	// fmt.Println(a.ID, a.AccountID, a.Content)

	return entity, nil
}

func (s *status) GetPublicTimelines(ctx context.Context, q object.Query) ([]object.Status, error) {
	sql := "SELECT id, account_id, content, create_at FROM status WHERE "

	//sinceIDより大きいIDのみ取得
	if q.SinceID == "" {
		q.SinceID = "0"
	}
	sql += "id > " + q.SinceID

	//maxIDより小さいIDのみ取得
	if q.MaxID != "" {
		sql += " AND id < " + q.MaxID
	}

	//limitのデフォルト値は40
	if q.Limit == "" {
		q.Limit = "40"
	}
	//limitの最大値は80
	if tmp, _ := strconv.Atoi(q.Limit); tmp > 80 {
		q.Limit = "80"
	}

	sql += " LIMIT " + q.Limit

	fmt.Println(sql)

	rows, err := s.db.QueryContext(ctx, sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	a := make([]object.Status, 0)

	for rows.Next() {
		var tmp object.Status
		rows.Scan(&tmp.ID, &tmp.AccountID, &tmp.Content, &tmp.CreateAt)
		a = append(a, tmp)
	}

	return a, nil

}
