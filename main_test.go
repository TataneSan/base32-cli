package main

import (
	"encoding/base32"
	"strings"
	"testing"
)

func TestEncodeStd(t *testing.T) {
	enc, err := buildEncoding("std")
	if err != nil {
		t.Fatal(err)
	}
	got := encodeText(enc, "hello")
	want := "NBSWY3DP"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestEncodeHex(t *testing.T) {
	enc, err := buildEncoding("hex")
	if err != nil {
		t.Fatal(err)
	}
	got := encodeText(enc, "hello")
	want := base32.HexEncoding.EncodeToString([]byte("hello"))
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestEncodeNoPad(t *testing.T) {
	enc, err := buildEncoding("nopad")
	if err != nil {
		t.Fatal(err)
	}
	got := encodeText(enc, "hi")
	if strings.Contains(got, "=") {
		t.Fatalf("nopad output contains padding: %q", got)
	}
}

func TestEncodeCrockford(t *testing.T) {
	enc, err := buildEncoding("crockford")
	if err != nil {
		t.Fatal(err)
	}
	got := encodeText(enc, "hello")
	if strings.ContainsAny(got, "ILOUilou=") {
		t.Fatalf("crockford output contains excluded characters: %q", got)
	}
	// Round trip
	back, err := decodeText(enc, "crockford", got)
	if err != nil {
		t.Fatal(err)
	}
	if back != "hello" {
		t.Fatalf("round trip failed: got %q", back)
	}
}

func TestDecodeCrockfordLowercase(t *testing.T) {
	enc, err := buildEncoding("crockford")
	if err != nil {
		t.Fatal(err)
	}
	upper := encodeText(enc, "hello")
	lower := strings.ToLower(upper)
	back, err := decodeText(enc, "crockford", lower)
	if err != nil {
		t.Fatal(err)
	}
	if back != "hello" {
		t.Fatalf("got %q want %q", back, "hello")
	}
}

func TestRoundTripAllVariants(t *testing.T) {
	for _, v := range []string{"std", "hex", "nopad", "crockford"} {
		enc, err := buildEncoding(v)
		if err != nil {
			t.Fatal(err)
		}
		in := "The quick brown fox jumps over 13 lazy dogs!"
		out := encodeText(enc, in)
		back, err := decodeText(enc, v, out)
		if err != nil {
			t.Fatalf("variant %s decode: %v", v, err)
		}
		if back != in {
			t.Fatalf("variant %s round trip: got %q want %q", v, back, in)
		}
	}
}

func TestUnknownVariant(t *testing.T) {
	if _, err := buildEncoding("bogus"); err == nil {
		t.Fatal("expected error for unknown variant")
	}
}

func TestDecodeInvalidBase32(t *testing.T) {
	enc, _ := buildEncoding("std")
	if _, err := decodeText(enc, "std", "!!!not-base32!!!"); err == nil {
		t.Fatal("expected decode error")
	}
}

func TestNormalizeCrockford(t *testing.T) {
	got := normalizeCrockford("iLo")
	if got != "110" {
		t.Fatalf("got %q want %q", got, "110")
	}
	got = normalizeCrockford("abc-def")
	if strings.ContainsAny(got, "- ") {
		t.Fatalf("separator kept: %q", got)
	}
}
