package handlers

import "net/http"

// Wrappers: adaptadores que exponen methods en *Server delegando
// a las funciones existentes (para mantener compatibilidad con main.go).

// Clientes
func (s *Server) ListarClientes(w http.ResponseWriter, r *http.Request)    { ListarClientes(w, r) }
func (s *Server) CrearCliente(w http.ResponseWriter, r *http.Request)      { CrearCliente(w, r) }
func (s *Server) ObtenerCliente(w http.ResponseWriter, r *http.Request)    { ObtenerCliente(w, r) }
func (s *Server) ActualizarCliente(w http.ResponseWriter, r *http.Request) { ActualizarCliente(w, r) }
func (s *Server) BorrarCliente(w http.ResponseWriter, r *http.Request)     { EliminarCliente(w, r) }

// Reservas
func (s *Server) ListarReservas(w http.ResponseWriter, r *http.Request)    { ListarReservas(w, r) }
func (s *Server) CrearReserva(w http.ResponseWriter, r *http.Request)      { CrearReserva(w, r) }
func (s *Server) ObtenerReserva(w http.ResponseWriter, r *http.Request)    { ObtenerReserva(w, r) }
func (s *Server) ActualizarReserva(w http.ResponseWriter, r *http.Request) { ActualizarReserva(w, r) }
func (s *Server) BorrarReserva(w http.ResponseWriter, r *http.Request)     { EliminarReserva(w, r) }

// Pagos
func (s *Server) ListarPagos(w http.ResponseWriter, r *http.Request)    { ListarPagos(w, r) }
func (s *Server) CrearPago(w http.ResponseWriter, r *http.Request)      { CrearPago(w, r) }
func (s *Server) ObtenerPago(w http.ResponseWriter, r *http.Request)    { ObtenerPago(w, r) }
func (s *Server) ActualizarPago(w http.ResponseWriter, r *http.Request) { ActualizarPago(w, r) }
func (s *Server) BorrarPago(w http.ResponseWriter, r *http.Request)     { EliminarPago(w, r) }

// Mantenimiento - Equipos
func (s *Server) ListarEquipos(w http.ResponseWriter, r *http.Request)    { ListarEquipos(w, r) }
func (s *Server) CrearEquipo(w http.ResponseWriter, r *http.Request)      { CrearEquipo(w, r) }
func (s *Server) ObtenerEquipo(w http.ResponseWriter, r *http.Request)    { ObtenerEquipo(w, r) }
func (s *Server) ActualizarEquipo(w http.ResponseWriter, r *http.Request) { ActualizarEquipo(w, r) }
func (s *Server) BorrarEquipo(w http.ResponseWriter, r *http.Request)     { EliminarEquipo(w, r) }

// Mantenimiento - Registros
func (s *Server) ListarRegistrosMantenimiento(w http.ResponseWriter, r *http.Request) {
	ListarRegistrosMantenimiento(w, r)
}
func (s *Server) CrearRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	CrearRegistroMantenimiento(w, r)
}
func (s *Server) ObtenerRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	ObtenerRegistroMantenimiento(w, r)
}
func (s *Server) ActualizarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	ActualizarRegistroMantenimiento(w, r)
}
func (s *Server) BorrarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	EliminarRegistroMantenimiento(w, r)
}

// Mantenimiento - Químicos (mapear nombres)
func (s *Server) ListarQuimicos(w http.ResponseWriter, r *http.Request) {
	ListarProductosQuimicos(w, r)
}
func (s *Server) CrearQuimico(w http.ResponseWriter, r *http.Request)   { CrearProductoQuimico(w, r) }
func (s *Server) ObtenerQuimico(w http.ResponseWriter, r *http.Request) { ObtenerProductoQuimico(w, r) }
func (s *Server) ActualizarQuimico(w http.ResponseWriter, r *http.Request) {
	ActualizarProductoQuimico(w, r)
}
func (s *Server) BorrarQuimico(w http.ResponseWriter, r *http.Request) { EliminarProductoQuimico(w, r) }
