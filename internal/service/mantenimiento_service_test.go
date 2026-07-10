package service

import (
	"testing"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type mantenimientoRepoMock struct {
	equipos            map[uint]models.Equipo
	registros          map[uint]models.RegistroMantenimiento
	quimicos           map[uint]models.ProductoQuimico
	crearEquipoLlamado bool
	crearRegistro      bool
	crearEquipoError   error
	crearRegistroError error
	crearQuimicoError  error
}

var _ storage.MantenimientoRepository = (*mantenimientoRepoMock)(nil)

func newMantenimientoRepoMock() *mantenimientoRepoMock {
	return &mantenimientoRepoMock{
		equipos:   make(map[uint]models.Equipo),
		registros: make(map[uint]models.RegistroMantenimiento),
		quimicos:  make(map[uint]models.ProductoQuimico),
	}
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

func (m *mantenimientoRepoMock) ListarRegistros() []models.RegistroMantenimiento {
	lista := make([]models.RegistroMantenimiento, 0, len(m.registros))
	for _, r := range m.registros {
		lista = append(lista, r)
	}
	return lista
}
func (m *mantenimientoRepoMock) BuscarRegistroPorID(id uint) (models.RegistroMantenimiento, bool) {
	r, ok := m.registros[id]
	return r, ok
}
func (m *mantenimientoRepoMock) CrearRegistro(r models.RegistroMantenimiento) (models.RegistroMantenimiento, error) {
	m.crearRegistro = true
	if m.crearRegistroError != nil {
		return models.RegistroMantenimiento{}, m.crearRegistroError
	}
	r.ID = uint(len(m.registros) + 1)
	m.registros[r.ID] = r
	return r, nil
}
func (m *mantenimientoRepoMock) ActualizarRegistro(id uint, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool) {
	r, ok := m.registros[id]
	if !ok {
		return models.RegistroMantenimiento{}, false
	}
	if datos.EquipoID != 0 {
		r.EquipoID = datos.EquipoID
	}
	if datos.Tipo != "" {
		r.Tipo = datos.Tipo
	}
	if datos.Descripcion != "" {
		r.Descripcion = datos.Descripcion
	}
	if datos.RealizadoPor != "" {
		r.RealizadoPor = datos.RealizadoPor
	}
	m.registros[id] = r
	return r, true
}
func (m *mantenimientoRepoMock) BorrarRegistro(id uint) bool {
	_, ok := m.registros[id]
	if !ok {
		return false
	}
	delete(m.registros, id)
	return true
}

func (m *mantenimientoRepoMock) ListarQuimicos() []models.ProductoQuimico {
	lista := make([]models.ProductoQuimico, 0, len(m.quimicos))
	for _, q := range m.quimicos {
		lista = append(lista, q)
	}
	return lista
}
func (m *mantenimientoRepoMock) BuscarQuimicoPorID(id uint) (models.ProductoQuimico, bool) {
	q, ok := m.quimicos[id]
	return q, ok
}
func (m *mantenimientoRepoMock) CrearQuimico(q models.ProductoQuimico) (models.ProductoQuimico, error) {
	if m.crearQuimicoError != nil {
		return models.ProductoQuimico{}, m.crearQuimicoError
	}
	q.ID = uint(len(m.quimicos) + 1)
	m.quimicos[q.ID] = q
	return q, nil
}
func (m *mantenimientoRepoMock) ActualizarQuimico(id uint, datos models.ProductoQuimico) (models.ProductoQuimico, bool) {
	q, ok := m.quimicos[id]
	if !ok {
		return models.ProductoQuimico{}, false
	}
	if datos.Nombre != "" {
		q.Nombre = datos.Nombre
	}
	m.quimicos[id] = q
	return q, true
}
func (m *mantenimientoRepoMock) BorrarQuimico(id uint) bool {
	_, ok := m.quimicos[id]
	if !ok {
		return false
	}
	delete(m.quimicos, id)
	return true
}

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

func TestClientesService_ListarEquiposDisponibles(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)

	_ = svc.ListarEquipos()
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
