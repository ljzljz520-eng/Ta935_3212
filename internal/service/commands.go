package service

import (
	"errors"
	"fmt"
	"strings"

	"engineering-document-vault/internal/domain"
	"engineering-document-vault/internal/search"
	"engineering-document-vault/internal/workflow13"
)

type CommandResult struct {
	Status  string
	Message string
	Count   int
}

func (a *Application) SubmitDocument(input workflow13.Submission) (CommandResult, error) {
	if a == nil || a.Workflow == nil {
		return CommandResult{}, errors.New("application is not initialized")
	}
	result, err := a.Workflow.Submit(input)
	if err != nil {
		return CommandResult{Status: "rejected", Message: err.Error()}, err
	}
	return CommandResult{Status: string(result.Document.Status), Message: result.Review.Reason, Count: 1}, nil
}

func (a *Application) VerifyDocument(id, reviewer string, sequence int) (CommandResult, error) {
	if strings.TrimSpace(id) == "" {
		return CommandResult{}, errors.New("document id is required")
	}
	result, err := a.Workflow.VerifyDocument(id, reviewer, sequence)
	if err != nil {
		return CommandResult{Status: string(result.Document.Status), Message: err.Error()}, err
	}
	return CommandResult{Status: string(result.Document.Status), Message: result.Review.Explanation(), Count: 1}, nil
}

func (a *Application) ArchiveDocument(id, checksum string) (CommandResult, error) {
	record, err := a.Workflow.ArchiveDocument(id, checksum)
	if err != nil {
		return CommandResult{Status: "rejected", Message: err.Error()}, err
	}
	return CommandResult{Status: string(domain.StatusArchived), Message: record.Location, Count: 1}, nil
}

func (a *Application) QueryDocuments(query search.Query) ([]domain.ProjectDocument, error) {
	if a == nil || a.Workflow == nil {
		return nil, errors.New("application is not initialized")
	}
	return a.Workflow.FindDocuments(query)
}

func (a *Application) QuerySummary(query search.Query) (CommandResult, error) {
	items, err := a.QueryDocuments(query)
	if err != nil {
		return CommandResult{}, err
	}
	return CommandResult{Status: "ok", Message: fmt.Sprintf("%d documents found", len(items)), Count: len(items)}, nil
}

func (a *Application) QueueSummary(query search.Query) (CommandResult, error) {
	items, err := a.QueryDocuments(query)
	if err != nil {
		return CommandResult{}, err
	}
	queue := workflow13.BuildQueue(items, query)
	return CommandResult{Status: "queue", Message: workflow13.QueueSummary(queue), Count: queue.Len()}, nil
}

func (a *Application) ExplainStatus(id string) (string, error) {
	document, err := a.Workflow.LoadDocument(id)
	if err != nil {
		return "", err
	}
	if document.Status == domain.StatusRejected {
		return "需要修订后重新提交", nil
	}
	if document.Status == domain.StatusApproved {
		return "已通过审核，可归档", nil
	}
	if document.Status == domain.StatusArchived {
		return "已归档，可查询", nil
	}
	return "等待审核", nil
}

func ValidateCommand(command, value string) error {
	if strings.TrimSpace(command) == "" {
		return errors.New("command is required")
	}
	if strings.TrimSpace(value) == "" {
		return errors.New("command value is required")
	}
	return nil
}
