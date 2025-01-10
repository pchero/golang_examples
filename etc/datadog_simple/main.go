package main

import (
	"context"
	"fmt"
	"os"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

func main() {
	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: os.Getenv("DD_CLIENT_API_KEY"),
			},
			"appKeyAuth": {
				Key: os.Getenv("DD_CLIENT_APP_KEY"),
			},
		},
	)

	body := *datadogV2.NewUserCreateRequest(*datadogV2.NewUserCreateData(*datadogV2.NewUserCreateAttributes("jane.doe@example.com"), datadogV2.UsersType("users")))

	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	usersApi := datadogV2.NewUsersApi(apiClient)

	resp, r, err := usersApi.CreateUser(ctx, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	responseData := resp.GetData()
	fmt.Fprintf(os.Stdout, "User ID: %s", responseData.GetId())
}
