// Command srccount counts, under each given path, the regular files whose
// extension is in --ext (all files when --ext is empty), their lines and
// bytes, and, with --contains, the lines that contain a substring.
// Directories named .git are skipped. bib/c-cpp.md quotes its output.
//
//	go run ./tools/srccount --ext .cpp,.h clang/lib/Sema
//	go run ./tools/srccount --contains 'gcc.c-torture/execute/' known_failures.go
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// osExit is replaced by tests.
var osExit = os.Exit

// Count is the tally for one path.
type Count struct {
	Files, Lines, Bytes, Matches int
}

// countLines returns the number of newline-terminated lines of b, plus one
// for a final unterminated line.
func countLines(b []byte) int {
	n := bytes.Count(b, []byte{'\n'})
	if len(b) > 0 && b[len(b)-1] != '\n' {
		n++
	}
	return n
}

// countMatches returns the number of lines of b that contain sub.
func countMatches(b []byte, sub string) int {
	n := 0
	for _, line := range strings.Split(string(b), "\n") {
		if strings.Contains(line, sub) {
			n++
		}
	}
	return n
}

// wanted reports whether name has one of exts; an empty exts accepts all.
func wanted(name string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}
	e := filepath.Ext(name)
	for _, x := range exts {
		if e == x {
			return true
		}
	}
	return false
}

// count walks root and tallies the wanted regular files.
func count(root string, exts []string, contains string) (Count, error) {
	var c Count
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() || !wanted(d.Name(), exts) {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		c.Files++
		c.Lines += countLines(b)
		c.Bytes += len(b)
		if contains != "" {
			c.Matches += countMatches(b, contains)
		}
		return nil
	})
	return c, err
}

func report(w io.Writer, paths, exts []string, contains string) error {
	for _, p := range paths {
		c, err := count(p, exts, contains)
		if err != nil {
			return err
		}
		if c.Files == 0 {
			return fmt.Errorf("%s: no file with extension in %v", p, exts)
		}
		fmt.Fprintf(w, "%s\tfiles %d\tlines %d\tbytes %d", p, c.Files, c.Lines, c.Bytes)
		if contains != "" {
			fmt.Fprintf(w, "\tcontaining %q %d", contains, c.Matches)
		}
		fmt.Fprintln(w)
	}
	return nil
}

// RootCmd is the srccount command.
func RootCmd() *cobra.Command {
	var exts []string
	var contains string
	cmd := &cobra.Command{
		Use:          "srccount [--ext .c,.h] [--contains text] <path>...",
		Short:        "Count files, lines and bytes under source paths",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return report(cmd.OutOrStdout(), args, exts, contains)
		},
	}
	cmd.Flags().StringSliceVar(&exts, "ext", nil, "extensions to count, with the dot (default: all files)")
	cmd.Flags().StringVar(&contains, "contains", "", "also count the lines containing this text")
	return cmd
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		osExit(1)
	}
}
