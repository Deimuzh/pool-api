package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"pool-api/internal/middleware"
	"pool-api/internal/models"
	"pool-api/internal/service"
)

type fakeAuthRepo struct {
	usuarios map[uint]models.Usuario
	seq      uint
}

func newFakeAuthRepo() *fakeAuthRepo {
	hash, _ := bcrypt.GenerateFromPassword([]byte("clave123"), bcrypt.DefaultCost)
	r := &fakeAuthRepo{usuarios: make(map[uint]models.Usuario), seq: 1}
	r.usuarios[1] = models.Usuario{ID: 1, Nombre: "Admin", Email: "admin@test.com", PasswordHash: string(hash), Rol: "admin"}
	return r
}

func (f *fakeAuthRepo) ListarUsuarios() []models.Usuario {
	lista := make([]models.Usuario, 0, len(f.usuarios))
	for _, u := range f.usuarios {
		lista = append(lista, u)
	}
	return lista
}

func (f *fakeAuthRepo) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	u, ok := f.usuarios[id]
	return u, ok
}

func (f *fakeAuthRepo) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	for _, u := range f.usuarios {
		if u.Email == email {
			return u, true
		}
	}
	return models.Usuario{}, false
}

func (f *fakeAuthRepo) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	f.seq++
	u.ID = f.seq
	f.usuarios[u.ID] = u
	return u, nil
}

func (f *fakeAuthRepo) ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool) {
	_, ok := f.usuarios[id]
	if !ok {
		return models.Usuario{}, false
	}
	datos.ID = id
	f.usuarios[id] = datos
	return datos, true
}

func (f *fakeAuthRepo) BorrarUsuario(id uint) bool {
	_, ok := f.usuarios[id]
	if !ok {
		return false
	}
	delete(f.usuarios, id)
	return true
}

func montarRouterAuthPrueba(t *testing.T) (chi.Router, string) {
	t.Helper()
	repo := newFakeAuthRepo()
	authSvc := service.NewAuthService(repo)

	token, _, err := authSvc.Login("admin@test.com", "clave123")
	if err != nil {
		t.Fatalf("no se pudo generar token: %v", err)
	}

	srv := NewServer(nil, nil, nil, authSvc)
	r := chi.NewRouter()

	r.Post("/api/v1/login", srv.Login)

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Route("/usuarios", func(r chi.Router) {
				r.Get("/", srv.ListarUsuarios)
				r.Post("/", srv.CrearUsuario)
				r.Get("/{id}", srv.ObtenerUsuario)
				r.Put("/{id}", srv.ActualizarUsuario)
				r.Delete("/{id}", srv.BorrarUsuario)
			})
		})
	})

	return r, token
}

func TestLogin_Exitoso_Devuelve200(t *testing.T) {
	router, _ := montarRouterAuthPrueba(t)

	body, _ := json.Marshal(map[string]string{"email": "admin@test.com", "password": "clave123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var resp loginResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Token == "" || resp.Email != "admin@test.com" {
		t.Fatalf("respuesta inesperada: %+v", resp)
	}
}

func TestLogin_JSONInvalido_Devuelve400(t *testing.T) {
	router, _ := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte("{json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestLogin_CredencialesInvalidas_Devuelve401(t *testing.T) {
	router, _ := montarRouterAuthPrueba(t)

	body, _ := json.Marshal(map[string]string{"email": "admin@test.com", "password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

func TestListarUsuarios_ConToken_OK(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}

	var lista []models.Usuario
	json.Unmarshal(rec.Body.Bytes(), &lista)
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 usuario, se obtuvieron %d", len(lista))
	}
}

func TestObtenerUsuario_ConToken_OK(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}

	var u models.Usuario
	json.Unmarshal(rec.Body.Bytes(), &u)
	if u.ID != 1 || u.Nombre != "Admin" {
		t.Fatalf("usuario inesperado: %+v", u)
	}
}

func TestObtenerUsuario_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestObtenerUsuario_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestCrearUsuario_ConToken_CreaYPersiste(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	body, _ := json.Marshal(usuarioRequest{Nombre: "Nuevo", Email: "nuevo@test.com", Password: "pass123", Rol: "user"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/usuarios/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var u models.Usuario
	json.Unmarshal(rec.Body.Bytes(), &u)
	if u.ID == 0 || u.Nombre != "Nuevo" {
		t.Fatalf("usuario creado inesperado: %+v", u)
	}
}

func TestActualizarUsuario_ConToken_OK(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	body, _ := json.Marshal(usuarioRequest{Nombre: "Admin Modificado", Email: "admin@test.com", Rol: "admin"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/usuarios/1", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var u models.Usuario
	json.Unmarshal(rec.Body.Bytes(), &u)
	if u.Nombre != "Admin Modificado" {
		t.Fatalf("nombre inesperado: %s", u.Nombre)
	}
}

func TestActualizarUsuario_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	body, _ := json.Marshal(usuarioRequest{Nombre: "Test", Email: "test@test.com", Rol: "user"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/usuarios/abc", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarUsuario_JSONInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/usuarios/1", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarUsuario_ConToken_OK(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/usuarios/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d", rec.Code)
	}
}

func TestBorrarUsuario_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/usuarios/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarUsuario_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterAuthPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/usuarios/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}
