package object

import "github.com/jmoiron/sqlx"

type (
	// Implementation for repository.Account
	object struct {
		db *sqlx.DB
	}
)

type (
	StatusID = int64

	// Account account
	Status struct {
		// The internal ID of the account
		ID StatusID `json:"-"`

		// The username of the account
		AccountID string `json:"account_id,omitempty"`

		// The username of the account
		Content string `json:"content" db:"content"`

		// The time the account was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
	}
)
