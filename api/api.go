package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/natasha-m-oliveira/crud-in-memory/db"
)

func NewHandler(ur db.UsersRepository) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", handleCreateUser(ur))
			r.Get("/", handleListUsers(ur))
			r.Get("/{id}", handleGetUser(ur))
			r.Delete("/{id}", handleDeleteUser(ur))
			r.Put("/{id}", handleUpdateUser(ur))
		})
	})

	return r
}

func handleCreateUser(ur db.UsersRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body UserBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Message: "Please provide firstName lastName and bio for the user"}, http.StatusBadRequest)
			return
		}

		user := ur.Insert(db.User{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Biography: body.Biography,
		})

		sendJSON(w, Response{Data: user}, http.StatusCreated)
	}
}

func handleListUsers(ur db.UsersRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := ur.FindAll()

		sendJSON(w, Response{Data: users}, http.StatusOK)
	}
}

func handleGetUser(ur db.UsersRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		user := ur.FindById(db.Id(id))
		if user == (db.User{}) {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}

func handleDeleteUser(ur db.UsersRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		if user := ur.FindById(db.Id(id)); user == (db.User{}) {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		ur.Delete(db.Id(id))

		sendJSON(w, Response{}, http.StatusNoContent)
	}
}

func handleUpdateUser(ur db.UsersRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body UserBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Message: "Please provide firstName lastName and bio for the user"}, http.StatusBadRequest)
			return
		}

		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		if user := ur.FindById(db.Id(id)); user == (db.User{}) {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		user := ur.Update(db.User{
			Id:        db.Id(id),
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Biography: body.Biography,
		})

		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}
