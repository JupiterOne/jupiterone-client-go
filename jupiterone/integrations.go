package jupiterone

import (
	"context"
	"encoding/json"
)

// IntegrationService handles the integration-related API requests.
type IntegrationService service

// IntegrationInstance represents an instance of an integration.
// Example: The account has an instance of the AWS integration.
type IntegrationInstance struct {
	ID                      string          `json:"id,omitempty"`
	Name                    string          `json:"name"`
	Description             string          `json:"description"`
	IntegrationDefinitionID string          `json:"integrationDefinitionId"`
	PollingInterval         string          `json:"pollingInterval,omitempty"`
	Config                  json.RawMessage `json:"config,omitempty"`
}

// IntegrationInstancesResponse is a slice of integration instances and
// pagination information.
type IntegrationInstancesResponse struct {
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
	ID               string          `json:"id"`
	IntegrationType  string          `json:"integrationType"`
	IntegrationClass []string        `json:"integrationClass"`
	Name             string          `json:"name"`
	Title            string          `json:"title"`
	RepoWebLink      string          `json:"repoWebLink"`
	ConfigFields     json.RawMessage `json:"configFields"`
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

// InvokeInstanceResult is the result of an InvokeInstance mutation.
type InvokeInstanceResult struct {
	Success bool `json:"success"`
}

// IntegrationJobsResponse represents the response from listing integration jobs.
// It contains a slice of IntegrationJobs and PageInfo.
type IntegrationJobsResponse struct {
	Jobs     []*IntegrationJob `json:"jobs"`
	PageInfo PageInfo          `json:"pageInfo"`
}

// IntegrationJob represent a single integration job.
type IntegrationJob struct {
	ID                    string `json:"id"`
	CreateDate            int    `json:"createDate"`
	EndDate               int    `json:"endDate"`
	ErrorsOccurred        bool   `json:"errorsOccurred"`
	Status                string `json:"status"`
	IntegrationInstanceID string `json:"integrationInstanceId"`
}

// IntegrationEventsResponse is the response from listing integration job events.
// It contains a slice of IntegrationJobEvents and PageInfo.
type IntegrationJobEventsResponse struct {
	Events   []*IntegrationJobEvent `json:"events"`
	PageInfo PageInfo               `json:"pageInfo"`
}

// IntegrationJobEvent represents a single event in an integration job's
// execution.
type IntegrationJobEvent struct {
	ID          string `json:"id"`
	JobID       string `json:"jobId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreateDate  int    `json:"createDate"`
}

// DeleteInstanceResult is the result of a DeleteInstance mutation.
type DeleteInstanceResult struct {
	Success bool `json:"success"`
}

// ListDefinitions lists all the IntegrationDefinitions in the current account.
// ListDefinitions returns a reference to an IntegrationDefinitionResponse object
// which contains the Definitions and PageInfo used to request additional
// definitions (if they exist).
//
// The first call to ListDefinitions should pass an empty string as the cursor. The caller
// should check PageInfo.HasNextPage to see if there is additional data available.
// If there is, the caller should pass PageInfo.Cursor on subsequent calls.
func (s *IntegrationService) ListDefinitions(cursor string) (*IntegrationDefinitionsResponse, error) {
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
					configFields {
						key
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}`)

	if cursor != "" {
		req.Var("cursor", cursor)
	}

	resp := struct {
		IntegrationDefinitionsResponse *IntegrationDefinitionsResponse `json:"integrationDefinitions"`
	}{
		IntegrationDefinitionsResponse: &IntegrationDefinitionsResponse{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.IntegrationDefinitionsResponse, nil
}

// GetDefinition gets a single Integration Definition by its id.
func (s *IntegrationService) GetDefinition(id string) (*IntegrationDefinition, error) {
	req := s.client.prepareRequest(`
    query getIntegrationDefinition($id: String) {
      integrationDefinition(id: $id) {
        id
        integrationType
        integrationClass
        name
        title
        repoWebLink
      }
  }`)

	req.Var("id", id)
	resp := struct {
		IntegrationDefinition *IntegrationDefinition `json:"integrationDefinition"`
	}{
		IntegrationDefinition: &IntegrationDefinition{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.IntegrationDefinition, nil
}

// ListInstances list the integration instances for the JupiterOne account.
// ListInstances returns a reference to an IntegrationInstanceResponse which
// contains the Instances and the PageInfo used to request additional instances.
// The first call to ListInstances should pass nil for the cursor. To paginate
// through all instances, the caller should check if PageInfo.HasNextPage
// is true and pass the PageInfo.Cursor from the response in subsequent calls.
func (s *IntegrationService) ListInstances(cursor string) (*IntegrationInstancesResponse, error) {
	req := s.client.prepareRequest(`
			query IntegrationInstances($cursor: String) {
				integrationInstances(cursor: $cursor) {
					instances {
						id
						name
						description
						integrationDefinitionId
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}`)

	if cursor != "" {
		req.Var("cursor", cursor)
	}

	resp := struct {
		IntegrationInstancesResponse *IntegrationInstancesResponse `json:"integrationInstances"`
	}{
		IntegrationInstancesResponse: &IntegrationInstancesResponse{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.IntegrationInstancesResponse, err
}

// CreateIntegrationInstance creates a new integration instance.
func (s *IntegrationService) CreateInstance(instance IntegrationInstance) (*IntegrationInstance, error) {
	req := s.client.prepareRequest(`
		mutation CreateInstance($instance: CreateIntegrationInstanceInput!) {
			createIntegrationInstance(instance: $instance)  {
				id
				name
				description
				pollingInterval
				integrationDefinitionId
				description
				config
			}
		}`)

	req.Var("instance", instance)

	resp := struct {
		CreateIntegrationInstance *IntegrationInstance `json:"createIntegrationInstance"`
	}{
		CreateIntegrationInstance: &IntegrationInstance{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.CreateIntegrationInstance, err
}

// InvokeInstance invokes an IntegrationInstance by its id.
func (s *IntegrationService) InvokeInstance(id string) (*InvokeInstanceResult, error) {
	req := s.client.prepareRequest(`
		mutation InvokeInstance($id: String!) {
			invokeIntegrationInstance(id: $id) {
				success
			}
		}`)

	req.Var("id", id)

	resp := struct {
		InvokeIntegrationInstance *InvokeInstanceResult `json:"invokeIntegrationInstance"`
	}{
		InvokeIntegrationInstance: &InvokeInstanceResult{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.InvokeIntegrationInstance, nil
}

// ListInstanceJobs lists the jobs for a specific integration with InstanceId, id.
// On the first call, a caller should pass the id of the integration, an empty cursor (""), and the size
// of the results.
//
// To paginate, the caller should first check PageInfo.HasNextPage. If true, then the caller should pass
// PageInfo.EndCursor as the cursor parameter.
//
// For default API response size behavior, pass 0 as the size.
func (s *IntegrationService) ListInstanceJobs(id string, cursor string, size int) (*IntegrationJobsResponse, error) {
	req := s.client.prepareRequest(`
		query ListJobs($integrationInstanceId: String!, $cursor: String, $size: Int) {
			integrationJobs(
				integrationInstanceId: $integrationInstanceId
				cursor: $cursor
				size: $size
			) {
				jobs {
					id
					createDate
					endDate
					errorsOccurred
					status
					integrationInstanceId
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
		}`)

	req.Var("integrationInstanceId", id)
	if cursor != "" {
		req.Var("cursor", cursor)
	}

	if size != 0 {
		req.Var("size", size)
	}

	resp := struct {
		IntegrationJobsResponse *IntegrationJobsResponse `json:"integrationJobs"`
	}{
		IntegrationJobsResponse: &IntegrationJobsResponse{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.IntegrationJobsResponse, nil
}

// ListJobEvents lists the events for a single integration job. On the first call
// the caller should pass the instance id of the integration, the integration job id,
// an empty cursor (""), and the size to limit the number of events returned (0 is default api behavior).
//
// To paginate, the caller should first check PageInfo.HasNextPage. If true, then the caller should pass
// PageInfo.EndCursor as the cursor parameter.
func (s *IntegrationService) ListJobEvents(instanceID string, jobID string, cursor string, size int) (*IntegrationJobEventsResponse, error) {
	req := s.client.prepareRequest(`
		query ListEvents (
			$jobId: String!,
			$integrationInstanceId: String!,
			$cursor: String,
			$size: Int,
		) {
			integrationEvents (
				size: $size,
				cursor: $cursor,
				jobId: $jobId,
				integrationInstanceId: $integrationInstanceId,
			) {
				events {
					id
					jobId
					name
					description
					createDate
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
		}
	`)

	resp := struct {
		IntegrationJobEventResponse *IntegrationJobEventsResponse `json:"integrationEvents"`
	}{
		IntegrationJobEventResponse: &IntegrationJobEventsResponse{},
	}

	req.Var("integrationInstanceId", instanceID)
	req.Var("jobId", jobID)
	if cursor != "" {
		req.Var("cursor", cursor)
	}

	if size != 0 {
		req.Var("size", size)
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.IntegrationJobEventResponse, nil
}

// DeleteInstance deletes an integration instance by its id.
func (s *IntegrationService) DeleteInstance(id string) (*DeleteInstanceResult, error) {
	req := s.client.prepareRequest(`
		mutation DeleteIntegrationInstance($id: String!) {
			deleteIntegrationInstance(id: $id) {
				success
			}
		}`)

	req.Var("id", id)

	resp := struct {
		DeleteInstanceResult *DeleteInstanceResult `json:"deleteIntegrationInstance"`
	}{
		DeleteInstanceResult: &DeleteInstanceResult{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.DeleteInstanceResult, nil
}
