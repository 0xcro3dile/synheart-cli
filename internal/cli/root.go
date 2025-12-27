package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "synheart",
	Short: "Synheart CLI - Mock HSI data generator for local development",
	Long: `Synheart CLI generates HSI-compatible sensor data streams
that mimic phone + wearable sources for local SDK development.

It eliminates dependency on physical devices during development,
providing repeatable scenarios for QA and demos.`,
	Example: strings.TrimSpace(`
synheart mock start
synheart mock list-scenarios
synheart doctor
synheart version
`),
	SilenceUsage:  true,
	SilenceErrors: true,
}

var ui *UI

// Execute runs the root command
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	rootCmd.SetHelpTemplate(helpTemplate())
	rootCmd.SetUsageTemplate(usageTemplate())

	if err := rootCmd.Execute(); err != nil {
		// At this point flags are parsed; UI is configured in init().
		if ui == nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		ui.Errorf("%v", err)
		ui.Printf("hint: run %s\n", ui.dim("synheart --help"))
		os.Exit(1)
	}
}

func init() {
	initRootFlags()
	rootCmd.AddCommand(mockCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
}

func initRootFlags() {
	rootCmd.PersistentFlags().StringVar(&globalOpts.Format, "format", "text", "Output format: text|json")
	rootCmd.PersistentFlags().BoolVar(&globalOpts.NoColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVarP(&globalOpts.Quiet, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&globalOpts.Verbose, "verbose", "v", false, "Verbose logging")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Cobra defaults to stdout/stderr; use the command writers so tests/redirects work.
		out := cmd.OutOrStdout()
		er := cmd.ErrOrStderr()
		ui = NewUI(out, er, globalOpts.NoColor, globalOpts.Quiet, globalOpts.Verbose)

		// Normalize format.
		globalOpts.Format = strings.ToLower(strings.TrimSpace(globalOpts.Format))
		if globalOpts.Format == "" {
			globalOpts.Format = "text"
		}
		if globalOpts.Format != "text" && globalOpts.Format != "json" {
			return fmt.Errorf("invalid --format %q (expected: text|json)", globalOpts.Format)
		}

		// Make sure help/usage can write to the right place.
		cmd.SetOut(out)
		cmd.SetErr(er)
		return nil
	}
}

func helpTemplate() string {
	// Minimal, modern-ish help output without extra dependencies.
	// Cobra will expand the template with command-specific values.
	return strings.TrimSpace(`
{{with (or .Long .Short)}}{{.}}{{end}}

Usage:
  {{.UseLine}}

{{if .HasAvailableSubCommands}}Commands:
{{range .Commands}}{{if (and (not .Hidden) .IsAvailableCommand)}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{end}}

{{if .HasAvailableLocalFlags}}Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasAvailableInheritedFlags}}Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasExample}}Examples:
{{.Example}}{{end}}
`)
}

func usageTemplate() string {
	return strings.TrimSpace(`
Usage:
  {{.UseLine}}

Run 'synheart --help' for more information.
`)
}
