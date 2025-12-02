# Security Module Overview

This document describes the components within the `internal/security/` directory. These modules collectively handle authentication, authorization, input sanitization, rate limiting, token lifecycle management, and request integrity.

All components are written for a Go + Gin backend and follow OWASP API Security best practices.

---

# 📁 Directory Structure

```
internal/security/
 ├── jwt.go
 ├── password.go
 ├── auth_middleware.go
 ├── token.go
 ├── input_sanitizer.go
 ├── rate_limiter.go
 └── signature.go
```

---

# 🔐 jwt.go — JWT Handling

Responsible for creating, signing, validating, and parsing JWT tokens.

### Features

* Generate access tokens
* Generate refresh tokens
* Validate token signature
* Check expiration and claims
* Support custom claims (user ID, roles, permissions)

### Best Practices

* Use short-lived access tokens
* Rotate refresh tokens
* Store signing secret securely (env variables or secret manager)

---

# 🔑 password.go — Password Hashing

Provides secure password hashing and validation using `bcrypt`.

### Features

* Hash plaintext passwords
* Compare stored hash with user input
* Avoids manual salt handling (bcrypt manages salt internally)

### Best Practices

* Use bcrypt cost 10–14
* Never log raw passwords
* Never store plaintext passwords in DB

---

# 🛡️ auth_middleware.go — Authentication Middleware

Gin middleware for validating `Authorization: Bearer <token>` headers.

### Responsibilities

* Extract token from header
* Validate token format
* Parse JWT and validate claims
* Attach authenticated user info to request context

### Typical Flow

1. Incoming request → middleware executes
2. Token extracted & validated
3. User ID/claims placed into context
4. Downstream handlers access user info
5. Return `401 Unauthorized` on failure

---

# 🔄 token.go — Refresh Token Lifecycle

Handles refresh token generation, storage, invalidation, and rotation.

### Features

* Issue refresh tokens alongside access tokens
* Validate refresh tokens during login or renewal
* Invalidate tokens after logout or expiration
* Support optional Redis storage for revocation lists

### Best Practices

* Store refresh tokens server-side (Redis) for invalidation
* Rotate refresh tokens on every refresh
* Use long-lived refresh tokens and short-lived access tokens

---

# 🧼 input_sanitizer.go — Input Sanitization

Prevents malicious input and protects against injection or XSS.

### Features

* Strip dangerous HTML/JS content
* Enforce allowed characters
* Sanitize path and query parameters
* Normalize whitespace and Unicode

### Best Practices

* Validate on both server and client
* Allow-list instead of block-list
* Avoid binding request bodies directly to DB models

---

# 🚦 rate_limiter.go — Rate Limiting

Implements API request throttling, typically backed by Redis.

### Features

* Per-user or per-IP rate limiting
* Sliding-window or token-bucket algorithms
* Prevents brute-force attacks and DoS attempts
* Returns `429 Too Many Requests` when exceeded

### Best Practices

* Use Redis for distributed rate limiting
* Apply stricter limits on sensitive endpoints (login, signup)
* Log rate limit violations for auditing

---

# ✒️ signature.go — Request Signing & Integrity

Provides cryptographic signature validation for secure assets or signed URLs.

### Features

* Generate signed URLs for temporary access
* Validate signatures to ensure data integrity
* Prevent tampering in file or media delivery

### Common Use Cases

* Secure image delivery
* Time-limited file access
* Webhook payload verification

---

# 🧭 Summary

The security module provides:

* Strong JWT-based authentication
* Robust password hashing
* Middleware-based authorization
* Refresh token lifecycle management
* Input sanitization and data protection
* Rate limiting and abuse prevention
* Signed URL and request integrity validation

It implements key protections aligned with:

* OWASP API Security Top 10
* Zero-trust API design principles

For additional security guidance, see:

* `docs/security/owasp-api-security.md`
