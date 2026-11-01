// Command rcount takes the counts of bib/r.md on an R source tree (the
// release tarball from CRAN, unpacked): the .Primitive/.Internal table of
// src/main/names.c, the grammar of src/main/gram.y, the byte-code opcodes of
// src/main/eval.c, lines per language under the directories that matter for
// a Go port, and the native-code call sites in the R code of the base
// packages.
package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// osExit is replaced by tests.
var osExit = os.Exit

// sizeDirs are the directories, relative to src/, whose lines are counted.
var sizeDirs = []string{
	"main", "nmath", "appl", "extra/blas", "modules/lapack",
	"library/base", "library/stats", "library/utils", "library/methods",
	"library/graphics", "library/grDevices", "library/grid", "library/tools",
	"library/compiler", "../tests",
}

// nativePkgs are the packages whose R/ code is scanned for native calls.
var nativePkgs = []string{"base", "stats", "utils", "methods", "graphics", "grDevices", "grid", "tools"}

// nativeKinds are the interfaces from R code to compiled code.
var nativeKinds = []string{".Internal", ".Primitive", ".Call", ".External", ".External2", ".C", ".Fortran"}

// FunTab is the summary of R_FunTab in names.c.
type FunTab struct {
	Entries, CFuns           int
	PrimBuiltin, PrimSpecial int
	IntBuiltin, IntSpecial   int
}

// The offset column is a number or a symbolic constant (PLUSOP, ...).
var entryRE = regexp.MustCompile(`(?m)^[ \t]*\{"([^"]+)",\s*(\w+),\s*[\w-]+,\s*(\d+),`)

// parseFunTab reads the entries of R_FunTab up to its {NULL terminator. The
// eval column is XYZ: Y=1 is .Internal, else .Primitive; Z=1 is a builtin
// (arguments evaluated), Z=0 a special (names.c, header comment).
func parseFunTab(src string) FunTab {
	var f FunTab
	start := strings.Index(src, "FUNTAB R_FunTab[] =")
	if start < 0 {
		return f
	}
	body := src[start:]
	if end := strings.Index(body, "{NULL"); end >= 0 {
		body = body[:end]
	}
	cfuns := map[string]bool{}
	for _, m := range entryRE.FindAllStringSubmatch(stripCComments(body), -1) {
		ev, _ := strconv.Atoi(m[3])
		f.Entries++
		cfuns[m[2]] = true
		internal, builtin := (ev%100)/10 == 1, ev%10 == 1
		switch {
		case internal && builtin:
			f.IntBuiltin++
		case internal:
			f.IntSpecial++
		case builtin:
			f.PrimBuiltin++
		default:
			f.PrimSpecial++
		}
	}
	f.CFuns = len(cfuns)
	return f
}

// stripCComments blanks /* */ comments, keeping their newlines; an
// unterminated comment runs to the end.
func stripCComments(s string) string {
	var b strings.Builder
	for {
		i := strings.Index(s, "/*")
		if i < 0 {
			b.WriteString(s)
			return b.String()
		}
		b.WriteString(s[:i])
		j := strings.Index(s[i+2:], "*/")
		if j < 0 {
			return b.String()
		}
		end := i + 2 + j + 2
		b.WriteString(strings.Repeat("\n", strings.Count(s[i:end], "\n")))
		s = s[end:]
	}
}

// Gram is the summary of a yacc grammar.
type Gram struct{ Lines, Tokens, Nonterminals, RuleLines, Epilogue int }

var ntRE = regexp.MustCompile(`^[a-z_]+\s*:`)

func parseGram(src string) Gram {
	var g Gram
	section := 0
	for _, ln := range strings.Split(strings.TrimSuffix(src, "\n"), "\n") {
		g.Lines++
		switch {
		case strings.TrimSpace(ln) == "%%":
			section++
		case section == 0 && strings.HasPrefix(ln, "%token"):
			for _, t := range strings.Fields(ln)[1:] {
				if !strings.HasPrefix(t, "<") {
					g.Tokens++
				}
			}
		case section == 1:
			g.RuleLines++
			if ntRE.MatchString(ln) {
				g.Nonterminals++
			}
		case section > 1:
			g.Epilogue++
		}
	}
	return g
}

var opRE = regexp.MustCompile(`(?m)^\s*([A-Z0-9_]+_OP),`)

