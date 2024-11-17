package api

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func (api *API) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("request_id")
		if id == "" {
			id = generateID()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(&rec, r)
		requestID := r.Context().Value(requestIDKey).(string)
		status := rec.statusCode
		log.Printf("time: %s IP: %s status: %d request ID: %v",
			time.Now().String(), r.RemoteAddr, status, requestID)
	})

}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func generateID() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	length := rng.Intn(5) + 6
	var id string

	for i := 0; i < length; i++ {
		id += strconv.Itoa(rng.Intn(10))
	}
	return id
}
