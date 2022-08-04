package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

func debugTable(ctx context.Context, r *account) {
	rows, _ := r.db.QueryContext(ctx, "SELECT id, username, password_hash FROM ACCOUNT")
	defer rows.Close()

	var a object.Account

	for rows.Next() {
		rows.Scan(&a.ID, &a.Username, &a.PasswordHash)
		fmt.Println(a.ID, a.Username, a.PasswordHash)
	}
}

func (r *account) CreateNewAccount(ctx context.Context, entity object.Account) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO ACCOUNT (id, username, password_hash) VALUES (?, ?, ?)",
		entity.ID, entity.Username, entity.PasswordHash,
	)

	if err != nil {
		return err
	}

	debugTable(ctx, r)

	return nil
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}

// FindByID : IDからユーザを取得
func (r *account) FindByID(ctx context.Context, ID int64) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where id = ?", ID).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}
