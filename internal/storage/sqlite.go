package storage

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

// AlmacenSQLite implementa la interfaz Almacen usando GORM sobre SQLite.
//
// Los métodos tienen exactamente las firmas que pide cada Repository en
// almacen.go. Por eso los services no se enteran de qué backend reciben.
type AlmacenSQLite struct {
	db *gorm.DB
}

// NuevoAlmacenSQLite envuelve una conexión *gorm.DB ya abierta.
func NuevoAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}

// =========================================================
// GUARDAVIDAS
// =========================================================

func (a *AlmacenSQLite) ListarGuardavidas() []models.Guardavida {
	var lista []models.Guardavida
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarGuardavidaPorID(id int) (models.Guardavida, bool) {
	var g models.Guardavida
	if err := a.db.First(&g, id).Error; err != nil {
		return models.Guardavida{}, false
	}
	return g, true
}

func (a *AlmacenSQLite) CrearGuardavida(g models.Guardavida) models.Guardavida {
	g.CreadoEn = time.Now()
	a.db.Create(&g)
	return g
}

func (a *AlmacenSQLite) ActualizarGuardavida(id int, datos models.Guardavida) (models.Guardavida, bool) {
	var existente models.Guardavida
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Guardavida{}, false
	}
	datos.ID = id
	datos.CreadoEn = existente.CreadoEn
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarGuardavida(id int) bool {
	res := a.db.Delete(&models.Guardavida{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// INCIDENTES
// =========================================================

func (a *AlmacenSQLite) ListarIncidentes() []models.Incidente {
	var lista []models.Incidente
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarIncidentePorID(id int) (models.Incidente, bool) {
	var i models.Incidente
	if err := a.db.First(&i, id).Error; err != nil {
		return models.Incidente{}, false
	}
	return i, true
}

func (a *AlmacenSQLite) CrearIncidente(i models.Incidente) models.Incidente {
	i.FechaHora = time.Now()
	a.db.Create(&i)
	return i
}

func (a *AlmacenSQLite) ActualizarIncidente(id int, datos models.Incidente) (models.Incidente, bool) {
	var existente models.Incidente
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Incidente{}, false
	}
	datos.ID = id
	datos.FechaHora = existente.FechaHora
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarIncidente(id int) bool {
	res := a.db.Delete(&models.Incidente{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// ACCESOS
// =========================================================

func (a *AlmacenSQLite) ListarAccesos() []models.AccesoCliente {
	var lista []models.AccesoCliente
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarAccesoPorID(id int) (models.AccesoCliente, bool) {
	var acc models.AccesoCliente
	if err := a.db.First(&acc, id).Error; err != nil {
		return models.AccesoCliente{}, false
	}
	return acc, true
}

func (a *AlmacenSQLite) CrearAcceso(acc models.AccesoCliente) models.AccesoCliente {
	acc.FechaHora = time.Now()
	a.db.Create(&acc)
	return acc
}

func (a *AlmacenSQLite) ActualizarAcceso(id int, datos models.AccesoCliente) (models.AccesoCliente, bool) {
	var existente models.AccesoCliente
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.AccesoCliente{}, false
	}
	datos.ID = id
	datos.FechaHora = existente.FechaHora
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarAcceso(id int) bool {
	res := a.db.Delete(&models.AccesoCliente{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// EQUIPOS
// =========================================================

func (a *AlmacenSQLite) ListarEquipos() []models.Equipo {
	var lista []models.Equipo
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarEquipoPorID(id int) (models.Equipo, bool) {
	var e models.Equipo
	if err := a.db.First(&e, id).Error; err != nil {
		return models.Equipo{}, false
	}
	return e, true
}

func (a *AlmacenSQLite) CrearEquipo(e models.Equipo) models.Equipo {
	a.db.Create(&e)
	return e
}

func (a *AlmacenSQLite) ActualizarEquipo(id int, datos models.Equipo) (models.Equipo, bool) {
	var existente models.Equipo
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Equipo{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarEquipo(id int) bool {
	res := a.db.Delete(&models.Equipo{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// REGISTROS DE MANTENIMIENTO
// =========================================================

func (a *AlmacenSQLite) ListarRegistros() []models.RegistroMantenimiento {
	var lista []models.RegistroMantenimiento
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarRegistroPorID(id int) (models.RegistroMantenimiento, bool) {
	var rm models.RegistroMantenimiento
	if err := a.db.First(&rm, id).Error; err != nil {
		return models.RegistroMantenimiento{}, false
	}
	return rm, true
}

func (a *AlmacenSQLite) CrearRegistro(rm models.RegistroMantenimiento) models.RegistroMantenimiento {
	rm.FechaHora = time.Now()
	a.db.Create(&rm)
	return rm
}

func (a *AlmacenSQLite) ActualizarRegistro(id int, datos models.RegistroMantenimiento) (models.RegistroMantenimiento, bool) {
	var existente models.RegistroMantenimiento
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.RegistroMantenimiento{}, false
	}
	datos.ID = id
	datos.FechaHora = existente.FechaHora
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarRegistro(id int) bool {
	res := a.db.Delete(&models.RegistroMantenimiento{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// PRODUCTOS QUÍMICOS
// =========================================================

func (a *AlmacenSQLite) ListarQuimicos() []models.ProductoQuimico {
	var lista []models.ProductoQuimico
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarQuimicoPorID(id int) (models.ProductoQuimico, bool) {
	var q models.ProductoQuimico
	if err := a.db.First(&q, id).Error; err != nil {
		return models.ProductoQuimico{}, false
	}
	return q, true
}

func (a *AlmacenSQLite) CrearQuimico(q models.ProductoQuimico) models.ProductoQuimico {
	a.db.Create(&q)
	return q
}

func (a *AlmacenSQLite) ActualizarQuimico(id int, datos models.ProductoQuimico) (models.ProductoQuimico, bool) {
	var existente models.ProductoQuimico
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.ProductoQuimico{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarQuimico(id int) bool {
	res := a.db.Delete(&models.ProductoQuimico{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// CLIENTES
// =========================================================

func (a *AlmacenSQLite) ListarClientes() []models.Cliente {
	var lista []models.Cliente
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarClientePorID(id int) (models.Cliente, bool) {
	var c models.Cliente
	if err := a.db.First(&c, id).Error; err != nil {
		return models.Cliente{}, false
	}
	return c, true
}

func (a *AlmacenSQLite) CrearCliente(c models.Cliente) models.Cliente {
	c.FechaRegistro = time.Now()
	a.db.Create(&c)
	return c
}

func (a *AlmacenSQLite) ActualizarCliente(id int, datos models.Cliente) (models.Cliente, bool) {
	var existente models.Cliente
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Cliente{}, false
	}
	datos.ID = id
	datos.FechaRegistro = existente.FechaRegistro
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarCliente(id int) bool {
	res := a.db.Delete(&models.Cliente{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// RESERVAS
// =========================================================

func (a *AlmacenSQLite) ListarReservas() []models.Reserva {
	var lista []models.Reserva
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarReservaPorID(id int) (models.Reserva, bool) {
	var rv models.Reserva
	if err := a.db.First(&rv, id).Error; err != nil {
		return models.Reserva{}, false
	}
	return rv, true
}

func (a *AlmacenSQLite) CrearReserva(rv models.Reserva) models.Reserva {
	rv.FechaHora = time.Now()
	a.db.Create(&rv)
	return rv
}

func (a *AlmacenSQLite) ActualizarReserva(id int, datos models.Reserva) (models.Reserva, bool) {
	var existente models.Reserva
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Reserva{}, false
	}
	datos.ID = id
	datos.FechaHora = existente.FechaHora
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarReserva(id int) bool {
	res := a.db.Delete(&models.Reserva{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// PAGOS
// =========================================================

func (a *AlmacenSQLite) ListarPagos() []models.Pago {
	var lista []models.Pago
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarPagoPorID(id int) (models.Pago, bool) {
	var p models.Pago
	if err := a.db.First(&p, id).Error; err != nil {
		return models.Pago{}, false
	}
	return p, true
}

func (a *AlmacenSQLite) CrearPago(p models.Pago) models.Pago {
	p.FechaHora = time.Now()
	a.db.Create(&p)
	return p
}

func (a *AlmacenSQLite) ActualizarPago(id int, datos models.Pago) (models.Pago, bool) {
	var existente models.Pago
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Pago{}, false
	}
	datos.ID = id
	datos.FechaHora = existente.FechaHora
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarPago(id int) bool {
	res := a.db.Delete(&models.Pago{}, id)
	return res.RowsAffected > 0
}

// ClienteTienePagoEntrada verifica si el cliente tiene al menos un pago
// con concepto "entrada" registrado. Lo usa el service de Seguridad para
// decidir si autoriza el acceso, sin acoplarse a la tabla Pago.
func (a *AlmacenSQLite) ClienteTienePagoEntrada(clienteID int) bool {
	var count int64
	a.db.Model(&models.Pago{}).
		Where("cliente_id = ? AND concepto = ?", clienteID, "entrada").
		Count(&count)
	return count > 0
}

// =========================================================
// USUARIOS (autenticación)
// =========================================================

func (a *AlmacenSQLite) ListarUsuarios() []models.Usuario {
	var lista []models.Usuario
	a.db.Find(&lista)
	return lista
}

func (a *AlmacenSQLite) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	var u models.Usuario
	if err := a.db.First(&u, id).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

func (a *AlmacenSQLite) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := a.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

func (a *AlmacenSQLite) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.CreadoEn = time.Now()
	if err := a.db.Create(&u).Error; err != nil {
		return models.Usuario{}, err
	}
	return u, nil
}

func (a *AlmacenSQLite) ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool) {
	var existente models.Usuario
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Usuario{}, false
	}
	datos.ID = id
	datos.CreadoEn = existente.CreadoEn
	// Si no se manda un nuevo hash, conservar el actual.
	if datos.PasswordHash == "" {
		datos.PasswordHash = existente.PasswordHash
	}
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarUsuario(id int) bool {
	res := a.db.Delete(&models.Usuario{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// SEEDS
// =========================================================

// SembrarSiVacio inserta datos iniciales solo si aún no hay clientes,
// y crea el usuario administrador por defecto si no existe ninguno.
func (a *AlmacenSQLite) SembrarSiVacio() {
	a.sembrarAdminPorDefecto()

	var n int64
	a.db.Model(&models.Cliente{}).Count(&n)
	if n > 0 {
		return
	}

	clientes := []models.Cliente{
		{ID: 1, Nombre: "Ana Reyes", Cedula: "1301234567", Email: "ana@correo.com", Membresia: "mensual"},
		{ID: 2, Nombre: "Luis Pino", Cedula: "1311234567", Email: "luis@correo.com", Membresia: "ninguna"},
		{ID: 3, Nombre: "Rosa Ávila", Cedula: "1321234567", Email: "rosa@correo.com", Membresia: "anual"},
	}
	a.db.Create(&clientes)

	guardavidas := []models.Guardavida{
		{ID: 1, Nombre: "Carlos Mendoza", Turno: "mañana", Certificado: "Cruz Roja Niv. 2", Activo: true},
		{ID: 2, Nombre: "María Suárez", Turno: "tarde", Certificado: "Cruz Roja Niv. 1", Activo: true},
	}
	a.db.Create(&guardavidas)

	equipos := []models.Equipo{
		{ID: 1, Nombre: "Bomba filtro A", Tipo: "bomba", Estado: "operativo"},
		{ID: 2, Nombre: "Clorador automático", Tipo: "quimico", Estado: "operativo"},
	}
	a.db.Create(&equipos)

	quimicos := []models.ProductoQuimico{
		{ID: 1, Nombre: "Cloro granulado", StockActual: 25, UnidadMedida: "kg", NivelMinimo: 5},
		{ID: 2, Nombre: "Ácido muriático", StockActual: 4, UnidadMedida: "litros", NivelMinimo: 5},
	}
	a.db.Create(&quimicos)

	pagos := []models.Pago{
		{ID: 1, ClienteID: 1, Monto: 5.0, Concepto: "entrada", Metodo: "efectivo"},
	}
	a.db.Create(&pagos)
}

// sembrarAdminPorDefecto crea el usuario admin@piscina.com / admin123
// la primera vez que se levanta el servidor, si no existe ningún usuario.
// Cámbiala apenas puedas desde el CRUD de usuarios.
func (a *AlmacenSQLite) sembrarAdminPorDefecto() {
	var n int64
	a.db.Model(&models.Usuario{}).Count(&n)
	if n > 0 {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	a.db.Create(&models.Usuario{
		Nombre:       "Administrador",
		Email:        "admin@piscina.com",
		PasswordHash: string(hash),
		Rol:          "admin",
		CreadoEn:     time.Now(),
	})
}

// Chequeo en tiempo de compilación: AlmacenSQLite debe cumplir Almacen.
var _ Almacen = (*AlmacenSQLite)(nil)
