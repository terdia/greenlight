package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {

	app := newTestApplication(t, 0)
	ts := newTestServer(t, app.routes())

	rs := ts.get(t, "/v1/healthcheck")

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	body := struct {
		Status string
		Data   map[string]map[string]string
	}{}

	_ = json.Unmarshal(rs.Body, &body)

	if body.Status != "success" {
		t.Errorf("want %s; got %s", "success", body.Status)
	}

	if body.Data["system_info"]["environment"] != "test" {
		t.Errorf("want %s; got %s", "test", body.Data["system_info"]["environment"])
	}

	if body.Data["system_info"]["version"] != "v1" {
		t.Errorf("want %s; got %s", "v1", body.Data["system_info"][version])
	}

}
