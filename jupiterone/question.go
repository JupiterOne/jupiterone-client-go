package jupiterone

import (
	"context"

	"github.com/mitchellh/mapstructure"
)

type QuestionService service

type QuestionQuery struct {
	Query   string `json:"query"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

type QuestionComplianceMetaData struct {
	Standard     string   `json:"standard"`
	Requirements []string `json:"requirements"`
	Controls     []string `json:"controls"`
}

type QuestionProperties struct {
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	Tags        []string                     `json:"tags"`
	Queries     []QuestionQuery              `json:"queries"`
	Compliance  []QuestionComplianceMetaData `json:"compliance"`
}

type Question struct {
	ID          string                       `json:"id"`
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	Tags        []string                     `json:"tags"`
	Queries     []QuestionQuery              `json:"queries"`
	Compliance  []QuestionComplianceMetaData `json:"compliance"`
}

func (s *QuestionService) Get(id string) (*Question, error) {
	req := s.client.prepareRequest(`
		query GetQuestionById ($id: ID!) {
			question(id: $id) {
				id
				title
				description
				queries {
					query
					version
				}
				tags
				compliance {
					type
					details {
						name
						description
					}
				}
				accountId
				integrationDefinitionId
			}
		}
	`)

	req.Var("id", id)

	var respData map[string]interface{}

	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	var question Question

	if err := mapstructure.Decode(respData["question"], &question); err != nil {
		return nil, err
	}

	return &question, nil
}

func (s *QuestionService) Create(properties QuestionProperties) (*Question, error) {
	req := s.client.prepareRequest(`
		mutation CreateQuestion($question: CreateQuestionInput!) {
			createQuestion(question: $question) {
				id
				title
				description
				queries {
					query
					version
				}
				tags
				variables {
					name
					required
					default
				}
				compliance {
					type
					details {
						name
						description
					}
				}
			}
		}
	`)

	req.Var("question", properties)

	var respData map[string]interface{}

	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	var question Question

	if err := mapstructure.Decode(respData["createQuestion"], &question); err != nil {
		return nil, err
	}

	return &question, nil
}

func (s *QuestionService) Update(id string, properties QuestionProperties) (*Question, error) {
	req := s.client.prepareRequest(`
		mutation UpdateQuestion ($id: ID!, $update: QuestionUpdate!) {
			updateQuestion(id: $id, update: $update) {
				id
				title
				description
				queries {
					query
					version
				}
				tags
				variables {
					name
					required
					default
				}
				compliance {
					type
					details {
						name
						description
					}
				}
			}
		}
	`)

	req.Var("id", id)
	req.Var("update", properties)

	var respData map[string]interface{}

	if err := s.client.graphqlClient.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	var question Question

	if err := mapstructure.Decode(respData["updateQuestion"], &question); err != nil {
		return nil, err
	}

	return &question, nil
}

func (s *QuestionService) Delete(id string) error {
	req := s.client.prepareRequest(`
		mutation DeleteQuestion($id: ID!) {
			deleteQuestion(id: $id) {
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
