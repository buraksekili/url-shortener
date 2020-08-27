package shortener

import (
	"testing"
)

func TestParseFlags(t *testing.T) {

	// Help operator, if there is no flag indicated
	op := ParseFlags(nil)
	_, ok := op.(HelpOp)
	if !ok {
		t.Fatalf("got=%T expected=HelpOp", op)
	}

	// Help operator, if args array is empty
	op = ParseFlags([]string{})
	_, ok = op.(HelpOp)
	if !ok {
		t.Fatalf("got=%T expected=HelpOp", op)
	}

	// Help operator, if flag is -h (shorthand for --help)
	op = ParseFlags([]string{"-h"})
	_, ok = op.(HelpOp)
	if !ok {
		t.Fatalf("got=%T expected=HelpOp", op)
	}

	// Help operator, if flag is --help
	op = ParseFlags([]string{"--help"})
	_, ok = op.(HelpOp)
	if !ok {
		t.Fatalf("got=%T expected=HelpOp", op)
	}

	// Entry operator, if flag is -e
	op = ParseFlags([]string{"-e=from-there"})
	v1, ok := op.(EntryOp)
	if !ok {
		t.Fatalf("got=%T expected=EntryOp", op)
	}
	expected1 := EntryOp{From: "from", Target: "there", FileName: "paths.json"}
	if v1 != expected1 {
		t.Fatalf("Target value mismatch, got=%#v expected=%#v", v1, expected1)
	}

	// Entry operator, if flag is -e
	op = ParseFlags([]string{"-e=from-there", "-f=foo.json"})
	v1, ok = op.(EntryOp)
	if !ok {
		t.Fatalf("got=%T expected=EntryOp", op)
	}
	expected1 = EntryOp{From: "from", Target: "there", FileName: "foo.json"}
	if v1 != expected1 {
		t.Fatalf("Target value mismatch, got=%#v expected=%#v", v1, expected1)
	}

	op = ParseFlags([]string{"-e=from-there", "--fname=foo.json"})
	v1, ok = op.(EntryOp)
	if !ok {
		t.Fatalf("got=%T expected=EntryOp", op)
	}
	expected1 = EntryOp{From: "from", Target: "there", FileName: "foo.json"}
	if v1 != expected1 {
		t.Fatalf("Target value mismatch, got=%#v expected=%#v", v1, expected1)
	}

	// Entry operator, if flag is -e
	op = ParseFlags([]string{"-e"})
	_, ok = op.(UnknownOp)
	if !ok {
		t.Fatalf("got=%T expected=UnknownOp", op)
	}

	// If Filename operator detected without filename, return UnknownOp
	op = ParseFlags([]string{"-f"})
	_, ok = op.(UnknownOp)
	if !ok {
		t.Fatalf("got=%T expected=UnknownOp", op)
	}

	// If Filename operator detected without filename, return UnknownOp
	op = ParseFlags([]string{"--fname"})
	_, ok = op.(UnknownOp)
	if !ok {
		t.Fatalf("got=%T expected=UnknownOp", op)
	}

	// Find operator, if flag is string that has no dash(-) as a prefix.
	op = ParseFlags([]string{"foo", "-f=foo.json"})
	v, ok := op.(FindOp)
	if !ok {
		t.Fatalf("got=%T expected=FindOp", op)
	}
	expected := FindOp{target: "foo", filename: "foo.json"}
	if v != expected {
		t.Fatalf("Target value mismatch, got=%#v expected=%#v", v, expected)
	}

	// Find operator, if flag is string that has no dash(-) as a prefix.
	op = ParseFlags([]string{"foo", "--fname=foo"})
	v, ok = op.(FindOp)
	if !ok {
		t.Fatalf("got=%T expected=FindOp", op)
	}
	expected = FindOp{target: "foo", filename: "foo"}
	if v != expected {
		t.Fatalf("Target value mismatch, got=%#v expected=%#v", v, expected)
	}

	// Find operator, if flag is string that has no dash(-) as a prefix.
	op = ParseFlags([]string{"foo"})
	v, ok = op.(FindOp)
	if !ok {
		t.Fatalf("got=%T expected=FindOp", op)
	}
	expected = FindOp{target: "foo", filename: "paths.json"}
	if v != expected {
		t.Fatalf("Target value mismatch, got=%#v expected=%#v", v, expected)
	}

	// Unrecognized flag
	op = ParseFlags([]string{"-x"})
	_, ok = op.(UnknownOp)
	if !ok {
		t.Fatalf("got=%T expected=UnknownOp", op)
	}

}
