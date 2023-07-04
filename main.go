package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	dir           string
	templateCache map[string]*template.Template

	infoLog *log.Logger
	errLog  *log.Logger

	sessionManager *scs.SessionManager

	riddle string
	answer string
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	content := flag.String("path", "/var/www/fog", "Path to vids")
	config := flag.String("config", "/etc/fog", "Path to configurations")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

	tc, err := loadTemplates()
	if err != nil {
		errLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = time.Hour * 24 * 365 * 5

	riddle, answer, err := loadRiddle(*config)
	if err != nil {
		errLog.Fatal(err)
	}

	app := &application{
		dir:            *content,
		templateCache:  tc,
		infoLog:        infoLog,
		errLog:         errLog,
		sessionManager: sessionManager,
		riddle:         riddle,
		answer:         answer,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("starting server on %s", *addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
