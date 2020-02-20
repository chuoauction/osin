package osin

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo"
)

func TestInfo(t *testing.T) {
	sconfig := NewServerConfig()
	server := NewServer(sconfig, NewTestingStorage())
	resp := server.NewResponse()

	req, err := http.NewRequest("GET", "http://localhost:14000/appauth", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Form = make(url.Values)
	req.Form.Set("code", "9999")

	w := httptest.NewRecorder()
	c := echo.New().NewContext(req, w)
	if ar := server.HandleInfoRequest(resp, c); ar != nil {
		server.FinishInfoRequest(resp, c, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		t.Fatalf("Error in response: %s", resp.InternalError)
	}

	if resp.IsError {
		t.Fatalf("Should not be an error")
	}

	if resp.Type != DATA {
		t.Fatalf("Response should be data")
	}

	if d := resp.Output["access_token"]; d != "9999" {
		t.Fatalf("Unexpected authorization code: %s", d)
	}
}

func TestInfoWhenCodeIsOnHeader(t *testing.T) {
	sconfig := NewServerConfig()
	server := NewServer(sconfig, NewTestingStorage())
	resp := server.NewResponse()

	req, err := http.NewRequest("GET", "http://localhost:14000/appauth", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer 9999")

	w := httptest.NewRecorder()
	c := echo.New().NewContext(req, w)
	if ar := server.HandleInfoRequest(resp, c); ar != nil {
		server.FinishInfoRequest(resp, c, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		t.Fatalf("Error in response: %s", resp.InternalError)
	}

	if resp.IsError {
		t.Fatalf("Should not be an error")
	}

	if resp.Type != DATA {
		t.Fatalf("Response should be data")
	}

	if d := resp.Output["access_token"]; d != "9999" {
		t.Fatalf("Unexpected authorization code: %s", d)
	}
}
