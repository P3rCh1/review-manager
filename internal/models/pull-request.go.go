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
	AssignedReviewers []string  `json:"assigned_reviewers" db:"assigned_reviewers"`
	CreatedAt         time.Time `json:"createdAt"          db:"created_at"`
	MergedAt          time.Time `json:"mergedAt,omitzero"  db:"merged_at"`
}

type PRCreateRequest struct {
	ID       string `json:"pull_request_id"   db:"id"`
	Title    string `json:"pull_request_name" db:"title"`
	AuthorID string `json:"author_id"         db:"author_id"`
}

type MergeRequest struct {
	ID string `json:"pull_request_id" db:"id"`
}

type ReassignRequest struct {
	PRID          string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type ReassignResponce struct {
	PR         PR     `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}

type UserReviewResponse struct {
	UserID       string    `json:"user_id"`
	PullRequests []PRShort `json:"pull_requests"`
}
