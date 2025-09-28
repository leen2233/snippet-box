package main

import (
  "fmt"
  "bytes"
  "net/http"
  "runtime/debug"
  "time"
	"errors"

	"github.com/go-playground/form/v4"
)


func (app *application) serverError(w http.ResponseWriter, err error) {
  trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
  app.errorLog.Output(2, trace)

  http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}


func (app *application) clientError(w http.ResponseWriter, statusCode int) {
  http.Error(w, http.StatusText(statusCode), statusCode)
}


func (app *application) notFound(w http.ResponseWriter) {
  app.clientError(w, http.StatusNotFound)
}


func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
  ts, ok := app.cachedTemplates[page]
  if !ok {
    err := fmt.Errorf("the template %s does not exist", page)
    app.serverError(w, err)
    return
  }

  buf := new(bytes.Buffer)
 
  err := ts.ExecuteTemplate(buf, "base", data)
  if err != nil {
    app.serverError(w, err)
  }

  w.WriteHeader(status)
  buf.WriteTo(w)
}


func (app *application) newTemplateData(r *http.Request) *templateData {
  return &templateData{
    CurrentYear: time.Now().Year(),
    Flash: app.sessionManager.PopString(r.Context(), "flash"),
    IsAuthenticated : app.isAuthenticated(r),
  }
}


func (app *application) decodePostForm(r *http.Request, dest any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dest, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
  return app.sessionManager.Exists(r.Context(), "authenticatedUserId")
}

