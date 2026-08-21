package service

import (
	"engineering-document-vault/internal/search"
	"engineering-document-vault/internal/store"
	"engineering-document-vault/internal/workflow13"
)

type Application struct{ Workflow *workflow13.Service }

func OpenApplication(path string) (*Application, error) {
	db, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	return &Application{Workflow: workflow13.NewService(db)}, nil
}

func (a *Application) Close() error { return a.Workflow.Store.Close() }

func (a *Application) Search(keyword, owner, category, status string) (int, error) {
	items, err := a.Workflow.FindDocuments(search.Parse(keyword, owner, category, status))
	return len(items), err
}
