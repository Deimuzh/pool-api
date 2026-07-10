// Package middleware contiene middlewares HTTP reutilizables por el router.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"pool-api/internal/service"
)

// ContextKey evita colisiones de claves al guardar valores en el context.
type ContextKey string

const (
	ContextKeyUsuarioID ContextKey = "usuarioID"
	ContextKeyRol       ContextKey = "rol"
)

// Auth devuelve un middleware que exige un JWT válido en el header
// "Authorization: Bearer <token>". Si falta o es inválido, responde
// 401 y corta la cadena sin llegar al handler protegido.
//
// Uso en main.go:
//
//	r.Group(func(r chi.Router) {
//	    r.Use(middleware.Auth(authSvc))
//	    r.Route("/guardavidas", ...)
//	})
func Auth(auth *service.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encabezado := r.Header.Get("Authorization")
			partes := strings.SplitN(encabezado, " ", 2)
			if len(partes) != 2 || partes[0] != "Bearer" {
				responderNoAutorizado(w)
				return
			}

			claims, err := auth.ValidarToken(partes[1])
			if err != nil {
				responderNoAutorizado(w)
				return
			}
// Guarda usuarioID y rol en el context para handlers posteriores
			ctx := context.WithValue(r.Context(), ContextKeyUsuarioID, claims.UsuarioID)
			ctx = context.WithValue(ctx, ContextKeyRol, claims.Rol)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SoloRoles permite proteger rutas por rol despues de validar el JWT.
// Debe usarse dentro de un grupo que ya tenga middleware.Auth(authSvc).
func SoloRoles(rolesPermitidos ...string) func(next http.Handler) http.Handler {
	permitidos := make(map[string]struct{}, len(rolesPermitidos))
	for _, rol := range rolesPermitidos {
		permitidos[rol] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rol, _ := r.Context().Value(ContextKeyRol).(string)
			if _, ok := permitidos[rol]; !ok {
				responderProhibido(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func responderNoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"token faltante o inválido, inicia sesión de nuevo"}`))
}

func responderProhibido(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":"no tienes permisos para esta ruta"}`))
}
