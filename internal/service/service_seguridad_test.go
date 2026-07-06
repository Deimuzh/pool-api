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
	guardavidas        map[int]models.Guardavida
}

func (m *mockSeguridadRepo) ListarGuardavidas() []models.Guardavida { return nil }
func (m *mockSeguridadRepo) BuscarGuardavidaPorID(id uint) (models.Guardavida, bool) {
	g, ok := m.guardavidas[int(id)]
	return g, ok

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
func (m *mockSeguridadRepo) CrearIncidente(i models.Incidente) models.Incidente {
	i.ID = 1
	return i
}
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
