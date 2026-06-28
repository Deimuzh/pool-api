package service

import (
	"pool-api/internal/models"
	"pool-api/internal/storage"
)

// MantenimientoService agrupa la lógica de negocio de Equipos, Registros
// de mantenimiento y Productos químicos.
type MantenimientoService struct {
	repo storage.MantenimientoRepository
}

func NewMantenimientoService(repo storage.MantenimientoRepository) *MantenimientoService {
	return &MantenimientoService{repo: repo}
}

// ─── EQUIPO ───────────────────────────────────────────────────────────────────

func (s *MantenimientoService) ListarEquipos() []models.Equipo {
	return s.repo.ListarEquipos()
}

func (s *MantenimientoService) ObtenerEquipo(id uint) (models.Equipo, bool) {
	return s.repo.BuscarEquipoPorID(id)
}

func (s *MantenimientoService) CrearEquipo(e models.Equipo) (models.Equipo, error) {
	if e.Nombre == "" || e.Tipo == "" {
		return models.Equipo{}, ErrCampoObligatorio
	}
	if e.Estado == "" {
		e.Estado = "operativo"
	}
	return s.repo.CrearEquipo(e), nil
}

func (s *MantenimientoService) ActualizarEquipo(id uint, e models.Equipo) (models.Equipo, error) {
	if e.Nombre == "" || e.Tipo == "" {
		return models.Equipo{}, ErrCampoObligatorio
	}
	actualizado, ok := s.repo.ActualizarEquipo(id, e)
	if !ok {
		return models.Equipo{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *MantenimientoService) BorrarEquipo(id uint) error {
	if !s.repo.BorrarEquipo(id) {
		return ErrNoEncontrado
	}
	return nil
}

// ─── REGISTRO MANTENIMIENTO ───────────────────────────────────────────────────

func (s *MantenimientoService) ListarRegistros() []models.RegistroMantenimiento {
	return s.repo.ListarRegistros()
}

func (s *MantenimientoService) ObtenerRegistro(id uint) (models.RegistroMantenimiento, bool) {
	return s.repo.BuscarRegistroPorID(id)
}

// CrearRegistro valida que el equipo exista antes de registrar el mantenimiento.
func (s *MantenimientoService) CrearRegistro(rm models.RegistroMantenimiento) (models.RegistroMantenimiento, error) {
	if rm.EquipoID == 0 || rm.Tipo == "" {
		return models.RegistroMantenimiento{}, ErrCampoObligatorio
	}
	if _, ok := s.repo.BuscarEquipoPorID(rm.EquipoID); !ok {
		return models.RegistroMantenimiento{}, ErrEquipoInvalido
	}
	return s.repo.CrearRegistro(rm), nil
}

func (s *MantenimientoService) ActualizarRegistro(id uint, rm models.RegistroMantenimiento) (models.RegistroMantenimiento, error) {
	if rm.EquipoID == 0 || rm.Tipo == "" {
		return models.RegistroMantenimiento{}, ErrCampoObligatorio
	}
	actualizado, ok := s.repo.ActualizarRegistro(id, rm)
	if !ok {
		return models.RegistroMantenimiento{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *MantenimientoService) BorrarRegistro(id uint) error {
	if !s.repo.BorrarRegistro(id) {
		return ErrNoEncontrado
	}
	return nil
}

// ─── PRODUCTO QUIMICO ─────────────────────────────────────────────────────────

func (s *MantenimientoService) ListarQuimicos() []models.ProductoQuimico {
	return s.repo.ListarQuimicos()
}

func (s *MantenimientoService) ObtenerQuimico(id uint) (models.ProductoQuimico, bool) {
	return s.repo.BuscarQuimicoPorID(id)
}

func (s *MantenimientoService) CrearQuimico(q models.ProductoQuimico) (models.ProductoQuimico, error) {
	if q.Nombre == "" {
		return models.ProductoQuimico{}, ErrNombreVacio
	}
	return s.repo.CrearQuimico(q), nil
}

func (s *MantenimientoService) ActualizarQuimico(id uint, q models.ProductoQuimico) (models.ProductoQuimico, error) {
	if q.Nombre == "" {
		return models.ProductoQuimico{}, ErrNombreVacio
	}
	actualizado, ok := s.repo.ActualizarQuimico(id, q)
	if !ok {
		return models.ProductoQuimico{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *MantenimientoService) BorrarQuimico(id uint) error {
	if !s.repo.BorrarQuimico(id) {
		return ErrNoEncontrado
	}
	return nil
}
