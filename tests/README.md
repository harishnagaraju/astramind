# AstraMind Testing Framework

This directory contains all testing artifacts for AstraMind.

## Structure

tests/

- manual/          Manual Test Plans
- unit/            Unit Tests
- integration/     Integration Tests
- testdata/        Sample test data
- reports/         Test execution reports

## Testing Levels

1. Unit Testing
2. Integration Testing
3. Manual Testing
4. Regression Testing

## Running Unit Tests

```bash
go test ./...
```

## Running Coverage

```bash
go test -cover ./...
```

## Verbose Output

```bash
go test -v ./...
```

Every new feature must include:

- Implementation
- Unit Tests
- Manual Test Cases
- Regression Test

No feature is considered complete without tests.