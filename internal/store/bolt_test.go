package store

import (
	"engineering-document-vault/internal/domain"
	"testing"
)

func TestStorePersistsDocuments(t *testing.T) {
	path := TempPath("store-test")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	document, _ := domain.NewDocument("d1", "title", "owner", []string{"design"}, []string{"plan.pdf"})
	if err := db.SaveDocument(document); err != nil {
		t.Fatal(err)
	}
	loaded, err := db.LoadDocument("d1")
	_ = db.Close()
	if err != nil || loaded.ID != "d1" {
		t.Fatalf("unexpected load: %+v %v", loaded, err)
	}
}
