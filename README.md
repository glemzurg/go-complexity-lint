# go-complexity-lint

A complexity linter for Go that measures four metrics with a three-zone severity model.<sup><a href="#cite1">1</a></sup> Yellow zone (warning) prints diagnostics but exits 0. Red zone (error) prints diagnostics and exits 1.

## Metrics

| Metric | What It Measures | Green | Yellow (warn) | Red (fail) |
|--------|-----------------|-------|---------------|------------|
| **nestdepth** | Max depth of control-flow nesting in a function | 1–4 | 5–6 | 7+ |
| **cyclo** | Cyclomatic complexity: 1 + 1 per branching/looping decision | 1–9 | 10–14 | 15+ |
| **params** | Number of function parameters | 0–4 | 5–6 | 7+ |
| **fanout** | Distinct non-builtin, non-stdlib function calls | 0–6 | 7–9 | 10+ |

### Counting Rules

**Cyclomatic complexity** counts: `if`, `for`, `range`, non-default `case`, non-default `select case`. Does not count: `else`, `default`, `&&`/`||`, `switch`/`select` themselves. Each `else if` counts as a new decision.

**Nesting depth** counts: `if`/`else`/`else if`, `for`, `range`, `switch`, `select`, `type switch`, func literals. Each level adds 1 to depth.

**Fan out** counts distinct function/method calls resolved via type information. Excludes builtins (`len`, `make`, etc.), type conversions, and standard library functions.

**Error guard clause exemption**: Both `nestdepth` and `cyclo` exempt the idiomatic Go error-handling pattern `if <ident> != nil { return ..., <ident> }` where the body is a single return statement with zero-valued results except the final error. The error variable can have any name (`err`, `e`, `dbErr`, etc.).

## Installation

```sh
go install github.com/glemzurg/go-complexity-lint/cmd/go-complexity-lint@latest
```

## Usage

### Standalone

```sh
go-complexity-lint ./...

# With custom thresholds (flags are namespaced by analyzer)
go-complexity-lint -nestdepth.warn=3 -nestdepth.fail=5 -cyclo.warn=12 -cyclo.fail=20 ./...
go-complexity-lint -params.warn=5 -params.fail=8 -fanout.warn=8 -fanout.fail=12 ./...
```

Thresholds must be non-negative and `warn` must not exceed `fail`.

### With `go vet`

```sh
go vet -vettool=$(which go-complexity-lint) ./...
```

Note: `go vet` treats all diagnostics as failures (exit 1) regardless of zone. It does not distinguish between warnings and errors. To suppress warnings and only fail on red-zone violations, set `warn` equal to `fail`:

```sh
go vet -vettool=$(which go-complexity-lint) \
  -nestdepth.warn=6 -cyclo.warn=14 -params.warn=6 -fanout.warn=9 ./...
```

For full severity-aware exit codes, use the standalone binary.

### Per-Function Overrides

Use doc comments to override thresholds for specific functions:

```go
//complexity:cyclo:warn=20,fail=30
//complexity:nestdepth:warn=8,fail=10
func ComplexRouter(input string) error {
    // ...
}
```

## golangci-lint Integration

### Plugin Configuration (`.custom-gcl.yml`)

```yaml
version: v2
plugins:
  - module: "github.com/glemzurg/go-complexity-lint"
    import: "github.com/glemzurg/go-complexity-lint/plugin"
    version: latest
```

### Linter Settings (`.golangci.yml`)

```yaml
linters:
  enable:
    - go-complexity-lint

linters-settings:
  custom:
    go-complexity-lint:
      type: module
      settings:
        nestdepth-warn: 3
        nestdepth-fail: 5
        cyclo-warn: 12
        cyclo-fail: 20
        params-warn: 5
        params-fail: 8
        fanout-warn: 8
        fanout-fail: 12
```

All settings are optional. Omitted values use the defaults from the metrics table above.

Build the custom binary:

```sh
golangci-lint custom
```

Then run:

```sh
./custom-gcl run ./...
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | No red-zone violations (warnings may be present) |
| 1 | Red-zone violations found or analysis error |

## References

<span id="cite1"><sup>1</sup></span> Tockey, S. (2019). *How to Engineer Software: A Model-Based Approach*. Wiley-IEEE Computer Society Press. Appendix N: Software Structural Complexity Metrics.