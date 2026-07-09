package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pool-api/internal/models"

	"github.com/go-chi/chi/v5"
)

// ─── EQUIPO ──────────────────────────────────────────────────────────────────

func (s *Server) ListarEquipos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Mantenimiento.ListarEquipos())
}

func (s *Server) ObtenerEquipo(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	eq, ok := s.Mantenimiento.ObtenerEquipo(uint(idInt))
	if !ok {
		RespondError(w, http.StatusNotFound, "equipo no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, eq)
}

func (s *Server) CrearEquipo(w http.ResponseWriter, r *http.Request) {
	var eq models.Equipo
	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Mantenimiento.CrearEquipo(eq)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ActualizarEquipo(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	var eq models.Equipo
	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarEquipo(uint(idInt), eq)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarEquipo(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Mantenimiento.BorrarEquipo(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "equipo eliminado"})
}

// ─── REGISTRO MANTENIMIENTO ──────────────────────────────────────────────────

func (s *Server) ListarRegistrosMantenimiento(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Mantenimiento.ListarRegistros())
}

func (s *Server) ObtenerRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	rm, ok := s.Mantenimiento.ObtenerRegistro(uint(idInt))
	if !ok {
		RespondError(w, http.StatusNotFound, "registro no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, rm)
}

func (s *Server) CrearRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	var rm models.RegistroMantenimiento
	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Mantenimiento.CrearRegistro(rm)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ActualizarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	var rm models.RegistroMantenimiento
	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarRegistro(uint(idInt), rm)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Mantenimiento.BorrarRegistro(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "registro eliminado"})
}

// ─── PRODUCTO QUIMICO ────────────────────────────────────────────────────────

func (s *Server) ListarQuimicos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Mantenimiento.ListarQuimicos())
}

func (s *Server) ObtenerQuimico(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	q, ok := s.Mantenimiento.ObtenerQuimico(uint(idInt))
	if !ok {
		RespondError(w, http.StatusNotFound, "producto no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, q)
}
//Recibe la request HTTP, parsea JSON/parámetros, llama al service y devuelve JSON.
func (s *Server) CrearQuimico(w http.ResponseWriter, r *http.Request) {
	var q models.ProductoQuimico
// 1. Decodificar JSON del body
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
// 2. Llamar al service (nunca toca la BD directamente)
	creado, err := s.Mantenimiento.CrearQuimico(q)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
// 3. Responder con el resultado
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ActualizarQuimico(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	var q models.ProductoQuimico
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarQuimico(uint(idInt), q)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarQuimico(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Mantenimiento.BorrarQuimico(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "producto eliminado"})
}
