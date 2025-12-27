package cli

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/synheart/synheart-cli/internal/scenario"
)

var listScenariosCmd = &cobra.Command{
	Use:   "list-scenarios",
	Short: "List available scenarios",
	Long:  `Lists all built-in scenarios with their descriptions.`,
	Aliases: []string{
		"scenarios",
		"ls",
	},
	RunE: runListScenarios,
}

func runListScenarios(cmd *cobra.Command, args []string) error {
	// Load scenarios
	registry := scenario.NewRegistry()
	if err := registry.LoadFromDir(getScenarioDir()); err != nil {
		return fmt.Errorf("failed to load scenarios: %w", err)
	}

	scenarios := registry.ListWithDescriptions()
	if len(scenarios) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No scenarios found")
		return nil
	}

	if globalOpts.Format == "json" {
		type row struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		names := make([]string, 0, len(scenarios))
		for name := range scenarios {
			names = append(names, name)
		}
		sort.Strings(names)
		out := make([]row, 0, len(names))
		for _, name := range names {
			out = append(out, row{Name: name, Description: scenarios[name]})
		}
		// UI may be nil if called in tests; fall back to stdout encoder.
		if ui != nil {
			return ui.PrintJSON(out)
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	// Sort by name
	names := make([]string, 0, len(scenarios))
	for name := range scenarios {
		names = append(names, name)
	}
	sort.Strings(names)

	if ui != nil {
		ui.Header("Available scenarios")
		ui.Println()
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "Available scenarios:")
		fmt.Fprintln(cmd.OutOrStdout())
	}
	for _, name := range names {
		fmt.Fprintf(cmd.OutOrStdout(), "  %-20s %s\n", name, scenarios[name])
	}
	fmt.Fprintln(cmd.OutOrStdout())

	return nil
}
