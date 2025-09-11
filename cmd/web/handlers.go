package main

import (
  "fmt"
  "net/http"
  "strconv"
  "html/template"
  "log"
)


func home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  files := []string{
    "./ui/html/base.tmpl",
    "./ui/html/partials/nav.tmpl",
    "./ui/html/pages/home.tmpl",
  }

  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    http.Error(w, "Internal sever error", 500)
    return
  }

  err = ts.ExecuteTemplate(w, "base", nil)
  if err != nil {
    log.Println(err.Error())
    http.Error(w, "Internal server error", 500)
  }
}


func ViewSnippet(w http.ResponseWriter, r *http.Request) {
  id, err := strconv.Atoi(r.URL.Query().Get("id"))

  if err != nil || id < 1 {
    http.NotFound(w, r)
    return
  }

  fmt.Fprintf(w, "View snippet with id: %d", id)
}


func CreateSnippet(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", "POST")
    http.Error(w, "Method not allowed", 405)
    return
  }

  w.Write([]byte("Creating new snippet"))
}

