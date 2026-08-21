package store

import (
	"encoding/json"
	"errors"
	"sort"

	"engineering-document-vault/internal/domain"
)

type Snapshot struct {
	Documents []domain.ProjectDocument
	Reviews   []domain.ReviewRecord
	Revisions []domain.DocumentRevision
	Archives  []domain.ArchiveRecord
}

func (s *BoltStore) Snapshot() (Snapshot, error) {
	documents, err := s.ListDocuments()
	if err != nil {
		return Snapshot{}, err
	}
	snapshot := Snapshot{Documents: documents, Reviews: []domain.ReviewRecord{}, Revisions: []domain.DocumentRevision{}, Archives: []domain.ArchiveRecord{}}
	for _, document := range documents {
		reviews, reviewErr := s.ListReviews(document.ID)
		if reviewErr != nil {
			return Snapshot{}, reviewErr
		}
		revisions, revisionErr := s.ListRevisions(document.ID)
		if revisionErr != nil {
			return Snapshot{}, revisionErr
		}
		archive, archiveErr := s.LoadArchive(document.ID)
		if archiveErr == nil {
			snapshot.Archives = append(snapshot.Archives, archive)
		}
		snapshot.Reviews = append(snapshot.Reviews, reviews...)
		snapshot.Revisions = append(snapshot.Revisions, revisions...)
	}
	sort.Slice(snapshot.Reviews, func(i, j int) bool { return snapshot.Reviews[i].ID < snapshot.Reviews[j].ID })
	return snapshot, nil
}

func EncodeSnapshot(snapshot Snapshot) ([]byte, error) { return json.MarshalIndent(snapshot, "", "  ") }

func DecodeSnapshot(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if len(data) == 0 {
		return Snapshot{}, errors.New("snapshot is empty")
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return Snapshot{}, err
	}
	return snapshot, nil
}

func (s *BoltStore) Restore(snapshot Snapshot) error {
	for _, document := range snapshot.Documents {
		if err := s.SaveDocument(document); err != nil {
			return err
		}
	}
	for _, review := range snapshot.Reviews {
		if err := s.SaveReview(review); err != nil {
			return err
		}
	}
	for _, revision := range snapshot.Revisions {
		if err := s.SaveRevision(revision); err != nil {
			return err
		}
	}
	for _, archive := range snapshot.Archives {
		if err := s.SaveArchive(archive); err != nil {
			return err
		}
	}
	return nil
}

func SnapshotDocumentCount(snapshot Snapshot) int { return len(snapshot.Documents) }

func SnapshotHasStatus(snapshot Snapshot, status domain.DocumentStatus) bool {
	for _, document := range snapshot.Documents {
		if document.Status == status {
			return true
		}
	}
	return false
}
