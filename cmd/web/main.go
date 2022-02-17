package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/gorilla/securecookie"
	"github.com/matthewlmitchell/tempshare/pkg/models"
	"github.com/matthewlmitchell/tempshare/pkg/models/mysql"
)

const version = "0.0.0001"

type config struct {
	port int    // For specifying port for the HTTP server to run on
	env  string // For launching server in development, staging, or production environment
	DB   struct {
		dsn                string
		maxOpenConnections int
		maxIdleConnections int
		maxIdleTime        string
	}
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	serverConfig  config
	templateCache map[string]*template.Template
	tempShare     interface {
		New(string, string, string) (*models.TempShare, error)
		Insert([]byte, string, string, int) error
		Get(string) (*models.TempShare, error)
		Delete(*models.TempShare) error
	}
}

func connectToDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verify that our connection to the database is alive
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	var servConfig config

	flag.IntVar(&servConfig.port, "port", 4000, "HTTP network address")
	flag.StringVar(&servConfig.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&servConfig.DB.dsn, "db-dsn", os.Getenv("TEMPSHARE_DSN"), "Specifies the MySQL database data source name (dsn)")
	flag.StringVar(&servConfig.DB.maxIdleTime, "db-max-idle-time", "5m", "MySQL maximum time allowed for an idle connection")
	flag.IntVar(&servConfig.DB.maxIdleConnections, "db-max-idle-conns", 25, "MySQL maximum number of idle connections")
	flag.IntVar(&servConfig.DB.maxOpenConnections, "db-max-open-conns", 25, "MySQL maximum number of open connections")

	// Generate a 32-bit key for securing our cookie session store
	secret := flag.String("secret", string(securecookie.GenerateRandomKey(32)), "Cookie store session secret")

	// We must parse all command line arguments before they can be used
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	db, err := connectToDatabase(servConfig.DB.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

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
		tempShare:     &mysql.TempShareModel{DB: db},
	}

	if err := app.initializeServer(); err != nil {
		app.errorLog.Fatalln(err)
	}
}
