package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Health_handler struct {
	Logger *zap.Logger
}

func (h Health_handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": "live"})

	h.Logger.Info("health endpoint pinged")
}
