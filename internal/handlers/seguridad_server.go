package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"pool-api/internal/models"

	"github.com/go-chi/chi/v5"
)

func (s *Server) ListarGuardavidas(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Seguridad.ListarGuardavidas())
}

func (s *Server) CrearGuardavida(w http.ResponseWriter, r *http.Request) {
	var g models.Guardavida
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado, err := s.Seguridad.CrearGuardavida(g)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerGuardavida(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	g, ok := s.Seguridad.ObtenerGuardavida(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "Guardavida no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, g)
}

func (s *Server) ActualizarGuardavida(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var g models.Guardavida
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	actualizado, err := s.Seguridad.ActualizarGuardavida(id, g)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarGuardavida(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Seguridad.BorrarGuardavida(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "guardavida eliminado"})
}

func (s *Server) ListarIncidentes(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Seguridad.ListarIncidentes())
}

func (s *Server) CrearIncidente(w http.ResponseWriter, r *http.Request) {
	var inc models.Incidente
	if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado, err := s.Seguridad.CrearIncidente(inc)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerIncidente(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	inc, ok := s.Seguridad.ObtenerIncidente(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "Incidente no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, inc)
}

func (s *Server) ActualizarIncidente(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var inc models.Incidente
	if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	actualizado, err := s.Seguridad.ActualizarIncidente(id, inc)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarIncidente(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Seguridad.BorrarIncidente(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "incidente eliminado"})
}

func (s *Server) ListarAccesos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Seguridad.ListarAccesos())
}

func (s *Server) CrearAcceso(w http.ResponseWriter, r *http.Request) {
	var acc models.AccesoCliente
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado, err := s.Seguridad.CrearAcceso(int(acc.ClienteID))
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) BorrarAcceso(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Seguridad.BorrarAcceso(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "acceso eliminado"})
}

func (s *Server) ListarEquipos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Mantenimiento.ListarEquipos())
}

func (s *Server) CrearEquipo(w http.ResponseWriter, r *http.Request) {
	var e models.Equipo
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado, err := s.Mantenimiento.CrearEquipo(e)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	g, ok := s.Mantenimiento.ObtenerEquipo(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "Equipo no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, g)
}

func (s *Server) ActualizarEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var e models.Equipo
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarEquipo(id, e)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Mantenimiento.BorrarEquipo(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "equipo eliminado"})
}

func (s *Server) ListarRegistrosMantenimiento(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Mantenimiento.ListarRegistros())
}

func (s *Server) CrearRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	var rm models.RegistroMantenimiento
	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado, err := s.Mantenimiento.CrearRegistro(rm)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	rm, ok := s.Mantenimiento.ObtenerRegistro(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "Registro no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, rm)
}

func (s *Server) ActualizarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var rm models.RegistroMantenimiento
	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarRegistro(id, rm)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Mantenimiento.BorrarRegistro(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "registro eliminado"})
}

func (s *Server) ListarQuimicos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Mantenimiento.ListarQuimicos())
}

func (s *Server) CrearQuimico(w http.ResponseWriter, r *http.Request) {
	var q models.ProductoQuimico
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado, err := s.Mantenimiento.CrearQuimico(q)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerQuimico(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	q, ok := s.Mantenimiento.ObtenerQuimico(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "Producto químico no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, q)
}

func (s *Server) ActualizarQuimico(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var q models.ProductoQuimico
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarQuimico(id, q)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarQuimico(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Mantenimiento.BorrarQuimico(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "químico eliminado"})
}

var _ = time.Now
