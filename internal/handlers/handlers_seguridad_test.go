package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"pool-api/internal/middleware"
	"pool-api/internal/models"
	"pool-api/internal/service"
)

type fakeSeguridadRepo struct {
	guardavidas map[uint]models.Guardavida
	incidentes  map[uint]models.Incidente
	accesos     map[uint]models.AccesoCliente
	seq         uint
}

func newFakeSeguridadRepo() *fakeSeguridadRepo {
	return &fakeSeguridadRepo{
		guardavidas: make(map[uint]models.Guardavida),
		incidentes:  make(map[uint]models.Incidente),
		accesos:     make(map[uint]models.AccesoCliente),
		seq:         0,
	}
}

func (f *fakeSeguridadRepo) ListarGuardavidas() []models.Guardavida {
	lista := make([]models.Guardavida, 0, len(f.guardavidas))
	for _, g := range f.guardavidas {
		lista = append(lista, g)
	}
	return lista
}
func (f *fakeSeguridadRepo) BuscarGuardavidaPorID(id uint) (models.Guardavida, bool) {
	g, ok := f.guardavidas[id]
	return g, ok
}
func (f *fakeSeguridadRepo) CrearGuardavida(g models.Guardavida) models.Guardavida {
	f.seq++
	g.ID = f.seq
	f.guardavidas[g.ID] = g
	return g
}
func (f *fakeSeguridadRepo) ActualizarGuardavida(id uint, datos models.Guardavida) (models.Guardavida, bool) {
	_, ok := f.guardavidas[id]
	if !ok {
		return models.Guardavida{}, false
	}
	datos.ID = id
	f.guardavidas[id] = datos
	return datos, true
}
func (f *fakeSeguridadRepo) BorrarGuardavida(id uint) bool {
	_, ok := f.guardavidas[id]
	if !ok {
		return false
	}
	delete(f.guardavidas, id)
	return true
}
func (f *fakeSeguridadRepo) ListarIncidentes() []models.Incidente {
	lista := make([]models.Incidente, 0, len(f.incidentes))
	for _, i := range f.incidentes {
		lista = append(lista, i)
	}
	return lista
}
func (f *fakeSeguridadRepo) BuscarIncidentePorID(id uint) (models.Incidente, bool) {
	i, ok := f.incidentes[id]
	return i, ok
}
func (f *fakeSeguridadRepo) CrearIncidente(i models.Incidente) models.Incidente {
	f.seq++
	i.ID = f.seq
	f.incidentes[i.ID] = i
	return i
}
func (f *fakeSeguridadRepo) ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool) {
	_, ok := f.incidentes[id]
	if !ok {
		return models.Incidente{}, false
	}
	datos.ID = id
	f.incidentes[id] = datos
	return datos, true
}
func (f *fakeSeguridadRepo) BorrarIncidente(id uint) bool {
	_, ok := f.incidentes[id]
	if !ok {
		return false
	}
	delete(f.incidentes, id)
	return true
}
func (f *fakeSeguridadRepo) ListarAccesos() []models.AccesoCliente {
	lista := make([]models.AccesoCliente, 0, len(f.accesos))
	for _, a := range f.accesos {
		lista = append(lista, a)
	}
	return lista
}
func (f *fakeSeguridadRepo) BuscarAccesoPorID(id uint) (models.AccesoCliente, bool) {
	a, ok := f.accesos[id]
	return a, ok
}
func (f *fakeSeguridadRepo) CrearAcceso(a models.AccesoCliente) models.AccesoCliente {
	f.seq++
	a.ID = f.seq
	f.accesos[a.ID] = a
	return a
}
func (f *fakeSeguridadRepo) ActualizarAcceso(id uint, datos models.AccesoCliente) (models.AccesoCliente, bool) {
	_, ok := f.accesos[id]
	if !ok {
		return models.AccesoCliente{}, false
	}
	datos.ID = id
	f.accesos[id] = datos
	return datos, true
}
func (f *fakeSeguridadRepo) BorrarAcceso(id uint) bool {
	_, ok := f.accesos[id]
	if !ok {
		return false
	}
	delete(f.accesos, id)
	return true
}

type fakeClienteRepo struct{}

