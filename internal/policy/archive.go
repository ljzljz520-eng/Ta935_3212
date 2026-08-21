package policy

import (
	"errors"
	"strings"

	"engineering-document-vault/internal/domain"
)

func ArchiveLocation(document domain.ProjectDocument) string {
	category := "general"
	if len(document.Categories) > 0 && strings.TrimSpace(document.Categories[0]) != "" {
		category = strings.ToLower(strings.TrimSpace(document.Categories[0]))
	}
	return "archive/" + category + "/" + document.ID
}

func ValidateArchiveRequest(document domain.ProjectDocument, checksum string) error {
	if err := CanArchive(document); err != nil {
		return err
	}
	if strings.TrimSpace(checksum) == "" {
		return errors.New("archive checksum is required")
	}
	return nil
}

func ArchiveNote(location string) string {
	return "stored:" + strings.TrimSpace(location)
}
