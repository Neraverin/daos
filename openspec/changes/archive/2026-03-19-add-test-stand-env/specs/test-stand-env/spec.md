# test-stand-env Specification

## Purpose

Stores test stand configuration in `.env` file for e2e tests, with proper security practices to prevent accidental credential commits.

## ADDED Requirements

### Requirement: .env file stores test stand credentials
The system SHALL store test stand configuration in a `.env` file.

#### Scenario: .env file format
- **WHEN** `.env` file is created for test stand
- **THEN** it contains `TEST_STAND_HOSTNAME`, `TEST_STAND_USERNAME`, `TEST_STAND_PASSWORD`
- **AND** values are in key=value format

#### Scenario: .env file is gitignored
- **WHEN** `.env` file is created
- **THEN** it is added to `.gitignore`
- **AND** credentials are never committed to version control

### Requirement: .env file is optional
The system SHALL work without `.env` file when running in local mode.

#### Scenario: Missing .env file
- **WHEN** `.env` file does not exist
- **AND** `TEST_STAND_MODE` is not set to `"remote"`
- **THEN** tests run locally without error
