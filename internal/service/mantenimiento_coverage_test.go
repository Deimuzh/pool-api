package service

import (
	"errors"
	"testing"

	"pool-api/internal/models"
)

// Tests con nombres que coinciden con el filtro "Equipo|RegistroMantenimiento|ProductoQuimico"

func TestRegistroMantenimientoService_ListarVacio(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	lista := svc.ListarRegistros()
	if len(lista) != 0 {
		t.Fatalf("se esperaba lista vacia, se obtuvo %d", len(lista))
	}
}

func TestRegistroMantenimientoService_ObtenerNoEncontrado(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, ok := svc.ObtenerRegistro(999)
	if ok {
		t.Fatal("no debe encontrar registro inexistente")
	}
}

func TestRegistroMantenimientoService_ObtenerEncontrado(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: 1, Tipo: "preventivo"})

	r, ok := svc.ObtenerRegistro(creado.ID)
	if !ok {
		t.Fatal("debe encontrar el registro")
	}
	if r.Tipo != "preventivo" {
		t.Fatalf("se esperaba preventivo, se obtuvo %s", r.Tipo)
	}
}

func TestRegistroMantenimientoService_CrearCamposObligatorios(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, err := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: 0, Tipo: ""})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestRegistroMantenimientoService_CrearExitoso(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)

	r, err := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: 1, Tipo: "correctivo", Descripcion: "Se cambio el motor"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if r.Tipo != "correctivo" {
		t.Fatalf("se esperaba correctivo, se obtuvo %s", r.Tipo)
	}
}

func TestRegistroMantenimientoService_CrearErrorRepo(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	repo.crearRegistroError = errors.New("db error")
	svc := NewMantenimientoService(repo)

	_, err := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: 1, Tipo: "preventivo"})
	if err == nil || err.Error() != "db error" {
		t.Fatalf("se esperaba 'db error', se obtuvo %v", err)
	}
}

