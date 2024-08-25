package http

import (
	"ddd"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type JiraFormService struct {
	Client *Client
}

type jiraFormResponse struct {
	QuestionKey string `json:"questionKey"`
	Answer      string `json:"answer"`
}

type jiraForm struct {
	BudgetAmount                 int      `json:"budgetAmount,string"`
	BudgetContact                []string `json:"budgetContact,omitempty"`
	EntraIDName                  string   `json:"pimApproverNew,omitempty"`
	Kostnadsoppfolger            string   `json:"kostnadsoppfolger"`
	L2Approver                   string   `json:"l2Approver"`
	ManagementTree               string   `json:"managementTree"`
	Environment                  string   `json:"environment"`
	SubscriptionName             string   `json:"subscriptionName"`
	VNetSize                     int      `json:"vnetSize"`
	BusinessBestillerReferanse   string   `json:"businessBestillerReferanse"`
	BusinessOrg                  string   `json:"businessOrg"`
	CreateNewPIM                 bool     `json:"createNewPim"`
	EntraIDGroup                 string   `json:"entraIDGroup"`
	Finansiering                 string   `json:"finansiering"`
	FinansieringVedProsjektslutt string   `json:"finansieringVedProsjektslutt"`
	Forretningsprodukt           string   `json:"forretningsprodukt"`
	ManagementGroup              string   `json:"managementGroup"`
	SecurityContact              []string `json:"securityContact"`
}

// Overload UnmarshalJSON for JiraForm to turn strings into []strings
// and to convert vnetSize from string to int
func (jf *jiraForm) UnmarshalJSON(data []byte) error {
	type Alias jiraForm
	tempStruct := &struct {
		BudsjettVarsling string `json:"budgetContact,omitempty"`
		SecurityContact  string `json:"securityContact,omitempty"`
		VNetStørrelse    string `json:"vnetSize,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(jf),
	}
	if err := json.Unmarshal(data, &tempStruct); err != nil {
		return err
	}
	// removes all white space, trims leading/trailing commas, and splits based on ","
	if tempStruct.BudsjettVarsling != "" {
		jf.BudgetContact = strings.Split(strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(tempStruct.BudsjettVarsling, " ", ""), ","), ","), ",")
	}
	// removes all white space, trims leading/trailing commas, and splits based on ","
	if tempStruct.SecurityContact != "" {
		jf.SecurityContact = strings.Split(strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(tempStruct.BudsjettVarsling, " ", ""), ","), ","), ",")
	}
	// remove / and convert to int
	if tempStruct.VNetStørrelse != "" {
		vnet, err := strconv.Atoi(strings.ReplaceAll(tempStruct.VNetStørrelse, "/", ""))
		if err != nil {
			return err
		}
		jf.VNetSize = vnet
	}
	return nil
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
	jiraForm := jiraForm{}
	path := fmt.Sprintf("jira/%s", key)
	request, err := service.Client.NewRequest(http.MethodGet, path)
	request.Header.Add("Accept", "application/json")
	if err != nil {
		return nil, err
	}

	_, err = service.Client.Do(request, &response)
	if err != nil {
		return nil, err
	}

	// Unpack the data from jira since each input is wrapped inside a jiraFormResponse
	jsonData := make(map[string]string)
	for _, item := range response {
		jsonData[item.QuestionKey] = strings.TrimSpace(item.Answer)
	}
	// Turn the map into a json object so that we can unmarshal it into the jiraForm struct
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &jiraForm)
	if err != nil {
		return nil, err
	}

	return &ddd.JiraForm{
		BudgetAmount:                 jiraForm.BudgetAmount,
		BudgetContact:                jiraForm.BudgetContact,
		EntraIDName:                  jiraForm.EntraIDName,
		Kostnadsoppfolger:            jiraForm.Kostnadsoppfolger,
		L2Approver:                   jiraForm.L2Approver,
		ManagementTree:               jiraForm.ManagementTree,
		Environment:                  jiraForm.Environment,
		SubscriptionName:             jiraForm.SubscriptionName,
		VNetSize:                     jiraForm.VNetSize,
		BusinessBestillerReferanse:   jiraForm.BusinessBestillerReferanse,
		BusinessOrg:                  jiraForm.BusinessOrg,
		CreateNewPIM:                 jiraForm.CreateNewPIM,
		EntraIDGroup:                 jiraForm.EntraIDGroup,
		Finansiering:                 jiraForm.Finansiering,
		FinansieringVedProsjektslutt: jiraForm.FinansieringVedProsjektslutt,
		Forretningsprodukt:           jiraForm.Forretningsprodukt,
		ManagementGroup:              jiraForm.ManagementGroup,
		SecurityContact:              jiraForm.SecurityContact,
	}, nil
}
