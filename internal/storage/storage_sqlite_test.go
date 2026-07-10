package storage

import (
	"testing"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

// abrirDBCompletaPrueba crea una base SQLite en memoria con todos los modelos
// migrados. Así cada test prueba el repositorio GORM real sin compartir estado
// con otros tests ni depender de la base local del proyecto.
func abrirDBCompletaPrueba(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite en memoria: %v", err)
	}
	if err := db.AutoMigrate(
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
	); err != nil {
		t.Fatalf("fallo AutoMigrate: %v", err)
	}
	return db
}

// TestAlmacenSQLite_ClientesCRUD verifica el ciclo completo de clientes en el
// repositorio real: crear, listar, buscar, actualizar, borrar y manejar IDs que
// no existen devolviendo ok=false/false.
func TestAlmacenSQLite_ClientesCRUD(t *testing.T) {
	repo := NuevoAlmacenSQLite(abrirDBCompletaPrueba(t))

	creado, err := repo.CrearCliente(models.Cliente{
		Nombre:    "Cliente Demo",
		Cedula:    "1310000001",
		Email:     "cliente@example.com",
		Telefono:  "0999999999",
		Membresia: "ninguna",
	})
	if err != nil {
		t.Fatalf("no se esperaba error al crear cliente: %v", err)
	}
	if creado.ID == 0 || creado.FechaRegistro.IsZero() {
		t.Fatalf("cliente creado inesperado: %+v", creado)
	}

	lista := repo.ListarClientes()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 cliente, se obtuvieron %d", len(lista))
	}

	buscado, ok := repo.BuscarClientePorID(creado.ID)
	if !ok || buscado.Cedula != "1310000001" {
		t.Fatalf("cliente buscado inesperado: %+v ok=%v", buscado, ok)
	}
	if _, ok := repo.BuscarClientePorID(999); ok {
		t.Fatal("no se esperaba encontrar cliente inexistente")
	}

	actualizado, ok := repo.ActualizarCliente(creado.ID, models.Cliente{
		Nombre:    "Cliente Actualizado",
		Cedula:    "1310000001",
		Membresia: "mensual",
	})
	if !ok || actualizado.Nombre != "Cliente Actualizado" || actualizado.FechaRegistro.IsZero() {
		t.Fatalf("cliente actualizado inesperado: %+v ok=%v", actualizado, ok)
	}
	if _, ok := repo.ActualizarCliente(999, models.Cliente{Nombre: "No existe"}); ok {
		t.Fatal("no se esperaba actualizar cliente inexistente")
	}

	if !repo.BorrarCliente(creado.ID) {
		t.Fatal("se esperaba borrar cliente")
	}
	if repo.BorrarCliente(999) {
		t.Fatal("no se esperaba borrar cliente inexistente")
	}
}

// TestAlmacenSQLite_ReservasCRUD prueba que las reservas se persistan con
// FechaHora automática y que los caminos de actualización/borrado inexistentes
// se reporten correctamente.
func TestAlmacenSQLite_ReservasCRUD(t *testing.T) {
	repo := NuevoAlmacenSQLite(abrirDBCompletaPrueba(t))

	creada, err := repo.CrearReserva(models.Reserva{ClienteID: 1, Duracion: 60, Estado: "pendiente"})
	if err != nil {
		t.Fatalf("no se esperaba error al crear reserva: %v", err)
	}
	if creada.ID == 0 || creada.FechaHora.IsZero() {
		t.Fatalf("reserva creada inesperada: %+v", creada)
	}
	if len(repo.ListarReservas()) != 1 {
		t.Fatal("se esperaba listar 1 reserva")
	}

	buscada, ok := repo.BuscarReservaPorID(creada.ID)
	if !ok || buscada.ClienteID != 1 {
		t.Fatalf("reserva buscada inesperada: %+v ok=%v", buscada, ok)
	}
	if _, ok := repo.BuscarReservaPorID(999); ok {
		t.Fatal("no se esperaba encontrar reserva inexistente")
	}

	actualizada, ok := repo.ActualizarReserva(creada.ID, models.Reserva{ClienteID: 1, Duracion: 120, Estado: "confirmada"})
	if !ok || actualizada.Estado != "confirmada" || actualizada.FechaHora.IsZero() {
		t.Fatalf("reserva actualizada inesperada: %+v ok=%v", actualizada, ok)
	}
	if _, ok := repo.ActualizarReserva(999, models.Reserva{ClienteID: 1}); ok {
		t.Fatal("no se esperaba actualizar reserva inexistente")
	}

	if !repo.BorrarReserva(creada.ID) {
		t.Fatal("se esperaba borrar reserva")
	}
	if repo.BorrarReserva(999) {
		t.Fatal("no se esperaba borrar reserva inexistente")
	}
}

