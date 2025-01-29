package cors

import (
	"net/http"

	"github.com/rs/cors"
)

type Cors struct {
	*cors.Cors
}

func New() *Cors {
	// TODO: Improve cors defaults.
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"*"},
		MaxAge:         86400,
	})

	return &Cors{c}
}
