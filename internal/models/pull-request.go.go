package models

import (
	"time"
)

type PRShort struct {
	ID       string `json:"pull_request_id"   db:"id"`
	Title    string `json:"pull_request_name" db:"title"`
	AuthorID string `json:"author_id"         db:"author_id"`
	Status   Status `json:"status"            db:"status"`
}
type PR struct {
	PRShort
	AssignedReviewers []string  `json:"assigned_reviewers"`
	CreatedAt         time.Time `json:"createdAt" db:"title"`
	MergedAt          time.Time `json:"mergedAt,omitzero" db:"title"`
}

type PRCreateRequest struct {
	ID       string `json:"pull_request_id"   db:"id"`
	Title    string `json:"pull_request_name" db:"title"`
	AuthorID string `json:"author_id"         db:"author_id"`
}

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type PullRequestResponse struct {
	PullRequestID     string    `json:"pull_request_id"`
	PullRequestName   string    `json:"pull_request_name"`
	AuthorID          string    `json:"author_id"`
	Status            Status    `json:"status"`
	AssignedReviewers []string  `json:"assigned_reviewers"`
	NeedMoreReviewers bool      `json:"need_more_reviewers"`
	CreatedAt         time.Time `json:"createdAt"`
	MergedAt          time.Time `json:"mergedAt,omitempty"`
}

type UserReviewResponse struct {
	UserID       string    `json:"user_id"`
	PullRequests []PRShort `json:"pull_requests"`
}
