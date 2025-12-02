# Repository Naming Guide

## Overview

Repository naming directly impacts readability, organizational consistency, and ease of discovery. This guide outlines conventions and best practices for naming repositories—especially Go projects and microservice-style repos.

---

## Core Principles

### 1. **Use lowercase, kebab-case**

Consistent with common open-source and Unix conventions.

* ✔️ `auth-service`
* ✔️ `payment-gateway`
* ❌ `AuthService`
* ❌ `payment_service`

### 2. **Avoid language-specific prefixes**

Do not include `go-` or `golang-` unless the repo is *specifically* Go tooling.

* ❌ `go-auth-service`
* ✔️ `auth-service`

### 3. **Name based on domain, not implementation**

Focus on *what* it does, not *how* it’s built.

* ✔️ `email-sender`
* ❌ `smtp-service`

### 4. **Avoid deployment details**

Infra choices shouldn’t appear in repo names.

* ✔️ `auth-service`
* ❌ `auth-service-k8s`

### 5. **Prefer singular nouns**

Unless used for exporters, collectors, or aggregators.

* ✔️ `order-service`
* ❌ `orders-service`

### 6. **Ensure repo name ≈ Go module name**

Example:

```
github.com/acme/auth-service
```

In `go.mod`:

```
module github.com/acme/auth-service
```

---

## Patterns You Can Use

### Microservice-style

* `user-service`
* `billing-service`
* `inventory-api`

### Domain-oriented

* `scheduler`
* `registry`
* `controller`

### Utility or tooling

* `schema-migrator`
* `task-runner`

---

## Special Cases

### Repos for Learning, Experiments, or POCs

* `go-playground`
* `go-sandbox`
* `go-experiments`
* `go-lab`

---

## Recommendations

For most Go projects:

* Use **lowercase**, **kebab-case**, domain-focused names
* Avoid technology or version-specific details
* Pick names that remain valid even if responsibilities expand slightly
