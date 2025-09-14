package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jaked0626/snippetbox/internal/db/models"
	"github.com/jaked0626/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

/* Home */

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.List(10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

/* Snippet View */

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// validate input
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.badRequest(w)
		return
	} else if id < 1 {
		app.notFound(w)
		return
	}

	// get from db
	s, err := app.snippets.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = s

	app.render(w, http.StatusOK, "view.html", data)
}

/* Snippet Create */

type snippetCreateForm struct {
	Title               string `form:"title"`
	Author              string `form:"author"`
	Work                string `form:"work"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (form *snippetCreateForm) validate() {
	// checkfield is a method of the embedded validator
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// initialize form data
	data.Form = &snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequest(w)
		return
	}

	// validate
	form.validate()
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	app.infoLog.Println("Creating following quote: ", form.Title, form.Author, form.Work, form.Content)

	// get user id from session
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userId == 0 {
		app.serverError(w, errors.New("unauthenticated user"))
	}

	id, err := app.snippets.Insert(form.Title, form.Author, form.Work, form.Content, form.Expires, userId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "toast", "Quote successfully submitted!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

/* Snippet List */

func (app *application) snippetList(w http.ResponseWriter, r *http.Request) {
	// validate input
	params := httprouter.ParamsFromContext(r.Context())
	limit, err := strconv.Atoi(params.ByName("limit"))
	if err != nil {
		app.badRequest(w)
		return
	} else if limit < 1 {
		app.notFound(w)
		return
	}

	// get from db
	snippets, err := app.snippets.List(limit)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// response
	res, err := json.Marshal(snippets)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

/* User Snippets */

func (app *application) snippetUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId, err := strconv.Atoi(params.ByName("userid"))
	if err != nil {
		app.badRequest(w)
		return
	} else if userId < 1 {
		app.notFound(w)
		return
	}

  snippets, err := app.snippets.ListUser(userId)
  if err != nil {
    app.serverError(w, err)
    return
  }

  data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetMine(w http.ResponseWriter, r *http.Request) {
  userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
  if userId <= 0 {
    app.serverError(w, errors.New("unauthenticated user"))
    return
  }

	http.Redirect(w, r, fmt.Sprintf("/snippet/user/%d", userId), http.StatusSeeOther)
}

/* User Signup */

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (form *userSignUpForm) validate() {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &userSignUpForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignUpForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequest(w)
		return
	}

	form.validate()
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "toast", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

/* User Log In */

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (form *userLoginForm) validate() {
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequest(w)
		return
	}

	form.validate()
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
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

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/* User Log Out */

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) { // Use the RenewToken() method on the current session to change the session // ID again.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "toast", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


