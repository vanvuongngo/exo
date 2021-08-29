package docker

import (
	"strings"

	"github.com/deref/exo/internal/providers/core"
	"github.com/deref/exo/internal/providers/docker/compose"
	dockerclient "github.com/docker/docker/client"
)

type ComponentBase struct {
	core.ComponentBase
	Docker *dockerclient.Client
}

func (c ComponentBase) GetExoLabels() map[string]string {
	return map[string]string{
		"io.deref.exo.workspace": c.WorkspaceID,
		"io.deref.exo.component": c.ComponentID,
	}
}

func (c *ComponentBase) UnmarshalSpec(v interface{}) error {
	r := strings.NewReader(c.ComponentSpec)
	env := compose.MapEnvironment(c.WorkspaceEnvironment)
	return compose.Unmarshal(r, v, env)
}
