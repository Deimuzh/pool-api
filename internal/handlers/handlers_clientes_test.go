package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"pool-api/internal/middleware"
	"pool-api/internal/models"
	"pool-api/internal/service"
	"pool-api/internal/storage"

	"github.com/go-chi/chi/v5"
)

// fakeClientesRepo implementa storage.ClientesModulo en memoria para pruebas
// de handlers de Cliente, Reserva y Pago sin necesidad de SQLite.
type fakeClientesRepo struct {
	clientes    map[uint]models.Cliente
	reservas    map[uint]models.Reserva
	pagos       map[uint]models.Pago
	siguienteID uint
}

func newFakeClientesRepo() *fakeClientesRepo {
	return &fakeClientesRepo{
		clientes:    make(map[uint]models.Cliente),
		reservas:    make(map[uint]models.Reserva),
		pagos:       make(map[uint]models.Pago),
		siguienteID: 0,
	}
}

// ─── Cliente ─────────────────────────────────────────────────────────────────

func (f *fakeClientesRepo) ListarClientes() []models.Cliente {
	lista := make([]models.Cliente, 0, len(f.clientes))
	for _, c := range f.clientes {
		lista = append(lista, c)
	}
	return lista
}

func (f *fakeClientesRepo) BuscarClientePorID(id uint) (models.Cliente, bool) {
	c, ok := f.clientes[id]
	return c, ok
}

func (f *fakeClientesRepo) CrearCliente(c models.Cliente) (models.Cliente, error) {
	for _, existing := range f.clientes {
		if existing.Cedula == c.Cedula {
			return models.Cliente{}, gorm.ErrDuplicatedKey
		}
	}
	f.siguienteID++
	c.ID = f.siguienteID
	f.clientes[c.ID] = c
	return c, nil
}

func (f *fakeClientesRepo) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	_, ok := f.clientes[id]
	if !ok {
		return models.Cliente{}, false
	}
	datos.ID = id
	f.clientes[id] = datos
	return datos, true
}

func (f *fakeClientesRepo) BorrarCliente(id uint) bool {
	_, ok := f.clientes[id]
	if !ok {
		return false
	}
	delete(f.clientes, id)
	return true
}

// ─── Reserva ─────────────────────────────────────────────────────────────────

func (f *fakeClientesRepo) ListarReservas() []models.Reserva {
	lista := make([]models.Reserva, 0, len(f.reservas))
	for _, r := range f.reservas {
		lista = append(lista, r)
	}
	return lista
}

func (f *fakeClientesRepo) BuscarReservaPorID(id uint) (models.Reserva, bool) {
	r, ok := f.reservas[id]
	return r, ok
}

func (f *fakeClientesRepo) CrearReserva(rv models.Reserva) (models.Reserva, error) {
	f.siguienteID++
	rv.ID = f.siguienteID
	f.reservas[rv.ID] = rv
	return rv, nil
}

func (f *fakeClientesRepo) ActualizarReserva(id uint, datos models.Reserva) (models.Reserva, bool) {
	_, ok := f.reservas[id]
	if !ok {
		return models.Reserva{}, false
	}
	datos.ID = id
	f.reservas[id] = datos
	return datos, true
}

func (f *fakeClientesRepo) BorrarReserva(id uint) bool {
	_, ok := f.reservas[id]
	if !ok {
		return false
	}
	delete(f.reservas, id)
	return true
}

// ─── Pago ────────────────────────────────────────────────────────────────────

func (f *fakeClientesRepo) ListarPagos() []models.Pago {
	lista := make([]models.Pago, 0, len(f.pagos))
	for _, p := range f.pagos {
		lista = append(lista, p)
	}
	return lista
}

func (f *fakeClientesRepo) BuscarPagoPorID(id uint) (models.Pago, bool) {
	p, ok := f.pagos[id]
	return p, ok
}

func (f *fakeClientesRepo) CrearPago(p models.Pago) (models.Pago, error) {
	f.siguienteID++
	p.ID = f.siguienteID
	f.pagos[p.ID] = p
	return p, nil
}

func (f *fakeClientesRepo) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	_, ok := f.pagos[id]
	if !ok {
		return models.Pago{}, false
	}
	datos.ID = id
	f.pagos[id] = datos
	return datos, true
}

func (f *fakeClientesRepo) BorrarPago(id uint) bool {
	_, ok := f.pagos[id]
	if !ok {
		return false
	}
	delete(f.pagos, id)
	return true
}

func (f *fakeClientesRepo) ClienteTienePagoEntrada(clienteID uint) bool {
	for _, p := range f.pagos {
		if p.ClienteID == clienteID {
			return true
		}
	}
	return false
}

var _ storage.ClientesModulo = (*fakeClientesRepo)(nil)

// ─── SETUP ───────────────────────────────────────────────────────────────────

