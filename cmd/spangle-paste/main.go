package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	api "github.com/adoublef/spangle-paste/internal/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	q := make(chan os.Signal, 1)
	signal.Notify(q, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-q
		cancel()
	}()

	if err := run(ctx); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(ctx context.Context) error {
	mux := chi.NewMux()

	logger := newLogger(AppName)
	mux.Use(httplog.RequestLogger(logger))

	api := api.NewService()
	mux.Mount("/", api)

	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	e := make(chan error, 1)
	go func() {
		e <- srv.ListenAndServe()
	}()

	select {
	case err := <-e: //failed to start server
		return errors.Join(ErrStartServer, err)
	case <-ctx.Done(): //server shutdown
		if err := srv.Shutdown(ctx); err != nil {
			return errors.Join(ErrShutdownServer, err)
		}
		return nil //ctx.Err()
	}
}
