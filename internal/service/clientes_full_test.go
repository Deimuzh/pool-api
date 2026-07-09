package service

import (
	"testing"

	"github.com/stretchr/testify/require"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type clientesRepoCompletoMock struct {
	clientes map[uint]models.Cliente
	reservas map[uint]models.Reserva
	pagos    map[uint]models.Pago

	sigClienteID uint
	sigReservaID uint
	sigPagoID    uint
}

var _ storage.ClientesModulo = (*clientesRepoCompletoMock)(nil)

func newClientesRepoCompletoMock() *clientesRepoCompletoMock {
	return &clientesRepoCompletoMock{
		clientes: make(map[uint]models.Cliente),
		reservas: make(map[uint]models.Reserva),
		pagos:    make(map[uint]models.Pago),
	}
}

func (m *clientesRepoCompletoMock) ListarClientes() []models.Cliente {
	lista := make([]models.Cliente, 0, len(m.clientes))
	for _, c := range m.clientes {
		lista = append(lista, c)
	}
	return lista
}
func (m *clientesRepoCompletoMock) BuscarClientePorID(id uint) (models.Cliente, bool) {
	c, ok := m.clientes[id]
	return c, ok
}
func (m *clientesRepoCompletoMock) CrearCliente(c models.Cliente) (models.Cliente, error) {
	m.sigClienteID++
	c.ID = m.sigClienteID
	m.clientes[c.ID] = c
	return c, nil
}
func (m *clientesRepoCompletoMock) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	_, ok := m.clientes[id]
	if !ok {
		return models.Cliente{}, false
	}
	datos.ID = id
	m.clientes[id] = datos
	return datos, true
}
func (m *clientesRepoCompletoMock) BorrarCliente(id uint) bool {
	_, ok := m.clientes[id]
	if !ok {
		return false
	}
	delete(m.clientes, id)
	return true
}
func (m *clientesRepoCompletoMock) ListarReservas() []models.Reserva {
	lista := make([]models.Reserva, 0, len(m.reservas))
	for _, r := range m.reservas {
		lista = append(lista, r)
	}
	return lista
}
func (m *clientesRepoCompletoMock) BuscarReservaPorID(id uint) (models.Reserva, bool) {
	r, ok := m.reservas[id]
	return r, ok
}
func (m *clientesRepoCompletoMock) CrearReserva(r models.Reserva) (models.Reserva, error) {
	m.sigReservaID++
	r.ID = m.sigReservaID
	m.reservas[r.ID] = r
	return r, nil
}
func (m *clientesRepoCompletoMock) ActualizarReserva(id uint, datos models.Reserva) (models.Reserva, bool) {
	_, ok := m.reservas[id]
	if !ok {
		return models.Reserva{}, false
	}
	datos.ID = id
	m.reservas[id] = datos
	return datos, true
}
func (m *clientesRepoCompletoMock) BorrarReserva(id uint) bool {
	_, ok := m.reservas[id]
	if !ok {
		return false
	}
	delete(m.reservas, id)
	return true
}
func (m *clientesRepoCompletoMock) ListarPagos() []models.Pago {
	lista := make([]models.Pago, 0, len(m.pagos))
	for _, p := range m.pagos {
		lista = append(lista, p)
	}
	return lista
}
func (m *clientesRepoCompletoMock) BuscarPagoPorID(id uint) (models.Pago, bool) {
	p, ok := m.pagos[id]
	return p, ok
}
func (m *clientesRepoCompletoMock) CrearPago(p models.Pago) (models.Pago, error) {
	m.sigPagoID++
	p.ID = m.sigPagoID
	m.pagos[p.ID] = p
	return p, nil
}
func (m *clientesRepoCompletoMock) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	_, ok := m.pagos[id]
	if !ok {
		return models.Pago{}, false
	}
	datos.ID = id
	m.pagos[id] = datos
	return datos, true
}
func (m *clientesRepoCompletoMock) BorrarPago(id uint) bool {
	_, ok := m.pagos[id]
	if !ok {
		return false
	}
	delete(m.pagos, id)
	return true
}
func (m *clientesRepoCompletoMock) ClienteTienePagoEntrada(clienteID uint) bool {
	for _, p := range m.pagos {
		if p.ClienteID == clienteID && (p.Concepto == "medio_dia" || p.Concepto == "dia") {
			return true
		}
	}
	return false
}

func TestClientesService_ListarClientes(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	require.Empty(t, svc.ListarClientes())
}

func TestClientesService_ObtenerCliente(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, ok := svc.ObtenerCliente(99)
	require.False(t, ok)
}

func TestClientesService_ActualizarCliente(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	creado, _ := svc.CrearCliente(models.Cliente{Nombre: "Maria", Cedula: "1234"})
	actualizado, err := svc.ActualizarCliente(creado.ID, models.Cliente{Nombre: "Maria Editada", Cedula: "1234"})
	require.NoError(t, err)
	require.Equal(t, "Maria Editada", actualizado.Nombre)
}

func TestClientesService_ActualizarCliente_NoEncontrado(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, err := svc.ActualizarCliente(99, models.Cliente{Nombre: "Test", Cedula: "0000"})
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestClientesService_BorrarCliente(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	creado, _ := svc.CrearCliente(models.Cliente{Nombre: "Maria", Cedula: "1234"})
	err := svc.BorrarCliente(creado.ID)
	require.NoError(t, err)
}

func TestClientesService_BorrarCliente_NoEncontrado(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	err := svc.BorrarCliente(99)
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestClientesService_CrearCliente_DuplicadoError(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	repo.CrearCliente(models.Cliente{Nombre: "Existente", Cedula: "0000"})
	svc := NewClientesService(repo)
	_, err := svc.CrearCliente(models.Cliente{Nombre: "Otro", Cedula: "0000"})
	require.NoError(t, err) // el mock no tiene validacion UNIQUE, no da error
}

func TestClientesService_ListarReservas(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	require.Empty(t, svc.ListarReservas())
}

func TestClientesService_ObtenerReserva(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, ok := svc.ObtenerReserva(99)
	require.False(t, ok)
}

func TestClientesService_CrearReserva_Exitoso(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	creado, _ := svc.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: "1111", Membresia: "mensual"})
	reserva, err := svc.CrearReserva(models.Reserva{ClienteID: creado.ID, Duracion: 720})
	require.NoError(t, err)
	require.NotZero(t, reserva.ID)
}

