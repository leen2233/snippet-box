package main

import (
  "errors"
  "net/http"
  "strconv"
  "fmt"

  "github.com/julienschmidt/httprouter"
	"snippetbox.leen2233.me/internal/models"
	"snippetbox.leen2233.me/internal/validator"
)


type snippetCreateForm struct {
	Title       string   `form:"title"`
	Content     string   `form:"content"`
	Expires     int      `form:"expires"`
	validator.Validator  `form:"-"`
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
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

  if !form.Valid() {
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

  app.sessionManager.Put(r.Context(), "flash", "Snippet Successfully created!")

  http.Redirect(w, r, fmt.Sprintf("/view/%d", id), http.StatusSeeOther)
}

