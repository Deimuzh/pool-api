package service

import (
	"testing"
)

func TestEquipoService_New(t *testing.T) {
	repo := newMantenimientoRepoMock()
	svc := NewMantenimientoService(repo)
	if svc == nil {
		t.Fatal("NewMantenimientoService no debe devolver nil")
	}
}
