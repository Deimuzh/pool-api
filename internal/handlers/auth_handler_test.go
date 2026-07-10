package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"pool-api/internal/models"
	"pool-api/internal/service"
)

type usuarioHandlerMock struct {
	usuarios   map[uint]models.Usuario
	emails     map[string]uint
	siguienteID uint
}

func newUsuarioHandlerMock() *usuarioHandlerMock {
	return &usuarioHandlerMock{
		usuarios: make(map[uint]models.Usuario),
		emails:   make(map[string]uint),
	}
}

func (m *usuarioHandlerMock) ListarUsuarios() []models.Usuario {
	lista := make([]models.Usuario, 0, len(m.usuarios))
	for _, u := range m.usuarios {
		lista = append(lista, u)
	}
	return lista
}

func (m *usuarioHandlerMock) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	u, ok := m.usuarios[id]
	return u, ok
}

func (m *usuarioHandlerMock) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	id, ok := m.emails[email]
	if !ok {
		return models.Usuario{}, false
	}
	return m.usuarios[id], true
}

func (m *usuarioHandlerMock) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	m.siguienteID++
	u.ID = m.siguienteID
	m.usuarios[u.ID] = u
	m.emails[u.Email] = u.ID
	return u, nil
}

func (m *usuarioHandlerMock) ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool) {
	_, ok := m.usuarios[id]
	if !ok {
		return models.Usuario{}, false
	}
	datos.ID = id
	m.usuarios[id] = datos
	return datos, true
}

func (m *usuarioHandlerMock) BorrarUsuario(id uint) bool {
	_, ok := m.usuarios[id]
	if !ok {
		return false
	}
	delete(m.usuarios, id)
	return true
}

func montarRouterAuthPrueba(t *testing.T) *chi.Mux {
	t.Helper()
	usuarios := newUsuarioHandlerMock()
	authSvc := service.NewAuthService(usuarios)
	// crear un usuario de prueba
	authSvc.CrearUsuario("Admin", "admin@test.com", "clave123", "admin")
	server := NewServer(nil, nil, nil, authSvc)
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", server.Login)
		r.Route("/usuarios", func(r chi.Router) {
			r.Get("/", server.ListarUsuarios)
			r.Post("/", server.CrearUsuario)
			r.Get("/{id}", server.ObtenerUsuario)
			r.Put("/{id}", server.ActualizarUsuario)
			r.Delete("/{id}", server.BorrarUsuario)
		})
	})
	return r
}

func TestLogin_Exitoso(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	body, _ := json.Marshal(loginRequest{Email: "admin@test.com", Password: "clave123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	var resp loginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}
	if resp.Token == "" {
		t.Error("se esperaba un token no vacio")
	}
	if resp.Nombre != "Admin" {
		t.Errorf("nombre inesperado: %s", resp.Nombre)
	}
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	body, _ := json.Marshal(loginRequest{Email: "admin@test.com", Password: "incorrecta"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

func TestLogin_JSONInvalido(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte("{malformado")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestListarUsuarios(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
	var usuarios []models.Usuario
	if err := json.Unmarshal(rec.Body.Bytes(), &usuarios); err != nil {
		t.Fatalf("error decodificando: %v", err)
	}
	if len(usuarios) == 0 {
		t.Error("se esperaba al menos un usuario")
	}
}

func TestCrearUsuario_Handler(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	body, _ := json.Marshal(usuarioRequest{Nombre: "Nuevo", Email: "nuevo@test.com", Password: "clave", Rol: "guardavida"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/usuarios/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCrearUsuario_JSONInvalido(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/usuarios/", bytes.NewReader([]byte("{mal")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestObtenerUsuario_Handler(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestObtenerUsuario_NoEncontrado(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestObtenerUsuario_IDInvalido(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarUsuario_Handler(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	body, _ := json.Marshal(usuarioRequest{Nombre: "Editado", Email: "editado@test.com", Rol: "admin"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/usuarios/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestBorrarUsuario_Handler(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/usuarios/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d", rec.Code)
	}
}

func TestBorrarUsuario_NoEncontrado(t *testing.T) {
	router := montarRouterAuthPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/usuarios/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}