func montarRouterClientesPrueba(t *testing.T) (router chi.Router, tokenValido string) {
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

	clientesRepo := newFakeClientesRepo()
	clientesRepo.clientes[1] = models.Cliente{ID: 1, Nombre: "Ana Reyes", Cedula: "1312345678", Membresia: "ninguna"}
	clientesRepo.clientes[2] = models.Cliente{ID: 2, Nombre: "Luis Mora", Cedula: "1398765432", Membresia: "mensual"}
	clientesRepo.reservas[1] = models.Reserva{ID: 1, ClienteID: 2, Duracion: 720, Estado: "pendiente"}
	clientesRepo.pagos[1] = models.Pago{ID: 1, ClienteID: 1, Monto: 50, Concepto: "membresia", Metodo: "efectivo"}
	clientesRepo.siguienteID = 2

	clientesSvc := service.NewClientesService(clientesRepo)
	srv := NewServer(nil, nil, clientesSvc, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Route("/clientes", func(r chi.Router) {
				r.Get("/", srv.ListarClientes)
				r.Post("/", srv.CrearCliente)
				r.Get("/{id}", srv.ObtenerCliente)
				r.Put("/{id}", srv.ActualizarCliente)
				r.Delete("/{id}", srv.BorrarCliente)
			})
			r.Route("/reservas", func(r chi.Router) {
				r.Get("/", srv.ListarReservas)
				r.Post("/", srv.CrearReserva)
				r.Get("/{id}", srv.ObtenerReserva)
				r.Put("/{id}", srv.ActualizarReserva)
				r.Delete("/{id}", srv.BorrarReserva)
			})
			r.Route("/pagos", func(r chi.Router) {
				r.Get("/", srv.ListarPagos)
				r.Post("/", srv.CrearPago)
				r.Get("/{id}", srv.ObtenerPago)
				r.Put("/{id}", srv.ActualizarPago)
				r.Delete("/{id}", srv.BorrarPago)
			})
		})
	})

	return r, token
}

// ─── AUTENTICACIÓN ──────────────────────────────────────────────────────────

func TestListarClientes_SinToken_Devuelve401(t *testing.T) {
	router, _ := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

// ─── LISTAR ─────────────────────────────────────────────────────────────────

func TestListarClientes_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var lista []models.Cliente
	if err := json.Unmarshal(rec.Body.Bytes(), &lista); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if len(lista) != 2 {
		t.Fatalf("se esperaban 2 clientes, se obtuvieron %d", len(lista))
	}
	nombres := make(map[string]bool)
	for _, c := range lista {
		nombres[c.Nombre] = true
	}
	if !nombres["Ana Reyes"] || !nombres["Luis Mora"] {
		t.Fatalf("nombres inesperados en la lista: %+v", lista)
	}
}
// ─── CREAR ──────────────────────────────────────────────────────────────────

