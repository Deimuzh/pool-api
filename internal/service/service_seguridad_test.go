package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pool-api/internal/models"
	"pool-api/internal/storage"
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
	guardavidas        map[int]models.Guardavida
	accesos            []models.AccesoCliente
	incidentes         map[int]models.Incidente
}

func (m *mockSeguridadRepo) ListarGuardavidas() []models.Guardavida { return nil }
func (m *mockSeguridadRepo) BuscarGuardavidaPorID(id uint) (models.Guardavida, bool) {
	g, ok := m.guardavidas[int(id)]
	return g, ok

}
func (m *mockSeguridadRepo) CrearGuardavida(g models.Guardavida) models.Guardavida {
	g.ID = 1
	return g
}
func (m *mockSeguridadRepo) ActualizarGuardavida(id uint, datos models.Guardavida) (models.Guardavida, bool) {
	if m.guardavidas == nil {
		return models.Guardavida{}, false
	}
	if _, ok := m.guardavidas[int(id)]; !ok {
		return models.Guardavida{}, false
	}
	datos.ID = id
	return datos, true
}
func (m *mockSeguridadRepo) BorrarGuardavida(id uint) bool { return false }

func (m *mockSeguridadRepo) ListarIncidentes() []models.Incidente { return nil }
func (m *mockSeguridadRepo) BuscarIncidentePorID(id uint) (models.Incidente, bool) {
	i, ok := m.incidentes[int(id)]
	return i, ok
}
func (m *mockSeguridadRepo) CrearIncidente(i models.Incidente) models.Incidente {
	i.ID = 1
	return i
}
func (m *mockSeguridadRepo) ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool) {
	if m.incidentes == nil {
		return models.Incidente{}, false
	}
	if _, ok := m.incidentes[int(id)]; !ok {
		return models.Incidente{}, false
	}
	datos.ID = id
	return datos, true
}
func (m *mockSeguridadRepo) BorrarIncidente(id uint) bool { return false }

func (m *mockSeguridadRepo) ListarAccesos() []models.AccesoCliente {
	return m.accesos
}
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
func (m *mockClienteRepo) CrearCliente(c models.Cliente) (models.Cliente, error) { return c, nil }
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

func (m *mockPagoRepo) ListarPagos() []models.Pago                   { return nil }
func (m *mockPagoRepo) BuscarPagoPorID(id uint) (models.Pago, bool)  { return models.Pago{}, false }
func (m *mockPagoRepo) CrearPago(p models.Pago) (models.Pago, error) { return p, nil }
func (m *mockPagoRepo) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}
func (m *mockPagoRepo) BorrarPago(id uint) bool { return false }
func (m *mockPagoRepo) ClienteTienePagoEntrada(clienteID uint) bool {
	return m.tienePago
}

// ─── TEST ───────────────────────────────────────────────────────────────────

// TestCrearAcceso_SinPagoNoLlegaAlRepo prueba que un cliente sin membresía ni
// pago no puede registrar acceso.
func TestCrearAcceso_SinPagoNoLlegaAlRepo(t *testing.T) {
	repo := &mockSeguridadRepo{}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			2: {ID: 2, Nombre: "Luis Pino"},
		},
	}
	pagos := &mockPagoRepo{tienePago: false} // el cliente 2 NO ha pagado entrada

	svc := NewSeguridadService(repo, clientes, pagos)

	_, err := svc.CrearAcceso(2)
	if err != ErrClienteSinAcceso {
		t.Fatalf("se esperaba ErrClienteSinAcceso, se obtuvo: %v", err)
	}

	if repo.crearAccesoLlamado {
		t.Error("CrearAcceso NO debió llegar al repositorio: el cliente no tiene pago")
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

// TestCrearAcceso_ClienteConMembresiaNoRegistraAcceso prueba que el flujo de
// acceso diario rechaza clientes con membresía, porque ellos no necesitan pagar
// entrada diaria.
func TestCrearAcceso_ClienteConMembresiaNoRegistraAcceso(t *testing.T) {
	repo := &mockSeguridadRepo{}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			1: {ID: 1, Nombre: "Ana Reyes", Membresia: "mensual"},
		},
	}
	pagos := &mockPagoRepo{tienePago: true}

	svc := NewSeguridadService(repo, clientes, pagos)

	_, err := svc.CrearAcceso(1)
	if err != ErrClienteConMembresia {
		t.Fatalf("se esperaba ErrClienteConMembresia, se obtuvo: %v", err)
	}

	if repo.crearAccesoLlamado {
		t.Error("CrearAcceso NO debió llegar al repositorio: el cliente tiene membresía")
	}
}

