package resolvers

import (
	"context"
	"database/sql"
	"errors"
)

type ComponentResolver struct {
	Q *QueryResolver
	ComponentRow
}

type ComponentRow struct {
	ID      string `db:"id"`
	StackID string `db:"stack_id"`
	Name    string `db:"name"`
}

func (r *QueryResolver) ComponentByID(ctx context.Context, args struct {
	ID string
}) (*ComponentResolver, error) {
	return r.componentByID(ctx, &args.ID)
}

func (r *QueryResolver) componentByID(ctx context.Context, id *string) (*ComponentResolver, error) {
	component := &ComponentResolver{}
	err := r.getRowByKey(ctx, &component.ComponentRow, `
		SELECT *
		FROM component
		WHERE id = ?
	`, id)
	if component.ID == "" {
		component = nil
	}
	return component, err
}

func (r *QueryResolver) componentByName(ctx context.Context, stack string, name string) (*ComponentResolver, error) {
	stackResolver, err := r.stackByRef(ctx, stack)
	if stackResolver == nil || err != nil {
		return nil, err
	}
	stackID := stackResolver.ID

	component := &ComponentResolver{}
	err = r.DB.GetContext(ctx, &component.ComponentRow, `
		SELECT *
		FROM component
		WHERE stack_id = ?
		AND name = ?
	`, stackID, name)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	if component.ID == "" {
		component = nil
	}
	return component, err
}

func (r *QueryResolver) ComponentByRef(ctx context.Context, args struct {
	Ref   string
	Stack *string
}) (*ComponentResolver, error) {
	return r.componentByRef(ctx, args.Ref, args.Stack)
}

func (r *QueryResolver) componentByRef(ctx context.Context, ref string, stack *string) (*ComponentResolver, error) {
	component, err := r.componentByID(ctx, &ref)
	if component != nil || err != nil {
		return component, err
	}
	if stack != nil {
		component, err = r.componentByName(ctx, *stack, ref)
	}
	return component, err
}

func (r *ComponentResolver) Stack(ctx context.Context) (*StackResolver, error) {
	return r.Q.stackByID(ctx, &r.StackID)
}

func (r *ComponentResolver) Resources(ctx context.Context) ([]*ResourceResolver, error) {
	return r.Q.resourcesByComponent(ctx, r.ID)
}
