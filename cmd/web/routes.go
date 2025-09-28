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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

  router.HandlerFunc(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.HandlerFunc(http.MethodGet, "/view/:id", dynamic.ThenFunc(app.ViewSnippet))
	router.HandlerFunc(http.MethodGet, "/create", dynamic.ThenFunc(app.CreateSnippetForm))
  router.HandlerFunc(http.MethodPost, "/create", dynamic.ThenFunc(app.CreateSnippetPost))

	standard := alice.New(app.recoverPanic, logRequest, secureHeaders)

  return standard.Then(router)
}

