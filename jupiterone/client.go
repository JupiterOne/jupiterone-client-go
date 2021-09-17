package jupiterone

import (
	"log"
	"net/http"
	"time"

	//"github.com/ardanlabs/graphql"
	graphql "github.com/hasura/go-graphql-client"
	// "github.com/machinebox/graphql"
)

const DefaultRegion string = "us"

type Config struct {
	APIKey     string
	AccountID  string
	Region     string
	HTTPClient *http.Client
}

type Client struct {
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	apiKey, accountID string
	graphqlClient     *graphql.Client
	RetryTimeout      time.Duration

	Entity   *EntityService
	Rule     *RuleService
	Question *QuestionService
}

type service struct {
	client *Client
}

func (c *Config) getRegion() string {
	region := c.Region

	if region == "" {
		region = DefaultRegion
	}

	log.Printf("[info] Utilizing region: %s", region)
	return region
}

func (c *Config) getGraphQLEndpoint() string {
	return "https://api." + c.getRegion() + ".jupiterone.io/graphql"
}

func NewClient(config *Config) (*Client, error) {
	endpoint := config.getGraphQLEndpoint()

	var client *graphql.Client

	if config.HTTPClient != nil {
		client = graphql.NewClient(endpoint, config.HTTPClient)
	} else {
		client = graphql.NewClient(endpoint, nil)
	}

	// client = graphql.NewClient(endpoint, graphql.)

	jupiterOneClient := &Client{
		apiKey:        config.APIKey,
		accountID:     config.AccountID,
		graphqlClient: client,
		RetryTimeout:  time.Minute,
	}

	// Pass around the single client to each service
	jupiterOneClient.common.client = jupiterOneClient
	jupiterOneClient.Entity = (*EntityService)(&jupiterOneClient.common)
	jupiterOneClient.Rule = (*RuleService)(&jupiterOneClient.common)
	jupiterOneClient.Question = (*QuestionService)(&jupiterOneClient.common)

	return jupiterOneClient, nil
}

func (c *Client) prepareRequest(query string) *graphql.Request {
	req := graphql.NewRequest(query)

	req.Header.Set("LifeOmic-Account", c.accountID)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	return req
}
