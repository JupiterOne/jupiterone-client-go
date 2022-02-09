package jupiterone

import (
	"context"
	"encoding/json"
	"fmt"
)

type EntityService service

type EntityProperties struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Class string `json:"class"`
}

type Entity struct {
	ID string `json:"id"`
}

func (s *EntityService) Create(properties EntityProperties) (*string, error) {
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

	req.Var("entityKey", properties.Key)
	req.Var("entityType", properties.Type)
	req.Var("entityClass", properties.Class)

	var respData map[string]interface{}
	// var respData string

	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	resp, nil := json.Marshal(respData)
	fmt.Println("Entity: " + fmt.Sprint(respData))
	respString := string(resp)
	return &respString, nil
}
