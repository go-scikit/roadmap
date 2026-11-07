package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// plant writes files under a fresh directory and returns it.
func plant(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// repoFiles is a repository with one function of every class, a harness
// module, private packages and external scripts.
var repoFiles = map[string]string{
	"go.mod": "module example.com/r\n\ngo 1.27\n",
	"r.go": `package r

import "errors"

// Ref is timed against the reference by perf_test.go.
func Ref(x float64) float64 { return x }

// Noted is timed by a benchmark whose comment names SciPy.
func Noted(x float64) float64 { return x }

// Bench is timed by a plain benchmark.
func Bench(x float64) float64 { return x }

// Untimed has a test but no timing.
func Untimed(x float64) float64 { return x }

// Harness is timed only by the harness module.
func Harness(x float64) float64 { return x }

// Gen is generic; its uses resolve to the origin.
func Gen[T any](x T) T { return x }

// T has methods.
type T struct{}

// M is a method.
func (T) M() float64 { return 0 }

// P has a pointer receiver.
func (*T) P() float64 { return 0 }

// Error is plumbing, not in the denominator.
func (T) Error() string { return "" }

func (T) private() {}

// A is an alias, whose methods are T's.
type A = T

// F is a named func type: no methods.
type F func()

var Err = errors.New("x")

func helper() {}
`,
	"internal/k/k.go":    "package k\n\nfunc K() {}\n",
	"cmd/tool/main.go":   "package main\n\nfunc Exported() {}\n\nfunc main() {}\n",
	"r_internal_test.go": "package r\n\nimport \"testing\"\n\nfunc TestUntimed(t *testing.T) { Untimed(1); helper(); T{}.private() }\n",
	"perf_test.go": `package r_test

import (
	"testing"
	"time"

	"example.com/r"
)

const program = "import numpy as np"

// cases are read by the timing loop.
var cases = []func(float64) float64{r.Ref}

func TestPerf(t *testing.T) {
	t0 := time.Now()
	for _, c := range cases {
		c(1)
	}
	_ = time.Since(t0)
	_ = program
}
`,
	"bench_test.go": `package r_test

import (
	"testing"

	"example.com/r"
)

// BenchmarkNoted is mirrored by a SciPy script.
func BenchmarkNoted(b *testing.B) {
	for range b.N {
		r.Noted(1)
	}
}
`,
	"plain_test.go": `package r_test

import (
	"errors"
	"testing"

	"example.com/r"
)

var table = []float64{1}

func run(b *testing.B, n *int, k int) {
	var anon interface{ M() float64 } = r.T{}
	_ = anon.M()
	_ = errors.Unwrap(r.Err).Error
	_ = r.Bench(table[0]) + r.Gen(1.0) + r.A{}.M() + (&r.T{}).P()
}

func BenchmarkPlain(b *testing.B) { run(b, nil, 0) }

func TestNotTimed(t *testing.T) { r.Untimed(1) }
`,
	"bench/go.mod": "module example.com/r/bench\n\ngo 1.27\n\nrequire example.com/r v0.0.0\n\nreplace example.com/r => ../\n",
	"bench/b.go":   "package bench\n\n// Server is the reference process.\nconst Server = \"bench/scipy_server.py\"\n",
	"bench/b_test.go": `package bench

import (
	"testing"
	"time"

	"example.com/r"
)

func TestHarness(t *testing.T) {
	_ = time.Now()
	r.Harness(1)
}
`,
	"scripts/ref.py":    "import time\nt = time.perf_counter()\n",
	"scripts/a.py":      "import time\ntime.perf_counter()\n",
	"scripts/b.sh":      "python -m timeit 'x'\n",
	"scripts/c.c":       "clock_gettime(0, 0);\n",
	"scripts/notes.sh":  "echo nothing timed\n",
	"scripts/notes.txt": "perf_counter in a text file is not a script\n",
	"testdata/x/go.mod": "module skipped\n",
	".hidden/go.mod":    "module skipped\n",
	"vendor/v/go.mod":   "module skipped\n",
}

func TestInventory(t *testing.T) {
	root := plant(t, repoFiles)
	r, err := inventory(root)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]Class{
		"example.com/r.Ref":     Ref,
		"example.com/r.Harness": Ref,
		"example.com/r.Noted":   Noted,
		"example.com/r.Bench":   Bench,
		"example.com/r.Gen":     Bench,
		"example.com/r.T.M":     Bench,
		"example.com/r.T.P":     Bench,
		"example.com/r.Untimed": None,
	}
	for k, c := range want {
		if got, ok := r.Funcs[k]; !ok || got != c {
			t.Errorf("%s: %v (declared %v), want %v", k, got, ok, c)
		}
	}
	if len(r.Funcs) != len(want) {
		t.Errorf("funcs = %v, want exactly %v", r.Funcs, want)
	}
	if r.Modules != 2 || r.ExternalHarness != 4 || len(r.ExternalExamples) != 3 {
		t.Errorf("modules %d, external %d %v", r.Modules, r.ExternalHarness, r.ExternalExamples)
	}
	if r.Why["example.com/r.Ref"] != "perf_test.go" || r.Why["example.com/r.Harness"] != filepath.Join("bench", "b_test.go") {
		t.Errorf("why = %v", r.Why)
	}
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
	root := plant(t, repoFiles)
	out, err := run(t, "--list", root)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{
		"| " + filepath.Base(root) + " | 2 | 8 | 2 | 1 | 4 | 1 | 4 |",
		"| **total** | | 8 | 2 | 1 | 4 | 1 | 4 |",
		"file\tperf_test.go\tref 1\tnoted 0\tbench 0\n",
		"none\texample.com/r.Untimed\t\n",
		"external\t",
	} {
		if !strings.Contains(out, s) {
			t.Errorf("output lacks %q:\n%s", s, out)
		}
	}
	out, err = run(t, root)
	if err != nil || strings.Contains(out, "## ") {
		t.Errorf("without --list: %v\n%s", err, out)
	}
}

