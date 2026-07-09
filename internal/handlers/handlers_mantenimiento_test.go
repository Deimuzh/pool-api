package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"pool-api/internal/middleware"
	"pool-api/internal/models"
	"pool-api/internal/service"
	"pool-api/internal/storage"
)

type fakeMantenimientoRepo struct {
	equipos  map[uint]models.Equipo
	registros map[uint]models.RegistroMantenimiento
	quimicos map[uint]models.ProductoQuimico
	seq      uint
}

func newFakeMantenimientoRepo() *fakeMantenimientoRepo {
	return &fakeMantenimientoRepo{
		equipos:   make(map[uint]models.Equipo),
		registros: make(map[uint]models.RegistroMantenimiento),
		quimicos:  make(map[uint]models.ProductoQuimico),
		seq:       0,
	}
}

func (f *fakeMantenimientoRepo) ListarEquipos() []models.Equipo {
	lista := make([]models.Equipo, 0, len(f.equipos))
	for _, e := range f.equipos {
		lista = append(lista, e)
	}
	return lista
}
func (f *fakeMantenimientoRepo) BuscarEquipoPorID(id uint) (models.Equipo, bool) {
	e, ok := f.equipos[id]
	return e, ok
}
func (f *fakeMantenimientoRepo) CrearEquipo(e models.Equipo) (models.Equipo, error) {
	f.seq++
	e.ID = f.seq
	f.equipos[e.ID] = e
	return e, nil
}
func (f *fakeMantenimientoRepo) ActualizarEquipo(id uint, datos models.Equipo) (models.Equipo, bool) {
	_, ok := f.equipos[id]
	if !ok {
		return models.Equipo{}, false
	}
	datos.ID = id
	f.equipos[id] = datos
	return datos, true
}
func (f *fakeMantenimientoRepo) BorrarEquipo(id uint) bool {
	_, ok := f.equipos[id]
	if !ok {
		return false
	}
	delete(f.equipos, id)
	return true
}
func (f *fakeMantenimientoRepo) ListarRegistros() []models.RegistroMantenimiento {
	lista := make([]models.RegistroMantenimiento, 0, len(f.registros))
	for _, r := range f.registros {
		lista = append(lista, r)
	}
	return lista
}
func (f *fakeMantenimientoRepo) BuscarRegistroPorID(id uint) (models.RegistroMantenimiento, bool) {
	r, ok := f.registros[id]
	return r, ok
}
func (f *fakeMantenimientoRepo) CrearRegistro(rm models.RegistroMantenimiento) (models.RegistroMantenimiento, error) {
	f.seq++
	rm.ID = f.seq
	f.registros[rm.ID] = rm
	return rm, nil
}
func (f *fakeMantenimientoRepo) ActualizarRegistro(id uint, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool) {
	_, ok := f.registros[id]
	if !ok {
		return models.RegistroMantenimiento{}, false
	}
	datos.ID = id
	f.registros[id] = datos
	return datos, true
}
func (f *fakeMantenimientoRepo) BorrarRegistro(id uint) bool {
	_, ok := f.registros[id]
	if !ok {
		return false
	}
	delete(f.registros, id)
	return true
}
func (f *fakeMantenimientoRepo) ListarQuimicos() []models.ProductoQuimico {
	lista := make([]models.ProductoQuimico, 0, len(f.quimicos))
	for _, q := range f.quimicos {
		lista = append(lista, q)
	}
	return lista
}
func (f *fakeMantenimientoRepo) BuscarQuimicoPorID(id uint) (models.ProductoQuimico, bool) {
	q, ok := f.quimicos[id]
	return q, ok
}
func (f *fakeMantenimientoRepo) CrearQuimico(q models.ProductoQuimico) (models.ProductoQuimico, error) {
	f.seq++
	q.ID = f.seq
	f.quimicos[q.ID] = q
	return q, nil
}
func (f *fakeMantenimientoRepo) ActualizarQuimico(id uint, datos models.ProductoQuimico) (models.ProductoQuimico, bool) {
	_, ok := f.quimicos[id]
	if !ok {
		return models.ProductoQuimico{}, false
	}
	datos.ID = id
	f.quimicos[id] = datos
	return datos, true
}
func (f *fakeMantenimientoRepo) BorrarQuimico(id uint) bool {
	_, ok := f.quimicos[id]
	if !ok {
		return false
	}
	delete(f.quimicos, id)
	return true
}

