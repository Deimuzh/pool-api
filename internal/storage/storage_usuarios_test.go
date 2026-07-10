package storage

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

func abrirDBUsuariosPrueba(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir la base de datos en memoria: %v", err)
	}
	if err := db.AutoMigrate(&models.Usuario{}); err != nil {
		t.Fatalf("fallo AutoMigrate: %v", err)
	}
	return db
}

func TestAlmacenSQLite_UsuarioCRUD(t *testing.T) {
	db := abrirDBUsuariosPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	nuevo, err := almacen.CrearUsuario(models.Usuario{
		Nombre: "Admin", Email: "admin@test.com", PasswordHash: "hash123", Rol: "admin",
	})
	if err != nil {
		t.Fatalf("CrearUsuario fallo: %v", err)
	}
	if nuevo.ID == 0 {
		t.Fatal("se esperaba ID")
	}
	if nuevo.Email != "admin@test.com" {
		t.Errorf("email inesperado: %s", nuevo.Email)
	}

	encontrado, ok := almacen.BuscarUsuarioPorID(nuevo.ID)
	if !ok {
		t.Fatal("no se encontro el usuario recien creado")
	}
	if encontrado.Nombre != "Admin" {
		t.Errorf("nombre inesperado: %s", encontrado.Nombre)
	}

	porEmail, ok := almacen.BuscarUsuarioPorEmail("admin@test.com")
	if !ok {
		t.Fatal("no se encontro usuario por email")
	}
	if porEmail.ID != nuevo.ID {
		t.Errorf("ID inesperado: %d", porEmail.ID)
	}

	lista := almacen.ListarUsuarios()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 usuario, se obtuvieron %d", len(lista))
	}

	actualizado, ok := almacen.ActualizarUsuario(nuevo.ID, models.Usuario{Nombre: "Admin Modificado", Email: "admin@test.com"})
	if !ok {
		t.Fatal("ActualizarUsuario devolvio false")
	}
	if actualizado.Nombre != "Admin Modificado" {
		t.Errorf("nombre inesperado: %s", actualizado.Nombre)
	}

	if !almacen.BorrarUsuario(nuevo.ID) {
		t.Fatal("BorrarUsuario devolvio false")
	}
	if almacen.BorrarUsuario(999) {
		t.Error("BorrarUsuario con ID inexistente devolvio true")
	}
}

func TestAlmacenSQLite_BuscarUsuarioPorEmailInexistente(t *testing.T) {
	db := abrirDBUsuariosPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	_, ok := almacen.BuscarUsuarioPorEmail("no-existe@test.com")
	if ok {
		t.Error("no se esperaba encontrar un usuario con ese email")
	}
}

func TestAlmacenSQLite_BuscarUsuarioPorIDInexistente(t *testing.T) {
	db := abrirDBUsuariosPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	_, ok := almacen.BuscarUsuarioPorID(999)
	if ok {
		t.Error("no se esperaba encontrar un usuario con ID 999")
	}
}
