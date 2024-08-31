package ddd

type AzureService interface {
	ValidateEmail(email string) error
}
