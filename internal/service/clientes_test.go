package service

import (
	"errors"
	"testing"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type clientesModuloMock struct {
	clientes map[uint]models.Cliente
	reservas map[uint]models.Reserva
	pagos    map[uint]models.Pago

	crearClienteLlamado bool
	crearReservaLlamado bool
	crearPagoLlamado    bool
	errorCrearCliente   error
}

var _ storage.ClientesModulo = (*clientesModuloMock)(nil)

func newClientesModuloMock() *clientesModuloMock {
	return &clientesModuloMock{
		clientes: make(map[uint]models.Cliente),
		reservas: make(map[uint]models.Reserva),
		pagos:    make(map[uint]models.Pago),
	}
}

func (m *clientesModuloMock) ListarClientes() []models.Cliente {
	lista := make([]models.Cliente, 0, len(m.clientes))
	for _, c := range m.clientes {
		lista = append(lista, c)
	}
	return lista
}

func (m *clientesModuloMock) BuscarClientePorID(id uint) (models.Cliente, bool) {
	c, ok := m.clientes[id]
	return c, ok
}

func (m *clientesModuloMock) CrearCliente(c models.Cliente) (models.Cliente, error) {
	m.crearClienteLlamado = true
	if m.errorCrearCliente != nil {
		return models.Cliente{}, m.errorCrearCliente
	}
	c.ID = uint(len(m.clientes) + 1)
	m.clientes[c.ID] = c
	return c, nil
}

func (m *clientesModuloMock) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	if _, ok := m.clientes[id]; !ok {
		return models.Cliente{}, false
	}
	datos.ID = id
	m.clientes[id] = datos
	return datos, true
}

func (m *clientesModuloMock) BorrarCliente(id uint) bool {
	if _, ok := m.clientes[id]; !ok {
		return false
	}
	delete(m.clientes, id)
	return true
}

func (m *clientesModuloMock) ListarReservas() []models.Reserva { return nil }
func (m *clientesModuloMock) BuscarReservaPorID(id uint) (models.Reserva, bool) {
	r, ok := m.reservas[id]
	return r, ok
}

func (m *clientesModuloMock) CrearReserva(r models.Reserva) (models.Reserva, error) {
	m.crearReservaLlamado = true
	r.ID = uint(len(m.reservas) + 1)
	m.reservas[r.ID] = r
	return r, nil
}

func (m *clientesModuloMock) ActualizarReserva(id uint, datos models.Reserva) (models.Reserva, bool) {
	if _, ok := m.reservas[id]; !ok {
		return models.Reserva{}, false
	}
	datos.ID = id
	m.reservas[id] = datos
	return datos, true
}

func (m *clientesModuloMock) BorrarReserva(id uint) bool { return false }

func (m *clientesModuloMock) ListarPagos() []models.Pago { return nil }
func (m *clientesModuloMock) BuscarPagoPorID(id uint) (models.Pago, bool) {
	p, ok := m.pagos[id]
	return p, ok
}

func (m *clientesModuloMock) CrearPago(p models.Pago) (models.Pago, error) {
	m.crearPagoLlamado = true
	p.ID = uint(len(m.pagos) + 1)
	m.pagos[p.ID] = p
	return p, nil
}

func (m *clientesModuloMock) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	if _, ok := m.pagos[id]; !ok {
		return models.Pago{}, false
	}
	datos.ID = id
	m.pagos[id] = datos
	return datos, true
}

func (m *clientesModuloMock) BorrarPago(id uint) bool { return false }

func (m *clientesModuloMock) ClienteTienePagoEntrada(clienteID uint) bool {
	for _, p := range m.pagos {
		if p.ClienteID == clienteID && (p.Concepto == "medio_dia" || p.Concepto == "dia") {
			return true
		}
	}
	return false
}

func TestClientesService_CrearCliente_AsignaMembresiaPorDefecto(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	creado, err := svc.CrearCliente(models.Cliente{Nombre: "Maria Garcia", Cedula: "1312345678"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if !repo.crearClienteLlamado {
		t.Fatal("se esperaba que CrearCliente llegara al repositorio")
	}
	if creado.Membresia != "ninguna" {
		t.Fatalf("se esperaba membresia por defecto 'ninguna', se obtuvo %q", creado.Membresia)
	}
}

func TestClientesService_CrearCliente_InvalidoNoLlegaAlRepo(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	_, err := svc.CrearCliente(models.Cliente{Nombre: "", Cedula: "1312345678"})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
	if repo.crearClienteLlamado {
		t.Fatal("no debe tocar el repositorio si faltan campos obligatorios")
	}
}

func TestClientesService_CrearCliente_CedulaDuplicada(t *testing.T) {
	repo := newClientesModuloMock()
	repo.errorCrearCliente = errors.New("UNIQUE constraint failed: clientes.cedula")
	svc := NewClientesService(repo)

	_, err := svc.CrearCliente(models.Cliente{Nombre: "Luis Pino", Cedula: "1312345678"})
	if err != ErrCedulaEnUso {
		t.Fatalf("se esperaba ErrCedulaEnUso, se obtuvo %v", err)
	}
}

func TestClientesService_CrearReserva_RequiereMembresia(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Membresia: "ninguna"}
	svc := NewClientesService(repo)

	_, err := svc.CrearReserva(models.Reserva{ClienteID: 1, Duracion: 720})
	if err != ErrClienteSinMembresia {
		t.Fatalf("se esperaba ErrClienteSinMembresia, se obtuvo %v", err)
	}
	if repo.crearReservaLlamado {
		t.Fatal("no debe crear reserva para cliente sin membresia")
	}
}

func TestClientesService_CrearPago_RechazaClienteConMembresia(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Ana Reyes", Membresia: "mensual"}
	svc := NewClientesService(repo)

	_, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 3, Concepto: "dia"})
	if err != ErrClienteConMembresia {
		t.Fatalf("se esperaba ErrClienteConMembresia, se obtuvo %v", err)
	}
	if repo.crearPagoLlamado {
		t.Fatal("no debe crear pago de entrada para cliente con membresia")
	}
}
