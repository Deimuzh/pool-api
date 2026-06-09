package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"piscina-comunitaria-api/internal/handlers"
)

func NuevoRouter() *chi.Mux {
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
			r.Post("/", handlers.CrearGuardavida)
			r.Get("/", handlers.ListarGuardavidas)
			r.Get("/{id}", handlers.ObtenerGuardavida)
			r.Patch("/{id}", handlers.ActualizarGuardavida)
			r.Delete("/{id}", handlers.EliminarGuardavida)
		})

		r.Route("/incidentes", func(r chi.Router) {
			r.Post("/", handlers.CrearIncidente)
			r.Get("/", handlers.ListarIncidentes)
			r.Get("/{id}", handlers.ObtenerIncidente)
			r.Patch("/{id}", handlers.ActualizarIncidente)
			r.Delete("/{id}", handlers.EliminarIncidente)
		})

		r.Route("/accesos", func(r chi.Router) {
			r.Post("/", handlers.CrearAcceso)
			r.Get("/", handlers.ListarAccesos)
			r.Get("/{id}", handlers.ObtenerAcceso)
			r.Patch("/{id}", handlers.ActualizarAcceso)
			r.Delete("/{id}", handlers.EliminarAcceso)
		})

		// ── MANTENIMIENTO ────────────────────────────────────────────────────
		r.Route("/equipos", func(r chi.Router) {
			r.Post("/", handlers.CrearEquipo)
			r.Get("/", handlers.ListarEquipos)
			r.Get("/{id}", handlers.ObtenerEquipo)
			r.Patch("/{id}", handlers.ActualizarEquipo)
			r.Delete("/{id}", handlers.EliminarEquipo)
		})

		r.Route("/mantenimientos", func(r chi.Router) {
			r.Post("/", handlers.CrearRegistroMantenimiento)
			r.Get("/", handlers.ListarRegistrosMantenimiento)
			r.Get("/{id}", handlers.ObtenerRegistroMantenimiento)
			r.Patch("/{id}", handlers.ActualizarRegistroMantenimiento)
			r.Delete("/{id}", handlers.EliminarRegistroMantenimiento)
		})

		r.Route("/quimicos", func(r chi.Router) {
			r.Post("/", handlers.CrearProductoQuimico)
			r.Get("/", handlers.ListarProductosQuimicos)
			r.Get("/{id}", handlers.ObtenerProductoQuimico)
			r.Patch("/{id}", handlers.ActualizarProductoQuimico)
			r.Delete("/{id}", handlers.EliminarProductoQuimico)
		})

		// ── CLIENTES ─────────────────────────────────────────────────────────
		r.Route("/clientes", func(r chi.Router) {
			r.Post("/", handlers.CrearCliente)
			r.Get("/", handlers.ListarClientes)
			r.Get("/{id}", handlers.ObtenerCliente)
			r.Patch("/{id}", handlers.ActualizarCliente)
			r.Delete("/{id}", handlers.EliminarCliente)
		})

		r.Route("/reservas", func(r chi.Router) {
			r.Post("/", handlers.CrearReserva)
			r.Get("/", handlers.ListarReservas)
			r.Get("/{id}", handlers.ObtenerReserva)
			r.Patch("/{id}", handlers.ActualizarReserva)
			r.Delete("/{id}", handlers.EliminarReserva)
		})

		r.Route("/pagos", func(r chi.Router) {
			r.Post("/", handlers.CrearPago)
			r.Get("/", handlers.ListarPagos)
			r.Get("/{id}", handlers.ObtenerPago)
			r.Patch("/{id}", handlers.ActualizarPago)
			r.Delete("/{id}", handlers.EliminarPago)
		})
	})

	return r
}
