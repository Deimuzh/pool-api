package main

import (
	"fmt"
	"log"
	"net/http"

	"pool-api/internal/handlers"
	"pool-api/internal/router"
	"pool-api/internal/service"
	"pool-api/internal/storage"
)

func main() {
	// Inicializar la base de datos SQLite
	storage.IniciarDB()

	// Crear el almacén SQLite
	almacen := storage.NuevoAlmacenSQLite(storage.DB)

	// Crear los services
	seguridadService := service.NewSeguridadService(almacen, almacen, almacen)
	mantenimientoService := service.NewMantenimientoService(almacen)
	clientesService := service.NewClientesService(almacen)
	authService := service.NewAuthService(almacen)

	// Crear el servidor de handlers
	server := handlers.NewServer(seguridadService, mantenimientoService, clientesService, authService)

	// Configurar el router con el server
	r := router.NuevoRouter(server)

	puerto := ":8080"
	fmt.Println("Servidor corriendo en http://localhost" + puerto)

	if err := http.ListenAndServe(puerto, r); err != nil {
		log.Fatal("Error al levantar el servidor: ", err)
	}
}
