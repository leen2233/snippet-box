package main

import (
  "database/sql"
  "flag"
  "log"
  "net/http"
  "html/template"
  "os"

  "snippetbox.leen2233.me/internal/models"

  _ "github.com/go-sql-driver/mysql"
)


// application-wide dependencies
type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  snippets *models.SnippetModel
  cachedTemplates map[string]*template.Template
}


func main() {
  // command-line argument
  var addr, dsn *string
  addr = flag.String("addr", ":4000", "HTTP network address")
  dsn = flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

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

  app := &application{
    errorLog: errorLog,
    infoLog: infoLog,
    snippets: &models.SnippetModel{DB: db},
    cachedTemplates: cachedTemplates,
  }

  // initilize custom Server
  srv := &http.Server{
    Addr:     *addr,
    ErrorLog: errorLog,
    Handler:  app.routes(),
  }

  infoLog.Printf("Listening on port http://localhost%s", *addr)
  err = srv.ListenAndServe()
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

