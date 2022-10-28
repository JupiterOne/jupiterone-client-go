package jupiterone

import (
	"log"
	"net/http"
	"time"

	"github.com/machinebox/graphql"
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
	httpClient        *http.Client
	httpBaseURL       string
	RetryTimeout      time.Duration

	Entity          *EntityService
	Rule            *RuleService
	Question        *QuestionService
	Relationship    *RelationshipService
	Integration     *IntegrationService
	Audit           *AuditService
	Synchronization *SynchronizationService
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

func (c *Config) getHTTPEndpoint() string {
	return "https://api." + c.getRegion() + ".jupiterone.io"
}

func NewClient(config *Config) (*Client, error) {
	endpoint := config.getGraphQLEndpoint()

	var client *graphql.Client
	var httpClient *http.Client

	if config.HTTPClient != nil {
		client = graphql.NewClient(endpoint, graphql.WithHTTPClient(config.HTTPClient))
		httpClient = config.HTTPClient
	} else {
		client = graphql.NewClient(endpoint)
		httpClient = &http.Client{}
	}

	jupiterOneClient := &Client{
		apiKey:        config.APIKey,
		accountID:     config.AccountID,
		graphqlClient: client,
		httpClient:    httpClient,
		httpBaseUrl:   config.getHTTPEndpoint(),
		RetryTimeout:  time.Minute,
	}

	// Pass around the single client to each service
	jupiterOneClient.common.client = jupiterOneClient
	jupiterOneClient.Entity = (*EntityService)(&jupiterOneClient.common)
	jupiterOneClient.Rule = (*RuleService)(&jupiterOneClient.common)
	jupiterOneClient.Question = (*QuestionService)(&jupiterOneClient.common)
	jupiterOneClient.Relationship = (*RelationshipService)(&jupiterOneClient.common)
	jupiterOneClient.Integration = (*IntegrationService)(&jupiterOneClient.common)
	jupiterOneClient.Audit = (*AuditService)(&jupiterOneClient.common)
	jupiterOneClient.Synchronization = (*SynchronizationService)(&jupiterOneClient.common)

	return jupiterOneClient, nil
}

func (c *Client) addAuthHeaders(req *http.Request) *http.Request {
	req.Header.Set("JupiterOne-Account", c.accountID)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	return req
}

func (c *Client) prepareRequest(query string) *graphql.Request {
	req := graphql.NewRequest(query)

	req.Header.Set("JupiterOne-Account", c.accountID)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", "JupiterOne-Client-Go")

	return req
}
