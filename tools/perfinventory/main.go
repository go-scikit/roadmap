// Command perfinventory counts, in each repository given as a directory, the
// exported functions and methods of its public packages, and how each is
// measured. docs/performance-campaign.md quotes its output.
//
//	go run ./tools/perfinventory [--list] <repo-dir>...
//
// The rules, which the campaign document restates:
//
//   - Every directory holding a go.mod is a module, except under testdata,
//     vendor and dot-directories; each is loaded with its tests through
//     golang.org/x/tools/go/packages, so names are resolved by the type
//     checker, not by spelling.
//   - Public: an exported function, or an exported method of an exported
//     named type, declared in a non-test file of a non-main package whose
//     import path has no internal, cmd, tools, bench, benchmarks, examples
//     or testdata element. Methods that satisfy a standard interface
//     (Error, String, Format, the Marshal family, Unwrap, Is) are left out.
//   - A use is a direct reference resolved to that function (generic
//     instantiations to their origin). A call through an interface, a
//     function value built elsewhere or a helper is not a use.
//   - A use is timed when it sits in a benchmark (named Benchmark… or taking
//     a *testing.B), in a function that calls time.Now or time.Since, in a
//     package-level declaration of a file that has one of those, or anywhere
//     in a harness package (an import path with a bench, benchmarks, cmd or
//     tools element) that times something.
//   - The classes, best first: ref (timed, and the file's code — string
//     literals, identifiers, the import "C" — names a reference
//     implementation, or the use is in a harness package whose code does);
//     noted (timed, and the reference is named only in the file's comments:
//     the comparison is run by a script outside Go); bench (timed, no
//     reference named); none.
//   - Non-Go files (.py, .sh, .c, .jl, .r) that contain a timing word are
//     counted per repository as external harnesses: their pairing with Go
//     functions is not machine-readable.
package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

// osExit is replaced by tests.
var osExit = os.Exit

var (
	referenceMarkers = []string{"numpy", "scipy", "mpmath", "statsmodels", "fftw",
		"pocketfft", "qdldl", "osqp", "shewchuk"}
	// literalMarkers are how a reference shows in a string literal that is
	// code (a Python program, a script path), not prose ("SciPy %g").
	literalMarkers = []string{"import numpy", "import scipy", "import mpmath", "import statsmodels",
		"np.", "scipy.", "mpmath.", "statsmodels.", "scipy_", "numpy_", "mpmath_", "fftw",
		"pocketfft", "qdldl", "osqp"}
	externalTiming = []string{"perf_counter", "timeit", "clock_gettime", "time.time(",
		"time_ns", "chrono::", "@btime", "system.time"}
	privateElems = map[string]bool{"internal": true, "cmd": true, "tools": true, "bench": true,
		"benchmarks": true, "examples": true, "testdata": true}
	plumbing = map[string]bool{"Error": true, "String": true, "GoString": true, "Format": true,
		"MarshalJSON": true, "UnmarshalJSON": true, "MarshalText": true, "UnmarshalText": true,
		"MarshalBinary": true, "UnmarshalBinary": true, "Unwrap": true, "Is": true}
	externalExts = map[string]bool{".py": true, ".sh": true, ".c": true, ".jl": true, ".r": true}
)

// Class is how a public function is measured.
type Class int

// The classes, from worst to best.
const (
	None Class = iota
	Bench
	Noted
	Ref
)

func (c Class) String() string { return [...]string{"none", "bench", "noted", "ref"}[c] }

// Repo is the tally of one repository.
type Repo struct {
	Dir              string
	Modules          int
	Funcs            map[string]Class  // key: import/path.Func or import/path.Type.Method
	Why              map[string]string // the file that gave the class
	ExternalHarness  int
	ExternalExamples []string
}

// Count returns how many functions are in class c.
func (r Repo) Count(c Class) int {
	n := 0
	for _, v := range r.Funcs {
		if v == c {
			n++
		}
	}
	return n
}

// containsAny reports whether s contains one of subs.
func containsAny(s string, subs []string) bool {
	for _, x := range subs {
		if strings.Contains(s, x) {
			return true
		}
	}
	return false
}

// hasElem reports whether the import path has one of the elements.
func hasElem(path string, elems map[string]bool) bool {
	for _, e := range strings.Split(path, "/") {
		if elems[e] {
			return true
		}
	}
	return false
}

var harnessElems = map[string]bool{"bench": true, "benchmarks": true, "cmd": true, "tools": true}

// public reports whether a package with this import path and name is public.
func public(path, name string) bool {
	return name != "main" && !hasElem(path, privateElems)
}

// skipDir reports whether a walk must not enter the directory (root excepted).
func skipDir(name string) bool {
	return name == "testdata" || name == "vendor" || strings.HasPrefix(name, ".")
}

