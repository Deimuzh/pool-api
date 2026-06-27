package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

// secretJWT firma y valida los tokens. En un entorno real esto debería
// venir de una variable de entorno, nunca quedar fijo en el código.
var secretJWT = []byte("piscina-los-ceibos-clave-secreta")

const duracionToken = 24 * time.Hour

// Claims es la información que viaja dentro del JWT.
type Claims struct {
	UsuarioID int    `json:"uid"`
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}

// AuthService maneja login, generación/validación de tokens y el CRUD de usuarios.
type AuthService struct {
	repo storage.UsuarioRepository
}

func NewAuthService(repo storage.UsuarioRepository) *AuthService {
	return &AuthService{repo: repo}
}

// ─── LOGIN ────────────────────────────────────────────────────────────────────

// Login verifica email+contraseña contra el hash guardado y, si coincide,
// devuelve un JWT firmado válido por 24 horas.
func (s *AuthService) Login(email, password string) (string, models.Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		return "", models.Usuario{}, ErrCredencialesInvalidas
	}

	u, existe := s.repo.BuscarUsuarioPorEmail(email)
	if !existe {
		return "", models.Usuario{}, ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", models.Usuario{}, ErrCredencialesInvalidas
	}

	token, err := s.generarToken(u)
	if err != nil {
		return "", models.Usuario{}, err
	}
	return token, u, nil
}

func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := &Claims{
		UsuarioID: u.ID,
		Rol:       u.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretJWT)
}

// ValidarToken decodifica y verifica un JWT. La usará el middleware de auth
// cuando decidan proteger rutas (por ahora no se usa en ninguna ruta).
func (s *AuthService) ValidarToken(tokenStr string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return secretJWT, nil
	})
	if err != nil || !token.Valid {
		return Claims{}, ErrCredencialesInvalidas
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return Claims{}, ErrCredencialesInvalidas
	}
	return *claims, nil
}

// ─── CRUD DE USUARIOS ─────────────────────────────────────────────────────────

func (s *AuthService) ListarUsuarios() []models.Usuario {
	return s.repo.ListarUsuarios()
}

func (s *AuthService) ObtenerUsuario(id int) (models.Usuario, bool) {
	return s.repo.BuscarUsuarioPorID(id)
}

// CrearUsuario registra una nueva cuenta. La contraseña llega en texto plano
// (vía HTTPS en producción) y se hashea con bcrypt antes de guardar.
func (s *AuthService) CrearUsuario(nombre, email, password, rol string) (models.Usuario, error) {
	nombre = strings.TrimSpace(nombre)
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if nombre == "" || email == "" || password == "" {
		return models.Usuario{}, ErrCampoObligatorio
	}
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}
	if rol == "" {
		rol = "admin"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}

	return s.repo.CrearUsuario(models.Usuario{
		Nombre:       nombre,
		Email:        email,
		PasswordHash: string(hash),
		Rol:          rol,
	})
}

// ActualizarUsuario edita nombre/email/rol. Si password viene vacío, conserva
// la contraseña actual (el repo ya maneja esa regla); si viene con texto,
// se vuelve a hashear aquí antes de pasarlo al repo.
func (s *AuthService) ActualizarUsuario(id int, nombre, email, password, rol string) (models.Usuario, error) {
	nombre = strings.TrimSpace(nombre)
	email = strings.TrimSpace(strings.ToLower(email))

	if nombre == "" || email == "" {
		return models.Usuario{}, ErrCampoObligatorio
	}

	datos := models.Usuario{Nombre: nombre, Email: email, Rol: rol}

	if password = strings.TrimSpace(password); password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return models.Usuario{}, err
		}
		datos.PasswordHash = string(hash)
	}

	actualizado, ok := s.repo.ActualizarUsuario(id, datos)
	if !ok {
		return models.Usuario{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *AuthService) BorrarUsuario(id int) error {
	if !s.repo.BorrarUsuario(id) {
		return ErrNoEncontrado
	}
	return nil
}
