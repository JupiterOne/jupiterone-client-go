package jupiterone

import (
	"context"
	"fmt"

	graphql "github.com/hasura/go-graphql-client"
	// "github.com/mitchellh/mapstructure"
)

type EntityService service

type EntityProperties struct {
	Key        graphql.String   `json:"key"`
	Type       graphql.String   `json:"type"`
	Class      []graphql.String `json:"class"`
	Properties graphql.String   `json:"properties"`
}

type Entity struct {
	Id string `json:"id"`
}

func (s *EntityService) Create(properties EntityProperties) (*string, error) {
	// req := s.client.prepareRequest(`
	// mutation CreateEntity(
	// 	$entityKey: String!
	// 	$entityType: String!
	// 	$entityClass: [String!]!
	// 	$properties: JSON
	//   ) {
	// 	createEntity(
	// 	  entityKey: $entityKey
	// 	  entityType: $entityType
	// 	  entityClass: $entityClass
	// 	  properties: $properties
	// 	) {
	// 	  entity {
	// 		_id
	// 	  }
	// 	  vertex {
	// 		id
	// 		entity {
	// 		  _id
	// 		}
	// 	  }
	// 	}
	//   }
	// `)

	var m struct {
		CreateEntity struct {
			Id graphql.String
		} `graphql:"createEntity(entityKey: $entityKey, entityType: $entityType, entityClass: $entityClass, properties: $properties)"`
	}
	variables := map[string]interface{}{
		"entityKey":   properties.Key,
		"entityType":  properties.Type,
		"entityClass": properties.Class,
		"properties":  properties.Properties,
	}

	err := s.client.graphqlClient.Mutate(context.Background(), &m, variables)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Created entity: %v", m.CreateEntity.Id)
	entityId := string(m.CreateEntity.Id)

	return &entityId, nil
}
