package jupiterone

import (
	"context"
	"encoding/json"
)

// EntityService is the service for creating, reading, updating,
// and deleting Entities in the JupiterOne graph.
type EntityService service

type EntityProperties struct {
	Key        string `json:"key"`
	Type       string `json:"type"`
	Class      string `json:"class"`
	Properties map[string]interface{}
}

type EntityRawDataResponse struct {
	EntityID string          `json:"entityId"`
	Payload  []EntityRawData `json:"payload"`
}

type EntityRawData struct {
	ContentType string           `json:"contentType"`
	Name        string           `json:"name"`
	JSONData    *json.RawMessage `json:"JSONData"`
	TextData    string           `json:"TextData"`
}

// GetRawData gets the default raw data an entity was created from. EntityId is the _id property of an entity.
func (s *EntityService) GetDefaultRawData(entityID string) (*EntityRawDataResponse, error) {
	req := s.client.prepareRequest(`
		query GetEntityRawData($entityId: String!, $source: String!, $name: String, $versionId: String)	 {
			entityRawDataLegacy(entityId: $entityId, source: $source, name: $name, versionId: $versionId) {
				entityId
				payload {
					... on RawDataJSONEntityLegacy {
						contentType
						name
						JSONData: data
					}
					... on RawDataTextEntityLegacy {
						contentType
						name
						TextData: data
					}
				}
			}
		}
	`)

	req.Var("entityId", entityID)
	// Source is required, but has no effect for default raw data.
	req.Var("source", "")

	resp := struct {
		EntityRawDataLegacy *EntityRawDataResponse `json:"entityRawDataLegacy"`
	}{
		EntityRawDataLegacy: &EntityRawDataResponse{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.EntityRawDataLegacy, nil
}

// GetRawData gets the default raw data an entity was created from. EntityId is the _id property of an entity.
// Passing 0 as the version will get the latest raw data.
func (s *EntityService) GetRawData(entityID string, name string, version int) (*EntityRawDataResponse, error) {
	req := s.client.prepareRequest(`
		query GetEntityRawData($entityId: String!, $source: String!, $name: String, $versionId: String)	 {
			entityRawDataLegacy(entityId: $entityId, source: $source, name: $name, versionId: $versionId) {
				entityId
				payload {
					... on RawDataJSONEntityLegacy {
						contentType
						name
						JSONData: data
					}
					... on RawDataTextEntityLegacy {
						contentType
						name
						TextData: data
					}
				}
			}
		}
	`)

	req.Var("entityId", entityID)
	// Source is required, but has no effect for default raw data.
	req.Var("source", "")
	req.Var("name", name)
	if version != 0 {
		req.Var("version", version)
	}

	resp := struct {
		EntityRawDataLegacy *EntityRawDataResponse `json:"entityRawDataLegacy"`
	}{
		EntityRawDataLegacy: &EntityRawDataResponse{},
	}

	err := s.client.graphqlClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.EntityRawDataLegacy, nil
}

// Create creates a new entity in the JupiterOne graph with
// the _key, _type, _class, and properties in the entity argument.
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
