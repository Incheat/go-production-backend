# Logger Module Overview

This document describes the components inside the `pkg/logger/` directory. The logger module provides a unified logging interface, Gin middleware for HTTP request logging, and tools for structured error tracking.

It is designed to support production-grade logging practices across all services in the project.

---

# 📁 Directory Structure

```
pkg/logger/
 ├── logger.go
 ├── middleware.go
 └── error_logger.go
```

---

# 📝 logger.go — Core Logger Wrapper

The central logging component. It wraps popular logging frameworks such as **Zap**, **Logrus**, or **Zerolog**, providing a consistent interface across your application.

### Features

* Create a global or request-scoped logger
* Structured logging with fields
* Log levels: debug, info, warn, error, fatal
* Optional JSON output for production
* Custom formatting for local development

### Example Usage

```go
log := logger.New()
log.Info("user logged in", logger.Fields{"user_id": user.ID})
```

### Best Practices

* Use structured logs instead of plain text
* Include contextual fields (request ID, user ID)
* Use JSON format in production for log aggregation tools

---

# 🌐 middleware.go — HTTP Request Logging (Gin)

Middleware for automatically logging HTTP request details.

### What It Logs

* HTTP method
* URL path
* Response status code
* Latency/duration
* IP address
* Optional request ID (if included via middleware)

### Example Output

```
{"method":"GET", "path":"/ping", "status":200, "duration":"2ms"}
```

### Example Middleware Registration

```go
r := gin.Default()
r.Use(logger.GinMiddleware())
```

### Best Practices

* Always log request latency
* Mask sensitive fields (password, tokens)
* Include request ID for traceability

---

# ❗ error_logger.go — Error Tracking & Reporting

Handles logging of errors in a structured, consistent format.

### Responsibilities

* Log unexpected errors with stack traces
* Log application-level errors (validation, DB issues)
* Integrate with external monitoring tools (optional)
* Distinguish between user errors (400-level) and server errors (500-level)

### Example Usage

```go
if err != nil {
    logger.Error("failed to fetch user", logger.Fields{"error": err.Error()})
}
```

### Best Practices

* Include contextual metadata when logging errors
* Avoid logging the same error multiple times
* Use error wrapping to capture call stack information

---

# 🧭 Summary

The logger module provides:

* Unified, framework-agnostic logging API
* Gin request logging middleware
* Centralized error tracking

It supports structured logs ideal for:

* Cloud environments (AWS, GCP)
* Log aggregators (ELK, CloudWatch, Stackdriver)
* Debugging request flows across microservices

For additional observability integrations, consider adding:

* Request IDs and correlation IDs
* Tracing (OpenTelemetry)
* Metrics instrumentation (Prometheus)

---

# 📚 References

* Zap: [https://github.com/uber-go/zap](https://github.com/uber-go/zap)
* Logrus: [https://github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)
* Zerolog: [https://github.co](https://github.co)
