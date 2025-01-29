package postgres

import (
	"fmt"

	thing "github.com/username/repo"
)

type thingService struct {
	psql *psqlService
}

func NewThingService(psql *psqlService) *thingService {
	return &thingService{psql}
}

func (ts thingService) CreateThing(t thing.Thing) error {
	_, err := ts.psql.db.Exec("INSERT INTO things (name, description, type) VALUES ($1, $2, $3)", t.Name, t.Description, t.Type)
	if err != nil {
		return fmt.Errorf("failed to insert new thing: %s", err.Error())
	}
	return nil
}

func (ts thingService) GetThing(id string) (thing.Thing, error) {
	var thing thing.Thing
	err := ts.psql.db.QueryRow("SELECT * FROM things WHERE thing_id = $1", id).Scan(&thing.ID, &thing.Name, &thing.Description, &thing.CreatedAt, &thing.Type)
	if err != nil {
		return thing, fmt.Errorf("failed to retrieve thing: %s", err.Error())
	}
	return thing, nil
}

func (ts thingService) GetAllThings() ([]thing.Thing, error) {
	rows, err := ts.psql.db.Query("SELECT * FROM things")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve things: %s", err.Error())
	}
	defer rows.Close()

	var things []thing.Thing
	for rows.Next() {
		var thing thing.Thing
		if err := rows.Scan(&thing.ID, &thing.Name, &thing.Description, &thing.CreatedAt, &thing.Type); err != nil {
			return nil, fmt.Errorf("failed to retrieve thing: %s", err.Error())
		}
		things = append(things, thing)
	}
	return things, nil
}

func (ts thingService) UpdateThing(id string, thing thing.Thing) error {
	_, err := ts.psql.db.Exec("UPDATE things SET name = $1, description = $2, type = $3 WHERE thing_id = $4", thing.Name, thing.Description, thing.Type, id)
	if err != nil {
		return fmt.Errorf("failed to update thing: %s", err.Error())
	}
	return nil
}

func (ts thingService) DeleteThing(id string) error {
	_, err := ts.psql.db.Exec("DELETE FROM things WHERE thing_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete thing: %s", err.Error())
	}
	return nil
}
