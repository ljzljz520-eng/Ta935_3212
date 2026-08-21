package policy

import (
	"fmt"
	"strings"

	"engineering-document-vault/internal/domain"
)

type Gate struct {
	Name        string
	Description string
	Evaluate    func(domain.ProjectDocument) bool
}

func StandardGates() []Gate {
	return []Gate{
		{Name: "identity", Description: "document identity", Evaluate: func(d domain.ProjectDocument) bool { return d.ID != "" }},
		{Name: "ownership", Description: "responsible owner", Evaluate: func(d domain.ProjectDocument) bool { return d.Owner != "" }},
		{Name: "classification", Description: "registered category", Evaluate: func(d domain.ProjectDocument) bool { return len(d.Categories) > 0 }},
		{Name: "evidence", Description: "attached artifact", Evaluate: func(d domain.ProjectDocument) bool { return len(d.Artifacts) > 0 }},
		{Name: "revision", Description: "positive revision", Evaluate: func(d domain.ProjectDocument) bool { return d.Revision > 0 }},
	}
}

func EvaluateGates(document domain.ProjectDocument) []string {
	failed := []string{}
	for _, gate := range StandardGates() {
		if !gate.Evaluate(document) {
			failed = append(failed, gate.Name)
		}
	}
	return failed
}

func GateSummary(document domain.ProjectDocument) string {
	failed := EvaluateGates(document)
	if len(failed) == 0 {
		return "all gates passed"
	}
	return "failed gates: " + strings.Join(failed, ",")
}

func ExplainDecision(decision domain.ReviewDecision, reason string) string {
	if decision == domain.DecisionApprove {
		return "approved: " + reason
	}
	if decision == domain.DecisionReject {
		return "rejected: " + reason
	}
	return fmt.Sprintf("unknown decision %q", decision)
}

func IsTerminal(status domain.DocumentStatus) bool {
	return status == domain.StatusArchived
}

func NextReviewSequence(reviews []domain.ReviewRecord) int {
	max := 0
	for _, review := range reviews {
		if review.Sequence > max {
			max = review.Sequence
		}
	}
	return max + 1
}
