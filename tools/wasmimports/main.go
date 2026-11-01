// Command wasmimports validates a WebAssembly module with wazero, with the
// post-2.0 features a clang/emscripten build may use enabled (exception
// handling, tail calls, extended constants, threads), and prints what a Go
// host must provide to instantiate it: the imported functions grouped by
// module, the imported memories, and the custom sections (dylink.0 marks an
// emscripten side module). bib/c-cpp.md quotes its output.
//
//	go run ./tools/wasmimports libclangCppInterOp.so
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
)

// osExit is replaced by tests.
var osExit = os.Exit

// features are the core features enabled for validation.
const features = api.CoreFeaturesV2 |
	experimental.CoreFeaturesExceptionHandling |
	experimental.CoreFeaturesTailCall |
	experimental.CoreFeaturesExtendedConst |
	experimental.CoreFeaturesThreads

func report(w io.Writer, path string, names bool) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ctx := context.Background()
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter().
		WithCoreFeatures(features).WithCustomSections(true))
	defer r.Close(ctx)
	m, err := r.CompileModule(ctx, b)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	fmt.Fprintf(w, "%s: valid, %d bytes\n", path, len(b))

	byModule := map[string][]string{}
	for _, f := range m.ImportedFunctions() {
		mod, name, _ := f.Import()
		byModule[mod] = append(byModule[mod], name)
	}
	mods := make([]string, 0, len(byModule))
	for k := range byModule {
		mods = append(mods, k)
	}
	sort.Strings(mods)
	fmt.Fprintf(w, "imported functions %d\n", len(m.ImportedFunctions()))
	for _, k := range mods {
		fmt.Fprintf(w, "  %s\t%d\n", k, len(byModule[k]))
		if names {
			sort.Strings(byModule[k])
			for _, n := range byModule[k] {
				fmt.Fprintf(w, "    %s\n", n)
			}
		}
	}
	for _, mem := range m.ImportedMemories() {
		mod, name, _ := mem.Import()
		max, _ := mem.Max()
		fmt.Fprintf(w, "imported memory %s.%s min %d pages max %d pages\n", mod, name, mem.Min(), max)
	}
	fmt.Fprintf(w, "exported functions %d\n", len(m.ExportedFunctions()))
	for _, s := range m.CustomSections() {
		fmt.Fprintf(w, "custom section %s\t%d bytes\n", s.Name(), len(s.Data()))
	}
	return nil
}

// RootCmd is the wasmimports command.
func RootCmd() *cobra.Command {
	var names bool
	cmd := &cobra.Command{
		Use:          "wasmimports [--names] <file.wasm>",
		Short:        "Validate a wasm module and list what its host must provide",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return report(cmd.OutOrStdout(), args[0], names)
		},
	}
	cmd.Flags().BoolVar(&names, "names", false, "list every imported function name")
	return cmd
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		osExit(1)
	}
}
