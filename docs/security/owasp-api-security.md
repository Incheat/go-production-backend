# OWASP API Security Top 10 (Summary & Best Practices)

The **OWASP API Security Top 10** highlights the most critical security risks affecting modern APIs. This document summarizes each risk, explains its meaning, and provides recommended mitigations tailored for Go + Gin applications.

---

# 🛡️ OWASP API Security Top 10 (2019)

## **API1:2019 — Broken Object Level Authorization (BOLA)**

**Description:**
Attackers manipulate object IDs in requests to access data that doesn’t belong to them.

**Impact:** Unauthorized access to user or resource data.

**Mitigation:**

* Validate object ownership for every resource access
* Always authorize based on authenticated user ID
* Avoid exposing sequential or guessable IDs

---

## **API2:2019 — Broken Authentication**

**Description:**
Weak authentication flows, poor session handling, or insecure tokens.

**Mitigation:**

* Use **JWT + Refresh Token** architecture
* Implement token expiration and rotation
* Reject invalid/expired tokens with proper middleware

---

## **API3:2019 — Excessive Data Exposure**

**Description:**
APIs returning more data than necessary, potentially leaking sensitive information.

**Mitigation:**

* Use response DTOs instead of raw DB models
* Avoid sending sensitive fields (password, token, internal IDs)
* Implement data sanitization at handler or controller level

---

## **API4:2019 — Lack of Resources & Rate Limiting**

**Description:**
APIs that allow unlimited requests may become vulnerable to denial-of-service attacks.

**Mitigation:**

* Implement rate limiting (Redis-based preferred)
* Throttling for expensive operations
* Return `429 Too Many Requests` when necessary

---

## **API5:2019 — Broken Function Level Authorization**

**Description:**
Incorrect role or privilege checks allow users to access admin-level operations.

**Mitigation:**

* Implement **RBAC** (Role-Based Access Control) or **ABAC**
* Enforce role checks in middleware or handlers
* Use separate route groups for admin endpoints

---

## **API6:2019 — Mass Assignment**

**Description:**
Automatic binding of JSON data into structs may cause attackers to overwrite restricted fields.

**Mitigation:**

* Use input DTOs with **whitelisted fields only**
* Avoid binding directly to DB entities
* Validate and sanitize input fields

---

## **API7:2019 — Security Misconfiguration**

**Description:**
Common mistakes such as permissive CORS, overly detailed error messages, or improper environment configuration.

**Mitigation:**

* Implement environment-specific CORS rules
* Avoid exposing stack traces in responses
* Use minimal privilege IAM roles (AWS/GCP)
* Keep secrets out of source control

---

## **API8:2019 — Injection**

**Description:**
SQL, NoSQL, or command injection via unsafe input handling.

**Mitigation:**

* Use ORM (GORM) with prepared statements
* Sanitize inputs
* Never build SQL queries through string concatenation

---

## **API9:2019 — Improper Assets Management**

**Description:**
Old API versions, outdated endpoints, or untracked routes expose you to vulnerabilities.

**Mitigation:**

* Version your API (`/v1`, `/v2`)
* Document endpoints properly
* Deprecate old versions gradually and track usage

---

## **API10:2019 — Insufficient Logging & Monitoring**

**Description:**
Lack of meaningful logs or failure alerts make incident response difficult.

**Mitigation:**

* Log authentication attempts, errors, and key operations
* Centralize logs (ELK, CloudWatch, GCP Logging)
* Add alerts for anomalies or high error rates

---

# 🔐 Security Features Recommended for This Project

Below are specific security practices that strengthen the project and directly map to OWASP API Top 10 items:

### ✔️ Environment-Specific + Path-Based CORS

(Addresses **API7: Security Misconfiguration**)

### ✔️ JWT + Refresh Token Lifecycle

(Addresses **API2: Broken Authentication**)

### ✔️ Redis-Based Rate Limiting Middleware

(Addresses **API4: Rate Limiting**)

### ✔️ RBAC Authorization Layer

(Addresses **API5: Broken Function Level Authorization**)

### ✔️ Response Data Sanitization

(Addresses **API3: Excessive Data Exposure**)

---

# 📚 References

* OWASP API Security Project: [https://owasp.org/www-project-api-security/](https://owasp.org/www-project-api-security/)
* OWASP API Security Top 10 (2019)
