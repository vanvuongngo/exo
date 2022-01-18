package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceRmCmd)
}

var resourceRmCmd = &cobra.Command{
	Use:   "rm <iri>", // TODO: Variadic.
	Short: "Remove a resource",
	Long:  "Dispose a resource via its provider, then stop tracking it.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		iri := args[0]

		var m struct {
			Job struct {
				ID string
			} `graphql:"disposeResource(iri: $iri)"`
		}
		if err := api.Mutate(ctx, svc, &m, map[string]interface{}{
			"iri": iri,
		}); err != nil {
			return err
		}
		return watchJob(ctx, m.Job.ID)
	},
}
