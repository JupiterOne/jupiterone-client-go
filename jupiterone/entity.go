package jupiterone

import (
	"context"
	"fmt"

	graphql "github.com/hasura/go-graphql-client"
	// "github.com/mitchellh/mapstructure"
)

type EntityService service

type EntityProperties struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Class string `json:"class"`
	// Properties []QuestionQuery              `json:"queries"`
}

type Entity struct {
	Id string `json:"id"`
}

func (s *EntityService) Create(properties EntityProperties) (*Entity, error) {
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
		} `graphql:"createEntity(entityKey: $entityKey
			entityType: $entityType
			entityClass: $entityClass
			properties: $properties)"`
	}
	variables := map[string]interface{}{
		"entityKey":   properties.Key,
		"entityType":  properties.Type,
		"entityClass": properties.Class,
	}

	req.Var("entityKey", properties.Key)
	req.Var("entityType", properties.Type)
	req.Var("entityClass", properties.Class)

	err := s.client.graphqlClient.Mutate(context.Background(), &m, variables)
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Created entity: %v", m.CreateEntity.Id)

	// var respData map[string]interface{}
	// //var respData string

	// if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
	// 	return nil, err
	// }

	// var entity Entity

	// if err := mapstructure.Decode(respData["createEntity"], &entity); err != nil {
	// 	return nil, err
	// }
	// fmt.Println(respData["id"])

	// resp, nil := json.Marshal(respData)
	// fmt.Println("Entity: " + fmt.Sprint(respData))
	// respString := string(resp)
	// return &respString, nil
	return &entity, nil
}
