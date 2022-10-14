package jupiterone

import (
	"context"
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

// IntegrationService handles the integration-related API requests.
type IntegrationService service

// IntegrationInstance represents an instance of an integration.
// Example: The account has an instance of the AWS integration.
type IntegrationInstance struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	IntegrationDefinitionID string `json:"integrationDefinitionId"`
}

// IntegrationInstanceResponse is a slice of integration instances and
// pagination information.
type IntegrationInstanceResponse struct {
	Instances []*IntegrationInstance `json:"instances"`
	PageInfo  PageInfo               `json:"pageInfo"`
}

// IntegrationDefinition defines information about a single integration. It
// does not mean there are configured instances of that integration, only that
// the integration is available to be configured.
//
// The ConfigFields of an integration represent the integration specific
// configuration needed to run that integration. It is different for
// each integration, so it uses the json.RawMessage type.
type IntegrationDefinition struct {
	ID               string            `json:"id"`
	IntegrationType  string            `json:"integrationType"`
	IntegrationClass []string          `json:"integrationClass"`
	Name             string            `json:"name"`
	Title            string            `json:"title"`
	RepoWebLink      string            `json:"repoWebLink"`
	ConfigFields     []json.RawMessage `json:"configFields"`
}

// IntegrationDefinitionsResponse is a slice of integration definitions
// and pagination information.
type IntegrationDefinitionsResponse struct {
	Definitions []*IntegrationDefinition `json:"definitions"`
	PageInfo    PageInfo                 `json:"pageInfo"`
}

// PageInfo is the pagination information for ListInstances and ListDefinitions.
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

// ListDefinitions lists all the IntegrationDefinitions in the current account.
// ListDefinitions returns a reference to an IntegrationDefinitionResponse object
// which contains the Definitions and PageInfo used to request additional
// definitions (if they exist).
//
// The first call to ListDefinitions should pass nil as the cursor. The caller
// should check PageInfo.HasNextPage to see if there is additional data available.
// If there is, the caller should pass PageInfo.Cursor on subsequent calls.
func (s *IntegrationService) ListDefinitions(cursor *string) (*IntegrationDefinitionsResponse, error) {
	req := s.client.prepareRequest(`
		query IntegrationDefinitions($cursor: String) {
			integrationDefinitions(cursor: $cursor) {
				definitions {
					id
					integrationType
					integrationClass
					name
					repoWebLink
					title
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}`)

	if cursor != nil {
		req.Var("cursor", cursor)
	}

	buf := map[string]interface{}{}
	err := s.client.graphqlClient.Run(context.Background(), req, &buf)
	if err != nil {
		return nil, err
	}
	resp := IntegrationDefinitionsResponse{}
	err = mapstructure.Decode(buf["integrationDefinitions"], &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListInstances list the integration instances for the JupiterOne account.
// ListInstances returns a reference to an IntegrationInstanceResponse which
// contains the Instances and the PageInfo used to request additional instances.
// The first call to ListInstances should pass nil for the cursor. To paginate
// through all instances, the caller should check if PageInfo.HasNextPage
// is true and pass the PageInfo.Cursor from the response in subsequent calls.
func (s *IntegrationService) ListInstances(cursor *string) (*IntegrationInstanceResponse, error) {
	req := s.client.prepareRequest(`
			query IntegrationInstances($cursor: String) {
				integrationInstances(cursor: $cursor) {
					instances {
						id
						name
						integrationDefinitionId
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}`)

	if cursor != nil {
		req.Var("cursor", *cursor)
	}

	buf := map[string]interface{}{}
	err := s.client.graphqlClient.Run(context.Background(), req, &buf)
	if err != nil {
		return nil, err
	}
	resp := IntegrationInstanceResponse{}
	err = mapstructure.Decode(buf["integrationInstances"], &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}
