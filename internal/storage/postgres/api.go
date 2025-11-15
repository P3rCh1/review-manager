package postgres

import (
	"context"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/p3rch1/review-manager/internal/config"
	"github.com/p3rch1/review-manager/internal/models"
)

type ReviewAPI interface {
	AddTeam(ctx context.Context, team *models.Team) error
	GetTeam(ctx context.Context, name string) (*models.Team, error)
	SetIsActive(ctx context.Context, req *models.SetActiveRequest) (*models.User, error)
	CreatePR(ctx context.Context, req *models.PRCreateRequest, reviewersCount int) (*models.PR, error)
	Merge(ctx context.Context, req *models.MergeRequest) (*models.PR, error)
	ReassignPR(ctx context.Context, req *models.ReassignRequest) (*models.ReassignResponce, error)
	GetReviews(ctx context.Context, id string) ([]models.PRShort, error)
	io.Closer
}

type reviewDB struct {
	db *sqlx.DB
}

func NewReviewAPI(cfg *config.Postgres) (ReviewAPI, error) {
	info := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("connect postgres fail: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres fail: %w", err)
	}

	return &reviewDB{db}, nil
}

func (s *reviewDB) Close() error {
	return s.db.Close()
}
