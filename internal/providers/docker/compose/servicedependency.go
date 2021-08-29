package compose

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

// ServiceDependenciesTemplate represents either the short (list) form of
// `depends_on` or the full (map) form that indicates the condition under which
// this service can start. Since there are two mutally exclusive
// representations of this structure, custom marshaling and unmarshing is
// implemented.
type ServiceDependenciesTemplate []ServiceDependencyTemplate

type ServiceDependencyTemplate struct {
	Service   string
	Condition string
}

type ServiceDependency struct {
	Condition string `yaml:"condition"`
}

func (sd ServiceDependenciesTemplate) MarshalYAML() (interface{}, error) {
	services := make(map[string]interface{}, len(sd))
	for _, service := range sd {
		services[service.Service] = map[string]interface{}{
			"condition": service.Condition,
		}
	}
	return services, nil
}

func (sd *ServiceDependenciesTemplate) UnmarshalYAML(b []byte) error {
	var asStrings []string
	if err := yaml.Unmarshal(b, &asStrings); err == nil {
		*sd = make([]ServiceDependencyTemplate, len(asStrings))
		for i, service := range asStrings {
			(*sd)[i] = ServiceDependencyTemplate{
				Service:   service,
				Condition: "service_started",
			}
		}
		return nil
	}

	asMap := make(map[string]struct {
		Condition string `yaml:"condition"`
	})
	if err := yaml.Unmarshal(b, &asMap); err != nil {
		return err
	}

	*sd = make([]ServiceDependencyTemplate, 0, len(asMap))
	for service, spec := range asMap {
		switch spec.Condition {
		case "service_started", "service_healthy", "service_completed_successfully":
			// Ok.
		case "":
			spec.Condition = "service_started"
		default:
			return fmt.Errorf("invalid condition %q for service dependency %q", spec.Condition, service)
		}
		*sd = append(*sd, ServiceDependencyTemplate{
			Service:   service,
			Condition: spec.Condition,
		})
	}

	return nil
}
