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

// ─── FAKE EN MEMORIA ────────────────────────────────────────────────────────
//
// A diferencia de un mock (que solo registra llamadas), este fake SÍ guarda
// los datos de verdad, pero en un slice en memoria en vez de una base real.
// Sirve para probar el handler de punta a punta (JSON → service → "repo")
// sin depender de SQLite.

type fakeSeguridadRepo struct {
	guardavidas []models.Guardavida
	siguienteID uint
}

func (f *fakeSeguridadRepo) ListarGuardavidas() []models.Guardavida { return f.guardavidas }
func (f *fakeSeguridadRepo) BuscarGuardavidaPorID(id uint) (models.Guardavida, bool) {
	for _, g := range f.guardavidas {
		if g.ID == id {
			return g, true
		}
	}
	return models.Guardavida{}, false
}
func (f *fakeSeguridadRepo) CrearGuardavida(g models.Guardavida) models.Guardavida {
	f.siguienteID++
	g.ID = f.siguienteID
	f.guardavidas = append(f.guardavidas, g)
	return g
}
func (f *fakeSeguridadRepo) ActualizarGuardavida(id uint, datos models.Guardavida) (models.Guardavida, bool) {
	return models.Guardavida{}, false
}
func (f *fakeSeguridadRepo) BorrarGuardavida(id uint) bool { return false }

func (f *fakeSeguridadRepo) ListarIncidentes() []models.Incidente { return nil }
func (f *fakeSeguridadRepo) BuscarIncidentePorID(id uint) (models.Incidente, bool) {
	return models.Incidente{}, false
}
func (f *fakeSeguridadRepo) CrearIncidente(i models.Incidente) models.Incidente { return i }
func (f *fakeSeguridadRepo) ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool) {
	return models.Incidente{}, false
}
func (f *fakeSeguridadRepo) BorrarIncidente(id uint) bool { return false }

func (f *fakeSeguridadRepo) ListarAccesos() []models.AccesoCliente { return nil }
func (f *fakeSeguridadRepo) BuscarAccesoPorID(id uint) (models.AccesoCliente, bool) {
	return models.AccesoCliente{}, false
}
func (f *fakeSeguridadRepo) CrearAcceso(a models.AccesoCliente) models.AccesoCliente { return a }
func (f *fakeSeguridadRepo) ActualizarAcceso(id uint, datos models.AccesoCliente) (models.AccesoCliente, bool) {
	return models.AccesoCliente{}, false
}
func (f *fakeSeguridadRepo) BorrarAcceso(id uint) bool { return false }

// fakeClienteRepo y fakePagoRepo: SeguridadService los necesita para
// construirse, pero estos tests de Guardavida no los usan, así que basta
// con que implementen la interfaz sin hacer nada relevante.
type fakeClienteRepo struct{}

func (f *fakeClienteRepo) ListarClientes() []models.Cliente { return nil }
func (f *fakeClienteRepo) BuscarClientePorID(id uint) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (f *fakeClienteRepo) CrearCliente(c models.Cliente) models.Cliente { return c }
func (f *fakeClienteRepo) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (f *fakeClienteRepo) BorrarCliente(id uint) bool { return false }

type fakePagoRepo struct{}

func (f *fakePagoRepo) ListarPagos() []models.Pago                  { return nil }
func (f *fakePagoRepo) BuscarPagoPorID(id uint) (models.Pago, bool) { return models.Pago{}, false }
func (f *fakePagoRepo) CrearPago(p models.Pago) models.Pago         { return p }
func (f *fakePagoRepo) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}
func (f *fakePagoRepo) BorrarPago(id uint) bool                     { return false }
func (f *fakePagoRepo) ClienteTienePagoEntrada(clienteID uint) bool { return false }

// fakeUsuarioRepo permite generar un JWT real vía AuthService.Login, igual
// que en producción, sin tocar una base de datos.
type fakeUsuarioRepo struct {
	usuario models.Usuario
}

func (f *fakeUsuarioRepo) ListarUsuarios() []models.Usuario { return []models.Usuario{f.usuario} }
func (f *fakeUsuarioRepo) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	if id == f.usuario.ID {
		return f.usuario, true
	}
	return models.Usuario{}, false
}
func (f *fakeUsuarioRepo) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	if email == f.usuario.Email {
		return f.usuario, true
	}
	return models.Usuario{}, false
}
func (f *fakeUsuarioRepo) CrearUsuario(u models.Usuario) (models.Usuario, error) { return u, nil }
func (f *fakeUsuarioRepo) ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (f *fakeUsuarioRepo) BorrarUsuario(id uint) bool { return false }

// ─── SETUP COMÚN ────────────────────────────────────────────────────────────

// montarRouterPrueba construye un router chi con SOLO la ruta /guardavidas
// protegida por el middleware.Auth real (el mismo que usa cmd/piscina-api),
// apuntando a un SeguridadService respaldado por el fake en memoria.
// Devuelve también un token JWT válido, generado por el AuthService real.
func montarRouterPrueba(t *testing.T) (router chi.Router, tokenValido string) {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("clave123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("no se pudo generar el hash de prueba: %v", err)
	}
	usuarios := &fakeUsuarioRepo{
		usuario: models.Usuario{ID: 1, Nombre: "Admin Prueba", Email: "admin@prueba.com", PasswordHash: string(hash), Rol: "admin"},
	}
	authSvc := service.NewAuthService(usuarios)

	token, _, err := authSvc.Login("admin@prueba.com", "clave123")
	if err != nil {
		t.Fatalf("no se pudo generar el token de prueba: %v", err)
	}

	seguridadSvc := service.NewSeguridadService(&fakeSeguridadRepo{}, &fakeClienteRepo{}, &fakePagoRepo{})
	srv := NewServer(seguridadSvc, nil, nil, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Route("/guardavidas", func(r chi.Router) {
				r.Get("/", srv.ListarGuardavidas)
				r.Post("/", srv.CrearGuardavida)
			})
		})
	})

	return r, token
}

// ─── TESTS ──────────────────────────────────────────────────────────────────

// TestListarGuardavidas_SinToken_Devuelve401 prueba que una ruta protegida
// (GET /api/v1/guardavidas) sin header Authorization responde 401, sin
// llegar nunca al handler ni al service. Si alguien quitara el middleware
// de esta ruta por error, este test fallaría porque recibiría 200.
func TestListarGuardavidas_SinToken_Devuelve401(t *testing.T) {
	router, _ := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

// TestCrearGuardavida_ConToken_CreaYPersisteEnFake prueba el camino feliz:
// con un token válido, POST /api/v1/guardavidas crea el guardavida y este
// queda reflejado en el fake en memoria (lo confirmamos listándolo después).
func TestCrearGuardavida_ConToken_CreaYPersisteEnFake(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(models.Guardavida{Nombre: "Sofía Loor", Turno: "mañana"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var creado models.Guardavida
	if err := json.Unmarshal(rec.Body.Bytes(), &creado); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if creado.Nombre != "Sofía Loor" || creado.ID == 0 {
		t.Errorf("guardavida creado inesperado: %+v", creado)
	}

	// Confirmar que quedó reflejado: listar debe traerlo de vuelta.
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	var lista []models.Guardavida
	if err := json.Unmarshal(rec2.Body.Bytes(), &lista); err != nil {
		t.Fatalf("no se pudo decodificar la lista: %v", err)
	}
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 guardavida en la lista, se obtuvieron %d", len(lista))
	}
}
