# elido — Go SDK

Sync HTTP client for the [Elido](https://elido.app) URL-shortener API.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/elidoapp/elido-go/elido"
)

func main() {
    c, err := elido.New(elido.Options{
        BaseURL: "https://api.elido.app",
        Token:   "elido_…",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    link, err := c.Links.Create(context.Background(), 42, elido.CreateLinkInput{
        DomainID:       1,
        DestinationURL: "https://example.com/landing",
        Tags:           []string{"spring-2026"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(link.Slug)
}
```

## Auth

In order of precedence:

| Field    | Header                    | Use case            |
| -------- | ------------------------- | ------------------- |
| `Token`  | `Authorization: Bearer …` | Production scripts  |
| `Cookie` | `Cookie: …`               | Server-side proxies |

## Multi-host setups

If you front Elido with separate hostnames per service, provide each URL:

```go
c, _ := elido.New(elido.Options{
    BaseURL:          "https://api.elido.app",
    AnalyticsBaseURL: "https://analytics.elido.app",
    BillingBaseURL:   "https://billing.elido.app",
    Token:            "elido_…",
})
```

## Errors

`*elido.APIError` carries `Status`, `Message`, `Body`. Use `elido.IsAPIError(err)` to discriminate from transport-level errors:

```go
if err := c.Workspaces.Get(ctx, 99); err != nil {
    var ae *elido.APIError
    if errors.As(err, &ae) && ae.Status == 404 {
        // soft-handle
    } else {
        return err
    }
}
```

## Install

```bash
go get github.com/elidoapp/elido-go@latest
```

Import the client package as `github.com/elidoapp/elido-go/elido`. Requires Go 1.22+.

## Versioning

This SDK is hand-written against Elido's public OpenAPI spec (<https://elido.app/openapi.json>). Adding fields to the wire shape stays compatible (Go decoder ignores unknown JSON keys); breaking changes bump the SDK minor.

## Links

- Documentation: <https://elido.app/docs>
- Issues: <https://github.com/elidoapp/elido-go/issues>
- Changelog: [CHANGELOG.md](./CHANGELOG.md)

## License

[MIT](./LICENSE) (c) Elido. The Elido service itself is a separate, hosted
product; this repository contains only the client library.
