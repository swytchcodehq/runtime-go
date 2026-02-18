package swytchcode

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

// Exec runs swytchcode exec <canonicalID> with optional JSON args on stdin.
// It returns parsed JSON as an interface{} (default) or raw stdout as []byte when opts.Raw is true.
// The swytchcode binary must be on PATH.
//
// The input argument is the kernel args object sent on stdin. Use this shape so the kernel
// builds the request correctly: body (request body), params (query/path params),
// Authorization (auth header value), headers (map of header name to value); other
// top-level keys are passed as query params. See the Swytchcode kernel documentation.
func Exec(ctx context.Context, canonicalID string, input any, opts *ExecOptions) (any, error) {
	canonicalID = strings.TrimSpace(canonicalID)
	if canonicalID == "" {
		return nil, &SwytchcodeError{Message: "canonicalID must be a non-empty string", Cause: nil}
	}

	if opts == nil {
		opts = &ExecOptions{}
	}

	flag := "--json"
	if opts.Raw {
		flag = "--raw"
	}
	args := []string{"exec", canonicalID, flag}
	if opts.DryRun {
		args = append(args, "--dry-run")
	}
	if opts.AllowRaw {
		args = append(args, "--allow-raw")
	}

	cmd := exec.CommandContext(ctx, "swytchcode", args...)
	cmd.Dir = opts.Cwd
	if cmd.Dir == "" {
		cmd.Dir, _ = os.Getwd()
	}
	cmd.Env = mergeEnv(opts.Env)

	var stdinBuf bytes.Buffer
	if input != nil {
		if err := json.NewEncoder(&stdinBuf).Encode(input); err != nil {
			return nil, &SwytchcodeError{Message: "failed to encode input", Cause: err}
		}
		cmd.Stdin = &stdinBuf
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = "swytchcode exec failed"
		}
		return nil, &SwytchcodeError{Message: strings.TrimSpace(msg), Cause: err}
	}

	if opts.Raw {
		return stdout.Bytes(), nil
	}

	if stdout.Len() == 0 {
		return nil, nil
	}

	var result any
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, &SwytchcodeError{Message: "invalid JSON output from swytchcode", Cause: stdout.String()}
	}
	return result, nil
}

func mergeEnv(extra map[string]string) []string {
	base := os.Environ()
	if len(extra) == 0 {
		return base
	}
	overrides := make(map[string]string)
	for _, e := range base {
		if i := strings.IndexByte(e, '='); i >= 0 {
			overrides[e[:i]] = e[i+1:]
		}
	}
	for k, v := range extra {
		overrides[k] = v
	}
	out := make([]string, 0, len(overrides))
	for k, v := range overrides {
		out = append(out, k+"="+v)
	}
	return out
}
