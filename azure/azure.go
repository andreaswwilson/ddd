package azure

import (
	"context"
	"ddd"
	log "ddd/logger"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
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

func newGraphServiceClient(armtenantID, servicePrincipalID, servicePrincipalSecret string) (*msgraphsdkgo.GraphServiceClient, error) {
	var graphClient *msgraphsdkgo.GraphServiceClient
	if servicePrincipalID == "" || servicePrincipalSecret == "" {
		// Log in with user credentials that are available after running "az login" for testing locally
		log.Debug("servicePrincipalID or servicePrincipalSecret are empty. Logging in using Default Credentials.")
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
		log.Info("Logged in as %s", *me.GetUserPrincipalName())
	} else {
		log.Debug("Logging in using service principal %s", servicePrincipalID)
		cred, err := azidentity.NewClientSecretCredential(armtenantID, servicePrincipalID, servicePrincipalSecret, nil)
		if err != nil {
			return nil, err
		}
		graphClient, err = msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
		if err != nil {
			return nil, err
		}
		log.Info("Logged in with client id %s", servicePrincipalID)
	}
	return graphClient, nil
}

// validate email against entra id
func (service *Service) ValidateMail(mail string) error {
	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")
	requestFilter := fmt.Sprintf("mail eq '%s'", mail)
	requestParameters := &graphusers.UsersRequestBuilderGetQueryParameters{
		Filter: &requestFilter,
	}
	configuration := &graphusers.UsersRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParameters,
	}

	result, err := service.GraphClient.Users().Get(context.Background(), configuration)
	if err != nil {
		return err
	}

	// Initialize iterator
	pageIterator, err := graphcore.NewPageIterator[*models.User](
		result,
		service.GraphClient.GetAdapter(),
		models.CreateMessageCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return fmt.Errorf("error creating page iterator: %v", err)
	}
	count, stopAfter := 0, 25

	mailExist := false
	// Iterate over all pages
	err = pageIterator.Iterate(
		context.Background(),
		func(message *models.User) bool {
			count++
			if message.GetMail() != nil {
				if strings.EqualFold(*message.GetMail(), mail) {
					log.Debug("Found %s in Entra ID", *message.GetMail())
					mailExist = true
					return false
				}
			}
			return count < stopAfter
		})
	if err != nil {
		return fmt.Errorf("error iterating over messages: %v", err)
	}
	if !mailExist {
		return fmt.Errorf("couldn't find '%s' in Entra ID", mail)
	}

	return nil
}
