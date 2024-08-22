package http

import (
	"ddd/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJiraForm(t *testing.T) {
	t.Parallel()

	for key, value := range mock.JiraFormInput {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
		}))
		defer server.Close()

		js, _ := NewJiraFormService("", WithBaseURL(server.URL))
		jiraForm, err := js.Get("")
		assert.Nil(t, err)
		assert.Equal(t, mock.JiraForm[key], jiraForm)
	}
}
