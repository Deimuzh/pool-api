package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pool-api/internal/models"
<<<<<<< HEAD
=======
	"pool-api/internal/storage"
>>>>>>> f4d782a (integrar lógica de clientes con GORM y tests automatizados)

	"github.com/go-chi/chi/v5"
)

// ─── GUARDAVIDA ──────────────────────────────────────────────────────────────

// ListarGuardavidas atiende GET /api/v1/guardavidas.
func (s *Server) ListarGuardavidas(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Seguridad.ListarGuardavidas())
}

// ObtenerGuardavida atiende GET /api/v1/guardavidas/{id}.
func (s *Server) ObtenerGuardavida(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	g, ok := s.Seguridad.ObtenerGuardavida(uint(idInt))
	if !ok {
		RespondError(w, http.StatusNotFound, "guardavida no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, g)
}

// CrearGuardavida atiende POST /api/v1/guardavidas.
func (s *Server) CrearGuardavida(w http.ResponseWriter, r *http.Request) {
	var g models.Guardavida
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Seguridad.CrearGuardavida(g)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

// ActualizarGuardavida atiende PUT /api/v1/guardavidas/{id}.
func (s *Server) ActualizarGuardavida(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	var g models.Guardavida
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Seguridad.ActualizarGuardavida(uint(idInt), g)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

// BorrarGuardavida atiende DELETE /api/v1/guardavidas/{id}.
func (s *Server) BorrarGuardavida(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Seguridad.BorrarGuardavida(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}

// ─── INCIDENTE ────────────────────────────────────────────────────────────────

// ListarIncidentes atiende GET /api/v1/incidentes.
func (s *Server) ListarIncidentes(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Seguridad.ListarIncidentes())
}

// ObtenerIncidente atiende GET /api/v1/incidentes/{id}.
func (s *Server) ObtenerIncidente(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	inc, ok := s.Seguridad.ObtenerIncidente(uint(idInt))
	if !ok {
		RespondError(w, http.StatusNotFound, "incidente no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, inc)
}

// CrearIncidente atiende POST /api/v1/incidentes.
func (s *Server) CrearIncidente(w http.ResponseWriter, r *http.Request) {
	var inc models.Incidente
	if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Seguridad.CrearIncidente(inc)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

// ActualizarIncidente atiende PUT /api/v1/incidentes/{id}.
func (s *Server) ActualizarIncidente(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	var inc models.Incidente
	if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Seguridad.ActualizarIncidente(uint(idInt), inc)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

// BorrarIncidente atiende DELETE /api/v1/incidentes/{id}.
func (s *Server) BorrarIncidente(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Seguridad.BorrarIncidente(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}

// ─── ACCESO CLIENTE ───────────────────────────────────────────────────────────

// accesoRequest es el body esperado para registrar un acceso: solo se
// necesita el cliente_id, porque CrearAcceso decide Autorizado/Motivo
// internamente según si el cliente tiene un pago de entrada registrado.
type accesoRequest struct {
	ClienteID uint `json:"cliente_id"`
}

// ListarAccesos atiende GET /api/v1/accesos.
func (s *Server) ListarAccesos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Seguridad.ListarAccesos())
}

// CrearAcceso atiende POST /api/v1/accesos.
// Esta es la regla de negocio central del módulo: el acceso solo se autoriza
// si el cliente existe y tiene un pago de entrada registrado.
func (s *Server) CrearAcceso(w http.ResponseWriter, r *http.Request) {
	var body accesoRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Seguridad.CrearAcceso(body.ClienteID)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

// BorrarAcceso atiende DELETE /api/v1/accesos/{id}.
func (s *Server) BorrarAcceso(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	if err := s.Seguridad.BorrarAcceso(uint(idInt)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
