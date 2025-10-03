package main

import (
  "database/sql"
  "flag"
  "log"
  "net/http"
  "html/template"
  "os"
  "time"
  "crypto/tls"

  "snippetbox.leen2233.me/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
  "github.com/alexedwards/scs/v2"

	"github.com/go-playground/form/v4"
  _ "github.com/go-sql-driver/mysql"
)


// application-wide dependencies
type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  snippets models.SnippetInterface
  users models.UserInterface
  cachedTemplates map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
  debugMode   bool
}


func main() {
  // command-line argument
  var addr, dsn *string
  var debug *bool
  addr = flag.String("addr", ":4000", "HTTP network address")
  dsn = flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
  debug = flag.Bool("debug", false, "In debug mode")

  flag.Parse()

  // setup custom logging
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  db, err := openDB(*dsn)
  if err != nil {
    errorLog.Fatal(err)
  }
  defer db.Close()

  cachedTemplates, err := newTemplateCache()
  if err != nil {
    errorLog.Fatal(err)
  }

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
  sessionManager.Cookie.Secure = true

  app := &application{
    errorLog: errorLog,
    infoLog: infoLog,
    snippets: &models.SnippetModel{DB: db},
    users: &models.UserModel{DB: db},
    cachedTemplates: cachedTemplates,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
    debugMode: *debug,
  }

  tlsConfig := &tls.Config{
    CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
  }

  // initilize custom Server
  srv := &http.Server{
    Addr:     *addr,
    ErrorLog: errorLog,
    Handler:  app.routes(),
    TLSConfig: tlsConfig,
    IdleTimeout: time.Minute,
    ReadTimeout: 5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  infoLog.Printf("Listening on port http://localhost%s  DEBUG = %v", *addr, *debug)
  err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
  errorLog.Fatal(err)
}


func openDB(dsn string) (*sql.DB, error) {
  var db *sql.DB
  var err error
  db, err = sql.Open("mysql", dsn)
  if err != nil {
    return nil, err
  }
  if err = db.Ping(); err != nil {
    return nil, err
  }

  return db, nil
}

