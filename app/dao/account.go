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
	rows, _ := r.db.QueryContext(
		ctx,
		"SELECT id, username, password_hash FROM ACCOUNT",
	)
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

	return nil
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(
		ctx,
		"SELECT * FROM account WHERE username = ?",
		username,
	).StructScan(entity)

	if err != nil {
		return nil, err
	}

	if entity.FollowersCount, err = r.GetFollowersCount(ctx, entity.ID); err != nil {
		return nil, err
	}
	if entity.FollowingCount, err = r.GetFollowingCount(ctx, entity.ID); err != nil {
		return nil, err
	}

	return entity, nil
}

// FindByID : IDからユーザを取得
func (r *account) FindByID(ctx context.Context, ID int64) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(
		ctx,
		"select * from account where id = ?",
		ID,
	).StructScan(entity)
	if err != nil {
		return nil, err
	}

	if entity.FollowersCount, err = r.GetFollowersCount(ctx, ID); err != nil {
		return nil, err
	}

	if entity.FollowingCount, err = r.GetFollowingCount(ctx, ID); err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *account) GetFollowersCount(ctx context.Context, ID int64) (int64, error) {
	const INF = 10000000
	query := object.FollowersQuery{
		MaxID:   INF,
		SinceID: -1,
		Limit:   INF,
	}

	var count int64
	err := r.db.QueryRowxContext(
		ctx,
		`SELECT count(follower_id)
		FROM relation
		WHERE followee_id = ? AND follower_id < ? AND follower_id > ?
		LIMIT ?`,
		ID, query.MaxID, query.SinceID, query.Limit,
	).Scan(&count)

	if err != nil {
		return -1, err
	}

	return count, nil
}
func (r *account) GetFollowingCount(ctx context.Context, ID int64) (int64, error) {
	const INF = 10000000
	var count int64
	err := r.db.QueryRowxContext(
		ctx,
		`SELECT count(followee_id)
		FROM relation
		WHERE follower_id = ?
		LIMIT ?`,
		ID, INF,
	).Scan(&count)

	if err != nil {
		return -1, err
	}

	return count, nil
}

func debugRelation(ctx context.Context, r *account) {
	rows, _ := r.db.QueryContext(
		ctx,
		"SELECT * FROM relation",
	)
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
	relation := &object.Relation{
		ID:         targetID,
		Following:  false,
		FollowedBy: false,
	}

	debugRelation(ctx, r)

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT *
		FROM relation
		WHERE follower_id = ? AND followee_id = ?`,
		myID, targetID)
	if err != nil {
		return nil, err
	}

	// 対象のrelationが存在しない場合
	if rows.Next() {
		relation.Following = true
	}

	rows, err = r.db.QueryContext(
		ctx,
		`SELECT *
		FROM relation
		WHERE follower_id = ? AND followee_id = ?`,
		targetID, myID)

	if err != nil {
		return nil, err
	}

	// 対象のrelationが存在しない場合
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
		`DELETE FROM relation 
		WHERE follower_id = ? AND followee_id = ?`,
		myID, targetID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *account) GetFollowing(ctx context.Context, ID int64, limit int) ([]object.Account, error) {
	following := make([]object.Account, 0)

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT followee_id
		FROM relation
		WHERE follower_id = ?
		LIMIT ?`,
		ID, limit,
	)
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var followee_id int64
		rows.Scan(&followee_id)

		entity, err := r.FindByID(ctx, followee_id)
		if err != nil {
			continue
		}
		following = append(following, *entity)
	}

	debugRelation(ctx, r)

	return following, nil
}

func (r *account) GetFollowers(ctx context.Context, ID int64, query object.FollowersQuery) ([]object.Account, error) {
	followers := make([]object.Account, 0)

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT follower_id
		FROM relation
		WHERE followee_id = ? AND follower_id < ? AND follower_id > ?
		LIMIT ?`,
		ID, query.MaxID, query.SinceID, query.Limit,
	)

	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var follower_id int64
		rows.Scan(&follower_id)

		entity, err := r.FindByID(ctx, follower_id)
		if err != nil {
			continue
		}
		followers = append(followers, *entity)
	}

	debugRelation(ctx, r)

	return followers, nil
}
