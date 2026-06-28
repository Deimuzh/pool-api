package storage

import "pool-api/internal/models"

// Almacen define QUÉ sabe hacer el almacén de la piscina, sin decir CÓMO.
//
// AlmacenSQLite (GORM) implementa esta interfaz. Los services dependen de
// las sub-interfaces (SeguridadRepository, MantenimientoRepository,
// ClienteRepository, UsuarioRepository), no de una implementación concreta:
// así se puede cambiar el backend de almacenamiento sin tocar un solo
// handler ni service.

// ─── SEGURIDAD ────────────────────────────────────────────────────────────────

type GuardavidaRepository interface {
	ListarGuardavidas() []models.Guardavida
	BuscarGuardavidaPorID(id uint) (models.Guardavida, bool)
	CrearGuardavida(g models.Guardavida) models.Guardavida
	ActualizarGuardavida(id uint, datos models.Guardavida) (models.Guardavida, bool)
	BorrarGuardavida(id uint) bool
}

type IncidenteRepository interface {
	ListarIncidentes() []models.Incidente
	BuscarIncidentePorID(id uint) (models.Incidente, bool)
	CrearIncidente(i models.Incidente) models.Incidente
	ActualizarIncidente(id uint, datos models.Incidente) (models.Incidente, bool)
	BorrarIncidente(id uint) bool
}

type AccesoRepository interface {
	ListarAccesos() []models.AccesoCliente
	BuscarAccesoPorID(id uint) (models.AccesoCliente, bool)
	CrearAcceso(a models.AccesoCliente) models.AccesoCliente
	ActualizarAcceso(id uint, datos models.AccesoCliente) (models.AccesoCliente, bool)
	BorrarAcceso(id uint) bool
}

type SeguridadRepository interface {
	GuardavidaRepository
	IncidenteRepository
	AccesoRepository
}

// ─── MANTENIMIENTO ────────────────────────────────────────────────────────────

type EquipoRepository interface {
	ListarEquipos() []models.Equipo
	BuscarEquipoPorID(id uint) (models.Equipo, bool)
	CrearEquipo(e models.Equipo) models.Equipo
	ActualizarEquipo(id uint, datos models.Equipo) (models.Equipo, bool)
	BorrarEquipo(id uint) bool
}

type RegistroMantenimientoRepository interface {
	ListarRegistros() []models.RegistroMantenimiento
	BuscarRegistroPorID(id uint) (models.RegistroMantenimiento, bool)
	CrearRegistro(rm models.RegistroMantenimiento) models.RegistroMantenimiento
	ActualizarRegistro(id uint, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool)
	BorrarRegistro(id uint) bool
}

type ProductoQuimicoRepository interface {
	ListarQuimicos() []models.ProductoQuimico
	BuscarQuimicoPorID(id uint) (models.ProductoQuimico, bool)
	CrearQuimico(q models.ProductoQuimico) models.ProductoQuimico
	ActualizarQuimico(id uint, datos models.ProductoQuimico) (models.ProductoQuimico, bool)
	BorrarQuimico(id uint) bool
}

type MantenimientoRepository interface {
	EquipoRepository
	RegistroMantenimientoRepository
	ProductoQuimicoRepository
}

// ─── CLIENTES ─────────────────────────────────────────────────────────────────

type ClienteRepository interface {
	ListarClientes() []models.Cliente
	BuscarClientePorID(id uint) (models.Cliente, bool)
	CrearCliente(c models.Cliente) models.Cliente
	ActualizarCliente(id uint, datos models.Cliente) (models.Cliente, bool)
	BorrarCliente(id uint) bool
}

type ReservaRepository interface {
	ListarReservas() []models.Reserva
	BuscarReservaPorID(id uint) (models.Reserva, bool)
	CrearReserva(rv models.Reserva) models.Reserva
	ActualizarReserva(id uint, datos models.Reserva) (models.Reserva, bool)
	BorrarReserva(id uint) bool
}

type PagoRepository interface {
	ListarPagos() []models.Pago
	BuscarPagoPorID(id uint) (models.Pago, bool)
	CrearPago(p models.Pago) models.Pago
	ActualizarPago(id uint, datos models.Pago) (models.Pago, bool)
	BorrarPago(id uint) bool
	// ClienteTienePagoEntrada permite a otros módulos (Seguridad) verificar
	// si un cliente ya pagó su entrada, sin acoplarse a la tabla Pago directamente.
	ClienteTienePagoEntrada(clienteID uint) bool
}

type ClientesModulo interface {
	ClienteRepository
	ReservaRepository
	PagoRepository
}

// ─── USUARIOS (autenticación) ─────────────────────────────────────────────────

type UsuarioRepository interface {
	ListarUsuarios() []models.Usuario
	BuscarUsuarioPorID(id uint) (models.Usuario, bool)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	ActualizarUsuario(id uint, datos models.Usuario) (models.Usuario, bool)
	BorrarUsuario(id uint) bool
}

// Almacen agrupa los 4 módulos. AlmacenSQLite debe cumplir esta interfaz completa.
type Almacen interface {
	SeguridadRepository
	MantenimientoRepository
	ClientesModulo
	UsuarioRepository
}
