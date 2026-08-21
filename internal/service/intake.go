package service

import (
	"errors"
	"sort"
	"strings"

	"engineering-document-vault/internal/catalog"
	"engineering-document-vault/internal/domain"
	"engineering-document-vault/internal/workflow13"
)

type IntakeSummary struct {
	Accepted int
	Rejected int
	Issues   []string
}

func (a *Application) ValidateIntake(input workflow13.Submission) IntakeSummary {
	issues := workflow13.ValidateSubmission(input)
	summary := IntakeSummary{Issues: []string{}}
	if len(issues) == 0 {
		summary.Accepted = 1
		return summary
	}
	summary.Rejected = 1
	for _, issue := range issues {
		summary.Issues = append(summary.Issues, issue.Field+":"+issue.Reason)
	}
	sort.Strings(summary.Issues)
	return summary
}

func (a *Application) NormalizeSubmission(input workflow13.Submission) workflow13.Submission {
	input.ID = strings.TrimSpace(input.ID)
	input.Title = strings.TrimSpace(input.Title)
	input.Owner = strings.TrimSpace(input.Owner)
	input.Categories = catalog.Normalize(input.Categories)
	if input.Sequence < 1 {
		input.Sequence = 1
	}
	artifacts := []string{}
	seen := map[string]bool{}
	for _, artifact := range input.Artifacts {
		value := strings.TrimSpace(artifact)
		if value != "" && !seen[value] {
			artifacts = append(artifacts, value)
			seen[value] = true
		}
	}
	input.Artifacts = artifacts
	return input
}

func (a *Application) Admit(input workflow13.Submission) (CommandResult, error) {
	normalized := a.NormalizeSubmission(input)
	if err := workflow13.SubmissionError(normalized); err != nil {
		return CommandResult{Status: "rejected", Message: err.Error()}, err
	}
	return a.SubmitDocument(normalized)
}

func (a *Application) AdmitOrExplain(input workflow13.Submission) CommandResult {
	result, err := a.Admit(input)
	if err == nil {
		return result
	}
	return CommandResult{Status: "rejected", Message: workflow13.FormatIssues(workflow13.ValidateSubmission(input))}
}

func (a *Application) ApplyRevision(id, summary, author string) error {
	if a == nil || a.Workflow == nil {
		return errors.New("application is not initialized")
	}
	document, err := a.Workflow.LoadDocument(id)
	if err != nil {
		return err
	}
	if document.Status == domain.StatusArchived {
		return errors.New("archived document cannot be revised")
	}
	if err := document.AddRevision(summary, author); err != nil {
		return err
	}
	revision := domain.DocumentRevision{ID: id + ":" + itoa(document.Revision), DocumentID: id, Version: document.Revision, Summary: strings.TrimSpace(summary), CreatedBy: strings.TrimSpace(author)}
	if err := a.Workflow.Store.SaveDocument(document); err != nil {
		return err
	}
	return a.Workflow.Store.SaveRevision(revision)
}

func (a *Application) RevisionCount(id string) (int, error) {
	revisions, err := a.RevisionHistory(id)
	return len(revisions), err
}

func (a *Application) CategoryLabel(category string) string { return catalog.Label(category) }

func (a *Application) CategoryRetention(category string) int { return catalog.RetentionYears(category) }

func (a *Application) SupportsCategory(category string) bool {
	return catalog.IsKnownCategory(category)
}

func (a *Application) DocumentIsTerminal(id string) (bool, error) {
	document, err := a.Workflow.LoadDocument(id)
	if err != nil {
		return false, err
	}
	return document.Status == domain.StatusArchived, nil
}

func (a *Application) ArchiveReady(id string) (bool, error) {
	document, err := a.Workflow.LoadDocument(id)
	if err != nil {
		return false, err
	}
	return document.Status == domain.StatusApproved && document.ArtifactCount() > 0, nil
}

func (a *Application) LatestReview(id string) (domain.ReviewRecord, error) {
	reviews, err := a.ReviewHistory(id)
	if err != nil {
		return domain.ReviewRecord{}, err
	}
	if len(reviews) == 0 {
		return domain.ReviewRecord{}, errors.New("review record not found")
	}
	return reviews[len(reviews)-1], nil
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	digits := []byte{}
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