func TestCrearCliente_ConToken_CreaYPersiste(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "Luis Pino", Cedula: "1311111111"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var creado models.Cliente
	if err := json.Unmarshal(rec.Body.Bytes(), &creado); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if creado.Nombre != "Luis Pino" || creado.ID == 0 {
		t.Fatalf("cliente creado inesperado: %+v", creado)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	var lista []models.Cliente
	json.Unmarshal(rec2.Body.Bytes(), &lista)
	if len(lista) != 3 {
		t.Fatalf("se esperaban 3 clientes, se obtuvieron %d", len(lista))
	}
}

func TestCrearCliente_JSONInvalido(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes/", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestCrearCliente_NombreVacio(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "", Cedula: "1311111111"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

// ─── OBTENER ────────────────────────────────────────────────────────────────

func TestObtenerCliente_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var c models.Cliente
	if err := json.Unmarshal(rec.Body.Bytes(), &c); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if c.Nombre != "Ana Reyes" {
		t.Fatalf("nombre inesperado: %s", c.Nombre)
	}
}

func TestObtenerCliente_IDInvalido(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestObtenerCliente_NoEncontrado(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// ─── ACTUALIZAR ─────────────────────────────────────────────────────────────

func TestActualizarCliente_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "Ana Actualizada", Cedula: "1312345678"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/1", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var c models.Cliente
	if err := json.Unmarshal(rec.Body.Bytes(), &c); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if c.Nombre != "Ana Actualizada" {
		t.Fatalf("nombre inesperado: %s", c.Nombre)
	}
}

func TestActualizarCliente_IDInvalido(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "Test", Cedula: "1311111111"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/no-numero", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarCliente_JSONInvalido(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/1", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarCliente_NoEncontrado(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "Test", Cedula: "1311111111"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/999", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// ─── BORRAR ─────────────────────────────────────────────────────────────────

func TestBorrarCliente_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/clientes/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["mensaje"] != "cliente eliminado" {
		t.Fatalf("mensaje inesperado: %s", resp["mensaje"])
	}
}

func TestBorrarCliente_IDInvalido(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/clientes/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarCliente_NoEncontrado(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/clientes/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// ─── ERRORES ADICIONALES ────────────────────────────────────────────────────

func TestCrearCliente_CedulaDuplicada_Devuelve409(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "Otro", Cedula: "1312345678"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("se esperaba 409, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCrearCliente_CedulaFormatoInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "Test", Cedula: "123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestActualizarCliente_NombreVacio_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "", Cedula: "1312345678"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/1", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

// ─── RESERVA ─────────────────────────────────────────────────────────────────

func TestListarReservas_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservas/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var lista []models.Reserva
	json.Unmarshal(rec.Body.Bytes(), &lista)
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 reserva, se obtuvieron %d", len(lista))
	}
}

func TestCrearReserva_ConToken_CreaYPersiste(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Reserva{ClienteID: 2, Duracion: 720, Estado: "pendiente"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var creada models.Reserva
	json.Unmarshal(rec.Body.Bytes(), &creada)
	if creada.ID == 0 || creada.ClienteID != 2 {
		t.Fatalf("reserva creada inesperada: %+v", creada)
	}
}

func TestObtenerReserva_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservas/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var rv models.Reserva
	json.Unmarshal(rec.Body.Bytes(), &rv)
	if rv.ID != 1 || rv.ClienteID != 2 {
		t.Fatalf("reserva inesperada: %+v", rv)
	}
}

func TestObtenerReserva_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservas/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestObtenerReserva_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservas/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestActualizarReserva_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Reserva{ClienteID: 2, Duracion: 720, Estado: "confirmada"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservas/1", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var rv models.Reserva
	json.Unmarshal(rec.Body.Bytes(), &rv)
	if rv.Estado != "confirmada" {
		t.Fatalf("estado inesperado: %s", rv.Estado)
	}
}

func TestBorrarReserva_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/reservas/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["mensaje"] != "reserva eliminada" {
		t.Fatalf("mensaje inesperado: %s", resp["mensaje"])
	}
}

// ─── PAGO ────────────────────────────────────────────────────────────────────

func TestListarPagos_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pagos/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var lista []models.Pago
	json.Unmarshal(rec.Body.Bytes(), &lista)
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 pago, se obtuvieron %d", len(lista))
	}
}

func TestCrearPago_ConToken_CreaYPersiste(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Pago{ClienteID: 1, Monto: 20, Concepto: "dia", Metodo: "efectivo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pagos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var creado models.Pago
	json.Unmarshal(rec.Body.Bytes(), &creado)
	if creado.ID == 0 || creado.ClienteID != 1 {
		t.Fatalf("pago creado inesperado: %+v", creado)
	}
}

func TestObtenerPago_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pagos/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var p models.Pago
	json.Unmarshal(rec.Body.Bytes(), &p)
	if p.ID != 1 || p.Monto != 50 {
		t.Fatalf("pago inesperado: %+v", p)
	}
}

func TestObtenerPago_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pagos/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestObtenerPago_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pagos/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestActualizarPago_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Pago{ClienteID: 1, Monto: 60, Concepto: "dia", Metodo: "transferencia"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/pagos/1", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var p models.Pago
	json.Unmarshal(rec.Body.Bytes(), &p)
	if p.Monto != 60 {
		t.Fatalf("monto inesperado: %f", p.Monto)
	}
}

func TestBorrarPago_ConToken_OK(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/pagos/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["mensaje"] != "pago eliminado" {
		t.Fatalf("mensaje inesperado: %s", resp["mensaje"])
	}
}

// ─── ERRORES RESERVA ─────────────────────────────────────────────────────────

func TestCrearReserva_JSONInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas/", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestCrearReserva_SinCliente_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Reserva{Duracion: 720})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestActualizarReserva_JSONInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservas/1", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarReserva_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Reserva{ClienteID: 2, Duracion: 720})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservas/no-numero", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarReserva_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Reserva{ClienteID: 2, Duracion: 720, Estado: "confirmada"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservas/999", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestActualizarPago_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Pago{ClienteID: 1, Monto: 20, Concepto: "dia", Metodo: "efectivo"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/pagos/999", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestBorrarReserva_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/reservas/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarReserva_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/reservas/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// ─── ERRORES PAGO ────────────────────────────────────────────────────────────

func TestCrearPago_JSONInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pagos/", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestCrearPago_ClienteInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Pago{ClienteID: 999, Monto: 20, Concepto: "dia", Metodo: "efectivo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pagos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestActualizarPago_JSONInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/pagos/1", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarPago_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Pago{ClienteID: 1, Monto: 20, Concepto: "dia", Metodo: "efectivo"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/pagos/no-numero", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarPago_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/pagos/no-numero", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBorrarPago_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/pagos/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}
