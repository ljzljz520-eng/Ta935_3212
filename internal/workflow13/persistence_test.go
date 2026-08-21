package workflow13

import (
	"engineering-document-vault/internal/search"
	"engineering-document-vault/internal/store"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := store.TempPath("reopen")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(db)
	_, err = service.Submit(Submission{ID: "d1", Title: "Archive plan", Owner: "alice", Categories: []string{"contract"}, Artifacts: []string{"contract.pdf"}, Reviewer: "reviewer", Sequence: 1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.ArchiveDocument("d1", "checksum"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	document, err := NewService(db).LoadDocument("d1")
	if err != nil || document.Status != "archived" {
		t.Fatalf("unexpected restored document: %+v %v", document, err)
	}
}

func TestArchivedDocumentRemainsDiscoverable(t *testing.T) {
	path := store.TempPath("archive-query")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db)
	_, err = service.Submit(Submission{ID: "d1", Title: "Archive plan", Owner: "alice", Categories: []string{"contract"}, Artifacts: []string{"contract.pdf"}, Reviewer: "reviewer", Sequence: 1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.ArchiveDocument("d1", "checksum"); err != nil {
		t.Fatal(err)
	}
	items, err := service.FindDocuments(search.Query{Keyword: "Archive", Status: "archived"})
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected query result: %+v %v", items, err)
	}
}
