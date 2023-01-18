# go-accessory

**accessory** is a database client for communication with supported database servers,
and some patterns for using databases in a multi-threaded and multi-host environment.

It provides a simple binding with go structures. 

## Database Connections

There are currently two implementations of a database connection. They are:

  * MongoDB
  * sqlite

Both are "in development" and require testing. In order to create a direct
database connection, use one of these:

|Import|Constructor|Description|
|------|-----------|-----------|
|`github.com/mutablelogic/go-accessory/pkg/mongodb`|`mongodb.New(ctx, url, opts...)`| Create a MongoDB connection to a remote server |
|`github.com/mutablelogic/go-accessory/pkg/sqlite`|`sqlite.New(ctx, url, opts...)`| Open a sqlite connection to file or in-memory database |

In either case, the `url` is a string that is parsed by the database driver, which should be of scheme `mongodb://`,  `mongodb+srv://`, `file://` or `sqlite://`. The `opts` are a list of options that are passed to the database driver (see the documentation for the driver for details).

## Mapping Collections to Structures

TODO

## Database Operations

TODO

## Connection Pools

TODO 

## Task Queues

TODO


