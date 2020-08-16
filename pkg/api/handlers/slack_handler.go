package handlers

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	Conf "github.com/annaworks/chatservice/pkg/conf"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Slack_handler struct {
	Logger *zap.Logger
	Conf Conf.Conf
	Api *slack.Client
}

func NewSlackHandler(logger *zap.Logger, conf Conf.Conf) Slack_handler {
	return Slack_handler {
		Logger: logger,
		Conf: conf,
		Api: slack.New(conf.SLACK_TOKEN),
	}
}

func (s Slack_handler) Events(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Received a slack event")

	buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %+v", err)))
			s.Logger.Error("Error in parsing event", zap.Error(err))
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				s.Api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			}
		}

}
