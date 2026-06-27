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
	BuscarGuardavidaPorID(id int) (models.Guardavida, bool)
	CrearGuardavida(g models.Guardavida) models.Guardavida
	ActualizarGuardavida(id int, datos models.Guardavida) (models.Guardavida, bool)
	BorrarGuardavida(id int) bool
}

type IncidenteRepository interface {
	ListarIncidentes() []models.Incidente
	BuscarIncidentePorID(id int) (models.Incidente, bool)
	CrearIncidente(i models.Incidente) models.Incidente
	ActualizarIncidente(id int, datos models.Incidente) (models.Incidente, bool)
	BorrarIncidente(id int) bool
}

type AccesoRepository interface {
	ListarAccesos() []models.AccesoCliente
	BuscarAccesoPorID(id int) (models.AccesoCliente, bool)
	CrearAcceso(a models.AccesoCliente) models.AccesoCliente
	ActualizarAcceso(id int, datos models.AccesoCliente) (models.AccesoCliente, bool)
	BorrarAcceso(id int) bool
}

type SeguridadRepository interface {
	GuardavidaRepository
	IncidenteRepository
	AccesoRepository
}

// ─── MANTENIMIENTO ────────────────────────────────────────────────────────────

type EquipoRepository interface {
	ListarEquipos() []models.Equipo
	BuscarEquipoPorID(id int) (models.Equipo, bool)
	CrearEquipo(e models.Equipo) models.Equipo
	ActualizarEquipo(id int, datos models.Equipo) (models.Equipo, bool)
	BorrarEquipo(id int) bool
}

type RegistroMantenimientoRepository interface {
	ListarRegistros() []models.RegistroMantenimiento
	BuscarRegistroPorID(id int) (models.RegistroMantenimiento, bool)
	CrearRegistro(rm models.RegistroMantenimiento) models.RegistroMantenimiento
	ActualizarRegistro(id int, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool)
	BorrarRegistro(id int) bool
}

type ProductoQuimicoRepository interface {
	ListarQuimicos() []models.ProductoQuimico
	BuscarQuimicoPorID(id int) (models.ProductoQuimico, bool)
	CrearQuimico(q models.ProductoQuimico) models.ProductoQuimico
	ActualizarQuimico(id int, datos models.ProductoQuimico) (models.ProductoQuimico, bool)
	BorrarQuimico(id int) bool
}

type MantenimientoRepository interface {
	EquipoRepository
	RegistroMantenimientoRepository
	ProductoQuimicoRepository
}

// ─── CLIENTES ─────────────────────────────────────────────────────────────────

type ClienteRepository interface {
	ListarClientes() []models.Cliente
	BuscarClientePorID(id int) (models.Cliente, bool)
	CrearCliente(c models.Cliente) models.Cliente
	ActualizarCliente(id int, datos models.Cliente) (models.Cliente, bool)
	BorrarCliente(id int) bool
}

type ReservaRepository interface {
	ListarReservas() []models.Reserva
	BuscarReservaPorID(id int) (models.Reserva, bool)
	CrearReserva(rv models.Reserva) models.Reserva
	ActualizarReserva(id int, datos models.Reserva) (models.Reserva, bool)
	BorrarReserva(id int) bool
}

type PagoRepository interface {
	ListarPagos() []models.Pago
	BuscarPagoPorID(id int) (models.Pago, bool)
	CrearPago(p models.Pago) models.Pago
	ActualizarPago(id int, datos models.Pago) (models.Pago, bool)
	BorrarPago(id int) bool
	// ClienteTienePagoEntrada permite a otros módulos (Seguridad) verificar
	// si un cliente ya pagó su entrada, sin acoplarse a la tabla Pago directamente.
	ClienteTienePagoEntrada(clienteID int) bool
}

type ClientesModulo interface {
	ClienteRepository
	ReservaRepository
	PagoRepository
}

// ─── USUARIOS (autenticación) ─────────────────────────────────────────────────

type UsuarioRepository interface {
	ListarUsuarios() []models.Usuario
	BuscarUsuarioPorID(id int) (models.Usuario, bool)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool)
	BorrarUsuario(id int) bool
}

// Almacen agrupa los 4 módulos. AlmacenSQLite debe cumplir esta interfaz completa.
type Almacen interface {
	SeguridadRepository
	MantenimientoRepository
	ClientesModulo
	UsuarioRepository
}
