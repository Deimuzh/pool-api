package service

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"pool-api/internal/models"
)

// mockUsuarioRepo es un repository falso en memoria para probar AuthService
// sin depender de GORM ni de una base de datos real.
type mockUsuarioRepo struct {
	usuarios map[uint]models.Usuario
	nextID   uint
}

// nuevoMockUsuarioRepo crea un repo vacío con contador de IDs autoincremental.
func nuevoMockUsuarioRepo() *mockUsuarioRepo {
	return &mockUsuarioRepo{
		usuarios: make(map[uint]models.Usuario),
		nextID:   1,
	}
}

// ListarUsuarios devuelve todos los usuarios guardados en memoria.
func (m *mockUsuarioRepo) ListarUsuarios() []models.Usuario {
	lista := make([]models.Usuario, 0, len(m.usuarios))
	for _, u := range m.usuarios {
		lista = append(lista, u)
	}
	return lista
}

// BuscarUsuarioPorID simula una búsqueda por clave primaria.
func (m *mockUsuarioRepo) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	u, ok := m.usuarios[id]
	return u, ok
}

// BuscarUsuarioPorEmail recorre el mapa y devuelve el usuario con email exacto.
func (m *mockUsuarioRepo) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	for _, u := range m.usuarios {
		if u.Email == email {
			return u, true
		}
	}
	return models.Usuario{}, false
}

// CrearUsuario simula el insert asignando ID y guardando en memoria.
func (m *mockUsuarioRepo) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = m.nextID
	m.nextID++
	m.usuarios[u.ID] = u
	return u, nil
}

// ActualizarUsuario simula update: si el ID existe reemplaza los datos.
func (m *mockUsuarioRepo) ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool) {
	if _, ok := m.usuarios[id]; !ok {
		return models.Usuario{}, false
	}
	datos.ID = id
	m.usuarios[id] = datos
	return datos, true
}

// BorrarUsuario simula delete y devuelve false si el ID no existe.
func (m *mockUsuarioRepo) BorrarUsuario(id uint) bool {
	if _, ok := m.usuarios[id]; !ok {
		return false
	}
	delete(m.usuarios, id)
	return true
}

// TestAuthService_LoginValidoGeneraToken verifica que credenciales correctas
// generen un JWT válido con el ID y rol del usuario.
func TestAuthService_LoginValidoGeneraToken(t *testing.T) {
	repo := nuevoMockUsuarioRepo()

	// Guardamos un hash real para probar bcrypt igual que en producción.
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	repo.usuarios[1] = models.Usuario{
		ID:           1,
		Nombre:       "Admin",
		Email:        "admin@piscina.com",
		PasswordHash: string(hash),
		Rol:          "admin",
	}

	// Inyectamos secreto/duración para poder validar el token generado.
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

	// ValidarToken debe recuperar los claims que Login firmó.
	claims, err := svc.ValidarToken(token)
	if err != nil {
		t.Fatalf("token generado debería ser válido: %v", err)
	}
	if claims.UsuarioID != 1 || claims.Rol != "admin" {
		t.Fatalf("claims inesperados: %+v", claims)
	}
}

