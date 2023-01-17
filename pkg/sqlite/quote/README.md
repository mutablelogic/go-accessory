
# Query Quoting

This package provides functions for quoting and unquoting strings for use in SQL queries, taking into account all sqlite quoting rules and reserved keywords.

## Quoting Identifiers

In order to quote an identifier, use the `QuoteIdentifier` function:

```go
    import "github.com/mutablelogic/go-accessory/pkg/sqlite/quote"

    func main() {
        // Quote an identifier
        fmt.Println(quote.QuoteIdentifier("foo"))
    }
```

This will ensure reserved words are quoted, and that the identifier is quoted if it contains spaces or special characters. Where several identifers need to be quoted as a list separated by commas, use the `QuoteIdentifiers` function:

```go
    import "github.com/mutablelogic/go-accessory/pkg/sqlite/quote"

    func main() {
        // Quote a list of identifiers
        fmt.Println(quote.QuoteIdentifiers("foo", "bar")
    }
```

## Quoting TEXT types

In order to quote a string for use in a TEXT type, use the `Quote` and `DoubleQuote` functions:

```go
    import "github.com/mutablelogic/go-accessory/pkg/sqlite/quote"

    func main() {
        // Quote a string 'foo'
        fmt.Println(quote.Quote("foo"))

        // Quote a string "foo"
        fmt.Println(quote.DoubleQuote("foo"))

}
```

## Reserved Words

You can use the `IsReservedWord` function to test if a string is a reserved word:

```go
    import "github.com/mutablelogic/go-accessory/pkg/sqlite/quote"

    func main() {
        // Test if a string is a reserved word, SELECT => true
        fmt.Println(quote.IsReservedWord("SELECT"))
    }
```

