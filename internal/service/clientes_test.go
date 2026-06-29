package service

import (
	"testing"

	"pool-api/internal/models"
)

// ─── MOCKS MANUALES PARA CLIENTES ──────────────────────────────────────────
//
// Mock manual de storage.ClientesModulo que implementa toda la interfaz
// 

type mockClientesModulo struct {
	clientes        map[uint]models.Cliente
	reservas        map[uint]models.Reserva
	pagos           map[uint]models.Pago
	ultimoClienteID uint
	ultimoReservaID uint
	ultimoPagoID    uint

	// flags para verificar qué se llamó
	crearClienteLlamado      bool
	actualizarClienteLlamado bool
	borrarClienteLlamado     bool
	crearReservaLlamado      bool
	crearPagoLlamado         bool
}

func newMockClientesModulo() *mockClientesModulo {
	return &mockClientesModulo{
		clientes: make(map[uint]models.Cliente),
		reservas: make(map[uint]models.Reserva),
		pagos:    make(map[uint]models.Pago),
	}
}

// ─── CLIENTE ──────────────────────────────────────────────────────────────────

func (m *mockClientesModulo) ListarClientes() []models.Cliente {
	result := make([]models.Cliente, 0, len(m.clientes))
	for _, c := range m.clientes {
		result = append(result, c)
	}
	return result
}

func (m *mockClientesModulo) BuscarClientePorID(id uint) (models.Cliente, bool) {
	c, ok := m.clientes[id]
	return c, ok
}

func (m *mockClientesModulo) CrearCliente(c models.Cliente) models.Cliente {
	m.crearClienteLlamado = true
	m.ultimoClienteID++
	c.ID = m.ultimoClienteID
	m.clientes[c.ID] = c
	return c
}

func (m *mockClientesModulo) ActualizarCliente(id uint, c models.Cliente) (models.Cliente, bool) {
	m.actualizarClienteLlamado = true
	if _, ok := m.clientes[id]; !ok {
		return models.Cliente{}, false
	}
	c.ID = id
	m.clientes[id] = c
	return c, true
}

func (m *mockClientesModulo) BorrarCliente(id uint) bool {
	m.borrarClienteLlamado = true
	if _, ok := m.clientes[id]; !ok {
		return false
	}
	delete(m.clientes, id)
	return true
}

// ─── RESERVA ──────────────────────────────────────────────────────────────────

func (m *mockClientesModulo) ListarReservas() []models.Reserva {
	result := make([]models.Reserva, 0, len(m.reservas))
	for _, r := range m.reservas {
		result = append(result, r)
	}
	return result
}

func (m *mockClientesModulo) BuscarReservaPorID(id uint) (models.Reserva, bool) {
	r, ok := m.reservas[id]
	return r, ok
}

func (m *mockClientesModulo) CrearReserva(r models.Reserva) models.Reserva {
	m.crearReservaLlamado = true
	m.ultimoReservaID++
	r.ID = m.ultimoReservaID
	m.reservas[r.ID] = r
	return r
}

func (m *mockClientesModulo) ActualizarReserva(id uint, r models.Reserva) (models.Reserva, bool) {
	if _, ok := m.reservas[id]; !ok {
		return models.Reserva{}, false
	}
	r.ID = id
	m.reservas[id] = r
	return r, true
}

func (m *mockClientesModulo) BorrarReserva(id uint) bool {
	if _, ok := m.reservas[id]; !ok {
		return false
	}
	delete(m.reservas, id)
	return true
}

// ─── PAGO ─────────────────────────────────────────────────────────────────────

func (m *mockClientesModulo) ListarPagos() []models.Pago {
	result := make([]models.Pago, 0, len(m.pagos))
	for _, p := range m.pagos {
		result = append(result, p)
	}
	return result
}

func (m *mockClientesModulo) BuscarPagoPorID(id uint) (models.Pago, bool) {
	p, ok := m.pagos[id]
	return p, ok
}

func (m *mockClientesModulo) CrearPago(p models.Pago) models.Pago {
	m.crearPagoLlamado = true
	m.ultimoPagoID++
	p.ID = m.ultimoPagoID
	m.pagos[p.ID] = p
	return p
}

func (m *mockClientesModulo) ActualizarPago(id uint, p models.Pago) (models.Pago, bool) {
	if _, ok := m.pagos[id]; !ok {
		return models.Pago{}, false
	}
	p.ID = id
	m.pagos[id] = p
	return p, true
}

func (m *mockClientesModulo) BorrarPago(id uint) bool {
	if _, ok := m.pagos[id]; !ok {
		return false
	}
	delete(m.pagos, id)
	return true
}

