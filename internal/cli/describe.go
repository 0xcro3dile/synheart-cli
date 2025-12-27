package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/synheart/synheart-cli/internal/scenario"
)

var describeCmd = &cobra.Command{
	Use:   "describe <scenario>",
	Short: "Describe a scenario in detail",
	Long:  `Shows detailed information about a scenario including signals, phases, and typical ranges.`,
	Aliases: []string{
		"scenario",
		"scen",
	},
	Args: cobra.ExactArgs(1),
	RunE: runDescribe,
}

func runDescribe(cmd *cobra.Command, args []string) error {
	scenarioName := args[0]

	// Load scenarios
	registry := scenario.NewRegistry()
	if err := registry.LoadFromDir(getScenarioDir()); err != nil {
		return fmt.Errorf("failed to load scenarios: %w", err)
	}

	// Get scenario
	scen, err := registry.Get(scenarioName)
	if err != nil {
		return fmt.Errorf("scenario not found: %w", err)
	}

	if globalOpts.Format == "json" {
		type outScenario struct {
			Name        string                            `json:"name"`
			Description string                            `json:"description"`
			Duration    string                            `json:"duration"`
			DefaultRate string                            `json:"default_rate"`
			Signals     map[string]*scenario.SignalConfig `json:"signals"`
			Phases      []scenario.Phase                  `json:"phases"`
		}
		payload := outScenario{
			Name:        scen.Name,
			Description: scen.Description,
			Duration:    scen.Duration,
			DefaultRate: scen.DefaultRate,
			Signals:     scen.Signals,
			Phases:      scen.Phases,
		}
		if ui != nil {
			return ui.PrintJSON(payload)
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}

	// Print details
	if ui != nil {
		ui.Header("Scenario")
		ui.KV("Name", scen.Name)
		ui.KV("Description", scen.Description)
		ui.KV("Duration", scen.Duration)
		ui.KV("Default rate", scen.DefaultRate)
		ui.Println()
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Scenario: %s\n", scen.Name)
		fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", scen.Description)
		fmt.Fprintf(cmd.OutOrStdout(), "Duration: %s\n", scen.Duration)
		fmt.Fprintf(cmd.OutOrStdout(), "Default Rate: %s\n\n", scen.DefaultRate)
	}

	// Print signals
	if ui != nil {
		ui.Section("Signals")
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "Signals:")
	}
	for name, config := range scen.Signals {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", name)
		if config.Baseline != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "    Baseline: %v\n", config.Baseline)
		}
		if config.Noise != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "    Noise: %v\n", config.Noise)
		}
		if config.Rate != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "    Rate: %s\n", config.Rate)
		}
		if config.Unit != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "    Unit: %s\n", config.Unit)
		}
	}

	// Print phases
	if len(scen.Phases) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "\nPhases:")
		for i, phase := range scen.Phases {
			fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s (duration: %s)\n", i+1, phase.Name, phase.Duration)
			if len(phase.Overrides) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "     Overrides:")
				for signal, override := range phase.Overrides {
					fmt.Fprintf(cmd.OutOrStdout(), "       %s:", signal)
					if override.Add != 0 {
						fmt.Fprintf(cmd.OutOrStdout(), " add=%.1f", override.Add)
					}
					if override.Multiply != 0 {
						fmt.Fprintf(cmd.OutOrStdout(), " multiply=%.1f", override.Multiply)
					}
					if override.Value != "" {
						fmt.Fprintf(cmd.OutOrStdout(), " value=%s", override.Value)
					}
					if override.Baseline != nil {
						fmt.Fprintf(cmd.OutOrStdout(), " baseline=%v", override.Baseline)
					}
					if override.Noise != nil {
						fmt.Fprintf(cmd.OutOrStdout(), " noise=%v", override.Noise)
					}
					fmt.Fprintln(cmd.OutOrStdout())
				}
			}
		}
	}

	fmt.Fprintln(cmd.OutOrStdout())
	return nil
}
