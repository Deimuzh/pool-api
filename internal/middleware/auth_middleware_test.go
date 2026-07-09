package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"pool-api/internal/models"
	"pool-api/internal/service"
)

type usuarioRepoParaTest struct {
	usuario models.Usuario
}

func (u *usuarioRepoParaTest) ListarUsuarios() []models.Usuario { return nil }
func (u *usuarioRepoParaTest) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (u *usuarioRepoParaTest) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	if email == u.usuario.Email {
		return u.usuario, true
	}
	return models.Usuario{}, false
}
func (u *usuarioRepoParaTest) CrearUsuario(usr models.Usuario) (models.Usuario, error) { return usr, nil }
func (u *usuarioRepoParaTest) ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (u *usuarioRepoParaTest) BorrarUsuario(id uint) bool { return false }

func nuevoAuthServiceTest() *service.AuthService {
	return service.NewAuthService(&usuarioRepoParaTest{
		usuario: models.Usuario{ID: 1, Nombre: "Admin", Email: "admin@test.com", PasswordHash: "$2a$10$dummy", Rol: "admin"},
	})
}

func TestAuth_SinToken_Devuelve401(t *testing.T) {
	authSvc := nuevoAuthServiceTest()
	handler := Auth(authSvc)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

func TestAuth_TokenInvalido_Devuelve401(t *testing.T) {
	authSvc := nuevoAuthServiceTest()
	handler := Auth(authSvc)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", "Bearer token-invalido")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

func TestAuth_FormatoInvalido_Devuelve401(t *testing.T) {
	authSvc := nuevoAuthServiceTest()
	handler := Auth(authSvc)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", "Basic token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

func TestSoloRoles_SinRolEnContexto_Devuelve403(t *testing.T) {
	handler := SoloRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyRol, "otro"))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("se esperaba 403, se obtuvo %d", rec.Code)
	}
}

func TestSoloRoles_MultiplesRolesPermitidos(t *testing.T) {
	handler := SoloRoles("admin", "guardavida")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyRol, "guardavida"))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d", rec.Code)
	}
}