func (m *mockClientesModulo) ClienteTienePagoEntrada(clienteID uint) bool {
	// Por defecto devolvemos false; los tests pueden sobrescribir si lo necesitan
	for _, p := range m.pagos {
		if p.ClienteID == clienteID && p.Concepto == "entrada" {
			return true
		}
	}
	return false
}

// ═════════════════════════════════════════════════════════════════════════════
// ─── TESTS PARA CLIENTE ────────────────────────────────────────────────────────
// ═════════════════════════════════════════════════════════════════════════════

// TestCrearCliente_Exitoso prueba que se crea un cliente con todos los campos.
func TestCrearCliente_Exitoso(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	cliente := models.Cliente{
		Nombre:    "Juan Pérez",
		Cedula:    "1234567890",
		Membresia: "anual",
	}

	resultado, err := svc.CrearCliente(cliente)
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}

	if !repo.crearClienteLlamado {
		t.Fatal("se esperaba que CrearCliente llegara al repositorio")
	}

	if resultado.ID == 0 {
		t.Error("se esperaba un ID asignado")
	}
	if resultado.Nombre != "Juan Pérez" {
		t.Errorf("nombre incorrecto: %s", resultado.Nombre)
	}
	if resultado.Cedula != "1234567890" {
		t.Errorf("cédula incorrecta: %s", resultado.Cedula)
	}
}

// TestCrearCliente_SinNombre prueba que se rechaza un cliente sin nombre.
func TestCrearCliente_SinNombre(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	cliente := models.Cliente{
		Nombre: "",
		Cedula: "1234567890",
	}

	_, err := svc.CrearCliente(cliente)
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo: %v", err)
	}

	if repo.crearClienteLlamado {
		t.Error("no se debería haber llegado al repositorio sin nombre")
	}
}

// TestCrearCliente_SinCedula prueba que se rechaza un cliente sin cédula.
func TestCrearCliente_SinCedula(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	cliente := models.Cliente{
		Nombre: "Juan Pérez",
		Cedula: "",
	}

	_, err := svc.CrearCliente(cliente)
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo: %v", err)
	}

	if repo.crearClienteLlamado {
		t.Error("no se debería haber llegado al repositorio sin cédula")
	}
}

// TestCrearCliente_MembresiaPorDefecto prueba que si no se especifica membresía,
// se asigna "ninguna" automáticamente.
func TestCrearCliente_MembresiaPorDefecto(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	cliente := models.Cliente{
		Nombre: "María García",
		Cedula: "9876543210",
		// Membresia vacía
	}

	resultado, err := svc.CrearCliente(cliente)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.Membresia != "ninguna" {
		t.Errorf("se esperaba membresía 'ninguna', se obtuvo: %s", resultado.Membresia)
	}
}

// TestActualizarCliente_Exitoso prueba que se actualiza un cliente existente.
func TestActualizarCliente_Exitoso(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	// Crear cliente primero
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "123"}

	actualizado := models.Cliente{
		Nombre: "Juan Actualizado",
		Cedula: "123",
	}

	resultado, err := svc.ActualizarCliente(1, actualizado)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.Nombre != "Juan Actualizado" {
		t.Errorf("nombre no actualizado: %s", resultado.Nombre)
	}
}

// TestActualizarCliente_NoEncontrado prueba que se rechaza la actualización
// de un cliente que no existe.
func TestActualizarCliente_NoEncontrado(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	cliente := models.Cliente{
		Nombre: "Juan",
		Cedula: "123",
	}

	_, err := svc.ActualizarCliente(999, cliente)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo: %v", err)
	}
}

// TestBorrarCliente_Exitoso prueba que se borra un cliente.
func TestBorrarCliente_Exitoso(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "123"}

	err := svc.BorrarCliente(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if _, ok := repo.clientes[1]; ok {
		t.Error("el cliente debería haber sido borrado")
	}
}

// TestBorrarCliente_NoEncontrado prueba que se rechaza el borrado de un
// cliente que no existe.
func TestBorrarCliente_NoEncontrado(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	err := svc.BorrarCliente(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo: %v", err)
	}
}

// ═════════════════════════════════════════════════════════════════════════════
// ─── TESTS PARA RESERVA ────────────────────────────────────────────────────────
// ═════════════════════════════════════════════════════════════════════════════

