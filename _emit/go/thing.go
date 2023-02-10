package thing

import "time"

type Thing struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
}

type ThingService interface {
	CreateThing(thing Thing) error
	GetThing(id string) (Thing, error)
	GetAllThings() ([]Thing, error)
	UpdateThing(id string, thing Thing) error
	DeleteThing(id string) error
}
