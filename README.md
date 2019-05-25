# Example HTTP API

This is a rough framework for how I've ended up writing JSON REST APIs in Go.
I've put this together as I spent quite a long time trying to figure out how I should structure code, mostly between the frontend and backend packages. I wanted to have an abstract database access interface to allow creating lightweight unit tests for API handlers (meaning I wanted mock data and didn't want to require a full RDBMS) so the backend needed to be easily swappable.

It's structured as follows:

```
api/              The package implementing the HTTP endpoints
cmd/
    apiserver/    The binary running the HTTP server
database/         The package with types and database interface. This package could equally be extracted to its own repo
    mock/         A mock backend returning sample data
    postgres/     A PostgreSQL backend
```

I came up with the `Storage` interface in the `database` package

```go
type Storage interface {
    // Database access methods
}
```

which each of the `mock` and `postgres` packages implement.

The `api` package creates a mux with its own routes, which are returned by the `Routes()` method. This allows versioning the API by "mounting" these routes in the API server at different paths.

The `apiserver` binary mounts these routes at `/` by default:

```go
apisvc := api.NewService(db)
r := chi.NewRouter()
r.Mount("/", apisvc.Routes())
```

So the `/users` endpoint defined in the `api` package appears as `/users`. If you created a new version of the API you could do:

```go
import (
    apiv1 "github.com/jamesog/example-http-api/api/v1"
    apiv2 "github.com/jamesog/example-http-api/api/v2"
)

apiv1 := apiv1.NewService(db)
apiv2 := apiv2.NewService(db)
r := chi.NewRouter()
r.Mount("/", apiv1.Routes())
r.Mount("/v2", apiv2.Routes())
```

This would provide `/users` and `/v2/users`.
