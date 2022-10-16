package jupiterone

import (
	"context"
  "encoding/json"
)

type AuditService service

type ListAuditEventsResponse struct {
	Items    []*AuditEvent `json:"items"`
	PageInfo struct {
		HasNextPage bool   `json:"hasNextPage"`
		Cursor      string `json:"endCursor"`
	}
}

type AuditEvent struct {
	Id                string
	ResourceType      string
	ResourceId        string
	Category          string
	Timestamp         uint
	PerformedByUserId string
	Data              json.RawMessage
}

func (as *AuditService) ListAuditEvents(limit int, cursor string) (*ListAuditEventsResponse, error) {
	req := as.client.prepareRequest(`
    query GetAuditEventsForAccount($limit: Int, $cursor: String) {
      getAuditEventsForAccount(limit: $limit, cursor: $cursor) {
        items {
          id
          resourceType
          resourceId
          category
          timestamp
          performedByUserId
          data
        }
        pageInfo {
          endCursor
          hasNextPage
        }
      }
    }`)

	if limit != 0 {
		req.Var("limit", limit)
	}

	if cursor != "" {
		req.Var("cursor", cursor)
	}

	resp := struct {
		GetAuditEventsForAccount *ListAuditEventsResponse `json:"getAuditEventsForAccount"`
	}{
		GetAuditEventsForAccount: &ListAuditEventsResponse{},
	}

	err := as.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.GetAuditEventsForAccount, nil

}
