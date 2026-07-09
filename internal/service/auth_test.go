package service

import (
	"testing"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type authRepoMock struct {
	usuarios map[uint]models.Usuario
}

var _ storage.UsuarioRepository = (*authRepoMock)(nil)

func newAuthRepoMock() *authRepoMock {
	return &authRepoMock{usuarios: make(map[uint]models.Usuario)}
}

func (m *authRepoMock) ListarUsuarios() []models.Usuario {
	lista := make([]models.Usuario, 0, len(m.usuarios))
	for _, u := range m.usuarios {
		lista = append(lista, u)
	}
	return lista
}

func (m *authRepoMock) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	u, ok := m.usuarios[id]
	return u, ok
}

func (m *authRepoMock) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	for _, u := range m.usuarios {
		if u.Email == email {
			return u, true
		}
	}
	return models.Usuario{}, false
}

func (m *authRepoMock) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = uint(len(m.usuarios)) + 1
	m.usuarios[u.ID] = u
	return u, nil
}

func (m *authRepoMock) ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool) {
	if _, ok := m.usuarios[id]; !ok {
		return models.Usuario{}, false
	}
	datos.ID = id
	m.usuarios[id] = datos
	return datos, true
}

func (m *authRepoMock) BorrarUsuario(id uint) bool {
	if _, ok := m.usuarios[id]; !ok {
		return false
	}
	delete(m.usuarios, id)
	return true
}

func TestClientesService_LoginClienteExitoso(t *testing.T) {
	repo := newAuthRepoMock()
	svc := NewAuthService(repo)

	u, _ := svc.CrearUsuario("Cliente Uno", "cliente@test.com", "pass123", "cliente")

	token, _, err := svc.Login(u.Email, "pass123")
	if err != nil {
		t.Fatalf("se esperaba login exitoso, se obtuvo %v", err)
	}
	if token == "" {
		t.Fatal("el token no debe estar vacio")
	}
}

func TestClientesService_LoginClienteCredencialesInvalidas(t *testing.T) {
	repo := newAuthRepoMock()
	svc := NewAuthService(repo)

	_, _, err := svc.Login("", "")
	if err != ErrCredencialesInvalidas {
		t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}
}

func TestPagoCliente_ActualizarUsuario(t *testing.T) {
	repo := newAuthRepoMock()
	svc := NewAuthService(repo)

	original, _ := svc.CrearUsuario("Original", "orig@test.com", "pass123", "admin")
	actualizado, err := svc.ActualizarUsuario(original.ID, "Actualizado", "orig@test.com", "", "admin")
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo %v", err)
	}
	if actualizado.Nombre != "Actualizado" {
		t.Fatal("el nombre deberia haberse actualizado")
	}
}
