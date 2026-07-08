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

	seguridadRepo := &fakeSeguridadRepo{
		guardavidas: []models.Guardavida{
			{ID: 1, Nombre: "Carlos Mendoza", Turno: "mañana"},
		},
		incidentes: []models.Incidente{
			{ID: 1, Tipo: "lesion", Gravedad: "leve", GuardavidaID: 1, ClienteID: 2},
		},
		accesos: []models.AccesoCliente{
			{ID: 1, ClienteID: 2, Autorizado: true},
		},
		siguienteID: 1,
	}
	clientesRepo := &fakeClienteRepo{
		clientes: map[uint]models.Cliente{
			2: {ID: 2, Nombre: "Luis Pino", Membresia: "ninguna"},
			3: {ID: 3, Nombre: "Ana Reyes", Membresia: "mensual"},
		},
	}
	pagosRepo := &fakePagoRepo{tienePago: true}

	seguridadSvc := service.NewSeguridadService(seguridadRepo, clientesRepo, pagosRepo)

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
				r.Get("/{id}", srv.ObtenerIncidente)
				r.Post("/", srv.CrearIncidente)
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
	if len(lista) != 2 {
		t.Fatalf("se esperaban 2 guardavidas en la lista, se obtuvieron %d", len(lista))
	}
}

// TestObtenerGuardavida_ConToken_OK verifica que el handler pueda devolver
// un guardavida existente cuando la petición tiene JWT válido.
func TestObtenerGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

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
	if guardavida.Nombre != "Carlos Mendoza" {
		t.Fatalf("nombre inesperado: %s", guardavida.Nombre)
	}
}

// TestActualizarGuardavida_ConToken_OK verifica que el handler permita actualizar
// un guardavida existente cuando el JWT es válido.
func TestActualizarGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	bodyActualizar, _ := json.Marshal(models.Guardavida{Nombre: "Carlos Actualizado", Turno: "noche"})
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
	if guardavida.Nombre != "Carlos Actualizado" || guardavida.Turno != "noche" {
		t.Fatalf("guardavida inesperado: %+v", guardavida)
	}
}

// TestBorrarGuardavida_ConToken_OK verifica que el handler permita eliminar
// un guardavida existente y responda 200.
func TestBorrarGuardavida_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/guardavidas/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

// TestListarAccesos_ConToken_OK verifica que el handler devuelva la lista
// de accesos enriquecida cuando la petición tiene JWT válido.
func TestListarAccesos_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/accesos/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var accesos []service.AccesoConNombre
	if err := json.Unmarshal(rec.Body.Bytes(), &accesos); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if len(accesos) == 0 {
		t.Fatal("se esperaba al menos un acceso")
	}
	if accesos[0].NombreCliente != "Luis Pino" {
		t.Fatalf("nombre de cliente inesperado: %s", accesos[0].NombreCliente)
	}
}

// TestCrearAcceso_ConToken_OK verifica que el handler permita crear un acceso
// para un cliente válido con pago registrado.
func TestCrearAcceso_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(map[string]uint{"cliente_id": 2})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accesos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var acceso service.AccesoConNombre
	if err := json.Unmarshal(rec.Body.Bytes(), &acceso); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if acceso.ClienteID != 2 || !acceso.Autorizado {
		t.Fatalf("acceso inesperado: %+v", acceso)
	}
}

// TestBorrarAcceso_ConToken_OK verifica que el handler permita eliminar
// un acceso existente y responda 204.
func TestBorrarAcceso_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/accesos/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

// TestListarIncidentes_ConToken_OK verifica que el handler devuelva incidentes
// enriquecidos con nombres de guardavida y cliente.
func TestListarIncidentes_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/incidentes/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var incidentes []service.IncidenteConNombre
	if err := json.Unmarshal(rec.Body.Bytes(), &incidentes); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if len(incidentes) == 0 {
		t.Fatal("se esperaba al menos un incidente")
	}
	if incidentes[0].NombreGuardavida != "Carlos Mendoza" {
		t.Fatalf("guardavida inesperado: %s", incidentes[0].NombreGuardavida)
	}
}

// TestObtenerIncidente_ConToken_OK verifica que el handler devuelva un incidente
// existente por ID cuando la petición tiene JWT válido.
func TestObtenerIncidente_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/incidentes/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var incidente service.IncidenteConNombre
	if err := json.Unmarshal(rec.Body.Bytes(), &incidente); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if incidente.ID != 1 {
		t.Fatalf("id inesperado: %d", incidente.ID)
	}
}

