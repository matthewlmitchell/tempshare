package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/gorilla/securecookie"
)

const version = "0.0.0001"

type config struct {
	port int    // For specifying port for the HTTP server to run on
	env  string // For launching server in development, staging, or production environment
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	serverConfig  config
	templateCache map[string]*template.Template
}

func main() {

	var servConfig config

	flag.IntVar(&servConfig.port, "port", 4000, "HTTP network address")
	flag.StringVar(&servConfig.env, "env", "development", "Environment (development|staging|production)")
	secret := flag.String("secret", string(securecookie.GenerateRandomKey(32)), "Cookie store session secret")

	// We must parse all command line arguments before they can be used
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	templateCache, err := initTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteLaxMode

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		serverConfig:  servConfig,
		templateCache: templateCache,
	}

	if err := app.initializeServer(); err != nil {
		app.errorLog.Fatalln(err)
	}
}
