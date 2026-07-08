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

// fakeSeguridadRepo es un fake en memoria del storage.SeguridadRepository:
// guarda los guardavidas en un slice real en vez de una base de datos, para
// poder probar el handler de punta a punta sin SQLite.
type fakeSeguridadRepo struct {
	guardavidas []models.Guardavida
	incidentes  []models.Incidente
	accesos     []models.AccesoCliente
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
	for i := range f.guardavidas {
		if f.guardavidas[i].ID == id {
			datos.ID = id
			f.guardavidas[i] = datos
			return datos, true
		}
	}
	return models.Guardavida{}, false
}
func (f *fakeSeguridadRepo) BorrarGuardavida(id uint) bool {
	for i := range f.guardavidas {
		if f.guardavidas[i].ID == id {
			f.guardavidas = append(f.guardavidas[:i], f.guardavidas[i+1:]...)
			return true
		}
	}
	return false
}

func (f *fakeSeguridadRepo) ListarIncidentes() []models.Incidente {
	return f.incidentes
}

func (f *fakeSeguridadRepo) BuscarIncidentePorID(id uint) (models.Incidente, bool) {
	for _, inc := range f.incidentes {
		if inc.ID == id {
			return inc, true
		}
	}
	return models.Incidente{}, false
}

func (f *fakeSeguridadRepo) CrearIncidente(i models.Incidente) models.Incidente {
	f.siguienteID++
	i.ID = f.siguienteID
	f.incidentes = append(f.incidentes, i)
	return i
}

func (f *fakeSeguridadRepo) ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool) {
	for i := range f.incidentes {
		if f.incidentes[i].ID == id {
			datos.ID = id
			f.incidentes[i] = datos
			return datos, true
		}
	}
	return models.Incidente{}, false
}
func (f *fakeSeguridadRepo) BorrarIncidente(id uint) bool {
	for i := range f.incidentes {
		if f.incidentes[i].ID == id {
			f.incidentes = append(f.incidentes[:i], f.incidentes[i+1:]...)
			return true
		}
	}
	return false
}

func (f *fakeSeguridadRepo) ListarAccesos() []models.AccesoCliente {
	return f.accesos
}

func (f *fakeSeguridadRepo) BuscarAccesoPorID(id uint) (models.AccesoCliente, bool) {
	for _, a := range f.accesos {
		if a.ID == id {
			return a, true
		}
	}
	return models.AccesoCliente{}, false
}

func (f *fakeSeguridadRepo) CrearAcceso(a models.AccesoCliente) models.AccesoCliente {
	f.siguienteID++
	a.ID = f.siguienteID
	f.accesos = append(f.accesos, a)
	return a
}

func (f *fakeSeguridadRepo) ActualizarAcceso(id uint, datos models.AccesoCliente) (models.AccesoCliente, bool) {
	for i := range f.accesos {
		if f.accesos[i].ID == id {
			datos.ID = id
			f.accesos[i] = datos
			return datos, true
		}
	}
	return models.AccesoCliente{}, false
}

func (f *fakeSeguridadRepo) BorrarAcceso(id uint) bool {
	for i := range f.accesos {
		if f.accesos[i].ID == id {
			f.accesos = append(f.accesos[:i], f.accesos[i+1:]...)
			return true
		}
	}
	return false
}

// fakeClienteRepo y fakePagoRepo: SeguridadService los necesita para
// construirse, pero estos tests de Guardavida no los usan, así que basta
// con que implementen la interfaz sin hacer nada relevante.
type fakeClienteRepo struct {
	clientes map[uint]models.Cliente
}

func (f *fakeClienteRepo) ListarClientes() []models.Cliente { return nil }

func (f *fakeClienteRepo) BuscarClientePorID(id uint) (models.Cliente, bool) {
	c, ok := f.clientes[id]
	return c, ok
}

func (f *fakeClienteRepo) CrearCliente(c models.Cliente) (models.Cliente, error) { return c, nil }

func (f *fakeClienteRepo) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	return models.Cliente{}, false
}

func (f *fakeClienteRepo) BorrarCliente(id uint) bool { return false }

type fakePagoRepo struct {
	tienePago bool
}

func (f *fakePagoRepo) ListarPagos() []models.Pago { return nil }

func (f *fakePagoRepo) BuscarPagoPorID(id uint) (models.Pago, bool) { return models.Pago{}, false }

func (f *fakePagoRepo) CrearPago(p models.Pago) (models.Pago, error) { return p, nil }

func (f *fakePagoRepo) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}

func (f *fakePagoRepo) BorrarPago(id uint) bool { return false }

func (f *fakePagoRepo) ClienteTienePagoEntrada(clienteID uint) bool {
	return f.tienePago
}

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
				r.Get("/{id}", srv.ObtenerGuardavida)
				r.Put("/{id}", srv.ActualizarGuardavida)
				r.Delete("/{id}", srv.BorrarGuardavida)
			})

			r.Route("/accesos", func(r chi.Router) {
				r.Get("/", srv.ListarAccesos)
				r.Post("/", srv.CrearAcceso)
				r.Delete("/{id}", srv.BorrarAcceso)
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

// TestObtenerGuardavida_ConToken_OK verifica que el handler pueda devolver
// un guardavida existente cuando la petición tiene JWT válido.
func TestObtenerGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(models.Guardavida{Nombre: "Sofía Loor", Turno: "mañana"})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(body))
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	reqCrear.Header.Set("Content-Type", "application/json")
	recCrear := httptest.NewRecorder()
	router.ServeHTTP(recCrear, reqCrear)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	var guardavida models.Guardavida
	if err := json.Unmarshal(rec.Body.Bytes(), &guardavida); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if guardavida.Nombre != "Sofía Loor" {
		t.Fatalf("nombre inesperado: %s", guardavida.Nombre)
	}
}

// TestActualizarGuardavida_ConToken_OK verifica que el handler permita actualizar
// un guardavida existente cuando el JWT es válido.
func TestActualizarGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	bodyCrear, _ := json.Marshal(models.Guardavida{Nombre: "Sofía Loor", Turno: "mañana"})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(bodyCrear))
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	reqCrear.Header.Set("Content-Type", "application/json")
	recCrear := httptest.NewRecorder()
	router.ServeHTTP(recCrear, reqCrear)

	bodyActualizar, _ := json.Marshal(models.Guardavida{Nombre: "Sofía Actualizada", Turno: "noche"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/guardavidas/1", bytes.NewReader(bodyActualizar))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var guardavida models.Guardavida
	if err := json.Unmarshal(rec.Body.Bytes(), &guardavida); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if guardavida.Nombre != "Sofía Actualizada" || guardavida.Turno != "noche" {
		t.Fatalf("guardavida inesperado: %+v", guardavida)
	}
}

// TestBorrarGuardavida_ConToken_OK verifica que el handler permita eliminar
// un guardavida existente y responda 200.
func TestBorrarGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	bodyCrear, _ := json.Marshal(models.Guardavida{Nombre: "Sofía Loor", Turno: "mañana"})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(bodyCrear))
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	reqCrear.Header.Set("Content-Type", "application/json")
	recCrear := httptest.NewRecorder()
	router.ServeHTTP(recCrear, reqCrear)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/guardavidas/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}
