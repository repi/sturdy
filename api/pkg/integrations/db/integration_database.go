package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/integrations"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var _ IntegrationsRepository = &configDatabase{}

type configDatabase struct {
	db *sqlx.DB
}

func NewIntegrationDatabase(db *sqlx.DB) IntegrationsRepository {
	return &configDatabase{db: db}
}

func (cd *configDatabase) Create(ctx context.Context, cfg *integrations.Integration) error {
	if _, err := cd.db.ExecContext(ctx, `
		INSERT INTO integrations
			(id, codebase_id, provider, provider_type, seed_files, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6)
		`, cfg.ID, cfg.CodebaseID, cfg.Provider, cfg.ProviderType, pq.Array(cfg.SeedFiles), cfg.CreatedAt, cfg.UpdatedAt); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (cd *configDatabase) Update(ctx context.Context, cfg *integrations.Integration) error {
	if _, err := cd.db.ExecContext(ctx, `
		UPDATE integrations
		SET
			provider = $2,
			seed_files = $3,
			updated_at = $4,
		    deleted_at = $5
		WHERE
			id = $1
		`, cfg.ID, cfg.Provider, pq.Array(cfg.SeedFiles), cfg.UpdatedAt, cfg.DeletedAt); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (cd *configDatabase) ListByCodebaseID(ctx context.Context, codebaseID codebases.ID) ([]*integrations.Integration, error) {
	rows, err := cd.db.QueryContext(ctx, `
		SELECT
			id, codebase_id, provider, provider_type, seed_files, created_at, updated_at
		FROM
			integrations
		WHERE
			codebase_id = $1
			AND deleted_at IS NULL
	`, codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	cfgs := []*integrations.Integration{}
	for rows.Next() {
		cfg := &integrations.Integration{}
		if err := rows.Scan(&cfg.ID, &cfg.CodebaseID, &cfg.Provider, &cfg.ProviderType, pq.Array(&cfg.SeedFiles), &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		cfgs = append(cfgs, cfg)
	}

	return cfgs, nil
}

func (cd *configDatabase) Get(ctx context.Context, id string) (*integrations.Integration, error) {
	row := cd.db.QueryRowContext(ctx, `
		SELECT
			id, codebase_id, provider, provider_type, seed_files, created_at, updated_at, deleted_at
		FROM
			integrations
		WHERE
			id = $1
	`, id)
	var cfg integrations.Integration
	if err := row.Scan(&cfg.ID, &cfg.CodebaseID, &cfg.Provider, &cfg.ProviderType, pq.Array(&cfg.SeedFiles), &cfg.CreatedAt, &cfg.UpdatedAt, &cfg.DeletedAt); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}
	return &cfg, nil
}
