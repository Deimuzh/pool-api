package storage

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"pool-api/internal/models"
)

func abrirDBClientesPrueba(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir la base de datos en memoria: %v", err)
	}
	if err := db.AutoMigrate(&models.Cliente{}, &models.Reserva{}, &models.Pago{}); err != nil {
		t.Fatalf("fallo AutoMigrate: %v", err)
	}
	return db
}

func TestAlmacenSQLite_CrearYListarReserva(t *testing.T) {
	db := abrirDBClientesPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	cliente := models.Cliente{Nombre: "Maria Perez", Cedula: "1311111111"}
	creado, _ := almacen.CrearCliente(cliente)

	reserva := models.Reserva{
		ClienteID: creado.ID,
		FechaHora: time.Date(2026, 7, 10, 14, 0, 0, 0, time.UTC),
		Duracion:  60,
		Estado:    "pendiente",
	}
	rv, err := almacen.CrearReserva(reserva)
	if err != nil {
		t.Fatalf("CrearReserva fallo: %v", err)
	}
	if rv.ID == 0 {
		t.Fatal("se esperaba que GORM asignara un ID")
	}
	if rv.Estado != "pendiente" {
		t.Errorf("estado inesperado: %s", rv.Estado)
	}

	encontrada, ok := almacen.BuscarReservaPorID(rv.ID)
	if !ok {
		t.Fatal("no se encontro la reserva recien creada")
	}
	if encontrada.ClienteID != creado.ID {
		t.Errorf("ClienteID inesperado: %d", encontrada.ClienteID)
	}

	lista := almacen.ListarReservas()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 reserva, se obtuvieron %d", len(lista))
	}
}

func TestAlmacenSQLite_ActualizarYBorrarReserva(t *testing.T) {
	db := abrirDBClientesPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	cliente := models.Cliente{Nombre: "Luis Torres", Cedula: "1322222222"}
	creado, _ := almacen.CrearCliente(cliente)

	reserva := models.Reserva{
		ClienteID: creado.ID,
		Duracion:  30,
		Estado:    "pendiente",
	}
	rv, _ := almacen.CrearReserva(reserva)

	actualizada, ok := almacen.ActualizarReserva(rv.ID, models.Reserva{Estado: "confirmada"})
	if !ok {
		t.Fatal("ActualizarReserva devolvio false")
	}
	if actualizada.Estado != "confirmada" {
		t.Errorf("estado inesperado: %s", actualizada.Estado)
	}

	if !almacen.BorrarReserva(rv.ID) {
		t.Fatal("BorrarReserva devolvio false")
	}
	if almacen.BorrarReserva(999) {
		t.Error("BorrarReserva con ID inexistente devolvio true")
	}
}

func TestAlmacenSQLite_CrearYListarPago(t *testing.T) {
	db := abrirDBClientesPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	cliente := models.Cliente{Nombre: "Ana Ruiz", Cedula: "1333333333"}
	creado, _ := almacen.CrearCliente(cliente)

	pago := models.Pago{
		ClienteID: creado.ID,
		Monto:     50.0,
		Concepto:  "membresia",
		Metodo:    "efectivo",
	}
	p, err := almacen.CrearPago(pago)
	if err != nil {
		t.Fatalf("CrearPago fallo: %v", err)
	}
	if p.ID == 0 {
		t.Fatal("se esperaba que GORM asignara un ID")
	}
	if p.Monto != 50.0 {
		t.Errorf("monto inesperado: %f", p.Monto)
	}

	encontrado, ok := almacen.BuscarPagoPorID(p.ID)
	if !ok {
		t.Fatal("no se encontro el pago recien creado")
	}
	if encontrado.ClienteID != creado.ID {
		t.Errorf("ClienteID inesperado: %d", encontrado.ClienteID)
	}

	lista := almacen.ListarPagos()
	if len(lista) != 1 {
		t.Fatalf("se esperaba 1 pago, se obtuvieron %d", len(lista))
	}
}

func TestAlmacenSQLite_ActualizarYBorrarPago(t *testing.T) {
	db := abrirDBClientesPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	cliente := models.Cliente{Nombre: "Carlos Diaz", Cedula: "1344444444"}
	creado, _ := almacen.CrearCliente(cliente)

	pago := models.Pago{
		ClienteID: creado.ID,
		Monto:     100.0,
		Concepto:  "entrada",
		Metodo:    "transferencia",
	}
	p, _ := almacen.CrearPago(pago)

	actualizado, ok := almacen.ActualizarPago(p.ID, models.Pago{Monto: 120.0})
	if !ok {
		t.Fatal("ActualizarPago devolvio false")
	}
	if actualizado.Monto != 120.0 {
		t.Errorf("monto inesperado: %f", actualizado.Monto)
	}

	if !almacen.BorrarPago(p.ID) {
		t.Fatal("BorrarPago devolvio false")
	}
	if almacen.BorrarPago(999) {
		t.Error("BorrarPago con ID inexistente devolvio true")
	}
}

func TestAlmacenSQLite_ClienteTienePagoEntrada(t *testing.T) {
	db := abrirDBClientesPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	cliente := models.Cliente{Nombre: "Sofia Vega", Cedula: "1355555555"}
	creado, _ := almacen.CrearCliente(cliente)

	if almacen.ClienteTienePagoEntrada(creado.ID) {
		t.Error("no deberia tener pago de entrada sin crear ninguno")
	}

	almacen.CrearPago(models.Pago{
		ClienteID: creado.ID,
		Monto:     5.0,
		Concepto:  "dia",
		Metodo:    "efectivo",
	})
	if !almacen.ClienteTienePagoEntrada(creado.ID) {
		t.Error("deberia tener pago de entrada despues de crearlo")
	}
}

func TestAlmacenSQLite_BuscarReservaPagoInexistente(t *testing.T) {
	db := abrirDBClientesPrueba(t)
	almacen := NuevoAlmacenSQLite(db)

	_, ok := almacen.BuscarReservaPorID(999)
	if ok {
		t.Error("no se esperaba encontrar reserva con ID 999")
	}

	_, ok = almacen.BuscarPagoPorID(999)
	if ok {
		t.Error("no se esperaba encontrar pago con ID 999")
	}
}
