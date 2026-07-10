package service

import (
	"testing"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type mantenimientoRepoMock struct {
	equipos            map[uint]models.Equipo
	crearEquipoLlamado bool
	crearRegistro      bool
	crearEquipoError   error
}

// Verificación en compilación de que implementa la interfaz:
var _ storage.MantenimientoRepository = (*mantenimientoRepoMock)(nil)

func newMantenimientoRepoMock() *mantenimientoRepoMock {
	return &mantenimientoRepoMock{equipos: make(map[uint]models.Equipo)}
}

func (m *mantenimientoRepoMock) ListarEquipos() []models.Equipo {
	lista := make([]models.Equipo, 0, len(m.equipos))
	for _, e := range m.equipos {
		lista = append(lista, e)
	}
	return lista
}
func (m *mantenimientoRepoMock) BuscarEquipoPorID(id uint) (models.Equipo, bool) {
	e, ok := m.equipos[id]
	return e, ok
}
func (m *mantenimientoRepoMock) CrearEquipo(e models.Equipo) (models.Equipo, error) {
	m.crearEquipoLlamado = true
	if m.crearEquipoError != nil {
		return models.Equipo{}, m.crearEquipoError
	}
	e.ID = uint(len(m.equipos) + 1)
	m.equipos[e.ID] = e
	return e, nil
}
func (m *mantenimientoRepoMock) ActualizarEquipo(id uint, datos models.Equipo) (models.Equipo, bool) {
	e, ok := m.equipos[id]
	if !ok {
		return models.Equipo{}, false
	}
	if datos.Nombre != "" {
		e.Nombre = datos.Nombre
	}
	if datos.Tipo != "" {
		e.Tipo = datos.Tipo
	}
	if datos.Estado != "" {
		e.Estado = datos.Estado
	}
	m.equipos[id] = e
	return e, true
}
func (m *mantenimientoRepoMock) BorrarEquipo(id uint) bool {
	_, ok := m.equipos[id]
	if !ok {
		return false
	}
	delete(m.equipos, id)
	return true
}

func (m *mantenimientoRepoMock) ListarRegistros() []models.RegistroMantenimiento { return nil }
func (m *mantenimientoRepoMock) BuscarRegistroPorID(id uint) (models.RegistroMantenimiento, bool) {
	return models.RegistroMantenimiento{}, false
}
func (m *mantenimientoRepoMock) CrearRegistro(r models.RegistroMantenimiento) (models.RegistroMantenimiento, error) {
	m.crearRegistro = true
	r.ID = 1
	return r, nil
}
func (m *mantenimientoRepoMock) ActualizarRegistro(id uint, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool) {
	return models.RegistroMantenimiento{}, false
}
func (m *mantenimientoRepoMock) BorrarRegistro(id uint) bool { return false }

func (m *mantenimientoRepoMock) ListarQuimicos() []models.ProductoQuimico { return nil }
func (m *mantenimientoRepoMock) BuscarQuimicoPorID(id uint) (models.ProductoQuimico, bool) {
	return models.ProductoQuimico{}, false
}
func (m *mantenimientoRepoMock) CrearQuimico(q models.ProductoQuimico) (models.ProductoQuimico, error) {
	q.ID = 1
	return q, nil
}
func (m *mantenimientoRepoMock) ActualizarQuimico(id uint, datos models.ProductoQuimico) (models.ProductoQuimico, bool) {
	return models.ProductoQuimico{}, false
}
func (m *mantenimientoRepoMock) BorrarQuimico(id uint) bool { return false }

func TestMantenimientoService_CrearEquipo_SinTipoNoLlegaAlRepo(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)

	_, err := svc.CrearEquipo(models.Equipo{Nombre: "Filtro"})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
	if repo.crearEquipoLlamado {
		t.Fatal("no debe crear equipo si falta el tipo")
	}
}

func TestMantenimientoService_CrearEquipo_AsignaEstadoOperativo(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)

	creado, err := svc.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.Estado != "operativo" {
		t.Fatalf("se esperaba estado operativo, se obtuvo %q", creado.Estado)
	}
}

func TestMantenimientoService_CrearRegistro_EquipoInvalido(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)

	_, err := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: 99, Tipo: "preventivo"})
	if err != ErrEquipoInvalido {
		t.Fatalf("se esperaba ErrEquipoInvalido, se obtuvo %v", err)
	}
	if repo.crearRegistro {
		t.Fatal("no debe crear registro si el equipo no existe")
	}
}
