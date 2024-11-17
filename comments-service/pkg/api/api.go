package api

import (
	"Anastasia/skillfactory/advanced/comments-service/pkg/models"
	"Anastasia/skillfactory/advanced/comments-service/pkg/postgres"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type key string

const requestIDKey key = "request_id"

type API struct {
	db     *postgres.Store
	router *mux.Router
}

func New(db *postgres.Store) *API {
	api := &API{
		db:     db,
		router: mux.NewRouter(),
	}
	api.router.Use(api.requestIDMiddleware)
	api.router.Use(api.log)
	api.endpoints()

	return api
}

func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.HandleFunc("/news/{id}/comments", api.commentsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/comments", api.createCommentHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (api *API) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var c models.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.db.CreateComment(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (api *API) commentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	newsID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments, err := api.db.Comments(newsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}
