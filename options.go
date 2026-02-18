package swytchcode

// ExecOptions configures an exec call.
type ExecOptions struct {
	// Cwd is the working directory for the swytchcode process. If empty, the current process directory is used.
	Cwd string
	// Env is merged with the current process environment. Keys and values are added or override.
	Env map[string]string
	// Raw, if true, passes --raw to the CLI and returns stdout as []byte instead of parsing JSON.
	Raw bool
	// DryRun, if true, passes --dry-run to the CLI; request details are output instead of calling the server.
	DryRun bool
	// AllowRaw, if true, passes --allow-raw to the CLI; required for executing raw methods (disabled by default in kernel).
	AllowRaw bool
}