func (f *fakeClienteRepo) ListarClientes() []models.Cliente {
	return []models.Cliente{{ID: 1, Nombre: "Cliente Test", Cedula: "1311111111", Membresia: "mensual"}}
}
func (f *fakeClienteRepo) BuscarClientePorID(id uint) (models.Cliente, bool) {
	if id == 1 {
		return models.Cliente{ID: 1, Nombre: "Cliente Test", Cedula: "1311111111", Membresia: "mensual"}, true
	}
	return models.Cliente{}, false
}
func (f *fakeClienteRepo) CrearCliente(c models.Cliente) (models.Cliente, error) { return c, nil }
func (f *fakeClienteRepo) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (f *fakeClienteRepo) BorrarCliente(id uint) bool { return false }

type fakePagoRepo struct{}

func (f *fakePagoRepo) ListarPagos() []models.Pago                  { return nil }
func (f *fakePagoRepo) BuscarPagoPorID(id uint) (models.Pago, bool) { return models.Pago{}, false }
func (f *fakePagoRepo) CrearPago(p models.Pago) (models.Pago, error) { return p, nil }
func (f *fakePagoRepo) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}
func (f *fakePagoRepo) BorrarPago(id uint) bool                     { return false }
func (f *fakePagoRepo) ClienteTienePagoEntrada(clienteID uint) bool { return false }

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

	seguridadSvc := service.NewSeguridadService(newFakeSeguridadRepo(), &fakeClienteRepo{}, &fakePagoRepo{})
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
			r.Route("/incidentes", func(r chi.Router) {
				r.Get("/", srv.ListarIncidentes)
				r.Post("/", srv.CrearIncidente)
				r.Get("/{id}", srv.ObtenerIncidente)
				r.Put("/{id}", srv.ActualizarIncidente)
				r.Delete("/{id}", srv.BorrarIncidente)
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

func TestListarGuardavidas_SinToken_Devuelve401(t *testing.T) {
	router, _ := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

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
	json.Unmarshal(rec.Body.Bytes(), &creado)
	if creado.Nombre != "Sofía Loor" || creado.ID == 0 {
		t.Errorf("guardavida creado inesperado: %+v", creado)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	var lista []models.Guardavida
	json.Unmarshal(rec2.Body.Bytes(), &lista)
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 guardavida en la lista, se obtuvieron %d", len(lista))
	}
}

func TestObtenerGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(models.Guardavida{Nombre: "Luis", Turno: "tarde"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.Guardavida
	json.Unmarshal(rec.Body.Bytes(), &creado)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/"+strconv.Itoa(int(creado.ID)), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec2.Code)
	}
}

func TestObtenerGuardavida_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestActualizarGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(models.Guardavida{Nombre: "Ana", Turno: "mañana"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.Guardavida
	json.Unmarshal(rec.Body.Bytes(), &creado)

	body2, _ := json.Marshal(models.Guardavida{Nombre: "Ana Modificada", Turno: "noche"})
	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/guardavidas/"+strconv.Itoa(int(creado.ID)), bytes.NewReader(body2))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec2.Code, rec2.Body.String())
	}
}

func TestActualizarGuardavida_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(models.Guardavida{Nombre: "Test", Turno: "mañana"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/guardavidas/abc", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(models.Guardavida{Nombre: "Test", Turno: "mañana"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.Guardavida
	json.Unmarshal(rec.Body.Bytes(), &creado)

	req2 := httptest.NewRequest(http.MethodDelete, "/api/v1/guardavidas/"+strconv.Itoa(int(creado.ID)), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d", rec2.Code)
	}
}

func TestBorrarGuardavida_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/guardavidas/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestListarIncidentes_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/incidentes/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func crearGuardavidaEnRouter(t *testing.T, router chi.Router, token string) uint {
	t.Helper()
	body, _ := json.Marshal(models.Guardavida{Nombre: "Guardia Test", Turno: "mañana"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var g models.Guardavida
	json.Unmarshal(rec.Body.Bytes(), &g)
	return g.ID
}

func TestCrearIncidente_ConToken_Crea(t *testing.T) {
	router, token := montarRouterPrueba(t)
	gID := crearGuardavidaEnRouter(t, router, token)

	body, _ := json.Marshal(models.Incidente{Tipo: "lesion", Gravedad: "leve", GuardavidaID: gID, ClienteID: 1, Descripcion: "Corte menor"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/incidentes/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestObtenerIncidente_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/incidentes/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestListarAccesos_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/accesos/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}
