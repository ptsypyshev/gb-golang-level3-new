package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/google/uuid"

	"github.com/ptsypyshev/gb-golang-level3-new/internal/database/users"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/env"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := runMain(ctx); err != nil {
		log.Panic(err)
	}
}

func runMain(ctx context.Context) error {
	e, err := env.Setup(ctx)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	_ = e
	create, err := e.UsersRepository.Create(
		ctx, users.CreateUserReq{
			ID:       uuid.New(),
			Username: "random",
			Password: "password",
		},
	)
	if err != nil {
		return err
	}

	found, err := e.UsersRepository.FindByID(ctx, create.ID)
	if err != nil {
		return err
	}

	foundBy, err := e.UsersRepository.FindByUsername(ctx, "random")
	if err != nil {
		return err
	}

	fmt.Println(create, found, foundBy)
	return nil
}
