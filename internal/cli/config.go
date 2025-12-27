package cli

// GlobalOptions are shared flags that apply across commands.
type GlobalOptions struct {
	Format  string
	NoColor bool
	Quiet   bool
	Verbose bool
}

var globalOpts = GlobalOptions{
	Format: "text",
}
