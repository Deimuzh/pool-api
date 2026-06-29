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

// ─── GUARDAVIDA ──────────────────────────────────────────────────────────────

// CrearGuardavida maneja POST /api/v1/guardavidas
func CrearGuardavida(w http.ResponseWriter, r *http.Request) {
	// Declaro una variable de tipo Guardavida para guardar los datos del request
	var g models.Guardavida

	// Decodifico el JSON del body y lo guardo en g
	// Si el JSON está mal formado, devuelvo 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validación básica: nombre y turno son obligatorios
	// Si vienen vacíos, devuelvo 400 porque es error del cliente
	if g.Nombre == "" || g.Turno == "" {
		http.Error(w, "nombre y turno son obligatorios", http.StatusBadRequest)
		return
	}

	// Asigno la fecha actual en el servidor, no la recibo del cliente
	// para que el timestamp sea confiable y no manipulable
	g.CreadoEn = time.Now()

	// Inserto el nuevo registro en la tabla guardavidas de la base de datos
	// Si falla la inserción, devuelvo 500 Internal Server Error
	if err := storage.DB.Create(&g).Error; err != nil {
		http.Error(w, "Error al guardar guardavida", http.StatusInternalServerError)
		return
	}

	// Indico al cliente que la respuesta es JSON
	w.Header().Set("Content-Type", "application/json")
	// Envío 201 Created - debe ir ANTES del Encode, si no Go lo ignora
	w.WriteHeader(http.StatusCreated)
	// Convierto el struct a JSON y lo envío como respuesta
	json.NewEncoder(w).Encode(g)
}

