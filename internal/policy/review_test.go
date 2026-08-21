package policy

import (
	"engineering-document-vault/internal/domain"
	"testing"
)

func TestReviewPolicyApprovesCompleteDocument(t *testing.T) {
	document, _ := domain.NewDocument("d1", "title", "owner", []string{"design"}, []string{"plan.pdf"})
	decision, reason := DefaultReviewPolicy().Evaluate(document, 1)
	if decision != domain.DecisionApprove || reason == "" {
		t.Fatalf("unexpected decision: %s %s", decision, reason)
	}
}

func TestReviewPolicyRejectsSeventhReview(t *testing.T) {
	document, _ := domain.NewDocument("d1", "title", "owner", []string{"design"}, []string{"plan.pdf"})
	decision, reason := DefaultReviewPolicy().Evaluate(document, 7)
	if decision != domain.DecisionReject || reason == "" {
		t.Fatalf("unexpected decision: %s %s", decision, reason)
	}
}
