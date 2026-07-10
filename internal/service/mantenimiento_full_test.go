package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"pool-api/internal/models"
	"pool-api/internal/storage"
)

type mantenimientoRepoCompletoMock struct {
	equipos         map[uint]models.Equipo
	registros       map[uint]models.RegistroMantenimiento
	quimicos        map[uint]models.ProductoQuimico
	sigEquipoID     uint
	sigRegistroID   uint
	sigQuimicoID    uint
	crearEquipoErr  error
	crearRegistroErr error
	crearQuimicoErr error
}

var _ storage.MantenimientoRepository = (*mantenimientoRepoCompletoMock)(nil)

func newMantenimientoRepoCompletoMock() *mantenimientoRepoCompletoMock {
	return &mantenimientoRepoCompletoMock{
		equipos:   make(map[uint]models.Equipo),
		registros: make(map[uint]models.RegistroMantenimiento),
		quimicos:  make(map[uint]models.ProductoQuimico),
	}
}

func (m *mantenimientoRepoCompletoMock) ListarEquipos() []models.Equipo {
	lista := make([]models.Equipo, 0, len(m.equipos))
	for _, e := range m.equipos {
		lista = append(lista, e)
	}
	return lista
}
func (m *mantenimientoRepoCompletoMock) BuscarEquipoPorID(id uint) (models.Equipo, bool) {
	e, ok := m.equipos[id]
	return e, ok
}
func (m *mantenimientoRepoCompletoMock) CrearEquipo(e models.Equipo) (models.Equipo, error) {
	if m.crearEquipoErr != nil {
		return models.Equipo{}, m.crearEquipoErr
	}
	m.sigEquipoID++
	e.ID = m.sigEquipoID
	m.equipos[e.ID] = e
	return e, nil
}
func (m *mantenimientoRepoCompletoMock) ActualizarEquipo(id uint, datos models.Equipo) (models.Equipo, bool) {
	_, ok := m.equipos[id]
	if !ok {
		return models.Equipo{}, false
	}
	datos.ID = id
	m.equipos[id] = datos
	return datos, true
}
func (m *mantenimientoRepoCompletoMock) BorrarEquipo(id uint) bool {
	_, ok := m.equipos[id]
	if !ok {
		return false
	}
	delete(m.equipos, id)
	return true
}

func (m *mantenimientoRepoCompletoMock) ListarRegistros() []models.RegistroMantenimiento {
	lista := make([]models.RegistroMantenimiento, 0, len(m.registros))
	for _, r := range m.registros {
		lista = append(lista, r)
	}
	return lista
}
func (m *mantenimientoRepoCompletoMock) BuscarRegistroPorID(id uint) (models.RegistroMantenimiento, bool) {
	r, ok := m.registros[id]
	return r, ok
}
func (m *mantenimientoRepoCompletoMock) CrearRegistro(r models.RegistroMantenimiento) (models.RegistroMantenimiento, error) {
	if m.crearRegistroErr != nil {
		return models.RegistroMantenimiento{}, m.crearRegistroErr
	}
	m.sigRegistroID++
	r.ID = m.sigRegistroID
	m.registros[r.ID] = r
	return r, nil
}
func (m *mantenimientoRepoCompletoMock) ActualizarRegistro(id uint, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool) {
	_, ok := m.registros[id]
	if !ok {
		return models.RegistroMantenimiento{}, false
	}
	datos.ID = id
	m.registros[id] = datos
	return datos, true
}
func (m *mantenimientoRepoCompletoMock) BorrarRegistro(id uint) bool {
	_, ok := m.registros[id]
	if !ok {
		return false
	}
	delete(m.registros, id)
	return true
}

func (m *mantenimientoRepoCompletoMock) ListarQuimicos() []models.ProductoQuimico {
	lista := make([]models.ProductoQuimico, 0, len(m.quimicos))
	for _, q := range m.quimicos {
		lista = append(lista, q)
	}
	return lista
}
func (m *mantenimientoRepoCompletoMock) BuscarQuimicoPorID(id uint) (models.ProductoQuimico, bool) {
	q, ok := m.quimicos[id]
	return q, ok
}
func (m *mantenimientoRepoCompletoMock) CrearQuimico(q models.ProductoQuimico) (models.ProductoQuimico, error) {
	if m.crearQuimicoErr != nil {
		return models.ProductoQuimico{}, m.crearQuimicoErr
	}
	m.sigQuimicoID++
	q.ID = m.sigQuimicoID
	m.quimicos[q.ID] = q
	return q, nil
}
func (m *mantenimientoRepoCompletoMock) ActualizarQuimico(id uint, datos models.ProductoQuimico) (models.ProductoQuimico, bool) {
	_, ok := m.quimicos[id]
	if !ok {
		return models.ProductoQuimico{}, false
	}
	datos.ID = id
	m.quimicos[id] = datos
	return datos, true
}
func (m *mantenimientoRepoCompletoMock) BorrarQuimico(id uint) bool {
	_, ok := m.quimicos[id]
	if !ok {
		return false
	}
	delete(m.quimicos, id)
	return true
}

func TestMantenimientoService_ListarEquipos(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	svc.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion"})
	svc.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	require.Len(t, svc.ListarEquipos(), 2)
}

