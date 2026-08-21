package search

import (
	"strings"

	"engineering-document-vault/internal/domain"
)

type Query struct {
	Keyword  string
	Owner    string
	Status   domain.DocumentStatus
	Category string
}

func (q Query) Match(document domain.ProjectDocument) bool {
	if q.Status != "" && document.Status != q.Status {
		return false
	}
	if q.Owner != "" && !strings.EqualFold(document.Owner, q.Owner) {
		return false
	}
	if q.Category != "" && !document.HasCategory(q.Category) {
		return false
	}
	if q.Keyword != "" {
		needle := strings.ToLower(q.Keyword)
		if !strings.Contains(strings.ToLower(document.Title), needle) && !strings.Contains(strings.ToLower(strings.Join(document.Artifacts, " ")), needle) {
			return false
		}
	}
	return document.IsVisible()
}

func Filter(documents []domain.ProjectDocument, query Query) []domain.ProjectDocument {
	result := make([]domain.ProjectDocument, 0, len(documents))
	for _, document := range documents {
		if query.Match(document) {
			result = append(result, domain.CloneDocument(document))
		}
	}
	domain.SortDocuments(result)
	return result
}

func Parse(keyword, owner, category, status string) Query {
	return Query{Keyword: strings.TrimSpace(keyword), Owner: strings.TrimSpace(owner), Category: strings.TrimSpace(category), Status: domain.DocumentStatus(strings.TrimSpace(status))}
}
