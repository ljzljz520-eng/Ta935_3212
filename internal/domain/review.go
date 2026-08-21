package domain

import "errors"

type ReviewDecision string

const (
	DecisionApprove ReviewDecision = "approve"
	DecisionReject  ReviewDecision = "reject"
)

type ReviewRecord struct {
	ID         string         `json:"id"`
	DocumentID string         `json:"document_id"`
	Sequence   int            `json:"sequence"`
	Decision   ReviewDecision `json:"decision"`
	Reason     string         `json:"reason"`
	Reviewer   string         `json:"reviewer"`
}

func NewReview(id, documentID, reviewer string, sequence int, decision ReviewDecision, reason string) (ReviewRecord, error) {
	if id == "" || documentID == "" || reviewer == "" || reason == "" {
		return ReviewRecord{}, errors.New("review fields are required")
	}
	if sequence < 1 {
		return ReviewRecord{}, errors.New("review sequence must be positive")
	}
	if decision != DecisionApprove && decision != DecisionReject {
		return ReviewRecord{}, errors.New("review decision is invalid")
	}
	return ReviewRecord{ID: id, DocumentID: documentID, Sequence: sequence, Decision: decision, Reason: reason, Reviewer: reviewer}, nil
}

func (r ReviewRecord) Approved() bool { return r.Decision == DecisionApprove }

func (r ReviewRecord) Rejected() bool { return r.Decision == DecisionReject }
