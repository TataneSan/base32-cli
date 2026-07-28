package main

import (
	"bufio"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

const version = "1.0.0"

const usage = `base32-cli - Base32 encoder/decoder (RFC 4648)

USAGE:
    base32-cli <command> [ARGS] [flags]

COMMANDS:
    encode TEXT        Encode text to Base32
    decode TEXT        Decode Base32 to text
    batch FILE         Encode each line of FILE ("-" for stdin)
    batch -d FILE      Decode each line of FILE ("-" for stdin)
    version            Show version

Variants (flag -variant):
    std       RFC 4648 standard alphabet (default): A-Z 2-7, padding '='
    hex       RFC 4648 "Extended Hex" alphabet: 0-9 A-V, padding '='
    nopad     Standard alphabet without padding
    crockford Crockford Base32 alphabet 0-9 A-H J-K M-N P-T V-Z (no padding,
              letters I L O U excluded; decoding accepts lowercase and
              substitutes I/L->1, O->0)
Flags:
    -d          Decode mode (used with batch)
    -variant V  Alphabet variant (std, hex, nopad, crockford)
    -json       JSON output
    -h, --help  Show help

EXAMPLES:
    base32-cli encode "hello"
    base32-cli decode "NBSWY3DP"
    base32-cli encode "hello" -variant hex
    echo "NBSWY3DP" | base32-cli batch -d -
`

type resultPayload struct {
	Mode    string `json:"mode"`
	Variant string `json:"variant"`
	Input   string `json:"input"`
	Output  string `json:"output"`
}

type batchItem struct {
	Line   int    `json:"line"`
	Input  string `json:"input"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func buildEncoding(variant string) (*base32.Encoding, error) {
	switch strings.ToLower(variant) {
	case "std", "standard":
		return base32.StdEncoding, nil
	case "hex":
		return base32.HexEncoding, nil
	case "nopad":
		return base32.StdEncoding.WithPadding(base32.NoPadding), nil
	case "crockford":
		return base32.NewEncoding(crockfordAlphabet).WithPadding(base32.NoPadding), nil
	default:
		return nil, fmt.Errorf("unknown variant %q (expected std, hex, nopad or crockford)", variant)
	}
}

func normalizeCrockford(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case 'i', 'I', 'l', 'L':
			b.WriteByte('1')
		case 'o', 'O':
			b.WriteByte('0')
		case 'u', 'U':
			// U is not part of the Crockford alphabet; skip it silently
		case ' ', '-', '_':
			// common visual separators, ignored
		default:
			if r >= 'a' && r <= 'z' {
				b.WriteRune(r - 32)
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

func encodeText(enc *base32.Encoding, text string) string {
	return enc.EncodeToString([]byte(text))
}

func decodeText(enc *base32.Encoding, variant, text string) (string, error) {
	data := strings.TrimSpace(text)
	if strings.ToLower(variant) == "crockford" {
		data = normalizeCrockford(data)
	}
	out, err := enc.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func readInput(arg string) (string, error) {
	if arg == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return strings.TrimRight(string(data), "\r\n"), nil
	}
	return arg, nil
}

func printResult(payload resultPayload, asJSON bool) {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(payload)
		return
	}
	fmt.Println(payload.Output)
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "base32-cli: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Print(usage)
		os.Exit(1)
	}

	cmd := args[0]
	rest := args[1:]

	var variant string
	var decodeFlag, jsonFlag bool
	var positional []string

	for i := 0; i < len(rest); i++ {
		a := rest[i]
		switch {
		case a == "-d":
			decodeFlag = true
		case a == "-json":
			jsonFlag = true
		case a == "-variant":
			if i+1 >= len(rest) {
				fatal("flag -variant requires a value")
			}
			i++
			variant = rest[i]
		case a == "-h" || a == "--help":
			fmt.Print(usage)
			os.Exit(0)
		case strings.HasPrefix(a, "-") && len(a) > 1 && a != "-":
			fatal("unknown flag %s", a)
		default:
			positional = append(positional, a)
		}
	}

	if variant == "" {
		variant = "std"
	}

	switch cmd {
	case "version":
		fmt.Printf("base32-cli %s\n", version)
		os.Exit(0)
	case "encode", "decode":
		if len(positional) < 1 {
			fatal("missing TEXT argument for %s", cmd)
		}
		enc, err := buildEncoding(variant)
		if err != nil {
			fatal("%v", err)
		}
		text, err := readInput(positional[0])
		if err != nil {
			fatal("reading input: %v", err)
		}
		if cmd == "encode" {
			out := encodeText(enc, text)
			printResult(resultPayload{Mode: "encode", Variant: variant, Input: text, Output: out}, jsonFlag)
		} else {
			out, err := decodeText(enc, variant, text)
			if err != nil {
				fatal("decode: %v", err)
			}
			printResult(resultPayload{Mode: "decode", Variant: variant, Input: text, Output: out}, jsonFlag)
		}
	case "batch":
		if len(positional) < 1 {
			fatal("missing FILE argument for batch")
		}
		enc, err := buildEncoding(variant)
		if err != nil {
			fatal("%v", err)
		}
		var reader io.Reader
		if positional[0] == "-" {
			reader = os.Stdin
		} else {
			f, err := os.Open(positional[0])
			if err != nil {
				fatal("open %s: %v", positional[0], err)
			}
			defer f.Close()
			reader = f
		}
		mode := "encode"
		if decodeFlag {
			mode = "decode"
		}
		scanner := bufio.NewScanner(reader)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
		items := make([]batchItem, 0)
		line := 0
		for scanner.Scan() {
			line++
			in := scanner.Text()
			item := batchItem{Line: line, Input: in}
			if decodeFlag {
				out, err := decodeText(enc, variant, in)
				if err != nil {
					item.Error = err.Error()
				} else {
					item.Output = out
				}
			} else {
				item.Output = encodeText(enc, in)
			}
			if !jsonFlag {
				if item.Error != "" {
					fmt.Printf("line %d: error: %v\n", line, item.Error)
				} else {
					fmt.Println(item.Output)
				}
			}
			items = append(items, item)
		}
		if err := scanner.Err(); err != nil {
			fatal("read: %v", err)
		}
		if jsonFlag {
			out := struct {
				Mode    string      `json:"mode"`
				Variant string      `json:"variant"`
				Count   int         `json:"count"`
				Items   []batchItem `json:"items"`
			}{Mode: mode, Variant: variant, Count: len(items), Items: items}
			e := json.NewEncoder(os.Stdout)
			e.SetIndent("", "  ")
			e.Encode(out)
		}
	default:
		fmt.Print(usage)
		fatal("unknown command %q", cmd)
	}
}