// TestCrearAcceso_ConPagoRegistraAcceso prueba el camino feliz: un cliente sin
// membresía pero con pago de entrada sí llega al repositorio y queda autorizado.
func TestCrearAcceso_ConPagoRegistraAcceso(t *testing.T) {
	repo := &mockSeguridadRepo{}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			2: {ID: 2, Nombre: "Luis Pino", Membresia: "ninguna"},
		},
	}
	pagos := &mockPagoRepo{tienePago: true}

	svc := NewSeguridadService(repo, clientes, pagos)

	creado, err := svc.CrearAcceso(2)
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}

	if !repo.crearAccesoLlamado {
		t.Fatal("CrearAcceso debió llegar al repositorio porque el cliente tiene pago")
	}

	if creado.ClienteID != 2 {
		t.Errorf("cliente_id inesperado: %d", creado.ClienteID)
	}

	if !creado.Autorizado {
		t.Error("se esperaba que el acceso quede autorizado")
	}

	if creado.NombreCliente != "Luis Pino" {
		t.Errorf("nombre_cliente inesperado: %s", creado.NombreCliente)
	}
}

// TestCrearIncidente_GuardavidaInvalidoNoLlegaAlRepo prueba que no se registra
// un incidente si el guardavida responsable no existe.
func TestCrearIncidente_GuardavidaInvalidoNoLlegaAlRepo(t *testing.T) {
	repo := &mockSeguridadRepo{}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			1: {ID: 1, Nombre: "Ana Reyes", Membresia: "mensual"},
		},
	}
	pagos := &mockPagoRepo{tienePago: false}

	svc := NewSeguridadService(repo, clientes, pagos)

	_, err := svc.CrearIncidente(models.Incidente{
		Tipo:         "lesion",
		Gravedad:     "leve",
		GuardavidaID: 99,
		ClienteID:    1,
	})
	if err != ErrGuardavidaInvalido {
		t.Fatalf("se esperaba ErrGuardavidaInvalido, se obtuvo: %v", err)
	}
}

// TestCrearIncidente_ClienteSinAccesoNoLlegaAlRepo prueba que un cliente sin
// membresía y sin acceso registrado no puede asociarse a un incidente.
func TestCrearIncidente_ClienteSinAccesoNoLlegaAlRepo(t *testing.T) {
	repo := &mockSeguridadRepo{
		guardavidas: map[int]models.Guardavida{
			1: {ID: 1, Nombre: "Carlos Mendoza"},
		},
	}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			2: {ID: 2, Nombre: "Luis Pino", Membresia: "ninguna"},
		},
	}
	pagos := &mockPagoRepo{tienePago: false}

	svc := NewSeguridadService(repo, clientes, pagos)

	_, err := svc.CrearIncidente(models.Incidente{
		Tipo:         "lesion",
		Gravedad:     "leve",
		GuardavidaID: 1,
		ClienteID:    2,
	})
	if err != ErrClienteSinAcceso {
		t.Fatalf("se esperaba ErrClienteSinAcceso, se obtuvo: %v", err)
	}
}

// TestCrearIncidente_ClienteConMembresiaRegistraIncidente prueba el camino
// feliz para incidentes: si el cliente tiene membresía, no necesita acceso diario.
func TestCrearIncidente_ClienteConMembresiaRegistraIncidente(t *testing.T) {
	repo := &mockSeguridadRepo{
		guardavidas: map[int]models.Guardavida{
			1: {ID: 1, Nombre: "Carlos Mendoza"},
		},
	}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			1: {ID: 1, Nombre: "Ana Reyes", Membresia: "mensual"},
		},
	}
	pagos := &mockPagoRepo{tienePago: false}

	svc := NewSeguridadService(repo, clientes, pagos)

	creado, err := svc.CrearIncidente(models.Incidente{
		Tipo:         "lesion",
		Gravedad:     "leve",
		GuardavidaID: 1,
		ClienteID:    1,
	})
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}

	if creado.ClienteID != 1 {
		t.Errorf("cliente_id inesperado: %d", creado.ClienteID)
	}
	if creado.NombreCliente != "Ana Reyes" {
		t.Errorf("nombre_cliente inesperado: %s", creado.NombreCliente)
	}
	if creado.NombreGuardavida != "Carlos Mendoza" {
		t.Errorf("nombre_guardavida inesperado: %s", creado.NombreGuardavida)
	}
}

