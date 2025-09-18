package main

import (
  "net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)


func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	
  fileserver := http.FileServer(http.Dir("./ui/static/"))
  router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))

  router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/view/:id", app.ViewSnippet)
	router.HandlerFunc(http.MethodGet, "/create", app.CreateSnippetForm)
  router.HandlerFunc(http.MethodPost, "/create", app.CreateSnippetPost)

	standard := alice.New(app.recoverPanic, logRequest, secureHeaders)

  return standard.Then(router)
}

