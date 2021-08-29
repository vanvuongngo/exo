// References:
// https://github.com/compose-spec/compose-spec/blob/master/spec.md
// https://docs.docker.com/compose/compose-file/compose-file-v3/
// https://github.com/docker/compose/blob/4a51af09d6cdb9407a6717334333900327bc9302/compose/config/compose_spec.json
//
// Fields enumerated as of July 17, 2021 with from the following spec file:
// <https://github.com/compose-spec/compose-spec/blob/5141aafafa6ea03fcf52eb2b44218408825ab480/spec.md>.

package compose

import (
	"fmt"
	"io"
	"strconv"

	"code.cloudfoundry.org/bytefmt"
	"github.com/deref/exo/internal/providers/docker/compose/interpolate"
	"github.com/goccy/go-yaml"
)

type Environment = interpolate.Environment
type MapEnvironment = interpolate.MapEnvironment

func Unmarshal(r io.Reader, v interface{}, env Environment) error {
	// Decode into generic data types.
	dec := yaml.NewDecoder(r,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(), // TODO: Handle this more gracefully.
	)
	var raw interface{}
	if err := dec.Decode(&raw); err != nil {
		return err
	}

	// Interpolate and then convert back to yaml.
	// Ideally, we could bypass the extra marshal/unmarshal pair, but the
	// interface unmarshalling internals are not exposed from the yaml library.
	if err := interpolate.Interpolate(raw, env); err != nil {
		return fmt.Errorf("interpolating: %w", err)
	}
	interpolatedBytes, err := yaml.Marshal(raw)
	if err != nil {
		// Should be unreachable, but potentially possible if there is some
		// weird structure that can be unmarshalled but not marshalled.
		return fmt.Errorf("intermediate remarshaling: %w", err)
	}

	// Decode into stronger data types.
	if err := yaml.Unmarshal(interpolatedBytes, v); err != nil {
		return err
	}
	return nil
}

type ProjectTemplate struct {
	Services map[string]ServiceTemplate `yaml:"services"`
	Volumes  map[string]VolumeTemplate  `yaml:"volumes"`
	Networks map[string]NetworkTemplate `yaml:"networks"`
	MapSlice yaml.MapSlice              `yaml:",inline"`
}

type Project struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]Network `yaml:"networks"`
	Volumes  map[string]Volume  `yaml:"volumes"`
	Configs  map[string]Config  `yaml:"configs"`
	Secrets  map[string]Secret  `yaml:"secrets"`
	// TODO: extensions with "x-" prefix.
}

// This is a temporary placeholder for fields that we presently don't support,
// but are safe to ignore.
// TODO: Eliminate all usages of this with actual parsing logic.
type IgnoredField struct{}

func (ignored *IgnoredField) UnmarshalYAML(b []byte) error {
	return nil
}

type MemoryField int64

func (memory *MemoryField) UnmarshalYAML(b []byte) error {
	memString := string(b)
	memBytes, err := strconv.ParseInt(memString, 10, 64)
	if err == nil {
		*memory = MemoryField(memBytes)
		return nil
	}

	uMemBytes, err := bytefmt.ToBytes(memString)
	if err == nil {
		*memory = MemoryField(uMemBytes)
		return nil
	}

	return fmt.Errorf("could not unmarshal memory value %s: %w", b, err)
}

type NetworkTemplate struct {
	Name     string        `yaml:"name"`
	Driver   string        `yaml:"driver"`
	MapSlice yaml.MapSlice `yaml:",inline"`
}

type Network struct {
	// Name is the actual name of the docker network. The docker-compose network name, which can
	// be referenced by individual services, is the component name.
	Name       string            `yaml:"name"`
	Driver     string            `yaml:"driver"`
	DriverOpts map[string]string `yaml:"driver_opts"`
	Attachable bool              `yaml:"attachable"`
	EnableIPv6 bool              `yaml:"enable_ipv6"`
	Internal   bool              `yaml:"internal"`
	Labels     Dictionary        `yaml:"labels"`
	External   bool              `yaml:"external"`
}

type VolumeTemplate struct {
	Name     string        `yaml:"name"`
	MapSlice yaml.MapSlice `yaml:",inline"`
}

type Volume struct {
	Driver     string            `yaml:"driver"`
	DriverOpts map[string]string `yaml:"driver_opts"`
	// TODO: external
	Labels Dictionary `yaml:"labels"`
	Name   string     `yaml:"name"`
}

type Config struct {
	File     string `yaml:"file"`
	External bool   `yaml:"external"`
	Name     string `yaml:"name"`
}

type Secret struct {
	File     string `yaml:"file"`
	External bool   `yaml:"external"`
	Name     string `yaml:"name"`
}
