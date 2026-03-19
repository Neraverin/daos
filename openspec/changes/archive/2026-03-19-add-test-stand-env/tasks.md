## 1. Environment Setup

- [x] 1.1 Create `.env` file with test stand credentials
- [x] 1.2 Add `.env` to `.gitignore` (already present)
- [x] 1.3 Add `.env.example` file with template credentials

## 2. TestStand Connection Manager

- [x] 2.1 Create `tests/e2e/teststand/connection.go` package
- [x] 2.2 Implement `LoadEnv()` function to read `.env` file
- [x] 2.3 Implement `TestStand` struct with SSH client
- [x] 2.4 Implement `Connect()` method to establish SSH connection
- [x] 2.5 Implement `Execute()` method to run commands remotely
- [x] 2.6 Add singleton pattern with `GetTestStand()` function

## 3. Execution Mode Support

- [x] 3.1 Define `Mode` type (Local, Remote, Auto)
- [x] 3.2 Implement `GetExecutionMode()` function
- [x] 3.3 Implement auto-detect logic (check .env exists + connection works)
- [x] 3.4 Add `GetTarget()` function returning Local or Remote

## 4. TestStand Integration

- [x] 4.1 Create `TestStand` wrapper that handles both modes
- [x] 4.2 Update existing e2e tests to use `GetTestStand()`
- [x] 4.3 Add `RunLocal()` / `RunRemote()` execution helpers

## 5. Testing

- [x] 5.1 Add unit tests for `LoadEnv()`
- [x] 5.2 Add unit tests for `GetExecutionMode()`
- [x] 5.3 Add integration test for SSH connection (local tests pass)
