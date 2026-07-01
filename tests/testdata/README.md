# AstraMind Test Data Repository

This directory contains reusable datasets used by automated tests.

## Structure

conversations/
Reusable conversation datasets.

expected/
Expected outputs for automated tests.

exports/
Reference export files.

sessions/
Reusable session JSON files.

## Purpose

The same datasets are reused across:

- Unit Tests
- Integration Tests
- Regression Tests
- Mock AI
- Search
- Streaming
- Knowledge Base

Keeping test data centralized avoids duplication and improves maintainability.