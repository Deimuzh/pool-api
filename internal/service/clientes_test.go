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
	c.ID = uint(len(m.clientes)) + 1
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
	r.ID = uint(len(m.reservas)) + 1
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

func (m *clientesModuloMock) BorrarReserva(id uint) bool {
	if _, ok := m.reservas[id]; !ok {
		return false
	}
	delete(m.reservas, id)
	return true
}

func (m *clientesModuloMock) ListarPagos() []models.Pago { return nil }
func (m *clientesModuloMock) BuscarPagoPorID(id uint) (models.Pago, bool) {
	p, ok := m.pagos[id]
	return p, ok
}

func (m *clientesModuloMock) CrearPago(p models.Pago) (models.Pago, error) {
	m.crearPagoLlamado = true
	p.ID = uint(len(m.pagos)) + 1
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

func (m *clientesModuloMock) BorrarPago(id uint) bool {
	if _, ok := m.pagos[id]; !ok {
		return false
	}
	delete(m.pagos, id)
	return true
}

func (m *clientesModuloMock) ClienteTienePagoEntrada(clienteID uint) bool {
	for _, p := range m.pagos {
		if p.ClienteID == clienteID && (p.Concepto == "medio_dia" || p.Concepto == "dia") {
			return true
		}
	}
	return false
}

func TestClientesService_ListarClientes_Vacio(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	lista := svc.ListarClientes()
	if len(lista) != 0 {
		t.Fatal("se esperaba lista vacia")
	}
}

func TestClientesService_validarEmail(t *testing.T) {
	if !validarEmail("test@example.com") {
		t.Fatal("se esperaba email valido")
	}
	if validarEmail("invalido") {
		t.Fatal("se esperaba email invalido")
	}
	if validarEmail("") {
		t.Fatal("se esperaba email vacio como invalido")
	}
}

func TestClientesService_CrearCliente_TrimEspacios(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	creado, err := svc.CrearCliente(models.Cliente{Nombre: "  Maria Garcia  ", Cedula: " 1312345678 "})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.Nombre != "Maria Garcia" {
		t.Fatalf("se esperaba nombre sin espacios, se obtuvo %q", creado.Nombre)
	}
}

func TestClientesService_CrearCliente_CedulaFormatoInvalido(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	_, err := svc.CrearCliente(models.Cliente{Nombre: "Test", Cedula: "123"})
	if err != ErrCedulaFormatoInvalido {
		t.Fatalf("se esperaba ErrCedulaFormatoInvalido, se obtuvo %v", err)
	}
}

func TestClientesService_CrearCliente_EmailFormatoInvalido(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	_, err := svc.CrearCliente(models.Cliente{Nombre: "Test", Cedula: "1312345678", Email: "invalido"})
	if err != ErrEmailFormatoInvalido {
		t.Fatalf("se esperaba ErrEmailFormatoInvalido, se obtuvo %v", err)
	}
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

func TestClientesService_ObtenerCliente_NoEncontrado(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	_, ok := svc.ObtenerCliente(999)
	if ok {
		t.Fatal("se esperaba ok=false para cliente inexistente")
	}
}

func TestClientesService_ActualizarCliente_Exitoso(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Cedula: "1312345678"}
	svc := NewClientesService(repo)

	actualizado, err := svc.ActualizarCliente(1, models.Cliente{Nombre: "Luis Actualizado", Cedula: "1312345678"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Nombre != "Luis Actualizado" {
		t.Fatalf("se esperaba nombre actualizado, se obtuvo %q", actualizado.Nombre)
	}
}

func TestClientesService_ActualizarCliente_CedulaDuplicada(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Cedula: "1312345678"}
	repo.clientes[2] = models.Cliente{ID: 2, Nombre: "Ana Reyes", Cedula: "1398765432"}
	svc := NewClientesService(repo)

	_, err := svc.ActualizarCliente(1, models.Cliente{Nombre: "Luis Actualizado", Cedula: "1398765432"})
	if err != ErrCedulaEnUso {
		t.Fatalf("se esperaba ErrCedulaEnUso al usar cedula de otro cliente, se obtuvo %v", err)
	}
}

func TestClientesService_ActualizarCliente_NoEncontrado(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	_, err := svc.ActualizarCliente(999, models.Cliente{Nombre: "Inexistente", Cedula: "1312345678"})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestClientesService_BorrarCliente_Exitoso(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Cedula: "1312345678"}
	svc := NewClientesService(repo)

	err := svc.BorrarCliente(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, ok := repo.clientes[1]; ok {
		t.Fatal("el cliente deberia haberse eliminado del repositorio")
	}
}

func TestClientesService_BorrarCliente_NoEncontrado(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	err := svc.BorrarCliente(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestClientesService_CrearReserva_Exitoso(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Ana Reyes", Membresia: "mensual"}
	svc := NewClientesService(repo)

	creado, err := svc.CrearReserva(models.Reserva{ClienteID: 1, Duracion: 720})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if !repo.crearReservaLlamado {
		t.Fatal("se esperaba que CrearReserva llegara al repositorio")
	}
	if creado.Estado != "pendiente" {
		t.Fatalf("se esperaba estado 'pendiente', se obtuvo %q", creado.Estado)
	}
}

func TestClientesService_CrearReserva_DuracionInvalida(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Ana Reyes", Membresia: "mensual"}
	svc := NewClientesService(repo)

	_, err := svc.CrearReserva(models.Reserva{ClienteID: 1, Duracion: 60})
	if err != ErrDuracionInvalida {
		t.Fatalf("se esperaba ErrDuracionInvalida, se obtuvo %v", err)
	}
}

func TestClientesService_CrearReserva_ClienteInvalido(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	_, err := svc.CrearReserva(models.Reserva{ClienteID: 999, Duracion: 720})
	if err != ErrClienteInvalido {
		t.Fatalf("se esperaba ErrClienteInvalido, se obtuvo %v", err)
	}
}

func TestClientesService_CrearReserva_ClienteIDCero(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	_, err := svc.CrearReserva(models.Reserva{ClienteID: 0, Duracion: 720})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestClientesService_ActualizarReserva_Exitoso(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Ana Reyes", Membresia: "mensual"}
	repo.reservas[1] = models.Reserva{ID: 1, ClienteID: 1, Duracion: 720, Estado: "pendiente"}
	svc := NewClientesService(repo)

	actualizado, err := svc.ActualizarReserva(1, models.Reserva{ClienteID: 1, Duracion: 1440, Estado: "confirmada"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Duracion != 1440 {
		t.Fatalf("se esperaba duracion 1440, se obtuvo %d", actualizado.Duracion)
	}
}

func TestClientesService_ActualizarReserva_NoEncontrada(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Ana Reyes", Membresia: "mensual"}
	svc := NewClientesService(repo)

	_, err := svc.ActualizarReserva(999, models.Reserva{ClienteID: 1, Duracion: 720})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestClientesService_BorrarReserva_Exitoso(t *testing.T) {
	repo := newClientesModuloMock()
	repo.reservas[1] = models.Reserva{ID: 1, ClienteID: 1, Duracion: 720}
	svc := NewClientesService(repo)

	err := svc.BorrarReserva(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, ok := repo.reservas[1]; ok {
		t.Fatal("la reserva deberia haberse eliminado del repositorio")
	}
}

func TestClientesService_BorrarReserva_NoEncontrada(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	err := svc.BorrarReserva(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestClientesService_CrearPago_Exitoso(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Membresia: "ninguna"}
	svc := NewClientesService(repo)

	creado, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 5, Concepto: "dia", Metodo: "efectivo"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if !repo.crearPagoLlamado {
		t.Fatal("se esperaba que CrearPago llegara al repositorio")
	}
	if creado.Concepto != "dia" {
		t.Fatalf("se esperaba concepto 'dia', se obtuvo %q", creado.Concepto)
	}
}

func TestClientesService_CrearPago_ClienteIDCero(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	_, err := svc.CrearPago(models.Pago{ClienteID: 0, Monto: 5, Concepto: "dia"})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestClientesService_CrearPago_ConceptoInvalido(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Membresia: "ninguna"}
	svc := NewClientesService(repo)

	_, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 5, Concepto: "semanal"})
	if err != ErrConceptoPagoInvalido {
		t.Fatalf("se esperaba ErrConceptoPagoInvalido, se obtuvo %v", err)
	}
}

func TestClientesService_CrearPago_MontoInvalido(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Membresia: "ninguna"}
	svc := NewClientesService(repo)

	_, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 0, Concepto: "dia"})
	if err != ErrMontoInvalido {
		t.Fatalf("se esperaba ErrMontoInvalido, se obtuvo %v", err)
	}
}

func TestClientesService_CrearPago_MetodoInvalido(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Membresia: "ninguna"}
	svc := NewClientesService(repo)

	_, err := svc.CrearPago(models.Pago{ClienteID: 1, Monto: 5, Concepto: "dia", Metodo: "tarjeta"})
	if err != ErrMetodoPagoInvalido {
		t.Fatalf("se esperaba ErrMetodoPagoInvalido, se obtuvo %v", err)
	}
}

func TestClientesService_ActualizarPago_Exitoso(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Membresia: "ninguna"}
	repo.pagos[1] = models.Pago{ID: 1, ClienteID: 1, Monto: 5, Concepto: "dia", Metodo: "efectivo"}
	svc := NewClientesService(repo)

	actualizado, err := svc.ActualizarPago(1, models.Pago{ClienteID: 1, Monto: 10, Concepto: "dia", Metodo: "transferencia"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Monto != 10 {
		t.Fatalf("se esperaba monto 10, se obtuvo %f", actualizado.Monto)
	}
	if actualizado.Metodo != "transferencia" {
		t.Fatalf("se esperaba metodo 'transferencia', se obtuvo %q", actualizado.Metodo)
	}
}

func TestClientesService_ActualizarPago_NoEncontrado(t *testing.T) {
	repo := newClientesModuloMock()
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Luis Pino", Membresia: "ninguna"}
	svc := NewClientesService(repo)

	_, err := svc.ActualizarPago(999, models.Pago{ClienteID: 1, Monto: 5, Concepto: "dia", Metodo: "efectivo"})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestClientesService_BorrarPago_Exitoso(t *testing.T) {
	repo := newClientesModuloMock()
	repo.pagos[1] = models.Pago{ID: 1, ClienteID: 1, Monto: 5, Concepto: "dia"}
	svc := NewClientesService(repo)

	err := svc.BorrarPago(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, ok := repo.pagos[1]; ok {
		t.Fatal("el pago deberia haberse eliminado del repositorio")
	}
}

func TestClientesService_BorrarPago_NoEncontrado(t *testing.T) {
	repo := newClientesModuloMock()
	svc := NewClientesService(repo)

	err := svc.BorrarPago(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}
