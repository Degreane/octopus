# AGENT.md - Octopus Development Guide

## Build/Test/Lint Commands
- `make build` - Build the server binary
- `make run` - Run the server directly
- `make dev` - Start development server with hot reload (using air)
- `make test` - Run all tests (`go test ./...`)
- `npm run build-css` - Build Tailwind CSS styles
- `go run cmd/server/main.go` - Direct server execution

## Architecture & Structure
- **Entry Point:** `cmd/server/main.go` - Fiber web server with HTML templates
- **Internal Packages:** `internal/{database,middleware,routes,service,utilities}`
- **Public Packages:** `pkg/api/` - External API interfaces
- **Frontend:** `views/` HTML templates, `public/css/` static assets with Tailwind CSS
- **Config:** YAML-based configuration with environment overrides
- **Tech Stack:** Go + Fiber framework + MongoDB + Redis + HTML templates + Lua scripting

## Code Style & Conventions
- **Naming:** PascalCase exports, camelCase locals, descriptive verb-noun functions
- **Error Handling:** Early returns, structured logging with component tags, error wrapping
- **Imports:** Standard library first, external packages, then internal packages
- **Comments:** Package-level docs, function descriptions, inline explanations
- **Structure:** Exported functions first, related functions grouped, helpers near usage
- **Config:** Environment variables via `.env`, YAML config files with validation
