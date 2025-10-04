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

type userSignupForm struct {
  Name        string   `form:"name"`
  Email       string   `form:"email"`
  Password    string   `form:"password"`
  validator.Validator  `form:-`
}

type userLoginForm struct {
  Email       string   `form:"email"`
  Password    string   `form:"password"`
  validator.Validator  `form:-`
}

type ChangePasswordForm struct {
  CurrentPassword     string `form:"currentPassword"`
  NewPassword         string `form:"newPassword"`
  PasswordConfirm     string `form:"passwordConfirm"`
  validator.Validator        `form:-`
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
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

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


// authentication

func (app *application) userSignupForm(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  form := userSignupForm{}
  data.Form = form

  app.render(w, 200, "signup.tmpl", data)
}


func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) { 
  var form userSignupForm

  err := app.decodePostForm(r, &form)
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
  form.CheckField(validator.MaxChars(form.Name, 50), "name", "This field cannot be more than 100 characters long")
  form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
  form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be valid email address")
  form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
  form.CheckField(validator.MinChars(form.Password, 8), "password", "This field cannot be less than 8 characters")

  if !form.Valid() {
    data := app.newTemplateData(r)
    data.Form = form

    app.render(w, 400, "signup.tmpl", data)
    return
  }

  err = app.users.Insert(form.Name, form.Email, form.Password)
  if err != nil {
    if errors.Is(err, models.ErrDuplicateEmail) {
      form.AddFieldError("email", "Email is already in use")
      
      data := app.newTemplateData(r)
      data.Form = form
      app.render(w, 400, "signup.tmpl", data)
    } else {
      app.serverError(w, err)
    }

    return
  }

  app.sessionManager.Put(r.Context(), "flash", "Your signup was successfull. Please log in.")
  http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r) 
  form := userLoginForm{}  
  data.Form = form

  app.render(w, 200, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
  var form userLoginForm

  err := app.decodePostForm(r, &form)
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  nextUrl := r.URL.Query().Get("next")
  if nextUrl == "" {
    nextUrl = "/"
  }

  id, err := app.users.Authenticate(form.Email, form.Password) 
  if err != nil {
    if errors.Is(err, models.ErrInvalidCredentials) {
      form.AddNonFieldError("Invalid Credentials")

      data := app.newTemplateData(r)
      data.Form = form
      app.render(w, 400, "login.tmpl", data)
    } else {
      app.serverError(w, err)
    }

    return
  }

  err = app.sessionManager.RenewToken(r.Context())
  if err != nil {
    app.serverError(w, err)
    return
  }

  // successful login
  app.sessionManager.Put(r.Context(), "authenticatedUserId", id)
  app.sessionManager.Put(r.Context(), "flash", "Login Successfull!")
  http.Redirect(w, r, nextUrl, http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
  err := app.sessionManager.RenewToken(r.Context())
  if err != nil {
    app.serverError(w, err)
    return
  }

  app.sessionManager.Remove(r.Context(), "authenticatedUserId")
  app.sessionManager.Put(r.Context(), "flash", "Logout successfull!")
  http.Redirect(w, r, "/", http.StatusSeeOther)
}


func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}


func (app *application) about(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  app.render(w, 200, "about.tmpl", data)
}


func (app *application) profile (w http.ResponseWriter, r *http.Request) {
  userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")

  user, err := app.users.Get(userId)
  if err != nil {
    app.serverError(w, err)
    return
  }

  data := app.newTemplateData(r)
  data.User = user

  app.render(w, 200, "profile.tmpl", data)
}


func (app *application) changePasswordForm (w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  form := ChangePasswordForm{}
  data.Form = form

  app.render(w, 200, "change_password.tmpl", data)
}


func (app *application) changePasswordPost (w http.ResponseWriter, r *http.Request) {
  var form ChangePasswordForm

  err := app.decodePostForm(r, &form)
  if err != nil {
    app.serverError(w, err)
  }

  form.CheckField(validator.MinChars(form.CurrentPassword, 8), "currentPassword", "Current Password should be at least 8 characters") 
  form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "New Password should be at least 8 characters")
  form.CheckField(validator.Equal(form.NewPassword, form.PasswordConfirm), "newPassword", "Passwords don't match'")
  form.CheckField(validator.Equal(form.NewPassword, form.PasswordConfirm), "passwordConfirm", "Passwords don't match'")

  if !form.Valid() {
    data := app.newTemplateData(r)
    data.Form = form

    app.render(w, 400, "change_password.tmpl", data)
  }

  userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
  err = app.users.ChangePassword(userId, form.CurrentPassword, form.NewPassword)
  if err != nil {
    if errors.Is(models.ErrInvalidCurrentPassword, err) {
      form.AddFieldError("currentPassword", "Current Password is invalid")

      data := app.newTemplateData(r)
      data.Form = form
      app.render(w, 400, "change_password.tmpl", data)
    } else {
      app.serverError(w, err)
      return
    }
  }

  app.sessionManager.Put(r.Context(), "flash", "Password is successfully changed")
  http.Redirect(w, r, "/", http.StatusSeeOther)  
}

