package workflow13

import (
	"fmt"
	"sort"
	"strings"

	"engineering-document-vault/internal/domain"
)

type AuditEvent struct {
	DocumentID string
	Kind       string
	Actor      string
	Detail     string
}

func EventFromReview(review domain.ReviewRecord) AuditEvent {
	return AuditEvent{DocumentID: review.DocumentID, Kind: string(review.Decision), Actor: review.Reviewer, Detail: review.Reason}
}

func EventFromDocument(document domain.ProjectDocument) []AuditEvent {
	events := make([]AuditEvent, 0, len(document.AuditTrail))
	for index, entry := range document.AuditTrail {
		events = append(events, AuditEvent{DocumentID: document.ID, Kind: fmt.Sprintf("event-%d", index+1), Detail: entry})
	}
	return events
}

func SummarizeEvents(events []AuditEvent) string {
	parts := make([]string, 0, len(events))
	for _, event := range events {
		parts = append(parts, event.Kind+":"+event.Detail)
	}
	return strings.Join(parts, " | ")
}

func SortEvents(events []AuditEvent) {
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].DocumentID == events[j].DocumentID {
			return events[i].Kind < events[j].Kind
		}
		return events[i].DocumentID < events[j].DocumentID
	})
}

func FilterEvents(events []AuditEvent, documentID string) []AuditEvent {
	result := []AuditEvent{}
	for _, event := range events {
		if documentID == "" || event.DocumentID == documentID {
			result = append(result, event)
		}
	}
	SortEvents(result)
	return result
}

func IsRejectionEvent(event AuditEvent) bool {
	return strings.EqualFold(event.Kind, string(domain.DecisionReject)) || strings.Contains(strings.ToLower(event.Detail), "rejected")
}

func RejectionCount(events []AuditEvent) int {
	count := 0
	for _, event := range events {
		if IsRejectionEvent(event) {
			count++
		}
	}
	return count
}

func ApprovalCount(events []AuditEvent) int {
	count := 0
	for _, event := range events {
		if strings.EqualFold(event.Kind, string(domain.DecisionApprove)) {
			count++
		}
	}
	return count
}
