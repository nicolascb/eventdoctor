# Contributing to EventDoctor

Thank you for your interest in contributing to EventDoctor! This guide will help you get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Pull Request Process](#pull-request-process)
- [Style Guidelines](#style-guidelines)
- [Reporting Bugs](#reporting-bugs)
- [Requesting Features](#requesting-features)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. Be kind, constructive, and professional in all interactions.

## How Can I Contribute?

- **Report bugs** — found something broken? [Open an issue](#reporting-bugs).
- **Suggest features** — have an idea? [Request a feature](#requesting-features).
- **Fix bugs** — browse [open issues](https://github.com/nicolascb/eventdoctor/issues) labeled `bug`.
- **Improve documentation** — typos, clarifications, and examples are always welcome.
- **Write code** — pick up an issue or propose a new change.

## Getting Started

1. **Fork** the repository on GitHub.
2. **Clone** your fork locally:
   ```bash
   git clone git@github.com:<your-username>/eventdoctor.git
   cd eventdoctor
   ```
3. **Add the upstream remote**:
   ```bash
   git remote add upstream git@github.com:nicolascb/eventdoctor.git
   ```
4. **Create a branch** for your work:
   ```bash
   git checkout -b feature/my-change
   ```

## Development Setup

### Prerequisites

- **Go** 1.22+ (with CGO support for the API)
- **Node.js** 20+ and npm
- **Docker & Docker Compose** (optional, for running the full stack)

### Quick Start with Docker Compose

```bash
docker-compose up
```

This starts the backend API on port `8087` and the frontend UI on port `5173`.

### Backend (Go)

```bash
cd backend

# Run the API with mock data
WITH_MOCK=1 go run ./cmd/api/main.go

# Run all tests
CGO_ENABLED=1 go test -v -race ./...

# Build binaries
CGO_ENABLED=1 go build -o eventdoctor ./cmd/api/main.go
CGO_ENABLED=0 go build -o eventdoctor-cli ./cmd/cli/*.go
```

### Frontend (React)

```bash
cd frontend

# Install dependencies
npm install

# Start dev server (http://localhost:5173)
npm run dev

# Lint
npm run lint

# Production build
npm run build
```

## Project Structure

```
eventdoctor/
├── backend/           # Go REST API, CLI, and auditor
│   ├── cmd/           # Entry points (api, cli, auditor)
│   └── internal/      # Application logic (api, service, db, commands)
│       └── db/        # Database and migrations
├── frontend/          # React + Vite web UI
│   └── src/
├── docs/              # Documentation and images
└── docker-compose.yml
```

## Making Changes

### Branch Naming

Use descriptive branch names with a prefix:

- `feature/` — new features (e.g., `feature/add-search-endpoint`)
- `fix/` — bug fixes (e.g., `fix/topic-detail-404`)
- `docs/` — documentation changes (e.g., `docs/update-api-reference`)
- `refactor/` — code refactoring (e.g., `refactor/service-layer`)

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) to ensure a consistent and structured commit history.

Guidelines:

- **Format:** `<type>[optional scope]: <description>`
- **Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- **Mood:** Use the imperative mood for the description (e.g., "add feature" not "added feature")
- **Length:** Keep the subject line under 72 characters
- **Body:** Add a body for complex changes explaining **why**, not just **what**

Example:

```
feat(api): add topic search endpoint

Implements full-text search across topic names and descriptions.
The search uses SQLite FTS5 for efficient querying.
```

### Testing

- **Backend**: write tests for new functionality. Run `CGO_ENABLED=1 go test -v -race ./...` from the `backend/` directory and make sure all tests pass.
- **Frontend**: ensure `npm run lint` and `npm run build` pass without errors.

### Before Submitting

1. Sync your branch with upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```
2. Run all checks:
   ```bash
   # Backend
   CGO_ENABLED=1 go test -C ./backend -v -race ./...

   # Frontend
   cd frontend && npm run lint && npm run build
   ```
3. Ensure your changes don't introduce regressions.

## Pull Request Process

1. Push your branch to your fork and open a pull request against `main`.
2. Fill in the PR description with:
   - **What** the change does
   - **Why** the change is needed
   - **How** to test it
3. A maintainer will add the `ci-approved` label to trigger CI checks.
4. Address any review feedback by pushing additional commits.
5. Once approved, a maintainer will merge your PR.

> **Note:** CI workflows for pull requests are triggered only when the `ci-approved` label is added.

## Style Guidelines

### Go

- Follow standard `gofmt` formatting.
- Use the `internal/logger` package for logging (not `fmt.Println` or `log.Printf`).
- Database functions should accept `SQLExecutor` interface, not `*sql.DB` directly.
- Add new endpoint methods to the `api.Service` interface first.

### TypeScript / React

- Follow the existing ESLint configuration.
- Use TypeScript types — avoid `any` where possible.
- Place API calls in `src/lib/api.ts` and types in `src/types/index.ts`.
- Create hooks in `src/hooks/` for data fetching.

### General

- Keep changes focused — one concern per PR.
- Write meaningful variable and function names.
- Comment only when the code needs clarification.

## Reporting Bugs

Open an issue with the following information:

- **Summary**: a clear, concise description of the bug.
- **Steps to Reproduce**: numbered steps to trigger the bug.
- **Expected Behavior**: what you expected to happen.
- **Actual Behavior**: what actually happened.
- **Environment**: OS, Go version, Node.js version, browser (if applicable).
- **Screenshots or Logs**: if relevant.

## Requesting Features

Open an issue describing:

- **Problem**: what problem does this feature solve?
- **Proposed Solution**: how you envision it working.
- **Alternatives Considered**: other approaches you thought about.
- **Additional Context**: mockups, examples, or references.

---

## License

By contributing to EventDoctor, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).