// TestBorrarGuardavida_InexistenteDevuelveNoEncontrado prueba que el service
// traduzca el false del repositorio a un error de dominio.
func TestBorrarGuardavida_InexistenteDevuelveNoEncontrado(t *testing.T) {
	repo := &mockSeguridadRepo{}
	clientes := &mockClienteRepo{clientes: map[int]models.Cliente{}}
	pagos := &mockPagoRepo{}

	svc := NewSeguridadService(repo, clientes, pagos)

	err := svc.BorrarGuardavida(99)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo: %v", err)
	}
}

// TestCrearIncidente_CamposObligatoriosNoLlegaAlRepo prueba que el service
// rechaza incidentes incompletos antes de validar guardavida o cliente.
func TestCrearIncidente_CamposObligatoriosNoLlegaAlRepo(t *testing.T) {
	repo := &mockSeguridadRepo{}
	clientes := &mockClienteRepo{clientes: map[int]models.Cliente{}}
	pagos := &mockPagoRepo{}

	svc := NewSeguridadService(repo, clientes, pagos)

	_, err := svc.CrearIncidente(models.Incidente{
		Tipo:         "",
		Gravedad:     "",
		GuardavidaID: 0,
		ClienteID:    0,
	})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo: %v", err)
	}
}

// TestCrearGuardavida_CamposValidos verifica que un guardavida con nombre y turno válidos
// se cree correctamente y reciba un ID desde el repositorio mock.
func TestCrearGuardavida_CamposValidos(t *testing.T) {
	repo := &mockSeguridadRepo{}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})

	creado, err := svc.CrearGuardavida(models.Guardavida{
		Nombre: "Pedro Salazar",
		Turno:  "tarde",
	})

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.ID != 1 {
		t.Fatalf("se esperaba ID 1, se obtuvo %d", creado.ID)
	}
}

// TestCrearGuardavida_CamposObligatorios verifica que el service rechace la creación
// cuando falta un campo obligatorio, en este caso el turno.
func TestCrearGuardavida_CamposObligatorios(t *testing.T) {
	repo := &mockSeguridadRepo{}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})

	_, err := svc.CrearGuardavida(models.Guardavida{
		Nombre: "Pedro Salazar",
	})

	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

// TestActualizarGuardavida_Valido verifica que un guardavida existente pueda actualizarse
// correctamente cuando los datos enviados son válidos.
func TestActualizarGuardavida_Valido(t *testing.T) {
	repo := &mockSeguridadRepo{
		guardavidas: map[int]models.Guardavida{
			1: {ID: 1, Nombre: "Pedro", Turno: "mañana"},
		},
	}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})

	actualizado, err := svc.ActualizarGuardavida(1, models.Guardavida{
		Nombre: "Pedro Actualizado",
		Turno:  "noche",
	})

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Nombre != "Pedro Actualizado" {
		t.Fatalf("nombre inesperado: %s", actualizado.Nombre)
	}
}

// TestActualizarGuardavida_Inexistente verifica que el service devuelva ErrNoEncontrado
// cuando se intenta actualizar un guardavida que no existe en el repositorio.
func TestActualizarGuardavida_Inexistente(t *testing.T) {
	repo := &mockSeguridadRepo{
		guardavidas: map[int]models.Guardavida{},
	}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})

	_, err := svc.ActualizarGuardavida(99, models.Guardavida{
		Nombre: "Pedro",
		Turno:  "tarde",
	})

	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

