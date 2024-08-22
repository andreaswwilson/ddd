package ddd

type JiraForm struct {
	Kostnadsoppfolger string `json:"kostnadsoppfolger"`
}

type JiraFormService interface {
	Get(key string) (*JiraForm, error)
}
