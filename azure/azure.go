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

func NewService(armtenantID, servicePrincipalID, servicePrincipalSecret string) (*Service, error) {
	client, err := newGraphServiceClient(armtenantID, servicePrincipalID, servicePrincipalSecret)
	if err != nil {
		return nil, err
	}
	return &Service{GraphClient: client}, nil
}

// newGraphServiceClient creates a new instance of GraphServiceClient.
// It uses either service principal credentials or default Azure credentials
// based on the provided input parameters.
//
// Parameters:
// - armtenantID: The Azure Resource Manager tenant ID. If empty, default credentials are used.
// - servicePrincipalID: The service principal ID. If empty, default credentials are used.
// - servicePrincipalSecret: The service principal secret. If empty, default credentials are used.
//
// Returns:
// - *msgraphsdkgo.GraphServiceClient: A pointer to the new GraphServiceClient instance.
// - error: An error object if any error occurs during the creation of the client.
func newGraphServiceClient(armtenantID, servicePrincipalClientID, servicePrincipalSecret string) (*msgraphsdkgo.GraphServiceClient, error) {
	log.Debug("Creating azure graph service client")
	var graphClient *msgraphsdkgo.GraphServiceClient
	if armtenantID == "" || servicePrincipalClientID == "" || servicePrincipalSecret == "" {
		// Log in with user credentials that are available after running "az login" for testing locally
		log.Debug("armTenantID, servicePrincipalID or servicePrincipalSecret are empty. Logging in using Default Credentials.")
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}
		graphClient, err = msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
		if err != nil {
			return nil, err
		}
		// Read user
		me, err := graphClient.Me().Get(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		log.Info("Authenticated against graph.microsoft.com/.default as %s", *me.GetUserPrincipalName())
	} else {
		log.Debug("Logging in using service principal id and password.")
		log.Debug("armtenantID: %s", armtenantID)
		log.Debug("servicePrincipalClientID: %s", servicePrincipalClientID)
		log.Debug("servicePrincipalSecret: %s******************", servicePrincipalSecret[0:3])
		cred, err := azidentity.NewClientSecretCredential(armtenantID, servicePrincipalClientID, servicePrincipalSecret, nil)
		if err != nil {
			return nil, err
		}
		graphClient, err = msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
		if err != nil {
			return nil, err
		}
		log.Info("Authenticated against graph.microsoft.com/.default with client_id %s", servicePrincipalClientID)
	}
	log.Debug("Azure graph service client created")
	return graphClient, nil
}

// ValidateEmail checks if the provided email exists in Entra ID.
//
// Parameters:
// - email: The email address to be validated.
//
// Returns:
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
