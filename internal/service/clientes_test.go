package service

import (
	"testing"

	"pool-api/internal/models"
)

type mockClientesRepo struct {
	crearClienteLlamado bool
	clienteRecibido     models.Cliente
}

func (m *mockClientesRepo) ListarClientes() []models.Cliente { return nil }
func (m *mockClientesRepo) BuscarClientePorID(id int) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (m *mockClientesRepo) CrearCliente(c models.Cliente) models.Cliente {
	m.crearClienteLlamado = true
	m.clienteRecibido = c
	c.ID = 1
	return c
}
func (m *mockClientesRepo) ActualizarCliente(id int, datos models.Cliente) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (m *mockClientesRepo) BorrarCliente(id int) bool { return false }

func (m *mockClientesRepo) ListarReservas() []models.Reserva { return nil }
func (m *mockClientesRepo) BuscarReservaPorID(id int) (models.Reserva, bool) {
	return models.Reserva{}, false
}
func (m *mockClientesRepo) CrearReserva(rv models.Reserva) models.Reserva { return rv }
func (m *mockClientesRepo) ActualizarReserva(id int, datos models.Reserva) (models.Reserva, bool) {
	return models.Reserva{}, false
}
func (m *mockClientesRepo) BorrarReserva(id int) bool { return false }

func (m *mockClientesRepo) ListarPagos() []models.Pago { return nil }
func (m *mockClientesRepo) BuscarPagoPorID(id int) (models.Pago, bool) {
	return models.Pago{}, false
}
func (m *mockClientesRepo) CrearPago(p models.Pago) models.Pago { return p }
func (m *mockClientesRepo) ActualizarPago(id int, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}
func (m *mockClientesRepo) BorrarPago(id int) bool                     { return false }
func (m *mockClientesRepo) ClienteTienePagoEntrada(clienteID int) bool { return false }

func TestCrearCliente_CedulaVaciaNoLlegaAlRepo(t *testing.T) {
	repo := &mockClientesRepo{}
	svc := NewClientesService(repo)

	_, err := svc.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: ""})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
	if repo.crearClienteLlamado {
		t.Fatal("no debería llegar al repositorio cuando la cédula está vacía")
	}
}
