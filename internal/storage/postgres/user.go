package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/p3rch1/review-manager/internal/models"
)

func (r *reviewDB) SetIsActive(ctx context.Context, request *models.SetActiveRequest) (*models.User, error) {
	const query = `
    	UPDATE users u
    	SET is_active = $1
    	FROM teams t
    	WHERE u.id = $2 AND u.team_id = t.id
    	RETURNING u.id, u.username, u.is_active, t.name AS team_name
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, request.IsActive, request.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("update user: %w", err)
	}

	return &user, nil

}
