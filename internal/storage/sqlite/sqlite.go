package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sso/internal/domain/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (storage *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (userID int64, err error) {
	panic("implement me")
}

func (storage *Storage) User(
	ctx context.Context,
	email string,
) (user models.User, err error) {
	panic("implement me")
}

func (storage *Storage) IsAdmin(
	ctx context.Context,
	userID int64,
) (isAdmin bool, err error) {
	panic("implement me")
}

func (storage *Storage) App(
	ctx context.Context,
	appID int,
) (app models.App, err error) {
	panic("implement me")
}
