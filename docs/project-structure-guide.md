# Project Structure Guide

This guide describes the directory structure of this Go monorepo.  
It supports multiple services, contract‑first APIs (OpenAPI + gRPC), observability, secure networking (mTLS), Kubernetes deployment, and multiple testing strategies.

---

## Directory Layout

```
/project-root
│── api/                      # API contracts (source of truth)
│   ├── auth/oapi             # Public HTTP API
│   └── user/
│       ├── grpc              # Internal service communication
│       └── oapi              # Private HTTP API
│
│── services/                 # Deployable microservices
│   ├── auth/
│   │   ├── cmd/              # Application entrypoint
│   │   ├── internal/
│   │   │   ├── handler/      # HTTP transport layer
│   │   │   ├── service/      # Business logic
│   │   │   ├── repository/   # Data access (memory/redis)
│   │   │   ├── gateway/      # Calls to other services (gRPC clients)
│   │   │   ├── middleware/   # HTTP middleware
│   │   │   ├── config/       # Runtime configuration
│   │   │   ├── token/        # JWT & opaque tokens
│   │   │   └── api/          # Generated OpenAPI server interfaces
│   │   └── pkg/model/        # Service-owned domain models
│   │
│   └── user/
│       ├── cmd/
│       ├── internal/
│       │   ├── handler/      # HTTP + gRPC endpoints
│       │   ├── service/      # Business logic
│       │   ├── repository/   # MySQL persistence
│       │   ├── interceptor/  # gRPC interceptors
│       │   ├── db/           # sqlc generated code
│       │   └── api/          # Generated OpenAPI server
│       ├── db/mysql/
│       │   ├── migrations/   # Database schema migrations
│       │   └── sqlc.yaml     # Query generation
│       └── pkg/model/
│
│── pkg/                      # Shared reusable libraries
│   └── obs/                  # Observability platform
│       ├── logging/
│       ├── metrics/
│       ├── tracing/
│       ├── correlation/
│       └── otel/
│
│── infra/                    # Runtime platform infrastructure
│   ├── envoy/                # Service mesh & gateway config
│   ├── obs/                  # Prometheus, Grafana, Loki, Tempo
│   └── security/             # mTLS certificates & CA
│
│── deploy/helm/              # Kubernetes deployment charts
│
│── make/                     # Modular Makefile system
│   ├── oapi.mk
│   ├── grpc.mk
│   ├── sqlc.mk
│   ├── migrate.mk
│   ├── helm.mk
│   └── security.mk
│
│── docs/                     # Engineering documentation
│── test/                     # Integration, BDD, and contract tests
│── docker-compose.yaml       # Local development environment
│── Makefile
│── go.mod
└── go.sum
```

---

## Structure Philosophy

### **api/**

Contains all API contracts.  
This is the source of truth — implementations must conform to these definitions.

Rules:

* OpenAPI → HTTP APIs
* gRPC → service‑to‑service APIs
* Generated code must never be edited
* Services communicate through generated clients only

---

### **services/**

Each service is independently deployable and internally layered.

Layers:

| Layer | Responsibility |
|------|------|
| handler | Transport (HTTP/gRPC) |
| service | Business logic |
| repository | Storage |
| gateway | External service calls |
| config | Runtime configuration |

This separation keeps business logic independent of transport and storage.

---

### **pkg/**

Reusable libraries shared across services.

Currently contains the observability platform:

* structured logging
* distributed tracing
* metrics
* correlation context

Anything here must be safe for reuse.

---

### **infra/**

Represents the runtime environment required for production‑like execution locally.

Includes:

* Envoy proxy
* mTLS certificate authority
* Prometheus / Grafana / Loki / Tempo
* OpenTelemetry collector

The goal is to make local behavior match production behavior.

---

### **deploy/**

Helm charts used for Kubernetes deployment.

All services and infrastructure are deployable together as a single platform.

---

### **make/**

Modular build system. Each file manages a specific concern.

| File | Purpose |
|----|----|
| oapi.mk | Generate OpenAPI code |
| grpc.mk | Generate protobuf |
| sqlc.mk | Generate DB queries |
| migrate.mk | Run migrations |
| helm.mk | Kubernetes deploy |
| security.mk | Certificates |

---

### **test/**

Supports multiple testing approaches:

* BDD (Gherkin features)
* Consumer contract tests (Pact)
* Provider verification
* Cross‑service integration

The repository prioritizes verifying service interaction correctness.

---

## Best Practices

* Do not import another service’s internal package
* Never modify generated code
* Keep business logic inside `service/`
* Handlers only translate transport ↔ domain
* Repositories must not contain business logic
* Inter‑service communication only through contracts
* If production requires infrastructure, run it locally too

---

This structure prioritizes scalability, safety, and integration stability across services.
