package main

import (
	"fmt"
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
	_, _ = w.Write([]byte("Hello World"))
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
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/404", pageMissing)
	mux.HandleFunc("/snippet/new", snippetNew)
	mux.HandleFunc("/snippet/view", snippetView)

	log.Println("starting server on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
