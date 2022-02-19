package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)

	r.Body.Close()

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
	return
}

func (s *server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("Request " + r.RequestURI + " from " + r.RemoteAddr)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *server) emptyresponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	defer r.Body.Close()
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {

	s.router.Use(s.loggingMiddleware)

	s.router.Methods("OPTIONS").HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, PUT, OPTIONS")
			rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
			rw.WriteHeader(http.StatusOK)
		})

	s.ConfigureNotificationRouter()
}

func (s *server) ConfigureNotificationRouter() {

	router := s.router.PathPrefix("/api/notification").Subrouter()
	router.HandleFunc("/user/{id}", s.HandleRequestAccept()).Methods("POST") // Отправить заявку
}

func (s *server) HandleRequestAccept() http.HandlerFunc {
	type Request struct {
		int `json:"id"`
	}
	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
			fmt.Println(err)
			return
		}

		var notification Notification
		json.NewDecoder(request.Body).Decode(&notification)

		if notifications, ok := s.BufferNotiff[userID]; ok {
			notifications.Lock.Lock()
			notifications.List = append(notifications.List, notification)
			notifications.Lock.Unlock()
		}

		//s.BufferNotiff[userID] = append()

		s.respond(w, request, http.StatusOK, "ok")
	}
}
