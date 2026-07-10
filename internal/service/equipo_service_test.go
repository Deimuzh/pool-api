package service

import (
	"errors"
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

func TestEquipoService_Crear_ErrorRepo(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.crearEquipoError = errors.New("db error")
	svc := NewMantenimientoService(repo)

	_, err := svc.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	if err == nil || err.Error() != "db error" {
		t.Fatalf("se esperaba 'db error', se obtuvo %v", err)
	}
}

func TestEquipoService_Crear_NombreVacio(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)

	_, err := svc.CrearEquipo(models.Equipo{Tipo: "bomba"})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoService_Crear_TipoVacio(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)

	_, err := svc.CrearEquipo(models.Equipo{Nombre: "Bomba"})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoService_Crear_AsignaEstadoOperativo(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)

	creado, err := svc.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.Estado != "operativo" {
		t.Fatalf("se esperaba estado operativo, se obtuvo %q", creado.Estado)
	}
}

func TestEquipoService_Crear_Exitoso(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)

	creado, err := svc.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "averiado"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.Nombre != "Bomba" {
		t.Fatalf("se esperaba Bomba, se obtuvo %s", creado.Nombre)
	}
	if creado.Estado != "averiado" {
		t.Fatalf("debe conservar el estado enviado, se obtuvo %q", creado.Estado)
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

func TestEquipoService_Actualizar_Exitoso(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "operativo"}
	svc := NewMantenimientoService(repo)

	actualizado, err := svc.ActualizarEquipo(1, models.Equipo{Nombre: "Bomba 2HP", Tipo: "bomba"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Nombre != "Bomba 2HP" {
		t.Fatalf("se esperaba 'Bomba 2HP', se obtuvo %s", actualizado.Nombre)
	}
}

func TestEquipoService_Actualizar_CampoObligatorio(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)

	_, err := svc.ActualizarEquipo(1, models.Equipo{Nombre: "", Tipo: ""})
	if err != ErrCampoObligatorio {
		t.Fatalf("se esperaba ErrCampoObligatorio, se obtuvo %v", err)
	}
}

func TestEquipoService_Actualizar_NoEncontrado(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	_, err := svc.ActualizarEquipo(999, models.Equipo{Nombre: "Nuevo", Tipo: "bomba"})
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestEquipoService_ListarConDatos(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	repo.equipos[2] = models.Equipo{Nombre: "Filtro", Tipo: "filtracion"}
	svc := NewMantenimientoService(repo)

	lista := svc.ListarEquipos()
	if len(lista) != 2 {
		t.Fatalf("se esperaban 2 equipos, se obtuvo %d", len(lista))
	}
}

func TestEquipoService_Borrar_Exitoso(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)

	err := svc.BorrarEquipo(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestEquipoService_Borrar_VerificaListaVaciaDespues(t *testing.T) {
	repo := newMantenimientoRepoMock()
	repo.equipos[1] = models.Equipo{Nombre: "Bomba", Tipo: "bomba"}
	svc := NewMantenimientoService(repo)

	svc.BorrarEquipo(1)

	lista := svc.ListarEquipos()
	if len(lista) != 0 {
		t.Fatalf("la lista debe estar vacia tras borrar, se obtuvo %d", len(lista))
	}
}

func TestEquipoService_Borrar_NoEncontrado(t *testing.T) {
	svc := NewMantenimientoService(newMantenimientoRepoMock())
	err := svc.BorrarEquipo(999)
	if err != ErrNoEncontrado {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}
