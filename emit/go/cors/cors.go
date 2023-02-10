package cors

import (
	"net/http"

	"github.com/rs/cors"
)

type Cors struct {
	*cors.Cors
}

func New() *Cors {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"*"},
		MaxAge:         86400,
	})

	return &Cors{c}
}

func (c Cors) Handler(h http.Handler) http.Handler {
	return c.Cors.Handler(h)
}
