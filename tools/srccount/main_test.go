package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	for in, want := range map[string]int{"": 0, "a": 1, "a\n": 1, "a\nb": 2, "\n\n": 2} {
		if got := countLines([]byte(in)); got != want {
			t.Errorf("countLines(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestCountMatches(t *testing.T) {
	if got := countMatches([]byte("x/execute/a.c\nx/compile/b.c\nexecute/execute/c.c\n"), "execute/"); got != 2 {
		t.Fatalf("got %d, want 2 (one per line, not per occurrence)", got)
	}
}

func TestWanted(t *testing.T) {
	for _, c := range []struct {
		name string
		exts []string
		want bool
	}{
		{"vector", nil, true},
		{"a.cpp", []string{".c", ".cpp"}, true},
		{"a.cppm", []string{".cpp"}, false},
		{"Makefile", []string{".c"}, false},
	} {
		if got := wanted(c.name, c.exts); got != c.want {
			t.Errorf("wanted(%q, %v) = %v, want %v", c.name, c.exts, got, c.want)
		}
	}
}

// tree plants files of every kind the walk must accept or skip.
func tree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"a.c":          "int x;\nint y;\n",
		"sub/b.h":      "#define Z 1",
		"sub/vector":   "// libc++ header, no extension\n",
		"notes.txt":    "not C\n",
		".git/HEAD":    "ref: x.c\n",
		".git/obj/x.c": "int planted_in_git;\n",
	}
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(filepath.Join(root, "a.c"), filepath.Join(root, "link.c")); err != nil {
		t.Fatal(err)
	}
	return root
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

func TestRootCmdExt(t *testing.T) {
	root := tree(t)
	out, err := run(t, "--ext", ".c,.h", root)
	if err != nil {
		t.Fatal(err)
	}
	// a.c (2 lines, 14 bytes) + sub/b.h (1 line, 11 bytes); .git and the
	// symlink are not counted.
	if want := "\tfiles 2\tlines 3\tbytes 25\n"; !strings.HasSuffix(out, want) {
		t.Fatalf("got %q, want suffix %q", out, want)
	}
}

func TestRootCmdAllFilesContains(t *testing.T) {
	root := tree(t)
	out, err := run(t, "--contains", "int", root, filepath.Join(root, "a.c"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		root + "\tfiles 4\tlines 5\tbytes 62\tcontaining \"int\" 2\n",
		filepath.Join(root, "a.c") + "\tfiles 1\tlines 2\tbytes 14\tcontaining \"int\" 2\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

// Nothing matched is a reading error, never a count of zero.
func TestRootCmdNoFile(t *testing.T) {
	if _, err := run(t, "--ext", ".rs", tree(t)); err == nil || !strings.Contains(err.Error(), "no file") {
		t.Fatalf("err = %v, want a no-file error", err)
	}
}

func TestRootCmdMissing(t *testing.T) {
	if _, err := run(t, filepath.Join(t.TempDir(), "absent")); err == nil {
		t.Fatal("want an error for a missing path")
	}
}

func TestRootCmdUnreadable(t *testing.T) {
	root := tree(t)
	p := filepath.Join(root, "a.c")
	if err := os.Chmod(p, 0); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(p, 0o644)
	if _, err := os.ReadFile(p); err == nil {
		t.Skip("running with privileges that ignore file modes")
	}
	if _, err := run(t, root); err == nil {
		t.Fatal("want an error for an unreadable file")
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

	os.Args = []string{"srccount", tree(t)}
	main()
	if code != -1 {
		t.Fatalf("exit %d on a good tree", code)
	}
	os.Args = []string{"srccount"}
	main()
	if code != 1 {
		t.Fatalf("exit %d without an argument, want 1", code)
	}
}
