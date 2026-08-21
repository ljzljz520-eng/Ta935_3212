package workflow13

import (
	"errors"
	"fmt"
	"strings"

	"engineering-document-vault/internal/catalog"
	"engineering-document-vault/internal/domain"
	"engineering-document-vault/internal/policy"
	"engineering-document-vault/internal/search"
	"engineering-document-vault/internal/store"
)

type Service struct {
	Store  *store.BoltStore
	Policy policy.ReviewPolicy
}

type Submission struct {
	ID         string
	Title      string
	Owner      string
	Categories []string
	Artifacts  []string
	Reviewer   string
	Sequence   int
}

type ReviewResult struct {
	Document domain.ProjectDocument
	Review   domain.ReviewRecord
}

func NewService(database *store.BoltStore) *Service {
	return &Service{Store: database, Policy: policy.DefaultReviewPolicy()}
}

func (s *Service) Submit(input Submission) (ReviewResult, error) {
	categories := catalog.Normalize(input.Categories)
	for _, category := range categories {
		if !catalog.IsKnownCategory(category) {
			return ReviewResult{}, fmt.Errorf("unknown category: %s", category)
		}
	}
	document, err := domain.NewDocument(input.ID, input.Title, input.Owner, categories, input.Artifacts)
	if err != nil {
		return ReviewResult{}, err
	}
	if input.Sequence < 1 {
		input.Sequence = 1
	}
	if err := s.Store.SaveDocument(document); err != nil {
		return ReviewResult{}, err
	}
	revision := domain.DocumentRevision{ID: document.ID + ":1", DocumentID: document.ID, Version: 1, Summary: document.Title, CreatedBy: document.Owner}
	if err := s.Store.SaveRevision(revision); err != nil {
		return ReviewResult{}, err
	}
	return s.VerifyDocument(document.ID, input.Reviewer, input.Sequence)
}

func (s *Service) VerifyDocument(documentID, reviewer string, sequence int) (result ReviewResult, err error) {
	if strings.TrimSpace(reviewer) == "" {
		return ReviewResult{}, errors.New("reviewer is required")
	}
	document, err := s.Store.LoadDocument(documentID)
	if err != nil {
		return ReviewResult{}, err
	}
	if err = policy.ValidateReview(document, sequence); err != nil {
		return ReviewResult{}, err
	}
	decision, reason := s.Policy.Evaluate(document, sequence)
	review, reviewErr := domain.NewReview(fmt.Sprintf("%s:%d", documentID, sequence), documentID, reviewer, sequence, decision, reason)
	if reviewErr != nil {
		return ReviewResult{}, reviewErr
	}
	if decision == domain.DecisionApprove {
		err = document.SetStatus(domain.StatusApproved, reason)
	} else {
		err = document.SetStatus(domain.StatusRejected, reason)
	}
	if err != nil {
		return ReviewResult{}, err
	}
	if err = s.Store.SaveDocument(document); err != nil {
		return ReviewResult{}, err
	}
	if err = s.Store.SaveReview(review); err != nil {
		return ReviewResult{}, err
	}
	if decision == domain.DecisionReject {
		// Surface a stable business rejection to the caller so the interface
		// reports the rejection status and its reason instead of masking it
		// behind a success result. The persisted document and review already
		// record the rejection, so only the returned error is conveyed.
		err = errors.New(reason)
	}
	result = ReviewResult{Document: document, Review: review}
	return result, err
}

func (s *Service) FindDocuments(query search.Query) ([]domain.ProjectDocument, error) {
	documents, err := s.Store.ListDocuments()
	if err != nil {
		return nil, err
	}
	return search.Filter(documents, query), nil
}

func (s *Service) ArchiveDocument(documentID, checksum string) (domain.ArchiveRecord, error) {
	document, err := s.Store.LoadDocument(documentID)
	if err != nil {
		return domain.ArchiveRecord{}, err
	}
	if err := policy.ValidateArchiveRequest(document, checksum); err != nil {
		return domain.ArchiveRecord{}, err
	}
	record, err := domain.NewArchive("archive:"+documentID, document, policy.ArchiveLocation(document), checksum)
	if err != nil {
		return domain.ArchiveRecord{}, err
	}
	if err := document.SetStatus(domain.StatusArchived, policy.ArchiveNote(record.Location)); err != nil {
		return domain.ArchiveRecord{}, err
	}
	if err := s.Store.SaveDocument(document); err != nil {
		return domain.ArchiveRecord{}, err
	}
	if err := s.Store.SaveArchive(record); err != nil {
		return domain.ArchiveRecord{}, err
	}
	return record, nil
}

func (s *Service) LoadDocument(id string) (domain.ProjectDocument, error) {
	return s.Store.LoadDocument(id)
}

func (s *Service) Reviews(id string) ([]domain.ReviewRecord, error) { return s.Store.ListReviews(id) }

func (s *Service) Revisions(id string) ([]domain.DocumentRevision, error) {
	return s.Store.ListRevisions(id)
}

func (s *Service) Archive(id string) (domain.ArchiveRecord, error) { return s.Store.LoadArchive(id) }

func IsBusinessRejection(err error) bool {
	return err != nil && strings.Contains(err.Error(), "rejected")
}
