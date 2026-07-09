package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"pool-api/internal/handlers"
	"pool-api/internal/models"
	"pool-api/internal/service"
	"pool-api/internal/storage"
)

type routerMockClientes struct{}

var _ storage.ClientesModulo = (*routerMockClientes)(nil)

func (m *routerMockClientes) ListarClientes() []models.Cliente { return nil }
func (m *routerMockClientes) BuscarClientePorID(id uint) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (m *routerMockClientes) CrearCliente(c models.Cliente) (models.Cliente, error) { return c, nil }
func (m *routerMockClientes) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (m *routerMockClientes) BorrarCliente(id uint) bool { return false }
func (m *routerMockClientes) ListarReservas() []models.Reserva { return nil }
func (m *routerMockClientes) BuscarReservaPorID(id uint) (models.Reserva, bool) {
	return models.Reserva{}, false
}
func (m *routerMockClientes) CrearReserva(rv models.Reserva) (models.Reserva, error) { return rv, nil }
func (m *routerMockClientes) ActualizarReserva(id uint, datos models.Reserva) (models.Reserva, bool) {
	return models.Reserva{}, false
}
func (m *routerMockClientes) BorrarReserva(id uint) bool { return false }
func (m *routerMockClientes) ListarPagos() []models.Pago { return nil }
func (m *routerMockClientes) BuscarPagoPorID(id uint) (models.Pago, bool) {
	return models.Pago{}, false
}
func (m *routerMockClientes) CrearPago(p models.Pago) (models.Pago, error) { return p, nil }
func (m *routerMockClientes) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}
func (m *routerMockClientes) BorrarPago(id uint) bool { return false }
func (m *routerMockClientes) ClienteTienePagoEntrada(clienteID uint) bool { return false }

type routerMockSeguridad struct{}

func (m *routerMockSeguridad) ListarGuardavidas() []models.Guardavida { return nil }
func (m *routerMockSeguridad) BuscarGuardavidaPorID(id uint) (models.Guardavida, bool) {
	return models.Guardavida{}, false
}
func (m *routerMockSeguridad) CrearGuardavida(g models.Guardavida) models.Guardavida { return g }
func (m *routerMockSeguridad) ActualizarGuardavida(id uint, datos models.Guardavida) (models.Guardavida, bool) {
	return models.Guardavida{}, false
}
func (m *routerMockSeguridad) BorrarGuardavida(id uint) bool { return false }
func (m *routerMockSeguridad) ListarIncidentes() []models.Incidente { return nil }
func (m *routerMockSeguridad) BuscarIncidentePorID(id uint) (models.Incidente, bool) {
	return models.Incidente{}, false
}
func (m *routerMockSeguridad) CrearIncidente(i models.Incidente) models.Incidente { return i }
func (m *routerMockSeguridad) ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool) {
	return models.Incidente{}, false
}
func (m *routerMockSeguridad) BorrarIncidente(id uint) bool { return false }
func (m *routerMockSeguridad) ListarAccesos() []models.AccesoCliente { return nil }
func (m *routerMockSeguridad) BuscarAccesoPorID(id uint) (models.AccesoCliente, bool) {
	return models.AccesoCliente{}, false
}
func (m *routerMockSeguridad) CrearAcceso(a models.AccesoCliente) models.AccesoCliente { return a }
func (m *routerMockSeguridad) ActualizarAcceso(id uint, datos models.AccesoCliente) (models.AccesoCliente, bool) {
	return models.AccesoCliente{}, false
}
func (m *routerMockSeguridad) BorrarAcceso(id uint) bool { return false }

type routerMockMantenimiento struct{}