// modules returns the directories under root that hold a go.mod, and the
// number of non-Go files carrying a timing word (with up to 3 of their
// paths, relative to root).
func modules(root string) (mods []string, ext int, examples []string, err error) {
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p != root && skipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "go.mod" {
			mods = append(mods, filepath.Dir(p))
			return nil
		}
		if !externalExts[strings.ToLower(filepath.Ext(p))] {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if containsAny(strings.ToLower(string(b)), externalTiming) {
			ext++
			if len(examples) < 3 {
				rel, _ := filepath.Rel(root, p)
				examples = append(examples, rel)
			}
		}
		return nil
	})
	return mods, ext, examples, err
}

// funcKey names a function: path.Name, or path.Type.Name for a method of a
// named type; ok is false for a method of anything else.
func funcKey(f *types.Func) (string, bool) {
	f = f.Origin()
	if f.Pkg() == nil {
		return "", false
	}
	sig := f.Type().(*types.Signature)
	if sig.Recv() == nil {
		return f.Pkg().Path() + "." + f.Name(), true
	}
	t := sig.Recv().Type()
	if p, isPtr := t.(*types.Pointer); isPtr {
		t = p.Elem()
	}
	n, isNamed := t.(*types.Named)
	if !isNamed {
		return "", false
	}
	return f.Pkg().Path() + "." + n.Obj().Name() + "." + f.Name(), true
}

// declare adds the public functions of pkg to funcs.
func declare(pkg *packages.Package, funcs map[string]Class) {
	if strings.Contains(pkg.ID, "[") || strings.HasSuffix(pkg.PkgPath, "_test") ||
		!public(pkg.PkgPath, pkg.Name) {
		return
	}
	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if !obj.Exported() {
			continue
		}
		switch o := obj.(type) {
		case *types.Func:
			k, _ := funcKey(o)
			funcs[k] = None
		case *types.TypeName:
			n, isNamed := o.Type().(*types.Named)
			if !isNamed || o.IsAlias() {
				continue
			}
			for i := 0; i < n.NumMethods(); i++ {
				if m := n.Method(i); m.Exported() && !plumbing[m.Name()] {
					k, _ := funcKey(m)
					funcs[k] = None
				}
			}
		}
	}
}

// benchmarkFunc reports whether fd is a Go benchmark: named Benchmark… or
// taking a *testing.B.
func benchmarkFunc(fd *ast.FuncDecl) bool {
	if strings.HasPrefix(fd.Name.Name, "Benchmark") {
		return true
	}
	for _, f := range fd.Type.Params.List {
		if s, isStar := f.Type.(*ast.StarExpr); isStar {
			if sel, isSel := s.X.(*ast.SelectorExpr); isSel && sel.Sel.Name == "B" {
				return true
			}
		}
	}
	return false
}

// callsClock reports whether n references time.Now or time.Since.
func callsClock(n ast.Node, info *types.Info) bool {
	found := false
	ast.Inspect(n, func(n ast.Node) bool {
		if id, isIdent := n.(*ast.Ident); isIdent {
			if f, isF := info.Uses[id].(*types.Func); isF && f.Pkg() != nil &&
				f.Pkg().Path() == "time" && (f.Name() == "Now" || f.Name() == "Since") {
				found = true
			}
		}
		return !found
	})
	return found
}

// codeNamesReference reports whether the code of file (string literals,
// identifiers, an import "C"), comments excluded, names a reference.
func codeNamesReference(file *ast.File) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			found = found || x.Path.Value == `"C"`
		case *ast.BasicLit:
			found = found || (x.Kind == token.STRING && containsAny(strings.ToLower(x.Value), literalMarkers))
		case *ast.Ident:
			low := strings.ToLower(x.Name)
			found = found || containsAny(low, referenceMarkers) || strings.HasPrefix(low, "_cfunc_")
		}
		return !found
	})
	return found
}

// commentsNameReference reports whether a comment of file names a reference.
func commentsNameReference(file *ast.File) bool {
	for _, g := range file.Comments {
		if containsAny(strings.ToLower(g.Text()), referenceMarkers) {
			return true
		}
	}
	return false
}

// raise sets r.Funcs[k] to c, and r.Why[k] to the file, when c is better
// than what it holds.
func raise(r *Repo, k string, c Class, file string) {
	if old, ok := r.Funcs[k]; ok && c > old {
		r.Funcs[k] = c
		r.Why[k] = file
	}
}

