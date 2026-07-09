package storage

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

func abrirDBCompleta(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(
		&models.Guardavida{},
		&models.Incidente{},
		&models.AccesoCliente{},
		&models.Equipo{},
		&models.RegistroMantenimiento{},
		&models.ProductoQuimico{},
		&models.Cliente{},
		&models.Reserva{},
		&models.Pago{},
		&models.Usuario{},
	)
	require.NoError(t, err)
	return db
}

func TestAlmacenSQLite_IncidenteCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	g := a.CrearGuardavida(models.Guardavida{Nombre: "Carlos", Turno: "mañana"})

	inc := a.CrearIncidente(models.Incidente{Tipo: "lesion", Gravedad: "leve", GuardavidaID: g.ID, ClienteID: 1})
	require.NotZero(t, inc.ID)

	lista := a.ListarIncidentes()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarIncidentePorID(inc.ID)
	require.True(t, ok)
	require.Equal(t, "lesion", encontrado.Tipo)

	_, ok = a.BuscarIncidentePorID(999)
	require.False(t, ok)

	actualizado, ok := a.ActualizarIncidente(inc.ID, models.Incidente{Tipo: "ahogamiento", Gravedad: "grave", GuardavidaID: g.ID})
	require.True(t, ok)
	require.Equal(t, "ahogamiento", actualizado.Tipo)

	ok = a.BorrarIncidente(inc.ID)
	require.True(t, ok)

	ok = a.BorrarIncidente(999)
	require.False(t, ok)
}

func TestAlmacenSQLite_AccesoCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	acc := a.CrearAcceso(models.AccesoCliente{ClienteID: 1, Autorizado: true})
	require.NotZero(t, acc.ID)

	lista := a.ListarAccesos()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarAccesoPorID(acc.ID)
	require.True(t, ok)
	require.True(t, encontrado.Autorizado)

	_, ok = a.BuscarAccesoPorID(999)
	require.False(t, ok)

	actualizado, ok := a.ActualizarAcceso(acc.ID, models.AccesoCliente{ClienteID: 1, Autorizado: false})
	require.True(t, ok)
	require.False(t, actualizado.Autorizado)

	ok = a.BorrarAcceso(acc.ID)
	require.True(t, ok)
}

func TestAlmacenSQLite_EquipoCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	e, err := a.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "operativo"})
	require.NoError(t, err)
	require.NotZero(t, e.ID)

	lista := a.ListarEquipos()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarEquipoPorID(e.ID)
	require.True(t, ok)
	require.Equal(t, "Bomba", encontrado.Nombre)

	_, ok = a.BuscarEquipoPorID(999)
	require.False(t, ok)
}

func TestAlmacenSQLite_RegistroCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	eq, _ := a.CrearEquipo(models.Equipo{Nombre: "Filtro", Tipo: "filtracion"})
	r, err := a.CrearRegistro(models.RegistroMantenimiento{EquipoID: eq.ID, Tipo: "preventivo"})
	require.NoError(t, err)
	require.NotZero(t, r.ID)

	lista := a.ListarRegistros()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarRegistroPorID(r.ID)
	require.True(t, ok)
	require.Equal(t, "preventivo", encontrado.Tipo)

	actualizado, ok := a.ActualizarRegistro(r.ID, models.RegistroMantenimiento{EquipoID: eq.ID, Tipo: "correctivo"})
	require.True(t, ok)
	require.Equal(t, "correctivo", actualizado.Tipo)

	ok = a.BorrarRegistro(r.ID)
	require.True(t, ok)
}

func TestAlmacenSQLite_QuimicoCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	q, err := a.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro", StockActual: 10, UnidadMedida: "kg"})
	require.NoError(t, err)
	require.NotZero(t, q.ID)

	lista := a.ListarQuimicos()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarQuimicoPorID(q.ID)
	require.True(t, ok)
	require.Equal(t, "Cloro", encontrado.Nombre)

	actualizado, ok := a.ActualizarQuimico(q.ID, models.ProductoQuimico{Nombre: "Cloro Granulado"})
	require.True(t, ok)
	require.Equal(t, "Cloro Granulado", actualizado.Nombre)

	ok = a.BorrarQuimico(q.ID)
	require.True(t, ok)
}

