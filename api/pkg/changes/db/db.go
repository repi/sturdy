package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) Get(ctx context.Context, id changes.ID) (*changes.Change, error) {
	var res changes.Change
	err := r.db.GetContext(ctx, &res, `SELECT id, codebase_id, title, updated_description, user_id, git_creator_name, git_creator_email, created_at, git_created_at, commit_id, parent_change_id, workspace_id
		FROM changes
		WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repo) GetByCommitID(ctx context.Context, commitID string, codebaseID codebases.ID) (*changes.Change, error) {
	var res changes.Change
	err := r.db.GetContext(ctx, &res, `SELECT id, codebase_id, title, updated_description, user_id, git_creator_name, git_creator_email, created_at, git_created_at, commit_id, parent_change_id, workspace_id
		FROM changes
		WHERE commit_id = $1 AND codebase_id = $2`, commitID, codebaseID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repo) Insert(ctx context.Context, ch changes.Change) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO changes
		(id, codebase_id, title, updated_description, user_id, git_creator_name, git_creator_email, created_at, git_created_at, commit_id, parent_change_id, workspace_id)
		VALUES(:id, :codebase_id, :title, :updated_description, :user_id, :git_creator_name, :git_creator_email, :created_at, :git_created_at, :commit_id, :parent_change_id, :workspace_id)
    	`, &ch)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (r *repo) Update(ctx context.Context, ch changes.Change) error {
	_, err := r.db.NamedExecContext(ctx, `UPDATE changes
    	SET updated_description = :updated_description,
    	    title = :title,
    	    user_id = :user_id,
    	    git_creator_name = :git_creator_name,
    	    git_creator_email = :git_creator_email,
    	    created_at = :created_at,
    	    git_created_at = :git_created_at,
			commit_id = :commit_id,
    	    parent_change_id = :parent_change_id,
			workspace_id = :workspace_id
    	WHERE id = :id`, &ch)
	if err != nil {
		return fmt.Errorf("failed to update change: %w", err)
	}
	return nil
}

func (r *repo) ListByIDs(ctx context.Context, ids ...changes.ID) ([]*changes.Change, error) {
	var res []*changes.Change
	err := r.db.SelectContext(ctx, &res, `
		SELECT
			id, codebase_id, title, updated_description, user_id, git_creator_name, git_creator_email, created_at, git_created_at, commit_id, parent_change_id, workspace_id
		FROM
			changes
		WHERE
			id = ANY($1)
	`, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return res, nil
}

func (r *repo) GetByParentChangeID(ctx context.Context, parentChangeID changes.ID) (*changes.Change, error) {
	res := &changes.Change{}
	if err := r.db.GetContext(ctx, res, `
		SELECT id, codebase_id, title, updated_description, user_id, git_creator_name, git_creator_email, created_at, git_created_at, commit_id, parent_change_id, workspace_id
		FROM changes
		WHERE parent_change_id = $1
	`, parentChangeID); err != nil {
		return nil, err
	}
	return res, nil
}
