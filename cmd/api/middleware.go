package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/internal/validator"
	"github.com/terdia/greenlight/src/users/entities"
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
				clients[ip] = &client{limiter: rate.NewLimiter(
					rate.Limit(app.config.Limiter.Rps),
					app.config.Limiter.Burst),
				}
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

func (app *application) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//log every request like so e.g. 172.18.0.1:60504 - HTTP/1.1 GET /snippet?id=4
		app.logger.PrintInfo("incoming request", map[string]string{
			"ip":     r.RemoteAddr,
			"proto":  r.Proto,
			"method": r.Method,
			"uri":    r.URL.RequestURI(),
		})

		next.ServeHTTP(rw, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, entities.AnonymousUser)
			next.ServeHTTP(rw, r)

			return
		}

		utils := app.registry.Services.SharedUtil

		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.InvalidAuthenticationTokenResponse(rw, r)

			return
		}

		token := parts[1]

		v := validator.New()

		v.Check(token != "", "token", "must be provided")
		v.Check(len(token) == 26, "token", "must be 26 bytes long")
		if !v.Valid() {
			utils.InvalidAuthenticationTokenResponse(rw, r)

			return
		}

		user, err := app.registry.Services.UserRepository.GetForToken(token, data.TokenScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				utils.InvalidAuthenticationTokenResponse(rw, r)
			default:
				utils.ServerErrorResponse(rw, r, err)
			}

			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(rw, r)

	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		utils := app.registry.Services.SharedUtil
		if user.IsAnonymous() {
			utils.AuthenticationRequiredResponse(rw, r)

			return
		}

		next.ServeHTTP(rw, r)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {

	fn := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		utils := app.registry.Services.SharedUtil
		if !user.Activated {
			utils.InactiveAccountResponse(rw, r)

			return
		}

		next.ServeHTTP(rw, r)
	})

	// Wrap fn with the requireAuthenticatedUser() middleware before returning it.
	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {

	fn := func(rw http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		utils := app.registry.Services.SharedUtil
		permissions, err := app.registry.Services.PermissionRepository.GetAllForUser(user.ID)

		if err != nil {
			utils.ServerErrorResponse(rw, r, err)

			return
		}

		if !permissions.Includes(code) {
			utils.NotPermittedRResponse(rw, r)

			return
		}

		next.ServeHTTP(rw, r)
	}

	return app.requireActivatedUser(fn)
}
