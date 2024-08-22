package http

import (
	"ddd"
	"encoding/json"
	"fmt"
	"net/http"
)

type JiraFormService struct {
	Client *Client
}

type jiraFormResponse struct {
	Label       string `json:"label"`
	QuestionKey string `json:"questionKey"`
	Answer      string `json:"answer"`
}

func NewJiraFormService(token string, options ...ClientOptionFunc) (*JiraFormService, error) {
	service := &JiraFormService{}

	client, err := NewClient(token, options...)
	if err != nil {
		return nil, err
	}
	service.Client = client

	return service, nil
}

func (service JiraFormService) Get(key string) (*ddd.JiraForm, error) {
	response := []jiraFormResponse{}
	u := ddd.JiraForm{}
	path := fmt.Sprintf("jira/%s", key)
	req, err := service.Client.NewRequest(http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	_, err = service.Client.Do(req, &response)
	if err != nil {
		return nil, err
	}

	// Unpack the data from jira since each input is wrapped inside a jiraFormResponse
	jsonData := make(map[string]string)
	for _, item := range response {
		jsonData[item.QuestionKey] = item.Answer
	}
	// Turn the map into a json object so that we can unmarshal it into the jiraForm struct
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