func TestAlmacenSQLite_ClienteCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	c, err := a.CrearCliente(models.Cliente{Nombre: "Ana", Cedula: "1234567890"})
	require.NoError(t, err)
	require.NotZero(t, c.ID)

	lista := a.ListarClientes()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarClientePorID(c.ID)
	require.True(t, ok)
	require.Equal(t, "Ana", encontrado.Nombre)

	actualizado, ok := a.ActualizarCliente(c.ID, models.Cliente{Nombre: "Ana Editada", Cedula: "1234567890"})
	require.True(t, ok)
	require.Equal(t, "Ana Editada", actualizado.Nombre)

	ok = a.BorrarCliente(c.ID)
	require.True(t, ok)
}

func TestAlmacenSQLite_ReservaCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	cliente, _ := a.CrearCliente(models.Cliente{Nombre: "Luis", Cedula: "0000"})
	r, err := a.CrearReserva(models.Reserva{ClienteID: cliente.ID, Duracion: 720})
	require.NoError(t, err)
	require.NotZero(t, r.ID)

	lista := a.ListarReservas()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarReservaPorID(r.ID)
	require.True(t, ok)
	require.Equal(t, 720, encontrado.Duracion)

	actualizado, ok := a.ActualizarReserva(r.ID, models.Reserva{ClienteID: cliente.ID, Duracion: 1440})
	require.True(t, ok)
	require.Equal(t, 1440, actualizado.Duracion)

	ok = a.BorrarReserva(r.ID)
	require.True(t, ok)
}

func TestAlmacenSQLite_PagoCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	cliente, _ := a.CrearCliente(models.Cliente{Nombre: "Luis", Cedula: "0000"})
	p, err := a.CrearPago(models.Pago{ClienteID: cliente.ID, Monto: 5, Concepto: "dia"})
	require.NoError(t, err)
	require.NotZero(t, p.ID)

	lista := a.ListarPagos()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarPagoPorID(p.ID)
	require.True(t, ok)
	require.Equal(t, 5.0, encontrado.Monto)

	actualizado, ok := a.ActualizarPago(p.ID, models.Pago{ClienteID: cliente.ID, Monto: 10, Concepto: "medio_dia"})
	require.True(t, ok)
	require.Equal(t, 10.0, actualizado.Monto)

	tienePago := a.ClienteTienePagoEntrada(cliente.ID)
	require.True(t, tienePago)

	ok = a.BorrarPago(p.ID)
	require.True(t, ok)
}

func TestAlmacenSQLite_UsuarioCRUD(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	u, err := a.CrearUsuario(models.Usuario{Nombre: "Admin", Email: "admin@test.com", PasswordHash: "hash", Rol: "admin"})
	require.NoError(t, err)
	require.NotZero(t, u.ID)

	lista := a.ListarUsuarios()
	require.Len(t, lista, 1)

	encontrado, ok := a.BuscarUsuarioPorID(u.ID)
	require.True(t, ok)
	require.Equal(t, "Admin", encontrado.Nombre)

	encontradoEmail, ok := a.BuscarUsuarioPorEmail("admin@test.com")
	require.True(t, ok)
	require.Equal(t, u.ID, encontradoEmail.ID)

	actualizado, ok := a.ActualizarUsuario(u.ID, models.Usuario{Nombre: "Admin Editado", Email: "admin@test.com", PasswordHash: ""})
	require.True(t, ok)
	require.Equal(t, "Admin Editado", actualizado.Nombre)

	ok = a.BorrarUsuario(u.ID)
	require.True(t, ok)
}

func TestAlmacenSQLite_ActualizarUsuario_ConservaHash(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	u, _ := a.CrearUsuario(models.Usuario{Nombre: "Admin", Email: "admin@test.com", PasswordHash: "hashoriginal", Rol: "admin"})
	actualizado, ok := a.ActualizarUsuario(u.ID, models.Usuario{Nombre: "Editado", Email: "admin@test.com", PasswordHash: ""})
	require.True(t, ok)
	require.Equal(t, "hashoriginal", actualizado.PasswordHash)
}

func TestAlmacenSQLite_BuscarPorEmail_NoEncontrado(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	_, ok := a.BuscarUsuarioPorEmail("no@existe.com")
	require.False(t, ok)
}

func TestAlmacenSQLite_BorrarGuaradavia_NoEncontrado(t *testing.T) {
	db := abrirDBCompleta(t)
	a := NuevoAlmacenSQLite(db)

	ok := a.BorrarGuardavida(999)
	require.False(t, ok)
}
