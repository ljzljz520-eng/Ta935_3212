package policy

import (
	"fmt"
	"strings"

	"engineering-document-vault/internal/catalog"
	"engineering-document-vault/internal/domain"
)

type Check struct {
	Code   string
	Passed bool
	Reason string
}

type Checklist struct {
	Checks []Check
}

func BuildChecklist(document domain.ProjectDocument, sequence int) Checklist {
	checks := []Check{}
	checks = append(checks, Check{Code: "identity", Passed: strings.TrimSpace(document.ID) != "", Reason: "document identity is present"})
	checks = append(checks, Check{Code: "title", Passed: strings.TrimSpace(document.Title) != "", Reason: "title is present"})
	checks = append(checks, Check{Code: "owner", Passed: strings.TrimSpace(document.Owner) != "", Reason: "owner is present"})
	checks = append(checks, Check{Code: "category", Passed: len(document.Categories) > 0, Reason: "at least one category is selected"})
	checks = append(checks, Check{Code: "artifact", Passed: len(document.Artifacts) > 0, Reason: "supporting artifact is attached"})
	checks = append(checks, Check{Code: "known-category", Passed: knownCategories(document.Categories), Reason: "categories are registered"})
	checks = append(checks, Check{Code: "sequence", Passed: sequence > 0, Reason: "review sequence is positive"})
	checks = append(checks, Check{Code: "revision", Passed: document.Revision > 0, Reason: "revision is positive"})
	checks = append(checks, Check{Code: "visibility", Passed: document.IsVisible(), Reason: "document can be returned to operators"})
	return Checklist{Checks: checks}
}

func knownCategories(categories []string) bool {
	if len(categories) == 0 {
		return false
	}
	for _, category := range categories {
		if !catalog.IsKnownCategory(category) {
			return false
		}
	}
	return true
}

func (c Checklist) Passed() bool {
	for _, check := range c.Checks {
		if !check.Passed {
			return false
		}
	}
	return true
}

func (c Checklist) FailedCodes() []string {
	result := []string{}
	for _, check := range c.Checks {
		if !check.Passed {
			result = append(result, check.Code)
		}
	}
	return result
}

func (c Checklist) Explanation() string {
	parts := []string{}
	for _, check := range c.Checks {
		state := "pass"
		if !check.Passed {
			state = "fail"
		}
		parts = append(parts, fmt.Sprintf("%s=%s", check.Code, state))
	}
	return strings.Join(parts, ",")
}

func RejectionMessage(checklist Checklist, sequence int) string {
	if checklist.Passed() {
		return ""
	}
	failed := strings.Join(checklist.FailedCodes(), ", ")
	return fmt.Sprintf("review %d rejected: %s", sequence, failed)
}

func ReviewOutcome(document domain.ProjectDocument, sequence int) (domain.ReviewDecision, string) {
	checklist := BuildChecklist(document, sequence)
	if !checklist.Passed() {
		return domain.DecisionReject, RejectionMessage(checklist, sequence)
	}
	if sequence%7 == 0 {
		return domain.DecisionReject, fmt.Sprintf("review %d rejected by compliance gate", sequence)
	}
	return domain.DecisionApprove, "all checklist controls passed"
}