// TestObtenerIncidente_Existente verifica que un incidente existente se obtenga
// enriquecido con el nombre del guardavida y del cliente.
func TestObtenerIncidente_Existente(t *testing.T) {
	repo := &mockSeguridadRepo{
		guardavidas: map[int]models.Guardavida{
			1: {ID: 1, Nombre: "Carlos Mendoza"},
		},
		incidentes: map[int]models.Incidente{
			10: {ID: 10, Tipo: "lesion", Gravedad: "leve", GuardavidaID: 1, ClienteID: 2},
		},
	}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			2: {ID: 2, Nombre: "Ana Reyes", Membresia: "mensual"},
		},
	}
	svc := NewSeguridadService(repo, clientes, &mockPagoRepo{})

	incidente, ok := svc.ObtenerIncidente(10)

	if !ok {
		t.Fatal("se esperaba encontrar el incidente")
	}
	if incidente.NombreGuardavida != "Carlos Mendoza" {
		t.Fatalf("nombre de guardavida inesperado: %s", incidente.NombreGuardavida)
	}
	if incidente.NombreCliente != "Ana Reyes" {
		t.Fatalf("nombre de cliente inesperado: %s", incidente.NombreCliente)
	}
}

// TestActualizarIncidente_Valido verifica que un incidente existente pueda actualizarse
// cuando el guardavida y el cliente cumplen las reglas de negocio.
func TestActualizarIncidente_Valido(t *testing.T) {
	repo := &mockSeguridadRepo{
		guardavidas: map[int]models.Guardavida{
			1: {ID: 1, Nombre: "Carlos Mendoza"},
		},
		incidentes: map[int]models.Incidente{
			5: {ID: 5, Tipo: "lesion", Gravedad: "leve", GuardavidaID: 1, ClienteID: 2},
		},
	}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			2: {ID: 2, Nombre: "Ana Reyes", Membresia: "mensual"},
		},
	}
	svc := NewSeguridadService(repo, clientes, &mockPagoRepo{})

	actualizado, err := svc.ActualizarIncidente(5, models.Incidente{
		Tipo:         "rescate",
		Gravedad:     "media",
		GuardavidaID: 1,
		ClienteID:    2,
	})

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Tipo != "rescate" {
		t.Fatalf("tipo inesperado: %s", actualizado.Tipo)
	}
	if actualizado.NombreCliente != "Ana Reyes" {
		t.Fatalf("nombre de cliente inesperado: %s", actualizado.NombreCliente)
	}
}

// TestListarAccesos_EnriqueceNombreYPago verifica que ListarAccesos devuelva
// accesos con el nombre del cliente y el estado de pago calculado.
func TestListarAccesos_EnriqueceNombreYPago(t *testing.T) {
	repo := &mockSeguridadRepo{
		accesos: []models.AccesoCliente{
			{ID: 1, ClienteID: 2, Autorizado: true},
		},
	}
	clientes := &mockClienteRepo{
		clientes: map[int]models.Cliente{
			2: {ID: 2, Nombre: "Luis Pino", Membresia: "ninguna"},
		},
	}
	pagos := &mockPagoRepo{tienePago: true}
	svc := NewSeguridadService(repo, clientes, pagos)

	accesos := svc.ListarAccesos()

	if len(accesos) != 1 {
		t.Fatalf("se esperaba 1 acceso, se obtuvieron %d", len(accesos))
	}
	if accesos[0].NombreCliente != "Luis Pino" {
		t.Fatalf("nombre de cliente inesperado: %s", accesos[0].NombreCliente)
	}
	if !accesos[0].PagoAlDia {
		t.Fatal("se esperaba pago al día")
	}
}

// TestBorrarAcceso_Inexistente verifica que el service devuelva ErrNoEncontrado
// cuando se intenta borrar un acceso inexistente.
func TestBorrarAcceso_Inexistente(t *testing.T) {
	repo := &mockSeguridadRepo{}
	svc := NewSeguridadService(repo, &mockClienteRepo{}, &mockPagoRepo{})

	err := svc.BorrarAcceso(99)

	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

// ─── TESTIFY MOCKS ──────────────────────────────────────────────────────────

type seguridadRepoTestifyMock struct {
	mock.Mock
}

var _ storage.SeguridadRepository = (*seguridadRepoTestifyMock)(nil)

func (m *seguridadRepoTestifyMock) ListarGuardavidas() []models.Guardavida {
	args := m.Called()
	return args.Get(0).([]models.Guardavida)
}

func (m *seguridadRepoTestifyMock) BuscarGuardavidaPorID(id uint) (models.Guardavida, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Guardavida), args.Bool(1)
}

