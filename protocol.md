# The lokey wire protocol

`lokey` speaks a RESP-style format: a client sends a request, the server sends
back a reply, and both are encoded the same way.

Two rules cover most of it:

1. **The first byte says what the value is.** `+` simple string, `-` error,
   `:` integer, `$` bulk string, `*` array.
2. **Every part ends with `\r\n`** (carriage return, then line feed — CRLF).

So a parser only ever has to look at one byte to know what it is about to read.

| Prefix | Type | Example |
| --- | --- | --- |
| `+` | Simple string | `+OK\r\n` |
| `-` | Error | `-ERR unknown command\r\n` |
| `:` | Integer | `:1234\r\n` |
| `$` | Bulk string | `$5\r\nhello\r\n` |
| `*` | Array | `*1\r\n$4\r\nPING\r\n` |

## Simple strings

A `+`, the text, then `\r\n`:

```
+OK\r\n
```

The text cannot contain `\r` or `\n`, because those are what end it. Use a bulk
string for anything that might.

## Errors

Identical in shape to a simple string, but prefixed with `-`:

```
-ERR unknown command\r\n
```

The leading byte is the only difference, and it is what tells the client this
reply is a failure rather than a value. Convention is an uppercase error code
first (`ERR`, `WRONGTYPE`), then a human-readable message.

## Integers

A `:`, the digits, then `\r\n`:

```
:1234\r\n
:0\r\n
:-42\r\n
```

Signed, 64-bit.

## Bulk strings

These carry arbitrary bytes, including `\r`, `\n`, and NUL. The trick is that
the **length comes first**, so the parser never has to scan for a terminator —
it reads the count, then blindly takes that many bytes.

```
$5\r\nhello\r\n
│ │   │     └── trailing CRLF
│ │   └── the 5 bytes of payload
│ └── CRLF after the length
└── length in bytes
```

Empty string — length `0`, then an empty payload:

```
$0\r\n\r\n
```

Null — length `-1`, and no payload or trailing CRLF at all. This is how you say
"no value", distinct from the empty string above:

```
$-1\r\n
```

Because the length is authoritative, `$5\r\nhe\r\no\r\n` is a perfectly valid
five-byte string containing a newline.

## Arrays

An `*`, the number of **elements** (not bytes), `\r\n`, then each element
encoded in its own format. Elements may be of mixed types, and an element can
itself be an array.

Encoding `["tony", 3000, "stark"]`:

```
*3\r\n$4\r\ntony\r\n:3000\r\n$5\r\nstark\r\n
```

Split one element per line, which is easier to read but the same bytes:

```
*3\r\n              3 elements follow
$4\r\ntony\r\n      bulk string, 4 bytes
:3000\r\n           integer
$5\r\nstark\r\n     bulk string, 5 bytes
```

Empty array:

```
*0\r\n
```

Null array — same idea as a null bulk string:

```
*-1\r\n
```

Arrays are how **commands** are sent. `SET name tony` goes over the wire as an
array of bulk strings:

```
*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$4\r\ntony\r\n
```

## Notes for the parser

The decoder lives in [core/resp.go](core/resp.go). Each reader has the shape:

```go
func readX(data []byte) (value, int, error)
```

That middle `int` is **how many bytes the value consumed**. It is what makes
arrays possible: after decoding one element you advance the cursor by that
amount to find where the next one begins. Without it there is no way to locate
element two, since elements vary in length and may nest.

```go
pos := 0
for i := 0; i < count; i++ {
    elem, n, err := DecodeOne(data[pos:])
    pos += n   // ← the cursor moves by whatever the element reported
}
```

It also handles pipelining: a client may pack several requests into one packet,
and the byte count is what lets the read loop find the second one.

Every reader is handed a slice that starts at its own first byte, so the count
is a **delta** (bytes consumed), never an absolute offset into the original
buffer. Keep that consistent — mixing the two conventions corrupts offsets
silently.

Two things a reader must not assume:

- **That a terminator exists.** Scanning for `\r` without a bounds check walks
  off the end of the slice and panics on malformed input like `+OK` (no CRLF).
- **That the whole value arrived.** TCP splits wherever it likes, so a read can
  end mid-value. A parser needs to report "incomplete, come back with more"
  separately from "malformed".
