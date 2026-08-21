package store

import (
	"errors"
	"strings"

	"engineering-document-vault/internal/domain"
	"go.etcd.io/bbolt"
)

type IntegrityIssue struct {
	Entity string
	Key    string
	Reason string
}

func (s *BoltStore) CheckIntegrity() ([]IntegrityIssue, error) {
	issues := []IntegrityIssue{}
	documents, err := s.ListDocuments()
	if err != nil {
		return nil, err
	}
	for _, document := range documents {
		if err := document.Validate(); err != nil {
			issues = append(issues, IntegrityIssue{Entity: "ProjectDocument", Key: document.ID, Reason: err.Error()})
		}
		reviews, reviewErr := s.ListReviews(document.ID)
		if reviewErr != nil {
			return nil, reviewErr
		}
		for _, review := range reviews {
			if review.DocumentID != document.ID {
				issues = append(issues, IntegrityIssue{Entity: "ReviewRecord", Key: review.ID, Reason: "document reference mismatch"})
			}
			if strings.TrimSpace(review.Reason) == "" {
				issues = append(issues, IntegrityIssue{Entity: "ReviewRecord", Key: review.ID, Reason: "review reason is empty"})
			}
		}
		revisions, revisionErr := s.ListRevisions(document.ID)
		if revisionErr != nil {
			return nil, revisionErr
		}
		if len(revisions) == 0 {
			issues = append(issues, IntegrityIssue{Entity: "DocumentRevision", Key: document.ID, Reason: "initial revision is missing"})
		}
		for _, revision := range revisions {
			if revision.Version < 1 {
				issues = append(issues, IntegrityIssue{Entity: "DocumentRevision", Key: revision.ID, Reason: "revision number is invalid"})
			}
		}
		if document.Status == domain.StatusArchived {
			if _, archiveErr := s.LoadArchive(document.ID); archiveErr != nil {
				issues = append(issues, IntegrityIssue{Entity: "ArchiveRecord", Key: document.ID, Reason: "archive index is missing"})
			}
		}
	}
	return issues, nil
}

func (s *BoltStore) RequireHealthy() error {
	issues, err := s.CheckIntegrity()
	if err != nil {
		return err
	}
	if len(issues) > 0 {
		return errors.New("store integrity check failed")
	}
	return nil
}

func (s *BoltStore) CountByStatus(status domain.DocumentStatus) (int, error) {
	documents, err := s.ListDocuments()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, document := range documents {
		if document.Status == status {
			count++
		}
	}
	return count, nil
}

func (s *BoltStore) CountByOwner(owner string) (int, error) {
	documents, err := s.ListDocuments()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, document := range documents {
		if strings.EqualFold(document.Owner, owner) {
			count++
		}
	}
	return count, nil
}

func (s *BoltStore) HasDocument(id string) (bool, error) {
	_, err := s.LoadDocument(id)
	if err == nil {
		return true, nil
	}
	if strings.Contains(err.Error(), "not found") {
		return false, nil
	}
	return false, err
}

func (s *BoltStore) DeleteDocument(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("document id is required")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(documentsBucket).Delete([]byte(id)) })
}
