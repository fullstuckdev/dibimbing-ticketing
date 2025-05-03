# Unit Tests for Ticketing System

This directory contains unit tests for the ticketing system API.

## Structure

The tests are organized by layer:

- `service/` - Tests for service layer functionality
- `repository/` - Tests for repository layer functionality
- `controller/` - Tests for controller layer functionality

## Running Tests

To run all tests:

```bash
go test ./tests/...
```

To run tests for a specific package:

```bash
go test ./tests/service
```

To run a specific test:

```bash
go test ./tests/service -run TestLogActivity_Success
```

To run tests with coverage:

```bash
go test ./tests/... -cover
```

For a detailed coverage report:

```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Writing Tests

When writing tests, follow these guidelines:

1. Create a mock for the dependencies of the component you're testing
2. Implement the mock's methods to return the expected values for your test
3. Use the `github.com/stretchr/testify/assert` package for assertions
4. Name your tests with the format `Test[FunctionName]_[Scenario]`
5. Split your tests into Setup, Execution, and Assertion sections for clarity

## Mocks

Common mocks are defined in the test files. For example:

- `service/user_service_test.go` contains `MockUserRepository`
- `service/audit_service_test.go` contains `MockAuditRepository`

## Testing Patterns

For service tests:

- Mock the repository layer
- Test both success and error cases
- Verify repository method calls with expectations

For controller tests:

- Mock the service layer
- Test HTTP response codes and body content
- Test validation errors

For repository tests:

- Use a test database (SQLite in-memory) or mock the DB connection
