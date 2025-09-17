package main

import (
  "net/http"
)


func (app *application) routes() http.Handler {
  mux := http.NewServeMux()

  fileserver := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileserver))

  mux.HandleFunc("/", app.home)
  mux.HandleFunc("/view", app.ViewSnippet)
  mux.HandleFunc("/create", app.CreateSnippet)

  return logRequest(secureHeaders(mux))
}

