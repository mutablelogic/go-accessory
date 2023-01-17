
# Connection Pool

This is an implementation of a connection pool. In order to create a connection pool, use the `New` function, with the URL for the database connection and optionally a number of client connection options:

```go

import "github.com/mutablelogic/go-accessory/pkg/pool"

func main() {
    var url *url.URL // This is the connection to the database
    var opts []pool.Option // This is a list of options for the pool

    pool := pool.New(context.TODO(), url, opts...)
    if pool == nil {
        panic("Unable to create pool")
    }
    defer pool.Close()

    // Use the connection pool
}
```

The connection pool URL can be of scheme `mongodb://` or `sqlite://` depending on your database. The
options you can pass to the pool are as follows:

| Option | Description | Usage |
|--------|-------------|-------|
| `pool.OptMaxSize(int64)` | Set the maximum number of connections allowed to be pooled |
| `pool.OptDatabase(string)` | Set the default database to use | MongoDB only |
| `pool.OptAttach(*url.URL, string)` | Attach additional databases | sqlite only |
| `pool.OptTimeout(time.Duration)` | Set the connection and operation timeout | MongoDB only |
| `pool.OptCollection(any, string)` | Map a struct prototype to a collection name |
| `pool.OptTrace(trace.Func)` | Trace database operations to a trace function. The trace function signature should be `func(context.Context, time.Duration, error)` |

## Getting a connection from the pool

In order to get a connection for use, use the `Get` function and return it to the pool with the `Put` function. You should always pair a `Get` with a `Put`:

```go
    conn, err := pool.Get()
    if err != nil {
        panic("Unable to get connection")
    }
    defer pool.Put(conn)
```

You should always test the `Get` function for returning `nil`. Typically this will be returned if a connection to the database could not be established or the maximum number of connections have been reached.

## Getting the pool size

The `Size` function returns the current number of connections in the pool:

```go
    size := pool.Size()
```

## Why use a connection pool

A connection pool is used to reduce the overhead of establishing a connection to a database. The connection pool will maintain a pool of connections to the database, and will reuse these connections when a new connection is requested. In addition,

  * It can be used to limit the number of connections to the database, to manage resources;
  * In a multi-threaded environment, it can be used to ensure that only one thread is using a connection at a time.

The ability to close idle connections after a time is not currently implemented.
