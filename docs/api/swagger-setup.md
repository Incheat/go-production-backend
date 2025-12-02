# Swagger / OpenAPI Setup Guide

This guide explains how to generate, configure, and serve Swagger (OpenAPI) documentation in a Go + Gin project using **swaggo/swag**.

---

# 📦 Installation

Install the required Swagger tools:

```
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

Install the latest binary:

```
go install github.com/swaggo/swag/cmd/swag@latest
```

Check that `swag` is available:

```
swag -v
```

---

# 🔧 Environment Setup

Ensure your `GOPATH/bin` is in the PATH:

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

If using zsh:

```
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

---

# 🏁 Initialize Swagger Docs

To generate Swagger files, run from project root:

```
swag init -g cmd/api/main.go
```

This creates:

```
docs/
  docs.go
  swagger.json
  swagger.yaml
```

These files include:

* OpenAPI specification (JSON/YAML)
* Auto-generated metadata from comments

---

# ✍️ Writing Swagger Annotations

Place Swagger comments above handlers or controllers.

Example:

```go
// @Summary Health Check
// @Description Returns API health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
    c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
```

Annotations describe:

* Endpoint summary
* Description
* Tags (grouping)
* Request/Response schemas
* HTTP methods & routes

---

# 🌐 Serve Swagger UI

To enable Swagger UI in the `main.go`:

```go
import (
    _ "github.com/yourname/project/docs" // Swagger docs
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
)

func main() {
    r := gin.Default()
    // register routes...

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.Run()
}
```

Open Swagger UI:

```
http://localhost:8080/swagger/index.html
```

---

# 🔄 Regenerating Swagger Docs

Whenever you:

* Add new routes
* Modify handler annotations
* Update models

Run:

```
swag init -g cmd/api/main.go
```

Consider adding a Makefile target:

```
make swagger
```

---

# 📁 Directory Structure Example

```
cmd/api/main.go
internal/handlers/...       # Contains annotated handlers
internal/models/...         # Contains response/request models
internal/...                
docs/                       # Auto-generated swagger docs
```

---

# 🧪 Verify Documented Endpoints

Swagger UI validates:

* Route availability
* Request parameters
* Response objects
* HTTP method correctness

Use it to confirm the API contract.

---

# ⭐ Best Practices

* Keep annotations near handler implementations
* Organize tags by service/domain
* Document all error responses (400, 404, 500)
* Use DTOs instead of DB models in Swagger
* Do not manually modify `docs.go` (auto-generated)

---

# 📚 References

* Swaggo Repository: [https://github.com/swaggo/swag](https://github.com/swaggo/swag)
* Swagger/OpenAPI Specification: [https://swagger.io/specification/](https://swagger.io/specification/)
