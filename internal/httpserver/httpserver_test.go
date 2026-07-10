package httpserver

import (
	"net/http"
	"testing"
	"time"
)

func TestNuevo(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := Nuevo(handler, ":9090", 5*time.Second, 10*time.Second)

	if srv.Addr != ":9090" {
		t.Errorf("Addr esperado :9090, obtuve %s", srv.Addr)
	}
	if srv.ReadTimeout != 5*time.Second {
		t.Errorf("ReadTimeout inesperado: %v", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 10*time.Second {
		t.Errorf("WriteTimeout inesperado: %v", srv.WriteTimeout)
	}
	if srv.Handler == nil {
		t.Error("Handler no deberia ser nil")
	}
}
