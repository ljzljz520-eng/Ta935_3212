package service

import (
	"fmt"
	"sort"
	"strings"

	"engineering-document-vault/internal/domain"
	"engineering-document-vault/internal/search"
)

type Report struct {
	Total         int
	ByStatus      map[domain.DocumentStatus]int
	ByOwner       map[string]int
	ByCategory    map[string]int
	ArtifactTotal int
}

func BuildReport(documents []domain.ProjectDocument) Report {
	report := Report{ByStatus: map[domain.DocumentStatus]int{}, ByOwner: map[string]int{}, ByCategory: map[string]int{}}
	for _, document := range documents {
		report.Total++
		report.ByStatus[document.Status]++
		report.ByOwner[strings.ToLower(strings.TrimSpace(document.Owner))]++
		report.ArtifactTotal += document.ArtifactCount()
		for _, category := range document.Categories {
			key := strings.ToLower(strings.TrimSpace(category))
			if key != "" {
				report.ByCategory[key]++
			}
		}
	}
	return report
}

func (r Report) CompletionRate() float64 {
	if r.Total == 0 {
		return 0
	}
	return float64(r.ByStatus[domain.StatusArchived]+r.ByStatus[domain.StatusApproved]) / float64(r.Total)
}

func (r Report) StatusLine() string {
	statuses := []domain.DocumentStatus{domain.StatusDraft, domain.StatusApproved, domain.StatusRejected, domain.StatusArchived}
	parts := make([]string, 0, len(statuses))
	for _, status := range statuses {
		parts = append(parts, fmt.Sprintf("%s=%d", status, r.ByStatus[status]))
	}
	return strings.Join(parts, ",")
}

func (r Report) OwnerRanking() []string {
	owners := make([]string, 0, len(r.ByOwner))
	for owner := range r.ByOwner {
		owners = append(owners, owner)
	}
	sort.Slice(owners, func(i, j int) bool {
		if r.ByOwner[owners[i]] == r.ByOwner[owners[j]] {
			return owners[i] < owners[j]
		}
		return r.ByOwner[owners[i]] > r.ByOwner[owners[j]]
	})
	return owners
}

func (a *Application) Report(keyword, owner, category, status string) (Report, error) {
	documents, err := a.Workflow.FindDocuments(search.Parse(keyword, owner, category, status))
	if err != nil {
		return Report{}, err
	}
	return BuildReport(documents), nil
}

func (a *Application) DocumentSummary(id string) (string, error) {
	document, err := a.Workflow.LoadDocument(id)
	if err != nil {
		return "", err
	}
	return document.Summary(), nil
}
