package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	if err := start(log); err != nil {
		log.Error("Failed to start RestWatch server", "error", err)
		os.Exit(1)
	}
}

func start(log *slog.Logger) error {
	log.Info("Starting RestWatch server")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := NewServer(NewServerOptions{
		Log:         log,
		OfflineMode: true,
	})

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.Start()
	})

	<-ctx.Done()
	log.Info("Shutting down RestWatch server")

	//eg.Go(func() error {
	//	return s.Stop()
	//})

	//if err := eg.Wait(); err != nil {
	//	log.Error("Error during shutdown", "error", err)
	//	return err
	//}
	//
	//log.Info("RestWatch Server stopped gracefully")

	return nil
}
