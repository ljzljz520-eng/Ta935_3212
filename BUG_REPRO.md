# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	engineering-document-vault/cmd/docvault	[no test files]
ok  	engineering-document-vault/internal/catalog	0.001s
ok  	engineering-document-vault/internal/domain	0.001s
ok  	engineering-document-vault/internal/policy	0.001s
ok  	engineering-document-vault/internal/search	0.001s
?   	engineering-document-vault/internal/service	[no test files]
ok  	engineering-document-vault/internal/store	0.008s
--- FAIL: TestWorkflow13BusinessInvariant (0.00s)
    workflow_test.go:30: expected stable rejection on seventh review
FAIL
FAIL	engineering-document-vault/internal/workflow13	0.011s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/docvault): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/docvault): exit `0`
