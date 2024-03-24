package users

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/ptsypyshev/gb-golang-level3-new/internal/database"
)

const (
	createQuery      = `INSERT INTO users(id,username,password) VALUES($1,$2,$3)`
	readByIDQuery    = `SELECT id, username, password, created_at, updated_at FROM users u WHERE id=$1`
	readByIDUsername = `SELECT id, username, password, created_at, updated_at FROM users u WHERE username=$1`
)

func New(userDB *pgx.Conn, timeout time.Duration) *Repository {
	return &Repository{userDB: userDB, timeout: timeout}
}

type Repository struct {
	userDB  *pgx.Conn
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateUserReq) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.userDB.Exec(ctx, createQuery, req.ID, req.Username, req.Password)
	if err != nil {
		return database.User{}, ErrCreateUser
	}
	u.ID = req.ID
	u.Username = req.Username
	u.Password = req.Password
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	return u, nil
}

func (r *Repository) FindByID(ctx context.Context, userID uuid.UUID) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	err := r.userDB.QueryRow(ctx, readByIDQuery, userID).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	switch {
	case err == nil:
		return u, nil
	case errors.Is(pgx.ErrNoRows, err):
		return database.User{}, ErrNoUserFound
	default:
		return database.User{}, err
	}
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	err := r.userDB.QueryRow(ctx, readByIDUsername, username).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	switch {
	case err == nil:
		return u, nil
	case errors.Is(pgx.ErrNoRows, err):
		return database.User{}, ErrNoUserFound
	default:
		return database.User{}, err
	}
}
