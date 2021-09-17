package jupiterone

import (
	"context"

	"github.com/mitchellh/mapstructure"
)

type RelationshipService service



type RelationshipProperties struct {
	RelationshipKey       	string                       `json:"key"`
	RelationshipType 		string                       `json:"type"`
	RelationshipClass 		string                       `json:"class"`
	FromEntityId string 					 `json:"fromEntityId"`
	ToEntityId string 					 `json:"toEntityId"`
	Properties  map[string]string  `json:"properties" bson:"properties,omitempty"`

}
type EdgeProperties struct {
	Id          string                       `json:"id"`
	ToVertexId          string                       `json:"toVertexId"`
	FromVertexId          string                       `json:"fromVertexId"`
	Relationship          map[string]string                       `json:"relationship"`
	Properties  map[string]string  `json:"properties" bson:"properties,omitempty"`
}

type Relationship struct {
	Relationship          map[string]string                       `json:"relationship"`
	Edge       EdgeProperties                       `json:"_fromVertexId"`
}



func (s *RelationshipService) Create(properties RelationshipProperties) (*Relationship, error) {
	req := s.client.prepareRequest(`
	mutation CreateRelationship(
		$relationshipKey: String!
		$relationshipType: String!
		$relationshipClass: String!
		$fromEntityId: String!
		$toEntityId: String!
		$properties: JSON
	  ) {
		createRelationship(
		  relationshipKey: $relationshipKey
		  relationshipType: $relationshipType
		  relationshipClass: $relationshipClass
		  fromEntityId: $fromEntityId
		  toEntityId: $toEntityId
		  properties: $properties
		) {
		  relationship {
			_id
		  }
		  edge {
			id
			toVertexId
			fromVertexId
			relationship {
			  _id
			}
			properties
		  }
		}
	  }
	`)

	req.Var("question", properties)

	var respData map[string]interface{}

	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	var relationship Relationship

	if err := mapstructure.Decode(respData["createRelationship"], &relationship); err != nil {
		return nil, err
	}

	return &relationship, nil
}



func (s *RelationshipService) Delete(id string ) error {
	req := s.client.prepareRequest(`
	mutation DeleteRelationship($relationshipId: String! $timestamp: Long) {
		deleteRelationship (relationshipId: $relationshipId, timestamp: $timestamp) {
		  relationship {
			_id
		  }
		  edge {
			id
			toVertexId
			fromVertexId
			relationship {
			  _id
			}
			properties
		  }
		}
	  }
	`)

	req.Var("id", id)

	if err := s.client.graphqlClient.Run(context.Background(), req, nil); err != nil {
		return err
	}

	return nil
}
