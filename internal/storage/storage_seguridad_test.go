package storage

import (
	"testing"

	"github.com/glebarez/sqlite" // mismo driver GORM pure-Go que usa main.go
	"gorm.io/gorm"

	"pool-api/internal/models"
)

// abrirDBPrueba abre una base de datos SQLite en memoria y migra los modelos
// del módulo Seguridad. Cada test la llama de cero, así que las pruebas no
// comparten estado entre sí (cada uno arranca con una base de datos limpia).
func abrirDBPrueba(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir la base de datos en memoria: %v", err)
	}
	if err := db.AutoMigrate(&models.Guardavida{}); err != nil {
		t.Fatalf("falló AutoMigrate: %v", err)
	}
	return db
}

// TestAlmacenSQLite_CrearYListarGuardavida prueba el repositorio real (sin
// mocks ni fakes) contra una base sqlite :memory:: crear un Guardavida debe
// reflejarse al buscarlo por ID y al listarlo. Si CrearGuardavida o
// BuscarGuardavidaPorID tuvieran un bug (por ejemplo, no guardar realmente
// en la tabla, o buscar en la tabla equivocada), este test fallaría porque
// no encontraría el registro recién creado.
func TestAlmacenSQLite_CrearYListarGuardavida(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	nuevo := models.Guardavida{
		Nombre:      "Pedro Salazar",
		Turno:       "tarde",
		Certificado: "Cruz Roja Niv. 1",
		Activo:      true,
	}

	creado := almacen.CrearGuardavida(nuevo)
	if creado.ID == 0 {
		t.Fatal("se esperaba que GORM asignara un ID al crear el guardavida")
	}

	// Buscar por ID debe reflejar lo creado.
	encontrado, ok := almacen.BuscarGuardavidaPorID(int(creado.ID))
	if !ok {
		t.Fatal("se esperaba encontrar el guardavida recién creado")
	}
	if encontrado.Nombre != "Pedro Salazar" || encontrado.Turno != "tarde" {
		t.Errorf("datos inesperados: %+v", encontrado)
	}

	// Listar también debe reflejarlo.
	lista := almacen.ListarGuardavidas()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 guardavida en la lista, se obtuvieron %d", len(lista))
	}
	if lista[0].Nombre != "Pedro Salazar" {
		t.Errorf("nombre inesperado en la lista: %s", lista[0].Nombre)
	}
}

// TestAlmacenSQLite_BuscarGuardavidaInexistente prueba que buscar un ID que
// no existe devuelve ok=false, en lugar de un guardavida vacío con ok=true.
func TestAlmacenSQLite_BuscarGuardavidaInexistente(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	_, ok := almacen.BuscarGuardavidaPorID(999)
	if ok {
		t.Error("no se esperaba encontrar un guardavida con ID 999")
	}
}
