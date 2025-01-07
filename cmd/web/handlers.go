package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	// serve 404 page when url path is not found
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func pageMissing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Page missing"))
}

func snippetNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// The api will work fine without it but it tells the user what method to use
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Snippet New\n"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}
	fmt.Fprintf(w, "Viewing snippet %d", id)
}
