# bulwark-auth-guard

Go client library for [bulwark-auth](https://github.com/latebit-io/bulwark-auth) authentication service.

## Installation

```bash
go get github.com/latebit-io/bulwark-auth-guard
```

## Usage

### Creating a Client

```go
import (
    "net/http"
    bulwark "github.com/latebit-io/bulwark-auth-guard"
)

client := &http.Client{}
guard := bulwark.NewGuard("http://localhost:8080", client)
```

### Account Management

#### Create Account
```go
ctx := context.Background()
email := "user@example.com"
password := "securePassword123!"

err := guard.Account.Create(ctx, email, password)
if err != nil {
    log.Fatal(err)
}
// User receives verification email
```

#### Verify Account
```go
verificationToken := "token-from-email"
err := guard.Account.Verify(ctx, email, verificationToken)
if err != nil {
    log.Fatal(err)
}
```

#### Change Password
```go
err := guard.Account.ChangePassword(ctx, email, "newPassword123!", accessToken)
if err != nil {
    log.Fatal(err)
}
```

### Authentication

#### Password Authentication
```go
clientID := "my-app-device-id"
authenticated, err := guard.Authenticate.Password(ctx, email, clientID, password)
if err != nil {
    log.Fatal(err)
}

// Important: Must acknowledge after authentication
err = guard.Authenticate.Acknowledge(ctx, authenticated)
if err != nil {
    log.Fatal(err)
}

// Access tokens
fmt.Println(authenticated.AccessToken)
fmt.Println(authenticated.RefreshToken)
```

#### Magic Code (Passwordless) Authentication
```go
// Request magic code
err := guard.Authenticate.RequestMagicCode(ctx, email)
if err != nil {
    log.Fatal(err)
}
// User receives email with magic code

// Authenticate with code
magicCode := "code-from-email"
authenticated, err := guard.Authenticate.MagicCode(ctx, email, clientID, magicCode)
if err != nil {
    log.Fatal(err)
}

// Important: Must acknowledge after authentication
err = guard.Authenticate.Acknowledge(ctx, authenticated)
if err != nil {
    log.Fatal(err)
}
```

### Token Management

#### Validate Access Token
```go
claims, err := guard.Authenticate.ValidateAccessToken(ctx, authenticated.AccessToken)
if err != nil {
    log.Fatal(err)
}

fmt.Println(claims.Subject)   // User email
fmt.Println(claims.Roles)     // User roles
fmt.Println(claims.ExpiresAt) // Token expiration
fmt.Println(claims.ClientID)  // Client identifier
```

#### Renew Access Token
```go
newAuth, err := guard.Authenticate.Renew(ctx, email, authenticated.RefreshToken)
if err != nil {
    log.Fatal(err)
}

fmt.Println(newAuth.AccessToken)
fmt.Println(newAuth.RefreshToken)
```

#### Revoke Session
```go
err := guard.Authenticate.Revoke(ctx, email, authenticated.AccessToken, clientID)
if err != nil {
    log.Fatal(err)
}
```

## Testing

Tests require running services:
- **bulwark-auth** server at `http://localhost:8080`
- **MailHog** at `http://localhost:8025`

```bash
go test -v
```

## License

MIT License - See LICENSE file for details
