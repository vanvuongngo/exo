package state

import "context"

type Store interface {
	DescribeComponents(context.Context, *DescribeComponentsInput) (*DescribeComponentsOutput, error)
	AddComponent(context.Context, *AddComponentInput) (*AddComponentOutput, error)
	PatchComponent(context.Context, *PatchComponentInput) (*PatchComponentOutput, error)
	RemoveComponent(context.Context, *RemoveComponentInput) (*RemoveComponentOutput, error)
}

type DescribeComponentsInput struct {
	ProjectID string `json:"projectId"`
}

type DescribeComponentsOutput struct {
	Components []ComponentDescription `json:"components"`
}

type ComponentDescription struct {
	ID          string                 `json:"id"`
	ProjectID   string                 `json:"projectId"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Spec        map[string]interface{} `json:"spec"`
	State       map[string]interface{} `json:"state"`
	Created     string                 `json:"created"`
	Initialized *string                `json:"initialized"`
	Disposed    *string                `json:"disposed"`
}

type AddComponentInput struct {
	ProjectID string                 `json:"projectId"`
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Spec      map[string]interface{} `json:"spec"`
	Created   string                 `json:"created"`
}

type AddComponentOutput struct{}

type PatchComponentInput struct {
	ID          string                 `json:"id"`
	State       map[string]interface{} `json:"state"`
	Initialized string                 `json:"initialized"`
	Disposed    string                 `json:"disposed"`
}

type PatchComponentOutput struct{}

type RemoveComponentInput struct {
	ID string `json:"id"`
}

type RemoveComponentOutput struct{}
