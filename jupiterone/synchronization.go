package jupiterone

import (
	"bytes"
	"encoding/json"
)

type SynchronizationService service

type StartParams struct {
	Source           string `json:"source"`
	Scope            string `json:"scope"`
	SyncMode         string `json:"syncMode"`
	InstanceId       string `json:"integrationInstanceId"`
	IgnoreDuplicates bool   `json:"-"`
}

type SynchronizationJobStatus struct {
	Source                   string
	Scope                    string
	ID                       string
	Status                   string
	StartTimestamp           string
	NumEntitiesUploaded      int
	NumEntitiesCreated       int
	NumEntitiesUpdated       int
	NumEntitiesDeleted       int
	NumRelationshipsUploaded int
	NumRelationshipsCreated  int
	NumRelationshipsUpdated  int
	NumRelationshipsDeleted  int
}

func (s *SynchronizationService) Start(params StartParams) (*SynchronizationJobStatus, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewBuffer(body)

  url := s.client.httpBaseUrl + "/persister/synchronization/jobs"
	resp, err := s.client.httpClient.Post(url, "application/json", bodyReader)
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

func (s *SynchronizationService) Finalize(id string) (*SynchronizationJobStatus, error) {
  url := s.client.httpBaseUrl + "/persister/synchronization/jobs/" + id + "/finalize"
	resp, err := s.client.httpClient.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	syncJobStatus := &SynchronizationJobStatus{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(syncJobStatus)
	if err != nil {
		return nil, err
	}
	return syncJobStatus, nil

}
