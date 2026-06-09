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

// ─── GUARDAVIDA ──────────────────────────────────────────────────────────────

// CrearGuardavida maneja POST /api/v1/guardavidas
func CrearGuardavida(w http.ResponseWriter, r *http.Request) {
	var g models.Guardavida

	// Decodificar el JSON del body
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validación básica de campos requeridos
	if g.Nombre == "" || g.Turno == "" {
		http.Error(w, "nombre y turno son obligatorios", http.StatusBadRequest)
		return
	}

	g.CreadoEn = time.Now()

	if err := storage.DB.Create(&g).Error; err != nil {
		http.Error(w, "Error al guardar guardavida", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(g)
}

// ListarGuardavidas maneja GET /api/v1/guardavidas
func ListarGuardavidas(w http.ResponseWriter, r *http.Request) {
	var guardavidas []models.Guardavida

	if err := storage.DB.Find(&guardavidas).Error; err != nil {
		http.Error(w, "Error al obtener guardavidas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(guardavidas)
}

// ObtenerGuardavida maneja GET /api/v1/guardavidas/{id}
func ObtenerGuardavida(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var g models.Guardavida
	if err := storage.DB.First(&g, id).Error; err != nil {
		http.Error(w, "Guardavida no encontrado", http.StatusNotFound) // 404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
}

// ActualizarGuardavida maneja PATCH /api/v1/guardavidas/{id}
func ActualizarGuardavida(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Verificar que existe
	var g models.Guardavida
	if err := storage.DB.First(&g, id).Error; err != nil {
		http.Error(w, "Guardavida no encontrado", http.StatusNotFound)
		return
	}

	// Decodificar los campos a actualizar
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Save(&g).Error; err != nil {
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

	if err := storage.DB.Delete(&models.Guardavida{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar guardavida", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"guardavida eliminado"}`))
}

// ─── INCIDENTE ───────────────────────────────────────────────────────────────

// CrearIncidente maneja POST /api/v1/incidentes
func CrearIncidente(w http.ResponseWriter, r *http.Request) {
	var inc models.Incidente

	if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if inc.Tipo == "" || inc.Gravedad == "" || inc.GuardavidaID == 0 {
		http.Error(w, "tipo, gravedad y guardavida_id son obligatorios", http.StatusBadRequest)
		return
	}

	inc.FechaHora = time.Now()

	if err := storage.DB.Create(&inc).Error; err != nil {
		http.Error(w, "Error al guardar incidente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inc)
}

// ListarIncidentes maneja GET /api/v1/incidentes
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

	if err := storage.DB.Save(&inc).Error; err != nil {
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

// CrearAcceso maneja POST /api/v1/accesos
func CrearAcceso(w http.ResponseWriter, r *http.Request) {
	var acc models.AccesoCliente

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if acc.ClienteID == 0 {
		http.Error(w, "cliente_id es obligatorio", http.StatusBadRequest)
		return
	}

	acc.FechaHora = time.Now()

	if err := storage.DB.Create(&acc).Error; err != nil {
		http.Error(w, "Error al registrar acceso", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acc)
}

// ListarAccesos maneja GET /api/v1/accesos
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

// ObtenerAcceso maneja GET /api/v1/accesos/{id}
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

	if err := storage.DB.Save(&acc).Error; err != nil {
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