func (m *seguridadRepoTestifyMock) CrearGuardavida(g models.Guardavida) models.Guardavida {
	args := m.Called(g)
	return args.Get(0).(models.Guardavida)
}

func (m *seguridadRepoTestifyMock) ActualizarGuardavida(id uint, datos models.Guardavida) (models.Guardavida, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Guardavida), args.Bool(1)
}

func (m *seguridadRepoTestifyMock) BorrarGuardavida(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *seguridadRepoTestifyMock) ListarIncidentes() []models.Incidente {
	args := m.Called()
	return args.Get(0).([]models.Incidente)
}

func (m *seguridadRepoTestifyMock) BuscarIncidentePorID(id uint) (models.Incidente, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Incidente), args.Bool(1)
}

func (m *seguridadRepoTestifyMock) CrearIncidente(i models.Incidente) models.Incidente {
	args := m.Called(i)
	return args.Get(0).(models.Incidente)
}

func (m *seguridadRepoTestifyMock) ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Incidente), args.Bool(1)
}

func (m *seguridadRepoTestifyMock) BorrarIncidente(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *seguridadRepoTestifyMock) ListarAccesos() []models.AccesoCliente {
	args := m.Called()
	return args.Get(0).([]models.AccesoCliente)
}

func (m *seguridadRepoTestifyMock) BuscarAccesoPorID(id uint) (models.AccesoCliente, bool) {
	args := m.Called(id)
	return args.Get(0).(models.AccesoCliente), args.Bool(1)
}

func (m *seguridadRepoTestifyMock) CrearAcceso(a models.AccesoCliente) models.AccesoCliente {
	args := m.Called(a)
	return args.Get(0).(models.AccesoCliente)
}

func (m *seguridadRepoTestifyMock) ActualizarAcceso(id uint, datos models.AccesoCliente) (models.AccesoCliente, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.AccesoCliente), args.Bool(1)
}

func (m *seguridadRepoTestifyMock) BorrarAcceso(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

type clienteRepoTestifyMock struct {
	mock.Mock
}

var _ storage.ClienteRepository = (*clienteRepoTestifyMock)(nil)

func (m *clienteRepoTestifyMock) ListarClientes() []models.Cliente {
	args := m.Called()
	return args.Get(0).([]models.Cliente)
}

func (m *clienteRepoTestifyMock) BuscarClientePorID(id uint) (models.Cliente, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Cliente), args.Bool(1)
}

func (m *clienteRepoTestifyMock) CrearCliente(c models.Cliente) (models.Cliente, error) {
	args := m.Called(c)
	return args.Get(0).(models.Cliente), args.Error(1)
}

func (m *clienteRepoTestifyMock) ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Cliente), args.Bool(1)
}

func (m *clienteRepoTestifyMock) BorrarCliente(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

type pagoRepoTestifyMock struct {
	mock.Mock
}

var _ storage.PagoRepository = (*pagoRepoTestifyMock)(nil)

func (m *pagoRepoTestifyMock) ListarPagos() []models.Pago {
	args := m.Called()
	return args.Get(0).([]models.Pago)
}

func (m *pagoRepoTestifyMock) BuscarPagoPorID(id uint) (models.Pago, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Pago), args.Bool(1)
}

func (m *pagoRepoTestifyMock) CrearPago(p models.Pago) (models.Pago, error) {
	args := m.Called(p)
	return args.Get(0).(models.Pago), args.Error(1)
}

func (m *pagoRepoTestifyMock) ActualizarPago(id uint, datos models.Pago) (models.Pago, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Pago), args.Bool(1)
}

func (m *pagoRepoTestifyMock) BorrarPago(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *pagoRepoTestifyMock) ClienteTienePagoEntrada(clienteID uint) bool {
	args := m.Called(clienteID)
	return args.Bool(0)
}

func TestSeguridadService_CrearAcceso_ConTestifyMock(t *testing.T) {
	seguridad := new(seguridadRepoTestifyMock)
	clientes := new(clienteRepoTestifyMock)
	pagos := new(pagoRepoTestifyMock)

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
