package main

import (
	"context"
  "net/http"
  "fmt"
  "time"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

    w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "deny")
    w.Header().Set("X-XSS-Protection", "0")

    next.ServeHTTP(w, r)
  })
}


func logRequest(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%v - %s - %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)

    next.ServeHTTP(w, r)
  })
}


func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err:= recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}


func (app *application) requireAuthentication(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if !app.isAuthenticated(r) {
      http.Redirect(w, r, "/user/login", http.StatusSeeOther)
    }

    w.Header().Set("Cache-Control", "no-store")
    next.ServeHTTP(w, r)
  })
}


func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:"/",
		Secure:true,
	})
	return csrfHandler
}


func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
		if userId == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.users.Exists(userId)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

