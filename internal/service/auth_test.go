package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type usuarioRepoMock struct {
	usuarios    map[uint]models.Usuario
	emails      map[string]uint
	siguienteID uint
	errorCrear  error
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

// ── Tests desde HEAD ──────────────────────────────────────────────────

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

// ── Tests desde origin/main ───────────────────────────────────────────

func TestAuthService_LoginValidoGeneraToken(t *testing.T) {
	repo := newUsuarioRepoMock()

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	repo.CrearUsuario(models.Usuario{
		Nombre:       "Admin",
		Email:        "admin@piscina.com",
		PasswordHash: string(hash),
		Rol:          "admin",
	})

	svc := NewAuthService(repo, WithSecreto([]byte("secreto-test")), WithDuracionToken(time.Hour))

	token, usuario, err := svc.Login(" ADMIN@PISCINA.COM ", " admin123 ")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if token == "" {
		t.Fatal("se esperaba token")
	}
	if usuario.ID != 1 || usuario.Rol != "admin" {
		t.Fatalf("usuario inesperado: %+v", usuario)
	}

	claims, err := svc.ValidarToken(token)
	if err != nil {
		t.Fatalf("token generado deberia ser valido: %v", err)
	}
	if claims.UsuarioID != 1 || claims.Rol != "admin" {
		t.Fatalf("claims inesperados: %+v", claims)
	}
}

func TestAuthService_LoginCredencialesInvalidas(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())

	casos := []struct {
		nombre   string
		email    string
		password string
	}{
		{"email vacio", "", "admin123"},
		{"usuario inexistente", "nadie@piscina.com", "admin123"},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			_, _, err := svc.Login(tc.email, tc.password)
			if !errors.Is(err, ErrCredencialesInvalidas) {
				t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
			}
		})
	}
}

func TestAuthService_ValidarTokenInvalido(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo, WithSecreto([]byte("secreto-correcto")))

	if _, err := svc.ValidarToken("token-malo"); !errors.Is(err, ErrCredencialesInvalidas) {
		t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	repo.CrearUsuario(models.Usuario{
		Nombre:       "Admin",
		Email:        "admin@piscina.com",
		PasswordHash: string(hash),
		Rol:          "admin",
	})

	otroSvc := NewAuthService(repo, WithSecreto([]byte("otro-secreto")))
	token, _, err := otroSvc.Login("admin@piscina.com", "admin123")
	if err != nil {
		t.Fatalf("no se pudo generar token: %v", err)
	}

	if _, err := svc.ValidarToken(token); !errors.Is(err, ErrCredencialesInvalidas) {
		t.Fatalf("se esperaba ErrCredencialesInvalidas por secreto distinto, se obtuvo %v", err)
	}
}

func TestAuthService_CrearUsuarioValido(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)

	creado, err := svc.CrearUsuario("Admin", " ADMIN@PISCINA.COM ", " clave123 ", "")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.Email != "admin@piscina.com" {
		t.Fatalf("email no normalizado: %s", creado.Email)
	}
	if creado.Rol != "admin" {
		t.Fatalf("rol por defecto inesperado: %s", creado.Rol)
	}
	if creado.PasswordHash == "" || creado.PasswordHash == "clave123" {
		t.Fatal("se esperaba password hasheado")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(creado.PasswordHash), []byte("clave123")); err != nil {
		t.Fatalf("hash no corresponde a la password: %v", err)
	}
}

func TestAuthService_CrearUsuarioErrores(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)

	if _, err := svc.CrearUsuario("", "admin@piscina.com", "clave123", "admin"); !errors.Is(err, ErrCampoObligatorio) {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}

	if _, err := svc.CrearUsuario("Admin", "admin@piscina.com", "clave123", "admin"); err != nil {
		t.Fatalf("no se esperaba error creando usuario inicial: %v", err)
	}

	if _, err := svc.CrearUsuario("Otro", "ADMIN@PISCINA.COM", "otra", "admin"); !errors.Is(err, ErrEmailEnUso) {
		t.Fatalf("se esperaba ErrEmailEnUso, se obtuvo %v", err)
	}
}

func TestAuthService_ListarYObtenerUsuarios(t *testing.T) {
	repo := newUsuarioRepoMock()
	repo.CrearUsuario(models.Usuario{Nombre: "Admin", Email: "admin@piscina.com", Rol: "admin"})

	svc := NewAuthService(repo)

	if got := len(svc.ListarUsuarios()); got != 1 {
		t.Fatalf("se esperaba 1 usuario, se obtuvo %d", got)
	}

	u, ok := svc.ObtenerUsuario(1)
	if !ok || u.Email != "admin@piscina.com" {
		t.Fatalf("usuario inesperado: %+v ok=%v", u, ok)
	}
	if _, ok := svc.ObtenerUsuario(999); ok {
		t.Fatal("no se esperaba encontrar usuario inexistente")
	}
}

func TestAuthService_ActualizarUsuario(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)

	creado, _ := svc.CrearUsuario("Admin", "admin@piscina.com", "clave123", "admin")

	if _, err := svc.ActualizarUsuario(creado.ID, "", "admin@piscina.com", "", "admin"); !errors.Is(err, ErrCampoObligatorio) {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}

	actualizado, err := svc.ActualizarUsuario(creado.ID, "Admin Editado", " ADMIN@PISCINA.COM ", "", "guardavida")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Nombre != "Admin Editado" || actualizado.Email != "admin@piscina.com" || actualizado.Rol != "guardavida" {
		t.Fatalf("usuario actualizado inesperado: %+v", actualizado)
	}

	actualizado, err = svc.ActualizarUsuario(creado.ID, "Admin Editado", "admin@piscina.com", "nueva123", "admin")
	if err != nil {
		t.Fatalf("no se esperaba error actualizando password: %v", err)
	}
	if actualizado.PasswordHash == "" || actualizado.PasswordHash == "nueva123" {
		t.Fatal("se esperaba nuevo hash de password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(actualizado.PasswordHash), []byte("nueva123")); err != nil {
		t.Fatalf("hash nuevo invalido: %v", err)
	}

	if _, err := svc.ActualizarUsuario(999, "No Existe", "no@piscina.com", "", "admin"); !errors.Is(err, ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestAuthService_BorrarUsuario(t *testing.T) {
	repo := newUsuarioRepoMock()
	svc := NewAuthService(repo)

	creado, _ := svc.CrearUsuario("Admin", "admin@piscina.com", "clave123", "admin")

	if err := svc.BorrarUsuario(creado.ID); err != nil {
		t.Fatalf("no se esperaba error borrando usuario: %v", err)
	}
	if err := svc.BorrarUsuario(999); !errors.Is(err, ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestClientesService_LoginClienteExitoso(t *testing.T) {
	repo := newUsuarioRepoMock()
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
	svc := NewAuthService(newUsuarioRepoMock())

	_, _, err := svc.Login("", "")
	if err != ErrCredencialesInvalidas {
		t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}
}

func TestPagoCliente_ActualizarUsuario(t *testing.T) {
	repo := newUsuarioRepoMock()
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
