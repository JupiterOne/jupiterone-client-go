package jupiterone

import (
	"context"
	"encoding/json"
	"fmt"
	// "github.com/mitchellh/mapstructure"
)

type EntityProperties struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Class string `json:"class"`
	// Properties []QuestionQuery              `json:"queries"`
}

type Entity struct {
	Id string `json:"id"`
}

func (c *JupiterOneClient) CreateEntity(properties EntityProperties) (*string, error) {
	req := c.prepareRequest(`
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
	//var respData string

	if err := c.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	// var entity Entity

	// if err := mapstructure.Decode(respData["createEntity"], &question); err != nil {
	// 	return nil, err
	// }

	resp, nil := json.Marshal(respData)
	fmt.Println("Entity: " + fmt.Sprint(respData))
	respString := string(resp)
	return &respString, nil
	//return &respData, nil
}