// TestCrearIncidente_ConToken_OK verifica que el handler permita crear un
// incidente válido asociado a guardavida y cliente existentes.
func TestCrearIncidente_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(models.Incidente{
		Tipo:         "rescate",
		Gravedad:     "media",
		GuardavidaID: 1,
		ClienteID:    2,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/incidentes/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var incidente service.IncidenteConNombre
	if err := json.Unmarshal(rec.Body.Bytes(), &incidente); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if incidente.Tipo != "rescate" {
		t.Fatalf("tipo inesperado: %s", incidente.Tipo)
	}
}

// TestActualizarIncidente_ConToken_OK verifica que el handler permita actualizar
// un incidente existente.
func TestActualizarIncidente_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(models.Incidente{
		Tipo:         "rescate",
		Gravedad:     "alta",
		GuardavidaID: 1,
		ClienteID:    2,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/incidentes/1", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var incidente service.IncidenteConNombre
	if err := json.Unmarshal(rec.Body.Bytes(), &incidente); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if incidente.Gravedad != "alta" {
		t.Fatalf("gravedad inesperada: %s", incidente.Gravedad)
	}
}

// TestBorrarIncidente_ConToken_OK verifica que el handler permita eliminar
// un incidente existente y responda 204.
func TestBorrarIncidente_ConToken_OK(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/incidentes/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

// TestObtenerGuardavida_IDInvalido verifica que el handler responda 400
// cuando el parámetro id no es numérico.
func TestObtenerGuardavida_IDInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

// TestObtenerGuardavida_NoEncontrado verifica que el handler responda 404
// cuando el guardavida solicitado no existe.
func TestObtenerGuardavida_NoEncontrado(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/guardavidas/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// TestCrearGuardavida_JSONInvalido verifica que el handler responda 400
// cuando el cuerpo de la petición no es JSON válido.
func TestCrearGuardavida_JSONInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

// TestActualizarGuardavida_IDInvalido verifica que actualizar con id no numérico
// devuelva 400 antes de llegar al service.
func TestActualizarGuardavida_IDInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(models.Guardavida{Nombre: "Carlos", Turno: "noche"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/guardavidas/no-numero", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

// TestActualizarGuardavida_NoEncontrado verifica que actualizar un guardavida
// inexistente devuelva 404.
func TestActualizarGuardavida_NoEncontrado(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(models.Guardavida{Nombre: "Carlos", Turno: "noche"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/guardavidas/999", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// TestBorrarGuardavida_IDInvalido verifica que borrar con id no numérico
// devuelva 400.
func TestBorrarGuardavida_IDInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/guardavidas/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

// TestBorrarGuardavida_NoEncontrado verifica que borrar un guardavida inexistente
// devuelva 404.
func TestBorrarGuardavida_NoEncontrado(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/guardavidas/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// ─── INCIDENTES: ERROR ─────────────────────────────────────────────────────

func TestObtenerIncidente_IDInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/incidentes/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestObtenerIncidente_NoEncontrado(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/incidentes/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestCrearIncidente_JSONInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/incidentes/", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestCrearIncidente_CamposVacios(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(models.Incidente{Tipo: "", Gravedad: "leve", GuardavidaID: 1, ClienteID: 2})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/incidentes/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarIncidente_IDInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(models.Incidente{Tipo: "rescate", Gravedad: "alta", GuardavidaID: 1, ClienteID: 2})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/incidentes/no-numero", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarIncidente_JSONInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/incidentes/1", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarIncidente_NoEncontrado(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(models.Incidente{Tipo: "rescate", Gravedad: "alta", GuardavidaID: 1, ClienteID: 2})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/incidentes/999", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestBorrarIncidente_IDInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/incidentes/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarIncidente_NoEncontrado(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/incidentes/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// TestCrearGuardavida_NombreVacio verifica que crear un guardavida sin nombre
// falle con 400 porque el service devuelve ErrNombreVacio.
func TestCrearGuardavida_NombreVacio(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(models.Guardavida{Nombre: "", Turno: "mañana"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/guardavidas/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

// TestActualizarGuardavida_JSONInvalido verifica que actualizar con cuerpo
// no JSON devuelva 400.
func TestActualizarGuardavida_JSONInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/guardavidas/1", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

// TestActualizarGuardavida_NombreVacio verifica que actualizar un guardavida
// seteando nombre vacío falle con 400.
func TestActualizarGuardavida_NombreVacio(t *testing.T) {
	router, token := montarRouterPrueba(t)

	body, _ := json.Marshal(models.Guardavida{Nombre: "", Turno: "noche"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/guardavidas/1", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

// ─── ACCESOS: ERROR ────────────────────────────────────────────────────────

func TestCrearAcceso_JSONInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accesos/", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestCrearAcceso_ClienteInexistente(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(map[string]uint{"cliente_id": 999})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accesos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestCrearAcceso_ClienteConMembresia(t *testing.T) {
	router, token := montarRouterPrueba(t)
	body, _ := json.Marshal(map[string]uint{"cliente_id": 3})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accesos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestBorrarAcceso_IDInvalido(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/accesos/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarAcceso_NoEncontrado(t *testing.T) {
	router, token := montarRouterPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/accesos/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}