func TestErrors(t *testing.T) {
	if _, err := run(t, filepath.Join(t.TempDir(), "absent")); err == nil {
		t.Error("a missing directory was counted")
	}
	broken := plant(t, map[string]string{
		"go.mod": "module example.com/b\n\ngo 1.27\n",
		"b.go":   "package b\n\nfunc F() int { return \"no\" }\n",
	})
	if _, err := run(t, broken); err == nil || !strings.Contains(err.Error(), "package errors") {
		t.Errorf("a package that does not type-check: %v", err)
	}
	saved := load
	defer func() { load = saved }()
	load = func(*packages.Config, ...string) ([]*packages.Package, error) { return nil, errors.New("boom") }
	if _, err := run(t, plant(t, repoFiles)); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Errorf("a loader failure: %v", err)
	}
}

func TestUnreadableScript(t *testing.T) {
	root := plant(t, map[string]string{"x.py": "perf_counter"})
	p := filepath.Join(root, "x.py")
	if err := os.Chmod(p, 0); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(p, 0o644)
	if _, err := os.ReadFile(p); err == nil {
		t.Skip("running with privileges that ignore file modes")
	}
	if _, err := inventory(root); err == nil {
		t.Error("an unreadable script was skipped")
	}
}

func TestClassString(t *testing.T) {
	if Ref.String() != "ref" || None.String() != "none" || Noted.String() != "noted" || Bench.String() != "bench" {
		t.Error("class names")
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
	os.Args = []string{"perfinventory", plant(t, map[string]string{"README": "no module"})}
	main()
	if code != -1 {
		t.Fatalf("exit %d on an empty repository", code)
	}
	os.Args = []string{"perfinventory"}
	main()
	if code != 1 {
		t.Fatalf("exit %d without an argument, want 1", code)
	}
}
