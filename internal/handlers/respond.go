package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"pool-api/internal/service"
)

// RespondJSON escribe data como JSON con el status HTTP indicado.
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error codificando JSON: %v", err)
	}
}

// RespondError escribe un error en un formato JSON consistente: {"error": "..."}.
func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

// statusDeError traduce un error de la capa service al código HTTP correspondiente.
func statusDeError(err error) int {
	switch {
	case errors.Is(err, service.ErrNombreVacio),
		errors.Is(err, service.ErrCampoObligatorio),
		errors.Is(err, service.ErrMontoInvalido),
		errors.Is(err, service.ErrGuardavidaInvalido),
		errors.Is(err, service.ErrClienteInvalido),
		errors.Is(err, service.ErrEquipoInvalido):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, service.ErrCedulaEnUso),
		errors.Is(err, service.ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, service.ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
