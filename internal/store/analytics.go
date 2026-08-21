package store

import (
	"strings"

	"engineering-document-vault/internal/domain"
)

type StoreAnalytics struct {
	Documents  int
	Reviews    int
	Revisions  int
	Archives   int
	Owners     map[string]int
	Categories map[string]int
}

func (s *BoltStore) Analytics() (StoreAnalytics, error) {
	analytics := StoreAnalytics{Owners: map[string]int{}, Categories: map[string]int{}}
	documents, err := s.ListDocuments()
	if err != nil {
		return analytics, err
	}
	analytics.Documents = len(documents)
	for _, document := range documents {
		analytics.Owners[strings.ToLower(strings.TrimSpace(document.Owner))]++
		for _, category := range document.Categories {
			analytics.Categories[strings.ToLower(strings.TrimSpace(category))]++
		}
		reviews, reviewErr := s.ListReviews(document.ID)
		if reviewErr != nil {
			return analytics, reviewErr
		}
		analytics.Reviews += len(reviews)
		revisions, revisionErr := s.ListRevisions(document.ID)
		if revisionErr != nil {
			return analytics, revisionErr
		}
		analytics.Revisions += len(revisions)
		if document.Status == domain.StatusArchived {
			if _, archiveErr := s.LoadArchive(document.ID); archiveErr == nil {
				analytics.Archives++
			}
		}
	}
	return analytics, nil
}

func (a StoreAnalytics) IsEmpty() bool { return a.Documents == 0 }

func (a StoreAnalytics) ReviewCoverage() float64 {
	if a.Documents == 0 {
		return 0
	}
	return float64(a.Reviews) / float64(a.Documents)
}

func (a StoreAnalytics) ArchiveCoverage() float64 {
	if a.Documents == 0 {
		return 0
	}
	return float64(a.Archives) / float64(a.Documents)
}

func (a StoreAnalytics) TopOwner() string {
	name := ""
	count := 0
	for owner, value := range a.Owners {
		if value > count || (value == count && owner < name) {
			name, count = owner, value
		}
	}
	return name
}

func (a StoreAnalytics) TopCategory() string {
	name := ""
	count := 0
	for category, value := range a.Categories {
		if value > count || (value == count && category < name) {
			name, count = category, value
		}
	}
	return name
}

func (s *BoltStore) DocumentsForCategory(category string) ([]domain.ProjectDocument, error) {
	documents, err := s.ListDocuments()
	if err != nil {
		return nil, err
	}
	result := []domain.ProjectDocument{}
	for _, document := range documents {
		if document.HasCategory(category) {
			result = append(result, domain.CloneDocument(document))
		}
	}
	domain.SortDocuments(result)
	return result, nil
}

func (s *BoltStore) DocumentsForOwner(owner string) ([]domain.ProjectDocument, error) {
	documents, err := s.ListDocuments()
	if err != nil {
		return nil, err
	}
	result := []domain.ProjectDocument{}
	for _, document := range documents {
		if strings.EqualFold(document.Owner, owner) {
			result = append(result, domain.CloneDocument(document))
		}
	}
	domain.SortDocuments(result)
	return result, nil
}
