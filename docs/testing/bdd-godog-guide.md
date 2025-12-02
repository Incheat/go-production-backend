# BDD Testing Guide (Godog)

This guide explains how to use **Godog** for Behavior-Driven Development (BDD) testing in this project. It covers installation, directory structure, running tests, feature definitions, and tagged execution.

---

# 🧰 Installation & Setup

## Ensure `GOPATH/bin` is in your PATH

Godog installs its binary into `$(go env GOPATH)/bin`, so you must add it to your shell PATH.

```
export PATH=$PATH:$(go env GOPATH)/bin
```

### Verify installation

```
godog --version
```

---

# 📁 Project Test Structure

The project organizes BDD tests under the `test/` directory:

```
test/
 ├── pact/                # Contract testing (Pact)
 ├── features/            # Gherkin .feature files
 ├── steps/               # Step definitions for Godog
 └── testutils/           # Shared utilities, fixtures, mocks
```

### **features/**

Contains `.feature` files written in Gherkin.

### **steps/**

Contains Go files implementing the step logic.

### **testutils/**

Provides shared test helpers, dummy data, mocks, and common utilities.

---

# 📝 Writing Feature Files

Feature files describe behavior using Gherkin syntax.

Example: `test/features/login.feature`

```
Feature: User Login
  Scenario: Successful login
    Given a registered user exists
    When the user logs in with valid credentials
    Then the user receives a valid access token
```

---

# 🧩 Writing Step Definitions

Each Gherkin step must be implemented in Go.

Example: `test/steps/login_steps.go`

```go
func (s *Suite) aRegisteredUserExists() error {
    // setup test user
    return nil
}
```

Bind steps to your suite:

```go
godog.SuiteInitializer(func(s *godog.Suite) {
    s.Step(`^a registered user exists$`, suite.aRegisteredUserExists)
})
```

---

# ▶️ Running Tests

### Run **all** Godog tests

```
godog
```

### Run a specific feature file

```
godog test/features/dummy.feature
```

### Run with pretty output

```
godog --format=pretty
```

### Run scenarios with specific tags

```
godog --tags=@api
```

Example:

```
godog test/features/calculator.feature --format=pretty
```

---

# 🏷️ Tagging Conventions

Use tags to group related scenarios.

Example in `.feature`:

```
@api @login
Scenario: Successful login
```

Run only `@api` tests:

```
godog --tags=@api
```

Run everything except @wip (work in progress):

```
godog --tags="~@wip"
```

---

# 🧹 Best Practices

* Keep steps **reusable** across scenarios.
* Avoid embedding business logic in step definitions—use helper functions.
* Use `testutils/` for common utilities.
* Organize features by domain (auth, user, payments, etc.).
* Prefer small, focused scenarios over long, multi-step flows.

---

# 📚 Additional Resources

* Godog documentation: [https://github.com/cucumber/godog](https://github.com/cucumber/godog)
* Gherkin reference: [https://cucumber.io/docs/gherkin/reference/](https://cucumber.io/docs/gherkin/reference/)
