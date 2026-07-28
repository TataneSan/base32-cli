# base32-cli

A CLI encoder/decoder for Base32 (RFC 4648) with several alphabet variants:
standard, extended-hex, unpadded, and Crockford Base32. Single static Go
binary, no dependencies.

## Features

- `encode` / `decode` with RFC 4648 standard alphabet (A-Z, 2-7, `=` padding)
- `-variant hex` for the RFC 4648 "Extended Hex" alphabet (0-9, A-V)
- `-variant nopad` for standard Base32 without padding
- `-variant crockford` for Crockford Base32 (no `I L O U`, decodes
  case-insensitively with `I/L->1`, `O->0`, ignores `-`, `_` and spaces)
- `batch` mode for line-by-line processing of a file or stdin
- JSON output (`-json`) on every command
- Reads text from an argument or from stdin (`-`)
- Zero dependencies

## Install

```sh
git clone https://github.com/TataneSan/base32-cli.git
cd base32-cli
go build -o base32-cli .
sudo install -m 0755 base32-cli /usr/local/bin/
```

## Usage

```
base32-cli <command> [ARGS] [flags]

COMMANDS:
    encode TEXT        Encode text to Base32
    decode TEXT        Decode Base32 to text
    batch FILE         Encode each line of FILE ("-" for stdin)
    batch -d FILE      Decode each line of FILE ("-" for stdin)
    version            Show version

Variants (flag -variant): std | hex | nopad | crockford
```

### Encode / decode

```sh
$ base32-cli encode "hello"
NBSWY3DP

$ base32-cli decode "NBSWY3DP"
hello

$ base32-cli encode "hello" -variant hex
D1MORR3P

$ base32-cli encode "hello world" -variant nopad
NBSWY3DPEB3W64TMMQ

$ base32-cli encode "hello" -variant crockford
D1JPRV3F
```

### Via stdin

```sh
$ echo -n "hello" | base32-cli encode -
NBSWY3DP
```

### Batch mode

```sh
$ printf "foo\nbar\n" | base32-cli batch -
MZXW6
MZXHO

$ echo "NBSWY3DP" | base32-cli batch -d -
hello
```

### JSON output

```sh
$ base32-cli encode "hello" -json
{
  "mode": "encode",
  "variant": "std",
  "input": "hello",
  "output": "NBSWY3DP"
}
```

## License

MIT
