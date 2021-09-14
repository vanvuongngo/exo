package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/core/state/api"
	"github.com/deref/exo/internal/util/pathutil"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

var _ api.Store = (*Store)(nil)

func Open(ctx context.Context, path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	sto := &Store{
		db: db,
	}
	if _, err := sto.db.ExecContext(ctx, ddlStatements); err != nil {
		return nil, err
	}
	return sto, nil
}

func (sto *Store) DescribeWorkspaces(ctx context.Context, input *api.DescribeWorkspacesInput) (*api.DescribeWorkspacesOutput, error) {
	var rows *sql.Rows
	var err error
	if input.IDs == nil {
		rows, err = sto.db.QueryContext(ctx, `
			SELECT id, root FROM workspace
			ORDER BY id
		`)
	} else {
		rows, err = sto.db.QueryContext(ctx, `
			SELECT id, root FROM workspace
			WHERE id IN (?)
			ORDER BY id
		`, input.IDs)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var output api.DescribeWorkspacesOutput
	for rows.Next() {
		var id, root string
		if err := rows.Scan(&id, &root); err != nil {
			return nil, fmt.Errorf("scanning workspace row: %w", err)
		}
	}
	return &output, nil
}

func (sto *Store) AddWorkspace(ctx context.Context, input *api.AddWorkspaceInput) (*api.AddWorkspaceOutput, error) {
	if _, err := sto.db.Exec(`
		INSERT INTO workspace ( id, root )
		VALUES ( ?, ? )
	`, input.ID, input.Root); err != nil {
		// XXX handle id/root conflicts, etc.
		return nil, err
	}
	return &api.AddWorkspaceOutput{}, nil
}

func (sto *Store) RemoveWorkspace(ctx context.Context, input *api.RemoveWorkspaceInput) (*api.RemoveWorkspaceOutput, error) {
	if _, err := sto.RemoveComponent(ctx, &api.RemoveComponentInput{
		ID: input.ID,
	}); err != nil {
		// XXX if is not exists, no-op.
		// XXX if non-empty, translate error to say cannot remove non-empty workspace.
		return nil, err
	}
	return &api.RemoveWorkspaceOutput{}, nil
}

func (sto *Store) ResolveWorkspace(ctx context.Context, input *api.ResolveWorkspaceInput) (*api.ResolveWorkspaceOutput, error) {
	// Resolve by ID.
	row := sto.db.QueryRow(`
		SELECT 1
		FROM workspace
		WHERE id = ?
	`, input.Ref)
	if row.Err() != sql.ErrNoRows {
		return &api.ResolveWorkspaceOutput{
			ID: &input.Ref,
		}, nil
	}

	// Load all workspaces, so we can search use custom search logic in memory.
	describeOutput, err := sto.DescribeWorkspaces(ctx, &api.DescribeWorkspacesInput{})
	if err != nil {
		return nil, err
	}
	workspaces := describeOutput.Workspaces

	// Resolve by path. Search for the deepest root prefix match.
	maxLen := 0
	found := ""
	for _, workspace := range workspaces {
		n := len(workspace.Root)
		if n > maxLen && pathutil.HasFilePathPrefix(input.Ref, workspace.Root) {
			found = workspace.ID
			maxLen = n
		}
	}
	var output api.ResolveWorkspaceOutput
	if maxLen > 0 {
		output.ID = &found
	}
	return &output, nil
}

func (sto *Store) Resolve(ctx context.Context, input *api.ResolveInput) (*api.ResolveOutput, error) {
	rows, err := sto.db.Query(`
		SELECT id
		FROM components
		WHERE parent_id = ?
		AND (
			id IN ( ? )
			OR name IN ( ? )
		)
	`, input.WorkspaceID, input.Refs, input.Refs)
}

func (sto *Store) DescribeComponents(ctx context.Context, input *api.DescribeComponentsInput) (*api.DescribeComponentsOutput, error) {
	if input.WorkspaceID == "" {
		return nil, errors.New("workspace-id is required")
	}
}

func (sto *Store) AddComponent(ctx context.Context, input *api.AddComponentInput) (*api.AddComponentOutput, error) {
	if input.WorkspaceID == "" {
		return nil, errors.New("workspace-id is required")
	}
	created := input.Created // XXX parse timestamp, convert to unix integer.
	// XXX transaction.
	_, err := sto.db.Exec(`
		INSERT INTO component ( id, parent_id, name, type, definition_version, state_version, created )
		VALUES ( ?, ?, ?, ?, 1, 1, ?)
	`, input.ID, input.WorkspaceID, input.Name, input.Type, created)
	// XXX handle error
	_, err := sto.db.Exec(`
		INSERT INTO component_definition ( component_id, version, parent_id, name, type, timestamp, spec, meta )
		VALUES ( ?, ? , ?, ?, ?, ?, ?, ? )
	`, input.ID, 1, input.WorkspaceID, input.Name, input.Type, created, spec, meta)
	// XXX handle error
	_, err := sto.db.Exec(`
		INSERT INTO component_state ( component_id, version, timestamp, value )
		VALUES ( ?, ? , ?, ?, ?, ?, ?, ? )
	`, input.ID, 1, created, state)
	// XXX handle error
	// XXX end transaction
	return &api.AddComponentOutput{}, nil
}

func (sto *Store) PatchComponent(ctx context.Context, input *api.PatchComponentInput) (*api.PatchComponentOutput, error) {
	if input.ID == "" {
		return nil, errors.New("component id is required")
	}
}

func (sto *Store) RemoveComponent(ctx context.Context, input *api.RemoveComponentInput) (*api.RemoveComponentOutput, error) {
}
