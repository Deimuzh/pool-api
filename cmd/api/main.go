package main

import (
	"fmt"
	"log"
	"net/http"

	"pool-api/internal/router"
	"pool-api/internal/storage"
)

func main() {
	// Inicializar la base de datos SQLite
	storage.IniciarDB()

	// Configurar el router con todas las rutas
	r := router.NuevoRouter()

	puerto := ":8080"
	fmt.Println("Servidor corriendo en http://localhost" + puerto)

	if err := http.ListenAndServe(puerto, r); err != nil {
		log.Fatal("Error al levantar el servidor: ", err)
	}
}