// parseOpcodes returns the byte-code opcodes of the enum that ends in OPCOUNT.
func parseOpcodes(src string) []string {
	end := strings.Index(src, "OPCOUNT\n")
	if end < 0 {
		return nil
	}
	start := strings.LastIndex(src[:end], "enum {")
	if start < 0 {
		return nil
	}
	var ops []string
	for _, m := range opRE.FindAllStringSubmatch(src[start:end], -1) {
		ops = append(ops, m[1])
	}
	return ops
}

var f77RE = regexp.MustCompile(`F77_CALL\(\s*(\w+)\s*\)`)

// f77Callees returns the sorted distinct Fortran routines named by F77_CALL
// outside comments.
func f77Callees(src string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range f77RE.FindAllStringSubmatch(stripCComments(src), -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	sort.Strings(out)
	return out
}

// countLines is wc -l plus one for a final unterminated line.
func countLines(b []byte) int {
	n := strings.Count(string(b), "\n")
	if len(b) > 0 && b[len(b)-1] != '\n' {
		n++
	}
	return n
}

// lang classifies a file by extension; "" means not counted.
func lang(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".r":
		return "R"
	case ".c", ".h":
		return "C"
	case ".f", ".f90":
		return "Fortran"
	case ".y":
		return "yacc"
	case ".rd":
		return "Rd"
	}
	return ""
}

// Stat is a file and line count.
type Stat struct{ Files, Lines int }

// walk calls fn with the contents of every counted file under root.
func walk(root string, fn func(path string, b []byte)) error {
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || lang(p) == "" {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		fn(p, b)
		return nil
	})
}

var (
	internalRE = regexp.MustCompile(`\.Internal\(\s*([A-Za-z.][\w.]*)\s*\(`)
	primRE     = regexp.MustCompile(`\.Primitive\(\s*"([^"]+)"`)
	callRE     = regexp.MustCompile(`(^|[^\w.])(\.Call|\.External2|\.External|\.C|\.Fortran)\(\s*([A-Za-z_.][\w.]*|"[^"]*")`)
)

// stripRComments drops # comments that are outside string literals.
func stripRComments(src string) string {
	lines := strings.Split(src, "\n")
	for k, ln := range lines {
		var q byte
		for i := 0; i < len(ln); i++ {
			c := ln[i]
			switch {
			case q != 0 && c == '\\':
				i++
			case q != 0 && c == q:
				q = 0
			case q != 0:
			case c == '"' || c == '\'' || c == '`':
				q = c
			case c == '#':
				ln = ln[:i]
			}
		}
		lines[k] = ln
	}
	return strings.Join(lines, "\n")
}

// Native counts call sites and distinct callee names per interface.
type Native struct {
	Sites map[string]int
	Names map[string]map[string]bool
}

func newNative() Native {
	n := Native{Sites: map[string]int{}, Names: map[string]map[string]bool{}}
	for _, k := range nativeKinds {
		n.Names[k] = map[string]bool{}
	}
	return n
}

func (n Native) add(src string) {
	src = stripRComments(src)
	for _, m := range internalRE.FindAllStringSubmatch(src, -1) {
		n.Sites[".Internal"]++
		n.Names[".Internal"][m[1]] = true
	}
	for _, m := range primRE.FindAllStringSubmatch(src, -1) {
		n.Sites[".Primitive"]++
		n.Names[".Primitive"][m[1]] = true
	}
	for _, m := range callRE.FindAllStringSubmatch(src, -1) {
		n.Sites[m[2]]++
		n.Names[m[2]][m[3]] = true
	}
}

