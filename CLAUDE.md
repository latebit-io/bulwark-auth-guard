# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## IMPORTANT: Always Start in Ask Mode

**CRITICAL:** When Claude Code opens this project, ALWAYS start in ask mode. Do NOT make any code changes without explicit user permission. Ask questions, provide guidance, and explain options - but never edit files unless the user specifically requests it.

## Project Overview

**bulwark-auth-guard** is a Go client library for the bulwark-auth authentication service. It provides account management and authentication flows including password-based auth, magic code (passwordless) auth, token refresh, and revocation.

## Development Commands

### Running Tests
Tests are **integration tests** that require running services:
- bulwark-auth server at `http://localhost:8080`
- MailHog at `http://localhost:8025` (for email verification testing)

```bash
# Run all tests
go test -v

# Run specific test
go test -run TestAccountCreate

# Standard Go commands
go fmt ./...
go vet ./...
go mod tidy
```

## Architecture

### Main Components

**Guard** - Main entry point that aggregates Account and Authenticate clients
- Created via `NewGuard(baseUrl, *http.Client)`
- Provides unified access to all functionality

**Account** - Manages account lifecycle
- `Create(email, password)` - Creates account, triggers verification email
- `Verify(email, verificationToken)` - Verifies account via email token
- `ChangePassword(email, newPassword, accessToken)` - Changes password

**Authenticate** - Handles authentication flows
- `Password(email, clientID, password)` - Traditional auth
- `MagicCode(email, clientID, magicCode)` - Passwordless auth
- `RequestMagicCode(email)` - Sends magic code email
- `Acknowledge(authenticated)` - **Must be called after every successful authentication**
- `ValidateAccessToken(accessToken)` - Returns JWT claims
- `Renew(email, refreshToken)` - Refreshes access token
- `Revoke(email, accessToken, clientID)` - Revokes session

### HTTP Helpers

Internal functions `doPost`, `doPut`, `doDelete` handle HTTP operations:
- Accept `context.Context` for cancellation/timeouts
- Return structured errors via `JsonError` type
- Auto-marshal/unmarshal JSON payloads

## Key Patterns

### Error Handling
- Status codes >= 300 decode `JsonError` from response body
- Errors formatted as: `"{Title} - {Detail}"`
- Check at line 33-38 in `do_post.go` and `do_delete.go`: the `if jsonError != nil` check is always true (pointer is never nil), but this matches the existing pattern

### Request Payloads
- Use anonymous structs for request bodies to avoid exposing internal types
- JSON tags follow camelCase convention (`clientId`, `accessToken`, etc.)

### Testing Strategy
- Tests use UUID-generated unique emails: `{uuid}@bulwark.io`
- Email verification uses MailHog client (github.com/latebit-io/go-hog)
- Tests follow full end-to-end flows
- `baseUri` and `mailHogUri` constants defined in `account_test.go`

### Context Usage
All public methods accept `context.Context` as first parameter for proper cancellation and timeout handling.

## Important Notes

1. **Acknowledge is required**: After any authentication method returns `Authenticated`, you must call `Acknowledge(authenticated)` to notify the server
2. **clientID parameter**: Represents the device/client identifier, required for Password and MagicCode authentication
3. **Test dependencies**: Tests require live bulwark-auth and MailHog instances
4. **Single package**: Entire library is in package `bulwark` with no subpackages
