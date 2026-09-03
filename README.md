# lokey

A Redis-like in-memory key-value server written in Go, built from scratch.

Right now `lokey` is a single-threaded TCP server that echoes back whatever a
client sends. The wire format it will speak is documented in
[protocol.md](protocol.md) and summarized below.

## Requirements

- Go 1.26.3 or newer (see [go.mod](go.mod))

## Running

```sh
go run .
```

The server listens on TCP port `1111`.

## Connecting

Any raw TCP client works:

```sh
nc localhost 1111
```

On connect the server sends a welcome banner, then echoes back every line you
send.

## Protocol

`lokey` uses a RESP-style serialization format. Every type is identified by its
first byte and every part is terminated with `\r\n`.

### Simple strings

Prefixed with `+`:

```
+OK\r\n
```

### Errors

Prefixed with `-`, followed by the error message and `\r\n`:

```
-ERR unknown command\r\n
```

Errors are framed like simple strings; the leading byte is what tells a client
the reply represents a failure rather than a value.

### Bulk strings

Prefixed with `$`, followed by the length of the payload in bytes, `\r\n`, the
payload itself, and a trailing `\r\n`. Because the length is given up front,
bulk strings can carry arbitrary binary data:

```
$5\r\nhello\r\n
```

Empty string:

```
$0\r\n\r\n
```

Null (a bulk string with a length of `-1`, carrying no payload):

```
$-1\r\n
```

### Integers

Prefixed with `:`:

```
:1234\r\n
```

### Arrays

Prefixed with `*`, followed by the number of elements and `\r\n`, then each
element encoded in its own format. Elements may be of mixed types.

For example, `["tony", 3000, "stark"]` encodes as:

```
*3\r\n$4\r\ntony\r\n:3000\r\n$5\r\nstark\r\n
```

Written out one part per line:

```
*3\r\n
$4\r\ntony\r\n
:3000\r\n
$5\r\nstark\r\n
```

## Layout

| Path | Description |
| --- | --- |
| [main.go](main.go) | Entry point; prints the banner and starts the server |
| [server/tcpsync.go](server/tcpsync.go) | Single-threaded TCP accept and echo loop |
| [protocol.md](protocol.md) | Wire format notes |

## Status

- [x] TCP server accepting connections
- [x] Welcome banner
- [x] Echo responses
- [ ] Protocol parser and serializer
- [ ] Commands (`GET`, `SET`, ...)
- [ ] Concurrent client handling
