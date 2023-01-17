# Query Builder

This package provides a query builder for sqlite. There are a small number of primitives and then a number of functions which build on these primitives to provide support for building queries.

It's not really intended you would use this query builder directly, but it would be incorpporated into a higher level packages.

## Primitives

The following primitives exist, which are used to build queries:

| Primitive | Description | Example |
|-----------|-------------|---------|
| `S()` | SELECT statement | `S("id", "name").From("users")` |
| `E()` | Expression | `E("string")` |
| `N()` | Table or column name | `N("column")` |
| `Q()` | Bare query | `Q("SELECT * FROM users")` |

## Query

The `Q(string)` primitive is used to create a bare query. This is useful if you want to use a query which is not supported by the query builder. In order to create an SQL query:

```go
    import (
        . "github.com/mutablelogic/go-accessory/pkg/sqlite/query"
    )

    sql := Q("SELECT * FROM users").Query()
```

Any query can include "flags" which are used to control the behaviour of the query. The `With` function should always be included at the end of an expression. For example,

```go
    import (
        . "github.com/mutablelogic/go-accessory/pkg/sqlite/query"
    )

    // SELECT DISTINCT * FROM users
    sql := S().From(N("users")).With(DISTINCT).Query()
```

## Name

The `N(string)` primitive is used to create a table or column name. This can be used as an argument to a `S()` primitive or `E()` primitive, for example. To create a table, use:

```go
    import (
        . "github.com/mutablelogic/go-accessory/pkg/sqlite/query"
    )

    // CREATE TABLE users (a, b)
    sql := N("users").CreateTable(N("a"), N("b"))
```

You can modify a name primitive with the following modifiers:

| Modifier | Description | Example |
|----------|-------------|---------|
| `N(string).As(string)` | Alias name | `N("column").As("alias")` |
| `N(string).WithSchema(string)` | Schema name | `N("column").WithSchema("main")` |
| `N(string).WithType(string)` | Declared column type | `N("column").WithType("TIMESTAMP")` |

## Expression

The `E(any)` primitive is used to create an expression which can be used in a `S()` primitive. This can be used to create a literal value, a column name, or a function. To create a literal value, use:

```go
    import (
        . "github.com/mutablelogic/go-accessory/pkg/sqlite/query"
    )

    // SELECT 1
    sql := S(E(1)).Query()

    // SELECT name AS uid FROM users
    sql := S(N("users").As("uid")).From("users").Query()
```

