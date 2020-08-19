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

type ButtonActionPayload struct {
	Question string
	Username string
	UserID string
	ButtonClicked string
	TriggerID string
}

func NewButtonActionPayload(question string, username string, userID string, buttonClicked string, triggerID string) ButtonActionPayload {
	return ButtonActionPayload {
		Question: question,
		Username: username,
		UserID: userID,
		ButtonClicked: buttonClicked,
		TriggerID: triggerID,
	}
}

func NewActionHandler(logger *zap.Logger, conf Conf.Conf) Action_handler {
	return Action_handler {
		Logger: logger,
		Conf: conf,
		Api: slack.New(conf.SLACK_TOKEN),
	}
}

func (p ButtonActionPayload) newViewRequest() slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject("plain_text", "View Answers", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", p.Question, false, false)
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

func (p ButtonActionPayload) newAnswerRequest() slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject("plain_text", "Add an answer", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	answerText := slack.NewTextBlockObject("plain_text", p.Question, false, false)
	answerPlaceholder := slack.NewTextBlockObject("plain_text", "Write something", false, false)
	answerElement := slack.NewPlainTextInputBlockElement(answerPlaceholder, "answer")
	answerElement.Multiline = true
	answer := slack.NewInputBlock("Answer", answerText, answerElement)

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

func getQuestion (payload *slack.InteractionCallback) string {
	return fmt.Sprintf("*%v*", payload.Message.Msg.Text)
}

func (s Action_handler) Events(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Received a slack action")

	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: Unknown action")))
		s.Logger.Error("Error receiving unknown slack action")
		return
	}

	// Get the question text from the Message payload
	var message string
	for _, b := range payload.Message.Msg.Blocks.BlockSet {
		switch b.BlockType() {
		case "section":
			s := b.(*slack.SectionBlock)
			message = s.Fields[0].Text
		default:
			fmt.Println("not section")
		}
	}

	// Configure the payload for use in action handler methods
	p := NewButtonActionPayload(
		message,
		payload.User.Name, 
		payload.User.ID,
		payload.ActionCallback.BlockActions[0].Value,
		payload.TriggerID,
	)

	switch payload.ActionCallback.BlockActions[0].Value {
	case "view_clicked":
		modalRequest := p.newViewRequest()
		_, err = s.Api.OpenView(payload.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
	case "answer_clicked":
		modalRequest := p.newAnswerRequest()
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
