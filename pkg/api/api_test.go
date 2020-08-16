package api

import (
	"fmt"
	"strings"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	Conf "github.com/annaworks/chatservice/pkg/conf"

	"go.uber.org/zap"
)

func TestApi(t *testing.T) {
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout"}
	logger, err := c.Build()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not init zap logger: %v", err))
	}
	defer logger.Sync()

	a := NewApi(logger, Conf.Conf{})
	a.Init()

	// health endpoint
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	a.router.ServeHTTP(rr, req)
	if code := rr.Code; code != http.StatusOK {
		t.Errorf("expected status %d but got %d\n", http.StatusOK, code)
	}

	// slack endpoint
	// success case
	slack_challenge := `{
		"token": "Jhj5dZrVaK7ZwHHjRyZWjbDl",
		"challenge": "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P",
		"type": "url_verification"
	}`
	challenge := "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P"
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/slack", strings.NewReader(slack_challenge))
	rr = httptest.NewRecorder()
	a.router.ServeHTTP(rr, req)
	if code := rr.Code; code != http.StatusOK {
		t.Errorf("expected status %d but got %d\n", http.StatusOK, code)
	}
	if resp := rr.Body.String(); resp != challenge {
		t.Errorf("Expected %s but got %s\n", challenge, resp)
	}
	// missing challenge field for url_verification
	const bad_slack_challenge = `{
		"token": "Jhj5dZrVaK7ZwHHjRyZWjbDl",
		"challenge_typo": "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P",
		"type": "url_verification"
	}`
	expected := ""
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/slack", strings.NewReader(bad_slack_challenge))
	rr = httptest.NewRecorder()
	a.router.ServeHTTP(rr, req)
	if resp := rr.Body.String(); resp != expected {
		t.Errorf("Expected %s but got %s\n", expected, resp)
	}
	// unknown event
	fake_slack_event := `{
		"type": "fake_event"
	}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/slack", strings.NewReader(fake_slack_event))
	rr = httptest.NewRecorder()
	a.router.ServeHTTP(rr, req)
	if code := rr.Code; code != http.StatusBadRequest {
		t.Errorf("expected status %d but got %d\n", http.StatusBadRequest, code)
	}
}
