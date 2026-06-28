package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"pool-api/internal/models"
	"pool-api/internal/storage"

	"github.com/go-chi/chi/v5"
)

// CLIENTE

func (s *Server) ListarClientes(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Clientes.ListarClientes())
}

func (s *Server) CrearCliente(w http.ResponseWriter, r *http.Request) {
	var c models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	creado, err := s.Clientes.CrearCliente(c)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerCliente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	cliente, ok := s.Clientes.ObtenerCliente(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "Cliente no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, cliente)
}

func (s *Server) ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var c models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	actualizado, err := s.Clientes.ActualizarCliente(id, c)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarCliente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Clientes.BorrarCliente(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "cliente eliminado"})
}

func CrearCliente(w http.ResponseWriter, r *http.Request) {
	panic("use Server.CrearCliente")
}
func ListarClientes(w http.ResponseWriter, r *http.Request) {
	panic("use Server.ListarClientes")
}
func ObtenerCliente(w http.ResponseWriter, r *http.Request) {
	panic("use Server.ObtenerCliente")
}
func ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	panic("use Server.ActualizarCliente")
}
func EliminarCliente(w http.ResponseWriter, r *http.Request) {
	panic("use Server.BorrarCliente")
}

// ─── RESERVA ─────────────────────────────────────────────────────────────────

func CrearReserva(w http.ResponseWriter, r *http.Request) {
	var rv models.Reserva

	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if rv.ClienteID == 0 {
		http.Error(w, "cliente_id es obligatorio", http.StatusBadRequest)
		return
	}

	if rv.Estado == "" {
		rv.Estado = "pendiente"
	}

	rv.FechaHora = time.Now()

	if err := storage.DB.Create(&rv).Error; err != nil {
		http.Error(w, "Error al guardar reserva", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rv)
}

func ListarReservas(w http.ResponseWriter, r *http.Request) {
	var reservas []models.Reserva

	if err := storage.DB.Find(&reservas).Error; err != nil {
		http.Error(w, "Error al obtener reservas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservas)
}

func ObtenerReserva(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var rv models.Reserva
	if err := storage.DB.First(&rv, id).Error; err != nil {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rv)
}

func ActualizarReserva(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var rv models.Reserva
	if err := storage.DB.First(&rv, id).Error; err != nil {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Save(&rv).Error; err != nil {
		http.Error(w, "Error al actualizar reserva", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rv)
}

func EliminarReserva(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Delete(&models.Reserva{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar reserva", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"reserva eliminada"}`))
}

// ─── PAGO ────────────────────────────────────────────────────────────────────

func CrearPago(w http.ResponseWriter, r *http.Request) {
	var p models.Pago

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if p.ClienteID == 0 || p.Monto <= 0 {
		http.Error(w, "cliente_id y monto son obligatorios", http.StatusBadRequest)
		return
	}

	p.FechaHora = time.Now()

	if err := storage.DB.Create(&p).Error; err != nil {
		http.Error(w, "Error al guardar pago", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func ListarPagos(w http.ResponseWriter, r *http.Request) {
	var pagos []models.Pago

	if err := storage.DB.Find(&pagos).Error; err != nil {
		http.Error(w, "Error al obtener pagos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pagos)
}

func ObtenerPago(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p models.Pago
	if err := storage.DB.First(&p, id).Error; err != nil {
		http.Error(w, "Pago no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func ActualizarPago(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p models.Pago
	if err := storage.DB.First(&p, id).Error; err != nil {
		http.Error(w, "Pago no encontrado", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Save(&p).Error; err != nil {
		http.Error(w, "Error al actualizar pago", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func EliminarPago(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Delete(&models.Pago{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar pago", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"pago eliminado"}`))
}
