package models

import "time"

// Usuario representa una cuenta de acceso al panel de administración
// de la piscina (no confundir con Cliente, que es un socio/residente).
type Usuario struct {
	ID           uint      `json:"id"            gorm:"primaryKey"`
	Nombre       string    `json:"nombre"        gorm:"not null"`
	Email        string    `json:"email"         gorm:"not null;uniqueIndex"`
	PasswordHash string    `json:"-"             gorm:"not null"`        // nunca se serializa en JSON
	Rol          string    `json:"rol"           gorm:"default:'admin'"` // "admin", "guardavida", etc.
	CreadoEn     time.Time `json:"creado_en"`
}
