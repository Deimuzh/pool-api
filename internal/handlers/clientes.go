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

// ─── CLIENTE ─────────────────────────────────────────────────────────────────

// CrearCliente lee JSON desde el body y crea un nuevo cliente en la base.
// Request: JSON con campos de `models.Cliente` (Nombre, Cedula, Membresia opcional).
// Responses: 201 Created con el cliente creado; 400 Bad Request en JSON inválido
// o campos requeridos; 500 Internal Server Error si falla la persistencia.
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

func ListarClientes(w http.ResponseWriter, r *http.Request) {
	var clientes []models.Cliente
	// ListarClientes devuelve todos los clientes en JSON.
	// Responses: 200 OK con slice de clientes; 500 si hay error en la consulta.

	if err := storage.DB.Find(&clientes).Error; err != nil {
		http.Error(w, "Error al obtener clientes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clientes)
}

func ObtenerCliente(w http.ResponseWriter, r *http.Request) {
	// ObtenerCliente obtiene un cliente por ID (param URL `id`).
	// Responses: 200 OK con el cliente; 400 si el id es inválido; 404 si no existe.
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

func ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	// ActualizarCliente actualiza los campos de un cliente existente.
	// Request: JSON con los campos a actualizar.
	// Responses: 200 OK con el cliente actualizado; 400 en JSON inválido o id inválido;
	// 404 si el cliente no existe; 500 en error de persistencia.
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

func EliminarCliente(w http.ResponseWriter, r *http.Request) {
	// EliminarCliente borra un cliente por ID.
	// Responses: 200 OK en caso de eliminación; 400 si id inválido; 500 si falla la eliminación.
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

// CrearReserva crea una nueva reserva asociada a un cliente.
// Request: JSON con `ClienteID` (obligatorio) y otros campos.
// Responses: 201 Created; 400 si faltan datos; 500 si error al persistir.
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

// ListarReservas devuelve todas las reservas en JSON.
// Responses: 200 OK; 500 si falla la consulta.
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

// ObtenerReserva obtiene una reserva por ID (param URL `id`).
// Responses: 200 OK; 400 si id inválido; 404 si no encontrada.
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

// ActualizarReserva actualiza una reserva existente.
// Responses: 200 OK; 400 si JSON inválido o id inválido; 404 si no existe; 500 si falla.
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

// EliminarReserva borra una reserva por ID.
// Responses: 200 OK; 400 si id inválido; 500 si falla la eliminación.
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

// CrearPago crea un registro de pago para un cliente.
// Request: JSON con `ClienteID` y `Monto` (obligatorios).
// Responses: 201 Created; 400 si faltan datos o monto inválido; 500 si falla persistencia.
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

// ListarPagos devuelve todos los pagos en JSON.
// Responses: 200 OK; 500 si falla la consulta.
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

// ObtenerPago obtiene un pago por ID (param URL `id`).
// Responses: 200 OK; 400 si id inválido; 404 si no existe.
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

// ActualizarPago actualiza un pago existente.
// Responses: 200 OK; 400 si JSON inválido o id inválido; 404 si no existe; 500 si falla.
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

// EliminarPago borra un pago por ID.
// Responses: 200 OK; 400 si id inválido; 500 si falla la eliminación.
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
