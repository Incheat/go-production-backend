# go-playground
A repo for experimenting with Go—practicing go modules, exploring frameworks, and trying out microservice patterns and small POCs.

## Example Project

This project is a Go-based backend service designed with clean architecture principles, modular service separation, automated API generation, and strong security practices. It includes a modern tech stack, testing utilities, containerized development, and CI/CD integration.

---

## 🚀 Tech Stack

### **Languages & Frameworks**

* **Go** — primary backend language
* **Gin** — high‑performance HTTP framework
* **GORM** — ORM for database access

### **Databases & Caching**

* **MySQL** — primary relational database
* **Redis (go-redis)** — caching & session management

### **Authentication & Security**

* **golang-jwt/jwt** — JWT authentication with refresh tokens
* Follows **OWASP API Security Top 10** best practices

### **Testing Tools**

* **Testify** — unit testing framework
* **GoMock** — mocking
* **Godog** — BDD testing
* **Pact** — consumer/provider contract testing

### **Deployment & DevOps**

* **Docker** — containerization
* **GitHub Actions / GitLab CI** — CI/CD

### **Cloud Providers**

* **AWS** (ECS/EKS, RDS, ElastiCache, S3)
  or
* **GCP** (Cloud Run, CloudSQL, Memorystore, GCS)

---

## 📦 Project Structure (High-Level)

For full details, see `docs/project-structure-guide.md`.

Key directories:

* `cmd/` — service entry points
* `config/` — environment configs
* `services/` — individual service modules
* `pkg/` — shared utilities
* `migrations/` — DB migrations
* `scripts/` — CI/CD & tooling scripts
* `test/` — BDD, contract tests, utilities

---

## ⚡ Quickstart

### Run the API service

```
go run cmd/api/main.go
```

### Health check

```
curl http://localhost:8080/health
```

---

## 🌎 Environment Switching

Use `APP_ENV` to start the application with different configurations.

### Test

```
APP_ENV=test go run services/{service_name}/cmd/main.go
```

### Staging

```
APP_ENV=staging go run services/{service_name}/cmd/main.go
```

### Production

```
APP_ENV=prod go run services/{service_name}/cmd/main.go
```

---

## 🧪 BDD Testing (Godog)

Ensure your `GOPATH/bin` is in your PATH:

```
export PATH=$PATH:$(go env GOPATH)/bin
```

Check installation:

```
godog --version
```

Run all tests:

```
godog
```

Run a specific feature:

```
godog test/features/dummy.feature
```

With formatting:

```
godog --format=pretty
```

With tags:

```
godog --tags=@api
```

---

## 📘 Swagger / OpenAPI Documentation

Swagger is automatically generated using `swag`.

### Install tools

```
go install github.com/swaggo/swag/cmd/swag@latest
```

Ensure PATH:

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Initialize Swagger docs

```
swag init -g cmd/api/main.go
```

This generates:

```
docs/
  docs.go
  swagger.json
  swagger.yaml
```

### Access Swagger UI

Run the server and open:

```
http://localhost:8080/swagger/index.html
```

---

## 🔧 OpenAPI Code Generation (oapi-codegen)

Install:

```
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

Check version:

```
oapi-codegen -version
```

Generate code:

```
oapi-codegen --config=service/ping/oapi-codegen.server.yaml service/ping/openapi.yaml
```

Makefile helpers:

```
make gen        # generate all services
make ping/gen   # generate ping service only
```

---

## 🔒 Security Module Overview

For details see: `docs/security-module.md`

Key components under `internal/security/`:

* `jwt.go` — token generation/parsing
* `password.go` — bcrypt hashing
* `auth_middleware.go` — authorization middleware
* `token.go` — refresh token flow
* `input_sanitizer.go` — sanitize inputs
* `rate_limiter.go` — rate limiting
* `signature.go` — signed URL & integrity validation

---

## 📝 Logging Module Overview

See `docs/logger-module.md`.

Located in `pkg/logger/`:

* `logger.go` — unified logger wrapper
* `middleware.go` — Gin request logging
* `error_logger.go` — error tracking & output sinks

---

## 📚 Additional Documentation

All extended documentation is located under `docs/`.

Recommended docs:

* `project-structure-guide.md`
* `repo-naming-guide.md`
* `license-choice-guide.md`
* `api/swagger-setup.md`
* `api/oapi-codegen-guide.md`
* `security/owasp-api-security.md`
* `testing/bdd-godog-guide.md`

---

## 📄 License

Choose the appropriate license for your project. See:
`docs/license-choice-guide.md`

---

## 🤝 Contributions

Contributions are welcome! Please open an issue or submit a pull request.

