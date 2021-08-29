package compose_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/stretchr/testify/assert"
)

func TestParseServiceTemplate(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected compose.ServiceTemplate
	}{
		{
			name: "volumes - full syntax",
			in: `volumes:
- type: volume
  source: mydata
  target: /data
  read_only: true
  volume:
    nocopy: true
- type: bind
  source: /path/a
  target: /path/b
  bind:
    propagation: rshared
    create_host_path: true
- type: tmpfs
  target: /data/buffer
  tmpfs:
    size: 208666624`,
			expected: compose.ServiceTemplate{
				Volumes: []compose.VolumeMountTemplate{
					{
						Type:     "volume",
						Source:   "mydata",
						Target:   "/data",
						ReadOnly: "true",
						Volume: &compose.VolumeOptionsTemplate{
							Nocopy: "true",
						},
					},
					{
						Type:   "bind",
						Source: "/path/a",
						Target: "/path/b",
						Bind: &compose.BindOptionsTemplate{
							Propagation:    "rshared",
							CreateHostPath: "true",
						},
					},
					{
						Type:   "tmpfs",
						Target: "/data/buffer",
						Tmpfs: &compose.TmpfsOptionsTemplate{
							Size: "208666624",
						},
					},
				},
			},
		},

		{
			name: "volumes: short syntax",
			in: `volumes:
- /var/myapp
- './data:/data'
- "/home/fred/.ssh:/root/.ssh:ro"
- '~/util:/usr/bin/util:rw'
- my-log-volume:/var/log/xyzzy`,
			expected: compose.ServiceTemplate{
				Volumes: []compose.VolumeMountTemplate{
					{
						Type:   "volume",
						Target: "/var/myapp",
					},
					{
						Type:   "bind",
						Source: "./data",
						Target: "/data",
						Bind: &compose.BindOptionsTemplate{
							CreateHostPath: "true",
						},
					},
					{
						Type:     "bind",
						Source:   "/home/fred/.ssh",
						Target:   "/root/.ssh",
						ReadOnly: "true",
						Bind: &compose.BindOptionsTemplate{
							CreateHostPath: "true",
						},
					},
					{
						Type:   "bind",
						Source: "~/util",
						Target: "/usr/bin/util",
						Bind: &compose.BindOptionsTemplate{
							CreateHostPath: "true",
						},
					},
					{
						Type:   "volume",
						Source: "my-log-volume",
						Target: "/var/log/xyzzy",
					},
				},
			},
		},

		{
			name: "service dependencies - short syntax",
			in: `depends_on:
- db
- messages`,
			expected: compose.ServiceTemplate{
				DependsOn: []compose.ServiceDependencyTemplate{
					{
						Service:   "db",
						Condition: "service_started",
					},
					{
						Service:   "messages",
						Condition: "service_started",
					},
				},
			},
		},

		{
			name: "service dependencies - extended syntax",
			in: `depends_on:
  db:
  messages:
    condition: service_healthy`,
			expected: compose.ServiceTemplate{
				DependsOn: []compose.ServiceDependencyTemplate{
					{
						Service:   "db",
						Condition: "service_started",
					},
					{
						Service:   "messages",
						Condition: "service_healthy",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		name := testCase.name
		in := testCase.in
		expected := testCase.expected
		t.Run(name, func(t *testing.T) {
			var content bytes.Buffer
			content.WriteString("services:\n  test-svc:\n")
			lines := strings.Split(in, "\n")
			for _, line := range lines {
				content.WriteString("    ")
				content.WriteString(line)
				content.WriteByte('\n')
			}

			var proj compose.ProjectTemplate
			yamlutil.MustUnmarshalString(content.String(), &proj)

			// Ignore MapSlice fields in test.
			proj.MapSlice = nil
			for k, service := range proj.Services {
				service.MapSlice = nil
				proj.Services[k] = service
			}
			for k, volume := range proj.Volumes {
				volume.MapSlice = nil
				proj.Volumes[k] = volume
			}
			for k, network := range proj.Networks {
				network.MapSlice = nil
				proj.Networks[k] = network
			}

			svc := proj.Services["test-svc"]
			assert.Equal(t, expected, svc)
		})
	}
}
