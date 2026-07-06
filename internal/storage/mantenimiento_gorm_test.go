package storage

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

func TestAlmacenSQLite_EquipoCrearListarYBuscar(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite en memoria: %v", err)
	}
	if err := db.AutoMigrate(&models.Equipo{}); err != nil {
		t.Fatalf("fallo AutoMigrate: %v", err)
	}

	repo := NuevoAlmacenSQLite(db)
	creado, err := repo.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion", Estado: "operativo"})
	if err != nil {
		t.Fatalf("no se esperaba error al crear equipo: %v", err)
	}
	if creado.ID == 0 {
		t.Fatal("se esperaba que GORM asigne un ID")
	}

	lista := repo.ListarEquipos()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 equipo, se obtuvieron %d", len(lista))
	}

	buscado, ok := repo.BuscarEquipoPorID(creado.ID)
	if !ok {
		t.Fatal("se esperaba encontrar el equipo creado")
	}
	if buscado.Nombre != "Filtro" {
		t.Fatalf("nombre inesperado: %s", buscado.Nombre)
	}
}
