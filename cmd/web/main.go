package main

import (
	"log"
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func main() {
	mux := http.NewServeMux()

	fileserver := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/new", snipNew)
	mux.HandleFunc("/snippet/view", snipView)
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	log.Println("starting the server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatalln(err)
}