// TestAuthService_LoginCredencialesInvalidas cubre email vacío, usuario
// inexistente y contraseña incorrecta.
func TestAuthService_LoginCredencialesInvalidas(t *testing.T) {
	repo := nuevoMockUsuarioRepo()

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	repo.usuarios[1] = models.Usuario{
		ID:           1,
		Email:        "admin@piscina.com",
		PasswordHash: string(hash),
		Rol:          "admin",
	}

	svc := NewAuthService(repo)

	casos := []struct {
		nombre   string
		email    string
		password string
	}{
		{"email vacio", "", "admin123"},
		{"usuario inexistente", "nadie@piscina.com", "admin123"},
		{"password incorrecto", "admin@piscina.com", "mal"},
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

// TestAuthService_ValidarTokenInvalido verifica que un token malformado o
// firmado con otro secreto sea rechazado.
func TestAuthService_ValidarTokenInvalido(t *testing.T) {
	repo := nuevoMockUsuarioRepo()
	svc := NewAuthService(repo, WithSecreto([]byte("secreto-correcto")))

	// Token malformado: no tiene estructura JWT válida.
	if _, err := svc.ValidarToken("token-malo"); !errors.Is(err, ErrCredencialesInvalidas) {
		t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	repo.usuarios[1] = models.Usuario{
		ID:           1,
		Email:        "admin@piscina.com",
		PasswordHash: string(hash),
		Rol:          "admin",
	}

	// Generamos un token con otro secreto para confirmar que el service actual
	// lo rechaza por firma inválida.
	otroSvc := NewAuthService(repo, WithSecreto([]byte("otro-secreto")))
	token, _, err := otroSvc.Login("admin@piscina.com", "admin123")
	if err != nil {
		t.Fatalf("no se pudo generar token: %v", err)
	}

	if _, err := svc.ValidarToken(token); !errors.Is(err, ErrCredencialesInvalidas) {
		t.Fatalf("se esperaba ErrCredencialesInvalidas por secreto distinto, se obtuvo %v", err)
	}
}

// TestAuthService_CrearUsuarioValido verifica que el service normalice el email,
// asigne rol admin por defecto y guarde password hasheado.
func TestAuthService_CrearUsuarioValido(t *testing.T) {
	repo := nuevoMockUsuarioRepo()
	svc := NewAuthService(repo)

	// Enviamos email/password con espacios para probar TrimSpace y ToLower.
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

// TestAuthService_CrearUsuarioErrores cubre campos obligatorios y email duplicado.
func TestAuthService_CrearUsuarioErrores(t *testing.T) {
	repo := nuevoMockUsuarioRepo()
	svc := NewAuthService(repo)

	// Nombre vacío debe cortar antes de guardar.
	if _, err := svc.CrearUsuario("", "admin@piscina.com", "clave123", "admin"); !errors.Is(err, ErrCampoObligatorio) {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}

	if _, err := svc.CrearUsuario("Admin", "admin@piscina.com", "clave123", "admin"); err != nil {
		t.Fatalf("no se esperaba error creando usuario inicial: %v", err)
	}

	// El service normaliza email, por eso ADMIN@... debe detectarse duplicado.
	if _, err := svc.CrearUsuario("Otro", "ADMIN@PISCINA.COM", "otra", "admin"); !errors.Is(err, ErrEmailEnUso) {
		t.Fatalf("se esperaba ErrEmailEnUso, se obtuvo %v", err)
	}
}

// TestAuthService_ListarYObtenerUsuarios verifica que el service delegue listado
// y búsqueda por ID al repository.
func TestAuthService_ListarYObtenerUsuarios(t *testing.T) {
	repo := nuevoMockUsuarioRepo()
	repo.usuarios[1] = models.Usuario{ID: 1, Nombre: "Admin", Email: "admin@piscina.com", Rol: "admin"}

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

// TestAuthService_ActualizarUsuario verifica actualización sin password y con
// password nueva, incluyendo el caso de ID inexistente.
func TestAuthService_ActualizarUsuario(t *testing.T) {
	repo := nuevoMockUsuarioRepo()
	repo.usuarios[1] = models.Usuario{
		ID:           1,
		Nombre:       "Admin",
		Email:        "admin@piscina.com",
		PasswordHash: "hash-viejo",
		Rol:          "admin",
	}
	svc := NewAuthService(repo)

	// Nombre vacío debe devolver error de validación.
	if _, err := svc.ActualizarUsuario(1, "", "admin@piscina.com", "", "admin"); !errors.Is(err, ErrCampoObligatorio) {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}

	// Sin password nueva: el service manda PasswordHash vacío y el repo conserva
	// la regla en producción. Este mock reemplaza datos, pero igual probamos
	// normalización de email y rol.
	actualizado, err := svc.ActualizarUsuario(1, "Admin Editado", " ADMIN@PISCINA.COM ", "", "guardavida")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Nombre != "Admin Editado" || actualizado.Email != "admin@piscina.com" || actualizado.Rol != "guardavida" {
		t.Fatalf("usuario actualizado inesperado: %+v", actualizado)
	}

	// Con password nueva: debe llegar hasheada al repo.
	actualizado, err = svc.ActualizarUsuario(1, "Admin Editado", "admin@piscina.com", "nueva123", "admin")
	if err != nil {
		t.Fatalf("no se esperaba error actualizando password: %v", err)
	}
	if actualizado.PasswordHash == "" || actualizado.PasswordHash == "nueva123" {
		t.Fatal("se esperaba nuevo hash de password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(actualizado.PasswordHash), []byte("nueva123")); err != nil {
		t.Fatalf("hash nuevo inválido: %v", err)
	}

	if _, err := svc.ActualizarUsuario(999, "No Existe", "no@piscina.com", "", "admin"); !errors.Is(err, ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

// TestAuthService_BorrarUsuario verifica borrado exitoso y error cuando el ID no existe.
func TestAuthService_BorrarUsuario(t *testing.T) {
	repo := nuevoMockUsuarioRepo()
	repo.usuarios[1] = models.Usuario{ID: 1, Nombre: "Admin"}

	svc := NewAuthService(repo)

	if err := svc.BorrarUsuario(1); err != nil {
		t.Fatalf("no se esperaba error borrando usuario: %v", err)
	}
	if err := svc.BorrarUsuario(999); !errors.Is(err, ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestClientesService_LoginClienteExitoso(t *testing.T) {
	repo := nuevoMockUsuarioRepo()
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
	repo := nuevoMockUsuarioRepo()
	svc := NewAuthService(repo)

	_, _, err := svc.Login("", "")
	if err != ErrCredencialesInvalidas {
		t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}
}

func TestPagoCliente_ActualizarUsuario(t *testing.T) {
	repo := nuevoMockUsuarioRepo()
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
