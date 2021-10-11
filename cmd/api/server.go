package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	gracePeriod = 5 * time.Second
)

func (app *application) serve() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.AppPort),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	// Start a background goroutine for implementing graceful shutdown mechanism.
	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// 5 seconds grace period before shutdown
		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		app.wg.Wait() // wait for background go routines to finish before shuting down
		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.Env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("graceful shoutdown of server completed", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
