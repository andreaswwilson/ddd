package mock

import "ddd"

var _ ddd.AzureService = (*AzureService)(nil)

type AzureService struct{}

func (s *AzureService) ValidateEmail(email string) error {
	return nil
}
