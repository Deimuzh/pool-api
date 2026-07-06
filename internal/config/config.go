package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config centraliza la configuración de arranque de la API.
type Config struct {
	Puerto       string
	DBDriver     string
	DBDSN        string
	RutaDB       string
	JWTSecreto   []byte
	JWTDuracion  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Cargar lee variables de entorno y, si existe, un archivo .env local.
func Cargar() Config {
	_ = godotenv.Load()

	return Config{
		Puerto:       conTexto("PUERTO", ":8080"),
		DBDriver:     strings.ToLower(conTexto("DB_DRIVER", "sqlite")),
		DBDSN:        conTexto("DB_DSN", ""),
		RutaDB:       conTexto("RUTA_DB", "piscina.db"),
		JWTSecreto:   []byte(conTexto("JWT_SECRETO", "piscina-los-ceibos-clave-secreta")),
		JWTDuracion:  conDuracion("JWT_DURACION", 24*time.Hour),
		ReadTimeout:  conDuracion("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: conDuracion("HTTP_WRITE_TIMEOUT", 10*time.Second),
	}
}

func conTexto(clave, porDefecto string) string {
	if valor := os.Getenv(clave); valor != "" {
		return valor
	}
	return porDefecto
}

func conDuracion(clave string, porDefecto time.Duration) time.Duration {
	valor := os.Getenv(clave)
	if valor == "" {
		return porDefecto
	}
	duracion, err := time.ParseDuration(valor)
	if err != nil {
		return porDefecto
	}
	return duracion
}
