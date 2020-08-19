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

func newViewRequest() slack.ModalViewRequest {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "View Answers", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerContent := fmt.Sprintf("*Question: *")
	headerText := slack.NewTextBlockObject("mrkdwn", headerContent, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	firstNameText := slack.NewTextBlockObject("plain_text", "First Name", false, false)
	firstNamePlaceholder := slack.NewTextBlockObject("plain_text", "Enter your first name", false, false)
	firstNameElement := slack.NewPlainTextInputBlockElement(firstNamePlaceholder, "firstName")
	// Notice that blockID is a unique identifier for a block
	firstName := slack.NewInputBlock("First Name", firstNameText, firstNameElement)

	lastNameText := slack.NewTextBlockObject("plain_text", "Last Name", false, false)
	lastNamePlaceholder := slack.NewTextBlockObject("plain_text", "Enter your first name", false, false)
	lastNameElement := slack.NewPlainTextInputBlockElement(lastNamePlaceholder, "lastName")
	lastName := slack.NewInputBlock("Last Name", lastNameText, lastNameElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			firstName,
			lastName,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func newAnswerRequest() slack.ModalViewRequest {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "Add an Answer", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	answerText := slack.NewTextBlockObject("plain_text", "How to setup a go api?", false, false)
	answerPlaceholder := slack.NewTextBlockObject("plain_text", "Write something", false, false)
	answerElement := slack.NewPlainTextInputBlockElement(answerPlaceholder, "answer")
	answer := slack.NewInputBlock("Answer 1", answerText, answerElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			answer,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
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

	switch payload.ActionCallback.BlockActions[0].Value {
	case "view_clicked":
		modalRequest := newViewRequest()
		_, err = s.Api.OpenView(payload.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
	case "answer_clicked":
		modalRequest := newAnswerRequest()
		_, err = s.Api.OpenView(payload.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Payload %+v", payload)
	fmt.Printf("Message button pressed by user %s with value %s", payload.User.Name, payload.ActionCallback.BlockActions[0].Value)
}
