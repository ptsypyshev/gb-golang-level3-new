package env

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/sethvargo/go-envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ptsypyshev/gb-golang-level3-new/internal/database/links"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/database/users"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/env/config"
)

type Env struct {
	UsersRepository *users.Repository
	LinksRepository *links.Repository
}

func Setup(ctx context.Context) (*Env, error) {
	var cfg config.Config
	env := &Env{}

	if err := envconfig.Process(ctx, &cfg); err != nil { //nolint:typecheck
		return nil, fmt.Errorf("env processing: %w", err)
	}

	linksDB, err := mongo.Connect(
		ctx, &options.ClientOptions{
			ConnectTimeout: &cfg.LinksDB.ConnectTimeout,
			Hosts:          []string{fmt.Sprintf("%s:%d", cfg.LinksDB.Host, cfg.LinksDB.Port)},
			MaxPoolSize:    &cfg.LinksDB.MaxPoolSize,
			MinPoolSize:    &cfg.LinksDB.MinPoolSize,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect: %w", err)
	}

	usersClient, err := pgx.Connect(ctx, cfg.UsersDB.ConnectionURL())
	if err != nil {
		return nil, err
	}

	usersRepository := users.New(usersClient, 5*time.Second)                        // вынести в конфиг duration
	linksRepository := links.New(linksDB.Database(cfg.LinksDB.Name), 5*time.Second) // вынести в конфиг duratino
	env.LinksRepository = linksRepository
	env.UsersRepository = usersRepository

	return env, nil
}
