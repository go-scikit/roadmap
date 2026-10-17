// Command checklinks fails when a relative link in a Markdown file of the
// tree points at a file that does not exist.
//
//	go run ./tools/checklinks [root]
package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// link matches the target of an inline Markdown link or image.
var link = regexp.MustCompile(`\]\(([^)\s]+)\)`)

// broken returns, for every Markdown file under root, the relative links
// whose target is missing, as "file: target".
func broken(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, m := range link.FindAllStringSubmatch(string(b), -1) {
			target := m[1]
			if strings.Contains(target, "://") || strings.HasPrefix(target, "#") || strings.HasPrefix(target, "mailto:") {
				continue
			}
			target, _, _ = strings.Cut(target, "#")
			if _, err := os.Stat(filepath.Join(filepath.Dir(path), target)); err != nil {
				out = append(out, path+": "+m[1])
			}
		}
		return nil
	})
	return out, err
}

func run(args []string, w io.Writer) int {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	bad, err := broken(root)
	if err != nil {
		fmt.Fprintln(w, err)
		return 2
	}
	for _, b := range bad {
		fmt.Fprintln(w, "broken link:", b)
	}
	if len(bad) > 0 {
		return 1
	}
	return 0
}

// exit is os.Exit, replaced by the tests.
var exit = os.Exit

func main() { exit(run(os.Args[1:], os.Stderr)) }
