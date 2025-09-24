package main

import (
  "errors"
  "net/http"
  "strconv"
  "fmt"
  "strings"
  "unicode/utf8"

  "github.com/julienschmidt/httprouter"
	"snippetbox.leen2233.me/internal/models"
)


type snippetCreateForm struct {
  Title       string
  Content     string
  Expires     int
  FieldErrors map[string]string
}


func (app *application) home(w http.ResponseWriter, r *http.Request) {
  snippets, err := app.snippets.Latest()
  if err != nil {
    app.serverError(w, err)
    return
  } 

  data := app.newTemplateData(r)
  data.Snippets = snippets

  app.render(w, 200, "home.tmpl", data)
}


func (app *application) ViewSnippet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
  if err != nil || id < 1 {
    app.notFound(w)
    return
  }

  snippet, err := app.snippets.Get(id)
  if err != nil {
    if errors.Is(err, models.ErrNoRecord) {
      app.notFound(w)
    } else {
      app.serverError(w, err)
    }
    return
  }

  data := app.newTemplateData(r)
  data.Snippet = snippet

  app.render(w, 200, "view.tmpl", data)
}


func (app *application) CreateSnippetForm(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  form := snippetCreateForm{}
  data.Form = form

	app.render(w, 200, "create.tmpl", data)
}

func (app *application) CreateSnippetPost(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()

  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  expires, err := strconv.Atoi(r.PostForm.Get("expires"))
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return 
  }

  form := snippetCreateForm{
    Title:        r.PostForm.Get("title"),
    Content:      r.PostForm.Get("content"),
    Expires:      expires,
    FieldErrors:  make(map[string]string),
  }

  if strings.TrimSpace(form.Title) == "" {
    form.FieldErrors["title"] = "This field cannot be null"
  } else if utf8.RuneCountInString(form.Title) > 100 {
    form.FieldErrors["title"] = "This field letter count can't exceed 100 characters'"
  }

  if strings.TrimSpace(form.Content) == "" {
    form.FieldErrors["content"] = "This field cannot be null"
  }

  if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
    form.FieldErrors["expires"] = "This field should be one of 1, 7, 365"
  }

  if len(form.FieldErrors) > 0 {
    data := app.newTemplateData(r)
    data.Form = form

    app.render(w, 400, "create.tmpl", data)
    return
  }

  id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
  if err != nil {
    app.serverError(w, err)
    return
  }
 
  http.Redirect(w, r, fmt.Sprintf("/view/%d", id), http.StatusSeeOther)
}

