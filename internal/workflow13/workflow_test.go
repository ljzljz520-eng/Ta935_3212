package workflow13

import (
	"engineering-document-vault/internal/store"
	"testing"
)

func TestDocumentAdmissionAndApproval(t *testing.T) {
	path := store.TempPath("workflow")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	result, err := NewService(db).Submit(Submission{ID: "d1", Title: "Bridge plan", Owner: "alice", Categories: []string{"design"}, Artifacts: []string{"plan.pdf"}, Reviewer: "reviewer", Sequence: 1})
	if err != nil || !result.Review.Approved() || result.Document.Status != "approved" {
		t.Fatalf("unexpected result: %+v %v", result, err)
	}
}

func TestWorkflow13BusinessInvariant(t *testing.T) {
	path := store.TempPath("workflow-invariant")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = NewService(db).Submit(Submission{ID: "d7", Title: "Safety gate", Owner: "alice", Categories: []string{"safety"}, Artifacts: []string{"safety.pdf"}, Reviewer: "reviewer", Sequence: 7})
	if err == nil {
		t.Fatal("expected stable rejection on seventh review")
	}
}
