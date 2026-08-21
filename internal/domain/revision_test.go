package domain

import "testing"

func TestRevisionAdvancesWithAudit(t *testing.T) {
	document, _ := NewDocument("d1", "title", "owner", []string{"design"}, []string{"plan.pdf"})
	if err := document.AddRevision("second issue", "engineer"); err != nil {
		t.Fatal(err)
	}
	if document.Revision != 2 || len(document.AuditTrail) != 2 {
		t.Fatalf("unexpected revision: %+v", document)
	}
}

func TestReviewRecordDecisions(t *testing.T) {
	review, err := NewReview("r1", "d1", "alice", 1, DecisionApprove, "ok")
	if err != nil || !review.Approved() || review.Rejected() {
		t.Fatalf("unexpected review: %+v %v", review, err)
	}
}
