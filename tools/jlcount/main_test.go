package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const protoH = `// comment XX(not_a_builtin,"nope")
#define JL_BUILTIN_FUNCTIONS(XX) \
    XX(_apply_iterate,"_apply_iterate") \
    XX(is,"===") \
    XX(setfield,"setfield!") \

#define DECLARE_BUILTIN(cname,jlname) \
    XX(after_the_body,"after")
`

const intrinsicsH = `#define INTRINSICS \
    /*  wrap and unwrap */ \
    ADD_I(bitcast, 2) \
    ALIAS(add_float_fast, add_float) \
    ADD_HIDDEN(cglobal_auto, 1)

enum intrinsic {
#define ADD_I(func, nargs) func,
`

func TestParseBuiltins(t *testing.T) {
	got := parseBuiltins(protoH)
	want := []Builtin{{"_apply_iterate", "_apply_iterate"}, {"is", "==="}, {"setfield", "setfield!"}}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestParseIntrinsics(t *testing.T) {
	got := parseIntrinsics(intrinsicsH)
	want := []Intrinsic{{"ADD_I", "bitcast"}, {"ALIAS", "add_float_fast"}, {"ADD_HIDDEN", "cglobal_auto"}}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestCountLines(t *testing.T) {
	for in, want := range map[string]int{"": 0, "a": 1, "a\n": 1, "a\nb": 2, "\n\n": 2} {
		if got := countLines([]byte(in)); got != want {
			t.Errorf("countLines(%q) = %d, want %d", in, got, want)
		}
	}
}

func tree(t *testing.T, proto string) string {
	t.Helper()
	root := t.TempDir()
	src := filepath.Join(root, "src")
	if err := os.Mkdir(src, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{"builtin_proto.h": proto, "intrinsics.h": intrinsicsH}
	for _, f := range coreFiles {
		files[f] = "one\ntwo\n"
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(src, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestRootCmd(t *testing.T) {
	var out bytes.Buffer
	cmd := RootCmd()
	cmd.SetArgs([]string{tree(t, protoH)})
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"builtins 3\n", "  ===\tis\n",
		"intrinsics 3 (ADD_I 1, ALIAS 1, ADD_HIDDEN 1)\n",
		"  subtype.c\t2\n", "core total 15 files 30 lines\n",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
}

// An empty table is a reading error, never a count of zero.
func TestRootCmdNoBuiltins(t *testing.T) {
	var out bytes.Buffer
	cmd := RootCmd()
	cmd.SetArgs([]string{tree(t, "// nothing here\n")})
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "no JL_BUILTIN_FUNCTIONS entry") {
		t.Fatalf("err = %v, want a no-entry error", err)
	}
}

// Each file the report reads is required: removing or emptying any one of
// them is an error.
func TestRootCmdBrokenTree(t *testing.T) {
	for name, edit := range map[string]func(src string) error{
		"no intrinsics.h": func(src string) error { return os.Remove(filepath.Join(src, "intrinsics.h")) },
		"empty INTRINSICS": func(src string) error {
			return os.WriteFile(filepath.Join(src, "intrinsics.h"), []byte("enum x {};\n"), 0o644)
		},
		"no core file":       func(src string) error { return os.Remove(filepath.Join(src, "gf.c")) },
		"no builtin_proto.h": func(src string) error { return os.Remove(filepath.Join(src, "builtin_proto.h")) },
	} {
		t.Run(name, func(t *testing.T) {
			root := tree(t, protoH)
			if err := edit(filepath.Join(root, "src")); err != nil {
				t.Fatal(err)
			}
			cmd := RootCmd()
			cmd.SetArgs([]string{root})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			if err := cmd.Execute(); err == nil {
				t.Fatal("want an error")
			}
		})
	}
}

func TestMain_(t *testing.T) {
	saved, savedExit, savedOut := os.Args, osExit, os.Stdout
	defer func() { os.Args, osExit, os.Stdout = saved, savedExit, savedOut }()
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer devnull.Close()
	os.Stdout = devnull
	code := -1
	osExit = func(c int) { code = c }

	os.Args = []string{"jlcount", tree(t, protoH)}
	main()
	if code != -1 {
		t.Fatalf("exit %d on a good tree", code)
	}
	os.Args = []string{"jlcount"}
	main()
	if code != 1 {
		t.Fatalf("exit %d without an argument, want 1", code)
	}
}

func TestRootCmdMissingTree(t *testing.T) {
	cmd := RootCmd()
	cmd.SetArgs([]string{filepath.Join(t.TempDir(), "absent")})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want an error for a missing tree")
	}
}
