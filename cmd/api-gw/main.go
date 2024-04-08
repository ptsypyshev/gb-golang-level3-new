package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ptsypyshev/gb-golang-level3-new/internal/env"
)

const ShutdownTimeout = 3 * time.Second

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := runMain(ctx); err != nil {
		log.Fatal(err)
	}
}

func runMain(ctx context.Context) error {
	e, c, err := env.Setup(ctx)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	httpServer := e.APIGWHTTPServer

	go func() {
		<-ctx.Done()
		// если посылаем сигнал завершения то завершаем работу нашего сервера
		httpServer.Close()
	}()

	go func() {
		defer wg.Done()

		slog.Info(fmt.Sprintf("api-gw http was started %s", e.Config.APIGWService.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			slog.Error("api-gw http server", slog.Any("err", err))
			return
		}

		httpServer.Close()
	}()

	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout) //nolint:contextcheck
	defer cancel()

	c.Close(ctx)

	return nil
}
