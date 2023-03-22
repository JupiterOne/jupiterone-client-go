package jupiterone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jupiterone/jupiterone-client-go/jupiterone/domain"
)

type SynchronizationService service

const (
	syncAPIStartPath         = "%s/persister/synchronization/jobs"
	syncAPIUploadPath        = "%s/persister/synchronization/jobs/%s/upload"
	syncAPIFinalizePath      = "%s/persister/synchronization/jobs/%s/finalize"
	syncAPIStatusPath        = "%s/persister/synchronization/jobs/%s"
	syncAPIEntitiesPath      = "%s/persister/synchronization/jobs/%s/entities"
	syncAPIRelationshipsPath = "%s/persister/synchronization/jobs/%s/relationships"
)

func (s *SynchronizationService) Start(params domain.StartParams) (*domain.SynchronizationJobOutput, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewBuffer(body)

	url := fmt.Sprintf(syncAPIStartPath, s.client.httpBaseURL)
	return s.syncHelper(url, http.MethodPost, bodyReader)
}

func (s *SynchronizationService) Status(id string) (*domain.SynchronizationJobOutput, error) {
	url := fmt.Sprintf(syncAPIStatusPath, s.client.httpBaseURL, id)
	return s.syncHelper(url, http.MethodGet, nil)
}

func (s *SynchronizationService) Finalize(id string) (*domain.SynchronizationJobOutput, error) {
	url := fmt.Sprintf(syncAPIFinalizePath, s.client.httpBaseURL, id)
	return s.syncHelper(url, http.MethodPost, nil)
}

func (s *SynchronizationService) Upload(id string, data domain.SyncPayload) (*domain.SynchronizationJobOutput, error) {
	url := fmt.Sprintf(syncAPIUploadPath, s.client.httpBaseURL, id)
	dataAsBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(dataAsBytes)

	return s.syncHelper(url, http.MethodPost, body)
}

func (s *SynchronizationService) UploadEntities(id string, data []byte) (*domain.SynchronizationJobOutput, error) {
	url := fmt.Sprintf(syncAPIEntitiesPath, s.client.httpBaseURL, id)
	body := bytes.NewBuffer(data)

	return s.syncHelper(url, http.MethodPost, body)
}

func (s *SynchronizationService) UploadRelationships(id string, data []byte) (*domain.SynchronizationJobOutput, error) {
	url := fmt.Sprintf(syncAPIRelationshipsPath, s.client.httpBaseURL, id)
	body := bytes.NewBuffer(data)

	return s.syncHelper(url, http.MethodPost, body)
}

func (s *SynchronizationService) marshalEntities(entities []interface{}) domain.SyncPayload {
	return domain.SyncPayload{
		Entities: entities,
	}
}

func (s *SynchronizationService) marshalRelationships(relationships []interface{}) domain.SyncPayload {
	return domain.SyncPayload{
		Relationships: relationships,
	}
}

type chunkUploadFunctions struct {
	marshalPayload func([]interface{}) domain.SyncPayload
	upload         func(string, domain.SyncPayload) (*domain.SynchronizationJobOutput, error)
}

// chunkUpload breaks apart the payload into chunks and uploads them so that the user
// is protected from uploading data that is too large at one time.
func (s *SynchronizationService) chunkUpload(jobID string, payloadItems []interface{}, fns chunkUploadFunctions) error {
	interval := 150

	for len(payloadItems) != 0 {
		if interval > len(payloadItems) {
			interval = len(payloadItems)
		}

		chunk := payloadItems[:interval]
		syncPayload := fns.marshalPayload(chunk)
		_, err := fns.upload(jobID, syncPayload)
		if err != nil {
			return err
		}

		payloadItems = payloadItems[interval:]
	}

	return nil
}

// ProcessSyncJob is a helper function that will start, upload, and finalize a sync job.
func (s *SynchronizationService) ProcessSyncJob(sp domain.StartParams, data domain.SyncPayload) (*domain.SynchronizationJobOutput, error) {
	syncJob, err := s.Start(sp)
	if err != nil {
		return nil, err
	}

	entityChunkUploadFunctions := chunkUploadFunctions{
		marshalPayload: s.marshalEntities,
		upload:         s.Upload,
	}

	err = s.chunkUpload(syncJob.ID, data.Entities, entityChunkUploadFunctions)
	if err != nil {
		return nil, err
	}

	relationshipChunkUploadFunctions := chunkUploadFunctions{
		marshalPayload: s.marshalRelationships,
		upload:         s.Upload,
	}
	err = s.chunkUpload(syncJob.ID, data.Relationships, relationshipChunkUploadFunctions)
	if err != nil {
		return nil, err
	}

	_, err = s.Finalize(syncJob.ID)
	if err != nil {
		return nil, err
	}

	return s.Status(syncJob.ID)
}

func (s *SynchronizationService) syncHelper(url string, method string, body io.Reader) (*domain.SynchronizationJobOutput, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req = s.client.addAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	syncJobOutput := struct {
		SyncJobOutput *domain.SynchronizationJobOutput `json:"job"`
	}{
		SyncJobOutput: &domain.SynchronizationJobOutput{},
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&syncJobOutput)
	if err != nil {
		return nil, err
	}

	return syncJobOutput.SyncJobOutput, nil
}