// TestAlmacenSQLite_PagosCRUDYEntrada cubre pagos y la consulta transversal que
// usa Seguridad para saber si un cliente tiene pago de entrada registrado.
func TestAlmacenSQLite_PagosCRUDYEntrada(t *testing.T) {
	repo := NuevoAlmacenSQLite(abrirDBCompletaPrueba(t))

	creado, err := repo.CrearPago(models.Pago{ClienteID: 7, Monto: 3.5, Concepto: "medio_dia", Metodo: "efectivo"})
	if err != nil {
		t.Fatalf("no se esperaba error al crear pago: %v", err)
	}
	if creado.ID == 0 || creado.FechaHora.IsZero() {
		t.Fatalf("pago creado inesperado: %+v", creado)
	}
	if len(repo.ListarPagos()) != 1 {
		t.Fatal("se esperaba listar 1 pago")
	}
	if !repo.ClienteTienePagoEntrada(7) {
		t.Fatal("se esperaba detectar pago de entrada")
	}
	if repo.ClienteTienePagoEntrada(99) {
		t.Fatal("no se esperaba detectar pago para otro cliente")
	}

	buscado, ok := repo.BuscarPagoPorID(creado.ID)
	if !ok || buscado.Monto != 3.5 {
		t.Fatalf("pago buscado inesperado: %+v ok=%v", buscado, ok)
	}
	if _, ok := repo.BuscarPagoPorID(999); ok {
		t.Fatal("no se esperaba encontrar pago inexistente")
	}

	actualizado, ok := repo.ActualizarPago(creado.ID, models.Pago{ClienteID: 7, Monto: 5, Concepto: "dia", Metodo: "transferencia"})
	if !ok || actualizado.Concepto != "dia" || actualizado.FechaHora.IsZero() {
		t.Fatalf("pago actualizado inesperado: %+v ok=%v", actualizado, ok)
	}
	if _, ok := repo.ActualizarPago(999, models.Pago{ClienteID: 7}); ok {
		t.Fatal("no se esperaba actualizar pago inexistente")
	}

	if !repo.BorrarPago(creado.ID) {
		t.Fatal("se esperaba borrar pago")
	}
	if repo.BorrarPago(999) {
		t.Fatal("no se esperaba borrar pago inexistente")
	}
}