// ListarGuardavidas maneja GET /api/v1/guardavidas
func ListarGuardavidas(w http.ResponseWriter, r *http.Request) {
	// Declaro un slice para guardar todos los guardavidas
	var guardavidas []models.Guardavida

	// Busco todos los registros de la tabla guardavidas
	// Si falla la consulta, devuelvo 500
	if err := storage.DB.Find(&guardavidas).Error; err != nil {
		http.Error(w, "Error al obtener guardavidas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(guardavidas)
}

// ObtenerGuardavida maneja GET /api/v1/guardavidas/{id}
func ObtenerGuardavida(w http.ResponseWriter, r *http.Request) {
	// Extraigo el {id} de la URL - Chi lo devuelve como string
	idStr := chi.URLParam(r, "id")
	// Convierto el string a entero para buscar en la base de datos
	// Si no es un número válido, devuelvo 400
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var g models.Guardavida
	// Busco el guardavida por ID - si no existe, GORM devuelve error
	// y respondo 404 Not Found
	if err := storage.DB.First(&g, id).Error; err != nil {
		http.Error(w, "Guardavida no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
}

// ActualizarGuardavida maneja PATCH /api/v1/guardavidas/{id}
func ActualizarGuardavida(w http.ResponseWriter, r *http.Request) {
	// Extraigo y convierto el ID de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Verifico que el guardavida existe antes de actualizar
	// Si no existe, devuelvo 404
	var g models.Guardavida
	if err := storage.DB.First(&g, id).Error; err != nil {
		http.Error(w, "Guardavida no encontrado", http.StatusNotFound)
		return
	}

	// Decodifico los campos a actualizar del body
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Uso Updates en lugar de Save para solo modificar los campos que llegaron
	// Save sobreescribia todos los campos con valores zero si no venían en el JSON
	if err := storage.DB.Model(&g).Updates(&g).Error; err != nil {
		http.Error(w, "Error al actualizar guardavida", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
}

// EliminarGuardavida maneja DELETE /api/v1/guardavidas/{id}
func EliminarGuardavida(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Elimino el registro por ID de la tabla guardavidas
	// Si falla, devuelvo 500
	if err := storage.DB.Delete(&models.Guardavida{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar guardavida", http.StatusInternalServerError)
		return
	}

	// Respondo 200 con un mensaje JSON confirmando la eliminación
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"guardavida eliminado"}`))
}

// ─── INCIDENTE ───────────────────────────────────────────────────────────────

// CrearIncidente registra un nuevo incidente vinculado a un guardavida
func CrearIncidente(w http.ResponseWriter, r *http.Request) {
	var inc models.Incidente

	// Decodifico el JSON del body
	if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Valido que los campos obligatorios no vengan vacíos
	if inc.Tipo == "" || inc.Gravedad == "" || inc.GuardavidaID == 0 {
		http.Error(w, "tipo, gravedad y guardavida_id son obligatorios", http.StatusBadRequest)
		return
	}

	// Verifico que el guardavida_id exista en la base de datos
	// GORM no enforza foreign keys automáticamente en SQLite
	// así que lo valido manualmente para mantener integridad referencial
	var g models.Guardavida
	if err := storage.DB.First(&g, inc.GuardavidaID).Error; err != nil {
		http.Error(w, "guardavida_id no existe", http.StatusBadRequest)
		return
	}

	// Asigno la fecha actual en el servidor
	inc.FechaHora = time.Now()

	// Inserto el incidente en la base de datos
	if err := storage.DB.Create(&inc).Error; err != nil {
		http.Error(w, "Error al guardar incidente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(inc)
}

// ListarIncidentes retorna todos los incidentes registrados
func ListarIncidentes(w http.ResponseWriter, r *http.Request) {
	var incidentes []models.Incidente

	if err := storage.DB.Find(&incidentes).Error; err != nil {
		http.Error(w, "Error al obtener incidentes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incidentes)
}

// ObtenerIncidente maneja GET /api/v1/incidentes/{id}
func ObtenerIncidente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var inc models.Incidente
	if err := storage.DB.First(&inc, id).Error; err != nil {
		http.Error(w, "Incidente no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(inc)
}

// ActualizarIncidente maneja PATCH /api/v1/incidentes/{id}
func ActualizarIncidente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var inc models.Incidente
	if err := storage.DB.First(&inc, id).Error; err != nil {
		http.Error(w, "Incidente no encontrado", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Updates solo modifica los campos que llegaron en el JSON
	if err := storage.DB.Model(&inc).Updates(&inc).Error; err != nil {
		http.Error(w, "Error al actualizar incidente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(inc)
}

// EliminarIncidente maneja DELETE /api/v1/incidentes/{id}
func EliminarIncidente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Delete(&models.Incidente{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar incidente", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"incidente eliminado"}`))
}

// ─── ACCESO CLIENTE ──────────────────────────────────────────────────────────

// CrearAcceso registra el ingreso de un cliente a la piscina
func CrearAcceso(w http.ResponseWriter, r *http.Request) {
	var acc models.AccesoCliente

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Valido que cliente_id no sea cero
	if acc.ClienteID == 0 {
		http.Error(w, "cliente_id es obligatorio", http.StatusBadRequest)
		return
	}

	// Asigno la fecha actual del acceso en el servidor
	acc.FechaHora = time.Now()

	if err := storage.DB.Create(&acc).Error; err != nil {
		http.Error(w, "Error al registrar acceso", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acc)
}

// ListarAccesos retorna todos los registros de acceso a la piscina
func ListarAccesos(w http.ResponseWriter, r *http.Request) {
	var accesos []models.AccesoCliente

	if err := storage.DB.Find(&accesos).Error; err != nil {
		http.Error(w, "Error al obtener accesos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accesos)
}

// ObtenerAcceso retorna un registro de acceso por su ID
func ObtenerAcceso(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var acc models.AccesoCliente
	if err := storage.DB.First(&acc, id).Error; err != nil {
		http.Error(w, "Acceso no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acc)
}

// ActualizarAcceso maneja PATCH /api/v1/accesos/{id}
func ActualizarAcceso(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var acc models.AccesoCliente
	if err := storage.DB.First(&acc, id).Error; err != nil {
		http.Error(w, "Acceso no encontrado", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Updates solo modifica los campos que llegaron en el JSON
	if err := storage.DB.Model(&acc).Updates(&acc).Error; err != nil {
		http.Error(w, "Error al actualizar acceso", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acc)
}

// EliminarAcceso maneja DELETE /api/v1/accesos/{id}
func EliminarAcceso(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Delete(&models.AccesoCliente{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar acceso", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"acceso eliminado"}`))
}