func TestMantenimientoService_ObtenerEquipo(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion"})
	e, ok := svc.ObtenerEquipo(creado.ID)
	require.True(t, ok)
	require.Equal(t, "Filtro", e.Nombre)
	_, ok = svc.ObtenerEquipo(99)
	require.False(t, ok)
}

func TestMantenimientoService_ActualizarEquipo(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion"})
	actualizado, err := svc.ActualizarEquipo(creado.ID, models.Equipo{Nombre: "Filtro Nuevo", Tipo: "filtracion"})
	require.NoError(t, err)
	require.Equal(t, "Filtro Nuevo", actualizado.Nombre)
}

func TestMantenimientoService_ActualizarEquipo_CamposObligatorios(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, err := svc.ActualizarEquipo(1, models.Equipo{Nombre: "", Tipo: ""})
	require.ErrorIs(t, err, ErrCampoObligatorio)
}

func TestMantenimientoService_ActualizarEquipo_NoEncontrado(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, err := svc.ActualizarEquipo(99, models.Equipo{Nombre: "Test", Tipo: "test"})
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestMantenimientoService_BorrarEquipo(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion"})
	err := svc.BorrarEquipo(creado.ID)
	require.NoError(t, err)
}

func TestMantenimientoService_BorrarEquipo_NoEncontrado(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	err := svc.BorrarEquipo(99)
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestMantenimientoService_CrearEquipo_ErrorRepo(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	repo.crearEquipoErr = errors.New("error db")
	svc := NewMantenimientoService(repo)
	_, err := svc.CrearEquipo(models.Equipo{Nombre: "Test", Tipo: "test", Estado: "operativo"})
	require.Error(t, err)
}

func TestMantenimientoService_CrearRegistro_Exitoso(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	equipo, _ := svc.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion"})
	creado, err := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: equipo.ID, Tipo: "preventivo"})
	require.NoError(t, err)
	require.NotZero(t, creado.ID)
}

func TestMantenimientoService_CrearRegistro_CamposObligatorios(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, err := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: 0, Tipo: ""})
	require.ErrorIs(t, err, ErrCampoObligatorio)
}

func TestMantenimientoService_ListarRegistros(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	require.Empty(t, svc.ListarRegistros())
}

func TestMantenimientoService_ObtenerRegistro(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, ok := svc.ObtenerRegistro(99)
	require.False(t, ok)
}

func TestMantenimientoService_ActualizarRegistro(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	equipo, _ := svc.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion"})
	creado, _ := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: equipo.ID, Tipo: "preventivo"})
	_, err := svc.ActualizarRegistro(creado.ID, models.RegistroMantenimiento{EquipoID: equipo.ID, Tipo: "correctivo"})
	require.NoError(t, err)
}

func TestMantenimientoService_ActualizarRegistro_NoEncontrado(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, err := svc.ActualizarRegistro(99, models.RegistroMantenimiento{EquipoID: 1, Tipo: "preventivo"})
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestMantenimientoService_BorrarRegistro(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	equipo, _ := svc.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion"})
	creado, _ := svc.CrearRegistro(models.RegistroMantenimiento{EquipoID: equipo.ID, Tipo: "preventivo"})
	err := svc.BorrarRegistro(creado.ID)
	require.NoError(t, err)
}

func TestMantenimientoService_BorrarRegistro_NoEncontrado(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	err := svc.BorrarRegistro(99)
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestMantenimientoService_CrearQuimico(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	creado, err := svc.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro"})
	require.NoError(t, err)
	require.NotZero(t, creado.ID)
}

func TestMantenimientoService_CrearQuimico_NombreVacio(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, err := svc.CrearQuimico(models.ProductoQuimico{Nombre: ""})
	require.ErrorIs(t, err, ErrNombreVacio)
}

func TestMantenimientoService_ListarQuimicos(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	require.Empty(t, svc.ListarQuimicos())
}

func TestMantenimientoService_ObtenerQuimico(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, ok := svc.ObtenerQuimico(99)
	require.False(t, ok)
}

func TestMantenimientoService_ActualizarQuimico(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro"})
	actualizado, err := svc.ActualizarQuimico(creado.ID, models.ProductoQuimico{Nombre: "Cloro Granulado"})
	require.NoError(t, err)
	require.Equal(t, "Cloro Granulado", actualizado.Nombre)
}

func TestMantenimientoService_ActualizarQuimico_NombreVacio(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, err := svc.ActualizarQuimico(1, models.ProductoQuimico{Nombre: ""})
	require.ErrorIs(t, err, ErrNombreVacio)
}

func TestMantenimientoService_ActualizarQuimico_NoEncontrado(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	_, err := svc.ActualizarQuimico(99, models.ProductoQuimico{Nombre: "Test"})
	require.ErrorIs(t, err, ErrNoEncontrado)
}

func TestMantenimientoService_BorrarQuimico(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	creado, _ := svc.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro"})
	err := svc.BorrarQuimico(creado.ID)
	require.NoError(t, err)
}

func TestMantenimientoService_BorrarQuimico_NoEncontrado(t *testing.T) {
	repo := newMantenimientoRepoCompletoMock()
	svc := NewMantenimientoService(repo)
	err := svc.BorrarQuimico(99)
	require.ErrorIs(t, err, ErrNoEncontrado)
}
