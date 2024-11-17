package api

import (
	"Anastasia/skillfactory/advanced/APIGateway/pkg/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type API struct {
	router *mux.Router
}

type key string

const requestIDKey key = "request_id"

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
	api.router.HandleFunc("/news", api.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/{id}", api.newsDetailedHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/comments", api.createCommentHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(requestIDKey)

	pageStr := r.URL.Query().Get("page")

	s := r.URL.Query().Get("s")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	url := fmt.Sprintf("http://localhost:8081/news?request_id=%s&s=%s&page=%d", id, s, page)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func (api *API) newsDetailedHandler(w http.ResponseWriter, r *http.Request) {
	reqid := r.Context().Value(requestIDKey)
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	wg := sync.WaitGroup{}
	chanResp := make(chan *http.Response, 2)
	chanErr := make(chan error, 2)

	client := &http.Client{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		url := fmt.Sprintf("http://localhost:8081/news/%d?request_id=%s", id, reqid)
		req, err := http.NewRequest(http.MethodGet, url, r.Body)
		if err != nil {
			chanErr <- err
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			chanErr <- err
			return
		}

		chanResp <- resp
	}()

	go func() {
		defer wg.Done()
		url := fmt.Sprintf("http://localhost:8082/news/%d/comments?request_id=%s", id, reqid)
		req, err := http.NewRequest(http.MethodGet, url, r.Body)
		if err != nil {
			chanErr <- err
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			chanErr <- err
			return
		}

		chanResp <- resp
	}()

	wg.Wait()
	close(chanResp)
	close(chanErr)

	for err := range chanErr {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for resp := range chanResp {
		if resp.StatusCode == 500 {
			w.WriteHeader(500)
		}
		io.Copy(w, resp.Body)
		resp.Body.Close()
	}

}

func (api *API) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(requestIDKey)
	c := models.Comment{}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	commentJSON, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urlModerate := fmt.Sprintf("http://localhost:8083/news/comments?request_id=%s", id)
	reqModerate, err := http.NewRequest(http.MethodPost, urlModerate, bytes.NewReader(commentJSON))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reqModerate.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	respModerate, err := client.Do(reqModerate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer respModerate.Body.Close()

	if respModerate.StatusCode == 400 {
		w.WriteHeader(400)
	}

	io.Copy(w, respModerate.Body)

	url := fmt.Sprintf("http://localhost:8082/news/comments?request_id=%s", id)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(commentJSON))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}
