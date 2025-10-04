package main

import (
  "net/http"

	"snippetbox.leen2233.me/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)


func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	
  fileserver := http.FileServer(http.FS(ui.Files))
  router.Handler(http.MethodGet, "/static/*filepath",fileserver)
	
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

  router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
  router.HandlerFunc(http.MethodGet, "/ping", ping)
	router.Handler(http.MethodGet, "/view/:id", dynamic.ThenFunc(app.ViewSnippet))
  router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignupForm))
  router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
  router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLoginForm))
  router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
  router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))

  protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/create", protected.ThenFunc(app.CreateSnippetForm))
  router.Handler(http.MethodPost, "/create", protected.ThenFunc(app.CreateSnippetPost))
  router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
  router.Handler(http.MethodGet, "/profile", protected.ThenFunc(app.profile))
  router.Handler(http.MethodGet, "/profile/change-password", protected.ThenFunc(app.changePasswordForm))
  router.Handler(http.MethodPost, "/profile/change-password", protected.ThenFunc(app.changePasswordPost))

	standard := alice.New(app.recoverPanic, logRequest, secureHeaders)

  return standard.Then(router)
}

