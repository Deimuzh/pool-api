package config

import (
	"os"
	"testing"
	"time"
)

func TestCargar_ValoresPorDefecto(t *testing.T) {
	cfg := Cargar()
	if cfg.Puerto != ":8080" {
		t.Errorf("Puerto por defecto esperado :8080, obtuve %s", cfg.Puerto)
	}
	if cfg.DBDriver != "sqlite" {
		t.Errorf("DBDriver por defecto esperado sqlite, obtuve %s", cfg.DBDriver)
	}
	if cfg.JWTDuracion != 24*time.Hour {
		t.Errorf("JWTDuracion por defecto esperada 24h, obtuve %v", cfg.JWTDuracion)
	}
}

func TestCargar_VariablesEntorno(t *testing.T) {
	os.Setenv("PUERTO", ":9090")
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("JWT_DURACION", "2h")
	defer func() {
		os.Unsetenv("PUERTO")
		os.Unsetenv("DB_DRIVER")
		os.Unsetenv("JWT_DURACION")
	}()
	cfg := Cargar()
	if cfg.Puerto != ":9090" {
		t.Errorf("Puerto esperado :9090, obtuve %s", cfg.Puerto)
	}
	if cfg.DBDriver != "postgres" {
		t.Errorf("DBDriver esperado postgres, obtuve %s", cfg.DBDriver)
	}
	if cfg.JWTDuracion != 2*time.Hour {
		t.Errorf("JWTDuracion esperada 2h, obtuve %v", cfg.JWTDuracion)
	}
}

func TestCargar_DuracionInvalidaUsaDefault(t *testing.T) {
	os.Setenv("JWT_DURACION", "no-valida")
	defer os.Unsetenv("JWT_DURACION")
	cfg := Cargar()
	if cfg.JWTDuracion != 24*time.Hour {
		t.Errorf("duracion invalida deberia usar default 24h, obtuve %v", cfg.JWTDuracion)
	}
}
