package service

import (
	"errors"
	"strings"

	"engineering-document-vault/internal/domain"
	"engineering-document-vault/internal/policy"
)

type LifecycleSummary struct {
	DocumentID string
	Status     domain.DocumentStatus
	Revision   int
	LastEvent  string
}

func (a *Application) Lifecycle(id string) (LifecycleSummary, error) {
	document, err := a.Workflow.LoadDocument(id)
	if err != nil {
		return LifecycleSummary{}, err
	}
	return LifecycleSummary{DocumentID: document.ID, Status: document.Status, Revision: document.Revision, LastEvent: document.LatestAudit()}, nil
}

func (a *Application) ValidateForArchive(id, checksum string) error {
	document, err := a.Workflow.LoadDocument(id)
	if err != nil {
		return err
	}
	return policy.ValidateArchiveRequest(document, checksum)
}

func (a *Application) ReviewHistory(id string) ([]domain.ReviewRecord, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("document id is required")
	}
	return a.Workflow.Reviews(id)
}

func (a *Application) RevisionHistory(id string) ([]domain.DocumentRevision, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("document id is required")
	}
	return a.Workflow.Revisions(id)
}

func (a *Application) ArchiveIndex(id string) (domain.ArchiveRecord, error) {
	if strings.TrimSpace(id) == "" {
		return domain.ArchiveRecord{}, errors.New("document id is required")
	}
	return a.Workflow.Archive(id)
}

func (a *Application) EnsureHealthy() error {
	if a == nil || a.Workflow == nil || a.Workflow.Store == nil {
		return errors.New("application is not initialized")
	}
	return a.Workflow.Store.RequireHealthy()
}

func (a *Application) HasDocument(id string) (bool, error) {
	if a == nil || a.Workflow == nil {
		return false, errors.New("application is not initialized")
	}
	return a.Workflow.Store.HasDocument(id)
}

func (a *Application) Count(status domain.DocumentStatus) (int, error) {
	return a.Workflow.Store.CountByStatus(status)
}

func (a *Application) OwnerCount(owner string) (int, error) {
	return a.Workflow.Store.CountByOwner(owner)
}

func (a *Application) CloseWithHealthCheck() error {
	if err := a.EnsureHealthy(); err != nil {
		return err
	}
	return a.Close()
}

func StatusRequiresArchive(status domain.DocumentStatus) bool {
	return status == domain.StatusApproved
}

func StatusAllowsQuery(status domain.DocumentStatus) bool {
	switch status {
	case domain.StatusDraft, domain.StatusApproved, domain.StatusArchived:
		return true
	default:
		return false
	}
}
