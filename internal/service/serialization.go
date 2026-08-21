package service

import (
	"encoding/json"
	"errors"
	"strings"

	"engineering-document-vault/internal/domain"
	"engineering-document-vault/internal/store"
)

type ExportEnvelope struct {
	Format  string         `json:"format"`
	Version int            `json:"version"`
	Data    store.Snapshot `json:"data"`
}

func (a *Application) Export() ([]byte, error) {
	snapshot, err := a.Workflow.Store.Snapshot()
	if err != nil {
		return nil, err
	}
	envelope := ExportEnvelope{Format: "engineering-document-vault", Version: 1, Data: snapshot}
	return json.MarshalIndent(envelope, "", "  ")
}

func (a *Application) Import(data []byte) error {
	if len(data) == 0 {
		return errors.New("import payload is empty")
	}
	var envelope ExportEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	if envelope.Format != "engineering-document-vault" || envelope.Version != 1 {
		return errors.New("unsupported import format")
	}
	return a.Workflow.Store.Restore(envelope.Data)
}

func EncodeDocument(document domain.ProjectDocument) ([]byte, error) { return json.Marshal(document) }

func DecodeDocument(data []byte) (domain.ProjectDocument, error) {
	var document domain.ProjectDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return domain.ProjectDocument{}, err
	}
	if err := document.Validate(); err != nil {
		return domain.ProjectDocument{}, err
	}
	return document, nil
}

func NormalizeExportName(name string) string {
	value := strings.TrimSpace(name)
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.ReplaceAll(value, "/", "-")
	if value == "" {
		return "document-vault"
	}
	return strings.ToLower(value)
}
