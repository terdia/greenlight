package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.registry.Services.SharedUtil.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map once
	// every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent any rate limiter checks from happening while
			// the cleanup is taking place.
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()

		}
	}()

	utils := app.registry.Services.SharedUtil

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if app.config.Limiter.Enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				utils.ServerErrorResponse(rw, r, err)
				return
			}

			// Lock the mutex to prevent concurrent execution.
			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				utils.RateLimitExceededResponse(rw, r)

				return
			}

			mu.Unlock()

			next.ServeHTTP(rw, r)
		}
	})
}
