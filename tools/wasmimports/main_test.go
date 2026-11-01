package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// module is a hand-assembled wasm binary: it imports env.f and env.memory,
// defines and exports g, and carries a dylink.0 custom section.
var module = []byte{
	0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00, // magic, version
	0x01, 0x04, 0x01, 0x60, 0x00, 0x00, // type: () -> ()
	0x02, 0x18, 0x02, // import, 2 entries
	0x03, 'e', 'n', 'v', 0x01, 'f', 0x00, 0x00, // env.f func type 0
	0x03, 'e', 'n', 'v', 0x06, 'm', 'e', 'm', 'o', 'r', 'y', 0x02, 0x01, 0x01, 0x02, // env.memory min 1 max 2
	0x03, 0x02, 0x01, 0x00, // function: g has type 0
	0x07, 0x05, 0x01, 0x01, 'g', 0x00, 0x01, // export g = func 1
	0x0a, 0x04, 0x01, 0x02, 0x00, 0x0b, // code: empty body
	0x00, 0x0b, 0x08, 'd', 'y', 'l', 'i', 'n', 'k', '.', '0', 'x', 'y', // custom dylink.0 "xy"
}

func write(t *testing.T, b []byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "m.wasm")
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd := RootCmd()
	cmd.SetArgs(args)
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	err := cmd.Execute()
	return out.String(), err
}

func TestReport(t *testing.T) {
	out, err := run(t, "--names", write(t, module))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"valid, 70 bytes\n",
		"imported functions 1\n  env\t1\n    f\n",
		"imported memory env.memory min 1 pages max 2 pages\n",
		"exported functions 1\n",
		"custom section dylink.0\t2 bytes\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

func TestReportWithoutNames(t *testing.T) {
	out, err := run(t, write(t, module))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "    f\n") {
		t.Fatalf("names listed without --names:\n%s", out)
	}
}

// A legacy exception-handling "try" (0x06) is refused, as wazero refuses
// it in the clang module bib/c-cpp.md measures.
func TestInvalid(t *testing.T) {
	// The code section (bytes 51-56) is replaced by a body holding
	// "try (empty) end": 0x06 0x40 0x0b.
	bad := append(append(append([]byte(nil), module[:51]...),
		0x0a, 0x07, 0x01, 0x05, 0x00, 0x06, 0x40, 0x0b, 0x0b), module[57:]...)
	if _, err := run(t, write(t, bad)); err == nil || !strings.Contains(err.Error(), "legacy exception handling") {
		t.Fatalf("err = %v, want the legacy exception-handling refusal", err)
	}
}

func TestMissing(t *testing.T) {
	if _, err := run(t, filepath.Join(t.TempDir(), "absent.wasm")); err == nil {
		t.Fatal("want an error for a missing file")
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

	os.Args = []string{"wasmimports", write(t, module)}
	main()
	if code != -1 {
		t.Fatalf("exit %d on a valid module", code)
	}
	os.Args = []string{"wasmimports"}
	main()
	if code != 1 {
		t.Fatalf("exit %d without an argument, want 1", code)
	}
}
