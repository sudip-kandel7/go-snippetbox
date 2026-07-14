package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"go-snippetbox.kandel.net/internal/models"
	"go-snippetbox.kandel.net/internal/validator"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	// FieldErrors map[string]string
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {   // becuase now matches exactly '/'
	// 	http.NotFound(w, r)
	// 	return
	// }

	// panic("OOPS! there is some problem!")

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }

	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	// app.errorLog.Println(err.Error())
	// 	// http.Error(w, "Internal Server Error", 500)
	// 	// return
	// 	app.serverError(w, err)
	// 	return
	// }

	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// err = ts.ExecuteTemplate(w, "base", data)

	// if err != nil {
	// 	// app.errorLog.Println(err.Error())
	// 	// http.Error(w, "Internal Server Error", 500)
	// 	// return
	// 	app.serverError(w, err)
	// }

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	// id, err := strconv.Atoi(r.URL.Query().Get("id")) // without httprouter used before
	id, err := strconv.Atoi(params.ByName("id")) // with httprouter

	if err != nil || id < 1 {
		// http.NotFound(w, r)
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

	app.render(w, http.StatusOK, "view.tmpl", data)

	// fmt.Fprintf(w, "%+v", snippet)
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// data := &templateData {
	// 	Snippet: snippet,
	// }

	// err = ts.ExecuteTemplate(w, "base", snippet)
	// err = ts.ExecuteTemplate(w, "base", data)

	// if err != nil {
	// 	app.serverError(w, err)
	// }

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
		// FieldErrors: map[string]string{},
	}

	// fieldErrors := make(map[string]string)

	// if strings.TrimSpace(form.Title) == "" {
	// 	fieldErrors["title"] = "This field cannot be blank"
	// } else if utf8.RuneCountInString(form.Title) > 100 {
	// 	fieldErrors["title"] = "This field cannot be more than 100 characters"
	// }

	// if strings.TrimSpace(form.Content) == "" {
	// 	fieldErrors["content"] = "This field cannot be blank"
	// }

	// if expires != 1 && expires != 7 && expires != 365 {
	// 	fieldErrors["expires"] = "This field must equal 1, 7, 365"
	// }

	// if len(fieldErrors) > 0 {
	// 	form.FieldErrors = fieldErrors
	// 	data := app.newTemplateData(r)
	// 	data.Form = form
	// 	app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
	// 	return
	// }

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
