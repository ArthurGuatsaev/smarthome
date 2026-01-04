package httpapi

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ArthurGuatsaev/smarthome/internal/storage"
)

type createDeviceReq struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Capabilities []string `json:"capabilities"`
	MQTTDeviceID string   `json:"mqttDeviceId"`
}

type deviceDTO struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Capabilities []string `json:"capabilities"`
	MQTTDeviceID string   `json:"mqttDeviceId"`
	CreatedAt    string   `json:"createdAt"`
}

func (s *Server) handleDevicesCreate(w http.ResponseWriter, r *http.Request) {
	var req createDeviceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}
	if req.Name == "" || req.Type == "" || req.MQTTDeviceID == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "name, type, mqttDeviceId required")
		return
	}

	capsJSON, _ := json.Marshal(req.Capabilities)

	d := storage.Device{
		ID:           newID(),
		Name:         req.Name,
		Type:         req.Type,
		MQTTDeviceID: req.MQTTDeviceID,
		Capabilities: string(capsJSON),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.app.Devices.Create(r.Context(), d); err != nil {
		// уникальность mqtt_device_id
		writeError(w, http.StatusConflict, "conflict", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toDeviceDTO(d))
}

func (s *Server) handleDevicesList(w http.ResponseWriter, r *http.Request) {
	items, err := s.app.Devices.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}

	out := make([]deviceDTO, 0, len(items))
	for _, d := range items {
		out = append(out, toDeviceDTO(d))
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleDevicesGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.app.Devices.Get(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "not_found", "device not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toDeviceDTO(d))
}

func (s *Server) handleDevicesDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.app.Devices.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toDeviceDTO(d storage.Device) deviceDTO {
	var caps []string
	_ = json.Unmarshal([]byte(d.Capabilities), &caps)

	return deviceDTO{
		ID:           d.ID,
		Name:         d.Name,
		Type:         d.Type,
		Capabilities: caps,
		MQTTDeviceID: d.MQTTDeviceID,
		CreatedAt:    d.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}
