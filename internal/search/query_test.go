package search

import (
	"engineering-document-vault/internal/domain"
	"testing"
)

func TestFilterMatchesOwnerCategoryAndKeyword(t *testing.T) {
	first, _ := domain.NewDocument("a", "Bridge Plan", "alice", []string{"design"}, []string{"bridge.pdf"})
	second, _ := domain.NewDocument("b", "Tunnel Plan", "bob", []string{"design"}, []string{"tunnel.pdf"})
	items := Filter([]domain.ProjectDocument{second, first}, Parse("bridge", "alice", "design", ""))
	if len(items) != 1 || items[0].ID != "a" {
		t.Fatalf("unexpected matches: %+v", items)
	}
}
