package cli

import (
	"strings"
	"testing"
)

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fatalf("expected %#v, actual: %#v\n", expected, actual)
	}
}

func TestExtractFlag(t *testing.T) {
	args := NewArgs([]string{"hi", "-abcd", "hello", "world"})
	boolFlag, args := args.ExtractFlag("-c", "--colorize", true)
	assertEqual(t, true, boolFlag.IsProvided())
	assertEqual(t, true, boolFlag.Bool())
	assertEqual(t, "-a -b -d hello world", args.String())

	stringFlag, args := args.ExtractFlag("-d", "--description", "TEXT")
	assertEqual(t, true, stringFlag.IsProvided())
	assertEqual(t, "hello", stringFlag.String())
	assertEqual(t, "-a -b world", args.String())
}

func TestExtractFlag_MultipleBool(t *testing.T) {
	args := NewArgs([]string{"hi", "-c", "hello", "--colorize=0", "world"})
	boolFlag, args := args.ExtractFlag("-c", "--colorize", true)
	assertEqual(t, true, boolFlag.IsProvided())
	assertEqual(t, false, boolFlag.Bool())
	assertEqual(t, "hello world", args.String())
}

func TestExtractFlag_MultipleString(t *testing.T) {
	args := NewArgs([]string{"hi", "-d", "hello", "one", "--description=cruel world", "two", "--description", "goodbye", "three"})
	stringFlag, args := args.ExtractFlag("-d", "--description", "TEXT")
	assertEqual(t, true, stringFlag.IsProvided())
	assertEqual(t, "hello, cruel world, goodbye", strings.Join(stringFlag.Strings(), ", "))
	assertEqual(t, "one two three", args.String())
}
