package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/p3rch1/review-manager/internal/models"
)

func (r *reviewDB) CreatePR(ctx context.Context, req *models.PRCreateRequest, reviewersCount int) (*models.PR, error) {
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
				ARRAY[]::text[]
			) AS assigned_reviewers
		FROM inserted_pr ip
		LEFT JOIN assigned_reviewers ar ON true
		GROUP BY ip.id, ip.author_id, ip.title, ip.status, ip.created_at;
	`

	var pr models.PR
	var reviewers pq.StringArray

	row := r.db.QueryRowContext(ctx, query, req.ID, req.Title, req.AuthorID, reviewersCount)
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

		return nil, fmt.Errorf("create pr DB: %w", err)
	}

	pr.AssignedReviewers = []string(reviewers)

	return &pr, nil
}

func (r *reviewDB) Merge(ctx context.Context, req *models.MergeRequest) (*models.PR, error) {
	const query = `
		WITH updated_pr AS (
			UPDATE pull_requests
			SET status = 'MERGED', merged_at = COALESCE(merged_at, now())
			WHERE id = $1
			RETURNING *
		)

		SELECT 
			up.id,
			up.title,
			up.author_id,
			up.status,
			up.created_at,
			up.merged_at,
			COALESCE(
				array_agg(prr.user_id) FILTER (WHERE prr.user_id IS NOT NULL),
				ARRAY[]::text[]
			) AS assigned_reviewers
		FROM updated_pr up
		LEFT JOIN pr_reviewers prr ON up.id = prr.pr_id
		GROUP BY up.id, up.title, up.author_id, up.status, up.created_at, up.merged_at
	`

	var assigned pq.StringArray
	var pr models.PR

	if err := r.db.QueryRowContext(ctx, query, req.ID).Scan(
		&pr.ID,
		&pr.Title,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
		&assigned,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPRNotFound
		}

		return nil, fmt.Errorf("fetch info: %w", err)
	}

	pr.AssignedReviewers = []string(assigned)
	return &pr, nil
}

func (r *reviewDB) ReassignPR(ctx context.Context, req *models.ReassignRequest) (*models.ReassignResponce, error) {
	const infoQuery = `
		SELECT
			pr.id,
			pr.title,
			pr.author_id,
			pr.status,
			pr.created_at,
			EXISTS(SELECT 1 FROM users WHERE id = $2),
			COALESCE(
				array_agg(prr.user_id) FILTER (WHERE prr.user_id IS NOT NULL),
				ARRAY[]::text[]
			) AS assigned_reviewers
		FROM pull_requests pr
		LEFT JOIN pr_reviewers prr ON pr.id = prr.pr_id
		WHERE pr.id = $1
		GROUP BY pr.id
	`

	const selectNew = `
		WITH candidates AS (
    		SELECT u.id
    		FROM users u
    		JOIN users old ON old.id = $1
    		WHERE u.team_id = old.team_id
      			AND u.is_active = TRUE
      			AND u.id <> ALL($3)
      			AND u.id <> $2
		)
		SELECT id
		FROM candidates
		OFFSET floor(random() * (SELECT count(*) FROM candidates))
		LIMIT 1
	`

	const update = `
		UPDATE pr_reviewers
		SET user_id = $1
		WHERE pr_id = $2 AND user_id = $3
	`

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var userExists bool
	var assigned pq.StringArray
	var resp models.ReassignResponce

	if err := tx.QueryRowContext(ctx, infoQuery, req.PRID, req.OldReviewerID).Scan(
		&resp.PR.ID,
		&resp.PR.Title,
		&resp.PR.AuthorID,
		&resp.PR.Status,
		&resp.PR.CreatedAt,
		&userExists,
		&assigned,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPRNotFound
		}

		return nil, fmt.Errorf("fetch info: %w", err)
	}

	if !userExists {
		return nil, models.ErrUserNotFound
	}

	if resp.PR.Status == models.StatusMerged {
		return nil, models.ErrPRMerged
	}

	resp.PR.AssignedReviewers = []string(assigned)
	replacingIndex := -1
	for i, r := range resp.PR.AssignedReviewers {
		if r == req.OldReviewerID {
			replacingIndex = i
			break
		}
	}
	if replacingIndex == -1 {
		return nil, models.ErrNotAssigned
	}

	err = tx.QueryRowContext(ctx, selectNew, req.OldReviewerID, resp.PR.AuthorID, pq.Array(resp.PR.AssignedReviewers)).Scan(&resp.ReplacedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoCandidate
		}

		return nil, fmt.Errorf("select new reviewer: %w", err)
	}

	if _, err := tx.ExecContext(ctx, update, resp.ReplacedBy, req.PRID, req.OldReviewerID); err != nil {
		return nil, fmt.Errorf("update reviewer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp.PR.AssignedReviewers[replacingIndex] = resp.ReplacedBy
	return &resp, nil
}
