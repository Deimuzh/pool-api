package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"pool-api/internal/models"
	"pool-api/internal/service"
	"pool-api/internal/storage"
)

type mantenimientoRepoMock struct{
	equipos map[uint]models.Equipo
}

func newMantenimientoRepoMock() *mantenimientoRepoMock {
	return &mantenimientoRepoMock{
		equipos: map[uint]models.Equipo{
			1: {ID: 1, Nombre: "Bomba", Tipo: "bomba"},
		},
	}
}

var _ storage.MantenimientoRepository = (*mantenimientoRepoMock)(nil)

func (m *mantenimientoRepoMock) ListarEquipos() []models.Equipo { return nil }
func (m *mantenimientoRepoMock) BuscarEquipoPorID(id uint) (models.Equipo, bool) {
	e, ok := m.equipos[id]
	return e, ok
}
func (m *mantenimientoRepoMock) CrearEquipo(e models.Equipo) (models.Equipo, error) { return e, nil }
func (m *mantenimientoRepoMock) ActualizarEquipo(id uint, datos models.Equipo) (models.Equipo, bool) {
	e, ok := m.equipos[id]
	if !ok {
		return models.Equipo{}, false
	}
	if datos.Nombre != "" {
		e.Nombre = datos.Nombre
	}
	if datos.Tipo != "" {
		e.Tipo = datos.Tipo
	}
	if datos.Estado != "" {
		e.Estado = datos.Estado
	}
	m.equipos[id] = e
	return e, true
}
func (m *mantenimientoRepoMock) BorrarEquipo(id uint) bool {
	_, ok := m.equipos[id]
	if !ok {
		return false
	}
	delete(m.equipos, id)
	return true
}
func (m *mantenimientoRepoMock) ListarRegistros() []models.RegistroMantenimiento { return nil }
func (m *mantenimientoRepoMock) BuscarRegistroPorID(id uint) (models.RegistroMantenimiento, bool) {
	return models.RegistroMantenimiento{}, false
}
func (m *mantenimientoRepoMock) CrearRegistro(rm models.RegistroMantenimiento) (models.RegistroMantenimiento, error) {
	return rm, nil
}
func (m *mantenimientoRepoMock) ActualizarRegistro(id uint, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool) {
	return models.RegistroMantenimiento{}, false
}
func (m *mantenimientoRepoMock) BorrarRegistro(id uint) bool { return false }
func (m *mantenimientoRepoMock) ListarQuimicos() []models.ProductoQuimico { return nil }
func (m *mantenimientoRepoMock) BuscarQuimicoPorID(id uint) (models.ProductoQuimico, bool) {
	return models.ProductoQuimico{}, false
}
func (m *mantenimientoRepoMock) CrearQuimico(q models.ProductoQuimico) (models.ProductoQuimico, error) {
	return q, nil
}
func (m *mantenimientoRepoMock) ActualizarQuimico(id uint, datos models.ProductoQuimico) (models.ProductoQuimico, bool) {
	return models.ProductoQuimico{}, false
}
func (m *mantenimientoRepoMock) BorrarQuimico(id uint) bool { return false }

func montarRouterMantenimientoPrueba() *chi.Mux {
	mock := newMantenimientoRepoMock()
	svc := service.NewMantenimientoService(mock)
	srv := NewServer(nil, svc, nil, nil)
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/equipos", func(r chi.Router) {
			r.Get("/", srv.ListarEquipos)
			r.Post("/", srv.CrearEquipo)
			r.Get("/{id}", srv.ObtenerEquipo)
			r.Put("/{id}", srv.ActualizarEquipo)
			r.Delete("/{id}", srv.BorrarEquipo)
		})
		r.Route("/mantenimientos", func(r chi.Router) {
			r.Get("/", srv.ListarRegistrosMantenimiento)
			r.Post("/", srv.CrearRegistroMantenimiento)
			r.Get("/{id}", srv.ObtenerRegistroMantenimiento)
			r.Put("/{id}", srv.ActualizarRegistroMantenimiento)
			r.Delete("/{id}", srv.BorrarRegistroMantenimiento)
		})
		r.Route("/quimicos", func(r chi.Router) {
			r.Get("/", srv.ListarQuimicos)
			r.Post("/", srv.CrearQuimico)
			r.Get("/{id}", srv.ObtenerQuimico)
			r.Put("/{id}", srv.ActualizarQuimico)
			r.Delete("/{id}", srv.BorrarQuimico)
		})
	})
	return r
}

func TestListarEquipos(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestCrearEquipo_Handler(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	body, _ := json.Marshal(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCrearEquipo_CampoObligatorio(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	body, _ := json.Marshal(models.Equipo{Nombre: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCrearEquipo_JSONInvalido(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader([]byte("{mal")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestObtenerEquipo_NoEncontrado(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestObtenerEquipo_Handler_Encontrado(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
	var eq models.Equipo
	json.NewDecoder(rec.Body).Decode(&eq)
	if eq.Nombre != "Bomba" {
		t.Fatalf("se esperaba Bomba, se obtuvo %s", eq.Nombre)
	}
}

func TestActualizarEquipo_Handler(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	body, _ := json.Marshal(models.Equipo{Nombre: "Bomba 2HP", Tipo: "bomba"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/equipos/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestActualizarEquipo_JSONInvalido(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/equipos/1", bytes.NewReader([]byte("{mal")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarEquipo_IDInvalido(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	body, _ := json.Marshal(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/equipos/abc", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestActualizarEquipo_CampoObligatorio(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	body, _ := json.Marshal(models.Equipo{Nombre: ""})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/equipos/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestBorrarEquipo_Handler(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/equipos/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestBorrarEquipo_NoEncontrado(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/equipos/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestBorrarEquipo_IDInvalido(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/equipos/abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestObtenerEquipo_IDInvalido(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestListarRegistrosMantenimiento(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/mantenimientos/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestCrearRegistroMantenimiento_Handler(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	body, _ := json.Marshal(models.RegistroMantenimiento{EquipoID: 1, Tipo: "preventivo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/mantenimientos/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestListarQuimicos(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quimicos/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestCrearQuimico_Handler(t *testing.T) {
	router := montarRouterMantenimientoPrueba()
	body, _ := json.Marshal(models.ProductoQuimico{Nombre: "Cloro"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/quimicos/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}
