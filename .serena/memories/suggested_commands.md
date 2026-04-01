# Suggested Commands

## Build & Run
```bash
go build -o phunter .       # Build the binary
./phunter                    # Run the TUI app
```

## Development
```bash
go mod tidy                  # Clean up dependencies
go vet ./...                 # Static analysis / linting
go fmt ./...                 # Format all Go files
```

## Testing
```bash
go test ./...                # Run all tests (no tests exist yet)
```

## Release
```bash
git tag vX.Y.Z               # Tag a release
git push origin vX.Y.Z       # Push tag — triggers GitHub Actions release workflow
goreleaser release --clean    # Manual local release (requires goreleaser v2)
```

## System Utilities (macOS / Darwin)
```bash
git status / git log / git diff   # Version control
lsof -i -n -P -sTCP:LISTEN       # List listening TCP ports (what phunter uses internally)
kill -9 <pid>                     # Force kill a process (what phunter does via SIGKILL)
```
