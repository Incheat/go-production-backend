# Contributing Guide

Thank you for considering contributing to this project! This document explains how to propose changes, follow project conventions, and submit high‑quality contributions.

---

# 📌 Code of Conduct

By participating in this project, you agree to maintain respectful and constructive communication. Please be kind, professional, and considerate toward other contributors.

---

# 🧱 Project Structure Overview

Before contributing, familiarize yourself with the structure:

* `cmd/` — application entry points
* `services/` — individual service modules
* `internal/` — service‑specific logic
* `pkg/` — shared utilities
* `docs/` — documentation
* `test/` — BDD, unit, and contract tests

For a detailed explanation, see:
`docs/project-structure-guide.md`

---

# 🛠️ Development Workflow

## 1. Fork and Clone the Repository

```
git clone https://github.com/yourname/repo.git
cd repo
```

## 2. Create a Feature Branch

```
git checkout -b feature/my-improvement
```

Use descriptive branch names:

* `feature/auth-refactor`
* `fix/rate-limiter-bug`
* `docs/update-swagger-guide`

---

# 🧪 Testing

All contributions must pass tests before submission.

### Run unit tests

```
go test ./...
```

### Run BDD tests

```
godog
```

### Run contract tests

```
go test ./test/pact/...
```

Ensure the following:

* All tests pass
* New features include unit tests
* API changes include updated BDD scenarios

---

# 🚦 Commit Message Guidelines

This project uses **Conventional Commits**.

Format:

```
type(scope): short description
```

Common types:

* `feat` — new feature
* `fix` — bug fix
* `docs` — documentation changes
* `refactor` — code restructuring with no behavior change
* `test` — adding or updating tests
* `chore` — CI/CD or tooling work

See: `docs/commit-message-guide.md`

---

# 🔄 Code Style & Formatting

Follow Go best practices:

* Run `go fmt ./...`
* Run linters (e.g., golangci-lint) if available
* Keep functions small and focused
* Prefer composition over inheritance
* Use meaningful names and avoid abbreviations

See: `docs/best-practices.md`

---

# 🧬 API Changes

If you modify or add API endpoints:

1. Update `openapi.yaml` for the affected service
2. Regenerate code with

   ```
   make gen
   ```
3. Update Swagger annotations as needed
4. Update BDD tests and any API documentation

See:

* `docs/api/oapi-codegen-guide.md`
* `docs/api/swagger-setup.md`

---

# 🛡️ Security & Data Protection

Follow project security standards:

* Avoid exposing sensitive fields in responses
* Validate and sanitize all input
* Use rate limiting for sensitive endpoints
* Apply RBAC where needed

Reference:
`docs/security/owasp-api-security.md`

---

# 📄 Submitting Changes

When you're ready:

1. Commit your changes following Conventional Commits
2. Push your branch
3. Open a Pull Request
4. Fill out the PR template (if available)
5. Explain the motivation, design, and test coverage
6. Request review from maintainers

Your PR will be reviewed for:

* Code quality
* Tests & documentation updates
* Architecture consistency

---

# 🙌 Thank You!

Your contributions make this project better.
Feel free to open issues for bugs, feature requests, or questions.
