package db

import (
	"context"
	"errors"
	"time"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/collections"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool           *pgxpool.Pool
	Auth           auth.Repository
	Collections    collections.Repository
	Templates      *TemplateRepository
	Usage          *UsageRepository
	Abuse          *AbuseRepository
	BlockedTargets *BlockedTargetRepository
}

func Open(ctx context.Context, databaseURL string, maxConns int32) (*Store, error) {
	if databaseURL == "" {
		return nil, nil
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	if maxConns > 0 {
		config.MaxConns = maxConns
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}

	return NewStore(pool), nil
}

func NewStore(pool *pgxpool.Pool) *Store {
	if pool == nil {
		return nil
	}

	return &Store{
		pool:           pool,
		Auth:           NewAuthRepository(pool),
		Collections:    NewCollectionRepository(pool),
		Templates:      NewTemplateRepository(pool),
		Usage:          NewUsageRepository(pool),
		Abuse:          NewAbuseRepository(pool),
		BlockedTargets: NewBlockedTargetRepository(pool),
	}
}

func (s *Store) Pool() *pgxpool.Pool {
	if s == nil {
		return nil
	}

	return s.pool
}

func (s *Store) Close() {
	if s == nil || s.pool == nil {
		return
	}

	s.pool.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.pool == nil {
		return errors.New("database store is not initialized")
	}

	return s.pool.Ping(ctx)
}
