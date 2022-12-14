package object

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type (
	AccountID    = int64
	PasswordHash = string

	// Account account
	Account struct {
		// The internal ID of the account
		ID AccountID `json:"-"`

		// The username of the account
		Username string `json:"username,omitempty"`

		// The username of the account
		PasswordHash string `json:"-" db:"password_hash"`

		FollowersCount int64 `json:"followers_count" db:"followers_count"`

		FollowingCount int64 `json:"following_count" db:"following_count"`

		// The account's display name
		DisplayName *string `json:"display_name,omitempty" db:"display_name"`

		// URL to the avatar image
		Avatar *string `json:"avatar,omitempty"`

		// URL to the header image
		Header *string `json:"header,omitempty"`

		// Biography of user
		Note *string `json:"note,omitempty"`

		// The time the account was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
	}

	Relation struct {
		ID         AccountID `json:"id"`
		Following  bool      `json:"following"`
		FollowedBy bool      `json:"followed_by"`
	}

	FollowersQuery struct {
		MaxID   int
		SinceID int
		Limit   int
	}
)

// Check if given password is match to account's password
func (a *Account) CheckPassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(pass)) == nil
}

// Hash password and set it to account object
func (a *Account) SetPassword(pass string) error {
	passwordHash, err := generatePasswordHash(pass)
	if err != nil {
		return fmt.Errorf("generate error: %w", err)
	}
	a.PasswordHash = passwordHash
	return nil
}

func generatePasswordHash(pass string) (PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password failed: %w", errors.WithStack(err))
	}
	return PasswordHash(hash), nil
}
