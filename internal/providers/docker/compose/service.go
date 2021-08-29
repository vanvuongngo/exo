package compose

import "github.com/goccy/go-yaml"

type ServiceTemplate struct {
	Name          string                      `yaml:"name"`
	ContainerName string                      `yaml:"containerName"`
	Networks      ServiceNetworksTemplate     `yaml:"networks"`
	Volumes       []VolumeMountTemplate       `yaml:"volumes"`
	DependsOn     ServiceDependenciesTemplate `yaml:"depends_on"`
	MapSlice      yaml.MapSlice               `yaml:",inline"`
}

type Service struct {
	Deploy IgnoredField `yaml:"deploy"`

	// Note that these two are only applicable to Windows.
	// TODO: cpu_count
	// TODO: cpu_percent

	CPUShares int64 `yaml:"cpu_shares"`
	CPUPeriod int64 `yaml:"cpu_period"`
	CPUQuota  int64 `yaml:"cpu_quota"`

	// CPURealtimeRuntime and CPURealtimePeriod can both be specified as either
	// strings or integers.
	//CPURealtimeRuntime int64 `yaml:"cpu_rt_runtime"`
	//CPURealtimePeriod  int64 `yaml:"cpu_rt_period"`

	// TODO: cpus
	// TODO: cpuset
	// TODO: blkio_config

	Build Build `yaml:"build"`
	// TODO: cap_add
	// TODO: cap_drop
	// TODO: cgroup_parent

	Command       Command  `yaml:"command"`
	Configs       []string `yaml:"configs"` // TODO: support long syntax.
	ContainerName string   `yaml:"container_name"`
	// TODO: credential_spec

	DependsOn map[string]ServiceDependency `yaml:"depends_on"`

	// TODO: device_cgroup_rules
	// TODO: devices
	// TODO: dns
	// TODO: dns_opt
	// TODO: dns_search
	Domainname string  `yaml:"domainname"`
	Entrypoint Command `yaml:"entrypoint"`
	// TODO: env_file
	Environment Dictionary   `yaml:"environment"`
	Expose      PortMappings `yaml:"expose"` // TODO: Validate target-only.
	// TODO: extends
	// TODO: external_links
	// TODO: extra_hosts
	// TODO: group_add
	Healthcheck *Healthcheck `yaml:"healthcheck"`
	Hostname    string       `yaml:"hostname"`
	Image       string       `yaml:"image"`
	// TODO: init
	// TODO: ipc
	// TODO: isolation
	Labels Dictionary `yaml:"labels"`
	// TODO: links
	Logging Logging `yaml:"logging"`
	// TODO: network_mode
	Networks   map[string]ServiceNetwork `yaml:"networks"`
	MacAddress string                    `yaml:"mac_address"`

	MemorySwappiness *int64 `yaml:"mem_swappiness"`

	// MemoryLimit and MemoryReservation can be specified either as strings or integers.
	MemoryLimit       MemoryField `yaml:"mem_limit"`
	MemoryReservation MemoryField `yaml:"mem_reservation"`

	// TODO: memswap_limit
	// TODO: oom_kill_disable
	// TODO: oom_score_adj
	// TODO: pid
	// TODO: pids_limit
	// TODO: platform

	Ports      PortMappings `yaml:"ports"`
	Privileged bool         `yaml:"privileged"`
	Profiles   IgnoredField `yaml:"profiles"`
	// TODO: pull_policy
	// TODO: read_only
	Restart string `yaml:"restart"`
	Runtime string `yaml:"runtime"`
	// TODO: scale
	Secrets []string `yaml:"secrets"` // TODO: support long syntax.
	// TODO: security_opt
	ShmSize         Bytes     `yaml:"shm_size"`
	StdinOpen       bool      `yaml:"stdin_open"`
	StopGracePeriod *Duration `yaml:"stop_grace_period"`
	StopSignal      string    `yaml:"stop_signal"`
	// TODO: storage_opt
	// TODO: sysctls
	// TODO: tmpfs
	TTY bool `yaml:"tty"`
	// TODO: ulimits
	User string `yaml:"user"`
	// TODO: userns_mode
	Volumes []VolumeMount `yaml:"volumes"`
	// TODO: volumes_from

	WorkingDir string `yaml:"working_dir"`
}

type Healthcheck struct {
	Test        Command  `yaml:"test"`
	Interval    Duration `yaml:"interval"`
	Timeout     Duration `yaml:"timeout"`
	Retries     int      `yaml:"retries"`
	StartPeriod Duration `yaml:"start_period"`
}

type Logging struct {
	Driver  string            `yaml:"driver"`
	Options map[string]string `yaml:"options"`
}

type ServiceNetworksTemplate map[string]ServiceNetworkTemplate

func (t *ServiceNetworksTemplate) UnmarshalYAML(b []byte) error {
	var long map[string]ServiceNetworkTemplate
	var short []string
	if err := yaml.Unmarshal(b, &short); err == nil {
		long = make(map[string]ServiceNetworkTemplate, len(short))
		for _, s := range short {
			long[s] = ServiceNetworkTemplate{}
		}
	} else if err := yaml.Unmarshal(b, &long); err != nil {
		return err
	}
	return nil
}

type ServiceNetworkTemplate struct {
	Aliases      []string `yaml:"aliases"`
	IPv4Address  string   `yaml:"ipv4_address"`
	IPv6Address  string   `yaml:"ipv6_address"`
	LinkLocalIPs []string `yaml:"link_local_ips"`
	Priority     string   `yaml:"priority"`
}

type ServiceNetwork struct {
	Aliases      []string `yaml:"aliases"`
	IPv4Address  string   `yaml:"ipv4_address"`
	IPv6Address  string   `yaml:"ipv6_address"`
	LinkLocalIPs []string `yaml:"link_local_ips"`
	// TODO: Should be an integer.  If ServiceNetworkTemplate.Priority is empty
	// string, what is the default priority?
	Priority string `yaml:"priority"`
}
