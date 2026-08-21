package domain

import "testing"

func TestDocumentValidationRequiresFields(t *testing.T) {
	if _, err := NewDocument("", "title", "owner", []string{"design"}, []string{"plan.pdf"}); err == nil {
		t.Fatal("expected id validation")
	}
	document, err := NewDocument("d1", "title", "owner", []string{"design"}, []string{"plan.pdf"})
	if err != nil || document.Status != StatusDraft {
		t.Fatalf("unexpected document: %+v %v", document, err)
	}
}

func TestDocumentStatusTransitions(t *testing.T) {
	document, _ := NewDocument("d1", "title", "owner", []string{"design"}, []string{"plan.pdf"})
	if err := document.SetStatus(StatusApproved, "ok"); err != nil {
		t.Fatal(err)
	}
	if err := document.SetStatus(StatusArchived, "stored"); err != nil {
		t.Fatal(err)
	}
	if document.Status != StatusArchived || len(document.AuditTrail) != 3 {
		t.Fatalf("unexpected state: %+v", document)
	}
}
