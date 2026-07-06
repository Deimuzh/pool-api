// Command piscina-api arranca el servidor HTTP de la Piscina Comunitaria
// "Los Ceibos". Arquitectura en 4 capas: storage (Almacen) → service →
// handlers (Server) → main (ensamblaje + rutas).
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glebarez/sqlite" // driver GORM (pure-Go)
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"pool-api/internal/config"
	"pool-api/internal/handlers"
	"pool-api/internal/httpserver"
	"pool-api/internal/middleware"
	"pool-api/internal/models"
	"pool-api/internal/service"
	"pool-api/internal/storage"
)

func main() {
	cfg := config.Cargar()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

// run concentra el arranque real de la aplicación: abre la base de datos,
// migra modelos, construye services/handlers, registra rutas y levanta el servidor.
func run(cfg config.Config) error {
	// 1. Abrir la base de datos y migrar todos los modelos.
	db, err := abrirDB(cfg)
	if err != nil {
		return fmt.Errorf("no se pudo abrir la base de datos: %w", err)
	}

	if err := db.AutoMigrate(
		&models.Guardavida{}, &models.Incidente{}, &models.AccesoCliente{},
		&models.Equipo{}, &models.RegistroMantenimiento{}, &models.ProductoQuimico{},
		&models.Cliente{}, &models.Reserva{}, &models.Pago{},
		&models.Usuario{},
	); err != nil {
		return fmt.Errorf("falló AutoMigrate: %w", err)
	}

	almacen := storage.NuevoAlmacenSQLite(db)
	almacen.SembrarSiVacio() // crea admin@piscina.com / admin123 si no existe ningún usuario

	// 2. Construir los services.
	seguridadSvc := service.NewSeguridadService(almacen, almacen, almacen)
	mantenimientoSvc := service.NewMantenimientoService(almacen)
	clientesSvc := service.NewClientesService(almacen)
	// AuthService recibe secreto y duración desde config para no hardcodear JWT.
	authSvc := service.NewAuthService(
		almacen,
		service.WithSecreto(cfg.JWTSecreto),
		service.WithDuracionToken(cfg.JWTDuracion),
	)

	// 3. Server con los services inyectados.
	servidor := handlers.NewServer(seguridadSvc, mantenimientoSvc, clientesSvc, authSvc)

	// 4. Router + middleware globales.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

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
	r.Get("/styles.css", func(w http.ResponseWriter, _ *http.Request) {
		css, err := os.ReadFile("./web/styles.css")
		if err != nil {
			http.Error(w, "No se encontró styles.css en ./web/", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(css)
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
				r.Use(middleware.SoloRoles("admin"))
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

	log.Printf("Motor de base de datos: %s", cfg.DBDriver)
	log.Printf("Servidor escuchando en http://localhost%s", cfg.Puerto)
	log.Println("Login: admin@piscina.com / admin123")

	// Servidor HTTP con timeouts para evitar conexiones colgadas.
	srv := httpserver.Nuevo(r, cfg.Puerto, cfg.ReadTimeout, cfg.WriteTimeout)
	errServidor := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()

	// Espera Ctrl+C o SIGTERM para apagar sin cortar requests en curso.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Señal de apagado recibida, cerrando servidor...")
	}

	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	if err := srv.Shutdown(ctxApagado); err != nil {
		return err
	}
	log.Println("Servidor detenido correctamente.")
	return nil
}

func abrirDB(cfg config.Config) (*gorm.DB, error) {
	switch cfg.DBDriver {
	case "postgres":
		if cfg.DBDSN == "" {
			return nil, fmt.Errorf("DB_DSN es obligatorio cuando DB_DRIVER=postgres")
		}
		return gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.RutaDB), &gorm.Config{})
	default:
		return nil, fmt.Errorf("DB_DRIVER inválido: %s", cfg.DBDriver)
	}
}