// TestAlmacenSQLite_UsuariosCRUDYBuscarEmail valida el repositorio de usuarios:
// búsqueda por ID/email, conservación del hash al actualizar sin password nuevo
// y eliminación por ID.
func TestAlmacenSQLite_UsuariosCRUDYBuscarEmail(t *testing.T) {
	repo := NuevoAlmacenSQLite(abrirDBCompletaPrueba(t))

	creado, err := repo.CrearUsuario(models.Usuario{
		Nombre:       "Admin Demo",
		Email:        "admin.demo@piscina.com",
		PasswordHash: "hash-inicial",
		Rol:          "admin",
	})
	if err != nil {
		t.Fatalf("no se esperaba error al crear usuario: %v", err)
	}
	if creado.ID == 0 || creado.CreadoEn.IsZero() {
		t.Fatalf("usuario creado inesperado: %+v", creado)
	}
	if len(repo.ListarUsuarios()) != 1 {
		t.Fatal("se esperaba listar 1 usuario")
	}

	porID, ok := repo.BuscarUsuarioPorID(creado.ID)
	if !ok || porID.Email != "admin.demo@piscina.com" {
		t.Fatalf("usuario por id inesperado: %+v ok=%v", porID, ok)
	}
	porEmail, ok := repo.BuscarUsuarioPorEmail("admin.demo@piscina.com")
	if !ok || porEmail.ID != creado.ID {
		t.Fatalf("usuario por email inesperado: %+v ok=%v", porEmail, ok)
	}
	if _, ok := repo.BuscarUsuarioPorID(999); ok {
		t.Fatal("no se esperaba encontrar usuario inexistente por id")
	}
	if _, ok := repo.BuscarUsuarioPorEmail("nadie@example.com"); ok {
		t.Fatal("no se esperaba encontrar usuario inexistente por email")
	}

	actualizado, ok := repo.ActualizarUsuario(creado.ID, models.Usuario{
		Nombre: "Admin Actualizado",
		Email:  "admin.demo@piscina.com",
		Rol:    "admin",
	})
	if !ok || actualizado.Nombre != "Admin Actualizado" || actualizado.PasswordHash != "hash-inicial" {
		t.Fatalf("usuario actualizado inesperado: %+v ok=%v", actualizado, ok)
	}
	if _, ok := repo.ActualizarUsuario(999, models.Usuario{Nombre: "No existe"}); ok {
		t.Fatal("no se esperaba actualizar usuario inexistente")
	}

	if !repo.BorrarUsuario(creado.ID) {
		t.Fatal("se esperaba borrar usuario")
	}
	if repo.BorrarUsuario(999) {
		t.Fatal("no se esperaba borrar usuario inexistente")
	}
}

// TestAlmacenSQLite_RegistrosYQuimicosCRUD cubre los repositorios de
// mantenimiento que estaban sin cobertura: registros y productos químicos.
func TestAlmacenSQLite_RegistrosYQuimicosCRUD(t *testing.T) {
	repo := NuevoAlmacenSQLite(abrirDBCompletaPrueba(t))

	registro, err := repo.CrearRegistro(models.RegistroMantenimiento{
		EquipoID:     1,
		Descripcion:  "Revision semanal",
		Tipo:         "preventivo",
		RealizadoPor: "Tecnico Demo",
		Costo:        25,
	})
	if err != nil {
		t.Fatalf("no se esperaba error al crear registro: %v", err)
	}
	if registro.ID == 0 || registro.FechaHora.IsZero() {
		t.Fatalf("registro creado inesperado: %+v", registro)
	}
	if len(repo.ListarRegistros()) != 1 {
		t.Fatal("se esperaba listar 1 registro")
	}
	if _, ok := repo.BuscarRegistroPorID(999); ok {
		t.Fatal("no se esperaba encontrar registro inexistente")
	}
	registroActualizado, ok := repo.ActualizarRegistro(registro.ID, models.RegistroMantenimiento{EquipoID: 1, Descripcion: "Cambio de filtro", Tipo: "correctivo"})
	if !ok || registroActualizado.Tipo != "correctivo" || registroActualizado.FechaHora.IsZero() {
		t.Fatalf("registro actualizado inesperado: %+v ok=%v", registroActualizado, ok)
	}
	if _, ok := repo.ActualizarRegistro(999, models.RegistroMantenimiento{EquipoID: 1}); ok {
		t.Fatal("no se esperaba actualizar registro inexistente")
	}
	if !repo.BorrarRegistro(registro.ID) || repo.BorrarRegistro(999) {
		t.Fatal("resultado inesperado al borrar registro")
	}

	quimico, err := repo.CrearQuimico(models.ProductoQuimico{Nombre: "Cloro", StockActual: 10, UnidadMedida: "kg", NivelMinimo: 3})
	if err != nil {
		t.Fatalf("no se esperaba error al crear quimico: %v", err)
	}
	if len(repo.ListarQuimicos()) != 1 {
		t.Fatal("se esperaba listar 1 quimico")
	}
	if _, ok := repo.BuscarQuimicoPorID(quimico.ID); !ok {
		t.Fatal("se esperaba encontrar quimico")
	}
	if _, ok := repo.BuscarQuimicoPorID(999); ok {
		t.Fatal("no se esperaba encontrar quimico inexistente")
	}
	quimicoActualizado, ok := repo.ActualizarQuimico(quimico.ID, models.ProductoQuimico{Nombre: "Cloro", StockActual: 20, UnidadMedida: "kg", NivelMinimo: 4})
	if !ok || quimicoActualizado.StockActual != 20 {
		t.Fatalf("quimico actualizado inesperado: %+v ok=%v", quimicoActualizado, ok)
	}
	if _, ok := repo.ActualizarQuimico(999, models.ProductoQuimico{Nombre: "No existe"}); ok {
		t.Fatal("no se esperaba actualizar quimico inexistente")
	}
	if !repo.BorrarQuimico(quimico.ID) || repo.BorrarQuimico(999) {
		t.Fatal("resultado inesperado al borrar quimico")
	}
}

