package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infolog,
		errorLog: errorlog,
	}

	mux := http.NewServeMux()

	fileserver := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/new", app.snipNew)
	mux.HandleFunc("/snippet/view", app.snipView)
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorlog,
		Handler:  mux,
	}

	infolog.Printf("starting the server %s", *addr)
	err := srv.ListenAndServe()
	errorlog.Fatalln(err)
}
