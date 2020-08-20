package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"
	"errors"

	"go.uber.org/zap"
	Conf "github.com/annaworks/chatservice/pkg/conf"
	// "github.com/annaworks/chatservice/pkg/es"
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

type Question struct {
	Question  string `json:"question"`
	User      string `json:"user"`
}

func NewQuestion(value, user string) *Question {
	return &Question{
		Question:  value,
		User:      user,
	}
}

const EsQuestionMapping = `{
	"settings": {
		"number_of_shards": 3,
		"number_of_replicas": 1
	},
	"mappings": {
		"question": {
			"properties": {
				"user": {"type": "text"},
				"value": {"type": "text"}
			}
		}
	}
}`


func newButtonActionPayload(question, username, userID, buttonClicked, triggerID string) ButtonActionPayload {
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

func getQuestionText(payload slack.InteractionCallback) (string, error)  {
	var message string
	loop:for _, b := range payload.Message.Msg.Blocks.BlockSet {
		switch b.BlockType() {
		case "section":
			s := b.(*slack.SectionBlock)
			message = s.Fields[0].Text
			break loop
		default:
			return message, errors.New(fmt.Sprintf("Error: Unknown block type found in block action: %v", b.BlockType()))
		}
	}
	return message, nil
}

func getQuestionFromMessageInput(payload slack.InteractionCallback) (string, error) {
	var message string
	loop:for _, b := range payload.View.Blocks.BlockSet {
		switch b.BlockType() {
		case slack.MBTInput:
			s := b.(*slack.InputBlock)
			message = s.Label.Text
			break loop
		default:
			return message, errors.New(fmt.Sprintf("Error: Unknown block type found in view submission: %v", b.BlockType()))
		}
	}
	return message, nil
}

func (p ButtonActionPayload) newViewRequest() slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject("plain_text", "View answers", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Add answer", false, false)

	headerContent := fmt.Sprintf("*Question:*\n %v", p.Question)
	headerText := slack.NewTextBlockObject("mrkdwn", headerContent, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Empty state on initial view
	emptyContentText := slack.NewTextBlockObject("plain_text", "There are currently no answers for this question.", false, false)
	emptyContentSection := slack.NewSectionBlock(emptyContentText, nil, nil)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			emptyContentSection,
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

	message, err := getQuestionText(payload)
	if err != nil {
		fmt.Printf("Error getting question from interaction callback: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: Unknown interaction callback")))
		s.Logger.Error("Error getting question from interaction callback")
		return
	}

	// Configure the payload for use in action handler methods
	p := newButtonActionPayload(
		message,
		payload.User.Name, 
		payload.User.ID,
		payload.ActionCallback.BlockActions[0].Value,
		payload.TriggerID,
	)

	switch payload.Type {
		case slack.InteractionTypeBlockActions:
			// Check the buttonID value of interaction
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
		
		case slack.InteractionTypeViewSubmission:
			// message, err := getQuestionFromMessageInput(payload)
			// if err != nil {
			// 	fmt.Printf("Error getting question from view submission: %v", err)
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	w.Write([]byte(fmt.Sprintf("Error: Unknown view submission")))
			// 	s.Logger.Error("Error getting question from view submission")
			// 	return
			// }

			//TODO: query es using the question from message input - to get question id to save message input with

		default:
			fmt.Printf("Unknown interaction callback type: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: Unknown interaction callback type")))
			s.Logger.Error("Error receiving unknown interaction callback type")
			return
	}

	// fmt.Printf("Payload %+v", payload)
}