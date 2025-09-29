package main

import (
	"net/http"
	"testing"

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


