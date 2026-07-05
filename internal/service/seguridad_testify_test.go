package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type seguridadRepoMock struct {
	mock.Mock
}

var _ storage.SeguridadRepository = (*seguridadRepoMock)(nil)

func (m *seguridadRepoMock) ListarGuardavidas() []models.Guardavida {
	args := m.Called()
	return args.Get(0).([]models.Guardavida)
}

func (m *seguridadRepoMock) BuscarGuardavidaPorID(id uint) (models.Guardavida, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Guardavida), args.Bool(1)
}

func (m *seguridadRepoMock) CrearGuardavida(g models.Guardavida) models.Guardavida {
	args := m.Called(g)
	return args.Get(0).(models.Guardavida)
}

func (m *seguridadRepoMock) ActualizarGuardavida(id uint, datos models.Guardavida) (models.Guardavida, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Guardavida), args.Bool(1)
}

func (m *seguridadRepoMock) BorrarGuardavida(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *seguridadRepoMock) ListarIncidentes() []models.Incidente {
	args := m.Called()
	return args.Get(0).([]models.Incidente)
}

func (m *seguridadRepoMock) BuscarIncidentePorID(id uint) (models.Incidente, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Incidente), args.Bool(1)
}

func (m *seguridadRepoMock) CrearIncidente(i models.Incidente) models.Incidente {
	args := m.Called(i)
	return args.Get(0).(models.Incidente)
}

func (m *seguridadRepoMock) ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Incidente), args.Bool(1)
}

func (m *seguridadRepoMock) BorrarIncidente(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *seguridadRepoMock) ListarAccesos() []models.AccesoCliente {
	args := m.Called()
	return args.Get(0).([]models.AccesoCliente)
}

func (m *seguridadRepoMock) BuscarAccesoPorID(id uint) (models.AccesoCliente, bool) {
	args := m.Called(id)
	return args.Get(0).(models.AccesoCliente), args.Bool(1)
}

func (m *seguridadRepoMock) CrearAcceso(a models.AccesoCliente) models.AccesoCliente {
	args := m.Called(a)
	return args.Get(0).(models.AccesoCliente)
}

func (m *seguridadRepoMock) ActualizarAcceso(id uint, datos models.AccesoCliente) (models.AccesoCliente, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.AccesoCliente), args.Bool(1)
}

func (m *seguridadRepoMock) BorrarAcceso(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

type clienteRepoMock struct {
	mock.Mock
}

var _ storage.ClienteRepository = (*clienteRepoMock)(nil)

func (m *clienteRepoMock) ListarClientes() []models.Cliente {
	args := m.Called()
	return args.Get(0).([]models.Cliente)
}

func (m *clienteRepoMock) BuscarClientePorID(id uint) (models.Cliente, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Cliente), args.Bool(1)
}

func (m *clienteRepoMock) CrearCliente(c models.Cliente) (models.Cliente, error) {
	args := m.Called(c)
	return args.Get(0).(models.Cliente), args.Error(1)
}

func (m *clienteRepoMock) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Cliente), args.Bool(1)
}

func (m *clienteRepoMock) BorrarCliente(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

type pagoRepoMock struct {
	mock.Mock
}

var _ storage.PagoRepository = (*pagoRepoMock)(nil)

func (m *pagoRepoMock) ListarPagos() []models.Pago {
	args := m.Called()
	return args.Get(0).([]models.Pago)
}

func (m *pagoRepoMock) BuscarPagoPorID(id uint) (models.Pago, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Pago), args.Bool(1)
}

func (m *pagoRepoMock) CrearPago(p models.Pago) (models.Pago, error) {
	args := m.Called(p)
	return args.Get(0).(models.Pago), args.Error(1)
}

func (m *pagoRepoMock) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Pago), args.Bool(1)
}

func (m *pagoRepoMock) BorrarPago(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *pagoRepoMock) ClienteTienePagoEntrada(clienteID uint) bool {
	args := m.Called(clienteID)
	return args.Bool(0)
}

func TestSeguridadService_CrearAcceso_ConTestifyMock(t *testing.T) {
	seguridad := new(seguridadRepoMock)
	clientes := new(clienteRepoMock)
	pagos := new(pagoRepoMock)

	cliente := models.Cliente{ID: 7, Nombre: "Carlos Mero", Membresia: "ninguna"}
	accesoEsperado := models.AccesoCliente{ID: 1, ClienteID: 7, Autorizado: true}

	clientes.On("BuscarClientePorID", uint(7)).Return(cliente, true)
	pagos.On("ClienteTienePagoEntrada", uint(7)).Return(true)
	seguridad.On("CrearAcceso", models.AccesoCliente{ClienteID: 7, Autorizado: true}).Return(accesoEsperado)

	svc := NewSeguridadService(seguridad, clientes, pagos)
	creado, err := svc.CrearAcceso(7)

	require.NoError(t, err)
	require.Equal(t, uint(1), creado.ID)
	require.Equal(t, "Carlos Mero", creado.NombreCliente)
	require.True(t, creado.PagoAlDia)
	seguridad.AssertExpectations(t)
	clientes.AssertExpectations(t)
	pagos.AssertExpectations(t)
}
