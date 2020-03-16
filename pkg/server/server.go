package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/CESARBR/knot-babeltower/docs" // This blank import is needed in order to documentation be provided by the server
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/controllers"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	port           int
	logger         logging.Logger
	userController *controllers.UserController
	srv            *http.Server
}

// Health represents the service's health status
type Health struct {
	Status string `json:"status"`
}

// NewServer creates a new server instance
func NewServer(port int, logger logging.Logger, userController *controllers.UserController) Server {
	return Server{port, logger, userController, nil}
}

// Start starts the http server
func (s *Server) Start(started chan bool) {
	routers := s.createRouters()
	s.logger.Infof("Listening on %d", s.port)
	started <- true
	s.srv = &http.Server{Addr: fmt.Sprintf(":%d", s.port), Handler: s.logRequest(routers)}
	err := s.srv.ListenAndServe()
	if err != nil {
		s.logger.Error(err)
		started <- false
	}
}

// Stop stops the server
func (s *Server) Stop() {
	err := s.srv.Shutdown(context.TODO())
	if err != nil {
		s.logger.Error(err)
	}
}

// @title Babeltower API
// @version 1.0
// @description This is the babeltower HTTP API documentation.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func (s *Server) createRouters() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", s.healthcheckHandler)
	r.HandleFunc("/users", s.userController.Create).Methods("POST")
	r.HandleFunc("/tokens", s.userController.CreateToken).Methods("POST")
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)).Methods("GET")
	return r
}

// Healthcheck godoc
// @Summary Verify the service health
// @Produce json
// @Success 200 {object} Health
// @Router /healthcheck [get]
func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(&Health{Status: "online"})
	_, err := w.Write(response)
	if err != nil {
		s.logger.Errorf("Error sending response, %s\n", err)
	}
}

func (s *Server) logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
