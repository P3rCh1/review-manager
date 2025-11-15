package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/p3rch1/review-manager/internal/models"
)

func (r *reviewDB) SetIsActive(ctx context.Context, req *models.SetActiveRequest) (*models.User, error) {
	const query = `
    	UPDATE users u
    	SET is_active = $1
    	FROM teams t
    	WHERE u.id = $2 AND u.team_id = t.id
    	RETURNING u.id, u.username, u.is_active, t.name AS team_name
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, req.IsActive, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, fmt.Errorf("update user: %w", err)
	}

	return &user, nil
}

func (r *reviewDB) GetReviews(ctx context.Context, id string) ([]models.PRShort, error) {
	const query = `
    	SELECT 	
			COALESCE(
				json_agg (
					json_build_object(
						'pull_request_id', 		pr.id,
						'author_id', 			pr.author_id,
						'status', 				pr.status,
						'pull_request_name', 	pr.title
					)
				),
				'[]'
			) reviews,
			EXISTS(SELECT 1 FROM users WHERE id = $1) user_exists		
		FROM pull_requests pr
		JOIN pr_reviewers r ON pr.id = r.pr_id
		WHERE r.user_id = $1
	`

	queryRes := struct {
		ReviewsJSON string `db:"reviews"`
		UserExists  bool   `db:"user_exists"`
	}{}
	if err := r.db.GetContext(ctx, &queryRes, query, id); err != nil {
		return nil, fmt.Errorf("select reviews: %w", err)
	}

	if !queryRes.UserExists {
		return nil, models.ErrUserNotFound
	}

	res := []models.PRShort{}
	if err := json.Unmarshal([]byte(queryRes.ReviewsJSON), &res); err != nil {
		return nil, fmt.Errorf("unmarshal reviews: %w", err)
	}

	return res, nil
}
