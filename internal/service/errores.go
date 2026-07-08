package service

import "errors"

var (
	ErrNombreVacio           = errors.New("nombre es requerido")
	ErrCampoObligatorio      = errors.New("falta un campo obligatorio")
	ErrNoEncontrado          = errors.New("recurso no encontrado")
	ErrCedulaEnUso           = errors.New("cédula ya registrada")
	ErrMontoInvalido         = errors.New("el monto debe ser mayor a cero")
	ErrGuardavidaInvalido    = errors.New("guardavida_id no existe")
	ErrClienteInvalido       = errors.New("cliente_id no existe")
	ErrClienteSinMembresia   = errors.New("el cliente no tiene membresía")
	ErrClienteConMembresia   = errors.New("el cliente ya tiene membresía")
	ErrClienteSinAcceso      = errors.New("el cliente no tiene membresía ni pago registrado")
	ErrEquipoInvalido        = errors.New("equipo_id no existe")
	ErrConceptoPagoInvalido  = errors.New("el concepto debe ser medio_dia o dia")
	ErrDuracionInvalida      = errors.New("la duración debe ser medio día o un día")
	ErrCredencialesInvalidas = errors.New("email o contraseña incorrectos")
	ErrEmailEnUso            = errors.New("ese email ya está registrado")
	ErrTelefonoEnUso         = errors.New("ese teléfono ya está registrado")
)
