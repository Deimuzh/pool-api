package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"piscina-comunitaria-api/internal/models"
	"piscina-comunitaria-api/internal/storage"

	"github.com/go-chi/chi/v5"
)

// CrearCliente maneja POST /api/v1/clientes

func CrearCliente(w http.ResponseWriter, r *http.Request) {
	var c models.Cliente

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if c.Nombre == "" || c.Cedula == "" {
		http.Error(w, "nombre y cedula son obligatorios", http.StatusBadRequest)
		return
	}

	if c.Membresia == "" {
		c.Membresia = "ninguna"
	}

	c.FechaRegistro = time.Now()

	if err := storage.DB.Create(&c).Error; err != nil {
		http.Error(w, "Error al guardar cliente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}
// ListarClientes obtiene y devuelve todos los clientes registrados en el sistema.

func ListarClientes(w http.ResponseWriter, r *http.Request) {
	var clientes []models.Cliente

	if err := storage.DB.Find(&clientes).Error; err != nil {
		http.Error(w, "Error al obtener clientes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clientes)
}

// ObtenerCliente busca y retorna la información de un cliente según su ID.

func ObtenerCliente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var c models.Cliente
	if err := storage.DB.First(&c, id).Error; err != nil {
		http.Error(w, "Cliente no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

// ActualizarCliente modifica los datos de un cliente existente.

func ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var c models.Cliente
	if err := storage.DB.First(&c, id).Error; err != nil {
		http.Error(w, "Cliente no encontrado", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Save(&c).Error; err != nil {
		http.Error(w, "Error al actualizar cliente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

// EliminarCliente elimina un cliente utilizando el ID proporcionado.

func EliminarCliente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Delete(&models.Cliente{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar cliente", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"cliente eliminado"}`))
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