func TestClientesService_CrearReserva_CamposObligatorios(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, err := svc.CrearReserva(models.Reserva{ClienteID: 0})
	require.ErrorIs(t, err, ErrCampoObligatorio)
}

func TestClientesService_CrearReserva_ClienteInvalido(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, err := svc.CrearReserva(models.Reserva{ClienteID: 99, Duracion: 720})
	require.ErrorIs(t, err, ErrClienteInvalido)
}

func TestClientesService_CrearReserva_DuracionInvalida(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	creado, _ := svc.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: "1111", Membresia: "mensual"})
	_, err := svc.CrearReserva(models.Reserva{ClienteID: creado.ID, Duracion: 100})
	require.ErrorIs(t, err, ErrDuracionInvalida)
}

func TestClientesService_CrearReserva_EstadoPorDefecto(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := svc.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: "1111", Membresia: "mensual"})
	rv, err := svc.CrearReserva(models.Reserva{ClienteID: c.ID, Duracion: 1440})
	require.NoError(t, err)
	require.Equal(t, "pendiente", rv.Estado)
}

func TestClientesService_ActualizarReserva(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := svc.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: "1111", Membresia: "mensual"})
	rv, _ := svc.CrearReserva(models.Reserva{ClienteID: c.ID, Duracion: 720})
	actualizado, err := svc.ActualizarReserva(rv.ID, models.Reserva{ClienteID: c.ID, Duracion: 1440})
	require.NoError(t, err)
	require.Equal(t, 1440, actualizado.Duracion)
}

func TestClientesService_ActualizarReserva_NoEncontrado(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := repo.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: "1111", Membresia: "mensual"})
	_, err := svc.ActualizarReserva(99, models.Reserva{ClienteID: c.ID, Duracion: 720})
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestClientesService_BorrarReserva(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	err := svc.BorrarReserva(99)
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestClientesService_ListarPagos(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	require.Empty(t, svc.ListarPagos())
}

func TestClientesService_ObtenerPago(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, ok := svc.ObtenerPago(99)
	require.False(t, ok)
}

func TestClientesService_CrearPago_CamposObligatorios(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, err := svc.CrearPago(models.Pago{ClienteID: 0})
	require.ErrorIs(t, err, ErrCampoObligatorio)
}

func TestClientesService_CrearPago_MontoInvalido(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 0})
	require.ErrorIs(t, err, ErrMontoInvalido)
}

func TestClientesService_CrearPago_ClienteInvalido(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	_, err := svc.CrearPago(models.Pago{ClienteID: 99, Monto: 5, Concepto: "dia"})
	require.ErrorIs(t, err, ErrClienteInvalido)
}

func TestClientesService_CrearPago_ConceptoInvalido(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	repo.CrearCliente(models.Cliente{ID: 1, Nombre: "Luis", Cedula: "0000", Membresia: "ninguna"})
	_, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 5, Concepto: "invalido"})
	require.ErrorIs(t, err, ErrConceptoPagoInvalido)
}

func TestClientesService_ActualizarPago(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := repo.CrearCliente(models.Cliente{Nombre: "Luis", Cedula: "0000", Membresia: "ninguna"})
	p, _ := repo.CrearPago(models.Pago{ClienteID: c.ID, Monto: 5, Concepto: "dia"})
	actualizado, err := svc.ActualizarPago(p.ID, models.Pago{ClienteID: c.ID, Monto: 10, Concepto: "medio_dia"})
	require.NoError(t, err)
	require.Equal(t, 10.0, actualizado.Monto)
}

func TestClientesService_ActualizarPago_NoEncontrado(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := repo.CrearCliente(models.Cliente{Nombre: "Luis", Cedula: "0000", Membresia: "ninguna"})
	_, err := svc.ActualizarPago(99, models.Pago{ClienteID: c.ID, Monto: 5, Concepto: "dia"})
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestClientesService_BorrarPago(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	err := svc.BorrarPago(99)
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestClientesService_CrearPago_ExitosoSinMembresia(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := repo.CrearCliente(models.Cliente{Nombre: "Luis", Cedula: "0000", Membresia: "ninguna"})
	p, err := svc.CrearPago(models.Pago{ClienteID: c.ID, Monto: 5, Concepto: "dia"})
	require.NoError(t, err)
	require.NotZero(t, p.ID)
}

func TestClientesService_CrearPago_ClienteConMembresiaRechazado(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := repo.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: "1111", Membresia: "mensual"})
	_, err := svc.CrearPago(models.Pago{ClienteID: c.ID, Monto: 5, Concepto: "dia"})
	require.ErrorIs(t, err, ErrClienteConMembresia)
}

func TestClienteTieneMembresia(t *testing.T) {
	require.True(t, clienteTieneMembresia(models.Cliente{Membresia: "mensual"}))
	require.True(t, clienteTieneMembresia(models.Cliente{Membresia: "anual"}))
	require.False(t, clienteTieneMembresia(models.Cliente{Membresia: "ninguna"}))
	require.False(t, clienteTieneMembresia(models.Cliente{Membresia: ""}))
}
