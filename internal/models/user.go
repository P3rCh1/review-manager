package models

import "github.com/google/uuid"

type TeamMember struct {
	ID       string `json:"user_id"   db:"id"`
	Username string `json:"username"  db:"username"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type User struct {
	TeamMember
	TeamName string `json:"team_name" db:"team_name"`
}

type UserDB struct {
	TeamMember
	TeamID uuid.UUID `db:"team_id"`
}
