package service

import (
	"testing"

	"pool-api/internal/models"
)

// ─── MOCKS MANUALES ─────────────────────────────────────────────────────────
//
// No usamos ninguna librería de mocking: implementamos a mano las interfaces
// que pide SeguridadService (storage.SeguridadRepository, storage.ClienteRepository,
// storage.PagoRepository). Cada mock solo implementa lo mínimo necesario para
// este test y deja constancia (con flags/contadores) de qué se llamó, para
// poder afirmar "no llegó al repositorio".

// mockSeguridadRepo es un mock de storage.SeguridadRepository.
// Solo nos interesa CrearAcceso para esta prueba: registra si fue invocado.
type mockSeguridadRepo struct {
	crearAccesoLlamado bool
	accesoRecibido     models.AccesoCliente
}

func (m *mockSeguridadRepo) ListarGuardavidas() []models.Guardavida { return nil }
func (m *mockSeguridadRepo) BuscarGuardavidaPorID(id uint) (models.Guardavida, bool) {
	return models.Guardavida{}, false
}
func (m *mockSeguridadRepo) CrearGuardavida(g models.Guardavida) models.Guardavida { return g }
func (m *mockSeguridadRepo) ActualizarGuardavida(id uint, datos models.Guardavida) (models.Guardavida, bool) {
	return models.Guardavida{}, false
}
func (m *mockSeguridadRepo) BorrarGuardavida(id uint) bool { return false }

func (m *mockSeguridadRepo) ListarIncidentes() []models.Incidente { return nil }
func (m *mockSeguridadRepo) BuscarIncidentePorID(id uint) (models.Incidente, bool) {
	return models.Incidente{}, false
}
func (m *mockSeguridadRepo) CrearIncidente(i models.Incidente) models.Incidente { return i }
func (m *mockSeguridadRepo) ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool) {
	return models.Incidente{}, false
}
func (m *mockSeguridadRepo) BorrarIncidente(id uint) bool { return false }

func (m *mockSeguridadRepo) ListarAccesos() []models.AccesoCliente { return nil }
func (m *mockSeguridadRepo) BuscarAccesoPorID(id uint) (models.AccesoCliente, bool) {
	return models.AccesoCliente{}, false
}

// CrearAcceso es el único método que de verdad importa: guarda lo que el
// service le mandó, para que el test pueda revisar el campo Autorizado.
func (m *mockSeguridadRepo) CrearAcceso(a models.AccesoCliente) models.AccesoCliente {
	m.crearAccesoLlamado = true
	m.accesoRecibido = a
	a.ID = 1
	return a
}
func (m *mockSeguridadRepo) ActualizarAcceso(id uint, datos models.AccesoCliente) (models.AccesoCliente, bool) {
	return models.AccesoCliente{}, false
}
func (m *mockSeguridadRepo) BorrarAcceso(id uint) bool { return false }

// mockClienteRepo es un mock de storage.ClienteRepository: solo necesitamos
// BuscarClientePorID para que CrearAcceso pueda validar que el cliente existe.
type mockClienteRepo struct {
	clientes map[int]models.Cliente
}

func (m *mockClienteRepo) ListarClientes() []models.Cliente { return nil }
func (m *mockClienteRepo) BuscarClientePorID(id uint) (models.Cliente, bool) {
	c, ok := m.clientes[int(id)]
	return c, ok
}
func (m *mockClienteRepo) CrearCliente(c models.Cliente) models.Cliente { return c }
func (m *mockClienteRepo) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	return models.Cliente{}, false
}
func (m *mockClienteRepo) BorrarCliente(id uint) bool { return false }

// mockPagoRepo es un mock de storage.PagoRepository: lo que de verdad
// controla este test es ClienteTienePagoEntrada, que devolvemos fijo en false
// para simular "el cliente NO pagó la entrada".
type mockPagoRepo struct {
	tienePago bool
}

func (m *mockPagoRepo) ListarPagos() []models.Pago                  { return nil }
func (m *mockPagoRepo) BuscarPagoPorID(id uint) (models.Pago, bool) { return models.Pago{}, false }
func (m *mockPagoRepo) CrearPago(p models.Pago) models.Pago         { return p }
func (m *mockPagoRepo) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}
func (m *mockPagoRepo) BorrarPago(id uint) bool { return false }
func (m *mockPagoRepo) ClienteTienePagoEntrada(clienteID uint) bool {
	return m.tienePago
}

// ─── TEST ───────────────────────────────────────────────────────────────────

// TestCrearAcceso_SinPagoNoAutoriza prueba la regla de negocio central del
// módulo Seguridad: si el cliente existe pero NO tiene un pago de entrada
// registrado, el acceso debe guardarse como NO autorizado (Autorizado=false)
// y con un motivo explicativo. Si esta regla se rompiera (por ejemplo, si
// alguien cambiara el service para autorizar siempre), el test fallaría
// porque accesoRecibido.Autorizado vendría en true.
func TestCrearAcceso_SinPagoNoAutoriza(t *testing.T) {
	repo := &mockSeguridadRepo{}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			2: {ID: 2, Nombre: "Luis Pino"},
		},
	}
	pagos := &mockPagoRepo{tienePago: false} // el cliente 2 NO ha pagado entrada

	svc := NewSeguridadService(repo, clientes, pagos)

	resultado, err := svc.CrearAcceso(2)
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}

	if !repo.crearAccesoLlamado {
		t.Fatal("se esperaba que CrearAcceso llegara al repositorio")
	}

	// La regla de negocio real: sin pago, NO se autoriza.
	if repo.accesoRecibido.Autorizado {
		t.Error("el acceso no debería haberse autorizado: el cliente no tiene pago de entrada")
	}
	if repo.accesoRecibido.Motivo == "" {
		t.Error("se esperaba un motivo explicando por qué no se autorizó el acceso")
	}
	if resultado.PagoAlDia {
		t.Error("PagoAlDia debería ser false: el cliente no tiene pago de entrada")
	}
}

// TestCrearAcceso_ClienteInexistenteNoLlegaAlRepo prueba que un dato inválido
// (un cliente_id que no existe) se rechaza ANTES de tocar el repositorio.
// Si esta validación se rompiera, CrearAcceso intentaría guardar un acceso
// para un cliente fantasma y crearAccesoLlamado terminaría en true.
func TestCrearAcceso_ClienteInexistenteNoLlegaAlRepo(t *testing.T) {
	repo := &mockSeguridadRepo{}
	clientes := &mockClienteRepo{clientes: map[int]models.Cliente{}} // sin clientes
	pagos := &mockPagoRepo{tienePago: true}

	svc := NewSeguridadService(repo, clientes, pagos)

	_, err := svc.CrearAcceso(99)
	if err != ErrClienteInvalido {
		t.Fatalf("se esperaba ErrClienteInvalido, se obtuvo: %v", err)
	}
	if repo.crearAccesoLlamado {
		t.Error("CrearAcceso NO debió llegar al repositorio: el cliente_id no existe")
	}
}
