package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/p3rch1/review-manager/internal/models"
)

const REVIEWERS_COUNT = 2

func (r *reviewDB) PRCreate(ctx context.Context, request *models.PRCreateRequest) (*models.PR, error) {
	const query = `
		WITH inserted_pr AS (
			INSERT INTO pull_requests (id, title, author_id)
			SELECT $1, $2, u.id
			FROM users u
			WHERE u.id = $3
			RETURNING id, author_id, title, status, created_at
		),

		assigned_reviewers AS (
			INSERT INTO pr_reviewers (pr_id, user_id)
			SELECT ip.id, u.id
			FROM inserted_pr ip
			JOIN users u ON u.team_id = (SELECT team_id FROM users WHERE id = $3)
			WHERE u.id <> $3 AND u.is_active = TRUE
			ORDER BY random()
			LIMIT $4
			RETURNING user_id
		)

		SELECT 
			ip.id,
			ip.author_id,
			ip.title,
			ip.status,
			ip.created_at,
			COALESCE(
				array_agg(ar.user_id) FILTER (WHERE ar.user_id IS NOT NULL),
				ARRAY[]::text[]) AS assigned_reviewers
		FROM inserted_pr ip
		LEFT JOIN assigned_reviewers ar ON true
		GROUP BY ip.id, ip.author_id, ip.title, ip.status, ip.created_at;
	`

	var pr models.PR
	var reviewers pq.StringArray

	row := r.db.QueryRowContext(ctx, query, request.ID, request.Title, request.AuthorID, REVIEWERS_COUNT)
	err := row.Scan(
		&pr.ID,
		&pr.AuthorID,
		&pr.Title,
		&pr.Status,
		&pr.CreatedAt,
		&reviewers,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		if isUniqueViolation(err) {
			return nil, models.ErrPRExists
		}
		return nil, fmt.Errorf("create pr: %w", err)
	}

	pr.AssignedReviewers = []string(reviewers)

	return &pr, nil
}