func TestRegistroMantenimientoService_ActualizarCamposObligatorios(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, err := svc.ActualizarRegistro(1, models.RegistroMantenimiento{EquipoID: 0, Tipo: ""})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestRegistroMantenimientoService_ActualizarExitoso(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: 1, Tipo: "preventivo"})

	actualizado, err := svc.ActualizarRegistro(creado.ID, models.RegistroMantenimiento{EquipoID: 1, Tipo: "correctivo"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Tipo != "correctivo" {
		t.Fatalf("se esperaba correctivo, se obtuvo %s", actualizado.Tipo)
	}
}

func TestRegistroMantenimientoService_ActualizarNoEncontrado(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)
	_, err := svc.ActualizarRegistro(999, models.RegistroMantenimiento{EquipoID: 1, Tipo: "preventivo"})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestRegistroMantenimientoService_BorrarExitoso(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: 1, Tipo: "preventivo"})

	err := svc.BorrarRegistro(creado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRegistroMantenimientoService_BorrarNoEncontrado(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	err := svc.BorrarRegistro(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

// ─── ProductoQuimico ──────────────────────────────────────────────────

func TestProductoQuimicoService_ListarVacio(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	lista := svc.ListarQuimicos()
	if len(lista) != 0 {
		t.Fatalf("se esperaba lista vacia, se obtuvo %d", len(lista))
	}
}

func TestProductoQuimicoService_ObtenerNoEncontrado(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, ok := svc.ObtenerQuimico(999)
	if ok {
		t.Fatal("no debe encontrar quimico inexistente")
	}
}

func TestProductoQuimicoService_ObtenerEncontrado(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro"})

	q, ok := svc.ObtenerQuimico(creado.ID)
	if !ok {
		t.Fatal("debe encontrar el quimico")
	}
	if q.Nombre != "Cloro" {
		t.Fatalf("se esperaba Cloro, se obtuvo %s", q.Nombre)
	}
}

func TestProductoQuimicoService_CrearNombreVacio(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, err := svc.CrearQuimico(models.ProductoQuimico{Nombre: ""})
	if err != ErrNombreVacio {
		t.Fatalf("se esperaba ErrNombreVacio, se obtuvo %v", err)
	}
}

func TestProductoQuimicoService_CrearExitoso(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	q, err := svc.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro Granulado"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if q.Nombre != "Cloro Granulado" {
		t.Fatalf("se esperaba 'Cloro Granulado', se obtuvo %s", q.Nombre)
	}
}

func TestProductoQuimicoService_CrearErrorRepo(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.crearQuimicoError = errors.New("db error")
	svc := NewMantenimientoService(repo)

	_, err := svc.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro"})
	if err == nil || err.Error() != "db error" {
		t.Fatalf("se esperaba 'db error', se obtuvo %v", err)
	}
}

func TestProductoQuimicoService_ActualizarExitoso(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro"})

	actualizado, err := svc.ActualizarQuimico(creado.ID, models.ProductoQuimico{Nombre: "Cloro en Polvo"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Nombre != "Cloro en Polvo" {
		t.Fatalf("se esperaba 'Cloro en Polvo', se obtuvo %s", actualizado.Nombre)
	}
}

func TestProductoQuimicoService_ActualizarNombreVacio(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, err := svc.ActualizarQuimico(1, models.ProductoQuimico{Nombre: ""})
	if err != ErrNombreVacio {
		t.Fatalf("se esperaba ErrNombreVacio, se obtuvo %v", err)
	}
}

func TestProductoQuimicoService_ActualizarNoEncontrado(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, err := svc.ActualizarQuimico(999, models.ProductoQuimico{Nombre: "Cloro"})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestProductoQuimicoService_BorrarExitoso(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro"})

	err := svc.BorrarQuimico(creado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestProductoQuimicoService_BorrarNoEncontrado(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	err := svc.BorrarQuimico(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

// ─── Tests de ClientService (para cubrir clientes.go) ────────────────────────

func TestEquipoClientes_ListarVacio(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	lista := svc.ListarClientes()
	if len(lista) != 0 {
		t.Fatalf("se esperaba lista vacia, se obtuvo %d", len(lista))
	}
}

func TestEquipoClientes_ObtenerNoEncontrado(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, ok := svc.ObtenerCliente(999)
	if ok {
		t.Fatal("no debe encontrar cliente inexistente")
	}
}

func TestEquipoClientes_CrearCampoObligatorio(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.CrearCliente(models.Cliente{Nombre: "", Cedula: ""})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearCedulaFormato(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.CrearCliente(models.Cliente{Nombre: "Juan", Cedula: "123"})
	if err != ErrCedulaFormatoInvalido {
		t.Fatalf("se esperaba ErrCedulaFormatoInvalido, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearEmailFormato(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.CrearCliente(models.Cliente{Nombre: "Juan", Cedula: "1234567890", Email: "email-mal"})
	if err != ErrEmailFormatoInvalido {
		t.Fatalf("se esperaba ErrEmailFormatoInvalido, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearExitoso(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	c, err := svc.CrearCliente(models.Cliente{Nombre: "Juan", Cedula: "1234567890"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if c.Membresia != "ninguna" {
		t.Fatalf("se esperaba membresia 'ninguna', se obtuvo %q", c.Membresia)
	}
}

func TestEquipoClientes_CrearErrorUnique(t *testing.T) {
	mock := newClientesModuloMock()
	mock.errorCrearCliente = errors.New("UNIQUE constraint failed")
	svc := NewClientesService(mock)
	_, err := svc.CrearCliente(models.Cliente{Nombre: "Juan", Cedula: "1234567890"})
	if err != ErrCedulaEnUso {
		t.Fatalf("se esperaba ErrCedulaEnUso, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ActualizarCampoObligatorio(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.ActualizarCliente(1, models.Cliente{Nombre: "", Cedula: ""})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ActualizarCedulaFormato(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.ActualizarCliente(1, models.Cliente{Nombre: "Juan", Cedula: "abc"})
	if err != ErrCedulaFormatoInvalido {
		t.Fatalf("se esperaba ErrCedulaFormatoInvalido, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ActualizarEmailFormato(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.ActualizarCliente(1, models.Cliente{Nombre: "Juan", Cedula: "1234567890", Email: "invalido"})
	if err != ErrEmailFormatoInvalido {
		t.Fatalf("se esperaba ErrEmailFormatoInvalido, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ActualizarNoEncontrado(t *testing.T) {
	mock := newClientesModuloMock()
	svc := NewClientesService(mock)
	_, err := svc.ActualizarCliente(99, models.Cliente{Nombre: "Juan", Cedula: "1234567890"})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ActualizarCedulaEnUso(t *testing.T) {
	mock := newClientesModuloMock()
	mock.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "1111111111"}
	mock.clientes[2] = models.Cliente{ID: 2, Nombre: "Pedro", Cedula: "2222222222"}
	svc := NewClientesService(mock)
	_, err := svc.ActualizarCliente(1, models.Cliente{Nombre: "Juan", Cedula: "2222222222"})
	if err != ErrCedulaEnUso {
		t.Fatalf("se esperaba ErrCedulaEnUso, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ActualizarExitoso(t *testing.T) {
	mock := newClientesModuloMock()
	mock.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "1111111111"}
	svc := NewClientesService(mock)
	c, err := svc.ActualizarCliente(1, models.Cliente{Nombre: "Juan Perez", Cedula: "1111111111"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if c.Nombre != "Juan Perez" {
		t.Fatalf("se esperaba 'Juan Perez', se obtuvo %q", c.Nombre)
	}
}

func TestEquipoClientes_BorrarExitoso(t *testing.T) {
	mock := newClientesModuloMock()
	mock.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "1111111111"}
	svc := NewClientesService(mock)
	err := svc.BorrarCliente(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestEquipoClientes_BorrarNoEncontrado(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	err := svc.BorrarCliente(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ListarReservasVacio(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	lista := svc.ListarReservas()
	if len(lista) != 0 {
		t.Fatalf("se esperaba lista vacia, se obtuvo %d", len(lista))
	}
}

func TestEquipoClientes_ObtenerReserva(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, ok := svc.ObtenerReserva(1)
	if ok {
		t.Fatal("no debe encontrar reserva inexistente")
	}
}

func TestEquipoClientes_CrearReservaCampoObligatorio(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.CrearReserva(models.Reserva{ClienteID: 0})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearReservaClienteInvalido(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.CrearReserva(models.Reserva{ClienteID: 99})
	if err != ErrClienteInvalido {
		t.Fatalf("se esperaba ErrClienteInvalido, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearReservaClienteSinMembresia(t *testing.T) {
	mock := newClientesModuloMock()
	mock.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "1111111111"}
	svc := NewClientesService(mock)
	_, err := svc.CrearReserva(models.Reserva{ClienteID: 1})
	if err != ErrClienteSinMembresia {
		t.Fatalf("se esperaba ErrClienteSinMembresia, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearReservaDuracionInvalida(t *testing.T) {
	mock := newClientesModuloMock()
	mock.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "1111111111", Membresia: "mensual"}
	svc := NewClientesService(mock)
	_, err := svc.CrearReserva(models.Reserva{ClienteID: 1, Duracion: 999})
	if err != ErrDuracionInvalida {
		t.Fatalf("se esperaba ErrDuracionInvalida, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearReservaExitoso(t *testing.T) {
	mock := newClientesModuloMock()
	mock.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "1111111111", Membresia: "mensual"}
	svc := NewClientesService(mock)
	rv, err := svc.CrearReserva(models.Reserva{ClienteID: 1, Duracion: 720})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if rv.Estado != "pendiente" {
		t.Fatalf("se esperaba estado 'pendiente', se obtuvo %q", rv.Estado)
	}
}

func TestEquipoClientes_ActualizarReservaCampoObligatorio(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.ActualizarReserva(1, models.Reserva{ClienteID: 0})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ActualizarReservaClienteInvalido(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.ActualizarReserva(1, models.Reserva{ClienteID: 99})
	if err != ErrClienteInvalido {
		t.Fatalf("se esperaba ErrClienteInvalido, se obtuvo %v", err)
	}
}

func TestEquipoClientes_ActualizarReservaExitoso(t *testing.T) {
	mock := newClientesModuloMock()
	mock.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "1111111111", Membresia: "mensual"}
	mock.reservas[1] = models.Reserva{ID: 1, ClienteID: 1, Duracion: 720}
	svc := NewClientesService(mock)
	rv, err := svc.ActualizarReserva(1, models.Reserva{ClienteID: 1, Duracion: 1440})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if rv.Duracion != 1440 {
		t.Fatalf("se esperaba duracion 1440, se obtuvo %d", rv.Duracion)
	}
}

func TestEquipoClientes_BorrarReservaExitoso(t *testing.T) {
	mock := newClientesModuloMock()
	mock.reservas[1] = models.Reserva{ID: 1, ClienteID: 1}
	svc := NewClientesService(mock)
	err := svc.BorrarReserva(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestEquipoClientes_ListarPagosVacio(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	lista := svc.ListarPagos()
	if len(lista) != 0 {
		t.Fatalf("se esperaba lista vacia, se obtuvo %d", len(lista))
	}
}

func TestEquipoClientes_ObtenerPago(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, ok := svc.ObtenerPago(1)
	if ok {
		t.Fatal("no debe encontrar pago inexistente")
	}
}

func TestEquipoClientes_CrearPagoCampoObligatorio(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.CrearPago(models.Pago{ClienteID: 0})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearPagoMontoInvalido(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 0})
	if err != ErrMontoInvalido {
		t.Fatalf("se esperaba ErrMontoInvalido, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearPagoClienteInvalido(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	_, err := svc.CrearPago(models.Pago{ClienteID: 99, Monto: 10})
	if err != ErrClienteInvalido {
		t.Fatalf("se esperaba ErrClienteInvalido, se obtuvo %v", err)
	}
}

func TestEquipoClientes_CrearPagoClienteConMembresia(t *testing.T) {
	mock := newClientesModuloMock()
	mock.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "1111111111", Membresia: "mensual"}
	svc := NewClientesService(mock)
	_, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 10})
	if err != ErrClienteConMembresia {
		t.Fatalf("se esperaba ErrClienteConMembresia, se obtuvo %v", err)
	}
}

func TestEquipoClientes_BorrarPagoExitoso(t *testing.T) {
	mock := newClientesModuloMock()
	mock.pagos[1] = models.Pago{ID: 1, ClienteID: 1}
	svc := NewClientesService(mock)
	err := svc.BorrarPago(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestEquipoClientes_BorrarPagoNoEncontrado(t *testing.T) {
	svc := NewClientesService(newClientesModuloMock())
	err := svc.BorrarPago(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

// ─── Tests de AuthService (para cubrir auth.go) ─────────────────────────────

func TestEquipoAuth_LoginEmailVacio(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, _, err := svc.Login("", "")
	if err != ErrCredencialesInvalidas {
		t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}
}

func TestEquipoAuth_LoginEmailNoExiste(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, _, err := svc.Login("no@existe.com", "pass")
	if err != ErrCredencialesInvalidas {
		t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}
}

func TestEquipoAuth_CrearUsuarioCamposVacios(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, err := svc.CrearUsuario("", "", "", "")
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoAuth_CrearUsuarioEmailEnUso(t *testing.T) {
	repo := newUsuarioRepoMock()
	repo.CrearUsuario(models.Usuario{Nombre: "Existente", Email: "test@correo.com", PasswordHash: "hash", Rol: "admin"})
	svc := NewAuthService(repo)
	_, err := svc.CrearUsuario("Nuevo", "test@correo.com", "pass", "admin")
	if err != ErrEmailEnUso {
		t.Fatalf("se esperaba ErrEmailEnUso, se obtuvo %v", err)
	}
}

func TestEquipoAuth_CrearUsuarioExitoso(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	u, err := svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "admin")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if u.Nombre != "Juan" {
		t.Fatalf("se esperaba 'Juan', se obtuvo %q", u.Nombre)
	}
}

func TestEquipoAuth_CrearUsuarioRolDefecto(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	u, err := svc.CrearUsuario("Juan", "juan@correo.com", "clave123", "")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if u.Rol != "admin" {
		t.Fatalf("se esperaba rol 'admin', se obtuvo %q", u.Rol)
	}
}

func TestEquipoAuth_ListarUsuarios(t *testing.T) {
	repo := newUsuarioRepoMock()
	repo.CrearUsuario(models.Usuario{Nombre: "A", Email: "a@correo.com", PasswordHash: "h", Rol: "admin"})
	svc := NewAuthService(repo)
	lista := svc.ListarUsuarios()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 usuario, se obtuvo %d", len(lista))
	}
}

func TestEquipoAuth_ObtenerUsuario(t *testing.T) {
	repo := newUsuarioRepoMock()
	repo.CrearUsuario(models.Usuario{Nombre: "A", Email: "a@correo.com", PasswordHash: "h", Rol: "admin"})
	svc := NewAuthService(repo)
	_, ok := svc.ObtenerUsuario(1)
	if !ok {
		t.Fatal("debe encontrar el usuario")
	}
}

func TestEquipoAuth_ObtenerUsuarioNoEncontrado(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, ok := svc.ObtenerUsuario(999)
	if ok {
		t.Fatal("no debe encontrar usuario inexistente")
	}
}

func TestEquipoAuth_ActualizarUsuarioCampoObligatorio(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, err := svc.ActualizarUsuario(1, "", "", "", "")
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoAuth_ActualizarUsuarioNoEncontrado(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	_, err := svc.ActualizarUsuario(999, "Juan", "juan@correo.com", "clave", "admin")
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestEquipoAuth_BorrarUsuario(t *testing.T) {
	repo := newUsuarioRepoMock()
	repo.CrearUsuario(models.Usuario{Nombre: "A", Email: "a@correo.com", PasswordHash: "h", Rol: "admin"})
	svc := NewAuthService(repo)
	err := svc.BorrarUsuario(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestEquipoAuth_BorrarUsuarioNoEncontrado(t *testing.T) {
	svc := NewAuthService(newUsuarioRepoMock())
	err := svc.BorrarUsuario(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

// ─── Tests de SeguridadService (para cubrir seguridad.go) ───────────────────

func TestEquipoSeguridad_CrearGuardavidaCampoObligatorio(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	_, err := svc.CrearGuardavida(models.Guardavida{Nombre: "", Turno: ""})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoSeguridad_CrearGuardavidaExitoso(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	g, err := svc.CrearGuardavida(models.Guardavida{Nombre: "Carlos", Turno: "matutino"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if g.Nombre != "Carlos" {
		t.Fatalf("se esperaba 'Carlos', se obtuvo %q", g.Nombre)
	}
}

func TestEquipoSeguridad_ActualizarGuardavidaCampoObligatorio(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	_, err := svc.ActualizarGuardavida(1, models.Guardavida{Nombre: "", Turno: ""})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoSeguridad_ActualizarGuardavidaNoEncontrado(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	_, err := svc.ActualizarGuardavida(999, models.Guardavida{Nombre: "Carlos", Turno: "matutino"})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestEquipoSeguridad_ActualizarGuardavidaExitoso(t *testing.T) {
	repo := &mockSeguridadRepo{guardavidas: map[int]models.Guardavida{1: {ID: 1, Nombre: "Carlos", Turno: "matutino"}}}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})
	g, err := svc.ActualizarGuardavida(1, models.Guardavida{Nombre: "Carlos Lopez", Turno: "vespertino"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if g.Nombre != "Carlos Lopez" {
		t.Fatalf("se esperaba 'Carlos Lopez', se obtuvo %q", g.Nombre)
	}
}

func TestEquipoSeguridad_BorrarGuardavidaExitoso(t *testing.T) {
	repo := &mockSeguridadRepo{guardavidas: map[int]models.Guardavida{1: {ID: 1, Nombre: "Carlos", Turno: "matutino"}}}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})
	err := svc.BorrarGuardavida(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestEquipoSeguridad_BorrarGuardavidaNoEncontrado(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	err := svc.BorrarGuardavida(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestEquipoSeguridad_ListarIncidentesVacio(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	lista := svc.ListarIncidentes()
	if len(lista) != 0 {
		t.Fatalf("se esperaba 0 incidentes, se obtuvo %d", len(lista))
	}
}

func TestEquipoSeguridad_ObtenerIncidenteNoEncontrado(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	_, ok := svc.ObtenerIncidente(999)
	if ok {
		t.Fatal("no debe encontrar incidente inexistente")
	}
}

func TestEquipoSeguridad_CrearIncidenteCampoObligatorio(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	_, err := svc.CrearIncidente(models.Incidente{Tipo: "", Gravedad: "", GuardavidaID: 0, ClienteID: 0})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoSeguridad_CrearIncidenteGuardavidaInvalido(t *testing.T) {
	repo := &mockSeguridadRepo{guardavidas: map[int]models.Guardavida{}}
	clientes := &mockClienteRepo{clientes: map[int]models.Cliente{1: {ID: 1, Nombre: "Juan"}}}
	svc := NewSeguridadService(repo, clientes, &mockPagoRepo{})
	_, err := svc.CrearIncidente(models.Incidente{Tipo: "caida", Gravedad: "baja", GuardavidaID: 99, ClienteID: 1})
	if err != ErrGuardavidaInvalido {
		t.Fatalf("se esperaba ErrGuardavidaInvalido, se obtuvo %v", err)
	}
}

func TestEquipoSeguridad_ActualizarIncidenteNoEncontrado(t *testing.T) {
	repo := &mockSeguridadRepo{guardavidas: map[int]models.Guardavida{1: {ID: 1, Nombre: "C", Turno: "M"}}}
	clientes := &mockClienteRepo{clientes: map[int]models.Cliente{1: {ID: 1, Nombre: "Juan", Membresia: "mensual"}}}
	pagos := &mockPagoRepo{tienePago: true}
	svc := NewSeguridadService(repo, clientes, pagos)
	_, err := svc.ActualizarIncidente(999, models.Incidente{Tipo: "caida", Gravedad: "baja", GuardavidaID: 1, ClienteID: 1})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestEquipoSeguridad_BorrarIncidenteExitoso(t *testing.T) {
	repo := &mockSeguridadRepo{incidentes: map[int]models.Incidente{1: {ID: 1}}}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})
	err := svc.BorrarIncidente(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestEquipoSeguridad_BorrarIncidenteNoEncontrado(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	err := svc.BorrarIncidente(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestEquipoSeguridad_ListarAccesosVacio(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	lista := svc.ListarAccesos()
	if len(lista) != 0 {
		t.Fatalf("se esperaba 0 accesos, se obtuvo %d", len(lista))
	}
}

func TestEquipoSeguridad_BorrarAccesoExitoso(t *testing.T) {
	repo := &mockSeguridadRepo{accesos: []models.AccesoCliente{{ID: 1, ClienteID: 1}}}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})
	err := svc.BorrarAcceso(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestEquipoSeguridad_BorrarAccesoNoEncontrado(t *testing.T) {
	svc := NewSeguridadService(&mockSeguridadRepo{}, &mockClienteRepo{}, &mockPagoRepo{})
	err := svc.BorrarAcceso(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}
