# base32-cli

Fast Base32 encoder/decoder for the terminal.

Supports:
- Standard Base32 (RFC 4648)
- Base32hex variant (RFC 4648 extended hex alphabet)
- With/without padding
- Input from stdin or arguments
- Optional JSON output

## Install

```bash
cd base32-cli
go build -o base32-cli .
sudo mv base32-cli /usr/local/bin/
```

Or run directly:

```bash
go run . "hello world"
```

## Usage

Encode from arguments:
```bash
base32-cli "hello world"
# NBSWY3DPEB3W64TMMQ======
```

Encode from stdin:
```bash
echo "hello world" | base32-cli
```

Decode:
```bash
base32-cli -d "NBSWY3DPEB3W64TMMQ======"
# hello world
```

Use base32hex variant:
```bash
base32-cli -hex "hello world"
# D1IMOR3F41RMUSFJCC======
```

No padding:
```bash
base32-cli -nopad "hello world"
# NBSWY3DPEB3W64TMMQ
```

JSON output:
```bash
base32-cli -json "hello world"
# {"action":"encode","input":"hello world","output":"NBSWY3DPEB3W64TMMQ======","variant":"std"}
```

Decode + JSON:
```bash
base32-cli -d -json "NBSWY3DPEB3W64TMMQ======"
```

## Flags

| Flag     | Description                              |
|----------|------------------------------------------|
| `-d`     | Decode instead of encode                 |
| `-hex`   | Use base32hex alphabet (0-9A-V)          |
| `-nopad` | Remove `=` padding                       |
| `-json`  | Output result as JSON (machine-readable) |

## Examples

```bash
# Encode a token for TOTP
base32-cli -nopad -hex "JBSWY3DPEHPK3PXP"

# Decode from a file
base32-cli -d < input.b32 > output.bin

# Pipe through
cat data.txt | base32-cli | base32-cli -d
```

## License

MIT
