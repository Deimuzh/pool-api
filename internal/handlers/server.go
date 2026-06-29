package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"pool-api/internal/models"
	"pool-api/internal/service"
)

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

func parseUintParam(r *http.Request, name string) (uint, error) {
	idStr := chi.URLParam(r, name)
	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(idUint), nil
}

func decodeJSONBody(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	RespondJSON(w, status, data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	RespondError(w, status, message)
}

func (s *Server) ListarGuardavidas(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Seguridad.ListarGuardavidas())
}

func (s *Server) CrearGuardavida(w http.ResponseWriter, r *http.Request) {
	var g models.Guardavida
	if err := decodeJSONBody(r, &g); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Seguridad.CrearGuardavida(g)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerGuardavida(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	g, ok := s.Seguridad.ObtenerGuardavida(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Guardavida no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (s *Server) ActualizarGuardavida(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var g models.Guardavida
	if err := decodeJSONBody(r, &g); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Seguridad.ActualizarGuardavida(id, g)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarGuardavida(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Seguridad.BorrarGuardavida(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "guardavida eliminado"})
}

func (s *Server) ListarIncidentes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Seguridad.ListarIncidentes())
}

func (s *Server) CrearIncidente(w http.ResponseWriter, r *http.Request) {
	var inc models.Incidente
	if err := decodeJSONBody(r, &inc); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Seguridad.CrearIncidente(inc)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerIncidente(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	inc, ok := s.Seguridad.ObtenerIncidente(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Incidente no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, inc)
}

func (s *Server) ActualizarIncidente(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var inc models.Incidente
	if err := decodeJSONBody(r, &inc); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Seguridad.ActualizarIncidente(id, inc)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarIncidente(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Seguridad.BorrarIncidente(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "incidente eliminado"})
}

func (s *Server) ListarAccesos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Seguridad.ListarAccesos())
}

func (s *Server) CrearAcceso(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ClienteID uint `json:"cliente_id"`
	}
	if err := decodeJSONBody(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Seguridad.CrearAcceso(body.ClienteID)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) BorrarAcceso(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Seguridad.BorrarAcceso(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "acceso eliminado"})
}

func (s *Server) ListarEquipos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Mantenimiento.ListarEquipos())
}

func (s *Server) CrearEquipo(w http.ResponseWriter, r *http.Request) {
	var eq models.Equipo
	if err := decodeJSONBody(r, &eq); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Mantenimiento.CrearEquipo(eq)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	eq, ok := s.Mantenimiento.ObtenerEquipo(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Equipo no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, eq)
}

func (s *Server) ActualizarEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var eq models.Equipo
	if err := decodeJSONBody(r, &eq); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarEquipo(id, eq)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Mantenimiento.BorrarEquipo(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "equipo eliminado"})
}

func (s *Server) ListarRegistrosMantenimiento(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Mantenimiento.ListarRegistros())
}

func (s *Server) CrearRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	var rm models.RegistroMantenimiento
	if err := decodeJSONBody(r, &rm); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Mantenimiento.CrearRegistro(rm)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	rm, ok := s.Mantenimiento.ObtenerRegistro(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Registro no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, rm)
}

func (s *Server) ActualizarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var rm models.RegistroMantenimiento
	if err := decodeJSONBody(r, &rm); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarRegistro(id, rm)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarRegistroMantenimiento(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Mantenimiento.BorrarRegistro(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "registro eliminado"})
}

func (s *Server) ListarQuimicos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Mantenimiento.ListarQuimicos())
}

func (s *Server) CrearQuimico(w http.ResponseWriter, r *http.Request) {
	var pq models.ProductoQuimico
	if err := decodeJSONBody(r, &pq); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Mantenimiento.CrearQuimico(pq)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerQuimico(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	pq, ok := s.Mantenimiento.ObtenerQuimico(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Producto no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, pq)
}

func (s *Server) ActualizarQuimico(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var pq models.ProductoQuimico
	if err := decodeJSONBody(r, &pq); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Mantenimiento.ActualizarQuimico(id, pq)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarQuimico(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Mantenimiento.BorrarQuimico(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "producto eliminado"})
}

func (s *Server) ListarClientes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Clientes.ListarClientes())
}

func (s *Server) CrearCliente(w http.ResponseWriter, r *http.Request) {
	var c models.Cliente
	if err := decodeJSONBody(r, &c); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Clientes.CrearCliente(c)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerCliente(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	c, ok := s.Clientes.ObtenerCliente(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Cliente no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var c models.Cliente
	if err := decodeJSONBody(r, &c); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Clientes.ActualizarCliente(id, c)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarCliente(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Clientes.BorrarCliente(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "cliente eliminado"})
}

func (s *Server) ListarReservas(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Clientes.ListarReservas())
}

func (s *Server) CrearReserva(w http.ResponseWriter, r *http.Request) {
	var rv models.Reserva
	if err := decodeJSONBody(r, &rv); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Clientes.CrearReserva(rv)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerReserva(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	rv, ok := s.Clientes.ObtenerReserva(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Reserva no encontrada")
		return
	}
	writeJSON(w, http.StatusOK, rv)
}

func (s *Server) ActualizarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var rv models.Reserva
	if err := decodeJSONBody(r, &rv); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Clientes.ActualizarReserva(id, rv)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Clientes.BorrarReserva(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "reserva eliminada"})
}

func (s *Server) ListarPagos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Clientes.ListarPagos())
}

func (s *Server) CrearPago(w http.ResponseWriter, r *http.Request) {
	var p models.Pago
	if err := decodeJSONBody(r, &p); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creado, err := s.Clientes.CrearPago(p)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (s *Server) ObtenerPago(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	p, ok := s.Clientes.ObtenerPago(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Pago no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) ActualizarPago(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var p models.Pago
	if err := decodeJSONBody(r, &p); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizado, err := s.Clientes.ActualizarPago(id, p)
	if err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarPago(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := s.Clientes.BorrarPago(id); err != nil {
		writeError(w, statusDeError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "pago eliminado"})
}
