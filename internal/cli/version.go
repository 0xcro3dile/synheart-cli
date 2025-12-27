package cli

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "0.0.1"
	// Commit is set at build time
	Commit = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Displays the version of the Synheart CLI.`,
	RunE:  runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()

	if globalOpts.Format == "json" {
		payload := map[string]any{
			"name":    "synheart",
			"version": Version,
			"commit":  Commit,
			"go":      runtime.Version(),
			"os":      runtime.GOOS,
			"arch":    runtime.GOARCH,
		}
		if ui != nil {
			return ui.PrintJSON(payload)
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}

	if ui != nil {
		ui.Header("Synheart CLI")
		ui.KV("Version", "v"+Version)
		ui.KV("Commit", Commit)
		ui.KV("Go", runtime.Version())
		ui.KV("OS/Arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
		return nil
	}

	fmt.Fprintf(out, "Synheart CLI v%s\n", Version)
	fmt.Fprintf(out, "Commit: %s\n", Commit)
	fmt.Fprintf(out, "Go: %s\n", runtime.Version())
	fmt.Fprintf(out, "OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	return nil
}
