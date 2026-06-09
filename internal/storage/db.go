package storage

import (
	"log"

	"piscina-comunitaria-api/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func IniciarDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("piscina.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar con la base de datos: ", err)
	}

	// AutoMigrate crea las tablas si no existen
	err = DB.AutoMigrate(
		&models.Guardavida{},
		&models.Incidente{},
		&models.AccesoCliente{},
		&models.Equipo{},
		&models.RegistroMantenimiento{},
		&models.ProductoQuimico{},
		&models.Cliente{},
		&models.Reserva{},
		&models.Pago{},
	)
	if err != nil {
		log.Fatal("Error en AutoMigrate: ", err)
	}

	log.Println("Base de datos lista")
}
