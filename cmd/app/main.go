package main

import (
	"ddd/azure"
	"ddd/http"
	"ddd/logger"
	"os"
)

func main() {
	azureService, err := azure.NewService()
	if err != nil {
		logger.Error("%s", err)
		os.Exit(1)
	}
	jiraFormService, err := http.NewService("token")
	jiraFormService.AzureService = azureService
}
