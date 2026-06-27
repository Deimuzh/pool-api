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

const ContextKeyUsuarioID ContextKey = "usuarioID"

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

			ctx := context.WithValue(r.Context(), ContextKeyUsuarioID, claims.UsuarioID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func responderNoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"token faltante o inválido, inicia sesión de nuevo"}`))
}
