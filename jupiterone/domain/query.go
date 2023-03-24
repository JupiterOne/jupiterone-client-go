package domain

type entity struct {
	Type                    []string    `json:"_type"`
	Deleted                 bool        `json:"_deleted"`
	Version                 int         `json:"_version"`
	CreatedOn               string      `json:"_createdOn"`
	DisplayName             string      `json:"displayName"`
	IntegrationName         string      `json:"_integrationName"`
	IntegrationType         string      `json:"_integrationType"`
	IntegrationClass        interface{} `json:"_integrationClass"`
	Source                  string      `json:"_source"`
	Scope                   string      `json:"_scope"`
	AccountID               string      `json:"_accountId"`
	Key                     string      `json:"_key"`
	Class                   []string    `json:"_class"`
	BeginOn                 string      `json:"_beginOn"`
	IntegrationInstanceID   interface{} `json:"_integrationInstanceId"`
	IntegrationDefinitionID string      `json:"_integrationDefinitionId"`
	ID                      string      `json:"_id"`
}

type relationship struct {
	ToEntityKey   string `json:"_toEntityKey"`
	Version       int    `json:"_version"`
	Key           string `json:"_key"`
	FromEntityID  string `json:"_fromEntityId"`
	ID            string `json:"_id"`
	FromEntityKey string `json:"_fromEntityKey"`
	CreatedOn     string `json:"_createdOn"`
	Deleted       bool   `json:"_deleted"`
	BeginOn       string `json:"_beginOn"`
	Source        string `json:"_source"`
	ToEntityID    string `json:"_toEntityId"`
	Class         string `json:"_class"`
	Scope         string `json:"_scope"`
	DisplayName   string `json:"displayName"`
	Type          string `json:"_type"`
	AccountID     string `json:"_accountId"`
}

type QueryDataVertex struct {
	Entity     entity                 `json:"entity"`
	ID         string                 `json:"id"`
	Properties map[string]interface{} `json:"properties"`
}

type QueryDataEdge struct {
	ID           string       `json:"id"`
	ToVertexID   string       `json:"toVertexId"`
	FromVertexID string       `json:"fromVertexId"`
	Relationship relationship `json:"relationship"`
	Properties   interface{}  `json:"properties"`
}

type QueryDataTreeResultFormat struct {
	Vertices []QueryDataVertex `json:"vertices"`
	Edges    []QueryDataEdge   `json:"edges"`
}

type QueryFormat interface {
	[]QueryDataVertex | QueryDataTreeResultFormat | []interface{}
}

type QueryResult[T QueryFormat] struct {
	Type   string
	Data   T
	Cursor string
}

type DeferredQueryURLResponse struct {
	URL    string
	Status string
}
