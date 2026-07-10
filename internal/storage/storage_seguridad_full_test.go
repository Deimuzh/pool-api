package storage

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

func abrirDBFullSeguridad(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir la base de datos en memoria: %v", err)
	}
	if err := db.AutoMigrate(&models.Guardavida{}, &models.Incidente{}, &models.AccesoCliente{}, &models.Cliente{}, &models.Pago{}); err != nil {
		t.Fatalf("fallo AutoMigrate: %v", err)
	}
	return db
}

func TestAlmacenSQLite_GuardavidaCRUD(t *testing.T) {
	db := abrirDBFullSeguridad(t)
	almacen := NuevoAlmacenSQLite(db)

	creado := almacen.CrearGuardavida(models.Guardavida{Nombre: "Carlos", Turno: "tarde"})
	if creado.ID == 0 {
		t.Fatal("se esperaba ID")
	}

	actualizado, ok := almacen.ActualizarGuardavida(creado.ID, models.Guardavida{Nombre: "Carlos Modificado", Turno: "noche"})
	if !ok {
		t.Fatal("ActualizarGuardavida devolvio false")
	}
	if actualizado.Nombre != "Carlos Modificado" {
		t.Errorf("nombre inesperado: %s", actualizado.Nombre)
	}

	if !almacen.BorrarGuardavida(creado.ID) {
		t.Fatal("BorrarGuardavida devolvio false")
	}
	if almacen.BorrarGuardavida(999) {
		t.Error("BorrarGuardavida con ID inexistente devolvio true")
	}

	_, ok = almacen.BuscarGuardavidaPorID(creado.ID)
	if ok {
		t.Error("no se esperaba encontrar el guardavida borrado")
	}
}

func TestAlmacenSQLite_IncidenteCRUD(t *testing.T) {
	db := abrirDBFullSeguridad(t)
	almacen := NuevoAlmacenSQLite(db)

	g := almacen.CrearGuardavida(models.Guardavida{Nombre: "Guardia", Turno: "mañana"})
	c := models.Cliente{Nombre: "Cliente", Cedula: "1311111111"}
	creado, _ := almacen.CrearCliente(c)

	nuevo := almacen.CrearIncidente(models.Incidente{
		Tipo: "lesion", Gravedad: "leve", GuardavidaID: g.ID, ClienteID: creado.ID,
		FechaHora: time.Now(), Descripcion: "Corte",
	})
	if nuevo.ID == 0 {
		t.Fatal("se esperaba ID")
	}

	encontrado, ok := almacen.BuscarIncidentePorID(nuevo.ID)
	if !ok {
		t.Fatal("no se encontro el incidente recien creado")
	}
	if encontrado.Tipo != "lesion" {
		t.Errorf("tipo inesperado: %s", encontrado.Tipo)
	}

	lista := almacen.ListarIncidentes()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 incidente, se obtuvieron %d", len(lista))
	}

	actualizado, ok := almacen.ActualizarIncidente(nuevo.ID, models.Incidente{Gravedad: "moderado"})
	if !ok {
		t.Fatal("ActualizarIncidente devolvio false")
	}
	if actualizado.Gravedad != "moderado" {
		t.Errorf("gravedad inesperada: %s", actualizado.Gravedad)
	}

	if !almacen.BorrarIncidente(nuevo.ID) {
		t.Fatal("BorrarIncidente devolvio false")
	}
	if almacen.BorrarIncidente(999) {
		t.Error("BorrarIncidente con ID inexistente devolvio true")
	}
}

func TestAlmacenSQLite_BuscarIncidenteInexistente(t *testing.T) {
	db := abrirDBFullSeguridad(t)
	almacen := NuevoAlmacenSQLite(db)

	_, ok := almacen.BuscarIncidentePorID(999)
	if ok {
		t.Error("no se esperaba encontrar incidente con ID 999")
	}
}

func TestAlmacenSQLite_AccesoCRUD(t *testing.T) {
	db := abrirDBFullSeguridad(t)
	almacen := NuevoAlmacenSQLite(db)

	c := models.Cliente{Nombre: "Ana", Cedula: "1322222222"}
	creado, _ := almacen.CrearCliente(c)

	nuevo := almacen.CrearAcceso(models.AccesoCliente{ClienteID: creado.ID, Autorizado: true, FechaHora: time.Now()})
	if nuevo.ID == 0 {
		t.Fatal("se esperaba ID")
	}

	encontrado, ok := almacen.BuscarAccesoPorID(nuevo.ID)
	if !ok {
		t.Fatal("no se encontro el acceso recien creado")
	}
	if encontrado.ClienteID != creado.ID {
		t.Errorf("ClienteID inesperado: %d", encontrado.ClienteID)
	}

	lista := almacen.ListarAccesos()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 acceso, se obtuvieron %d", len(lista))
	}

	actualizado, ok := almacen.ActualizarAcceso(nuevo.ID, models.AccesoCliente{Autorizado: false, Motivo: "Denegado"})
	if !ok {
		t.Fatal("ActualizarAcceso devolvio false")
	}
	if actualizado.Autorizado {
		t.Error("se esperaba Autorizado=false")
	}

	if !almacen.BorrarAcceso(nuevo.ID) {
		t.Fatal("BorrarAcceso devolvio false")
	}
	if almacen.BorrarAcceso(999) {
		t.Error("BorrarAcceso con ID inexistente devolvio true")
	}
}
