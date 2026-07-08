package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"pool-api/internal/models"

	"github.com/go-chi/chi/v5"
)

func (s *Server) ListarReservas(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Clientes.ListarReservas())
}

func (s *Server) CrearReserva(w http.ResponseWriter, r *http.Request) {
	var rv models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creada, err := s.Clientes.CrearReserva(rv)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

func (s *Server) ObtenerReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	rv, ok := s.Clientes.ObtenerReserva(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "Reserva no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, rv)
}

func (s *Server) ActualizarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var rv models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	actualizada, err := s.Clientes.ActualizarReserva(id, rv)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

func (s *Server) BorrarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Clientes.BorrarReserva(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "reserva eliminada"})
}

func (s *Server) ListarPagos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Clientes.ListarPagos())
}

func (s *Server) CrearPago(w http.ResponseWriter, r *http.Request) {
	var p models.Pago
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado, err := s.Clientes.CrearPago(p)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerPago(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	p, ok := s.Clientes.ObtenerPago(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "Pago no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, p)
}

func (s *Server) ActualizarPago(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var p models.Pago
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	actualizado, err := s.Clientes.ActualizarPago(id, p)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarPago(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Clientes.BorrarPago(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "pago eliminado"})
}

var _ = time.Now
