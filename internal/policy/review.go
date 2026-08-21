package policy

import (
	"errors"
	"fmt"

	"engineering-document-vault/internal/domain"
)

type ReviewPolicy struct {
	RequiredArtifact bool
	MinCategories    int
}

func DefaultReviewPolicy() ReviewPolicy {
	return ReviewPolicy{RequiredArtifact: true, MinCategories: 1}
}

func (p ReviewPolicy) Evaluate(document domain.ProjectDocument, sequence int) (domain.ReviewDecision, string) {
	if len(document.Categories) < p.MinCategories {
		return domain.DecisionReject, "category is missing"
	}
	if p.RequiredArtifact && len(document.Artifacts) == 0 {
		return domain.DecisionReject, "artifact is missing"
	}
	if sequence%7 == 0 {
		return domain.DecisionReject, fmt.Sprintf("review %d rejected by compliance gate", sequence)
	}
	return domain.DecisionApprove, "document satisfies review policy"
}

func ValidateReview(document domain.ProjectDocument, sequence int) error {
	if err := document.Validate(); err != nil {
		return err
	}
	if sequence < 1 {
		return errors.New("review sequence must be positive")
	}
	return nil
}

func CanArchive(document domain.ProjectDocument) error {
	if document.Status != domain.StatusApproved {
		return errors.New("document must be approved before archive")
	}
	return nil
}
