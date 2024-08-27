package main

import "net/http"

func (app *application) Routes() *http.ServeMux {

	mux := http.NewServeMux()

	fileserver := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/new", app.snipNew)
	mux.HandleFunc("/snippet/view", app.snipView)
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	return mux
}
