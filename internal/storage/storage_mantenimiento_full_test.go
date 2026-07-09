package storage

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

func abrirDBFullMantenimiento(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir la base de datos en memoria: %v", err)
	}
	if err := db.AutoMigrate(&models.Equipo{}, &models.RegistroMantenimiento{}, &models.ProductoQuimico{}); err != nil {
		t.Fatalf("fallo AutoMigrate: %v", err)
	}
	return db
}

func TestAlmacenSQLite_EquipoCRUD(t *testing.T) {
	db := abrirDBFullMantenimiento(t)
	almacen := NuevoAlmacenSQLite(db)

	creado, err := almacen.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "operativo"})
	if err != nil {
		t.Fatalf("CrearEquipo fallo: %v", err)
	}

	actualizado, ok := almacen.ActualizarEquipo(creado.ID, models.Equipo{Nombre: "Bomba V2", Tipo: "bomba", Estado: "en reparacion"})
	if !ok {
		t.Fatal("ActualizarEquipo devolvio false")
	}
	if actualizado.Estado != "en reparacion" {
		t.Errorf("estado inesperado: %s", actualizado.Estado)
	}

	if !almacen.BorrarEquipo(creado.ID) {
		t.Fatal("BorrarEquipo devolvio false")
	}
	if almacen.BorrarEquipo(999) {
		t.Error("BorrarEquipo con ID inexistente devolvio true")
	}
}

func TestAlmacenSQLite_BuscarEquipoInexistente(t *testing.T) {
	db := abrirDBFullMantenimiento(t)
	almacen := NuevoAlmacenSQLite(db)

	_, ok := almacen.BuscarEquipoPorID(999)
	if ok {
		t.Error("no se esperaba encontrar equipo con ID 999")
	}
}

func TestAlmacenSQLite_RegistroMantenimientoCRUD(t *testing.T) {
	db := abrirDBFullMantenimiento(t)
	almacen := NuevoAlmacenSQLite(db)

	eq, _ := almacen.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtro"})

	nuevo, err := almacen.CrearRegistro(models.RegistroMantenimiento{
		EquipoID: eq.ID, Tipo: "preventivo", Descripcion: "Limpieza", FechaHora: time.Now(),
	})
	if err != nil {
		t.Fatalf("CrearRegistro fallo: %v", err)
	}
	if nuevo.ID == 0 {
		t.Fatal("se esperaba ID")
	}
	if nuevo.Tipo != "preventivo" {
		t.Errorf("tipo inesperado: %s", nuevo.Tipo)
	}

	encontrado, ok := almacen.BuscarRegistroPorID(nuevo.ID)
	if !ok {
		t.Fatal("no se encontro el registro recien creado")
	}
	if encontrado.EquipoID != eq.ID {
		t.Errorf("EquipoID inesperado: %d", encontrado.EquipoID)
	}

	lista := almacen.ListarRegistros()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 registro, se obtuvieron %d", len(lista))
	}

	actualizado, ok := almacen.ActualizarRegistro(nuevo.ID, models.RegistroMantenimiento{Tipo: "correctivo", Descripcion: "Reparacion"})
	if !ok {
		t.Fatal("ActualizarRegistro devolvio false")
	}
	if actualizado.Tipo != "correctivo" {
		t.Errorf("tipo inesperado: %s", actualizado.Tipo)
	}

	if !almacen.BorrarRegistro(nuevo.ID) {
		t.Fatal("BorrarRegistro devolvio false")
	}
	if almacen.BorrarRegistro(999) {
		t.Error("BorrarRegistro con ID inexistente devolvio true")
	}
}

func TestAlmacenSQLite_ProductoQuimicoCRUD(t *testing.T) {
	db := abrirDBFullMantenimiento(t)
	almacen := NuevoAlmacenSQLite(db)

	nuevo, err := almacen.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro", StockActual: 100, UnidadMedida: "kg", NivelMinimo: 10})
	if err != nil {
		t.Fatalf("CrearQuimico fallo: %v", err)
	}
	if nuevo.ID == 0 {
		t.Fatal("se esperaba ID")
	}

	encontrado, ok := almacen.BuscarQuimicoPorID(nuevo.ID)
	if !ok {
		t.Fatal("no se encontro el quimico recien creado")
	}
	if encontrado.Nombre != "Cloro" {
		t.Errorf("nombre inesperado: %s", encontrado.Nombre)
	}

	lista := almacen.ListarQuimicos()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 quimico, se obtuvieron %d", len(lista))
	}

	actualizado, ok := almacen.ActualizarQuimico(nuevo.ID, models.ProductoQuimico{Nombre: "Cloro Plus", StockActual: 80})
	if !ok {
		t.Fatal("ActualizarQuimico devolvio false")
	}
	if actualizado.StockActual != 80 {
		t.Errorf("stock inesperado: %f", actualizado.StockActual)
	}

	if !almacen.BorrarQuimico(nuevo.ID) {
		t.Fatal("BorrarQuimico devolvio false")
	}
	if almacen.BorrarQuimico(999) {
		t.Error("BorrarQuimico con ID inexistente devolvio true")
	}
}
