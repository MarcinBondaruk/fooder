package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSwaggerApiHandler(t *testing.T) {
	handler := SwaggerApiHandler()

	req := httptest.NewRequest(http.MethodGet, "/docs/api", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", ct, "text/html; charset=utf-8")
	}

	if rec.Body.String() != SwaggerApiHTML {
		t.Errorf("body does not match SwaggerApiHTML constant")
	}
}

func TestRawApiHandler(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })

	if _, err := os.Stat("api/openapi/openapi.yaml"); os.IsNotExist(err) {
		t.Skip("api/openapi/openapi.yaml not found at project root")
	}

	handler := RawApiHandler()

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "text/yaml" {
		t.Errorf("Content-Type = %q, want %q", ct, "text/yaml")
	}

	if rec.Body.Len() == 0 {
		t.Error("expected non-empty body")
	}
}
