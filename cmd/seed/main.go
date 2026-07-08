// Command seed ejecuta migraciones y datos iniciales sin levantar el servidor HTTP.
package main

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"pool-api/internal/config"
	"pool-api/internal/models"
	"pool-api/internal/storage"
)

func main() {
	cfg := config.Cargar()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
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
		return fmt.Errorf("fallo AutoMigrate: %w", err)
	}

	almacen := storage.NuevoAlmacenSQLite(db)
	almacen.SembrarSiVacio()
	log.Println("Seeders ejecutados correctamente")
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
