package cli

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/synheart/synheart-cli/internal/scenario"
)

var (
	doctorHost string
	doctorPort int
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment and print connection info",
	Long:  `Validates the local environment, checks port availability, and provides connection examples.`,
	RunE:  runDoctor,
}

func init() {
	doctorCmd.Flags().StringVar(&doctorHost, "host", "127.0.0.1", "Host to bind to")
	doctorCmd.Flags().IntVar(&doctorPort, "port", 8787, "Port to check")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()

	type doctorJSON struct {
		GoVersion    string   `json:"go_version"`
		OS           string   `json:"os"`
		Arch         string   `json:"arch"`
		ScenariosDir string   `json:"scenarios_dir"`
		Scenarios    []string `json:"scenarios"`
		Port         int      `json:"port"`
		Host         string   `json:"host"`
		PortFree     bool     `json:"port_free"`
		WebSocketURL string   `json:"websocket_url"`
	}

	if globalOpts.Format == "text" {
		if ui != nil {
			ui.Header("Synheart doctor")
			ui.Println()
		} else {
			fmt.Fprintln(out, "Synheart doctor")
			fmt.Fprintln(out)
		}
	}

	// Check Go version
	if globalOpts.Format == "text" {
		if ui != nil {
			ui.KV("Go", runtime.Version())
			ui.KV("OS/Arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
			ui.Println()
		} else {
			fmt.Fprintf(out, "Go:      %s\n", runtime.Version())
			fmt.Fprintf(out, "OS/Arch: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
		}
	}

	// Check scenarios directory
	scenariosDir := getScenarioDir()
	scenarios := []string(nil)
	if _, err := os.Stat(scenariosDir); err == nil {
		// Count scenarios
		registry := scenario.NewRegistry()
		if err := registry.LoadFromDir(scenariosDir); err == nil {
			scenarios = registry.List()
		}
		if globalOpts.Format == "text" {
			if ui != nil {
				ui.Successf("scenarios directory found: %s", scenariosDir)
				if len(scenarios) > 0 {
					ui.KV("Scenarios", scenarios)
				}
				ui.Println()
			} else {
				fmt.Fprintf(out, "Scenarios directory: %s\n", scenariosDir)
				if len(scenarios) > 0 {
					fmt.Fprintf(out, "Scenarios: %v\n", scenarios)
				}
				fmt.Fprintln(out)
			}
		}
	} else {
		if globalOpts.Format == "text" {
			if ui != nil {
				ui.Errorf("scenarios directory not found: %s", scenariosDir)
				ui.Println()
			} else {
				fmt.Fprintf(out, "Scenarios directory missing: %s\n\n", scenariosDir)
			}
		}
	}

	// Check default port availability
	portFree := isPortAvailable(doctorHost, doctorPort)
	wsURL := fmt.Sprintf("ws://%s:%d/hsi", doctorHost, doctorPort)
	if globalOpts.Format == "text" {
		if portFree {
			if ui != nil {
				ui.Successf("port %d is available on %s", doctorPort, doctorHost)
				ui.KV("WebSocket", wsURL)
				ui.Println()
			} else {
				fmt.Fprintf(out, "Port %d is available on %s\n", doctorPort, doctorHost)
				fmt.Fprintf(out, "WebSocket: %s\n\n", wsURL)
			}
		} else {
			if ui != nil {
				ui.Warnf("port %d is in use on %s", doctorPort, doctorHost)
				ui.KV("WebSocket", wsURL)
				ui.Println()
			} else {
				fmt.Fprintf(out, "Port %d is in use on %s\n", doctorPort, doctorHost)
				fmt.Fprintf(out, "WebSocket: %s\n\n", wsURL)
			}
		}

		if ui != nil {
			ui.Section("Connection examples")
		} else {
			fmt.Fprintln(out, "Connection examples:")
		}
		fmt.Fprintln(out)

		fmt.Fprintln(out, "JavaScript/Node.js:")
		fmt.Fprintf(out, "  const ws = new WebSocket('%s');\n", wsURL)
		fmt.Fprintln(out, "  ws.onmessage = (event) => {")
		fmt.Fprintln(out, "    const data = JSON.parse(event.data);")
		fmt.Fprintln(out, "    console.log(data);")
		fmt.Fprintln(out, "  };")
		fmt.Fprintln(out)

		fmt.Fprintln(out, "Python:")
		fmt.Fprintln(out, "  import websocket")
		fmt.Fprintln(out, "  import json")
		fmt.Fprintln(out, "  ws = websocket.WebSocket()")
		fmt.Fprintf(out, "  ws.connect('%s')\n", wsURL)
		fmt.Fprintln(out, "  while True:")
		fmt.Fprintln(out, "    data = json.loads(ws.recv())")
		fmt.Fprintln(out, "    print(data)")
		fmt.Fprintln(out)

		fmt.Fprintln(out, "Go:")
		fmt.Fprintf(out, "  conn, _, err := websocket.DefaultDialer.Dial(%q, nil)\n", wsURL)
		fmt.Fprintln(out, "  for {")
		fmt.Fprintln(out, "    _, message, err := conn.ReadMessage()")
		fmt.Fprintln(out, "    // decode message into your Event type")
		fmt.Fprintln(out, "    _ = message")
		fmt.Fprintln(out, "    _ = err")
		fmt.Fprintln(out, "  }")
		fmt.Fprintln(out)

		return nil
	}

	payload := doctorJSON{
		GoVersion:    runtime.Version(),
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		ScenariosDir: scenariosDir,
		Scenarios:    scenarios,
		Host:         doctorHost,
		Port:         doctorPort,
		PortFree:     portFree,
		WebSocketURL: wsURL,
	}
	if ui != nil {
		return ui.PrintJSON(payload)
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func isPortAvailable(host string, port int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
