package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router
	shoppingItems []Item
}

func NewServer() *Server {
	s := Server{
		mux.NewRouter(),
		[]Item{},
	}
	s.routes()
	return &s
}

func (s *Server) routes() {
	s.HandleFunc("/shopping-items", s.createShoppingItem()).Methods("POST")
	s.HandleFunc("/shopping-items", s.listShoppingItem()).Methods("GET")
	s.HandleFunc("/shopping-items/{id}", s.deleteShoppingItem()).Methods("DELETE")
	s.HandleFunc("/shopping-items/{id}", s.updateShoppingItem()).Methods("PATCH")
}

func (s *Server) createShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i Item
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, item := range s.shoppingItems {
			if item.Name == i.Name {
				http.Error(w, "Duplicated name", http.StatusBadRequest)
				return
			}
		}

		i.ID = uuid.New()
		s.shoppingItems = append(s.shoppingItems, i)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.shoppingItems); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) deleteShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for i, item := range s.shoppingItems {
			if item.ID == id {
				s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
				return
			}
		}

		http.Error(w, "No id matched", http.StatusBadRequest)
	}
}

func (s *Server) updateShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		index, err := s.inDatabase(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.shoppingItems[index].Name = item.Name
	}
}

func (s *Server) inDatabase(id uuid.UUID) (int, error) {
	for i, item := range s.shoppingItems {
		if item.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("No id matched")
}
