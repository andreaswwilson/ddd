//go:build integration
// +build integration

package azure_test

import (
	"ddd/azure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		email         string
		expectedError bool
	}{
		{"Test that a valid email returns nil", "skattejakten@skatteetaten.no", false},
		{"Test that a invalid email returns error", "asfasf", true},
	}

	client, err := azure.NewGraphServiceClient("", "", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err = azure.ValidateMail(client, test.email)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
