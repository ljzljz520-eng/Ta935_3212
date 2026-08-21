package policy

import (
	"engineering-document-vault/internal/domain"
	"testing"
)

func TestArchiveEligibilityRequiresApproval(t *testing.T) {
	document, _ := domain.NewDocument("d1", "title", "owner", []string{"design"}, []string{"plan.pdf"})
	if err := CanArchive(document); err == nil {
		t.Fatal("expected archive rejection")
	}
	_ = document.SetStatus(domain.StatusApproved, "ok")
	if err := ValidateArchiveRequest(document, "abc"); err != nil {
		t.Fatal(err)
	}
}
