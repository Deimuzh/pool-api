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
	if err := db.AutoMigrate(
		&models.Guardavida{},
		&models.Incidente{},
		&models.AccesoCliente{},
	); err != nil {
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
	encontrado, ok := almacen.BuscarGuardavidaPorID(creado.ID)
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

// TestAlmacenSQLite_CrearYListarAcceso prueba que el repositorio real pueda
// guardar y listar accesos usando GORM con SQLite en memoria.
func TestAlmacenSQLite_CrearYListarAcceso(t *testing.T) {
	db := abrirDBPrueba(t)
	if err := db.AutoMigrate(&models.AccesoCliente{}); err != nil {
		t.Fatalf("falló AutoMigrate de AccesoCliente: %v", err)
	}

	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearAcceso(models.AccesoCliente{
		ClienteID:  2,
		Autorizado: true,
	})

	if creado.ID == 0 {
		t.Fatal("se esperaba que GORM asignara un ID al acceso")
	}

	lista := almacen.ListarAccesos()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 acceso, se obtuvieron %d", len(lista))
	}

	if lista[0].ClienteID != 2 {
		t.Errorf("cliente_id inesperado: %d", lista[0].ClienteID)
	}

	if !lista[0].Autorizado {
		t.Error("se esperaba que el acceso esté autorizado")
	}
}

// TestAlmacenSQLite_BuscarAccesoPorID prueba que un acceso recién creado pueda
// recuperarse por su ID desde el repositorio real.
func TestAlmacenSQLite_BuscarAccesoPorID(t *testing.T) {
	db := abrirDBPrueba(t)
	if err := db.AutoMigrate(&models.AccesoCliente{}); err != nil {
		t.Fatalf("falló AutoMigrate de AccesoCliente: %v", err)
	}

	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearAcceso(models.AccesoCliente{
		ClienteID:  3,
		Autorizado: true,
	})

	encontrado, ok := almacen.BuscarAccesoPorID(creado.ID)
	if !ok {
		t.Fatal("se esperaba encontrar el acceso recién creado")
	}

	if encontrado.ClienteID != 3 {
		t.Errorf("cliente_id inesperado: %d", encontrado.ClienteID)
	}
}
// TestAlmacenSQLite_ActualizarGuardavida verifica que el repositorio real pueda
// actualizar los datos de un guardavida existente usando SQLite en memoria.
func TestAlmacenSQLite_ActualizarGuardavida(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearGuardavida(models.Guardavida{
		Nombre: "Pedro Salazar",
		Turno:  "mañana",
	})

	actualizado, ok := almacen.ActualizarGuardavida(creado.ID, models.Guardavida{
		Nombre: "Pedro Actualizado",
		Turno:  "noche",
	})

	if !ok {
		t.Fatal("se esperaba actualizar el guardavida")
	}
	if actualizado.Nombre != "Pedro Actualizado" || actualizado.Turno != "noche" {
		t.Fatalf("datos inesperados: %+v", actualizado)
	}
}

// TestAlmacenSQLite_BorrarGuardavida verifica que el repositorio real pueda
// eliminar un guardavida y que luego ya no pueda encontrarse por ID.
func TestAlmacenSQLite_BorrarGuardavida(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearGuardavida(models.Guardavida{
		Nombre: "Pedro Salazar",
		Turno:  "mañana",
	})

	if !almacen.BorrarGuardavida(creado.ID) {
		t.Fatal("se esperaba borrar el guardavida")
	}
	if _, ok := almacen.BuscarGuardavidaPorID(creado.ID); ok {
		t.Fatal("no se esperaba encontrar el guardavida borrado")
	}
}

// TestAlmacenSQLite_CrearListarYBuscarIncidente verifica que el repositorio real
// pueda crear, listar y buscar incidentes usando SQLite en memoria.
func TestAlmacenSQLite_CrearListarYBuscarIncidente(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	incidente := almacen.CrearIncidente(models.Incidente{
		Tipo:         "lesion",
		Gravedad:     "leve",
		GuardavidaID: 1,
		ClienteID:    2,
	})

	if incidente.ID == 0 {
		t.Fatal("se esperaba que GORM asigne ID al incidente")
	}

	lista := almacen.ListarIncidentes()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 incidente, se obtuvieron %d", len(lista))
	}

	encontrado, ok := almacen.BuscarIncidentePorID(incidente.ID)
	if !ok {
		t.Fatal("se esperaba encontrar el incidente")
	}
	if encontrado.Tipo != "lesion" {
		t.Fatalf("tipo inesperado: %s", encontrado.Tipo)
	}
}

// TestAlmacenSQLite_ActualizarIncidente verifica que el repositorio real pueda
// actualizar un incidente existente y persistir sus nuevos datos.
func TestAlmacenSQLite_ActualizarIncidente(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearIncidente(models.Incidente{
		Tipo:         "lesion",
		Gravedad:     "leve",
		GuardavidaID: 1,
		ClienteID:    2,
	})

	actualizado, ok := almacen.ActualizarIncidente(creado.ID, models.Incidente{
		Tipo:         "rescate",
		Gravedad:     "media",
		GuardavidaID: 1,
		ClienteID:    2,
	})

	if !ok {
		t.Fatal("se esperaba actualizar el incidente")
	}
	if actualizado.Tipo != "rescate" || actualizado.Gravedad != "media" {
		t.Fatalf("datos inesperados: %+v", actualizado)
	}
}

// TestAlmacenSQLite_BorrarIncidente verifica que el repositorio real pueda
// eliminar un incidente y que ya no pueda encontrarse por ID.
func TestAlmacenSQLite_BorrarIncidente(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearIncidente(models.Incidente{
		Tipo:         "lesion",
		Gravedad:     "leve",
		GuardavidaID: 1,
		ClienteID:    2,
	})

	if !almacen.BorrarIncidente(creado.ID) {
		t.Fatal("se esperaba borrar el incidente")
	}
	if _, ok := almacen.BuscarIncidentePorID(creado.ID); ok {
		t.Fatal("no se esperaba encontrar el incidente borrado")
	}
}

// TestAlmacenSQLite_ActualizarAcceso verifica que el repositorio real pueda
// actualizar el estado de autorización de un acceso existente.
func TestAlmacenSQLite_ActualizarAcceso(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearAcceso(models.AccesoCliente{
		ClienteID:  2,
		Autorizado: false,
		Motivo:     "Sin pago",
	})

	actualizado, ok := almacen.ActualizarAcceso(creado.ID, models.AccesoCliente{
		ClienteID:  2,
		Autorizado: true,
		Motivo:     "",
	})

	if !ok {
		t.Fatal("se esperaba actualizar el acceso")
	}
	if !actualizado.Autorizado || actualizado.Motivo != "" {
		t.Fatalf("datos inesperados: %+v", actualizado)
	}
}

// TestAlmacenSQLite_BorrarAcceso verifica que el repositorio real pueda eliminar
// un acceso y que luego ya no pueda encontrarse por ID.
func TestAlmacenSQLite_BorrarAcceso(t *testing.T) {
	db := abrirDBPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearAcceso(models.AccesoCliente{
		ClienteID:  2,
		Autorizado: true,
	})

	if !almacen.BorrarAcceso(creado.ID) {
		t.Fatal("se esperaba borrar el acceso")
	}
	if _, ok := almacen.BuscarAccesoPorID(creado.ID); ok {
		t.Fatal("no se esperaba encontrar el acceso borrado")
	}
}