func (m *routerMockMantenimiento) ListarEquipos() []models.Equipo { return nil }
func (m *routerMockMantenimiento) BuscarEquipoPorID(id uint) (models.Equipo, bool) {
	return models.Equipo{}, false
}
func (m *routerMockMantenimiento) CrearEquipo(e models.Equipo) (models.Equipo, error) { return e, nil }
func (m *routerMockMantenimiento) ActualizarEquipo(id uint, datos models.Equipo) (models.Equipo, bool) {
	return models.Equipo{}, false
}
func (m *routerMockMantenimiento) BorrarEquipo(id uint) bool { return false }
func (m *routerMockMantenimiento) ListarRegistros() []models.RegistroMantenimiento { return nil }
func (m *routerMockMantenimiento) BuscarRegistroPorID(id uint) (models.RegistroMantenimiento, bool) {
	return models.RegistroMantenimiento{}, false
}
func (m *routerMockMantenimiento) CrearRegistro(rm models.RegistroMantenimiento) (models.RegistroMantenimiento, error) {
	return rm, nil
}
func (m *routerMockMantenimiento) ActualizarRegistro(id uint, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool) {
	return models.RegistroMantenimiento{}, false
}
func (m *routerMockMantenimiento) BorrarRegistro(id uint) bool { return false }
func (m *routerMockMantenimiento) ListarQuimicos() []models.ProductoQuimico { return nil }
func (m *routerMockMantenimiento) BuscarQuimicoPorID(id uint) (models.ProductoQuimico, bool) {
	return models.ProductoQuimico{}, false
}
func (m *routerMockMantenimiento) CrearQuimico(q models.ProductoQuimico) (models.ProductoQuimico, error) {
	return q, nil
}
func (m *routerMockMantenimiento) ActualizarQuimico(id uint, datos models.ProductoQuimico) (models.ProductoQuimico, bool) {
	return models.ProductoQuimico{}, false
}
func (m *routerMockMantenimiento) BorrarQuimico(id uint) bool { return false }

type routerMockUsuario struct{}

func (m *routerMockUsuario) ListarUsuarios() []models.Usuario { return nil }
func (m *routerMockUsuario) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (m *routerMockUsuario) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (m *routerMockUsuario) CrearUsuario(u models.Usuario) (models.Usuario, error) { return u, nil }
func (m *routerMockUsuario) ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (m *routerMockUsuario) BorrarUsuario(id uint) bool { return false }

func TestNuevoRouter_Raiz(t *testing.T) {
	seguridadSvc := service.NewSeguridadService(&routerMockSeguridad{}, &routerMockClientes{}, &routerMockClientes{})
	mantenimientoSvc := service.NewMantenimientoService(&routerMockMantenimiento{})
	clientesSvc := service.NewClientesService(&routerMockClientes{})
	authSvc := service.NewAuthService(&routerMockUsuario{})
	srv := handlers.NewServer(seguridadSvc, mantenimientoSvc, clientesSvc, authSvc)

	r := NuevoRouter(srv)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200 en raiz, se obtuvo %d", rec.Code)
	}
}

func TestNuevoRouter_APIClientes(t *testing.T) {
	seguridadSvc := service.NewSeguridadService(&routerMockSeguridad{}, &routerMockClientes{}, &routerMockClientes{})
	mantenimientoSvc := service.NewMantenimientoService(&routerMockMantenimiento{})
	clientesSvc := service.NewClientesService(&routerMockClientes{})
	authSvc := service.NewAuthService(&routerMockUsuario{})
	srv := handlers.NewServer(seguridadSvc, mantenimientoSvc, clientesSvc, authSvc)

	r := NuevoRouter(srv)

	tests := []struct {
		method string
		path   string
		code   int
	}{
		{"GET", "/api/v1/clientes", 200},
		{"GET", "/api/v1/equipos", 200},
		{"GET", "/api/v1/reservas", 200},
		{"GET", "/api/v1/pagos", 200},
		{"GET", "/api/v1/guardavidas", 200},
		{"GET", "/api/v1/incidentes", 200},
		{"GET", "/api/v1/accesos", 200},
		{"GET", "/api/v1/mantenimientos", 200},
		{"GET", "/api/v1/quimicos", 200},
	}
	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code != tc.code {
				t.Errorf("%s %s: esperado %d, obtuve %d - body: %s", tc.method, tc.path, tc.code, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestNuevoRouter_404(t *testing.T) {
	seguridadSvc := service.NewSeguridadService(&routerMockSeguridad{}, &routerMockClientes{}, &routerMockClientes{})
	mantenimientoSvc := service.NewMantenimientoService(&routerMockMantenimiento{})
	clientesSvc := service.NewClientesService(&routerMockClientes{})
	authSvc := service.NewAuthService(&routerMockUsuario{})
	srv := handlers.NewServer(seguridadSvc, mantenimientoSvc, clientesSvc, authSvc)

	r := NuevoRouter(srv)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/no-existe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}
