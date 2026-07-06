package handlers

import "pool-api/internal/service"

// Server agrupa las dependencias compartidas por todos los handlers.
// Cada campo es un service de un módulo; los handlers nunca tocan storage
// directamente, siempre pasan por el service correspondiente.
type Server struct {
	Seguridad     *service.SeguridadService
	Mantenimiento *service.MantenimientoService
	Clientes      *service.ClientesService
	Auth          *service.AuthService
}

// NewServer construye un Server listo para usar, con los services inyectados.
func NewServer(
	seguridad *service.SeguridadService,
	mantenimiento *service.MantenimientoService,
	clientes *service.ClientesService,
	auth *service.AuthService,
) *Server {
	return &Server{
		Seguridad:     seguridad,
		Mantenimiento: mantenimiento,
		Clientes:      clientes,
		Auth:          auth,
	}
}
