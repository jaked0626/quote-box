package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jaked0626/snippetbox/internal/db/models"
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
	app.render(w, http.StatusOK, "create.html", data)
	return
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.badRequest(w)
	}
	title := r.PostForm.Get("title")
	// author := r.PostForm.Get("author")
	// work := r.PostForm.Get("work")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.badRequest(w)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	return
}
