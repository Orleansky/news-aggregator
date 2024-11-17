package api

import (
	"Anastasia/skillfactory/advanced/censoring-service/models"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type key string

const requestIDKey key = "request_id"

type API struct {
	router *mux.Router
}

func New() *API {
	api := API{}
	api.router = mux.NewRouter()
	api.router.Use(api.requestIDMiddleware)
	api.router.Use(api.log)
	api.endpoints()
	return &api
}

func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.HandleFunc("/news/comments", api.moderateHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (api *API) moderateHandler(w http.ResponseWriter, r *http.Request) {
	var data models.Comment
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var bannedWords = []string{"qwerty", "йцукен", "zxvbnm"}
	for _, word := range bannedWords {
		if strings.Contains(strings.ToLower(data.Content), word) {
			http.Error(w, "Ваш комментарий не прошёл модерацию", http.StatusBadRequest)
			return
		}
	}
}
