# OpenAPI Code Generation Guide (oapi-codegen)

This guide explains how to install, configure, and use **oapi-codegen** to generate server and client code from OpenAPI specifications in this project.

---

# 📦 Installation

Install the latest version of **oapi-codegen**:

```
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

Verify installation:

```
oapi-codegen -version
```

Ensure the binary path is available:

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

# 📁 Project Structure for Codegen

Typical service folder layout:

```
services/
  └── ping/
        ├── internal/
        │     ├── api/        # Generated server interfaces
        │     ├── client/     # Generated client code
        │     ├── handler/    # Handlers implementing the API
        │     └── repository/ # DB & cache interactions
        ├── cmd/main.go
        ├── oapi-codegen.server.yaml
        ├── oapi-codegen.client.yaml
        └── openapi.yaml      # API specification
```

---

# 🛠️ Generating Server Code

Example command:

```
oapi-codegen --config=services/ping/oapi-codegen.server.yaml services/ping/openapi.yaml
```

A typical `oapi-codegen.server.yaml`:

```yaml
generate:
  types: true
  server: true
output: internal/api/server.gen.go
package: api
```

This produces interfaces such as:

```go
type ServerInterface interface {
    GetPing(ctx echo.Context) error
}
```

---

# 📡 Generating Client Code

Command:

```
oapi-codegen --config=services/ping/oapi-codegen.client.yaml services/ping/openapi.yaml
```

Example `oapi-codegen.client.yaml`:

```yaml
generate:
  client: true
  types: true
output: internal/client/client.gen.go
package: client
```

This produces a typed client:

```go
client := client.NewClientWithResponses("http://localhost:8080")
res, err := client.GetPingWithResponse(ctx)
```

---

# 📄 Example OpenAPI Spec (openapi.yaml)

```yaml
openapi: 3.0.0
info:
  title: Ping API
  version: 1.0.0
paths:
  /ping:
    get:
      summary: Ping endpoint
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
```

---

# 🛠️ Makefile Integration

Add shortcuts to regenerate code:

```Makefile
.PHONY: gen

gen: ping/gen

ping/gen:
	oapi-codegen --config=services/ping/oapi-codegen.server.yaml services/ping/openapi.yaml
	oapi-codegen --config=services/ping/oapi-codegen.client.yaml services/ping/openapi.yaml
```

Run everything:

```
make gen
```

Run for a specific service:

```
make ping/gen
```

---

# 🧩 Implementing the Generated Interfaces

After generating server interfaces, implement them in your handler layer:

```go
type PingHandler struct{}

func (h *PingHandler) GetPing(ctx echo.Context) error {
    return ctx.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
```

Register with the router:

```go
api.RegisterHandlers(router, &PingHandler{})
```

---

# 🌐 Updating the OpenAPI Spec

Whenever you:

* Add new routes
* Modify request or response models
* Change API paths or methods

Update `openapi.yaml` and regenerate code:

```
make gen
```

---

# ⭐ Best Practices

* Keep `openapi.yaml` under version control
* Use separate `.yaml` configs for server and client
* Regenerate code whenever API schemas change
* Keep generated files in `internal` to avoid leaking API types
* Avoid editing generated files manually

---

# 📚 References

* oapi-codegen: [https://github.com/oapi-codegen/oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
* OpenAPI Specification: [https://swagger.io/specification/](https://swagger.io/specification/)
