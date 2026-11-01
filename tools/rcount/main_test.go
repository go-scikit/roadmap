package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// namesC plants every eval-column case, a commented-out entry, an entry
// before the table and one after its terminator: only the six real ones count.
const namesC = `{"before", do_before, 0, 1, 1, {PP_FUNCALL, PREC_FN, 0}},
FUNTAB R_FunTab[] =
{
/* {"commented", do_c, 0, 11, 1, {PP_FUNCALL, PREC_FN, 0}}, */
{"if",		do_if,		0,	200,	-1,	{PP_IF,	     PREC_FN,	  1}},
{"c",		do_c,		0,	1,	-1,	{PP_FUNCALL, PREC_FN,	0}},
{"paste",	do_paste,	0,	11,	-1,	{PP_FUNCALL, PREC_FN,	0}},
{"eval",	do_eval,	0,	211,	3,	{PP_FUNCALL, PREC_FN,	0}},
{"quit",	do_quit,	PLUSOP,	111,	3,	{PP_FUNCALL, PREC_FN,	0}},
{"sp",		do_paste,	1,	10,	3,	{PP_FUNCALL, PREC_FN,	0}},
{NULL,		NULL,		0,	0,	0,	{PP_INVALID, PREC_FN,	0}},
{"after", do_after, 0, 1, 1, {PP_FUNCALL}},
`

const gramY = "%token A B\n%token <x> C\n%%\nprog : A\n | B\n ;\nexpr: C ;\n%%\nint x;\nint y;\n"

const evalC = "enum { X_OP, };\nenum {\n  BCMISMATCH_OP,\n  RETURN_OP,\n  DECLNKSTK_OP,\n  OPCOUNT\n};\n"

const lapackC = "F77_CALL(dgesv)(&n);\n/* F77_CALL(dcomment)(x) */\nF77_CALL( dgeqp3 )(a); F77_CALL(dgesv)(&m);\n"

func TestF77Callees(t *testing.T) {
	if got := strings.Join(f77Callees(lapackC), " "); got != "dgeqp3 dgesv" {
		t.Fatalf("got %q", got)
	}
}

func TestParseFunTab(t *testing.T) {
	f := parseFunTab(namesC)
	want := FunTab{Entries: 6, CFuns: 5, PrimBuiltin: 1, PrimSpecial: 1, IntBuiltin: 3, IntSpecial: 1}
	if f != want {
		t.Fatalf("got %+v, want %+v", f, want)
	}
	if f := parseFunTab("no table"); f.Entries != 0 {
		t.Fatalf("no table: %+v", f)
	}
	// Without a terminator the table runs to the end.
	if f := parseFunTab("FUNTAB R_FunTab[] =\n{\"c\", do_c, 0, 1, -1, {PP_FUNCALL}},\n"); f.Entries != 1 {
		t.Fatalf("unterminated: %+v", f)
	}
}

func TestStripCComments(t *testing.T) {
	if got := stripCComments("a/* x\ny */b/* open"); got != "a\nb" {
		t.Fatalf("got %q", got)
	}
}

func TestParseGram(t *testing.T) {
	g := parseGram(gramY)
	want := Gram{Lines: 10, Tokens: 3, Nonterminals: 2, RuleLines: 4, Epilogue: 2}
	if g != want {
		t.Fatalf("got %+v, want %+v", g, want)
	}
}

func TestParseOpcodes(t *testing.T) {
	ops := parseOpcodes(evalC)
	if strings.Join(ops, " ") != "BCMISMATCH_OP RETURN_OP DECLNKSTK_OP" {
		t.Fatalf("got %v", ops)
	}
	if parseOpcodes("enum { A_OP, };") != nil {
		t.Fatal("no OPCOUNT must give nil")
	}
	if parseOpcodes("  A_OP,\n  OPCOUNT\n") != nil {
		t.Fatal("no enum must give nil")
	}
}

