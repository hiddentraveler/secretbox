package main

import (
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := "0.0.0.0" + ":" + port

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infolog,
		errorLog: errorlog,
	}

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorlog,
		Handler:  app.Routes(),
	}

	infolog.Printf("starting the server %s", port)
	err := srv.ListenAndServe()
	errorlog.Fatalln(err)
}
