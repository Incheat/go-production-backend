# Project Structure Guide

This guide describes a clean, scalable directory structure for Go applications—supporting multiple services, OpenAPI integration, shared utilities, migrations, testing, and development tooling.

---

## Directory Layout

```
/project-root
│── cmd/                 # Entry points for top‑level binaries (e.g., api, worker)
│    └── api/
│         └── main.go
│
│── config/              # Configuration files and environment settings
│    └── config.yaml
│
│── services/             # Each service has its own isolated structure
│    └── {service_name}/
│         ├── cmd/
│         │     └── main.go
│         ├── internal/
│         │     ├── api/        # OpenAPI-generated server interfaces
│         │     ├── client/     # OpenAPI-generated client code
│         │     ├── controller/ # Business logic / domain controllers
│         │     ├── handler/    # API handlers (HTTP, gRPC)
│         │     ├── security/   # Auth, RBAC, middleware
│         │     └── repository/ # Database & Redis implementations
│         ├── oapi-codegen.client.yaml
│         ├── oapi-codegen.server.yaml
│         └── openapi.yaml      # API definition
│
│── pkg/                 # Shared utilities (logger, middleware, helpers)
│
│── migrations/          # Database migrations (goose, migrate, etc.)
│
│── scripts/             # CI/CD scripts, build automation, deploy tooling
│
│── test/                # Integration, contract, and BDD test structure
│    ├── pact/           # Consumer/provider contract tests
│    │    ├── consumer/
│    │    └── provider/
│    ├── features/       # Gherkin feature files
│    ├── steps/          # Step definitions for Godog
│    └── testutils/      # Shared test data, fixtures, helpers
│
│── Makefile             # Build shortcuts and developer tasks
│── README.md            # Project overview and instructions
│── docker-compose.yaml  # Local dev environment setup
│── go.mod
└── go.sum
```

---

## Structure Philosophy

### **cmd/**

Contains the entry points for your top-level executables. Each folder represents a binary.

### **config/**

Stores static configuration files such as YAML, JSON, or environment presets.

### **services/**

Each service is self-contained and follows a mini clean architecture layout. This helps when scaling to multiple microservices while sharing common patterns.

### **pkg/**

Contains libraries intended to be reusable across the entire project or externally. Anything placed here must be safe for reuse.

### **migrations/**

All database schema changes belong here. Works with common migration tools.

### **scripts/**

Automates tasks such as CI pipelines, code generation, dev setup, and deployment workflows.

### **test/**

A comprehensive testing layout supporting:

* BDD (Godog)
* Contract tests (Pact)
* Shared testing utilities

---

## Best Practices

* Keep business logic out of `cmd/`—place it in `internal/`.
* Never import from another service’s `internal/` directory.
* Keep `pkg/` small—only place code intended for reuse.
* Store all API definitions and generated code alongside each service.
* Maintain a clear separation between handlers, controllers, and repositories.
* Use Makefile targets to standardize builds and common tasks.

This structure ensures scalability, modularity, testability, and team-friendly workflows.