func TestCountLines(t *testing.T) {
	for in, want := range map[string]int{"": 0, "a": 1, "a\n": 1, "a\nb": 2} {
		if got := countLines([]byte(in)); got != want {
			t.Errorf("countLines(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestLang(t *testing.T) {
	for in, want := range map[string]string{
		"a.R": "R", "a.r": "R", "a.c": "C", "a.h": "C", "a.f": "Fortran", "a.f90": "Fortran",
		"gram.y": "yacc", "x.Rd": "Rd", "Makefile.in": "", "a.Rout.save": "",
	} {
		if got := lang(in); got != want {
			t.Errorf("lang(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNative(t *testing.T) {
	n := newNative()
	n.add(`f <- function(x) .Internal(paste0(x, NULL))
# .Internal(commented(x))
g <- function() .Call(C_foo, 1); h <- .External2(C_bar)
s <- "a # not a comment .Internal(instr(1))"; e <- 'it\'s # still a string'
k <- .Primitive("c"); my.C(1); .C("cfun", 1L); .Fortran(C_lminfl, 2) # .Call(C_no)
.Call(C_foo, 2)`)
	want := map[string]int{".Internal": 2, ".Primitive": 1, ".Call": 2, ".External2": 1, ".External": 0, ".C": 1, ".Fortran": 1}
	for k, v := range want {
		if n.Sites[k] != v {
			t.Errorf("%s sites = %d, want %d", k, n.Sites[k], v)
		}
	}
	if len(n.Names[".Call"]) != 1 {
		t.Errorf(".Call names = %v, want one", n.Names[".Call"])
	}
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func tree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	src := filepath.Join(root, "src")
	write(t, filepath.Join(src, "main/names.c"), namesC)
	write(t, filepath.Join(src, "main/gram.y"), gramY)
	write(t, filepath.Join(src, "main/eval.c"), evalC)
	write(t, filepath.Join(src, "modules/lapack/Lapack.c"), lapackC)
	for _, d := range sizeDirs {
		write(t, filepath.Join(src, d, "README"), "not counted\n")
	}
	for _, p := range nativePkgs {
		write(t, filepath.Join(src, "library", p, "R", "a.R"), "f <- function() .Call(C_x)\n")
	}
	write(t, filepath.Join(src, "nmath/pnorm.c"), "a\nb\nc\n")
	write(t, filepath.Join(src, "appl/x.f"), "a\n")
	return root
}

func run(t *testing.T, root string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd := RootCmd()
	cmd.SetArgs([]string{root})
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	err := cmd.Execute()
	return out.String(), err
}

func TestRootCmd(t *testing.T) {
	out, err := run(t, tree(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"R_FunTab 6 entries, 5 distinct C functions\n",
		"  .Primitive 2 (builtin 1, special 1)\n",
		"  .Internal 4 (builtin 3, special 1)\n",
		"gram.y 10 lines: 3 tokens, 2 nonterminals, 4 rule lines, 2 epilogue lines\n",
		"byte-code opcodes 3 (BCMISMATCH_OP .. DECLNKSTK_OP)\n",
		"LAPACK/BLAS routines called by modules/lapack/Lapack.c 2: dgeqp3 dgesv\n",
		"  modules/lapack\tC\t1 files\t3 lines\n",
		"  main\tC\t2 files\t",
		"  main\tyacc\t1 files\t10 lines\n",
		"  nmath\tC\t1 files\t3 lines\n",
		"  appl\tFortran\t1 files\t1 lines\n",
		"  library/base\tR\t1 files\t1 lines\n",
		"  base\t.Internal 0/0\t.Primitive 0/0\t.Call 1/1\t",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

// Every file and directory the report reads is required, and an empty table
// is a reading error, never a count of zero.
func TestRootCmdBrokenTree(t *testing.T) {
	for name, edit := range map[string]func(src string) error{
		"no names.c":    func(src string) error { return os.Remove(filepath.Join(src, "main/names.c")) },
		"empty FunTab":  func(src string) error { return os.WriteFile(filepath.Join(src, "main/names.c"), []byte("x"), 0o644) },
		"no gram.y":     func(src string) error { return os.Remove(filepath.Join(src, "main/gram.y")) },
		"empty grammar": func(src string) error { return os.WriteFile(filepath.Join(src, "main/gram.y"), []byte("%%\n"), 0o644) },
		"no eval.c":     func(src string) error { return os.Remove(filepath.Join(src, "main/eval.c")) },
		"no opcodes":    func(src string) error { return os.WriteFile(filepath.Join(src, "main/eval.c"), []byte("x"), 0o644) },
		"no Lapack.c":   func(src string) error { return os.Remove(filepath.Join(src, "modules/lapack/Lapack.c")) },
		"no F77_CALL": func(src string) error {
			return os.WriteFile(filepath.Join(src, "modules/lapack/Lapack.c"), []byte("x"), 0o644)
		},
		"no size dir": func(src string) error { return os.RemoveAll(filepath.Join(src, "library/compiler")) },
		"no R dir":    func(src string) error { return os.RemoveAll(filepath.Join(src, "library/grid/R")) },
		"unreadable":  func(src string) error { return os.Chmod(filepath.Join(src, "nmath/pnorm.c"), 0) },
	} {
		t.Run(name, func(t *testing.T) {
			root := tree(t)
			if err := edit(filepath.Join(root, "src")); err != nil {
				t.Fatal(err)
			}
			if name == "unreadable" && os.Getuid() == 0 {
				t.Skip("root reads a mode-0 file")
			}
			if _, err := run(t, root); err == nil {
				t.Fatal("want an error")
			}
		})
	}
}

// cranIndex plants: a pure package on base only (a), a pure package whose
// dependency compiles (b -> c), a compiled one (c), an Rcpp user (d), a
// pure package on a missing dependency (e), a pure cycle (f <-> g), a pure
// chain on a pure package (h -> a, twice so the cache answers), a record
// without Package and a line without a colon.
const cranIndex = `Package: a
Depends: R (>= 4.1), stats
NeedsCompilation: no

Package: b
Imports: c (>=
        1.0), a
NeedsCompilation: no

Package: c
NeedsCompilation: yes

Package: d
LinkingTo: Rcpp (>= 1.0), BH
NeedsCompilation: yes

Package: e
Imports: bioconductorThing
NeedsCompilation: no

Package: f
Imports: g

Package: g
Imports: f

Package: i
NeedsCompilation: yes

Package: h
Imports: a, a

Version: 1.0
no colon here
`

func TestCranStats(t *testing.T) {
	got := cranStats(cranIndex)
	want := CRAN{Packages: 9, Compiled: 3, Rcpp: 1, PureClosure: 2}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if recs := parseDCF("  orphan continuation\nA: 1\n"); len(recs) != 1 || recs[0]["A"] != "1" {
		t.Fatalf("orphan continuation: %v", recs)
	}
}

func TestRootCmdCran(t *testing.T) {
	root := tree(t)
	idx := filepath.Join(t.TempDir(), "PACKAGES")
	write(t, idx, cranIndex)
	var out bytes.Buffer
	cmd := RootCmd()
	cmd.SetArgs([]string{"--cran", idx, root})
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"CRAN 9 packages: NeedsCompilation yes 3, LinkingTo Rcpp 1, pure R with pure-R hard dependencies 2\n",
		"  top compiled blockers: c d i\n",             // d and i tie at zero dependents: by name
		"  runnable with top 1 provided natively: 4\n", // a, h, and now c and b
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "top 5 provided") {
		t.Errorf("a step beyond the compiled packages must not print:\n%s", out.String())
	}
	for name, args := range map[string][]string{
		"missing index": {"--cran", filepath.Join(t.TempDir(), "absent"), root},
		"empty index":   {"--cran", filepath.Join(root, "src/main/gram.y"), root},
		"broken tree":   {"--cran", idx, t.TempDir()},
	} {
		cmd := RootCmd()
		cmd.SetArgs(args)
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		if err := cmd.Execute(); err == nil {
			t.Errorf("%s: want an error", name)
		}
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

	os.Args = []string{"rcount", tree(t)}
	main()
	if code != -1 {
		t.Fatalf("exit %d on a good tree", code)
	}
	os.Args = []string{"rcount"}
	main()
	if code != 1 {
		t.Fatalf("exit %d without an argument, want 1", code)
	}
}
