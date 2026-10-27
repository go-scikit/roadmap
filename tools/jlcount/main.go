// Command jlcount counts, in a julia source tree, the runtime surface a Go
// port must re-implement: the builtins of src/builtin_proto.h, the
// intrinsics of src/intrinsics.h by kind, and the lines of the C files of
// the semantic core. bib/julia-runtime.md quotes its output.
//
//	go run ./tools/jlcount <julia-source-root>
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// osExit is replaced by tests.
var osExit = os.Exit

var (
	// builtinRe matches one XX(cname,"jlname") entry of JL_BUILTIN_FUNCTIONS.
	builtinRe = regexp.MustCompile(`^\s*XX\(\s*(\w+)\s*,\s*"([^"]+)"\s*\)`)
	// intrinsicRe matches one ADD_I, ALIAS or ADD_HIDDEN entry of INTRINSICS.
	intrinsicRe = regexp.MustCompile(`^\s*(ADD_I|ALIAS|ADD_HIDDEN)\(\s*(\w+)\s*,`)
)

// coreFiles are the C files of the semantic core, in the order of the
// runtime design document.
var coreFiles = []string{
	"jltypes.c", "datatype.c", "subtype.c", "gf.c", "typemap.c", "method.c",
	"builtins.c", "runtime_intrinsics.c", "interpreter.c", "toplevel.c",
	"module.c", "genericmemory.c", "array.c", "symbol.c", "simplevector.c",
}

// Builtin is one julia-callable builtin function.
type Builtin struct{ CName, JLName string }

// Intrinsic is one entry of the INTRINSICS table.
type Intrinsic struct{ Kind, Name string }

// parseBuiltins returns the entries of JL_BUILTIN_FUNCTIONS, in order.
// Only lines inside the macro body are read: the body ends at the first line
// that does not end in a backslash.
func parseBuiltins(text string) []Builtin {
	var out []Builtin
	in := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#define JL_BUILTIN_FUNCTIONS(") {
			in = true
			continue
		}
		if !in {
			continue
		}
		if m := builtinRe.FindStringSubmatch(line); m != nil {
			out = append(out, Builtin{m[1], m[2]})
		}
		if !strings.HasSuffix(strings.TrimRight(line, " \t"), `\`) {
			break
		}
	}
	return out
}

// parseIntrinsics returns the entries of the INTRINSICS macro, in order.
func parseIntrinsics(text string) []Intrinsic {
	var out []Intrinsic
	in := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#define INTRINSICS") {
			in = true
			continue
		}
		if !in {
			continue
		}
		if m := intrinsicRe.FindStringSubmatch(line); m != nil {
			out = append(out, Intrinsic{m[1], m[2]})
		}
		if !strings.HasSuffix(strings.TrimRight(line, " \t"), `\`) {
			break
		}
	}
	return out
}

// countLines returns the number of newline-terminated lines of b, plus one
// for a final unterminated line, as wc -l plus that line.
func countLines(b []byte) int {
	n := strings.Count(string(b), "\n")
	if len(b) > 0 && b[len(b)-1] != '\n' {
		n++
	}
	return n
}

func report(w io.Writer, root string) error {
	src := filepath.Join(root, "src")
	b, err := os.ReadFile(filepath.Join(src, "builtin_proto.h"))
	if err != nil {
		return err
	}
	builtins := parseBuiltins(string(b))
	if len(builtins) == 0 {
		return fmt.Errorf("%s: no JL_BUILTIN_FUNCTIONS entry found", filepath.Join(src, "builtin_proto.h"))
	}
	b, err = os.ReadFile(filepath.Join(src, "intrinsics.h"))
	if err != nil {
		return err
	}
	intrinsics := parseIntrinsics(string(b))
	if len(intrinsics) == 0 {
		return fmt.Errorf("%s: no INTRINSICS entry found", filepath.Join(src, "intrinsics.h"))
	}

	fmt.Fprintf(w, "builtins %d\n", len(builtins))
	for _, x := range builtins {
		fmt.Fprintf(w, "  %s\t%s\n", x.JLName, x.CName)
	}
	kinds := map[string]int{}
	for _, x := range intrinsics {
		kinds[x.Kind]++
	}
	fmt.Fprintf(w, "intrinsics %d (ADD_I %d, ALIAS %d, ADD_HIDDEN %d)\n",
		len(intrinsics), kinds["ADD_I"], kinds["ALIAS"], kinds["ADD_HIDDEN"])
	for _, x := range intrinsics {
		fmt.Fprintf(w, "  %s\t%s\n", x.Name, x.Kind)
	}
	total := 0
	fmt.Fprintln(w, "core files")
	for _, f := range coreFiles {
		b, err := os.ReadFile(filepath.Join(src, f))
		if err != nil {
			return err
		}
		n := countLines(b)
		total += n
		fmt.Fprintf(w, "  %s\t%d\n", f, n)
	}
	fmt.Fprintf(w, "core total %d files %d lines\n", len(coreFiles), total)
	return nil
}

// RootCmd is the jlcount command.
func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "jlcount <julia-source-root>",
		Short:        "Count the builtins, intrinsics and core C lines of a julia source tree",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return report(cmd.OutOrStdout(), args[0])
		},
	}
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		osExit(1)
	}
}
