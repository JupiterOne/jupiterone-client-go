package jupiterone

import (
	"context"
	"encoding/json"
)

// EntityService is the service for creating, reading, updating,
// and deleting Entities in the JupiterOne graph
type EntityService service

type EntityProperties struct {
	Key        string `json:"key"`
	Type       string `json:"type"`
	Class      string `json:"class"`
	Properties map[string]interface{}
}

// Create creates a new entity in the JupiterOne graph with
// the _key, _type, _class, and properties in the entity arguement
func (s *EntityService) Create(entity EntityProperties) (*string, error) {
	req := s.client.prepareRequest(`
	mutation CreateEntity(
		$entityKey: String!
		$entityType: String!
		$entityClass: [String!]!
		$properties: JSON
	  ) {
		createEntity(
		  entityKey: $entityKey
		  entityType: $entityType
		  entityClass: $entityClass
		  properties: $properties
		) {
		  entity {
			_id
		  }
		  vertex {
			id
			entity {
			  _id
			}
		  }
		}
	  }
	`)

	req.Var("entityKey", entity.Key)
	req.Var("entityType", entity.Type)
	req.Var("entityClass", entity.Class)
	req.Var("properties", entity.Properties)

	var respData map[string]interface{}

	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	resp, nil := json.Marshal(respData)
	respString := string(resp)
	return &respString, nil
}
