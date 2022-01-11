package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stackCmd)

	stackCmd.AddCommand(makeHelpSubcmd())
}

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Create, inspect and modify stacks",
	Long: `Contains subcommands for operating on stacks.

If no subcommand is given, describes the current stack of the current
workspace.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		joshClient := newClient()
		workspace := requireCurrentWorkspace(ctx, joshClient)
		workspaceID := workspace.ID()

		gqlClient, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Workspace struct {
				Stack *struct {
					ID   string
					Name string
				}
			} `graphql:"workspaceById(id: $id)"`
		}
		err := gqlClient.Query(ctx, &q, map[string]interface{}{
			"id": graphql.String(workspaceID),
		})
		if err != nil {
			return fmt.Errorf("querying: %w", err)
		}
		stack := q.Workspace.Stack
		if stack == nil {
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", stack.ID)
		_ = w.Flush()
		return nil
	},
}
