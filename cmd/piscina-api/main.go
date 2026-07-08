// Command piscina-api arranca el servidor HTTP de la Piscina Comunitaria
// "Los Ceibos". Arquitectura en 4 capas: storage (Almacen) → service →
// handlers (Server) → main (ensamblaje + rutas).
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/glebarez/sqlite" // driver GORM (pure-Go)
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"pool-api/internal/handlers"
	"pool-api/internal/middleware"
	"pool-api/internal/models"
	"pool-api/internal/service"
	"pool-api/internal/storage"
)

func main() {
	// 1. Abrir la base de datos y migrar todos los modelos.
	db, err := gorm.Open(sqlite.Open("piscina.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}
	if err := db.AutoMigrate(
		&models.Guardavida{}, &models.Incidente{}, &models.AccesoCliente{},
		&models.Equipo{}, &models.RegistroMantenimiento{}, &models.ProductoQuimico{},
		&models.Cliente{}, &models.Reserva{}, &models.Pago{},
		&models.Usuario{},
	); err != nil {
		log.Fatal("falló AutoMigrate: ", err)
	}

	almacen := storage.NuevoAlmacenSQLite(db)
	almacen.SembrarSiVacio() // crea admin@piscina.com / admin123 si no existe ningún usuario

	// 2. Construir los services.
	seguridadSvc := service.NewSeguridadService(almacen, almacen, almacen)
	mantenimientoSvc := service.NewMantenimientoService(almacen)
	clientesSvc := service.NewClientesService(almacen)
	authSvc := service.NewAuthService(almacen)

	// 3. Server con los services inyectados.
	servidor := handlers.NewServer(seguridadSvc, mantenimientoSvc, clientesSvc, authSvc)

	// 4. Router + middleware globales.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// 5. Sirve el frontend (index.html) en la raíz. Sin proteger: el HTML
	//    necesita cargar sin token para poder mostrar la pantalla de login.
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		html, err := os.ReadFile("./web/index.html")
		if err != nil {
			http.Error(w, "No se encontró index.html en ./web/", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	})

	// 6. Rutas versionadas /api/v1/.
	r.Route("/api/v1", func(r chi.Router) {

		// ── PÚBLICA: login ───────────────────────────────────────────────────
		// Esta es la única ruta de /api/v1 que NO exige token, porque es
		// justamente la que entrega el token.
		r.Post("/login", servidor.Login)

		// ── PROTEGIDAS: todo lo demás exige "Authorization: Bearer <token>" ──
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			// Usuarios (gestión de cuentas de administrador)
			r.Route("/usuarios", func(r chi.Router) {
				r.Get("/", servidor.ListarUsuarios)
				r.Post("/", servidor.CrearUsuario)
				r.Get("/{id}", servidor.ObtenerUsuario)
				r.Put("/{id}", servidor.ActualizarUsuario)
				r.Delete("/{id}", servidor.BorrarUsuario)
			})

			// Seguridad
			r.Route("/guardavidas", func(r chi.Router) {
				r.Get("/", servidor.ListarGuardavidas)
				r.Post("/", servidor.CrearGuardavida)
				r.Get("/{id}", servidor.ObtenerGuardavida)
				r.Put("/{id}", servidor.ActualizarGuardavida)
				r.Delete("/{id}", servidor.BorrarGuardavida)
			})
			r.Route("/incidentes", func(r chi.Router) {
				r.Get("/", servidor.ListarIncidentes)
				r.Post("/", servidor.CrearIncidente)
				r.Get("/{id}", servidor.ObtenerIncidente)
				r.Put("/{id}", servidor.ActualizarIncidente)
				r.Delete("/{id}", servidor.BorrarIncidente)
			})
			r.Route("/accesos", func(r chi.Router) {
				r.Get("/", servidor.ListarAccesos)
				r.Post("/", servidor.CrearAcceso)
				r.Delete("/{id}", servidor.BorrarAcceso)
			})

			// Mantenimiento
			r.Route("/equipos", func(r chi.Router) {
				r.Get("/", servidor.ListarEquipos)
				r.Post("/", servidor.CrearEquipo)
				r.Get("/{id}", servidor.ObtenerEquipo)
				r.Put("/{id}", servidor.ActualizarEquipo)
				r.Delete("/{id}", servidor.BorrarEquipo)
			})
			r.Route("/mantenimientos", func(r chi.Router) {
				r.Get("/", servidor.ListarRegistrosMantenimiento)
				r.Post("/", servidor.CrearRegistroMantenimiento)
				r.Get("/{id}", servidor.ObtenerRegistroMantenimiento)
				r.Put("/{id}", servidor.ActualizarRegistroMantenimiento)
				r.Delete("/{id}", servidor.BorrarRegistroMantenimiento)
			})
			r.Route("/quimicos", func(r chi.Router) {
				r.Get("/", servidor.ListarQuimicos)
				r.Post("/", servidor.CrearQuimico)
				r.Get("/{id}", servidor.ObtenerQuimico)
				r.Put("/{id}", servidor.ActualizarQuimico)
				r.Delete("/{id}", servidor.BorrarQuimico)
			})

			// Clientes
			r.Route("/clientes", func(r chi.Router) {
				r.Get("/", servidor.ListarClientes)
				r.Post("/", servidor.CrearCliente)
				r.Get("/{id}", servidor.ObtenerCliente)
				r.Put("/{id}", servidor.ActualizarCliente)
				r.Delete("/{id}", servidor.BorrarCliente)
			})
			r.Route("/reservas", func(r chi.Router) {
				r.Get("/", servidor.ListarReservas)
				r.Post("/", servidor.CrearReserva)
				r.Get("/{id}", servidor.ObtenerReserva)
				r.Put("/{id}", servidor.ActualizarReserva)
				r.Delete("/{id}", servidor.BorrarReserva)
			})
			r.Route("/pagos", func(r chi.Router) {
				r.Get("/", servidor.ListarPagos)
				r.Post("/", servidor.CrearPago)
				r.Get("/{id}", servidor.ObtenerPago)
				r.Put("/{id}", servidor.ActualizarPago)
				r.Delete("/{id}", servidor.BorrarPago)
			})
		})
	})

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Println("Login: admin@piscina.com / admin123")
	log.Fatal(http.ListenAndServe(":8080", r))
}
