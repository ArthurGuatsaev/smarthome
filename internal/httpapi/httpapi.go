package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ArthurGuatsaev/smarthome/internal/app"
	"github.com/ArthurGuatsaev/smarthome/internal/buildinfo"
)

type Server struct {
	mux    *http.ServeMux
	ready  *ReadyState
	app    *app.App
	apiKey string
}

type ReadyState struct {
	// пока просто флаг, позже сюда подключишь проверки БД/MQTT
	isReady bool
}

func NewServer(a *app.App, apiKey string) *Server {
	rs := &ReadyState{isReady: true}
	mux := http.NewServeMux()

	s := &Server{mux: mux, ready: rs, app: a, apiKey: apiKey}

	// system
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("GET /readyz", s.handleReadyz)
	mux.HandleFunc("GET /api/v1/version", s.handleVersion)

	// devices
	mux.HandleFunc("GET /api/v1/devices", s.handleDevicesList)
	mux.HandleFunc("POST /api/v1/devices", s.handleDevicesCreate)
	mux.HandleFunc("GET /api/v1/devices/{id}", s.handleDevicesGet)
	mux.HandleFunc("DELETE /api/v1/devices/{id}", s.handleDevicesDelete)

	return s
}

func (s *Server) Handler() http.Handler {
	// тут подключаем middleware
	return Chain(s.mux,
		RequestID(),
		AccessLog(),
		Recoverer(),
		Timeout(8*time.Second),
		RequireAPIKey(s.apiKey),
	)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	if !s.ready.isReady {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"version": buildinfo.Version,
		"commit":  buildinfo.Commit,
		"date":    buildinfo.Date,
	})
}

func (s *Server) SetReady(v bool) {
	s.ready.isReady = v
}
