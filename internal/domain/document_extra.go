package domain

import "sort"

func CloneDocument(input ProjectDocument) ProjectDocument {
	input.Categories = append([]string(nil), input.Categories...)
	input.Artifacts = append([]string(nil), input.Artifacts...)
	input.AuditTrail = append([]string(nil), input.AuditTrail...)
	return input
}

func SortDocuments(documents []ProjectDocument) {
	sort.Slice(documents, func(i, j int) bool {
		if documents[i].Title == documents[j].Title {
			return documents[i].ID < documents[j].ID
		}
		return documents[i].Title < documents[j].Title
	})
}

func StatusLabel(status DocumentStatus) string {
	switch status {
	case StatusDraft:
		return "Draft"
	case StatusApproved:
		return "Approved"
	case StatusRejected:
		return "Rejected"
	case StatusArchived:
		return "Archived"
	default:
		return "Unknown"
	}
}
