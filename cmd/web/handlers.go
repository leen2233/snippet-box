package main

import (
  "errors"
  "net/http"
  "strconv"
  "fmt"

  "snippetbox.leen2233.me/internal/models"
)


func (app *application) home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    app.notFound(w)
    return
  }

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
  id, err := strconv.Atoi(r.URL.Query().Get("id"))

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


func (app *application) CreateSnippet(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", "POST")
    app.clientError(w, 405)
    return
  }

  title := "Story"
  content := "Story about loser"
  expires := 7

  id, err := app.snippets.Insert(title, content, expires)
  if err != nil {
    app.serverError(w, err)
    return
  }
 
  http.Redirect(w, r, fmt.Sprintf("/view?id=%d", id), http.StatusSeeOther)
}

