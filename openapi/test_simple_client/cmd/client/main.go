package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testoapi/gens/client"
)

type APIClient interface {
}

type apiClient struct{}

// func (h *apiClient) GetActiveflows(ctx context.Context) (*client.FlowManagerActiveflow, error) {
// 	hc := http.Client{}

// 	cl, err := client.NewClientWithResponses("http://localhost:1234", client.WithHTTPClient(&hc))
// 	if err != nil {
// 		return nil, err
// 	}

// 	tmp, err := cl.GetAccesskeysWithResponse(ctx, &client.GetAccesskeysParams{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	c, err := client.NewClient("http://localhost:1234", client.WithHTTPClient(&hc))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// c.GetAccesskeys

// 	req := &client.GetActiveflowsParams{}
// 	res, err := c.GetActiveflows(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	af := client.FlowManagerActiveflow{}
// 	json.NewDecoder(res.Body).Decode(&af)

// 	return &af, nil
// }

func main() {

	ctx := context.Background()

	// custom HTTP client
	hc := http.Client{}

	// with a raw http.Response
	{
		c, err := client.NewClient("https://voipbin.net", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		req := &client.GetActiveflowsParams{}
		res, err := c.GetActiveflows(ctx, req)
		if err != nil {
			return
		}

		af := client.FlowManagerActiveflow{}
		json.NewDecoder(res.Body).Decode(&af)

		return &af, nil

		res.Body.Read()

		resp, err := c.GetClient(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Expected HTTP 200 but received %d", resp.StatusCode)
		}
	}

	// or to get a struct with the parsed response body
	{
		c, err := client.NewClientWithResponses("http://localhost:1234", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := c.GetClientWithResponse(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode() != http.StatusOK {
			log.Fatalf("Expected HTTP 200 but received %d", resp.StatusCode())
		}

		fmt.Printf("resp.JSON200: %v\n", resp.JSON200)
	}
}
