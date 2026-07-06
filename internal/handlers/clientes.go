package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pool-api/internal/models"

	"github.com/go-chi/chi/v5"
)

// ─── CLIENTE ─────────────────────────────────────────────────────────────────

func (s *Server) ListarClientes(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Clientes.ListarClientes())
}

func (s *Server) ObtenerCliente(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	c, ok := s.Clientes.ObtenerCliente(uint(idInt))
	if !ok {
		RespondError(w, http.StatusNotFound, "cliente no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, c)
}

func (s *Server) CrearCliente(w http.ResponseWriter, r *http.Request) {
	var c models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Clientes.CrearCliente(c)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	var c models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Clientes.ActualizarCliente(uint(idInt), c)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarCliente(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Clientes.BorrarCliente(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "cliente eliminado"})
}

// ─── RESERVA ─────────────────────────────────────────────────────────────────

func (s *Server) ListarReservas(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Clientes.ListarReservas())
}

func (s *Server) ObtenerReserva(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	rv, ok := s.Clientes.ObtenerReserva(uint(idInt))
	if !ok {
		RespondError(w, http.StatusNotFound, "reserva no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, rv)
}

func (s *Server) CrearReserva(w http.ResponseWriter, r *http.Request) {
	var rv models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Clientes.CrearReserva(rv)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ActualizarReserva(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	var rv models.Reserva
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Clientes.ActualizarReserva(uint(idInt), rv)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarReserva(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Clientes.BorrarReserva(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "reserva eliminada"})
}

// ─── PAGO ────────────────────────────────────────────────────────────────────

func (s *Server) ListarPagos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Clientes.ListarPagos())
}

func (s *Server) ObtenerPago(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	p, ok := s.Clientes.ObtenerPago(uint(idInt))
	if !ok {
		RespondError(w, http.StatusNotFound, "pago no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, p)
}

func (s *Server) CrearPago(w http.ResponseWriter, r *http.Request) {
	var p models.Pago
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Clientes.CrearPago(p)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ActualizarPago(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	var p models.Pago
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Clientes.ActualizarPago(uint(idInt), p)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarPago(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Clientes.BorrarPago(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "pago eliminado"})
}
