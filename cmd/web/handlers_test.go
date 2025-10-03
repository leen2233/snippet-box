package main

import (
	"net/http"
	"testing"
  "net/url"

	"snippetbox.leen2233.me/internal/assert"
)


func TestPing(t *testing.T) {
  app := newTestApplication(t)
  ts := newTestServer(t, app.routes())
  defer ts.Close()

  statusCode, _, body := ts.get(t, "/ping")

  assert.Equal(t, statusCode, http.StatusOK)
  assert.Equal(t, body, "ok")
}


func TestSnippetView(t *testing.T) {
  app := newTestApplication(t)
  ts := newTestServer(t, app.routes())
  defer ts.Close()

  // valid snippet
  statusCode, _, body := ts.get(t, "/view/1")
  assert.Equal(t, statusCode, http.StatusOK)
  assert.StringContains(t, body, "mocked data")

  // invalid id
  statusCode, _, _ = ts.get(t, "/view/2")
  assert.Equal(t, statusCode, http.StatusNotFound)

}


