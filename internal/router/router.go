package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pool-api/internal/handlers"
)

func NuevoRouter(server *handlers.Server) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares globales
	r.Use(middleware.Logger)    // muestra cada request en la terminal
	r.Use(middleware.Recoverer) // evita que un panic baje el servidor

	// Ruta base de salud
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Piscina Comunitaria - OK"))
	})

	// ─── /api/v1 ────────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {

		// ── SEGURIDAD ────────────────────────────────────────────────────────
		r.Route("/guardavidas", func(r chi.Router) {
			r.Post("/", server.CrearGuardavida)
			r.Get("/", server.ListarGuardavidas)
			r.Get("/{id}", server.ObtenerGuardavida)
			r.Patch("/{id}", server.ActualizarGuardavida)
			r.Delete("/{id}", server.BorrarGuardavida)
		})

		r.Route("/incidentes", func(r chi.Router) {
			r.Post("/", server.CrearIncidente)
			r.Get("/", server.ListarIncidentes)
			r.Get("/{id}", server.ObtenerIncidente)
			r.Patch("/{id}", server.ActualizarIncidente)
			r.Delete("/{id}", server.BorrarIncidente)
		})

		r.Route("/accesos", func(r chi.Router) {
			r.Post("/", server.CrearAcceso)
			r.Get("/", server.ListarAccesos)
			r.Delete("/{id}", server.BorrarAcceso)
		})

		// ── MANTENIMIENTO ────────────────────────────────────────────────────
		r.Route("/equipos", func(r chi.Router) {
			r.Post("/", server.CrearEquipo)
			r.Get("/", server.ListarEquipos)
			r.Get("/{id}", server.ObtenerEquipo)
			r.Patch("/{id}", server.ActualizarEquipo)
			r.Delete("/{id}", server.BorrarEquipo)
		})

		r.Route("/mantenimientos", func(r chi.Router) {
			r.Post("/", server.CrearRegistroMantenimiento)
			r.Get("/", server.ListarRegistrosMantenimiento)
			r.Get("/{id}", server.ObtenerRegistroMantenimiento)
			r.Patch("/{id}", server.ActualizarRegistroMantenimiento)
			r.Delete("/{id}", server.BorrarRegistroMantenimiento)
		})

		r.Route("/quimicos", func(r chi.Router) {
			r.Post("/", server.CrearQuimico)
			r.Get("/", server.ListarQuimicos)
			r.Get("/{id}", server.ObtenerQuimico)
			r.Patch("/{id}", server.ActualizarQuimico)
			r.Delete("/{id}", server.BorrarQuimico)
		})

		// ── CLIENTES ─────────────────────────────────────────────────────────
		r.Route("/clientes", func(r chi.Router) {
			r.Post("/", server.CrearCliente)
			r.Get("/", server.ListarClientes)
			r.Get("/{id}", server.ObtenerCliente)
			r.Patch("/{id}", server.ActualizarCliente)
			r.Delete("/{id}", server.BorrarCliente)
		})

		r.Route("/reservas", func(r chi.Router) {
			r.Post("/", server.CrearReserva)
			r.Get("/", server.ListarReservas)
			r.Get("/{id}", server.ObtenerReserva)
			r.Patch("/{id}", server.ActualizarReserva)
			r.Delete("/{id}", server.BorrarReserva)
		})

		r.Route("/pagos", func(r chi.Router) {
			r.Post("/", server.CrearPago)
			r.Get("/", server.ListarPagos)
			r.Get("/{id}", server.ObtenerPago)
			r.Patch("/{id}", server.ActualizarPago)
			r.Delete("/{id}", server.BorrarPago)
		})
	})

	return r
}
