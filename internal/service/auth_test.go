package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type usuarioRepoMock struct {
	usuarios   map[uint]models.Usuario
	emails     map[string]uint
	siguienteID uint
	errorCrear error
}

var _ storage.UsuarioRepository = (*usuarioRepoMock)(nil)

func newUsuarioRepoMock() *usuarioRepoMock {
	return &usuarioRepoMock{
		usuarios: make(map[uint]models.Usuario),
		emails:   make(map[string]uint),
	}
}

func (m *usuarioRepoMock) ListarUsuarios() []models.Usuario {
	lista := make([]models.Usuario, 0, len(m.usuarios))
	for _, u := range m.usuarios {
		lista = append(lista, u)
	}
	return lista
}

func (m *usuarioRepoMock) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	u, ok := m.usuarios[id]
	return u, ok
}

func (m *usuarioRepoMock) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	id, ok := m.emails[email]
	if !ok {
		return models.Usuario{}, false
	}
	return m.usuarios[id], true
}

func (m *usuarioRepoMock) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	if m.errorCrear != nil {
		return models.Usuario{}, m.errorCrear
	}
	m.siguienteID++
	u.ID = m.siguienteID
	m.usuarios[u.ID] = u
	m.emails[u.Email] = u.ID
	return u, nil
}

func (m *usuarioRepoMock) ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool) {
	_, ok := m.usuarios[id]
	if !ok {
		return models.Usuario{}, false
	}
	datos.ID = id
	m.usuarios[id] = datos
	return datos, true
}

func (m *usuarioRepoMock) BorrarUsuario(id uint) bool {
	_, ok := m.usuarios[id]
	if !ok {
		return false
	}
	delete(m.usuarios, id)
	return true
}

func TestAuthService_Login_Exitoso(t *testing.T) {
	repo := newUsuarioRepoMock()
	hash, _ := bcrypt.GenerateFromPassword([]byte("clave123"), bcrypt.DefaultCost)
	repo.CrearUsuario(models.Usuario{Nombre: "Test", Email: "test@correo.com", PasswordHash: string(hash), Rol: "admin"})
	svc := NewAuthService(repo)

	token, u, err := svc.Login("test@correo.com", "clave123")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, "Test", u.Nombre)
}

func TestAuthService_Login_CredencialesInvalidas(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)

	_, _, err := svc.Login("no@existe.com", "clave")
	require.ErrorIs(t, err, ErrCredencialesInvalidas)
}

func TestAuthService_Login_PasswordIncorrecto(t *testing.T) {
	repo := newUsuarioRepoMock()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correcta"), bcrypt.DefaultCost)
	repo.CrearUsuario(models.Usuario{Nombre: "Test", Email: "test@correo.com", PasswordHash: string(hash), Rol: "admin"})
	svc := NewAuthService(repo)

	_, _, err := svc.Login("test@correo.com", "incorrecta")
	require.ErrorIs(t, err, ErrCredencialesInvalidas)
}

func TestAuthService_Login_CamposVacios(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, _, err := svc.Login("", "")
	require.ErrorIs(t, err, ErrCredencialesInvalidas)
}

func TestAuthService_CrearUsuario_Exitoso(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)

	u, err := svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "admin")
	require.NoError(t, err)
	require.NotZero(t, u.ID)
	require.Equal(t, "juan@correo.com", u.Email)
	require.NotEmpty(t, u.PasswordHash)
}

func TestAuthService_CrearUsuario_CamposVacios(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, err := svc.CrearUsuario("", "", "", "")
	require.ErrorIs(t, err, ErrCampoObligatorio)
}

func TestAuthService_CrearUsuario_RolPorDefecto(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)

	u, err := svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "")
	require.NoError(t, err)
	require.Equal(t, "admin", u.Rol)
}

func TestAuthService_CrearUsuario_EmailDuplicado(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)

	svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "admin")
	_, err := svc.CrearUsuario("Pedro", "juan@correo.com", "otra123", "admin")
	require.ErrorIs(t, err, ErrEmailEnUso)
}

func TestAuthService_ActualizarUsuario_Exitoso(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)
	creado, _ := svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "admin")

	actualizado, err := svc.ActualizarUsuario(creado.ID, "Juan Editado", "juan@correo.com", "", "guardavida")
	require.NoError(t, err)
	require.Equal(t, "Juan Editado", actualizado.Nombre)
}

func TestAuthService_ActualizarUsuario_CamposVacios(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, err := svc.ActualizarUsuario(1, "", "", "", "")
	require.ErrorIs(t, err, ErrCampoObligatorio)
}

func TestAuthService_ActualizarUsuario_NoEncontrado(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, err := svc.ActualizarUsuario(99, "Nombre", "email@correo.com", "", "admin")
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestAuthService_BorrarUsuario_Exitoso(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)
	creado, _ := svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "admin")

	err := svc.BorrarUsuario(creado.ID)
	require.NoError(t, err)
}

func TestAuthService_BorrarUsuario_NoEncontrado(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	err := svc.BorrarUsuario(99)
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestAuthService_ListarUsuarios(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)
	svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "admin")
	svc.CrearUsuario("Ana", "ana@correo.com", "clave456", "guardavida")

	usuarios := svc.ListarUsuarios()
	require.Len(t, usuarios, 2)
}

func TestAuthService_ObtenerUsuario(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)
	creado, _ := svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "admin")

	u, ok := svc.ObtenerUsuario(creado.ID)
	require.True(t, ok)
	require.Equal(t, "Juan", u.Nombre)

	_, ok = svc.ObtenerUsuario(99)
	require.False(t, ok)
}

func TestAuthService_ValidarToken_Exitoso(t *testing.T) {
	repo := newUsuarioRepoMock()
	hash, _ := bcrypt.GenerateFromPassword([]byte("clave123"), bcrypt.DefaultCost)
	repo.CrearUsuario(models.Usuario{Nombre: "Test", Email: "test@correo.com", PasswordHash: string(hash), Rol: "admin"})
	svc := NewAuthService(repo)

	token, _, _ := svc.Login("test@correo.com", "clave123")
	claims, err := svc.ValidarToken(token)
	require.NoError(t, err)
	require.Equal(t, "admin", claims.Rol)
}

func TestAuthService_ValidarToken_Invalido(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, err := svc.ValidarToken("token-malformado")
	require.ErrorIs(t, err, ErrCredencialesInvalidas)
}

func TestAuthService_WithOpciones(t *testing.T) {
	svc := NewAuthService(
		newUsuarioRepoMock(),
		WithSecreto([]byte("mi-secreto-personalizado")),
		WithDuracionToken(1*time.Hour),
	)
	require.NotNil(t, svc)
}

func TestAuthService_WithSecretoVacioNoCambia(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock(), WithSecreto([]byte{}))
	require.NotNil(t, svc)
}

func TestAuthService_WithDuracionCeroNoCambia(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock(), WithDuracionToken(0))
	require.NotNil(t, svc)
}
