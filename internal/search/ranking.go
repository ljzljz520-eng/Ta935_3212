package search

import (
	"sort"
	"strings"

	"engineering-document-vault/internal/domain"
)

type Result struct {
	Document domain.ProjectDocument
	Score    int
	Reasons  []string
}

func Score(document domain.ProjectDocument, query Query) Result {
	result := Result{Document: domain.CloneDocument(document), Reasons: []string{}}
	needle := strings.ToLower(strings.TrimSpace(query.Keyword))
	if needle != "" && strings.Contains(strings.ToLower(document.Title), needle) {
		result.Score += 10
		result.Reasons = append(result.Reasons, "title")
	}
	if needle != "" && strings.Contains(strings.ToLower(strings.Join(document.Artifacts, " ")), needle) {
		result.Score += 4
		result.Reasons = append(result.Reasons, "artifact")
	}
	if query.Owner != "" && strings.EqualFold(document.Owner, query.Owner) {
		result.Score += 3
		result.Reasons = append(result.Reasons, "owner")
	}
	if query.Category != "" && document.HasCategory(query.Category) {
		result.Score += 2
		result.Reasons = append(result.Reasons, "category")
	}
	if document.Status == domain.StatusArchived {
		result.Score++
		result.Reasons = append(result.Reasons, "archived")
	}
	return result
}

func Ranked(documents []domain.ProjectDocument, query Query) []Result {
	results := make([]Result, 0, len(documents))
	for _, document := range documents {
		if query.Match(document) {
			results = append(results, Score(document, query))
		}
	}
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Document.ID < results[j].Document.ID
		}
		return results[i].Score > results[j].Score
	})
	return results
}

func GroupByStatus(documents []domain.ProjectDocument) map[domain.DocumentStatus][]domain.ProjectDocument {
	groups := map[domain.DocumentStatus][]domain.ProjectDocument{}
	for _, document := range documents {
		groups[document.Status] = append(groups[document.Status], domain.CloneDocument(document))
	}
	for status := range groups {
		domain.SortDocuments(groups[status])
	}
	return groups
}

func GroupByOwner(documents []domain.ProjectDocument) map[string][]domain.ProjectDocument {
	groups := map[string][]domain.ProjectDocument{}
	for _, document := range documents {
		key := strings.ToLower(strings.TrimSpace(document.Owner))
		groups[key] = append(groups[key], domain.CloneDocument(document))
	}
	return groups
}

func ContainsStatus(documents []domain.ProjectDocument, status domain.DocumentStatus) bool {
	for _, document := range documents {
		if document.Status == status {
			return true
		}
	}
	return false
}

func CountArtifacts(documents []domain.ProjectDocument) int {
	total := 0
	for _, document := range documents {
		total += document.ArtifactCount()
	}
	return total
}
