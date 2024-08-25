package ddd

type JiraForm struct {
	BudgetAmount                 int
	BudgetContact                []string
	EntraIDName                  string
	Kostnadsoppfolger            string
	L2Approver                   string
	ManagementTree               string
	Environment                  string
	SubscriptionName             string
	VNetSize                     int
	BusinessBestillerReferanse   string
	BusinessOrg                  string
	CreateNewPIM                 bool
	EntraIDGroup                 string
	Finansiering                 string
	FinansieringVedProsjektslutt string
	Forretningsprodukt           string
	ManagementGroup              string
	SecurityContact              []string
}

type JiraFormService interface {
	Get(key string) (*JiraForm, error)
}
