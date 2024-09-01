package ddd

type NewSubscriptionOrder struct {
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

type NewSubscriptionOrderService interface {
	Get(key string) (*NewSubscriptionOrder, error)
}
