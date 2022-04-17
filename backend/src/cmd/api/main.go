package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
)

const (
	ExitOK int = iota
	ExitError
)

const (
	defaultPort = "8080"
)

func startHttpServer() int {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to prepare logger: %s\n", err)
		return ExitError
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	addr := ":" + port
	srv := NewHttpServer(addr, logger)
	go func() {
		logger.Info("start server", zap.String("addr", addr))
		if err := srv.Start(ctx); err != nil {
			logger.Error("start server", zap.Error(err))
			os.Exit(ExitError)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logger.Info("stop server")
	if err := srv.Stop(ctx); err != nil {
		logger.Error("faield to stop server", zap.Error(err))
		return ExitError
	}

	logger.Info("successfully shutdown server")

	return ExitOK
}

func main() {
	os.Exit(startHttpServer())
}
