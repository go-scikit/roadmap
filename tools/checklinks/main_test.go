package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path, s string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.md"), "[ok](b/c.md) [anchor](b/c.md#x) [web](https://x.org/y) [self](#top) [mail](mailto:a@b)\n")
	write(t, filepath.Join(dir, "b", "c.md"), "[up](../a.md)\n")
	write(t, filepath.Join(dir, "b", "notes.txt"), "[ignored](nothing.md)\n")
	write(t, filepath.Join(dir, ".git", "x.md"), "[skipped](nothing.md)\n")
	var out strings.Builder
	if code := run([]string{dir}, &out); code != 0 {
		t.Fatalf("clean tree: code %d, %q", code, out.String())
	}
	write(t, filepath.Join(dir, "b", "d.md"), "[gone](missing.md)\n")
	out.Reset()
	if code := run([]string{dir}, &out); code != 1 || !strings.Contains(out.String(), "missing.md") {
		t.Fatalf("broken tree: code %d, %q", code, out.String())
	}
	out.Reset()
	if code := run([]string{filepath.Join(dir, "nope")}, &out); code != 2 {
		t.Fatalf("missing root: code %d", code)
	}
	if err := os.Chmod(filepath.Join(dir, "b", "d.md"), 0); err == nil && os.Getuid() != 0 {
		out.Reset()
		if code := run([]string{dir}, &out); code != 2 {
			t.Fatalf("unreadable file: code %d", code)
		}
	}
}

func TestMain_(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.md"), "[gone](missing.md)\n")
	code := -1
	exit = func(c int) { code = c }
	defer func() { exit = os.Exit }()
	args := os.Args
	os.Args = []string{"checklinks", dir}
	defer func() { os.Args = args }()
	main()
	if code != 1 {
		t.Fatalf("main exit code %d, want 1", code)
	}
}

func TestDefaultRoot(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.md"), "no links\n")
	t.Chdir(dir)
	var out strings.Builder
	if code := run(nil, &out); code != 0 {
		t.Fatalf("code %d", code)
	}
}
