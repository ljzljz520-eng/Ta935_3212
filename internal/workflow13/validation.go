package workflow13

import (
	"errors"
	"fmt"
	"strings"

	"engineering-document-vault/internal/catalog"
	"engineering-document-vault/internal/domain"
)

type SubmissionIssue struct {
	Field  string
	Reason string
}

func ValidateSubmission(input Submission) []SubmissionIssue {
	issues := []SubmissionIssue{}
	if strings.TrimSpace(input.ID) == "" {
		issues = append(issues, SubmissionIssue{Field: "id", Reason: "required"})
	}
	if strings.TrimSpace(input.Title) == "" {
		issues = append(issues, SubmissionIssue{Field: "title", Reason: "required"})
	}
	if strings.TrimSpace(input.Owner) == "" {
		issues = append(issues, SubmissionIssue{Field: "owner", Reason: "required"})
	}
	if len(input.Categories) == 0 {
		issues = append(issues, SubmissionIssue{Field: "categories", Reason: "at least one category required"})
	}
	if len(input.Artifacts) == 0 {
		issues = append(issues, SubmissionIssue{Field: "artifacts", Reason: "at least one artifact required"})
	}
	for _, category := range input.Categories {
		if !catalog.IsKnownCategory(category) {
			issues = append(issues, SubmissionIssue{Field: "categories", Reason: "unknown category: " + category})
		}
	}
	if input.Sequence < 0 {
		issues = append(issues, SubmissionIssue{Field: "sequence", Reason: "cannot be negative"})
	}
	return issues
}

func SubmissionError(input Submission) error {
	issues := ValidateSubmission(input)
	if len(issues) == 0 {
		return nil
	}
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, issue.Field+" "+issue.Reason)
	}
	return errors.New(strings.Join(parts, "; "))
}

func ReviewInvariant(result ReviewResult) error {
	if result.Document.ID == "" || result.Review.DocumentID == "" {
		return errors.New("review result has no document")
	}
	if result.Document.ID != result.Review.DocumentID {
		return errors.New("review target mismatch")
	}
	if result.Review.Sequence < 1 {
		return errors.New("review sequence is invalid")
	}
	if result.Review.Decision == domain.DecisionApprove && result.Document.Status != domain.StatusApproved {
		return errors.New("approved review has wrong document status")
	}
	if result.Review.Decision == domain.DecisionReject && result.Document.Status != domain.StatusRejected {
		return errors.New("rejected review has wrong document status")
	}
	return nil
}

func ArchiveInvariant(document domain.ProjectDocument, archive domain.ArchiveRecord) error {
	if archive.DocumentID != document.ID {
		return errors.New("archive target mismatch")
	}
	if document.Status != domain.StatusArchived {
		return errors.New("document is not archived")
	}
	if strings.TrimSpace(archive.Location) == "" {
		return errors.New("archive location is empty")
	}
	return nil
}

func FormatIssues(issues []SubmissionIssue) string {
	if len(issues) == 0 {
		return "valid"
	}
	parts := make([]string, 0, len(issues))
	for index, issue := range issues {
		parts = append(parts, fmt.Sprintf("%d.%s:%s", index+1, issue.Field, issue.Reason))
	}
	return strings.Join(parts, ", ")
}

func IsRetryableReview(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "temporary") || strings.Contains(message, "sequence")
}
