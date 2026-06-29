package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"pool-api/internal/middleware"
	"pool-api/internal/models"
	"pool-api/internal/service"
)

type fakeClientesRepo struct {
	clientes    []models.Cliente
	siguienteID int
}

func (f *fakeClientesRepo) ListarClientes() []models.Cliente { return f.clientes }
func (f *fakeClientesRepo) BuscarClientePorID(id int) (models.Cliente, bool) {
	for _, c := range f.clientes {
		if int(c.ID) == id {
			return c, true
		}
	}
	return models.Cliente{}, false
}
func (f *fakeClientesRepo) CrearCliente(c models.Cliente) models.Cliente {
	f.siguienteID++
	c.ID = uint(f.siguienteID)
	f.clientes = append(f.clientes, c)
	return c
}
func (f *fakeClientesRepo) ActualizarCliente(id int, datos models.Cliente) (models.Cliente, bool) {
	for i, c := range f.clientes {
		if int(c.ID) == id {
			datos.ID = c.ID
			f.clientes[i] = datos
			return f.clientes[i], true
		}
	}
	return models.Cliente{}, false
}
func (f *fakeClientesRepo) BorrarCliente(id int) bool {
	for i, c := range f.clientes {
		if int(c.ID) == id {
			f.clientes = append(f.clientes[:i], f.clientes[i+1:]...)
			return true
		}
	}
	return false
}

func (f *fakeClientesRepo) ListarReservas() []models.Reserva { return nil }
func (f *fakeClientesRepo) BuscarReservaPorID(id int) (models.Reserva, bool) {
	return models.Reserva{}, false
}
func (f *fakeClientesRepo) CrearReserva(rv models.Reserva) models.Reserva { return rv }
func (f *fakeClientesRepo) ActualizarReserva(id int, datos models.Reserva) (models.Reserva, bool) {
	return models.Reserva{}, false
}
func (f *fakeClientesRepo) BorrarReserva(id int) bool                  { return false }
func (f *fakeClientesRepo) ListarPagos() []models.Pago                 { return nil }
func (f *fakeClientesRepo) BuscarPagoPorID(id int) (models.Pago, bool) { return models.Pago{}, false }
func (f *fakeClientesRepo) CrearPago(p models.Pago) models.Pago        { return p }
func (f *fakeClientesRepo) ActualizarPago(id int, datos models.Pago) (models.Pago, bool) {
	return models.Pago{}, false
}
func (f *fakeClientesRepo) BorrarPago(id int) bool                     { return false }
func (f *fakeClientesRepo) ClienteTienePagoEntrada(clienteID int) bool { return false }

type fakeUsuariosRepo struct {
	usuario models.Usuario
}

func (f *fakeUsuariosRepo) ListarUsuarios() []models.Usuario { return []models.Usuario{f.usuario} }
func (f *fakeUsuariosRepo) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	if int(f.usuario.ID) == id {
		return f.usuario, true
	}
	return models.Usuario{}, false
}
func (f *fakeUsuariosRepo) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	if f.usuario.Email == email {
		return f.usuario, true
	}
	return models.Usuario{}, false
}
func (f *fakeUsuariosRepo) CrearUsuario(u models.Usuario) (models.Usuario, error) { return u, nil }
func (f *fakeUsuariosRepo) ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (f *fakeUsuariosRepo) BorrarUsuario(id int) bool { return false }

func montarRouterClientesPrueba(t *testing.T) (chi.Router, string) {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("clave123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("no se pudo generar el hash: %v", err)
	}
	usuarios := &fakeUsuariosRepo{usuario: models.Usuario{ID: 1, Nombre: "Admin", Email: "admin@prueba.com", PasswordHash: string(hash), Rol: "admin"}}
	authSvc := service.NewAuthService(usuarios)
	token, _, err := authSvc.Login("admin@prueba.com", "clave123")
	if err != nil {
		t.Fatalf("no se pudo generar el token: %v", err)
	}

	clientesSvc := service.NewClientesService(&fakeClientesRepo{})
	srv := NewServer(nil, nil, clientesSvc, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Route("/clientes", func(r chi.Router) {
				r.Get("/", srv.ListarClientes)
				r.Post("/", srv.CrearCliente)
			})
		})
	})
	return r, token
}

func TestListarClientes_SinToken_Devuelve401(t *testing.T) {
	router, _ := montarRouterClientesPrueba(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", rec.Code)
	}
}

func TestCrearCliente_ConToken_CreaYLoDevuelve(t *testing.T) {
	router, token := montarRouterClientesPrueba(t)

	body, _ := json.Marshal(models.Cliente{Nombre: "Marta", Cedula: "1712345678", Membresia: "mensual"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes/", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	var creado models.Cliente
	if err := json.Unmarshal(rec.Body.Bytes(), &creado); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if creado.ID == 0 || creado.Nombre != "Marta" {
		t.Fatalf("cliente creado inesperado: %+v", creado)
	}
}
