package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/p3rch1/review-manager/internal/models"
)

func (r *reviewDB) AddTeam(ctx context.Context, team *models.Team) error {
	const (
		insertTeamQuery = `
			INSERT INTO teams (name)
			VALUES ($1)
			RETURNING id
		`
		upsertUsersQuery = `
			INSERT INTO users (id, username, is_active, team_id)
			VALUES (:id, :username, :is_active, :team_id)
			ON CONFLICT (id)
			DO UPDATE SET
				username = EXCLUDED.username,
				is_active = EXCLUDED.is_active,
				team_id  = EXCLUDED.team_id
		`
	)

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer tx.Rollback()

	var teamID uuid.UUID
	if err := tx.QueryRowContext(ctx, insertTeamQuery, team.TeamName).Scan(&teamID); err != nil {
		if isUniqueViolation(err) {
			return models.ErrTeamExists
		}
		return fmt.Errorf("insert team: %w", err)
	}

	members := make([]models.UserDB, len(team.Members))
	for i, m := range team.Members {
		members[i] = models.UserDB{
			TeamMember: m,
			TeamID:     teamID,
		}
	}

	if _, err := tx.NamedExecContext(ctx, upsertUsersQuery, members); err != nil {
		return fmt.Errorf("upsert users: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (r *reviewDB) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	const selecQuery = `
		SELECT u.id, u.username, u.is_active
		FROM teams t
		LEFT JOIN users u ON t.id = u.team_id
		WHERE t.name = $1
	`
	const checkExistsQuery = `
		SELECT EXISTS(
			SELECT 1 FROM teams
			WHERE name = $1
		)
	`
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}
	defer tx.Rollback()

	var members []models.TeamMember
	if err := tx.SelectContext(ctx, &members, selecQuery, name); err != nil {
		return nil, fmt.Errorf("select members: %w", err)
	}

	if len(members) == 0 {
		var exists bool
		err := tx.GetContext(ctx, &exists, checkExistsQuery, name)

		if err != nil {
			return nil, fmt.Errorf("check team exists: %w", err)
		}

		if !exists {
			return nil, models.ErrNotFound
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &models.Team{
		TeamName: name,
		Members:  members,
	}, nil
}
