package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"synchro/internal/models"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *slog.Logger
	events *models.EventModel
	users *models.UserModel
	templateCache map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {


	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	// sessionManager.Store = **** make it so session stored in sql database rather than in memory
	sessionManager.Lifetime = 12 * time.Hour
	// sessionManager.Cookie.Secure = true

	app := &application{
		logger: logger,
		events: &models.EventModel{DB: db},
		users: &models.UserModel{DB: db},
		templateCache: templateCache,
		sessionManager: sessionManager,
	}

	routes := app.routes()

	logger.Info("starting server", slog.String("addr", ":4000"))


	err = http.ListenAndServe(":4000", routes)
	logger.Error(err.Error())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}