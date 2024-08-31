package http_test

import (
	"ddd"
	dddhttp "ddd/http"
	"ddd/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Service struct {
	*dddhttp.Service

	// Mock services
	AzureService mock.AzureService
}

func TestGetJiraForm(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		description  string
		expected     *ddd.JiraForm
		jsonResponse string
	}{
		{
			description:  "a full and valid json response returns a fully populated jiraForm struct",
			jsonResponse: `[{"label":"Forstått?","answer":"Jeg har lest og forstått mitt ansvar for å avklare finansiering for ressurser i sky"},{"label":"Kostnadsoppfølger","questionKey":"kostnadsoppfolger","answer":"lukas.nilsen@skatteetaten.no"},{"label":"Maks mnd budsjett","questionKey":"budgetAmount","answer":"100000"},{"label":"Finansiering","questionKey":"finansiering","answer":"Test"},{"label":"Finansiering ved prosjektslutt","questionKey":"finansieringVedProsjektslutt","answer":"Test"},{"label":"Avslutningsdato for prosjekt","answer":""},{"label":"BDM","answer":"lukas.nilsen@skatteetaten.no"},{"label":"Forretningsprodukt","questionKey":"forretningsprodukt","answer":"Test"},{"label":"Forstått?","answer":"Jeg bekrefter at jeg har ført opp korrekt informasjon om finansiering, og forstår at BDM må godkjenne budsjett"},{"label":"Subscription Navn","questionKey":"subscriptionName","answer":"testcase1"},{"label":"Subscription eier (organisatorisk)","questionKey":"l2Approver","answer":"lukas.nilsen@skatteetaten.no"},{"label":"Sikkerhetsansvarlig","questionKey":"securityContact","answer":"lukas.nilsen@skatteetaten.no"},{"label":"Business Bestiller Referanse","questionKey":"businessBestillerReferanse","answer":"3500OLESKA"},{"label":"Business Group","questionKey":"businessGroup","answer":"NA"},{"label":"Business Organisasjon","questionKey":"businessOrg","answer":"Test"},{"label":"Management tre","questionKey":"managementTree","answer":"Prod"},{"label":"Management Group","questionKey":"managementGroup","answer":"Corp"},{"label":"Benytte eksisterende tilgangsrolle?","questionKey":"newPimGroup","answer":"Nei, jeg ønsker at det opprettes en ny tilgangsrolle"},{"label":"Navn til Entra ID Gruppe for PIM","questionKey":"pimApproverNew","answer":"[PIM-APR] testcase1"},{"label":"Miljø","questionKey":"environment","answer":"synt"},{"label":"VNet Størrelse","questionKey":"vnetSize","answer":"/20"},{"label":"Budsjettvarsling","questionKey":"budgetContact","answer":"lukas.nilsen@skatteetaten.no,andreaswago.wilson@skatteetaten.no"}]`,
			expected: &ddd.JiraForm{
				BudgetAmount:                 100000,
				BudgetContact:                []string{"lukas.nilsen@skatteetaten.no", "andreaswago.wilson@skatteetaten.no"},
				BusinessBestillerReferanse:   "3500OLESKA",
				BusinessOrg:                  "Test",
				EntraIDName:                  "[PIM-APR] testcase1",
				Environment:                  "synt",
				Finansiering:                 "Test",
				FinansieringVedProsjektslutt: "Test",
				Forretningsprodukt:           "Test",
				Kostnadsoppfolger:            "lukas.nilsen@skatteetaten.no",
				L2Approver:                   "lukas.nilsen@skatteetaten.no",
				ManagementGroup:              "Corp",
				ManagementTree:               "Prod",
				SecurityContact:              []string{"lukas.nilsen@skatteetaten.no", "andreaswago.wilson@skatteetaten.no"},
				SubscriptionName:             "testcase1",
				VNetSize:                     20,
			},
		},
		{
			description:  "json response where budgetContact with extra whitespace and starting and ending with ',' cleans the data and return a valid slice of strings",
			jsonResponse: `[{"label":"Budsjettvarsling","questionKey":"budgetContact","answer":",lukas.nilsen@skatteetaten.no,      andreaswago.wilson@skatteetaten.no     ,"}]`,
			expected:     &ddd.JiraForm{BudgetContact: []string{"lukas.nilsen@skatteetaten.no", "andreaswago.wilson@skatteetaten.no"}},
		},
		{
			description:  "json response where vnet size without '/' returns valid data",
			jsonResponse: `[{"label":"VNet Størrelse","questionKey":"vnetSize","answer":"20"}]`,
			expected:     &ddd.JiraForm{VNetSize: 20},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			server := setupTestServer(test.jsonResponse, http.StatusOK)
			defer server.Close()

			// Create regular service
			newService, _ := dddhttp.NewService("", dddhttp.WithBaseURL(server.URL))
			s := &Service{Service: newService}
			// Assign mock service
			s.Service.AzureService = &s.AzureService
			// Always return nil for validateEmail
			s.AzureService.ValidateEmailFn = func(email string) error {
				return nil
			}
			actual, err := s.Get("")

			assert.Nil(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
