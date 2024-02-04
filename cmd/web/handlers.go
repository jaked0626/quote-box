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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.List(10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
	return
}

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

	return
}

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
	return
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// initialize form data
	data.Form = &SnippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", data)
	return
}

type SnippetCreateForm struct {
	Title               string
	Author              string
	Work                string
	Content             string
	Expires             int
	validator.Validator // embedded
}

/* Validate form inputs */
func (form *SnippetCreateForm) validate() {
	// checkfield is a method of the embeeded validator
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	return
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.badRequest(w)
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.badRequest(w)
		return
	}

	author := r.PostForm.Get("author")
	if author == "" {
		author = "Unknown"
	}

	work := r.PostForm.Get("work")
	if work == "" {
		work = "Unknown"
	}

	form := &SnippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Author:  author,
		Work:    work,
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// validate
	form.validate()
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Author, form.Work, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "toast", "Quote successfully submitted!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	return
}
