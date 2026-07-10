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

// Tests unicos

func TestClientesService_CrearPago_ExitosoSinMembresia(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := repo.CrearCliente(models.Cliente{Nombre: "Luis", Cedula: "0000", Membresia: "ninguna"})
	p, err := svc.CrearPago(models.Pago{ClienteID: c.ID, Monto: 5, Concepto: "dia", Metodo: "efectivo"})
	require.NoError(t, err)
	require.NotZero(t, p.ID)
}

func TestClientesService_CrearPago_ClienteConMembresiaRechazado(t *testing.T) {
	repo := newClientesRepoCompletoMock()
	svc := NewClientesService(repo)
	c, _ := repo.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: "1111", Membresia: "mensual"})
	_, err := svc.CrearPago(models.Pago{ClienteID: c.ID, Monto: 5, Concepto: "dia", Metodo: "efectivo"})
	require.ErrorIs(t, err, ErrClienteConMembresia)
}

func TestClienteTieneMembresia(t *testing.T) {
	require.True(t, clienteTieneMembresia(models.Cliente{Membresia: "mensual"}))
	require.True(t, clienteTieneMembresia(models.Cliente{Membresia: "anual"}))
	require.False(t, clienteTieneMembresia(models.Cliente{Membresia: "ninguna"}))
	require.False(t, clienteTieneMembresia(models.Cliente{Membresia: ""}))
}
