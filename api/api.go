package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jamesog/example-http-api/database"
)

// API is an instance of the API.
type API struct {
	db database.Storage
}

// NewService returns a new API object.
func NewService(db database.Storage) *API {
	return &API{db: db}
}

// Routes returns the API routes.
func (api *API) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/users", api.GetUsers)

	return r
}

// GetUsers handles GET /users.
func (api *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := api.db.Users(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	j, _ := json.Marshal(users)
	w.Write(j)
}
