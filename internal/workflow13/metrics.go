package workflow13

import (
	"sort"
	"strings"

	"engineering-document-vault/internal/domain"
)

type WorkflowMetrics struct {
	Submitted int
	Approved  int
	Rejected  int
	Archived  int
	Artifacts int
}

func Metrics(documents []domain.ProjectDocument) WorkflowMetrics {
	metrics := WorkflowMetrics{}
	for _, document := range documents {
		metrics.Submitted++
		metrics.Artifacts += document.ArtifactCount()
		switch document.Status {
		case domain.StatusApproved:
			metrics.Approved++
		case domain.StatusRejected:
			metrics.Rejected++
		case domain.StatusArchived:
			metrics.Archived++
		}
	}
	return metrics
}

func (m WorkflowMetrics) CompletionPercent() int {
	if m.Submitted == 0 {
		return 0
	}
	return (m.Approved + m.Archived) * 100 / m.Submitted
}

func (m WorkflowMetrics) HasRejections() bool { return m.Rejected > 0 }

func (m WorkflowMetrics) NeedsAttention() bool {
	if m.Submitted == 0 {
		return false
	}
	return m.Rejected > 0 || m.CompletionPercent() < 50
}

func StatusOrder(documents []domain.ProjectDocument) []domain.ProjectDocument {
	result := make([]domain.ProjectDocument, len(documents))
	copy(result, documents)
	sort.SliceStable(result, func(i, j int) bool {
		weight := func(status domain.DocumentStatus) int {
			switch status {
			case domain.StatusRejected:
				return 0
			case domain.StatusDraft:
				return 1
			case domain.StatusApproved:
				return 2
			case domain.StatusArchived:
				return 3
			default:
				return 4
			}
		}
		if weight(result[i].Status) == weight(result[j].Status) {
			return result[i].ID < result[j].ID
		}
		return weight(result[i].Status) < weight(result[j].Status)
	})
	return result
}

func DescribeMetrics(metrics WorkflowMetrics) string {
	parts := []string{"submitted=" + itoa(metrics.Submitted), "approved=" + itoa(metrics.Approved), "rejected=" + itoa(metrics.Rejected), "archived=" + itoa(metrics.Archived), "artifacts=" + itoa(metrics.Artifacts)}
	return strings.Join(parts, ",")
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
