package azure

import (
	"context"
	"ddd"
	log "ddd/logger"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// Make sure the service implements all methods
var _ ddd.AzureService = (*Service)(nil)

type Service struct {
	GraphClient *msgraphsdkgo.GraphServiceClient
}

func NewService() (*Service, error) {
	client, err := newGraphServiceClient()
	if err != nil {
		return nil, err
	}
	return &Service{GraphClient: client}, nil
}

// newGraphServiceClient creates a new instance of GraphServiceClient.
// Log in based on environmental variables
// https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication?tabs=bash#option-1-define-environment-variables
// Returns:
// - *msgraphsdkgo.GraphServiceClient: A pointer to the new GraphServiceClient instance.
// - error: An error object if any error occurs during the creation of the client.
func newGraphServiceClient() (*msgraphsdkgo.GraphServiceClient, error) {
	log.Debug("Creating azure graph service client")
	// Log in with user credentials that are available after running "az login" for testing locally
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	graphClient, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil, err
	}
	log.Debug("Azure graph service client created")
	return graphClient, nil
}

// ValidateEmail checks if the provided email exists in Entra ID.
//
// Parameters:
// - email: The email address to be validated.
//
// Returns: nil if email exist Entra ID
// - error: An error object if the email is not found or if any error occurs during the validation process.
func (service *Service) ValidateEmail(email string) error {
	requestFilter := fmt.Sprintf("mail eq '%s'", email)
	requestParameters := &users.UsersRequestBuilderGetQueryParameters{
		Filter: &requestFilter,
	}
	configuration := &users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	result, err := service.GraphClient.Users().Get(context.Background(), configuration)
	if err != nil {
		return err
	}

	if len(result.GetValue()) == 0 {
		return fmt.Errorf("couldn't find '%s' in Entra ID", email)
	}

	for _, user := range result.GetValue() {
		if strings.EqualFold(*user.GetMail(), email) {
			return nil
		}
	}

	return fmt.Errorf("couldn't find '%s' in Entra ID", email)
}

func (s *Service) GetSubscriptionID(displayname string) (string, error) {
	subscriptions, err := s.GraphClient.Subscriptions().Get(context.Background(), nil)
	if err != nil {
		return "", fmt.Errorf("GetSubscriptionID: couldn't get subscriptions: %w", err)
	}
	fmt.Println(subscriptions)
	return "", nil
}
