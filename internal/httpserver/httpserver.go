package httpserver

import (
	"net/http"
	"time"
)

// Nuevo construye el servidor HTTP con timeouts configurables.
// Esto evita usar http.ListenAndServe directo y permite un apagado ordenado.
func Nuevo(handler http.Handler, puerto string, readTimeout, writeTimeout time.Duration) *http.Server {
	return &http.Server{
		Addr:         puerto,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}
