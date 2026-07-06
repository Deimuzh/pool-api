package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSoloRoles_AdminPermitido(t *testing.T) {
	handler := SoloRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	ctx := context.WithValue(req.Context(), ContextKeyRol, "admin")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d", rec.Code)
	}
}

func TestSoloRoles_RolNoPermitidoDevuelve403(t *testing.T) {
	handler := SoloRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	ctx := context.WithValue(req.Context(), ContextKeyRol, "guardavida")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("se esperaba 403, se obtuvo %d", rec.Code)
	}
}
