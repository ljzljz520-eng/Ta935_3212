package domain

import (
	"errors"
	"strings"
)

type LifecycleEvent struct {
	DocumentID string
	From       DocumentStatus
	To         DocumentStatus
	Actor      string
	Reason     string
}

func (d *ProjectDocument) ApplyEvent(event LifecycleEvent) error {
	if event.DocumentID != d.ID {
		return errors.New("lifecycle event targets another document")
	}
	if strings.TrimSpace(event.Actor) == "" {
		return errors.New("lifecycle actor is required")
	}
	if strings.TrimSpace(event.Reason) == "" {
		return errors.New("lifecycle reason is required")
	}
	if event.From != d.Status {
		return errors.New("lifecycle event has stale source state")
	}
	if !CanTransition(event.From, event.To) {
		return errors.New(TransitionReason(event.From, event.To))
	}
	d.Status = event.To
	d.AuditTrail = append(d.AuditTrail, string(event.To)+":"+event.Actor+":"+event.Reason)
	return nil
}

func (d ProjectDocument) LatestAudit() string {
	if len(d.AuditTrail) == 0 {
		return ""
	}
	return d.AuditTrail[len(d.AuditTrail)-1]
}

func (d ProjectDocument) AuditContains(fragment string) bool {
	needle := strings.ToLower(strings.TrimSpace(fragment))
	if needle == "" {
		return false
	}
	for _, entry := range d.AuditTrail {
		if strings.Contains(strings.ToLower(entry), needle) {
			return true
		}
	}
	return false
}

func (d ProjectDocument) ArtifactCount() int {
	count := 0
	for _, artifact := range d.Artifacts {
		if strings.TrimSpace(artifact) != "" {
			count++
		}
	}
	return count
}

func (d ProjectDocument) CategoryCount() int {
	seen := map[string]bool{}
	count := 0
	for _, category := range d.Categories {
		key := strings.ToLower(strings.TrimSpace(category))
		if key != "" && !seen[key] {
			seen[key] = true
			count++
		}
	}
	return count
}

func (d ProjectDocument) Summary() string {
	return d.ID + "|" + d.Title + "|" + string(d.Status) + "|" + d.Owner
}

func (r DocumentRevision) IsFor(documentID string) bool {
	return strings.TrimSpace(documentID) != "" && r.DocumentID == documentID
}

func (r DocumentRevision) SummaryLine() string {
	return r.DocumentID + "#" + string(rune('0'+r.Version)) + ":" + r.Summary
}

func (r ReviewRecord) Explanation() string {
	return string(r.Decision) + " by " + r.Reviewer + ": " + r.Reason
}

func (a ArchiveRecord) ShelfKey() string {
	return a.Classification + "/" + a.DocumentID
}
