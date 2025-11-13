package models

import "time"

type PullRequestDB struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Status    Status    `db:"status"`
	AuthorID  string    `db:"author_id"`
	CreatedAt time.Time `db:"created_at"`
	MergedAt  time.Time `db:"merged_at"`
}

type ReviewerDB struct {
	PRID   string `db:"pr_id"`
	UserID string `db:"user_id"`
}