// TestCrearReserva_Exitoso prueba que se crea una reserva con cliente válido.
func TestCrearReserva_Exitoso(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	// Crear cliente primero
	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "123"}

	reserva := models.Reserva{
		ClienteID: 1,
		Duracion:  120,
	}

	resultado, err := svc.CrearReserva(reserva)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if !repo.crearReservaLlamado {
		t.Fatal("se esperaba que CrearReserva llegara al repositorio")
	}

	if resultado.ID == 0 {
		t.Error("se esperaba un ID asignado")
	}
}

// TestCrearReserva_SinClienteID prueba que se rechaza una reserva sin ClienteID.
func TestCrearReserva_SinClienteID(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	reserva := models.Reserva{
		ClienteID: 0,
	}

	_, err := svc.CrearReserva(reserva)
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo: %v", err)
	}
}

// TestCrearReserva_ClienteInvalido prueba que se rechaza una reserva con
// ClienteID que no existe en la base de datos.
func TestCrearReserva_ClienteInvalido(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	reserva := models.Reserva{
		ClienteID: 999,
	}

	_, err := svc.CrearReserva(reserva)
	if err != ErrClienteInvalido {
		t.Fatalf("se esperaba ErrClienteInvalido, se obtuvo: %v", err)
	}

	if repo.crearReservaLlamado {
		t.Error("no se debería llegar al repositorio con cliente inválido")
	}
}

// TestCrearReserva_EstadoPorDefecto prueba que si no se especifica estado,
// se asigna "pendiente" automáticamente.
func TestCrearReserva_EstadoPorDefecto(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "123"}

	reserva := models.Reserva{
		ClienteID: 1,
		Estado:    "",
	}

	resultado, err := svc.CrearReserva(reserva)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.Estado != "pendiente" {
		t.Errorf("se esperaba estado 'pendiente', se obtuvo: %s", resultado.Estado)
	}
}

// ═════════════════════════════════════════════════════════════════════════════
// ─── TESTS PARA PAGO ───────────────────────────────────────────────────────────
// ═════════════════════════════════════════════════════════════════════════════

// TestCrearPago_Exitoso prueba que se crea un pago con datos válidos.
func TestCrearPago_Exitoso(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	repo.clientes[1] = models.Cliente{ID: 1, Nombre: "Juan", Cedula: "123"}

	pago := models.Pago{
		ClienteID: 1,
		Monto:     100.0,
		Concepto:  "entrada",
	}

	resultado, err := svc.CrearPago(pago)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if !repo.crearPagoLlamado {
		t.Fatal("se esperaba que CrearPago llegara al repositorio")
	}

	if resultado.ID == 0 {
		t.Error("se esperaba un ID asignado")
	}
}

// TestCrearPago_SinClienteID prueba que se rechaza un pago sin ClienteID.
func TestCrearPago_SinClienteID(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	pago := models.Pago{
		ClienteID: 0,
		Monto:     100.0,
	}

	_, err := svc.CrearPago(pago)
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo: %v", err)
	}
}

// TestCrearPago_MontoInvalido prueba que se rechaza un pago con monto <= 0.
func TestCrearPago_MontoInvalido(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	pago := models.Pago{
		ClienteID: 1,
		Monto:     -50.0,
	}

	_, err := svc.CrearPago(pago)
	if err != ErrMontoInvalido {
		t.Fatalf("se esperaba ErrMontoInvalido, se obtuvo: %v", err)
	}
}

// TestCrearPago_MontoZero prueba que se rechaza un pago con monto 0.
func TestCrearPago_MontoZero(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	pago := models.Pago{
		ClienteID: 1,
		Monto:     0,
	}

	_, err := svc.CrearPago(pago)
	if err != ErrMontoInvalido {
		t.Fatalf("se esperaba ErrMontoInvalido, se obtuvo: %v", err)
	}
}

// TestCrearPago_ClienteInvalido prueba que se rechaza un pago con ClienteID
// que no existe en la base de datos.
func TestCrearPago_ClienteInvalido(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	pago := models.Pago{
		ClienteID: 999,
		Monto:     100.0,
	}

	_, err := svc.CrearPago(pago)
	if err != ErrClienteInvalido {
		t.Fatalf("se esperaba ErrClienteInvalido, se obtuvo: %v", err)
	}

	if repo.crearPagoLlamado {
		t.Error("no se debería llegar al repositorio con cliente inválido")
	}
}

// TestBorrarPago_NoEncontrado prueba que se rechaza el borrado de un pago
// que no existe.
func TestBorrarPago_NoEncontrado(t *testing.T) {
	repo := newMockClientesModulo()
	svc := NewClientesService(repo)

	err := svc.BorrarPago(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo: %v", err)
	}
}
