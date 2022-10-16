package jupiterone

import (
	"context"
	"encoding/json"
)

// AuditService provides the Audit API functions.
type AuditService service

// ListAuditEventsResponse represents the response from
// a GetAuditEventsForAccount query. It contains a slice
// of AuditEvents and pagination info in PageInfo.
type ListAuditEventsResponse struct {
	Items    []*AuditEvent `json:"items"`
	PageInfo struct {
		HasNextPage bool   `json:"hasNextPage"`
		Cursor      string `json:"endCursor"`
	}
}

// Audit Event represents a single Audit event in a JupiterOne account.
type AuditEvent struct {
	ID                string
	ResourceType      string
	ResourceID        string
	Category          string
	Timestamp         uint
	PerformedByUserID string
	Data              json.RawMessage
}

// ListAuditEvents lists the audit events in the JupiterOne account.
//
// The limit for items returned can be set through the limit parameter.
// A limit of 0 uses the default behavior of the API.
//
// The first call should use an empty string for cursor.
// To paginate through all items, first check if PageInfo.HasNextPage
// is true. If it is, call ListAuditEvents again and pass the
// PageInfo.Cursor as the cursor parameter.
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
