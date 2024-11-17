package api

import (
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/postgres"
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

// Конструктор API
func New(db *postgres.Store) *API {
	api := API{
		db: db,
	}

	api.router = mux.NewRouter()
	api.router.Use(api.requestIDMiddleware)
	api.router.Use(api.log)
	api.endpoints()
	return &api
}

// Возвращает маршрутизатор для использования
// в качестве аргумента HTTP-сервера
func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.HandleFunc("/news", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/{id}", api.postDetailedHandler).Methods(http.MethodGet, http.MethodOptions)
}

func (api *API) postDetailedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s := mux.Vars(r)["id"]

	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	post, err := api.db.PostDetailed(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s := r.URL.Query().Get("s")
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	posts, pagination, err := api.db.Posts(page, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(pagination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
