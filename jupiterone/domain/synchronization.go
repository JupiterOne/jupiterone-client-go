package domain

type StartParams struct {
	Source           string `json:"source,omitempty"`
	Scope            string `json:"scope,omitempty"`
	SyncMode         string `json:"syncMode,omitempty"`
	InstanceID       string `json:"integrationInstanceId,omitempty"`
	IgnoreDuplicates bool   `json:"-"`
}

type SynchronizationJobOutput struct {
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

type SyncPayload struct {
	Entities      []interface{} `json:"entities,omitempty"`
	Relationships []interface{} `json:"relationships,omitempty"`
}
