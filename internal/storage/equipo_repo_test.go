package storage

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

func abrirDBEquipo(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Equipo{}))
	return db
}

func TestEquipoRepo_ListarVacio(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)
	lista := repo.ListarEquipos()
	require.Empty(t, lista)
}

func TestEquipoRepo_Crear_ErrorDB(t *testing.T) {
	db := abrirDBEquipo(t)
	sqlDB, _ := db.DB()
	sqlDB.Close()

	repo := NuevoAlmacenSQLite(db)
	_, err := repo.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	require.Error(t, err)
}

func TestEquipoRepo_CrearYListar(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	creado, err := repo.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "operativo"})
	require.NoError(t, err)
	require.NotZero(t, creado.ID)

	lista := repo.ListarEquipos()
	require.Len(t, lista, 1)
	require.Equal(t, "Bomba", lista[0].Nombre)
}

func TestEquipoRepo_BuscarPorID(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	creado, _ := repo.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion", Estado: "operativo"})

	encontrado, ok := repo.BuscarEquipoPorID(creado.ID)
	require.True(t, ok)
	require.Equal(t, "Filtro", encontrado.Nombre)
}

func TestEquipoRepo_Actualizar_CambiaTipoYEstado(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	creado, _ := repo.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion", Estado: "operativo"})

	actualizado, ok := repo.ActualizarEquipo(creado.ID, models.Equipo{Nombre: "Filtro", Tipo: "bomba", Estado: "averiado"})
	require.True(t, ok)
	require.Equal(t, "bomba", actualizado.Tipo)
	require.Equal(t, "averiado", actualizado.Estado)
}

func TestEquipoRepo_Actualizar(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	creado, _ := repo.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "operativo"})

	actualizado, ok := repo.ActualizarEquipo(creado.ID, models.Equipo{Nombre: "Bomba 2HP"})
	require.True(t, ok)
	require.Equal(t, "Bomba 2HP", actualizado.Nombre)
}

func TestEquipoRepo_Actualizar_NoEncontrado(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	_, ok := repo.ActualizarEquipo(999, models.Equipo{Nombre: "N/A"})
	require.False(t, ok)
}

func TestEquipoRepo_Borrar(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	creado, _ := repo.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion", Estado: "operativo"})

	ok := repo.BorrarEquipo(creado.ID)
	require.True(t, ok)

	_, encontrado := repo.BuscarEquipoPorID(creado.ID)
	require.False(t, encontrado)
}

func TestEquipoRepo_Borrar_NoEncontrado(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	ok := repo.BorrarEquipo(999)
	require.False(t, ok)
}

func TestEquipoRepo_BuscarPorID_NoEncontrado(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	_, ok := repo.BuscarEquipoPorID(999)
	require.False(t, ok)
}
