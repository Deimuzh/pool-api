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
	ErrEquipoInvalido        = errors.New("equipo_id no existe")
	ErrCredencialesInvalidas = errors.New("email o contraseña incorrectos")
	ErrEmailEnUso            = errors.New("ese email ya está registrado")
)
