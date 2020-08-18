package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	Conf "github.com/annaworks/chatservice/pkg/conf"
	"github.com/slack-go/slack"
)

type Action_handler struct {
	Logger *zap.Logger
	Conf Conf.Conf
	Api *slack.Client
}

func NewActionHandler(logger *zap.Logger, conf Conf.Conf) Action_handler {
	return Action_handler {
		Logger: logger,
		Conf: conf,
		Api: slack.New(conf.SLACK_TOKEN),
	}
}

func (s Action_handler) Events(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Received a slack action")

	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: Unknown event")))
		s.Logger.Error("Error receiving unknown slack event")
		return
	}

	fmt.Printf("Payload %+v", payload)
	fmt.Printf("Message button pressed by user %s with value %s", payload.User.Name, payload.ActionCallback.BlockActions[0].Value)
}
