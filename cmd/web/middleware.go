package main

import (
  "net/http"
  "fmt"
  "time"
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

