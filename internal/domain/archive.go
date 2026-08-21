package domain

import "errors"

type ArchiveRecord struct {
	ID             string `json:"id"`
	DocumentID     string `json:"document_id"`
	Location       string `json:"location"`
	Checksum       string `json:"checksum"`
	Classification string `json:"classification"`
}

func NewArchive(id string, document ProjectDocument, location, checksum string) (ArchiveRecord, error) {
	if id == "" || document.ID == "" || location == "" || checksum == "" {
		return ArchiveRecord{}, errors.New("archive fields are required")
	}
	if document.Status != StatusApproved {
		return ArchiveRecord{}, errors.New("only approved documents can be archived")
	}
	classification := "uncategorized"
	if len(document.Categories) > 0 {
		classification = document.Categories[0]
	}
	return ArchiveRecord{ID: id, DocumentID: document.ID, Location: location, Checksum: checksum, Classification: classification}, nil
}