// TestAlmacenSQLite_EquiposActualizarYBorrar completa la cobertura de equipos
// verificando actualización, eliminación y sus caminos para ID inexistente.
func TestAlmacenSQLite_EquiposActualizarYBorrar(t *testing.T) {
	repo := NuevoAlmacenSQLite(abrirDBCompletaPrueba(t))

	creado, err := repo.CrearEquipo(models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "operativo"})
	if err != nil {
		t.Fatalf("no se esperaba error al crear equipo: %v", err)
	}
	actualizado, ok := repo.ActualizarEquipo(creado.ID, models.Equipo{Nombre: "Bomba", Tipo: "bomba", Estado: "en reparacion"})
	if !ok || actualizado.Estado != "en reparacion" {
		t.Fatalf("equipo actualizado inesperado: %+v ok=%v", actualizado, ok)
	}
	if _, ok := repo.ActualizarEquipo(999, models.Equipo{Nombre: "No existe"}); ok {
		t.Fatal("no se esperaba actualizar equipo inexistente")
	}
	if !repo.BorrarEquipo(creado.ID) {
		t.Fatal("se esperaba borrar equipo")
	}
	if repo.BorrarEquipo(999) {
		t.Fatal("no se esperaba borrar equipo inexistente")
	}
}

// TestAlmacenSQLite_SembrarSiVacio prueba los seeders del repositorio: deben
// crear datos iniciales una sola vez y generar el admin por defecto con hash
// bcrypt válido para la contraseña documentada.
func TestAlmacenSQLite_SembrarSiVacio(t *testing.T) {
	repo := NuevoAlmacenSQLite(abrirDBCompletaPrueba(t))

	repo.SembrarSiVacio()
	repo.SembrarSiVacio()

	if got := len(repo.ListarClientes()); got != 3 {
		t.Fatalf("se esperaban 3 clientes sembrados, se obtuvieron %d", got)
	}
	if got := len(repo.ListarGuardavidas()); got != 2 {
		t.Fatalf("se esperaban 2 guardavidas sembrados, se obtuvieron %d", got)
	}
	if got := len(repo.ListarEquipos()); got != 2 {
		t.Fatalf("se esperaban 2 equipos sembrados, se obtuvieron %d", got)
	}
	if got := len(repo.ListarQuimicos()); got != 2 {
		t.Fatalf("se esperaban 2 quimicos sembrados, se obtuvieron %d", got)
	}
	if got := len(repo.ListarPagos()); got != 1 {
		t.Fatalf("se esperaba 1 pago sembrado, se obtuvieron %d", got)
	}

	admin, ok := repo.BuscarUsuarioPorEmail("admin@piscina.com")
	if !ok {
		t.Fatal("se esperaba usuario admin sembrado")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte("admin123")); err != nil {
		t.Fatalf("password admin sembrado inesperado: %v", err)
	}
	if got := len(repo.ListarUsuarios()); got != 1 {
		t.Fatalf("se esperaba que SembrarSiVacio no duplique admin, usuarios=%d", got)
	}
}
