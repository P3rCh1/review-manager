package models

type StatsResponse struct {
	Service *ServiceStats `json:"service"`
	User    *UserStats    `json:"user"`
	PR      *PRStats      `json:"pr"`
}

type ServiceStats struct {
	TotalUsers  int `json:"total_users"  db:"total_users"`
	ActiveUsers int `json:"active_users" db:"active_users"`
	TotalTeams  int `json:"total_teams"  db:"total_teams"`
	MergedPRs   int `json:"merged_prs"   db:"merged_prs"`
	OpenPRs     int `json:"open_prs"     db:"open_prs"`
}

type UserStats struct {
	AvgReviewsPerActiveUser    float64 `json:"avg_reviews_per_active_user"    db:"avg_reviews_per_active_user"`
	MaxReviewsOnUser           int     `json:"max_reviews_on_user"            db:"max_reviews_on_user"`
	ActiveUsersWithZeroReviews int     `json:"active_users_with_zero_reviews" db:"active_users_with_zero_reviews"`
}

type PRStats struct {
	OpenPRsWithZeroReviewers int `json:"open_prs_with_0_reviewers" db:"open_prs_with_0_reviewers"`
	OpenPRsWithOneReviewer   int `json:"open_prs_with_1_reviewer"  db:"open_prs_with_1_reviewer"`
	OpenPRsWithTwoReviewers  int `json:"open_prs_with_2_reviewers" db:"open_prs_with_2_reviewers"`
}
