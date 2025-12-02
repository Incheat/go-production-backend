# Commit Message Guide

This project follows the **Conventional Commits** specification.
More details: [https://www.conventionalcommits.org/en/v1.0.0/](https://www.conventionalcommits.org/en/v1.0.0/)

---

## Format

A commit message should follow this structure:

```
type(scope): short description
```

### Common Types

* **feat** — new functionality
* **fix** — bug fix
* **test** — adding or improving tests
* **refactor** — restructuring code without changing behavior
* **docs** — documentation updates
* **chore** — tooling, CI, or maintenance tasks
* **perf** — performance improvements

### Rules

* Use **lowercase** for the type.
* Use the **imperative mood** for the description (e.g., "add", "update", "fix").
* Keep the subject line **brief and clear**.
* Add a body when additional context is helpful.

---

## Examples

### Adding a new feature scenario

```
test(godog): add scenario for user login flow
```

Body:

```
Adds a Gherkin scenario covering successful and failed login cases.
Defines expected behavior before implementing handlers.
```

### Implementing steps for a feature

```
feat(auth): implement login handlers to satisfy godog steps
```

Body:

```
Adds handler logic and supporting service code required for the login feature
scenarios to pass.
```
