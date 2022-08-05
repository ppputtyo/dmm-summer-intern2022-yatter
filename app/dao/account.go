package dao

import (
	"context"
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

	// debugTable(ctx, r)

	return nil
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	entity.FollowersCount = r.GetFollowersCount(ctx, entity.ID)
	entity.FollowingCount = r.GetFollowingCount(ctx, entity.ID)

	return entity, nil
}

// FindByID : IDからユーザを取得
func (r *account) FindByID(ctx context.Context, ID int64) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where id = ?", ID).StructScan(entity)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	entity.FollowersCount = r.GetFollowersCount(ctx, ID)
	entity.FollowingCount = r.GetFollowingCount(ctx, ID)

	return entity, nil
}

func (r *account) GetFollowersCount(ctx context.Context, ID int64) int64 {
	query := object.FollowersQuery{
		MaxID:   1000000,
		SinceID: -1,
		Limit:   1000000,
	}
	var count int64

	r.db.QueryRowxContext(
		ctx,
		"select count(follower_id) from relation where followee_id = ? AND follower_id < ? AND follower_id > ? LIMIT ?",
		ID, query.MaxID, query.SinceID, query.Limit).Scan(&count)

	return count
}

func (r *account) GetFollowingCount(ctx context.Context, ID int64) int64 {
	var count int64
	r.db.QueryRowxContext(ctx, "select count(followee_id) from relation where follower_id = ? LIMIT ?", ID, 1000000).Scan(&count)

	return count
}

func debugRelation(ctx context.Context, r *account) {
	rows, _ := r.db.QueryContext(ctx, "SELECT * FROM relation")
	defer rows.Close()

	for rows.Next() {
		var follower_id, followee_id int64
		rows.Scan(&follower_id, &followee_id)
		res, _ := r.FindByID(ctx, follower_id)
		follower_name := res.Username

		res, _ = r.FindByID(ctx, followee_id)
		followee_name := res.Username

		fmt.Printf("%s => %s\n", follower_name, followee_name)
	}
}

func (r *account) GetRelation(ctx context.Context, myID, targetID int64) (*object.Relation, error) {
	relation := new(object.Relation)
	relation.ID = targetID
	relation.Following = false
	relation.FollowedBy = false

	debugRelation(ctx, r)

	rows, err := r.db.QueryContext(ctx, "select * from relation where follower_id = ? AND followee_id = ?", myID, targetID)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		relation.Following = true
	}

	rows, err = r.db.QueryContext(ctx, "select * from relation where follower_id = ? AND followee_id = ?", targetID, myID)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		relation.FollowedBy = true
	}

	return relation, nil
}

func (r *account) Follow(ctx context.Context, myID, targetID int64) error {
	if myID == targetID {
		return fmt.Errorf("cannot follow yourself")
	}

	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO relation (follower_id, followee_id) VALUES (?, ?)",
		myID, targetID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *account) Unfollow(ctx context.Context, myID, targetID int64) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM relation WHERE follower_id = ? AND followee_id = ?",
		myID, targetID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *account) GetFollowing(ctx context.Context, ID int64, limit int) ([]object.Account, error) {
	res := make([]object.Account, 0)

	rows, err := r.db.QueryContext(ctx, "select followee_id from relation where follower_id = ? LIMIT ?", ID, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var followee_id int64
		rows.Scan(&followee_id)

		entity, err := r.FindByID(ctx, followee_id)
		if err != nil {
			continue
		}
		res = append(res, *entity)
	}

	debugRelation(ctx, r)

	return res, nil
}

func (r *account) GetFollowers(ctx context.Context, ID int64, query object.FollowersQuery) ([]object.Account, error) {
	res := make([]object.Account, 0)

	rows, err := r.db.QueryContext(
		ctx,
		"select follower_id from relation where followee_id = ? AND follower_id < ? AND follower_id > ? LIMIT ?",
		ID, query.MaxID, query.SinceID, query.Limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var follower_id int64
		rows.Scan(&follower_id)

		entity, err := r.FindByID(ctx, follower_id)
		if err != nil {
			continue
		}
		res = append(res, *entity)
	}

	debugRelation(ctx, r)

	return res, nil
}
