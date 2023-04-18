package jupiterone

import (
	"context"

	"github.com/jupiterone/jupiterone-client-go/jupiterone/graphql"
)

// IntegrationService handles the integration-related API requests.
type IntegrationService service

// ListDefinitions lists all the IntegrationDefinitions in the current account.
// ListDefinitions returns a reference to an IntegrationDefinitionResponse object
// which contains the Definitions and PageInfo used to request additional
// definitions (if they exist).
//
// The first call to ListDefinitions should pass an empty string as the cursor. The caller
// should check PageInfo.HasNextPage to see if there is additional data available.
// If there is, the caller should pass PageInfo.Cursor on subsequent calls.
func (s *IntegrationService) ListDefinitions(cursor string) (*graphql.IntegrationDefinitionsResponse, error) {
	return graphql.IntegrationDefinitions(context.Background(), s.client.gqlClient, cursor)
}

// GetDefinition gets a single Integration Definition by its id.
func (s *IntegrationService) GetDefinition(id string) (*graphql.GetIntegrationDefinitionResponse, error) {
	return graphql.GetIntegrationDefinition(context.Background(), s.client.gqlClient, id)
}

// ListInstances list the integration instances for the JupiterOne account.
// ListInstances returns a reference to ListIntegrationInstancesResponse which
// contains the Instances and the PageInfo used to request additional instances.
// The first call to ListInstances should pass nil for the cursor. To paginate
// through all instances, the caller should check if PageInfo.HasNextPage
// is true and pass the PageInfo.Cursor from the response in subsequent calls.
func (s *IntegrationService) ListInstances(cursor string) (*graphql.ListIntegrationInstancesResponse, error) {
	return graphql.ListIntegrationInstances(context.Background(), s.client.gqlClient, cursor)
}

// CreateAnIntegrationInstance creates a new integration instance.
func (s *IntegrationService) CreateInstance(instance graphql.CreateIntegrationInstanceInput) (*graphql.CreateInstanceResponse, error) {
	return graphql.CreateInstance(context.Background(), s.client.gqlClient, instance)
}

// InvokeInstance invokes an AnIntegrationInstance by its id.
func (s *IntegrationService) InvokeInstance(id string) (*graphql.InvokeInstanceResponse, error) {
	return graphql.InvokeInstance(context.Background(), s.client.gqlClient, id)
}

// ListInstanceJobs lists the jobs for a specific integration with InstanceId, id.
// On the first call, a caller should pass the id of the integration, an empty cursor (""), and the size
// of the results.
//
// To paginate, the caller should first check PageInfo.HasNextPage. If true, then the caller should pass
// PageInfo.EndCursor as the cursor parameter.
//
// For default API response size behavior, pass 0 as the size.
func (s *IntegrationService) ListInstanceJobs(id string, cursor string, size int) (*graphql.ListJobsResponse, error) {
	return graphql.ListJobs(context.Background(), s.client.gqlClient, id, cursor, size)
}

// ListJobEvents lists the events for a single integration job. On the first call
// the caller should pass the instance id of the integration, the integration job id,
// an empty cursor (""), and the size to limit the number of events returned (0 is default api behavior).
//
// To paginate, the caller should first check PageInfo.HasNextPage. If true, then the caller should pass
// PageInfo.EndCursor as the cursor parameter.
func (s *IntegrationService) ListJobEvents(instanceID string, jobID string, cursor string, size int) (*graphql.ListEventsResponse, error) {
	return graphql.ListEvents(context.Background(), s.client.gqlClient, jobID, instanceID, cursor, size)
}

// DeleteInstance deletes an integration instance by its id.
func (s *IntegrationService) DeleteInstance(id string) (*graphql.DeleteIntegrationInstanceResponse, error) {
	return graphql.DeleteIntegrationInstance(context.Background(), s.client.gqlClient, id)
}

func (s *IntegrationService) GetIntegrationInstance(id string) (*graphql.GetIntegrationInstanceResponse, error) {
	return graphql.GetIntegrationInstance(context.Background(), s.client.gqlClient, id)
}

func (s *IntegrationService) UpdateIntegrationInstance(id string, payload graphql.UpdateIntegrationInstanceInput) (*graphql.UpdateIntegrationInstanceResponse, error) {
	return graphql.UpdateIntegrationInstance(context.Background(), s.client.gqlClient, id, payload)
}
