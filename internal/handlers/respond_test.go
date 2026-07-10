package handlers

import (
	"errors"
	"net/http"
	"testing"

	"pool-api/internal/service"
)

func TestStatusDeError(t *testing.T) {
	tests := []struct {
		err      error
		expected int
	}{
		{service.ErrCampoObligatorio, http.StatusBadRequest},
		{service.ErrNombreVacio, http.StatusBadRequest},
		{service.ErrMontoInvalido, http.StatusBadRequest},
		{service.ErrGuardavidaInvalido, http.StatusBadRequest},
		{service.ErrClienteInvalido, http.StatusBadRequest},
		{service.ErrClienteSinMembresia, http.StatusBadRequest},
		{service.ErrClienteConMembresia, http.StatusBadRequest},
		{service.ErrClienteSinAcceso, http.StatusBadRequest},
		{service.ErrConceptoPagoInvalido, http.StatusBadRequest},
		{service.ErrDuracionInvalida, http.StatusBadRequest},
		{service.ErrEquipoInvalido, http.StatusBadRequest},
		{service.ErrNoEncontrado, http.StatusNotFound},
		{service.ErrCedulaEnUso, http.StatusConflict},
		{service.ErrEmailEnUso, http.StatusConflict},
		{service.ErrCredencialesInvalidas, http.StatusUnauthorized},
		{errors.New("otro error"), http.StatusInternalServerError},
	}
	for _, tc := range tests {
		got := statusDeError(tc.err)
		if got != tc.expected {
			t.Errorf("statusDeError(%v) = %d, se esperaba %d", tc.err, got, tc.expected)
		}
	}
}
