package models

import "time"

// ─── MÓDULO SEGURIDAD ────────────────────────────────────────────────────────

type Guardavida struct {
	ID          uint      `json:"id"           gorm:"primaryKey;autoIncrement"`
	Nombre      string    `json:"nombre"       gorm:"not null"`
	Turno       string    `json:"turno"        gorm:"not null"` // "mañana", "tarde", "noche"
	Certificado string    `json:"certificado"`
	Activo      bool      `json:"activo"       gorm:"default:true"`
	CreadoEn   time.Time `json:"creado_en"`
}

type Incidente struct {
	ID           uint      `json:"id"            gorm:"primaryKey;autoIncrement"`
	Tipo         string    `json:"tipo"          gorm:"not null"` // "ahogamiento", "lesion", "pelea"
	Descripcion  string    `json:"descripcion"`
	Gravedad     string    `json:"gravedad"      gorm:"not null"` // "leve", "moderado", "grave"
	GuardavidaID uint      `json:"guardavida_id" gorm:"not null"`
	FechaHora    time.Time `json:"fecha_hora"`
	Resuelto     bool      `json:"resuelto"      gorm:"default:false"`
}

type AccesoCliente struct {
	ID         uint      `json:"id"          gorm:"primaryKey;autoIncrement"`
	ClienteID  uint      `json:"cliente_id"  gorm:"not null"`
	FechaHora  time.Time `json:"fecha_hora"`
	Autorizado bool      `json:"autorizado"  gorm:"default:true"`
	Motivo     string    `json:"motivo"` // razón si fue denegado
}

// ─── MÓDULO MANTENIMIENTO ────────────────────────────────────────────────────

type Equipo struct {
	ID             uint      `json:"id"              gorm:"primaryKey;autoIncrement"`
	Nombre         string    `json:"nombre"          gorm:"not null"`
	Tipo           string    `json:"tipo"            gorm:"not null"` // "filtro", "bomba", "iluminacion"
	Estado         string    `json:"estado"          gorm:"default:'operativo'"` // "operativo", "en reparacion", "fuera de servicio"
	UltimaRevision time.Time `json:"ultima_revision"`
}

type RegistroMantenimiento struct {
	ID           uint      `json:"id"            gorm:"primaryKey;autoIncrement"`
	EquipoID     uint      `json:"equipo_id"     gorm:"not null"`
	Descripcion  string    `json:"descripcion"`
	Tipo         string    `json:"tipo"          gorm:"not null"` // "preventivo", "correctivo"
	FechaHora    time.Time `json:"fecha_hora"`
	RealizadoPor string    `json:"realizado_por"`
	Costo        float64   `json:"costo"`
}

type ProductoQuimico struct {
	ID           uint    `json:"id"            gorm:"primaryKey;autoIncrement"`
	Nombre       string  `json:"nombre"        gorm:"not null"`
	StockActual  float64 `json:"stock_actual"`
	UnidadMedida string  `json:"unidad_medida"` // "litros", "kg"
	NivelMinimo  float64 `json:"nivel_minimo"`
}

// ─── MÓDULO CLIENTES ─────────────────────────────────────────────────────────

type Cliente struct {
	ID             uint      `json:"id"              gorm:"primaryKey;autoIncrement"`
	Nombre         string    `json:"nombre"          gorm:"not null"`
	Cedula         string    `json:"cedula"          gorm:"uniqueIndex;not null"`
	Email          string    `json:"email"`
	Telefono       string    `json:"telefono"`
	Membresia      string    `json:"membresia"       gorm:"default:'ninguna'"` // "mensual", "trimestral", "anual", "ninguna"
	FechaRegistro  time.Time `json:"fecha_registro"`
}

type Reserva struct {
	ID        uint      `json:"id"         gorm:"primaryKey;autoIncrement"`
	ClienteID uint      `json:"cliente_id" gorm:"not null"`
	FechaHora time.Time `json:"fecha_hora"`
	Duracion  int       `json:"duracion"`  // en minutos
	Estado    string    `json:"estado"     gorm:"default:'pendiente'"` // "pendiente", "confirmada", "cancelada"
}

type Pago struct {
	ID        uint      `json:"id"         gorm:"primaryKey;autoIncrement"`
	ClienteID uint      `json:"cliente_id" gorm:"not null"`
	Monto     float64   `json:"monto"      gorm:"not null"`
	Concepto  string    `json:"concepto"`  // "membresia", "entrada", "reserva"
	FechaHora time.Time `json:"fecha_hora"`
	Metodo    string    `json:"metodo"`    // "efectivo", "transferencia"
}
