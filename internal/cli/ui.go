package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type UI struct {
	out io.Writer
	err io.Writer

	colorEnabled bool
	quiet        bool
	verbose      bool
}

func NewUI(out io.Writer, err io.Writer, noColor bool, quiet bool, verbose bool) *UI {
	colorEnabled := isTTY(os.Stdout) && !noColor && os.Getenv("TERM") != "dumb"
	return &UI{
		out:          out,
		err:          err,
		colorEnabled: colorEnabled,
		quiet:        quiet,
		verbose:      verbose,
	}
}

func isTTY(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func (u *UI) Header(title string) {
	if u.quiet {
		return
	}
	fmt.Fprintln(u.out, u.bold(title))
}

func (u *UI) Section(title string) {
	if u.quiet {
		return
	}
	fmt.Fprintln(u.out, u.cyan(title))
}

func (u *UI) KV(key string, value any) {
	if u.quiet {
		return
	}
	fmt.Fprintf(u.out, "%-14s %v\n", key+":", value)
}

func (u *UI) Printf(format string, args ...any) {
	if u.quiet {
		return
	}
	fmt.Fprintf(u.out, format, args...)
}

func (u *UI) Println(args ...any) {
	if u.quiet {
		return
	}
	fmt.Fprintln(u.out, args...)
}

func (u *UI) Debugf(format string, args ...any) {
	if !u.verbose {
		return
	}
	fmt.Fprintf(u.err, u.dim("debug: ")+format+"\n", args...)
}

func (u *UI) Warnf(format string, args ...any) {
	if u.quiet {
		return
	}
	fmt.Fprintf(u.err, u.yellow("warn: ")+format+"\n", args...)
}

func (u *UI) Errorf(format string, args ...any) {
	fmt.Fprintf(u.err, u.red("error: ")+format+"\n", args...)
}

func (u *UI) Successf(format string, args ...any) {
	if u.quiet {
		return
	}
	fmt.Fprintf(u.out, u.green("ok: ")+format+"\n", args...)
}

func (u *UI) PrintJSON(v any) error {
	enc := json.NewEncoder(u.out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (u *UI) bold(s string) string   { return u.wrap("1", s) }
func (u *UI) dim(s string) string    { return u.wrap("2", s) }
func (u *UI) red(s string) string    { return u.wrap("31", s) }
func (u *UI) green(s string) string  { return u.wrap("32", s) }
func (u *UI) yellow(s string) string { return u.wrap("33", s) }
func (u *UI) cyan(s string) string   { return u.wrap("36", s) }

func (u *UI) wrap(code string, s string) string {
	if !u.colorEnabled {
		return s
	}
	// Avoid coloring multiline blocks (keeps copy/paste nice).
	if strings.ContainsRune(s, '\n') {
		return s
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", code, s)
}