func report(w io.Writer, root string) error {
	src := filepath.Join(root, "src")
	read := func(rel string) (string, error) {
		b, err := os.ReadFile(filepath.Join(src, rel))
		return string(b), err
	}
	text, err := read("main/names.c")
	if err != nil {
		return err
	}
	ft := parseFunTab(text)
	if ft.Entries == 0 {
		return fmt.Errorf("main/names.c: no R_FunTab entry found")
	}
	fmt.Fprintf(w, "R_FunTab %d entries, %d distinct C functions\n", ft.Entries, ft.CFuns)
	fmt.Fprintf(w, "  .Primitive %d (builtin %d, special %d)\n", ft.PrimBuiltin+ft.PrimSpecial, ft.PrimBuiltin, ft.PrimSpecial)
	fmt.Fprintf(w, "  .Internal %d (builtin %d, special %d)\n", ft.IntBuiltin+ft.IntSpecial, ft.IntBuiltin, ft.IntSpecial)

	if text, err = read("main/gram.y"); err != nil {
		return err
	}
	g := parseGram(text)
	if g.Nonterminals == 0 {
		return fmt.Errorf("main/gram.y: no grammar rule found")
	}
	fmt.Fprintf(w, "gram.y %d lines: %d tokens, %d nonterminals, %d rule lines, %d epilogue lines\n",
		g.Lines, g.Tokens, g.Nonterminals, g.RuleLines, g.Epilogue)

	if text, err = read("main/eval.c"); err != nil {
		return err
	}
	ops := parseOpcodes(text)
	if len(ops) == 0 {
		return fmt.Errorf("main/eval.c: no byte-code opcode found")
	}
	fmt.Fprintf(w, "byte-code opcodes %d (%s .. %s)\n", len(ops), ops[0], ops[len(ops)-1])

	if text, err = read("modules/lapack/Lapack.c"); err != nil {
		return err
	}
	callees := f77Callees(text)
	if len(callees) == 0 {
		return fmt.Errorf("modules/lapack/Lapack.c: no F77_CALL found")
	}
	fmt.Fprintf(w, "LAPACK/BLAS routines called by modules/lapack/Lapack.c %d: %s\n", len(callees), strings.Join(callees, " "))

	fmt.Fprintln(w, "lines per language")
	for _, dir := range sizeDirs {
		st := map[string]Stat{}
		if err := walk(filepath.Join(src, dir), func(p string, b []byte) {
			s := st[lang(p)]
			s.Files++
			s.Lines += countLines(b)
			st[lang(p)] = s
		}); err != nil {
			return err
		}
		for _, l := range []string{"R", "C", "Fortran", "yacc", "Rd"} {
			if s, ok := st[l]; ok {
				fmt.Fprintf(w, "  %s\t%s\t%d files\t%d lines\n", dir, l, s.Files, s.Lines)
			}
		}
	}

	fmt.Fprintln(w, "native call sites in R code (sites/distinct names)")
	for _, pkg := range nativePkgs {
		n := newNative()
		if err := walk(filepath.Join(src, "library", pkg, "R"), func(_ string, b []byte) { n.add(string(b)) }); err != nil {
			return err
		}
		fmt.Fprintf(w, "  %s", pkg)
		for _, k := range nativeKinds {
			fmt.Fprintf(w, "\t%s %d/%d", k, n.Sites[k], len(n.Names[k]))
		}
		fmt.Fprintln(w)
	}
	return nil
}

// basePkgs are the packages shipped with R itself (src/library), plus the
// pseudo-dependency "R".
var basePkgs = map[string]bool{
	"R": true, "base": true, "compiler": true, "datasets": true, "graphics": true, "grDevices": true,
	"grid": true, "methods": true, "parallel": true, "splines": true, "stats": true, "stats4": true,
	"tcltk": true, "tools": true, "utils": true,
}

// CRAN summarises a CRAN src/contrib/PACKAGES index.
type CRAN struct {
	Packages, Compiled, Rcpp int
	// PureClosure counts packages that need no compilation and whose hard
	// dependencies (Depends, Imports, LinkingTo), transitively, need none
	// either; a dependency absent from the index counts as compiled.
	PureClosure int
}

var depVersionRE = regexp.MustCompile(`\([^)]*\)`)

// parseDCF splits a DCF file into records of fields; continuation lines
// (starting with a blank) are joined to the previous field.
func parseDCF(text string) []map[string]string {
	var recs []map[string]string
	var cur map[string]string
	last := ""
	for _, ln := range strings.Split(text, "\n") {
		switch {
		case strings.TrimSpace(ln) == "":
			cur = nil
		case ln[0] == ' ' || ln[0] == '\t':
			if cur != nil && last != "" {
				cur[last] += " " + strings.TrimSpace(ln)
			}
		default:
			k, v, ok := strings.Cut(ln, ":")
			if !ok {
				continue
			}
			if cur == nil {
				cur = map[string]string{}
				recs = append(recs, cur)
			}
			last = k
			cur[k] = strings.TrimSpace(v)
		}
	}
	return recs
}

func deps(rec map[string]string) []string {
	var out []string
	for _, f := range []string{"Depends", "Imports", "LinkingTo"} {
		for _, d := range strings.Split(depVersionRE.ReplaceAllString(rec[f], ""), ",") {
			if d = strings.TrimSpace(d); d != "" {
				out = append(out, d)
			}
		}
	}
	return out
}

// index is a parsed CRAN PACKAGES file.
type index map[string]map[string]string

func parseIndex(text string) index {
	ix := index{}
	for _, r := range parseDCF(text) {
		if r["Package"] != "" {
			ix[r["Package"]] = r
		}
	}
	return ix
}

