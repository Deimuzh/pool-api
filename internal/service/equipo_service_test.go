package service

import (
	"testing"

	"pool-api/internal/models"
)

func TestEquipoService_New(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)
	if svc == nil {
		t.Fatal("NewMantenimientoService no debe devolver nil")
	}
}

func TestEquipoService_ListarVacio(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	resultado := svc.ListarEquipos()
	if len(resultado) != 0 {
		t.Fatalf("se esperaba lista vacia, se obtuvo %d", len(resultado))
	}
}

func TestEquipoService_Obtener_NoEncontrado(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, ok := svc.ObtenerEquipo(999)
	if ok {
		t.Fatal("no debe encontrar equipo inexistente")
	}
}

func TestEquipoService_Obtener_Encontrado(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)

	e, ok := svc.ObtenerEquipo(1)
	if !ok {
		t.Fatal("debe encontrar el equipo")
	}
	if e.Nombre != "Bomba" {
		t.Fatalf("se esperaba Bomba, se obtuvo %s", e.Nombre)
	}
}
