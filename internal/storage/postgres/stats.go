package postgres

import (
	"fmt"

	"github.com/p3rch1/review-manager/internal/models"
)

func (r *reviewDB) ServiceStats() (*models.ServiceStats, error) {
	const query = `
        SELECT
            COUNT(*) AS total_users,
            COUNT(*) FILTER (WHERE is_active = true) AS active_users,
            COUNT(DISTINCT team_id) AS total_teams,
            (SELECT COUNT(*) FROM pull_requests WHERE status = 'MERGED') AS merged_prs,
            (SELECT COUNT(*) FROM pull_requests WHERE status = 'OPEN') AS open_prs
        FROM users;
    `

	var stats models.ServiceStats
	if err := r.db.Get(&stats, query); err != nil {
		return nil, fmt.Errorf("get service stats: %w", err)
	}

	return &stats, nil
}

func (r *reviewDB) UserStats() (*models.UserStats, error) {
	const query = `
        WITH reviews_per_user AS (
            SELECT user_id, COUNT(*) AS cnt
            FROM pr_reviewers
            GROUP BY user_id
        )

        SELECT
            COALESCE(AVG(cnt), 0) AS avg_reviews_per_active_user,
            COALESCE(MAX(cnt), 0) AS max_reviews_on_user,
            (
				SELECT COUNT(*) 
             	FROM users u
             	WHERE is_active = true AND NOT EXISTS (
                	SELECT 1 FROM reviews_per_user rpu WHERE rpu.user_id = u.id
            	)
			) AS active_users_with_zero_reviews
        FROM reviews_per_user;
    `
	var stats models.UserStats
	if err := r.db.Get(&stats, query); err != nil {
		return nil, fmt.Errorf("get user stats: %w", err)
	}

	return &stats, nil
}

func (r *reviewDB) PRStats() (*models.PRStats, error) {
	const query = `
        WITH reviews_per_pr AS (
            SELECT pr_id, COUNT(*) AS cnt
            FROM pr_reviewers
            GROUP BY pr_id
        )
        SELECT
            (
				SELECT COUNT(*) 
             	FROM pull_requests p
             	WHERE status = 'OPEN' AND NOT EXISTS (
                	SELECT 1 FROM reviews_per_pr rpp WHERE rpp.pr_id = p.id
             	)
			) AS open_prs_with_0_reviewers,

            (
				SELECT COUNT(*) 
             	FROM pull_requests p
             	WHERE status = 'OPEN' AND EXISTS (
                	SELECT 1 FROM reviews_per_pr rpp WHERE rpp.pr_id = p.id AND rpp.cnt = 1
             	)
			) AS open_prs_with_1_reviewer,

            (
				SELECT COUNT(*) 
             	FROM pull_requests p
             	WHERE status = 'OPEN' AND EXISTS (
                 	SELECT 1 FROM reviews_per_pr rpp WHERE rpp.pr_id = p.id AND rpp.cnt = 2
            	)
			) AS open_prs_with_2_reviewers;
    `

	var stats models.PRStats
	if err := r.db.Get(&stats, query); err != nil {
		return nil, fmt.Errorf("get PR stats: %w", err)
	}

	return &stats, nil
}