// closure counts the packages runnable when every package in avail (and
// every base package) is provided natively: a package is runnable if it is
// in avail, or needs no compilation and all its hard dependencies are
// runnable. A dependency absent from the index is not runnable; a
// dependency cycle is not runnable (CRAN forbids hard cycles).
func (ix index) closure(avail map[string]bool) int {
	state := map[string]int{} // 1 runnable, 2 not, 3 being visited
	var ok func(p string) bool
	ok = func(p string) bool {
		if basePkgs[p] || avail[p] {
			return true
		}
		switch state[p] {
		case 1:
			return true
		case 2, 3:
			return false
		}
		r, found := ix[p]
		if !found || r["NeedsCompilation"] == "yes" {
			state[p] = 2
			return false
		}
		state[p] = 3
		for _, d := range deps(r) {
			if !ok(d) {
				state[p] = 2
				return false
			}
		}
		state[p] = 1
		return true
	}
	n := 0
	for p := range ix {
		if ok(p) {
			n++
		}
	}
	return n
}

// blockers ranks the compiled packages by how many packages depend on them,
// directly or transitively, through hard dependencies (ties by name).
func (ix index) blockers() []string {
	rev := map[string][]string{}
	for p, r := range ix {
		for _, d := range deps(r) {
			rev[d] = append(rev[d], p)
		}
	}
	reach := map[string]int{}
	var compiled []string
	for p, r := range ix {
		if r["NeedsCompilation"] != "yes" {
			continue
		}
		compiled = append(compiled, p)
		seen := map[string]bool{p: true}
		queue := []string{p}
		for len(queue) > 0 {
			q := queue[0]
			queue = queue[1:]
			for _, u := range rev[q] {
				if !seen[u] {
					seen[u] = true
					queue = append(queue, u)
				}
			}
		}
		reach[p] = len(seen) - 1
	}
	sort.Slice(compiled, func(i, j int) bool {
		a, b := compiled[i], compiled[j]
		if reach[a] != reach[b] {
			return reach[a] > reach[b]
		}
		return a < b
	})
	return compiled
}

func cranStats(text string) CRAN {
	var c CRAN
	ix := parseIndex(text)
	for _, r := range ix {
		c.Packages++
		if r["NeedsCompilation"] == "yes" {
			c.Compiled++
		}
		for _, d := range strings.Split(depVersionRE.ReplaceAllString(r["LinkingTo"], ""), ",") {
			if strings.TrimSpace(d) == "Rcpp" {
				c.Rcpp++
			}
		}
	}
	c.PureClosure = ix.closure(nil)
	return c
}

// unlockSteps are the numbers of top blockers provided natively in the
// unlock table.
var unlockSteps = []int{1, 5, 10, 20, 50, 100}

// unlockTable writes how many packages become runnable when the top-k
// blockers are provided natively, for each k of unlockSteps.
func unlockTable(w io.Writer, ix index) {
	order := ix.blockers()
	top := 10
	if len(order) < top {
		top = len(order)
	}
	fmt.Fprintf(w, "  top compiled blockers: %s\n", strings.Join(order[:top], " "))
	for _, k := range unlockSteps {
		if k > len(order) {
			break
		}
		avail := map[string]bool{}
		for _, p := range order[:k] {
			avail[p] = true
		}
		fmt.Fprintf(w, "  runnable with top %d provided natively: %d\n", k, ix.closure(avail))
	}
}

// RootCmd is the rcount command.
func RootCmd() *cobra.Command {
	var cranFile string
	cmd := &cobra.Command{
		Use:          "rcount <R-source-root>",
		Short:        "Count the primitives, internals, grammar, opcodes and lines of an R source tree",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := report(cmd.OutOrStdout(), args[0]); err != nil {
				return err
			}
			if cranFile == "" {
				return nil
			}
			b, err := os.ReadFile(cranFile)
			if err != nil {
				return err
			}
			c := cranStats(string(b))
			if c.Packages == 0 {
				return fmt.Errorf("%s: no package record found", cranFile)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "CRAN %d packages: NeedsCompilation yes %d, LinkingTo Rcpp %d, pure R with pure-R hard dependencies %d\n",
				c.Packages, c.Compiled, c.Rcpp, c.PureClosure)
			unlockTable(cmd.OutOrStdout(), parseIndex(string(b)))
			return nil
		},
	}
	cmd.Flags().StringVar(&cranFile, "cran", "", "a CRAN src/contrib/PACKAGES index to summarise")
	return cmd
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		osExit(1)
	}
}
