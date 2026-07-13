package main

import (
	"net/http"
	"github.com/justinas/alice"
)

// func (app *application) routes() *http.ServeMux{
func (app *application) routes() http.Handler{
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	standard := alice.New(app.panicRecovery, app.logRequest, secureHeaders)

	// return  app.panicRecovery(app.logRequest(secureHeaders(mux)))
	return standard.Then(mux)

}