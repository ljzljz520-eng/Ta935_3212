package workflow13

import (
	"errors"
	"sort"
	"strings"

	"engineering-document-vault/internal/domain"
	"engineering-document-vault/internal/search"
)

type WorkItem struct {
	DocumentID string
	Priority   int
	Reason     string
	Status     string
}

type ReviewQueue struct{ Items []WorkItem }

func NewQueue() *ReviewQueue { return &ReviewQueue{Items: []WorkItem{}} }

func (q *ReviewQueue) Add(document domain.ProjectDocument, reason string, priority int) error {
	if q == nil {
		return errors.New("review queue is nil")
	}
	if document.ID == "" {
		return errors.New("document id is required")
	}
	if priority < 0 {
		return errors.New("priority cannot be negative")
	}
	q.Items = append(q.Items, WorkItem{DocumentID: document.ID, Priority: priority, Reason: strings.TrimSpace(reason), Status: string(document.Status)})
	q.Sort()
	return nil
}

func (q *ReviewQueue) Sort() {
	if q == nil {
		return
	}
	sort.SliceStable(q.Items, func(i, j int) bool {
		if q.Items[i].Priority == q.Items[j].Priority {
			return q.Items[i].DocumentID < q.Items[j].DocumentID
		}
		return q.Items[i].Priority > q.Items[j].Priority
	})
}

func (q *ReviewQueue) Pop() (WorkItem, bool) {
	if q == nil || len(q.Items) == 0 {
		return WorkItem{}, false
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item, true
}

func (q *ReviewQueue) Len() int {
	if q == nil {
		return 0
	}
	return len(q.Items)
}

func (q *ReviewQueue) Contains(documentID string) bool {
	if q == nil {
		return false
	}
	for _, item := range q.Items {
		if item.DocumentID == documentID {
			return true
		}
	}
	return false
}

func (q *ReviewQueue) Remove(documentID string) bool {
	if q == nil {
		return false
	}
	for index, item := range q.Items {
		if item.DocumentID == documentID {
			q.Items = append(q.Items[:index], q.Items[index+1:]...)
			return true
		}
	}
	return false
}

func BuildQueue(documents []domain.ProjectDocument, query search.Query) *ReviewQueue {
	queue := NewQueue()
	for _, document := range documents {
		if !query.Match(document) {
			continue
		}
		priority := 1
		if document.Status == domain.StatusRejected {
			priority = 3
		}
		if document.Status == domain.StatusDraft {
			priority = 2
		}
		_ = queue.Add(document, "pending review", priority)
	}
	return queue
}

func QueueIDs(queue *ReviewQueue) []string {
	if queue == nil {
		return nil
	}
	ids := make([]string, 0, len(queue.Items))
	for _, item := range queue.Items {
		ids = append(ids, item.DocumentID)
	}
	return ids
}

func QueueHasRejected(queue *ReviewQueue) bool {
	if queue == nil {
		return false
	}
	for _, item := range queue.Items {
		if item.Status == string(domain.StatusRejected) {
			return true
		}
	}
	return false
}

func QueueSummary(queue *ReviewQueue) string {
	if queue == nil {
		return "empty"
	}
	return strings.Join(QueueIDs(queue), ",")
}
