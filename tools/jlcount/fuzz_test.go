package main

import (
	"strings"
	"testing"
)

// FuzzParse checks that the table readers never panic and only return
// entries that their own patterns match.
func FuzzParse(f *testing.F) {
	f.Add(protoH)
	f.Add(intrinsicsH)
	f.Add("#define JL_BUILTIN_FUNCTIONS(XX) \\")
	f.Add("#define INTRINSICS \\\n    ADD_I(x, 1) \\")
	f.Fuzz(func(t *testing.T, text string) {
		for _, b := range parseBuiltins(text) {
			if b.CName == "" || b.JLName == "" || strings.Contains(b.JLName, `"`) {
				t.Fatalf("bad builtin %+v", b)
			}
		}
		for _, x := range parseIntrinsics(text) {
			if x.Name == "" || (x.Kind != "ADD_I" && x.Kind != "ALIAS" && x.Kind != "ADD_HIDDEN") {
				t.Fatalf("bad intrinsic %+v", x)
			}
		}
	})
}
