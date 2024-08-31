package mock

import "ddd"

var _ ddd.AzureService = (*AzureService)(nil)

type AzureService struct {
	ValidateEmailFn func(email string) error
}

func (s *AzureService) ValidateEmail(email string) error {
	return s.ValidateEmailFn(email)
}
