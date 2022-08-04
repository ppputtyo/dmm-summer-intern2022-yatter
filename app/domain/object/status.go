package object

import (
	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	object struct {
		db *sqlx.DB
	}

	Query struct {
		OnlyMedia string
		MaxID     string
		SinceID   string
		Limit     string
	}
)

type (
	StatusID = int64

	// Account account
	Status struct {
		// The internal ID of the account
		ID StatusID `json:"id" db:"id"`

		// The username of the account
		AccountID AccountID `json:"account_id,omitempty" db:"account_id"`

		// The username of the account
		Content string `json:"content" db:"content"`

		// The time the account was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
	}
)
