package domain

import (
	"errors"
	"strings"
)

type DocumentStatus string

const (
	StatusDraft    DocumentStatus = "draft"
	StatusApproved DocumentStatus = "approved"
	StatusRejected DocumentStatus = "rejected"
	StatusArchived DocumentStatus = "archived"
)

type ProjectDocument struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	Revision   int            `json:"revision"`
	Status     DocumentStatus `json:"status"`
	Owner      string         `json:"owner"`
	Categories []string       `json:"categories"`
	Artifacts  []string       `json:"artifacts"`
	AuditTrail []string       `json:"audit_trail"`
}

type DocumentRevision struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	Version    int    `json:"version"`
	Summary    string `json:"summary"`
	CreatedBy  string `json:"created_by"`
}

func NewDocument(id, title, owner string, categories, artifacts []string) (ProjectDocument, error) {
	d := ProjectDocument{ID: strings.TrimSpace(id), Title: strings.TrimSpace(title), Revision: 1, Status: StatusDraft, Owner: strings.TrimSpace(owner), Categories: append([]string(nil), categories...), Artifacts: append([]string(nil), artifacts...)}
	if err := d.Validate(); err != nil {
		return ProjectDocument{}, err
	}
	d.AuditTrail = []string{"created"}
	return d, nil
}

func (d ProjectDocument) Validate() error {
	if d.ID == "" {
		return errors.New("document id is required")
	}
	if d.Title == "" {
		return errors.New("document title is required")
	}
	if d.Owner == "" {
		return errors.New("document owner is required")
	}
	if len(d.Categories) == 0 {
		return errors.New("at least one category is required")
	}
	if d.Revision < 1 {
		return errors.New("revision must be positive")
	}
	return nil
}

func (d *ProjectDocument) AddRevision(summary, author string) error {
	if strings.TrimSpace(summary) == "" || strings.TrimSpace(author) == "" {
		return errors.New("revision summary and author are required")
	}
	d.Revision++
	d.AuditTrail = append(d.AuditTrail, "revision:"+strings.TrimSpace(author))
	return nil
}

func (d *ProjectDocument) SetStatus(status DocumentStatus, note string) error {
	if status != StatusDraft && status != StatusApproved && status != StatusRejected && status != StatusArchived {
		return errors.New("unknown document status")
	}
	if strings.TrimSpace(note) == "" {
		return errors.New("status note is required")
	}
	d.Status = status
	d.AuditTrail = append(d.AuditTrail, string(status)+":"+strings.TrimSpace(note))
	return nil
}

func (d ProjectDocument) IsVisible() bool {
	return d.Status == StatusApproved || d.Status == StatusArchived || d.Status == StatusDraft
}

func (d ProjectDocument) HasCategory(category string) bool {
	for _, item := range d.Categories {
		if strings.EqualFold(item, category) {
			return true
		}
	}
	return false
}
