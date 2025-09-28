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

  router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/view/:id", dynamic.ThenFunc(app.ViewSnippet))
  router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignupForm))
  router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
  router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLoginForm))
  router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

  protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/create", protected.ThenFunc(app.CreateSnippetForm))
  router.Handler(http.MethodPost, "/create", protected.ThenFunc(app.CreateSnippetPost))
  router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, logRequest, secureHeaders)

  return standard.Then(router)
}

