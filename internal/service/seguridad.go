package service

import (
	"pool-api/internal/models"
	"pool-api/internal/storage"
)

// SeguridadService agrupa la lógica de negocio de Guardavidas, Incidentes y Accesos.
//
// Depende de storage.SeguridadRepository para sus propias tablas, y de
// storage.ClienteRepository + storage.PagoRepository para validar y enriquecer
// datos que pertenecen al módulo Clientes (vinculación entre módulos).
type SeguridadService struct {
	repo     storage.SeguridadRepository
	clientes storage.ClienteRepository
	pagos    storage.PagoRepository
}

func NewSeguridadService(repo storage.SeguridadRepository, clientes storage.ClienteRepository, pagos storage.PagoRepository) *SeguridadService {
	return &SeguridadService{repo: repo, clientes: clientes, pagos: pagos}
}

// ─── GUARDAVIDA ───────────────────────────────────────────────────────────────

func (s *SeguridadService) ListarGuardavidas() []models.Guardavida {
	return s.repo.ListarGuardavidas()
}

func (s *SeguridadService) ObtenerGuardavida(id uint) (models.Guardavida, bool) {
	return s.repo.BuscarGuardavidaPorID(id)
}

func (s *SeguridadService) CrearGuardavida(g models.Guardavida) (models.Guardavida, error) {
	if g.Nombre == "" || g.Turno == "" {
		return models.Guardavida{}, ErrCampoObligatorio
	}
	return s.repo.CrearGuardavida(g), nil
}

func (s *SeguridadService) ActualizarGuardavida(id uint, g models.Guardavida) (models.Guardavida, error) {
	if g.Nombre == "" || g.Turno == "" {
		return models.Guardavida{}, ErrCampoObligatorio
	}
	actualizado, ok := s.repo.ActualizarGuardavida(id, g)
	if !ok {
		return models.Guardavida{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *SeguridadService) BorrarGuardavida(id uint) error {
	if !s.repo.BorrarGuardavida(id) {
		return ErrNoEncontrado
	}
	return nil
}

// ─── INCIDENTE ────────────────────────────────────────────────────────────────

// IncidenteConNombre enriquece un Incidente con los nombres legibles
// del guardavida responsable y, si aplica, del cliente involucrado.
type IncidenteConNombre struct {
	models.Incidente
	NombreGuardavida string `json:"nombre_guardavida"`
	NombreCliente    string `json:"nombre_cliente,omitempty"`
}

func (s *SeguridadService) ListarIncidentes() []IncidenteConNombre {
	incidentes := s.repo.ListarIncidentes()
	resultado := make([]IncidenteConNombre, 0, len(incidentes))
	for _, inc := range incidentes {
		resultado = append(resultado, s.enriquecerIncidente(inc))
	}
	return resultado
}

func (s *SeguridadService) ObtenerIncidente(id uint) (IncidenteConNombre, bool) {
	inc, ok := s.repo.BuscarIncidentePorID(id)
	if !ok {
		return IncidenteConNombre{}, false
	}
	return s.enriquecerIncidente(inc), true
}

func (s *SeguridadService) enriquecerIncidente(inc models.Incidente) IncidenteConNombre {
	nombreGuardavida := ""
	if g, ok := s.repo.BuscarGuardavidaPorID(int(inc.GuardavidaID)); ok {
		nombreGuardavida = g.Nombre
	}
	return IncidenteConNombre{
		Incidente:        inc,
		NombreGuardavida: nombreGuardavida,
	}
}

// CrearIncidente valida que el guardavida exista, y si se indicó un cliente
// involucrado (ClienteID != 0), valida que ese cliente también exista.
func (s *SeguridadService) CrearIncidente(inc models.Incidente) (IncidenteConNombre, error) {
	if inc.Tipo == "" || inc.Gravedad == "" || inc.GuardavidaID == 0 {
		return IncidenteConNombre{}, ErrCampoObligatorio
	}
	if _, ok := s.repo.BuscarGuardavidaPorID(int(inc.GuardavidaID)); !ok {
		return IncidenteConNombre{}, ErrGuardavidaInvalido
	}
	creado := s.repo.CrearIncidente(inc)
	return s.enriquecerIncidente(creado), nil
}

func (s *SeguridadService) ActualizarIncidente(id uint, inc models.Incidente) (IncidenteConNombre, error) {
	if inc.Tipo == "" || inc.Gravedad == "" || inc.GuardavidaID == 0 {
		return IncidenteConNombre{}, ErrCampoObligatorio
	}
	actualizado, ok := s.repo.ActualizarIncidente(id, inc)
	if !ok {
		return IncidenteConNombre{}, ErrNoEncontrado
	}
	return s.enriquecerIncidente(actualizado), nil
}

func (s *SeguridadService) BorrarIncidente(id uint) error {
	if !s.repo.BorrarIncidente(id) {
		return ErrNoEncontrado
	}
	return nil
}

// ─── ACCESO CLIENTE ───────────────────────────────────────────────────────────

// AccesoConNombre enriquece un AccesoCliente con el nombre del cliente
// y si tiene o no un pago de entrada registrado (regla de negocio clave).
type AccesoConNombre struct {
	models.AccesoCliente
	NombreCliente string `json:"nombre_cliente"`
	PagoAlDia     bool   `json:"pago_al_dia"`
}

func (s *SeguridadService) ListarAccesos() []AccesoConNombre {
	accesos := s.repo.ListarAccesos()
	resultado := make([]AccesoConNombre, 0, len(accesos))
	for _, a := range accesos {
		resultado = append(resultado, s.enriquecerAcceso(a))
	}
	return resultado
}

func (s *SeguridadService) enriquecerAcceso(a models.AccesoCliente) AccesoConNombre {
	nombre := ""
	if c, ok := s.clientes.BuscarClientePorID(int(a.ClienteID)); ok {
		nombre = c.Nombre
	}
	return AccesoConNombre{
		AccesoCliente: a,
		NombreCliente: nombre,
		PagoAlDia:     s.pagos.ClienteTienePagoEntrada(int(a.ClienteID)),
	}
}

// CrearAcceso es la regla de negocio central de la vinculación Seguridad↔Clientes:
// antes de autorizar el ingreso, verifica que el cliente exista y tenga un pago
// de entrada registrado. Si no lo tiene, el acceso se guarda como NO autorizado.
func (s *SeguridadService) CrearAcceso(clienteID uint) (AccesoConNombre, error) {
	if clienteID == 0 {
		return AccesoConNombre{}, ErrCampoObligatorio
	}
	cliente, ok := s.clientes.BuscarClientePorID(clienteID)
	if !ok {
		return AccesoConNombre{}, ErrClienteInvalido
	}

	acc := models.AccesoCliente{ClienteID: uint(clienteID)}
	if s.pagos.ClienteTienePagoEntrada(clienteID) {
		acc.Autorizado = true
		acc.Motivo = ""
	} else {
		acc.Autorizado = false
		acc.Motivo = "Sin pago de entrada registrado"
	}

	creado := s.repo.CrearAcceso(acc)
	return AccesoConNombre{
		AccesoCliente: creado,
		NombreCliente: cliente.Nombre,
		PagoAlDia:     acc.Autorizado,
	}, nil
}

func (s *SeguridadService) BorrarAcceso(id uint) error {
	if !s.repo.BorrarAcceso(id) {
		return ErrNoEncontrado
	}
	return nil
}
