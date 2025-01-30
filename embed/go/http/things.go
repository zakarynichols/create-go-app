package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	thing "github.com/username/repo"
)

// TODO: Generate services with user provided input instead of "thing".

func (s *Server) RegisterThingRoutes(ctx context.Context) {
	s.router.Handle("/things", handleCreateThing(ctx, s.thingService)).Methods("POST")
	s.router.Handle("/things/{id}", handleGetThing(ctx, s.thingService)).Methods("GET")
	s.router.Handle("/things", handleGetAllThings(ctx, s.thingService)).Methods("GET")
	s.router.Handle("/things/{id}", handleUpdateThing(ctx, s.thingService)).Methods("PUT")
	s.router.Handle("/things/{id}", handleDeleteThing(ctx, s.thingService)).Methods("DELETE")
}

func handleCreateThing(ctx context.Context, ts thing.ThingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var thing thing.Thing

		if err := json.NewDecoder(r.Body).Decode(&thing); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err := ts.CreateThing(thing)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
	}
}

func handleGetThing(ctx context.Context, ts thing.ThingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, "ID not found in URL", http.StatusBadRequest)
			return
		}

		thing, err := ts.GetThing(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(thing)
	}
}

func handleGetAllThings(ctx context.Context, ts thing.ThingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		things, err := ts.GetAllThings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(things)
	}
}

func handleUpdateThing(ctx context.Context, ts thing.ThingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, "ID not found in URL", http.StatusBadRequest)
			return
		}

		var thing thing.Thing
		if err := json.NewDecoder(r.Body).Decode(&thing); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err := ts.UpdateThing(id, thing)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func handleDeleteThing(ctx context.Context, ts thing.ThingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		err := ts.DeleteThing(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
