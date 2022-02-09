package jupiterone

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mitchellh/mapstructure"
)

type RuleService service

type RuleQuestion struct {
	Queries []QuestionQuery `json:"queries"`
}

type RuleOperation struct {
	When    []map[string]interface{} `json:"when"`
	Actions []string                 `json:"actions"`
}

type QuestionRuleInstance struct {
	BaseQuestionRuleInstanceProperties
	ID        string `json:"id"`
	AccountID string `json:"accountId"`
	Version   int    `json:"version"`
	Latest    bool   `json:"latest"`
	Deleted   bool   `json:"deleted"`
	Type      string `json:"type"`
}

type UpdateQuestionRuleInstanceProperties struct {
	BaseQuestionRuleInstanceProperties
	ID      string `json:"id"`
	Version int    `json:"version"`
}

type BaseQuestionRuleInstanceProperties struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	SpecVersion     int                    `json:"specVersion"`
	PollingInterval string                 `json:"pollingInterval"`
	Outputs         []string               `json:"outputs"`
	Operations      string                 `json:"operations"`
	Question        RuleQuestion           `json:"question"`
	Templates       map[string]interface{} `json:"templates"`
}

type CreateQuestionRuleInstanceInput struct {
	BaseQuestionRuleInstanceProperties
	Operations []map[string]interface{} `json:"operations"`
}

type UpdateQuestionRuleInstanceInput struct {
	UpdateQuestionRuleInstanceProperties
	Operations []map[string]interface{} `json:"operations"`
}

// GetQuestionRuleInstanceByID - Fetches the QuestionRuleInstance by unique id.
func (s *RuleService) GetByID(id string) (*QuestionRuleInstance, error) {
	req := s.client.prepareRequest(`
		query GetQuestionRuleInstance($id: ID!) {
			questionRuleInstance (id: $id) {
				id
				name
				description
				version
				specVersion
				latest
				pollingInterval
				deleted
				accountId
				type
				templates
				question {
					queries {
						name
						query
						version
					}
				}
				operations {
					when
					actions
				}
				outputs
			}
		}
	`)

	req.Var("id", id)

	var respData map[string]interface{}
	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	var decodedQuestionRuleInstance QuestionRuleInstance
	err := mapstructure.Decode(respData["questionRuleInstance"], &decodedQuestionRuleInstance)
	if err != nil {
		return nil, err
	}

	return &decodedQuestionRuleInstance, nil
}

// CreateQuestionRuleInstance - Creates a question rule instance.
func (s *RuleService) Create(createQuestionRuleInstanceInput BaseQuestionRuleInstanceProperties) (*QuestionRuleInstance, error) {
	log.Println("Create question rule instance: " + createQuestionRuleInstanceInput.Name)

	req := s.client.prepareRequest(`
		mutation CreateQuestionRuleInstance ($instance: CreateQuestionRuleInstanceInput!) {
			createQuestionRuleInstance (
				instance: $instance
			) {
				id
				name
				description
				version
				specVersion
				latest
				deleted
				accountId
				type
				pollingInterval
				templates
				question {
					queries {
						name
						query
						version
					}
				}
				operations {
					when
					actions
				}
				outputs
			}
		}
	`)

	var input CreateQuestionRuleInstanceInput
	input.Name = createQuestionRuleInstanceInput.Name
	input.Description = createQuestionRuleInstanceInput.Description
	input.SpecVersion = createQuestionRuleInstanceInput.SpecVersion
	input.PollingInterval = createQuestionRuleInstanceInput.PollingInterval
	input.Outputs = createQuestionRuleInstanceInput.Outputs
	input.Question = createQuestionRuleInstanceInput.Question
	input.Templates = createQuestionRuleInstanceInput.Templates

	var deserializedOperationsMap []map[string]interface{}

	err := json.Unmarshal([]byte(createQuestionRuleInstanceInput.Operations), &deserializedOperationsMap)
	if err != nil {
		return nil, err
	}

	input.Operations = deserializedOperationsMap

	req.Var("instance", input)

	var respData map[string]interface{}

	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	var questionRuleInstance QuestionRuleInstance

	if err := mapstructure.Decode(respData["createQuestionRuleInstance"], &questionRuleInstance); err != nil {
		return nil, err
	}

	return &questionRuleInstance, nil
}

func (s *RuleService) Update(properties UpdateQuestionRuleInstanceProperties) (*QuestionRuleInstance, error) {
	log.Println("Updating question rule instance: " + properties.Name)

	req := s.client.prepareRequest(`
		mutation UpdateQuestionRuleInstance ($instance: UpdateQuestionRuleInstanceInput!) {
			updateQuestionRuleInstance (
				instance: $instance
			) {
				id
				name
				description
				version
				specVersion
				latest
				deleted
				accountId
				type
				pollingInterval
				templates
				question {
					queries {
						name
						query
						version
					}
				}
				operations {
					when
					actions
				}
				outputs
			}
		}
	`)

	var input UpdateQuestionRuleInstanceInput
	input.ID = properties.ID
	input.Version = properties.Version
	input.Name = properties.Name
	input.Description = properties.Description
	input.SpecVersion = properties.SpecVersion
	input.PollingInterval = properties.PollingInterval
	input.Outputs = properties.Outputs
	input.Question = properties.Question
	input.Templates = properties.Templates

	var deserializedOperationsMap []map[string]interface{}

	err := json.Unmarshal([]byte(properties.Operations), &deserializedOperationsMap)
	if err != nil {
		return nil, err
	}

	input.Operations = deserializedOperationsMap

	req.Var("instance", input)
	var respData map[string]interface{}

	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	var questionRuleInstance QuestionRuleInstance

	if err := mapstructure.Decode(respData["updateQuestionRuleInstance"], &questionRuleInstance); err != nil {
		return nil, err
	}

	return &questionRuleInstance, nil
}

func (s *RuleService) Delete(id string) error {
	req := s.client.prepareRequest(`
		mutation DeleteRuleInstance ($id: ID!) {
			deleteRuleInstance (id: $id) {
				id
			}
	      }
	`)

	req.Var("id", id)

	if err := s.client.graphqlClient.Run(context.Background(), req, nil); err != nil {
		return err
	}

	return nil
}
