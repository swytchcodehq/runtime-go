# swytchcode Go runtime

Thin runtime wrapper around the Swytchcode CLI. Calls `swytchcode exec` for you so you can stay in Go without shell boilerplate.

**Requires:** The `swytchcode` binary must be installed and on your `PATH`.

## Install

```bash
go get gitlab.com/swytchcode/go-runtime
```

## Use

### JSON mode (default)

```go
package main

import (
	"context"
	"fmt"
	swytchcode "gitlab.com/swytchcode/go-runtime"
)

func main() {
	result, err := swytchcode.Exec(context.Background(), "api.account.create", map[string]string{
		"email": "test@example.com",
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
```

Equivalent to: `swytchcode exec api.account.create --json` with args on stdin.

**Request input (args):** The third argument to `Exec` is the kernel **args** object (sent as JSON on stdin). Use this shape so the kernel builds the request correctly:
- **body** — Request body (object).
- **params** — Query/path params (e.g. `map[string]string{"id": "cluster-123"}`).
- **Authorization** — Auth header value (e.g. `"Bearer token123"`).
- **headers** — Additional request headers (map of header name to value).
- Other top-level keys are passed as query params.

Example with body, params, and headers:

```go
swytchcode.Exec(ctx, "api.cluster.get", map[string]any{
	"params":         map[string]string{"id": "cluster-123"},
	"Authorization": "Bearer token123",
	"headers":       map[string]string{"X-Request-Id": "abc-123"},
}, nil)
```

### Raw mode

Get stdout as bytes instead of parsing JSON:

```go
	result, err := swytchcode.Exec(ctx, "api.report.export", map[string]string{"id": "123"}, &swytchcode.ExecOptions{Raw: true})
	// result is []byte
```

### Options

- **Cwd** – Working directory for the process (default: current directory).
- **Env** – Extra environment variables (merged with the current process env).
- **Raw** – If true, use `--raw` and return stdout as `[]byte`.
- **DryRun** – If true, pass `--dry-run` to the CLI; request details (method, url, headers, body) are output instead of calling the server.
- **AllowRaw** – If true, pass `--allow-raw` to the CLI; required for executing raw methods (kernel has this disabled by default).

This runtime invokes `swytchcode exec [canonical_id]` with the flags above. For full exec behavior, see the Swytchcode kernel documentation.

### Error handling

Failures return a `*SwytchcodeError` (spawn error, non-zero exit, invalid JSON). Use the helper to narrow:

```go
	result, err := swytchcode.Exec(ctx, "api.account.create", input, nil)
	if err != nil {
		if swytchcode.IsSwytchcodeError(err) {
			fmt.Println(err.(*swytchcode.SwytchcodeError).Message, err.(*swytchcode.SwytchcodeError).Cause)
		}
		return err
	}
```

## What this library is

- A thin process wrapper over the Swytchcode shell.
- A stable call surface: `Exec(ctx, canonicalID, input, opts)`.
- A convenience for developers and IDE agents.

## What this library is not

- An HTTP SDK or API client.
- A replacement for the shell. All authority stays in the Swytchcode CLI.

This library **must not** and **does not**:

- Read `.swytchcode/*` or parse `tooling.json`.
- Validate schemas or add retries/timeouts.
- Call the registry or embed auth logic.

It is a pipe, not a brain.

## For Cursor / Claude / IDE agents

Generate code that uses this runtime instead of subprocess boilerplate.

**Preferred:**

```go
import swytchcode "gitlab.com/swytchcode/go-runtime"

func createAccount(ctx context.Context, input map[string]string) (any, error) {
	return swytchcode.Exec(ctx, "api.account.create", input, nil)
}
```

**Avoid:** os/exec boilerplate, HTTP calls, or config parsing. Use `Exec(...)` and let the CLI handle execution and policy.
