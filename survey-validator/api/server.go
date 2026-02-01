package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/survey-validator/engine"
	"github.com/survey-validator/models"
)

type Server struct {
	engine *engine.Engine
	addr   string
}

func NewServer(addr string) *Server {
	return &Server{
		engine: engine.NewEngine(),
		addr:   addr,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/v1/validate", s.handleValidate)

	handler := s.loggingMiddleware(mux)

	log.Printf("Starting server on %s", s.addr)
	return http.ListenAndServe(s.addr, handler)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "survey-validator",
	})
}

func (s *Server) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed. Use POST.")
		return
	}

	var surveyData models.SurveyData
	if err := json.NewDecoder(r.Body).Decode(&surveyData); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	defer r.Body.Close()

	report := s.engine.Validate(&surveyData)
	s.respondJSON(w, http.StatusOK, report)
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, map[string]string{
		"error": message,
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}
