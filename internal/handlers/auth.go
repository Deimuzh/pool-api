package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ─── LOGIN ────────────────────────────────────────────────────────────────────

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token  string `json:"token"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
}

// Login atiende POST /api/v1/login.
// Body esperado: {"email": "...", "password": "..."}
// Responde con un JWT que el frontend debe guardar y reenviar en
// el header "Authorization: Bearer <token>" en próximas peticiones
// (aunque, por ahora, ninguna ruta exige ese header todavía).
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	token, usuario, err := s.Auth.Login(body.Email, body.Password)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, loginResponse{
		Token:  token,
		Nombre: usuario.Nombre,
		Email:  usuario.Email,
		Rol:    usuario.Rol,
	})
}

// ─── CRUD USUARIOS ────────────────────────────────────────────────────────────

type usuarioRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"` // solo obligatorio al crear; opcional al actualizar
	Rol      string `json:"rol"`
}

// ListarUsuarios atiende GET /api/v1/usuarios.
func (s *Server) ListarUsuarios(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Auth.ListarUsuarios())
}

// ObtenerUsuario atiende GET /api/v1/usuarios/{id}.
func (s *Server) ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	idUint, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	id := uint(idUint)
	u, ok := s.Auth.ObtenerUsuario(id)
	if !ok {
		RespondError(w, http.StatusNotFound, "usuario no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, u)
}

// CrearUsuario atiende POST /api/v1/usuarios.
func (s *Server) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	var body usuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Auth.CrearUsuario(body.Nombre, body.Email, body.Password, body.Rol)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

// ActualizarUsuario atiende PUT /api/v1/usuarios/{id}.
// Si "password" llega vacío en el body, se conserva la contraseña actual.
func (s *Server) ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	idUint, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	id := uint(idUint)
	var body usuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Auth.ActualizarUsuario(id, body.Nombre, body.Email, body.Password, body.Rol)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

// BorrarUsuario atiende DELETE /api/v1/usuarios/{id}.
func (s *Server) BorrarUsuario(w http.ResponseWriter, r *http.Request) {
	idUint, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}
	id := uint(idUint)
	if err := s.Auth.BorrarUsuario(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
