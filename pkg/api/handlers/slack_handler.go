package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	Conf "github.com/annaworks/chatservice/pkg/conf"
)

type Slack_handler struct {
	Logger *zap.Logger
	Conf Conf.Conf
}

func (s Slack_handler) Events(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("slack endpoint evented")
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": s.Conf.SLACK_TOKEN,
		"verification_token": s.Conf.SLACK_VERIFICATION_TOKEN,
	})

}
