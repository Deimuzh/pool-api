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

func TestEquipoRepo_BuscarPorID_NoEncontrado(t *testing.T) {
	db := abrirDBEquipo(t)
	repo := NuevoAlmacenSQLite(db)

	_, ok := repo.BuscarEquipoPorID(999)
	require.False(t, ok)
}
