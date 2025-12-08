# A Practical Guide to Test Doubles  
*Dummy • Stub • Spy • Mock • Fake*  
*Classical vs. Modern Testing Perspectives*

## 1. What Is a Test Double?

A **Test Double** is any object used in testing to **replace a real collaborator**.  
The word comes from filmmaking: a *stunt double* or *body double* performs in place of a real actor.

A Test Double stands in for a real dependency to:
- control test behavior  
- isolate business logic  
- avoid external unreliability  
- improve determinism and performance  
- record interactions  

Subtypes include Dummy, Stub, Spy, Mock, and Fake.

---

## 2. Types of Test Doubles

### 2.1 Dummy
- Fills parameter slots  
- Never used  
- No logic  

```python
class DummyLogger:
    def log(self, msg):
        pass
```

---

### 2.2 Stub
- Returns predetermined outputs  
- Ignores input  
- No logic or call verification  

```python
class UserRepoStub:
    def get_user(self, user_id):
        return User(id=1, name="StubUser")
```

---

### 2.3 Spy
- Records interactions  
- May have real behavior  
- Verification happens *after* execution  

```python
class EmailServiceSpy:
    def __init__(self):
        self.sent = []

    def send(self, to, body):
        self.sent.append((to, body))
```

---

### 2.4 Mock
- Strict behavior verification  
- Expectations set *before* execution  

```python
mock_repo.expect("save").called_once_with(user)
```

---

### 2.5 Fake
- Working but simplified implementation  
- Contains real logic  

```python
class FakeDB:
    def __init__(self):
        self.data = {}

    def save(self, key, value):
        self.data[key] = value

    def get(self, key):
        return self.data.get(key)
```

---

## 3. Comparison Table

| Type | Provides Data? | Records Calls? | Contains Logic? | Purpose |
|------|----------------|----------------|------------------|---------|
| Dummy | ❌ | ❌ | ❌ | Fill params |
| Stub | ✔️ | ❌ | ❌ | Control outputs |
| Spy | Optional | ✔️ | Optional | Record interactions |
| Mock | Optional | ✔️ (strict) | ❌ | Behavior verification |
| Fake | ✔️ | ❌ | ✔️ | Functional replacement |

---

## 4. Classical vs. Modern Testing

### 4.1 Classical View (Fowler / GOOS)
- Unit tests → heavy use of doubles  
- Integration tests → use real components  
- E2E → all real systems  

Focus: theoretical purity & isolation.

---

### 4.2 Modern View (Microservices / Cloud Era)

Real-world constraints mean:
- Real systems are slow, unstable, or expensive  
- DBs hard to reset  
- External APIs rate-limited  

So **modern Integration Tests often use doubles**, including:
- Fake DBs (SQLite, Testcontainers)  
- Stubbed HTTP APIs (WireMock)  
- Fake message queues (Localstack, in-memory Kafka)  

Only E2E sticks to fully real infrastructure.

---

## 5. When to Use What

### Unit Tests  
Use: Stub, Mock, Spy, Fake, Dummy  
Avoid: real DB, real external API  

### Integration / Service Tests  
Use: Fake DB, stubbed API, Localstack, Testcontainers  
Avoid: mocking internal business logic  

### E2E  
Use: all real systems  

---

## 6. Summary

> Test Doubles replace real dependencies during testing.  
> Each subtype serves a different purpose.  
> Modern engineering uses doubles across multiple layers—not just Unit Tests—because real cloud/microservice components are often impractical in tests.
