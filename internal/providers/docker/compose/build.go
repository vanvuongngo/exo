package compose

type BuildTemplate ExpandedBuildTemplate

func (b BuildTemplate) MarshalYAML() (interface{}, error) {
	return ExpandedBuildTemplate(b), nil
}

func (dict *BuildTemplate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	var expanded ExpandedBuildTemplate
	err := unmarshal(&s)
	if err == nil {
		expanded.Context = s
	} else if err := unmarshal(&expanded); err != nil {
		return nil
	}
	*dict = BuildTemplate(expanded)
	return nil
}

type ExpandedBuildTemplate struct {
	Context    string     `yaml:"context"`
	Dockerfile string     `yaml:"dockerfile"`
	Args       Dictionary `yaml:"args"`
	CacheFrom  []string   `yaml:"cache_from"`
	ExtraHosts []string   `yaml:"extra_hosts"`
	Isolation  string     `yaml:"isolation"`
	Labels     Dictionary `yaml:"labels"`
	ShmSize    string     `yaml:"shm_size"`
	Target     string     `yaml:"target"`
}

type Build struct {
	Context    string     `yaml:"context"`
	Dockerfile string     `yaml:"dockerfile"`
	Args       Dictionary `yaml:"args"`
	CacheFrom  []string   `yaml:"cache_from"`
	ExtraHosts []string   `yaml:"extra_hosts"`
	Isolation  string     `yaml:"isolation"`
	Labels     Dictionary `yaml:"labels"`
	ShmSize    Bytes      `yaml:"shm_size"`
	Target     string     `yaml:"target"`
}