var _ storage.MantenimientoRepository = (*fakeMantenimientoRepo)(nil)

func montarRouterMantenimientoPrueba(t *testing.T) (chi.Router, string) {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("clave123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("no se pudo generar el hash: %v", err)
	}
	usuarios := &fakeUsuarioRepo{
		usuario: models.Usuario{ID: 1, Nombre: "Admin", Email: "admin@test.com", PasswordHash: string(hash), Rol: "admin"},
	}
	authSvc := service.NewAuthService(usuarios)

	token, _, err := authSvc.Login("admin@test.com", "clave123")
	if err != nil {
		t.Fatalf("no se pudo generar token: %v", err)
	}

	mantSvc := service.NewMantenimientoService(newFakeMantenimientoRepo())
	srv := NewServer(nil, mantSvc, nil, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Route("/equipos", func(r chi.Router) {
				r.Get("/", srv.ListarEquipos)
				r.Post("/", srv.CrearEquipo)
				r.Get("/{id}", srv.ObtenerEquipo)
				r.Put("/{id}", srv.ActualizarEquipo)
				r.Delete("/{id}", srv.BorrarEquipo)
			})
			r.Route("/registros-mantenimiento", func(r chi.Router) {
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
	})

	return r, token
}

// ─── EQUIPO ──────────────────────────────────────────────────────────────────

func TestListarEquipos_SinToken_Devuelve401(t *testing.T) {
	router, _ := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

func TestListarEquipos_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestCrearEquipo_ConToken_Crea(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	body, _ := json.Marshal(models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "operativo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	var e models.Equipo
	json.Unmarshal(rec.Body.Bytes(), &e)
	if e.ID == 0 {
		t.Fatal("se esperaba ID generado")
	}
}

func TestObtenerEquipo_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	body, _ := json.Marshal(models.Equipo{Nombre: "Filtro", Tipo: "filtro"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var creado models.Equipo
	json.Unmarshal(rec.Body.Bytes(), &creado)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/"+strconv.Itoa(int(creado.ID)), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec2.Code)
	}
}

func TestObtenerEquipo_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestActualizarEquipo_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	body, _ := json.Marshal(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.Equipo
	json.Unmarshal(rec.Body.Bytes(), &creado)

	body2, _ := json.Marshal(models.Equipo{Nombre: "Bomba V2", Tipo: "bomba", Estado: "operativo"})
	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/equipos/"+strconv.Itoa(int(creado.ID)), bytes.NewReader(body2))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec2.Code, rec2.Body.String())
	}
}

func TestBorrarEquipo_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	body, _ := json.Marshal(models.Equipo{Nombre: "Bomba", Tipo: "bomba"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.Equipo
	json.Unmarshal(rec.Body.Bytes(), &creado)

	req2 := httptest.NewRequest(http.MethodDelete, "/api/v1/equipos/"+strconv.Itoa(int(creado.ID)), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec2.Code)
	}
}

func TestBorrarEquipo_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/equipos/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestObtenerEquipo_IDInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipos/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestCrearEquipo_JSONInvalido_Devuelve400(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader([]byte("{json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

// ─── REGISTRO MANTENIMIENTO ─────────────────────────────────────────────────

func TestListarRegistros_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/registros-mantenimiento/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func crearEquipoEnRouter(t *testing.T, router chi.Router, token string) uint {
	t.Helper()
	body, _ := json.Marshal(models.Equipo{Nombre: "Equipo Test", Tipo: "filtro", Estado: "operativo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/equipos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var e models.Equipo
	json.Unmarshal(rec.Body.Bytes(), &e)
	return e.ID
}

func TestCrearRegistro_ConToken_Crea(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	eqID := crearEquipoEnRouter(t, router, token)

	body, _ := json.Marshal(models.RegistroMantenimiento{EquipoID: eqID, Descripcion: "Cambio aceite", Tipo: "preventivo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/registros-mantenimiento/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestObtenerRegistro_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	eqID := crearEquipoEnRouter(t, router, token)
	body, _ := json.Marshal(models.RegistroMantenimiento{EquipoID: eqID, Descripcion: "Test", Tipo: "preventivo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/registros-mantenimiento/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.RegistroMantenimiento
	json.Unmarshal(rec.Body.Bytes(), &creado)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/registros-mantenimiento/"+strconv.Itoa(int(creado.ID)), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec2.Code)
	}
}

func TestObtenerRegistro_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/registros-mantenimiento/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestActualizarRegistro_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	eqID := crearEquipoEnRouter(t, router, token)
	body, _ := json.Marshal(models.RegistroMantenimiento{EquipoID: eqID, Descripcion: "Test", Tipo: "preventivo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/registros-mantenimiento/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.RegistroMantenimiento
	json.Unmarshal(rec.Body.Bytes(), &creado)

	body2, _ := json.Marshal(models.RegistroMantenimiento{EquipoID: 1, Descripcion: "Modificado", Tipo: "correctivo"})
	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/registros-mantenimiento/"+strconv.Itoa(int(creado.ID)), bytes.NewReader(body2))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec2.Code, rec2.Body.String())
	}
}

func TestBorrarRegistro_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	eqID := crearEquipoEnRouter(t, router, token)
	body, _ := json.Marshal(models.RegistroMantenimiento{EquipoID: eqID, Descripcion: "Test", Tipo: "preventivo"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/registros-mantenimiento/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.RegistroMantenimiento
	json.Unmarshal(rec.Body.Bytes(), &creado)

	req2 := httptest.NewRequest(http.MethodDelete, "/api/v1/registros-mantenimiento/"+strconv.Itoa(int(creado.ID)), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec2.Code)
	}
}

func TestBorrarRegistro_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/registros-mantenimiento/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

// ─── PRODUCTO QUIMICO ────────────────────────────────────────────────────────

func TestListarQuimicos_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quimicos/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
}

func TestCrearQuimico_ConToken_Crea(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	body, _ := json.Marshal(models.ProductoQuimico{Nombre: "Cloro", StockActual: 100, UnidadMedida: "kg", NivelMinimo: 10})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/quimicos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	var q models.ProductoQuimico
	json.Unmarshal(rec.Body.Bytes(), &q)
	if q.ID == 0 {
		t.Fatal("se esperaba ID generado")
	}
}

func TestObtenerQuimico_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	body, _ := json.Marshal(models.ProductoQuimico{Nombre: "Cloro", StockActual: 50})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/quimicos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.ProductoQuimico
	json.Unmarshal(rec.Body.Bytes(), &creado)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/quimicos/"+strconv.Itoa(int(creado.ID)), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec2.Code)
	}
}

func TestActualizarQuimico_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	body, _ := json.Marshal(models.ProductoQuimico{Nombre: "Cloro", StockActual: 50})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/quimicos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.ProductoQuimico
	json.Unmarshal(rec.Body.Bytes(), &creado)

	body2, _ := json.Marshal(models.ProductoQuimico{Nombre: "Cloro Plus", StockActual: 80})
	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/quimicos/"+strconv.Itoa(int(creado.ID)), bytes.NewReader(body2))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec2.Code, rec2.Body.String())
	}
}

func TestBorrarQuimico_ConToken_OK(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	body, _ := json.Marshal(models.ProductoQuimico{Nombre: "Cloro", StockActual: 50})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/quimicos/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var creado models.ProductoQuimico
	json.Unmarshal(rec.Body.Bytes(), &creado)

	req2 := httptest.NewRequest(http.MethodDelete, "/api/v1/quimicos/"+strconv.Itoa(int(creado.ID)), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec2.Code)
	}
}

func TestBorrarQuimico_NoEncontrado_Devuelve404(t *testing.T) {
	router, token := montarRouterMantenimientoPrueba(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/quimicos/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}


