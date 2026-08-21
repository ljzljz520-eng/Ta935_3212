package domain

import (
	"fmt"
	"strings"
)

type FieldSpec struct {
	Name     string
	Required bool
	Kind     string
	Purpose  string
}

var documentFieldSpecs = []FieldSpec{
	{Name: "ID", Required: true, Kind: "identifier", Purpose: "stable document identity"},
	{Name: "Title", Required: true, Kind: "text", Purpose: "human-readable title"},
	{Name: "Revision", Required: true, Kind: "number", Purpose: "current revision number"},
	{Name: "Status", Required: true, Kind: "enum", Purpose: "lifecycle state"},
	{Name: "Owner", Required: true, Kind: "principal", Purpose: "responsible engineer"},
	{Name: "Categories", Required: true, Kind: "set", Purpose: "engineering classification"},
	{Name: "Artifacts", Required: false, Kind: "set", Purpose: "linked artifacts"},
	{Name: "AuditTrail", Required: true, Kind: "log", Purpose: "business event history"},
}

var revisionFieldSpecs = []FieldSpec{
	{Name: "ID", Required: true, Kind: "identifier", Purpose: "revision identity"},
	{Name: "DocumentID", Required: true, Kind: "identifier", Purpose: "parent document"},
	{Name: "Version", Required: true, Kind: "number", Purpose: "revision sequence"},
	{Name: "Summary", Required: true, Kind: "text", Purpose: "change summary"},
	{Name: "CreatedBy", Required: true, Kind: "principal", Purpose: "revision author"},
}

var reviewFieldSpecs = []FieldSpec{
	{Name: "ID", Required: true, Kind: "identifier", Purpose: "review identity"},
	{Name: "DocumentID", Required: true, Kind: "identifier", Purpose: "review target"},
	{Name: "Sequence", Required: true, Kind: "number", Purpose: "review attempt"},
	{Name: "Decision", Required: true, Kind: "enum", Purpose: "approval outcome"},
	{Name: "Reason", Required: true, Kind: "text", Purpose: "explainable decision"},
	{Name: "Reviewer", Required: true, Kind: "principal", Purpose: "reviewing person"},
}

var archiveFieldSpecs = []FieldSpec{
	{Name: "ID", Required: true, Kind: "identifier", Purpose: "archive identity"},
	{Name: "DocumentID", Required: true, Kind: "identifier", Purpose: "archived document"},
	{Name: "Location", Required: true, Kind: "path", Purpose: "archive location"},
	{Name: "Checksum", Required: true, Kind: "digest", Purpose: "content integrity"},
	{Name: "Classification", Required: true, Kind: "text", Purpose: "archive shelf"},
}

func FieldSpecs(entity string) []FieldSpec {
	switch strings.ToLower(strings.TrimSpace(entity)) {
	case "projectdocument":
		return append([]FieldSpec(nil), documentFieldSpecs...)
	case "documentrevision":
		return append([]FieldSpec(nil), revisionFieldSpecs...)
	case "reviewrecord":
		return append([]FieldSpec(nil), reviewFieldSpecs...)
	case "archiverecord":
		return append([]FieldSpec(nil), archiveFieldSpecs...)
	default:
		return nil
	}
}

func RequiredFields(entity string) []string {
	fields := FieldSpecs(entity)
	result := make([]string, 0, len(fields))
	for _, field := range fields {
		if field.Required {
			result = append(result, field.Name)
		}
	}
	return result
}

func FieldDescription(entity, name string) string {
	for _, field := range FieldSpecs(entity) {
		if strings.EqualFold(field.Name, name) {
			return fmt.Sprintf("%s (%s): %s", field.Name, field.Kind, field.Purpose)
		}
	}
	return ""
}

func ValidateFieldSet(entity string, values map[string]string) []string {
	missing := []string{}
	for _, field := range FieldSpecs(entity) {
		if field.Required && strings.TrimSpace(values[field.Name]) == "" {
			missing = append(missing, field.Name)
		}
	}
	return missing
}

func EntityNames() []string {
	return []string{"ProjectDocument", "ReviewRecord", "DocumentRevision", "ArchiveRecord"}
}

func IsLifecycleStatus(status DocumentStatus) bool {
	switch status {
	case StatusDraft, StatusApproved, StatusRejected, StatusArchived:
		return true
	default:
		return false
	}
}

func CanTransition(from, to DocumentStatus) bool {
	if !IsLifecycleStatus(from) || !IsLifecycleStatus(to) {
		return false
	}
	if from == to {
		return true
	}
	switch from {
	case StatusDraft:
		return to == StatusApproved || to == StatusRejected
	case StatusApproved:
		return to == StatusArchived
	case StatusRejected:
		return to == StatusDraft
	case StatusArchived:
		return false
	default:
		return false
	}
}

func TransitionReason(from, to DocumentStatus) string {
	if !CanTransition(from, to) {
		return "transition is not allowed"
	}
	if from == to {
		return "state is unchanged"
	}
	return fmt.Sprintf("%s to %s is permitted", from, to)
}