// use records, for every file of pkg, the public functions it references.
func use(pkg *packages.Package, r *Repo) {
	harness := hasElem(pkg.PkgPath, harnessElems)
	pkgRef, pkgTimes := false, false
	for _, f := range pkg.Syntax {
		pkgRef = pkgRef || codeNamesReference(f)
		pkgTimes = pkgTimes || callsClock(f, pkg.TypesInfo)
	}
	for i, file := range pkg.Syntax {
		rel, _ := filepath.Rel(r.Dir, pkg.CompiledGoFiles[i])
		perfFile := strings.HasPrefix(filepath.Base(rel), "perf")
		best := Bench
		switch {
		case codeNamesReference(file) || ((harness || perfFile) && pkgRef):
			best = Ref
		case commentsNameReference(file):
			best = Noted
		}
		// A harness package or a perf* file is a measurement throughout
		// (its case tables are built in init or helpers); elsewhere only
		// the timing functions and the package-level tables of a file that
		// times count.
		whole := (harness && pkgTimes) || (perfFile && callsClock(file, pkg.TypesInfo))
		fileTimes := whole
		for _, decl := range file.Decls {
			if fd, isFunc := decl.(*ast.FuncDecl); isFunc && (benchmarkFunc(fd) || callsClock(fd, pkg.TypesInfo)) {
				fileTimes = true
			}
		}
		for _, decl := range file.Decls {
			timed := whole
			if fd, isFunc := decl.(*ast.FuncDecl); isFunc {
				timed = timed || benchmarkFunc(fd) || callsClock(fd, pkg.TypesInfo)
			} else {
				timed = timed || fileTimes
			}
			if !timed {
				continue
			}
			ast.Inspect(decl, func(n ast.Node) bool {
				if id, isIdent := n.(*ast.Ident); isIdent {
					if f, isF := pkg.TypesInfo.Uses[id].(*types.Func); isF {
						if k, ok := funcKey(f); ok {
							raise(r, k, best, rel)
						}
					}
				}
				return true
			})
		}
	}
}

// load is replaced by tests that need a loader failure.
var load = packages.Load

// inventory tallies one repository.
func inventory(dir string) (Repo, error) {
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	r := Repo{Dir: dir, Funcs: map[string]Class{}, Why: map[string]string{}}
	mods, ext, examples, err := modules(dir)
	if err != nil {
		return r, err
	}
	r.Modules, r.ExternalHarness, r.ExternalExamples = len(mods), ext, examples
	var all []*packages.Package
	for _, m := range mods {
		cfg := &packages.Config{
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
				packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
			Dir:   m,
			Tests: true,
		}
		pkgs, err := load(cfg, "./...")
		if err != nil {
			return r, fmt.Errorf("%s: %w", m, err)
		}
		if n := packages.PrintErrors(pkgs); n > 0 {
			return r, fmt.Errorf("%s: %d package errors", m, n)
		}
		all = append(all, pkgs...)
	}
	for _, p := range all {
		declare(p, r.Funcs)
	}
	for _, p := range all {
		use(p, &r)
	}
	return r, nil
}

func report(w io.Writer, dirs []string, list bool) error {
	var tot [4]int
	total, ext := 0, 0
	fmt.Fprintln(w, "| repository | modules | public functions | ref | noted | bench only | none | external timing scripts |")
	fmt.Fprintln(w, "|---|---:|---:|---:|---:|---:|---:|---:|")
	var repos []Repo
	for _, d := range dirs {
		r, err := inventory(d)
		if err != nil {
			return err
		}
		repos = append(repos, r)
		var c [4]int
		for i := range c {
			c[i] = r.Count(Class(i))
			tot[i] += c[i]
		}
		total += len(r.Funcs)
		ext += r.ExternalHarness
		fmt.Fprintf(w, "| %s | %d | %d | %d | %d | %d | %d | %d |\n", filepath.Base(r.Dir), r.Modules,
			len(r.Funcs), c[Ref], c[Noted], c[Bench], c[None], r.ExternalHarness)
	}
	fmt.Fprintf(w, "| **total** | | %d | %d | %d | %d | %d | %d |\n", total, tot[Ref], tot[Noted], tot[Bench], tot[None], ext)
	if !list {
		return nil
	}
	for _, r := range repos {
		fmt.Fprintf(w, "\n## %s\n", filepath.Base(r.Dir))
		for _, e := range r.ExternalExamples {
			fmt.Fprintf(w, "external\t%s\n", e)
		}
		byFile := map[string]map[Class]int{}
		for k, f := range r.Why {
			if byFile[f] == nil {
				byFile[f] = map[Class]int{}
			}
			byFile[f][r.Funcs[k]]++
		}
		files := make([]string, 0, len(byFile))
		for f := range byFile {
			files = append(files, f)
		}
		sort.Strings(files)
		for _, f := range files {
			fmt.Fprintf(w, "file\t%s\tref %d\tnoted %d\tbench %d\n", f, byFile[f][Ref], byFile[f][Noted], byFile[f][Bench])
		}
		keys := make([]string, 0, len(r.Funcs))
		for k := range r.Funcs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.Funcs[k], k, r.Why[k])
		}
	}
	return nil
}

// RootCmd is the perfinventory command.
func RootCmd() *cobra.Command {
	var list bool
	cmd := &cobra.Command{
		Use:          "perfinventory [--list] <repo-dir>...",
		Short:        "Count public functions timed against a reference, timed alone, or not timed",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return report(cmd.OutOrStdout(), args, list)
		},
	}
	cmd.Flags().BoolVar(&list, "list", false, "also list every public function with its class and file")
	return cmd
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		osExit(1)
	}
}
