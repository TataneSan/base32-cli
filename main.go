package main

import (
	"encoding/base32"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type result struct {
	Action  string `json:"action"`
	Input   string `json:"input"`
	Output  string `json:"output"`
	Variant string `json:"variant"`
}

func main() {
	var (
		decode  = flag.Bool("d", false, "decode base32")
		hex     = flag.Bool("hex", false, "use base32hex variant (RFC 4648)")
		nopad   = flag.Bool("nopad", false, "no padding")
		js      = flag.Bool("json", false, "output JSON")
	)
	flag.Parse()

	args := flag.Args()
	var data []byte
	var err error

	if len(args) > 0 {
		data = []byte(strings.Join(args, " "))
	} else {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		data = []byte(strings.TrimRight(string(data), "\r\n"))
	}

	input := string(data)
	var output string
	var enc *base32.Encoding

	if *hex {
		if *nopad {
			enc = base32.HexEncoding.WithPadding(base32.NoPadding)
		} else {
			enc = base32.HexEncoding
		}
	} else {
		if *nopad {
			enc = base32.StdEncoding.WithPadding(base32.NoPadding)
		} else {
			enc = base32.StdEncoding
		}
	}

	action := "encode"
	if *decode {
		action = "decode"
		out, err := enc.DecodeString(strings.ToUpper(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
			os.Exit(1)
		}
		output = string(out)
	} else {
		output = enc.EncodeToString(data)
	}

	variant := "std"
	if *hex {
		variant = "hex"
	}

	if *js {
		r := result{Action: action, Input: input, Output: output, Variant: variant}
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
	} else {
		fmt.Println(output)
	}
}
