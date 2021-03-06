package jupiterone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGraphqlEndpointSetsDefaultValue(t *testing.T) {
	config := Config{}
	endpoint := config.getGraphQLEndpoint()
	assert.Equal(t, endpoint, "https://api.us.jupiterone.io/graphql", "Endpoints should match")
}

func TestGetGraphqlEndpointNoOverride(t *testing.T) {
	config := Config{
		Region: "dev",
	}

	endpoint := config.getGraphQLEndpoint()
	assert.Equal(t, endpoint, "https://api.dev.jupiterone.io/graphql", "Endpoints should match")
}
