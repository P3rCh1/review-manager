package models

type Team struct {
	TeamName string       `json:"team_name" db:"name"`
	Members  []TeamMember `json:"members"   db:"members"`
}
