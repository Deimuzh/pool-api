package service

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"gorm.io/gorm"

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

func validarEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var reCedula = regexp.MustCompile(`^\d{10}$`)

func (s *ClientesService) CrearCliente(c models.Cliente) (models.Cliente, error) {
	c.Nombre = strings.TrimSpace(c.Nombre)
	c.Cedula = strings.TrimSpace(c.Cedula)
	c.Email = strings.TrimSpace(c.Email)
	c.Telefono = strings.TrimSpace(c.Telefono)
	if c.Nombre == "" || c.Cedula == "" {
		return models.Cliente{}, ErrCampoObligatorio
	}
	if !reCedula.MatchString(c.Cedula) {
		return models.Cliente{}, ErrCedulaFormatoInvalido
	}
	if c.Email != "" && !validarEmail(c.Email) {
		return models.Cliente{}, ErrEmailFormatoInvalido
	}
	if c.Membresia == "" {
		c.Membresia = "ninguna"
	}
	creado, err := s.repo.CrearCliente(c)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") ||
			strings.Contains(err.Error(), "SQLSTATE 23505") ||
			errors.Is(err, gorm.ErrDuplicatedKey) {
			return models.Cliente{}, ErrCedulaEnUso
		}
		return models.Cliente{}, err
	}
	return creado, nil
}

func (s *ClientesService) ActualizarCliente(id uint, c models.Cliente) (models.Cliente, error) {
	c.Nombre = strings.TrimSpace(c.Nombre)
	c.Cedula = strings.TrimSpace(c.Cedula)
	c.Email = strings.TrimSpace(c.Email)
	c.Telefono = strings.TrimSpace(c.Telefono)
	if c.Nombre == "" || c.Cedula == "" {
		return models.Cliente{}, ErrCampoObligatorio
	}
	if !reCedula.MatchString(c.Cedula) {
		return models.Cliente{}, ErrCedulaFormatoInvalido
	}
	if c.Email != "" && !validarEmail(c.Email) {
		return models.Cliente{}, ErrEmailFormatoInvalido
	}
	existente, ok := s.repo.BuscarClientePorID(id)
	if !ok {
		return models.Cliente{}, ErrNoEncontrado
	}
	if existente.Cedula != c.Cedula {
		for _, cl := range s.repo.ListarClientes() {
			if cl.Cedula == c.Cedula {
				return models.Cliente{}, ErrCedulaEnUso
			}
		}
	}
	actualizado, ok := s.repo.ActualizarCliente(id, c)
	if !ok {
		return models.Cliente{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *ClientesService) BorrarCliente(id uint) error {
	for _, rv := range s.repo.ListarReservas() {
		if rv.ClienteID == id {
			return ErrClienteConReservas
		}
	}
	for _, p := range s.repo.ListarPagos() {
		if p.ClienteID == id {
			return ErrClienteConPagos
		}
	}
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
	cliente, ok := s.repo.BuscarClientePorID(rv.ClienteID)
	if !ok {
		return models.Reserva{}, ErrClienteInvalido
	}
	if !clienteTieneMembresia(cliente) {
		return models.Reserva{}, ErrClienteSinMembresia
	}
	if rv.Duracion != 720 && rv.Duracion != 1440 {
		return models.Reserva{}, ErrDuracionInvalida
	}
	if rv.Estado == "" {
		rv.Estado = "pendiente"
	}
	creado, err := s.repo.CrearReserva(rv)
	if err != nil {
		return models.Reserva{}, err
	}
	return creado, nil
}

func (s *ClientesService) ActualizarReserva(id uint, rv models.Reserva) (models.Reserva, error) {
	if rv.ClienteID == 0 {
		return models.Reserva{}, ErrCampoObligatorio
	}
	cliente, ok := s.repo.BuscarClientePorID(rv.ClienteID)
	if !ok {
		return models.Reserva{}, ErrClienteInvalido
	}
	if !clienteTieneMembresia(cliente) {
		return models.Reserva{}, ErrClienteSinMembresia
	}
	if rv.Duracion != 720 && rv.Duracion != 1440 {
		return models.Reserva{}, ErrDuracionInvalida
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
	cliente, ok := s.repo.BuscarClientePorID(p.ClienteID)
	if !ok {
		return models.Pago{}, ErrClienteInvalido
	}
	if clienteTieneMembresia(cliente) {
		return models.Pago{}, ErrClienteConMembresia
	}
	p.Concepto = strings.TrimSpace(strings.ToLower(p.Concepto))
	if p.Concepto != "medio_dia" && p.Concepto != "dia" {
		return models.Pago{}, ErrConceptoPagoInvalido
	}
	p.Metodo = strings.TrimSpace(strings.ToLower(p.Metodo))
	if p.Metodo != "efectivo" && p.Metodo != "transferencia" {
		return models.Pago{}, ErrMetodoPagoInvalido
	}
	creado, err := s.repo.CrearPago(p)
	if err != nil {
		return models.Pago{}, err
	}
	return creado, nil
}

func (s *ClientesService) ActualizarPago(id uint, p models.Pago) (models.Pago, error) {
	if p.ClienteID == 0 {
		return models.Pago{}, ErrCampoObligatorio
	}
	if p.Monto <= 0 {
		return models.Pago{}, ErrMontoInvalido
	}
	cliente, ok := s.repo.BuscarClientePorID(p.ClienteID)
	if !ok {
		return models.Pago{}, ErrClienteInvalido
	}
	if clienteTieneMembresia(cliente) {
		return models.Pago{}, ErrClienteConMembresia
	}
	p.Concepto = strings.TrimSpace(strings.ToLower(p.Concepto))
	if p.Concepto != "medio_dia" && p.Concepto != "dia" {
		return models.Pago{}, ErrConceptoPagoInvalido
	}
	p.Metodo = strings.TrimSpace(strings.ToLower(p.Metodo))
	if p.Metodo != "efectivo" && p.Metodo != "transferencia" {
		return models.Pago{}, ErrMetodoPagoInvalido
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

func clienteTieneMembresia(c models.Cliente) bool {
	m := strings.TrimSpace(strings.ToLower(c.Membresia))
	return m != "" && m != "ninguna"
}
