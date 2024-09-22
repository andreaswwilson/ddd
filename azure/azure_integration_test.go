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
		{"Test that a valid email returns nil", "andreasw@gmail.com", false},
		{"Test that a invalid email returns error", "asfasf", true},
	}

	azureService, err := azure.NewService()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err = azureService.ValidateEmail(test.email)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
