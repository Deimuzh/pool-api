package service

import (
	"pool-api/internal/models"
	"pool-api/internal/storage"
)

// ClientesService agrupa la lógica de negocio de Clientes, Reservas y Pagos.
type ClientesService struct {
	repo storage.ClientesModulo
}

func NewClientesService(repo storage.ClientesModulo) *ClientesService {
	return &ClientesService{repo: repo}
}

// ─── CLIENTE ──────────────────────────────────────────────────────────────────

func (s *ClientesService) ListarClientes() []models.Cliente {
	return s.repo.ListarClientes()
}

func (s *ClientesService) ObtenerCliente(id uint) (models.Cliente, bool) {
	return s.repo.BuscarClientePorID(id)
}

func (s *ClientesService) CrearCliente(c models.Cliente) (models.Cliente, error) {
	if c.Nombre == "" || c.Cedula == "" {
		return models.Cliente{}, ErrCampoObligatorio
	}
	if c.Membresia == "" {
		c.Membresia = "ninguna"
	}
	return s.repo.CrearCliente(c), nil
}

func (s *ClientesService) ActualizarCliente(id uint, c models.Cliente) (models.Cliente, error) {
	if c.Nombre == "" || c.Cedula == "" {
		return models.Cliente{}, ErrCampoObligatorio
	}
	actualizado, ok := s.repo.ActualizarCliente(id, c)
	if !ok {
		return models.Cliente{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *ClientesService) BorrarCliente(id uint) error {
	if !s.repo.BorrarCliente(id) {
		return ErrNoEncontrado
	}
	return nil
}

// ─── RESERVA ──────────────────────────────────────────────────────────────────

func (s *ClientesService) ListarReservas() []models.Reserva {
	return s.repo.ListarReservas()
}

func (s *ClientesService) ObtenerReserva(id uint) (models.Reserva, bool) {
	return s.repo.BuscarReservaPorID(id)
}

func (s *ClientesService) CrearReserva(rv models.Reserva) (models.Reserva, error) {
	if rv.ClienteID == 0 {
		return models.Reserva{}, ErrCampoObligatorio
	}
	if _, ok := s.repo.BuscarClientePorID(rv.ClienteID); !ok {
		return models.Reserva{}, ErrClienteInvalido
	}
	if rv.Estado == "" {
		rv.Estado = "pendiente"
	}
	return s.repo.CrearReserva(rv), nil
}

func (s *ClientesService) ActualizarReserva(id uint, rv models.Reserva) (models.Reserva, error) {
	if rv.ClienteID == 0 {
		return models.Reserva{}, ErrCampoObligatorio
	}
	actualizado, ok := s.repo.ActualizarReserva(id, rv)
	if !ok {
		return models.Reserva{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *ClientesService) BorrarReserva(id uint) error {
	if !s.repo.BorrarReserva(id) {
		return ErrNoEncontrado
	}
	return nil
}

// ─── PAGO ─────────────────────────────────────────────────────────────────────

func (s *ClientesService) ListarPagos() []models.Pago {
	return s.repo.ListarPagos()
}

func (s *ClientesService) ObtenerPago(id uint) (models.Pago, bool) {
	return s.repo.BuscarPagoPorID(id)
}

func (s *ClientesService) CrearPago(p models.Pago) (models.Pago, error) {
	if p.ClienteID == 0 {
		return models.Pago{}, ErrCampoObligatorio
	}
	if p.Monto <= 0 {
		return models.Pago{}, ErrMontoInvalido
	}
	if _, ok := s.repo.BuscarClientePorID(p.ClienteID); !ok {
		return models.Pago{}, ErrClienteInvalido
	}
	return s.repo.CrearPago(p), nil
}

func (s *ClientesService) ActualizarPago(id uint, p models.Pago) (models.Pago, error) {
	if p.ClienteID == 0 || p.Monto <= 0 {
		return models.Pago{}, ErrCampoObligatorio
	}
	actualizado, ok := s.repo.ActualizarPago(id, p)
	if !ok {
		return models.Pago{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *ClientesService) BorrarPago(id uint) error {
	if !s.repo.BorrarPago(id) {
		return ErrNoEncontrado
	}
	return nil
}
