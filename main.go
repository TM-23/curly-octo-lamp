package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoftgraph/msgraph-sdk-go"
	msgraphauth "github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"log"
)

func main() {
	ctx := context.Background()
	ListApplications(ctx)

}

func ListApplications(ctx context.Context) {
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain default credential: %v", err)
	}
	authProvider, err := msgraphauth.NewAzureIdentityAuthenticationProvider(creds)
	if err != nil {
		log.Fatalf("fail to obtain Azure identity authentication provider: %v", err)
	}
	adapter, err := msgraphsdkgo.NewGraphRequestAdapter(authProvider)
	if err != nil {
		log.Fatalf("failed to obtain a graph request adapter: %v", err)
	}
	client := msgraphsdkgo.NewGraphServiceClient(adapter)

	appRequestBuilder := client.Applications()

	queryParams := &applications.ApplicationsRequestBuilderGetQueryParameters{
		Select: []string{"id", "displayName"},
	}
	requestConfig := &applications.ApplicationsRequestBuilderGetRequestConfiguration{
		QueryParameters: queryParams,
	}
	appCollection, err := appRequestBuilder.Get(ctx, requestConfig)
	if err != nil {
		var oDataErr odataerrors.ODataErrorable
		if errors.As(err, &oDataErr) {
			log.Fatalf("failed to get list of applications: %v", oDataErr)
		} else {
			log.Fatalf("failed to list applications: %v", err)
		}
	}

	if appCollection.GetValue() != nil {
		for _, app := range appCollection.GetValue() {
			appID := *app.GetAppId()
			displayName := *app.GetDisplayName()
			fmt.Printf("Processing Application: %s (ID: %s)\n", displayName, appID)
		}
	} else {
		log.Fatalf("No applications found in the directory")
	}

}
