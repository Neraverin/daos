## Why

E2E tests require a test stand to run against. Credentials need to be stored securely and we need a reusable connection manager so multiple tests can share the same stand connection.

## What Changes

- Create `.env` file with test stand credentials (hostname, username, password)
- Add `.env` to `.gitignore` to prevent accidental commits
- Create `TestStand` connection manager that reads credentials from `.env` and provides reusable connections
- Use `TestStand` in e2e tests for consistent test environment

## Capabilities

### New Capabilities
- `test-stand-connection`: Reusable test stand connection manager for e2e tests
- `test-stand-env`: Store test stand configuration in `.env` file

### Modified Capabilities
- `api-testing`: Update to use `TestStand` connection manager

## Impact

- New file: `.env` (not committed to git)
- Modified: `.gitignore` (add `.env` entry)
- Modified: e2e test configuration to read from `.env`
