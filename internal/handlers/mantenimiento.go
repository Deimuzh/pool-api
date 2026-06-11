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

func CrearEquipo(w http.ResponseWriter, r *http.Request) {
	var eq models.Equipo

	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if eq.Nombre == "" || eq.Tipo == "" {
		http.Error(w, "nombre y tipo son obligatorios", http.StatusBadRequest)
		return
	}

	if eq.Estado == "" {
		eq.Estado = "operativo"
	}

	if err := storage.DB.Create(&eq).Error; err != nil {
		http.Error(w, "Error al guardar equipo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(eq)
}

// ListarEquipos retorna todos los equipos registrados
func ListarEquipos(w http.ResponseWriter, r *http.Request) {
	var equipos []models.Equipo

	if err := storage.DB.Find(&equipos).Error; err != nil {
		http.Error(w, "Error al obtener equipos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(equipos)
}

func ObtenerEquipo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var eq models.Equipo
	if err := storage.DB.First(&eq, id).Error; err != nil {
		http.Error(w, "Equipo no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(eq)
}

// ActualizarEquipo maneja PATCH /api/v1/equipos/{id}
func ActualizarEquipo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var eq models.Equipo
	if err := storage.DB.First(&eq, id).Error; err != nil {
		http.Error(w, "Equipo no encontrado", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Save(&eq).Error; err != nil {
		http.Error(w, "Error al actualizar equipo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(eq)
}

func EliminarEquipo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Delete(&models.Equipo{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar equipo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"equipo eliminado"}`))
}

// ─── REGISTRO MANTENIMIENTO ──────────────────────────────────────────────────

// CrearRegistroMantenimiento maneja POST /api/v1/mantenimientos
func CrearRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	var rm models.RegistroMantenimiento

	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if rm.EquipoID == 0 || rm.Tipo == "" {
		http.Error(w, "equipo_id y tipo son obligatorios", http.StatusBadRequest)
		return
	}

	rm.FechaHora = time.Now()

	if err := storage.DB.Create(&rm).Error; err != nil {
		http.Error(w, "Error al guardar registro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rm)
}

// ListarRegistrosMantenimiento retorna todos los registros
func ListarRegistrosMantenimiento(w http.ResponseWriter, r *http.Request) {
	var registros []models.RegistroMantenimiento

	if err := storage.DB.Find(&registros).Error; err != nil {
		http.Error(w, "Error al obtener registros", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(registros)
}

func ObtenerRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var rm models.RegistroMantenimiento
	if err := storage.DB.First(&rm, id).Error; err != nil {
		http.Error(w, "Registro no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rm)
}

// ActualizarRegistroMantenimiento maneja PATCH /api/v1/mantenimientos/{id}
func ActualizarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var rm models.RegistroMantenimiento
	if err := storage.DB.First(&rm, id).Error; err != nil {
		http.Error(w, "Registro no encontrado", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Save(&rm).Error; err != nil {
		http.Error(w, "Error al actualizar registro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rm)
}

func EliminarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Delete(&models.RegistroMantenimiento{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar registro", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"registro eliminado"}`))
}

// ─── PRODUCTO QUIMICO ────────────────────────────────────────────────────────

// CrearProductoQuimico maneja POST /api/v1/quimicos
func CrearProductoQuimico(w http.ResponseWriter, r *http.Request) {
	var pq models.ProductoQuimico

	if err := json.NewDecoder(r.Body).Decode(&pq); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if pq.Nombre == "" {
		http.Error(w, "nombre es obligatorio", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Create(&pq).Error; err != nil {
		http.Error(w, "Error al guardar producto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pq)
}

// ListarProductosQuimicos retorna todos los productos quimicos
func ListarProductosQuimicos(w http.ResponseWriter, r *http.Request) {
	var productos []models.ProductoQuimico

	if err := storage.DB.Find(&productos).Error; err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(productos)
}

func ObtenerProductoQuimico(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var pq models.ProductoQuimico
	if err := storage.DB.First(&pq, id).Error; err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pq)
}

// ActualizarProductoQuimico maneja PATCH /api/v1/quimicos/{id}
func ActualizarProductoQuimico(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var pq models.ProductoQuimico
	if err := storage.DB.First(&pq, id).Error; err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&pq); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Save(&pq).Error; err != nil {
		http.Error(w, "Error al actualizar producto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pq)
}

// EliminarProductoQuimico maneja DELETE /api/v1/quimicos/{id}
func EliminarProductoQuimico(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Delete(&models.ProductoQuimico{}, id).Error; err != nil {
		http.Error(w, "Error al eliminar producto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"producto eliminado"}`))
}
