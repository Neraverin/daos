## Context

E2E tests currently run on the local machine and test package installation. Future tests need to connect to a remote test stand for more comprehensive testing scenarios.

Test stand credentials:
- Hostname: `****`
- Username: `****`
- Password: `****`

## Goals / Non-Goals

**Goals:**
- Store test stand credentials in `.env` file
- Create reusable `TestStand` connection manager
- Support SSH connections to remote hosts
- Allow multiple e2e tests to share the same connection

**Non-Goals:**
- Implementing full SSH key management (password-based for now)
- Automatic connection pooling
- Connection health checks or reconnection logic

## Decisions

### 1. Use golang.org/x/crypto/ssh for SSH connections

**Decision:** Use `golang.org/x/crypto/ssh` package for SSH connectivity.

**Rationale:** Standard Go library, no external dependencies beyond Go modules.

### 2. TestStand singleton pattern

**Decision:** Create `TestStand` as a singleton that initializes once and can be reused across tests.

**Rationale:** Avoids repeated SSH connections in each test. Tests can share the same connection.

```go
var stand *TestStand

func GetTestStand(t *testing.T) *TestStand {
    if stand == nil {
        stand = NewTestStand(hostname, username, password)
        stand.Connect()
    }
    return stand
}
```

### 3. .env file structure

**Decision:** Store credentials in `.env` file with simple key-value format.

```
TEST_STAND_HOSTNAME=****
TEST_STAND_USERNAME=****
TEST_STAND_PASSWORD=****
```

**Rationale:** Simple, standard format. Go's `os.Setenv()` or a lightweight parser can read it.

### 4. Execution mode with auto-detect

**Decision:** Check `TEST_STAND_MODE` environment variable:
- `"remote"` → require `.env`, fail test if connection fails
- `"local"` → always run locally (ignore `.env`)
- Not set → auto-detect: use test stand if `.env` exists AND connection succeeds, else run locally

**Rationale:** Flexible default (auto-detect) while allowing explicit override for CI/CD or debugging.

```go
func GetExecutionMode() Mode {
    switch os.Getenv("TEST_STAND_MODE") {
    case "remote":
        return ModeRemote
    case "local":
        return ModeLocal
    default:
        return ModeAuto
    }
}
```

## Risks / Trade-offs

- **[Risk]** Password in plain text in `.env` → **Mitigation:** `.env` is in `.gitignore`, credentials never committed
- **[Risk]** Stand becomes unavailable during test run → **Mitigation:** Connection checked at test start, fails fast if in remote mode
