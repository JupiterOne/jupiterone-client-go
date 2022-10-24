package jupiterone

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type SynchronizationService service

type StartParams struct {
	Source           string `json:"source,omitempty"`
	Scope            string `json:"scope,omitempty"`
	SyncMode         string `json:"syncMode,omitempty"`
	InstanceID       string `json:"integrationInstanceId,omitempty"`
	IgnoreDuplicates bool   `json:"-"`
}

type SynchronizationJobStatus struct {
	Source         string
	Scope          string
	AccountID      string
	ID             string `json:"id"`
	Status         string
	StartTimestamp int
	DurationMs     int
	DeletionMode   string
	Done           bool
	TTL            int

	NumEntitiesUploaded         int
	NumStreamedEntitiesUploaded int
	NumEntitiesToDelete         int
	NumEntitiesCreated          int
	NumEntitiesUpdated          int
	NumEntitiesDeleted          int

	NumEntityCreateErrors int
	NumEntityUpdateErrors int
	NumEntityDeleteErrors int

	NumEntityRawDataEntriesUploaded int
	NumEntityRawDataEntriesCreated  int
	NumEntityRawDataEntriesUpdated  int
	NumEntityRawDataEntriesDeleted  int

	NumRelationshipsUploaded         int
	NumStreamedRelationshipsUploaded int
	NumRelationshipsToDelete         int
	NumRelationshipsCreated          int
	NumRelationshipsUpdated          int
	NumRelationshipsDeleted          int

	NumRelationshipCreateErrors int
	NumRelationshipUpdateErrors int
	NumRelationshipDeleteErrors int

	NumRelationshipRawDataEntriesUploaded int
	NumRelationshipRawDataEntriesCreated  int
	NumRelationshipRawDataEntriesUpdated  int
	NumRelationshipRawDataEntriesDeleted  int

	NumRelationshipRawDataEntryCreateErrors int
	NumRelationshipRawDataEntryUpdateErrors int
	NumRelationshipRawDataEntryDeleteErrors int

	NumMappedRelationshipsCreated int
	NumMappedRelationshipsUpdated int
	NumMappedRelationshipsDeleted int

	NumMappedRelationshipCreateErrors int
	NumMappedRelationshipUpdateErrors int
	NumMappedRelationshipDeleteErrors int

	NumMutationsSubmitted int
	NumMutationsCompleted int
}

func (s *SynchronizationService) Start(params StartParams) (*SynchronizationJobStatus, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewBuffer(body)

	url := s.client.httpBaseURL + "/persister/synchronization/jobs"
	return s.syncHelper(url, http.MethodPost, bodyReader)
}

func (s *SynchronizationService) Status(id string) (*SynchronizationJobStatus, error) {
	url := s.client.httpBaseURL + "/persister/synchronization/jobs/" + id
	return s.syncHelper(url, http.MethodGet, nil)
}

func (s *SynchronizationService) Finalize(id string) (*SynchronizationJobStatus, error) {
	url := s.client.httpBaseURL + "/persister/synchronization/jobs/" + id + "/finalize"
	return s.syncHelper(url, http.MethodPost, nil)
}

func (s *SynchronizationService) Upload(id string, data []byte) (*SynchronizationJobStatus, error) {
  url := s.client.httpBaseURL + "/persister/synchronization/jobs/" + id + "/upload"
  body := bytes.NewBuffer(data)

  return s.syncHelper(url, http.MethodPost, body)
}

func (s *SynchronizationService) UploadEntities(id string, data []byte) (*SynchronizationJobStatus, error) {
  url := s.client.httpBaseURL + "/persister/synchronization/jobs/" + id + "/entities"
  body := bytes.NewBuffer(data)

  return s.syncHelper(url,http.MethodPost, body)
}

func (s *SynchronizationService) UploadRelationships(id string, data []byte) (*SynchronizationJobStatus, error) {
  url := s.client.httpBaseURL + "/persister/synchronization/jobs/" + id + "/relationships"
  body := bytes.NewBuffer(data)

  return s.syncHelper(url, http.MethodPost, body)
}


func (s *SynchronizationService) syncHelper(url string, method string, body io.Reader) (*SynchronizationJobStatus, error) {
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

	syncJobStatus := struct {
		SyncJobStatus *SynchronizationJobStatus `json:"job"`
	}{
		SyncJobStatus: &SynchronizationJobStatus{},
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&syncJobStatus)
	if err != nil {
		return nil, err
	}

	return syncJobStatus.SyncJobStatus, nil
}
