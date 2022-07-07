package jupiterone

import (
	"context"
	"encoding/json"
)

// AlertService is the service for creating, reading, updating and
// deleting Alert rules
type AlertService service

type Alert struct {
}

func (s *AlertService) Create() {}

// Get retrieves an Alert by its id
func (s *AlertService) Get(id string) (*Alert, error) {
	req := s.client.prepareRequest(`
    query AlertInstanceById($id: ID) {
      alertInstance(id: $id) {
        id
        accountId
        ruleId
        level
        status
        lastUpdatedOn
        lastEvaluationBeginOn
        lastEvaluationEndOn
        createdOn
        dismissedOn
        lastEvaluationResult {
          rawDataDescriptors {
            name
            persistedResultType
            recordCount
          }
          outputs {
            name
            value
          }
        }
        questionRuleInstance {
          ...QuestionRuleInstanceWithoutTagsFragment
        }
        reportRuleInstance {
          ...ReportRuleInstanceFragment
        }
      }
    }`)

	req.Var("id", id)
	resp := make([]json.RawMessage, 0)
	err := s.client.common.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}

	return &Alert{}, nil
}

// TODO: add filter by status and option to change limit
func (s *AlertService) List() ([]*Alert, error) {
	req := s.client.prepareRequest(`
    query ListAlertInstances($alertStatus: AlertStatus, $limit: Int, $cursor: String) {
      listAlertInstances(alertStatus: $alertStatus, limit: $limit, cursor: $cursor) {
        instances {
          id
          accountId
          ruleId
          level
          status
          lastUpdatedOn
          lastEvaluationBeginOn
          lastEvaluationEndOn
          createdOn
          dismissedOn
          lastEvaluationResult {
            rawDataDescriptors {
              recordCount
            }
          }
          questionRuleInstance {
            id
            name
            description
          }
        }
        pageInfo { endCursor hasNextPage }
      }
    }`)

	req.Var("limit", 10)

	var respData []*Alert
	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	return respData, nil
}

func (s *AlertService) Update() {}

func (s *AlertService) Delete() {}
