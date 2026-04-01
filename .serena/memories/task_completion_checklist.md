# Task Completion Checklist

When a task is completed, run the following:

1. **Format**: `go fmt ./...`
2. **Vet/Lint**: `go vet ./...`
3. **Build**: `go build -o phunter .` (ensure it compiles)
4. **Test**: `go test ./...` (if tests exist)
5. **Tidy**: `go mod tidy` (if dependencies changed)
