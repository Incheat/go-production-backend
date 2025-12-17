# Project Structure Guide

This guide describes a clean, scalable directory structure for Go applications—supporting multiple services, OpenAPI integration, shared utilities, migrations, testing, and development tooling.

---

## Directory Layout

```
/project-root
│── api/                # OpenAPI specs (source of truth)
│   └── helloworld/
│       └── oapi 
│         ├── internal.yaml # API definition for internal
│         ├── public.yaml   # API definition for public
│         ├── internal.client.yaml 
│         ├── internal.server.yaml 
│         ├── internal.type.yaml 
│         └── gen/
│             ├── client.gen.go
│             └── type/
│                 └── types.gen.go
│
│── config/              # Configuration files and environment settings
│    └── config.yaml
│
│── services/             # Each service has its own isolated structure
│    └── {service_name}/
│         ├── cmd/
│         │     └── main.go
│         ├── internal/
│         │     ├── api/             # OpenAPI-generated server interfaces
│         │     │   ├── oapi/         
│         │     │   │   └── gen/    # oapi-codegen output (ignored by git)
│         │     │   │       ├── public/
│         │     │   │       │   └── server/
│         │     │   │       │       └── api_gen.go
│         │     │   │       └── private/
│         │     │   │           └── server/
│         │     │   │               └── api_gen.go
│         │     │   └── router.go    # glue between generated interfaces and handlers
│         │     ├── db/
│         │     │   └── mysql/
│         │     │       └── gen/
│         │     │           └── db.go
│         │     ├── config/
│         │     │   ├── config.go    # your Config struct
│         │     │   └── loader.go    # your Load / MustLoad
│         │     ├── service(controller)/      # Business logic / domain controllers
│         │     ├── gateway/     
│         │     ├── handler/  # API handlers or StrictServerInterface implementation for (HTTP, gRPC)
│         │     ├── security/        # Auth, RBAC, middleware
│         │     ├── repository/      # Database & Redis implementations
│         │     └── test/      
│         │         └── provider/    # auth_provider_pact_test.go # verifies pact files from all consumer in one provider test
│         ├── db/
│         │   └── mysql/
│         │       ├── migrations/
│         │       │   ├── 0001_init.up.sql
│         │       │   └── 0001_init.down.sql
│         │       ├── query.sql
│         │       └── sqlc.yaml
│         └── config/                # YAML files, mounted in Docker, etc.
│             ├── config.yaml
│             ├── config.dev.yaml
│             └── config.prod.yaml
│── internal/            # Shared utilities (logger, middleware, helpers)
│
│── deploy/             # CI/CD scripts, build automation, deploy tooling, helm
│
│── test/                # Integration, contract, and BDD test structure
│    ├── pacts/           # Consumer/provider contract tests
│    │   └── consumer/
│    │       └── auth/
│    │           ├── auth_member_pact_test.go # consumer_provider_pact_test.go, contains all internal api provided by member 
│    │           └── pacts/
│    │               └── auth-member.json   # generated pact file
│    ├── features/       # Gherkin feature files
│    ├── steps/          # Step definitions for Godog
│    └── testutils/      # Shared test data, fixtures, helpers
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
