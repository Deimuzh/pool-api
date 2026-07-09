package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"pool-api/internal/models"
	"pool-api/internal/service"
	"pool-api/internal/storage"
)

type clientesHandlerMock struct{
	clientes map[uint]models.Cliente
}

var _ storage.ClientesModulo = (*clientesHandlerMock)(nil)

func newClientesHandlerMock() *clientesHandlerMock {
	return &clientesHandlerMock{
	clientes: map[uint]models.Cliente{
		1: {ID: 1, Nombre: "Luis", Cedula: "0000", Membresia: "ninguna"},
		2: {ID: 2, Nombre: "Ana", Cedula: "1111", Membresia: "mensual"},
	},
	}
}

func (m *clientesHandlerMock) ListarClientes() []models.Cliente { return nil }
func (m *clientesHandlerMock) BuscarClientePorID(id uint) (models.Cliente, bool) {
	c, ok := m.clientes[id]
	return c, ok
}
func (m *clientesHandlerMock) CrearCliente(c models.Cliente) (models.Cliente, error) { return c, nil }
func (m *clientesHandlerMock) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (m *clientesHandlerMock) BorrarCliente(id uint) bool { return false }
func (m *clientesHandlerMock) ListarReservas() []models.Reserva { return nil }
func (m *clientesHandlerMock) BuscarReservaPorID(id uint) (models.Reserva, bool) {
	return models.Reserva{}, false
}
func (m *clientesHandlerMock) CrearReserva(rv models.Reserva) (models.Reserva, error) { return rv, nil }
func (m *clientesHandlerMock) ActualizarReserva(id uint, datos models.Reserva) (models.Reserva, bool) {
	return models.Reserva{}, false
}
func (m *clientesHandlerMock) BorrarReserva(id uint) bool { return false }
func (m *clientesHandlerMock) ListarPagos() []models.Pago { return nil }
func (m *clientesHandlerMock) BuscarPagoPorID(id uint) (models.Pago, bool) {
	return models.Pago{}, false
}
func (m *clientesHandlerMock) CrearPago(p models.Pago) (models.Pago, error) { return p, nil }
func (m *clientesHandlerMock) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}
func (m *clientesHandlerMock) BorrarPago(id uint) bool { return false }
func (m *clientesHandlerMock) ClienteTienePagoEntrada(clienteID uint) bool { return false }

func montarRouterClientesPrueba() *chi.Mux {
	mock := newClientesHandlerMock()
	svc := service.NewClientesService(mock)
	srv := NewServer(nil, nil, svc, nil)
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
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
	return r
}

func TestListarClientes_Handler(t *testing.T) {
	router := montarRouterClientesPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestCrearCliente_Handler(t *testing.T) {
	router := montarRouterClientesPrueba()
	body, _ := json.Marshal(models.Cliente{Nombre: "Maria", Cedula: "1234567890"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestObtenerCliente_Handler_NoEncontrado(t *testing.T) {
	router := montarRouterClientesPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestListarReservas_Handler(t *testing.T) {
	router := montarRouterClientesPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservas/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestCrearReserva_Handler(t *testing.T) {
	router := montarRouterClientesPrueba()
	body, _ := json.Marshal(models.Reserva{ClienteID: 2, Duracion: 720})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestListarPagos_Handler(t *testing.T) {
	router := montarRouterClientesPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/pagos/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestCrearPago_Handler(t *testing.T) {
	router := montarRouterClientesPrueba()
	body, _ := json.Marshal(models.Pago{ClienteID: 1, Monto: 5, Concepto: "dia"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pagos/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}
