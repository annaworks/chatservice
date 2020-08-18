package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	Conf "github.com/annaworks/chatservice/pkg/conf"
	"github.com/slack-go/slack"
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

const slash_command = "/annabot"

func (s Slack_handler) Events(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Received a slack event")

	slash, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %+v", err)))
		s.Logger.Error("Error in parsing slash command", zap.Error(err))
		return
	}

	switch slash.Command {
		case slash_command:
			fmt.Printf("%+v\n", slash.Text)
			fmt.Printf("%+v\n", slash.UserName)

			// Header Section
			// headerText := slack.NewTextBlockObject("mrkdwn", "You asked a question", false, false)
			// headerSection := slack.NewSectionBlock(headerText, nil, nil)

			// Divider 
			// divSection := slack.NewDividerBlock()

			// Fields
			questionText := fmt.Sprintf("*%v asked:*\n%v", slash.UserName, slash.Text)
			questionField := slack.NewTextBlockObject("mrkdwn", questionText, false, false)

			fieldSlice := make([]*slack.TextBlockObject, 0)
			fieldSlice = append(fieldSlice, questionField)

			fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

			// Action Buttons
			viewBtnTxt := slack.NewTextBlockObject("plain_text", "View", false, false)
			viewBtn := slack.NewButtonBlockElement("", "view_clicked", viewBtnTxt)

			answerBtnTxt := slack.NewTextBlockObject("plain_text", "Answer", false, false)
			answerBtn := slack.NewButtonBlockElement("", "answer_clicked", answerBtnTxt).WithStyle("primary")

			actionBlock := slack.NewActionBlock("", viewBtn, answerBtn)

			// Build Message with blocks created above
			msg := slack.NewBlockMessage(
				// headerSection,
				// divSection,
				fieldsSection,
				actionBlock,
			)
			msg.ResponseType = slack.ResponseTypeInChannel

			b, err := json.MarshalIndent(msg, "", "    ")

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(b)

			s.Logger.Info("Message with buttons sucessfully sent")
			return 
		default:
			return
	}
}