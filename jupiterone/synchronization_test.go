package jupiterone

import (
	"fmt"
	"testing"

	"github.com/jupiterone/jupiterone-client-go/jupiterone/domain"
)

type syncTests struct {
	title      string
	iterations int
}

func createFakeData(iterations int) []interface{} {
	output := []interface{}{}
	for i := 0; i < iterations; i++ {
		output = append(output, "hi!")
	}
	return output
}

func TestChunkUpload(t *testing.T) {
	tests := []syncTests{
		{
			title:      "large number of iterations",
			iterations: 1000000,
		},
		{
			title:      "zero iterations",
			iterations: 0,
		},
		{
			title:      "fair amount of iterations",
			iterations: 100,
		},
	}

	for _, tv := range tests {
		fmt.Printf("Running %s\n", tv.title)

		c := &Config{
			APIKey:     "a",
			AccountID:  "a",
			Region:     "a",
			HTTPClient: nil,
		}
		client, err := NewClient(c)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		fakeData := createFakeData(tv.iterations)
		fakeUploadFn := func(id string, data domain.SyncPayload) (*domain.SynchronizationJobOutput, error) {
			return &domain.SynchronizationJobOutput{}, nil
		}
		fakeMarshalFn := func(smtg []interface{}) domain.SyncPayload {
			return domain.SyncPayload{
				Entities: smtg,
			}
		}
		fns := chunkUploadFunctions{
			marshalPayload: fakeMarshalFn,
			upload:         fakeUploadFn,
		}

		err = client.Synchronization.chunkUpload("a", fakeData, fns)
		if err != nil {
			t.Fatalf("failed to chunk: %v", err)
		}
	}
}
