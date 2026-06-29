package storage

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

func abrirDBClientesPrueba(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir la base de datos en memoria: %v", err)
	}
	if err := db.AutoMigrate(&models.Cliente{}); err != nil {
		t.Fatalf("falló AutoMigrate: %v", err)
	}
	return db
}

func TestAlmacenSQLite_CrearYListarCliente(t *testing.T) {
	db := abrirDBClientesPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	nuevo := models.Cliente{
		Nombre:    "Carlos",
		Cedula:    "1712345678",
		Email:     "carlos@example.com",
		Membresia: "mensual",
	}

	creado := almacen.CrearCliente(nuevo)
	if creado.ID == 0 {
		t.Fatal("se esperaba que GORM asignara un ID al crear el cliente")
	}

	encontrado, ok := almacen.BuscarClientePorID(int(creado.ID))
	if !ok {
		t.Fatal("se esperaba encontrar el cliente recién creado")
	}
	if encontrado.Nombre != "Carlos" || encontrado.Cedula != "1712345678" {
		t.Fatalf("datos inesperados: %+v", encontrado)
	}

	lista := almacen.ListarClientes()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 cliente en la lista, se obtuvieron %d", len(lista))
	}
	if lista[0].Nombre != "Carlos" {
		t.Fatalf("nombre inesperado en la lista: %s", lista[0].Nombre)
	}
}
