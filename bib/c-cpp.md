# C and C++ — bibliography and feasibility plan for go-embedded-c and go-embedded-cpp

Written 2026-11-01. Facts only, each with a source; counts come from the sources, not from
estimates. Goal: ` ```c ` and ` ```cpp ` cells in go-notebook and its WebAssembly playground,
backed by Go code built with `CGO_ENABLED=0`, embeddable as a library and runnable as wasm, in
the spirit of [rbgo](https://github.com/go-embedded-ruby/ruby) and
[go-embedded-julia/julia](https://github.com/go-embedded-julia/julia). The organisations
[go-embedded-c](https://github.com/go-embedded-c) and [go-embedded-cpp](https://github.com/go-embedded-cpp)
were created on 2026-11-01 (`gh api orgs/…` `created_at` 10:01:16Z and 10:02:06Z) and hold no
repository (`gh repo list` returned `[]` for both).

**Decisions.**

- **C: reuse the front end, write the back end.** Parse and type-check with
  [modernc.org/cc/v4](https://pkg.go.dev/modernc.org/cc/v4) (BSD-3, pure Go, already runs under
  wasm, §2.2); write in Go a **C → wasm32 code generator** and a **cell linker** that loads each
  cell as a module sharing one linear memory and one function table, executed by
  [wazero](https://github.com/tetratelabs/wazero) on the desktop and by the browser's own engine
  in the playground. libc is **wasi-libc** (musl, MIT), compiled once, not rewritten. Do not
  transpile to Go (ccgo) and do not write a tree interpreter (§2.3).
- **C++: no front end of our own.** Run the upstream incremental compiler, clang-repl as packaged
  by CppInterOp/xeus-cpp-lite, **as a WebAssembly module** on the cell runtime built for C
  (§3.3, option b). go-embedded-cpp holds a host and packaging, not a compiler. It starts only
  after go-embedded-c M3, and behind a measured gate (§3.4): today wazero refuses the clang
  module (legacy exception handling), and the module is 71 MB.

## 0. Sources read, and how the counts were made

| source | revision read | what was used |
|---|---|---|
| WG14 [N1570](https://www.open-std.org/jtc1/sc22/wg14/www/docs/n1570.pdf) | "Committee Draft — April 12, 2011", ISO/IEC 9899:201x (cover page read) | C11 text, the language reference |
| WG14 [N2310](https://www.open-std.org/jtc1/sc22/wg14/www/docs/n2310.pdf) | November 6, 2018 (cover read) | a C2x working draft whose diff marks are "from the previous draft (ISO/IEC 9899:2018)", i.e. the C17 text plus marked changes |
| WG14 [N2176](https://www.open-std.org/jtc1/sc22/wg14/www/docs/n2176.pdf) | downloaded | the C17 draft; **password-protected**, not read |
| WG14 [N3220](https://www.open-std.org/jtc1/sc22/wg14/www/docs/n3220.pdf) | "ISO/IEC 9899:2024 (en) — N3220 working draft" (cover read) | C23, for what v1 leaves out |
| [modernc.org/cc/v4](https://gitlab.com/cznic/cc) | v4.30.5 (module cache) | `cc.go` (`Config`, `NewConfig`, `Translate`), `abi.go`, `parser.yy`, `all_test.go`, `LICENSE` |
| [modernc.org/ccgo/v4](https://gitlab.com/cznic/ccgo) | v4.38.0 | `README.md`, `lib/all_test.go`, `lib/known_failures_*_test.go` |
| [modernc.org/libc](https://gitlab.com/cznic/libc) | v1.79.1 | `README.md`, file list (targets) |
| [modernc.org/ccorpus2](https://gitlab.com/cznic/ccorpus2) | v1.6.0 | the corpus ccgo and cc test against (`assets/`) |
| [gotranspile/cxgo](https://github.com/gotranspile/cxgo) | `e682671cd3c9`, 2025-03-16 | `README.md`, `go.mod` |
| [hajimehoshi/goc](https://github.com/hajimehoshi/goc) | `142a171d62cb`, 2018-09-27 | `README.md`, tree (the only "C interpreter in Go" a GitHub search returned) |
| [llvm/llvm-project](https://github.com/llvm/llvm-project) | tag `llvmorg-23.1.4`, commit `7a941dfb2877`, 2026-10-19 (latest release, published 2026-10-20) | `clang/lib/{Sema,Parse,AST,Lex,CodeGen,Interpreter}`, `clang/lib/Interpreter/Wasm.cpp`, `clang/docs/ClangRepl.rst`, `clang/lib/Basic/Targets/WebAssembly.h`, `libcxx/include`, `llvm/lib/Target/WebAssembly` (sparse checkout) |
| [compiler-research/xeus-cpp](https://github.com/compiler-research/xeus-cpp) | tag 0.10.0, `c90867a2fd48`, 2026-04-02 | `README.md`, `CMakeLists.txt`, `share/jupyter/kernels`, `LICENSE` |
| [compiler-research/CppInterOp](https://github.com/compiler-research/CppInterOp) | tag v1.9.0, `7b1258aa5d82`, 2026-03-26 | `README.md`, `lib/CppInterOp/CMakeLists.txt`, `LICENSE.txt` |
| xeus-cpp-lite deployment, [compiler-research.org/xeus-cpp](https://compiler-research.org/xeus-cpp/lab/index.html) | fetched 2026-11-01 | `xeus/kernels.json`, `xeus-cpp-wasm-host/empack_env_meta.json`, sizes of `bin/xcpp.{wasm,js,data}` and of `kernel_packages/cppinterop-1.9.0-h0b0027f_0.tar.gz` (downloaded, listed, the wasm inside extracted) |
| [emscripten-core/emscripten](https://github.com/emscripten-core/emscripten) | tag 4.0.9 (the `emscripten-abi` of that deployment), `bbf1caa6e24f`, 2025-05-17 | `src/lib/libdylink.js` |
| [root-project/cling](https://github.com/root-project/cling) | `6d15097652f2`, 2026-10-30 | `README.md`, `LICENSE.TXT`, `lib/Interpreter/*` file list, `DeclExtractor.h`, `DefinitionShadower.h` |
| [tetratelabs/wazero](https://github.com/tetratelabs/wazero) | v1.12.0 | `README.md`, `config.go`, `api/wasm.go`, `experimental/features.go`, `internal/wasm/instruction.go`; the `wazero` CLI v1.12.0 |
| [WebAssembly/wasi-libc](https://github.com/WebAssembly/wasi-libc) | `d9ad2d11991c`, 2026-10-30 | `README.md`, `LICENSE`, tree |
| gcc ([gcc-mirror/gcc](https://github.com/gcc-mirror/gcc)) | tag `releases/gcc-15.2.0`, `5115c7e447fc`, 2025-08-08 | `gcc/testsuite/gcc.c-torture`, `gcc/testsuite/gcc.dg/torture` (sparse checkout) |
| [c-testsuite/c-testsuite](https://github.com/c-testsuite/c-testsuite) | `5c7275656d75`, 2020-03-09 | `tests/single-exec`, `LICENSE`, `tests/LICENSE` |
| [csmith-project/csmith](https://github.com/csmith-project/csmith) | `0cdc710315cf`, 2026-03-01 | `README.md`, `COPYING` |
| LLVM 23.1.4 release assets | `gh api repos/llvm/llvm-project/releases/tags/llvmorg-23.1.4` | asset sizes |

Counting method. Clones live under `_bibliography/c-cpp/` of the workspace; nothing was built
from them. File and line counts are `go run ./tools/srccount` (this repository; a test plants
the files to skip: `.git`, symlinks, other extensions). The wasm verdicts are
`go run ./tools/wasmimports` (this repository; validates with wazero v1.12.0 and the
exception-handling, tail-call, extended-const and threads features on) and the `wazero compile`
CLI. One throw-away probe (§2.2) built a Go program around cc/v4 for `GOOS=js` and
`GOOS=wasip1` with Go 1.27.3.

## 1. What has to be executed

A C cell is a translation-unit fragment; a C++ cell the same in C++. The notebook contract is
the one of `bib/julia.md` §5.1 (copied from the Ruby page): one session per language, cells
define the top-level names they declare and use the names they read, only variables cross to Go,
errors are `file:line: message` at the notebook line, a worker with time and memory limits.
Two things C adds: **declarations persist across cells** (a `struct` or a function defined in
cell 3 is used in cell 7), and **top-level statements** are not C (C allows only external
declarations at file scope), so a cell like `int x = 2; printf("%d\n", x * x);` needs the
treatment REPLs give it (§2.5).

## 2. C

### 2.1 Reference

- **Language**: C11, text N1570. C17 is C11 plus defect fixes and is readable as the base text
  of N2310 (its final draft N2176 is password-protected). C23 (N3220) is out of v1: cc/v4's
  grammar has every C11 keyword (`_Alignas`, `_Atomic`, `_Generic`, `_Noreturn`,
  `_Static_assert`, `_Thread_local`: `parser.yy` lines 29-118) but no `_BitInt` or
  `constexpr`; its package comment says "a C99 compiler front end" and cites C99's n1256.
- **ABI**: wasm32, as clang defines it (`clang/lib/Basic/Targets/WebAssembly.h:91-113`): ILP32,
  `size_t` = `unsigned long`, `ptrdiff_t`/`intptr_t` = `long`, `long double` = 128-bit IEEE quad
  (`LongDoubleFormat = IEEEquad`), unaligned access allowed. Using the target's ABI, not the
  host's, makes `sizeof`, layout and overflow identical in the desktop kernel, the playground
  and the oracle.
- **Oracle compilers**: clang **23.1.4** (release binaries: `LLVM-23.1.4-macOS-ARM64.tar.zst`
  873 640 882 bytes, `LLVM-23.1.4-Linux-X64.tar.zst` 1 190 060 068 bytes) targeting
  `wasm32-wasip1` with wasi-libc as sysroot, the output run under the same wazero; and gcc
  15.2.0 natively for differential runs where the program has no target-dependent behaviour.
  clang→wasm32 is the primary oracle because it shares ABI, libc and engine with our kernel, so
  a difference is a compiler difference.

### 2.2 Existing pure-Go pieces

| piece | licence | what it is | measured | verdict |
|---|---|---|---|---|
| **modernc.org/cc/v4** v4.30.5 | BSD-3 (`LICENSE`) | preprocessor, parser, type checker; `Translate(cfg, sources)` returns a checked AST | `parser.go` 3 894, `check.go` 5 680, `ast.go` 5 287, `cpp.go` 3 096, `type.go` 3 007, `abi.go` 824 lines (21 788); tested on gcc-9.1.0 `gcc.c-torture`, CompCert-3.6 `test/c`, the SQLite amalgamation, tcc-0.9.27 `tests` (`all_test.go:788-796`, `1177-1185`) | **reuse** |
| modernc.org/ccgo/v4 v4.38.0 | BSD-3 | C → Go source; pointers become `uintptr` into memory allocated by modernc libc | `lib/` 51 files, 32 237 lines; its known-failure map lists 151 entries under gcc-9.1.0 `gcc.c-torture/execute/` on linux/amd64 and 121 on darwin/arm64, out of 1 700 `.c` files in that directory of ccorpus2 v1.6.0 (EXEC, BUILD and COMPILE failures together) | not for cells (§2.3) |
| modernc.org/libc v1.79.1 | BSD-3 (+ musl, `LICENSE-3RD-PARTY.md`) | "a partial reimplementation of C libc in pure Go", "to ccgo what the `runtime` package is to Go" (README) | 1 423 `.go` files, 2 187 951 lines (per-target generated code); `capi_*.go` targets: darwin, freebsd, illumos, linux, netbsd, openbsd, windows — **no `js` or `wasip1`** | not reusable: tied to ccgo's ABI and to the host address space |
| gotranspile/cxgo | MIT | C → readable Go, on cc **v3** | README: "experimental"; TCC 62/89, GCC 783/1236 (63 %); `static`, bitfields, C-layout unions, `setjmp` partial | no |
| hajimehoshi/goc | Apache-2.0 | "A C interpreter in Go … Work in progress" | lexer and preprocessor only (`internal/{lex,preprocess,ctype,_parse}`), last commit 2018 | no |

**cc/v4 runs without a host compiler and under wasm.** `NewConfig` execs the host `cc` to learn
predefined macros and include paths (`cc.go:317-333`, `exec.Command` at 669 and 706), but
`Config` is a plain struct (`cc.go:772-818`: `ABI`, `FS fs.FS`, `Predefined`, include paths,
`HasBuiltin`), so a kernel fills it itself, with headers from an `embed.FS`. The probe built
`&cc.Config{ABI: abi}` by hand, passed a 16-line predefined-macro block, `cc.Builtin` and a
two-function cell to `Translate`: 178 external declarations, the same natively and as a
`GOOS=wasip1` binary under `wazero run`. Linked into `GOOS=js` with `-ldflags='-s -w'` the binary
is 6 369 703 bytes against 2 460 764 for a hello world: **cc/v4 costs ≈ 3.9 MB raw** in the
playground. Gaps: `NewABI` knows no wasm target (`abi.go:775-784` refuses an unknown arch), but
`ABI` has exported `ByteOrder`, `Types` and `SignedChar` fields; the unexported `goos`/`goarch`
appear only in one error message (`check.go:2636`), so a wasm32 ABI can be built as a literal and
upstreamed later. There is no incremental API: `Translate` checks one whole translation unit
(§2.5 works around it).

### 2.3 Execution: three options

| | (1) tree/IR interpreter in Go | (2) transpile to Go (ccgo) | (3) compile to wasm32, run on wazero / the browser |
|---|---|---|---|
| runs in the wasm playground | yes | **no**: needs `go build` at cell time, and modernc libc has no `js`/`wasip1` target | yes: the browser runs the module natively through `WebAssembly.instantiate` from Go (`syscall/js`); wazero on the desktop |
| speed vs clang `-O2` (the fleet's reference) | an interpreter loop per operation; structurally an open defect | native Go speed | wazero compiler mode "often by order of magnitude (10x) or more" faster than its interpreter (README §Compiler), compiler on amd64 and arm64 only; in the browser, the JS engine's wasm tiers |
| host memory safety | yes if written in safe Go | **no**: `uintptr` arithmetic on the C heap; README: "Never pass a Go pointer to a transpiled C function" | yes: every access is bounds-checked against the guest's linear memory |
| C-level UB detection | possible (each object a Go slice with provenance) | no | no (UB stays inside linear memory, as in native C); optional bounds instrumentation later |
| libc | must be written against the interpreter's memory model | modernc libc | wasi-libc as is |
| what we write | interpreter + libc shim | nothing usable | code generator (cc/v4 AST → wasm), cell linker, ~20-40 WASI host calls |
| shared with C++ | no | no | **yes**: clang-repl also emits wasm modules loaded into one memory (§3.3) |

**Recommendation: (3).** It is the only option that is fast in both homes, memory-safe for the
host, uses an unmodified libc, and gives C++ its runtime for free. The cost is a back end:
structured control flow from C's `goto`/`switch` (LLVM's wasm back end needs two passes for it,
`WebAssemblyCFGStackify.cpp` "inserts BLOCK, LOOP, TRY, and TRY_TABLE markers" and
`WebAssemblyFixIrreducibleControlFlow.cpp` turns "multiple-entry loops" into single-entry ones;
the whole `llvm/lib/Target/WebAssembly` is 113 files, 36 205 lines), `setjmp`/`longjmp` through
wasm exception handling (wazero v1.12.0 implements the exnref form, `try_table` = opcode `0x1f`,
`throw_ref` = `0x0a`, behind `experimental.CoreFeaturesExceptionHandling`), varargs, bitfields,
`long double` as software binary128. Codegen quality is measured against clang's own
`wasm32-wasip1` output **in the same engine**, which isolates our compiler from the engine.

### 2.4 libc: wasi-libc, compiled once

wasi-libc is "a libc for WebAssembly programs built on top of WASI system calls" (README),
multi-licensed Apache-2.0 WITH LLVM-exception / Apache-2.0 / MIT, its `libc-top-half/musl/` MIT,
`dlmalloc/` CC0 (`LICENSE`). Size of the C to compile: musl `src` 1 828 files, 80 755 lines
(`math` 499 files, 21 342 lines; `stdio` 118 files, 4 784 lines), `libc-bottom-half` 210 files,
37 336 lines. Plan: build it with clang 23.1.4 into a pinned `libc.wasm` + headers (a release
artefact, checksummed, embedded), the kernel providing the `wasi_snapshot_preview1` imports
itself (stdout/stderr to the cell output, a read-only in-memory FS, clock, random; no sockets).
Once our back end compiles musl, rebuild `libc.wasm` with it and drop clang from the build: the
same dogfooding gate the julia plan uses for Base. `libm` results are then musl's on every
platform, which is what the oracle (clang + wasi-libc) produces too.

### 2.5 Notebook semantics: incremental definitions

- **Prior art.** clang-repl: "an interactive C++ interpreter that allows for incremental
  compilation", input compiled to "an AST representation", optionally transformed, lowered to IR
  and executed (`ClangRepl.rst:5-48`), with `%undo` (`:213-221`); its wasm executor emits one
  object per partial translation unit, links it with `wasm-ld -shared --import-memory
  --stack-first --allow-undefined` and `dlopen`s it (`clang/lib/Interpreter/Wasm.cpp:78-147`).
  Cling wraps statements in a function and hoists declarations back to file scope, with an
  init-order rewrite so that `int i = 1; i++; int j = i;` gives `j == 2`
  (`DeclExtractor.h:64-78`), and lets a name be redefined by nesting each definition in an
  inline namespace and hiding the previous one from lookup (`DefinitionShadower.h:29-49`).
- **Design for C.** The session keeps a symbol table of what earlier cells defined. Each cell is
  translated by cc/v4 as one TU made of (a) a synthesized prelude — every earlier `typedef`,
  `struct`/`union`/`enum`, macro, and `extern` declarations of earlier variables and functions —
  then (b) the cell's external declarations, then (c) its top-level statements wrapped in
  `void __cell_N(void) { … }`, called after instantiation. Splitting (b) from (c) is the cling
  problem; v1 rule: a top-level chunk that parses as an external declaration is one (cc/v4 is the
  judge, the typedef names of the prelude in scope), the rest is wrapped, and file-scope objects
  declared among statements are hoisted with cling's init-order rewrite.
- **Redefinition** (re-running a cell): C has one namespace and no inline namespaces; a new
  definition gets a fresh link name (`f` → `f__c7`), the session maps `f` to the newest, and
  later cells bind to it. Earlier cells keep calling the version they were compiled against
  (documented, as cling documents qualified access to old definitions).
- **Linking.** All cells share one memory and one funcref table owned by the session; a cell
  module imports `env.memory`, `env.__indirect_function_table`, its data base and table base,
  and the symbols it uses. Unlike emscripten (§3.3) we control both sides, so imports are
  resolved at instantiation from the session's table: no lazy proxies.
- **Values to Go.** After a cell runs, each file-scope variable it defined crosses by its C type:
  integers and floating types to the same-width Go type, `_Bool` to `bool`, `char*` to `string`
  on request, arrays of arithmetic type to `*ndarray.Array` — a **copy**, or a view built on
  wazero's `Memory.Read`, which returns a slice whose writes "are visible to Wasm" but which
  "disconnects on any capacity change" (`api/wasm.go:630-648`), so a view is valid only while the
  memory does not grow; structs and pointers stay opaque handles.

### 2.6 Test corpora

| corpus | revision | size | licence | use |
|---|---|---|---|---|
| gcc `gcc.c-torture/execute` | gcc 15.2.0 | 1 890 `.c`, 79 491 lines | GPL (gcc's testsuite) — run as data, never ported | exit status 0 = pass; `abort()` = fail |
| gcc `gcc.c-torture/compile`, `gcc.dg/torture` | gcc 15.2.0 | 1 988 and 1 962 `.c` | same | must compile (compile/) or pass (dg/torture, with its dg-directives filtered) |
| the same at gcc-9.1.0, in ccorpus2 v1.6.0 | v1.6.0 | 1 700 `.c` in `execute/` | same | the baseline cc/v4 and ccgo already report against |
| c-testsuite `tests/single-exec` | `5c7275656d75` | 220 `.c`, 6 118 lines | test harness MIT; each test's licence in its `.otag` (`tests/LICENSE`) | exit status and expected output |
| CompCert-3.6 `test/c`, tcc-0.9.27 `tests2` | in ccorpus2 | — | per directory | programs with expected output |
| csmith | `0cdc710315cf` | generator | BSD (`COPYING`) | "outputs C programs free of undefined behaviors" (README): differential fuzzing, our kernel vs clang→wasm32, checksum printed by each program; ccgo already drives it (`lib/all_test.go:41-78`, `818-1026`) |

"Correct" means: same exit status, same stdout/stderr bytes as the clang→wasm32-wasip1 oracle in
the same engine. Report per corpus: pass / fail-only-ours / fail-both, by record (file), not by
total, as `bib/julia.md` §8 prescribes.

## 3. C++

### 3.1 Why our own front end is out of scale

Measured at `llvmorg-23.1.4` (`.cpp`, `.h`, `.td`, `.def` files):

| part | files | lines |
|---|---|---|
| `clang/lib/Sema` | 95 | 333 996 |
| `clang/lib/AST` | 153 | 197 518 |
| `clang/lib/Parse` | 19 | 49 811 |
| `clang/lib/Lex` | 27 | 32 214 |
| front end total (Sema + AST + Parse + Lex) | 294 | **613 539** |
| `clang/lib/CodeGen` | 148 | 211 641 |
| `libcxx/include` (all files: the headers have no extension) | 1 688 | 330 575 (12 685 851 bytes) |
| `clang/lib/Interpreter` + `include/clang/Interpreter` + `tools/clang-repl` | 23 | 5 158 |

Sema alone is 15× the 21 788 lines of cc/v4's C front end. A notebook user's C++ starts with
`#include <vector>` and `std::cout`, i.e. templates, overload resolution and the 330 575 lines
of libc++ headers that must be parsed, instantiated and compiled. The interesting number is the
last row: clang's incremental layer is 5 158 lines **on top of** a complete compiler. The work is
in the compiler, not in the REPL.

### 3.2 Prior art

- **Cling** (CERN ROOT): "an interactive C++ interpreter, built on top of Clang and LLVM",
  "implemented as a small extension" (README:9-14); it builds against a patched LLVM branch,
  `root-project/llvm-project` `cling-latest` (README:45-50); `lib/Interpreter` is 62 `.cpp`/`.h`
  files, 24 750 lines
  (`DeclExtractor`, `DefinitionShadower`, `NullDerefProtectionTransformer`, `ValuePrinter`,
  `IncrementalJIT`, …). Licence "UI/NCSAOSL" or "LGPL", at the user's choice (`LICENSE.TXT:12-13`).
- **clang-repl** (upstream LLVM) is where Cling's parts are moving: "Clang-Repl allows to move
  parts of Cling upstream" (`ClangRepl.rst:11-16`); Apache-2.0 WITH LLVM-exception.
- **CppInterOp** v1.9.0 "exposes API from Clang and LLVM in a backward compatible way … compiler
  as a service" (README); Apache-2.0 WITH LLVM-exception; built for emscripten as
  `-s SIDE_MODULE=1` (`lib/CppInterOp/CMakeLists.txt:133-140`).
- **xeus-cpp** 0.10.0 (BSD-3): a Jupyter kernel on CppInterOp; ships kernels `xc11`, `xc17`,
  `xc23`, `xc23-omp`, `xcpp17`, `xcpp20`, `xcpp23`, `xcpp23-omp` (`share/jupyter/kernels`) — C
  and C++; the wasm kernel is linked with `-fwasm-exceptions` and the sysroot headers preloaded
  (`CMakeLists.txt:84-85`, `432-448`).

### 3.3 Options

**(a) A declared C++ subset on the C kernel** (classes, references, overloading, no templates).
Rejected: without templates there is no `std::vector`, `std::string` or `<iostream>`, so it is a
different language under the C++ name, and every step towards templates is a step into the
333 996 lines of Sema (two-phase lookup, deduction, SFINAE, concepts). It also breaks the fleet
rule that a port has an equality oracle: no compiler implements "our subset".

**(b) clang-repl compiled to wasm, hosted by Go.** What the published xeus-cpp-lite runs
(`kernels.json`: `xc23`, `xc11`, `xc17`, `xcpp17`, `xcpp20`, `xcpp23`; environment pinning
cppinterop 1.9.0, clang-resource-headers 21.1.8, emscripten-abi 4.0.9), measured on 2026-11-01:

| artefact | bytes |
|---|---|
| `bin/xcpp.wasm` (main module: xeus kernel) | 2 252 745 |
| `bin/xcpp.js` (emscripten JS runtime) | 2 217 806 |
| `bin/xcpp.data` (preloaded sysroot `include/` + clang resource dir) | 26 074 947 |
| `kernel_packages/cppinterop-1.9.0-h0b0027f_0.tar.gz` (holds the module twice, as `.so` and `.so.21.1`) | 45 636 892 |
| `lib/libclangCppInterOp.so` inside it: clang + lld + CppInterOp as one wasm side module | **71 256 623** |

How the pieces load. Each cell is compiled to an object, linked by the in-module `wasm-ld` into a
**shared** wasm module importing the main memory, and `dlopen`ed (`Wasm.cpp:120-147`).
Emscripten's `dlopen` is `libdylink.js`: it reads the `dylink.0` custom section
(`libdylink.js:443-506`), allocates the module's data in the shared memory and grows the shared
table by its size (`:638-667`), and instantiates it with `env` served by a proxy that answers
`__memory_base`/`__table_base` and returns lazily-resolving stubs for symbols defined by modules
loaded later (`:694-745`), plus `GOT.mem`/`GOT.func` proxies (`:746-748`); then it relocates the
exports and registers `EM_ASM` bodies (`:753-816`). **A Go host must provide**: that loader (the
`dylink.0` parser, memory and table allocation, GOT globals, symbol resolution across modules —
without JS proxies, i.e. by reading each module's import section before instantiating it); the
imports `xcpp.js` provides today (syscalls, the virtual FS holding `xcpp.data`, clocks); and
wasm exception handling.

Measured on the 71 MB module: `wazero compile` (v1.12.0, default features) refuses it at
`import[4250] tag[env.__c_longjmp]: tag imports require exception handling feature`; with
exception handling enabled (`tools/wasmimports`), validation fails at
`function[70505] export["_ZN4llvm20CrashRecoveryContext9RunSafelyENS_12function_refIFvvEEE"]:
legacy exception handling instruction 0x6 not supported; recompile with wasm-opt
--translate-to-exnref`. So the module uses the legacy `try` encoding; wazero has only the exnref
form. A Binaryen `wasm-opt --translate-to-exnref` pass at packaging time is the documented
remedy; not tried here.

**(c) Delegate to a native clang-repl.** Rejected by the doctrine: it needs an installed LLVM
(the 23.1.4 release is 873 640 882 bytes compressed for macOS arm64), an `exec` of a host binary,
does not exist in the wasm playground, and JIT-compiles user code into native code running with
the user's privileges, outside any sandbox. It is not embeddable and not `CGO_ENABLED=0` in
spirit (the dependency is a native compiler instead of a C library).

### 3.4 Recommendation for C++

**(b)**, staged on go-embedded-c:

1. go-embedded-c M3 delivers the cell runtime (shared memory and table, cell linker, WASI host,
   limits). Its loader speaks our own linking convention, chosen as a subset of emscripten's
   `dylink.0` so that the same code can later load emscripten side modules.
2. go-embedded-cpp M0 is a **gate**, measured before any further work: (i) `wasm-opt
   --translate-to-exnref` on the pinned module, then wazero validation; (ii) wazero compile time
   and resident memory for the 71 MB module on one core; (iii) time to first result of
   `#include <vector>` + `std::cout` on the desktop; (iv) download size for the playground, where
   the module is handed to the browser engine (wazero-in-Go-wasm would interpret clang, an
   interpreter inside an interpreter). Budgets: start-up within the Julia/Ruby cores' order on
   the same machine (`bib/julia.md` §6), download stated next to xeus-cpp-lite's.
3. If the gate passes: go-embedded-cpp = the emscripten-side-module host (dylink loader, the
   imports `xcpp.js` provides, the sysroot as an `fs.FS`), packaging of a pinned CppInterOp build
   (rebuilt from source in CI, not downloaded), and the notebook bridge. If it fails, `cpp` cells
   are declared unsupported in the playground and the org stays empty; nothing else in the plan
   depends on it.

Reasons: it executes real clang, so C++ semantics come from the reference and the equality
oracle is clang itself; it reuses go-embedded-c's runtime; licences are permissive (Apache-2.0
WITH LLVM-exception, BSD-3); and it is the only option that reaches the playground.

## 4. Validation and performance

- C: the clang 23.1.4 → `wasm32-wasip1` + wasi-libc oracle in the same wazero (§2.1), driven
  from a Go harness that compiles each corpus file with both compilers and compares exit status
  and output bytes; csmith for differential fuzzing; `go test -fuzz` on the cell splitter and on
  the module loader.
- C++: clang-repl is the implementation; validation is of the host: xeus-cpp's own tests
  (`test/`), run through the Go host, and the C++ cells of a go-notebook example notebook.
- Performance, per the fleet rule (reference tool, one core, interleaved rounds): C cells
  against native `clang -O2` and against clang→wasm32 under the same engine; the first gap is the
  engine, the second our code generator, each an open defect when slower. Nothing was timed for
  this file.

## 5. Proposed home and plan

Effort scale of `bib/`: S days, M weeks, L months, XL several months.

| work item | org/repo | depends on | size | priority |
|---|---|---|---|---|
| M0 oracle harness: pinned clang 23.1.4 + wasi-libc sysroot, corpus runner (gcc torture 15.2.0, c-testsuite, ccorpus2), golden outputs, csmith driver | go-embedded-c/c (`cmd/cref`) | — | S | 1 |
| M1 front end: cc/v4 with a wasm32 `ABI` literal, embedded wasi-libc headers, `HasBuiltin`; cell splitter and prelude (§2.5); upstream the ABI to cc/v4 | go-embedded-c/c | M0 | M | 1 |
| M2 code generator: cc/v4 AST → wasm32 (expressions, structured control flow incl. `goto`/`switch`, varargs, bitfields, `long double` binary128, `setjmp` via exnref) | go-embedded-c/c (`wasmgen`) | M1 | L | 1 |
| M3 cell runtime: shared memory/table session, cell linker (dylink.0 subset), WASI host, limits; `libc.wasm` built with clang, then with M2 | go-embedded-c/c (`session`) | M2 | M | 1 |
| M4 notebook bridge: value table §2.5, `file:line` errors, go-notebook ` ```c ` cells (desktop on wazero, playground on the browser engine) | go-embedded-c/c + go-notebook | M3 | M | 1 |
| M5 conformance: gcc torture execute ≥ the pass rate of ccgo on the same files, c-testsuite, csmith clean for a fixed budget; performance table | go-embedded-c/c | M4 | M | 2 |
| site, landing, brand | go-embedded-c/{docs, go-embedded-c.github.io, .github, brand} via `ghnew` | M4 | S | 2 |
| C++ M0 gate (§3.4) | go-embedded-cpp/cpp (spike branch) | go-embedded-c M3 | S | 2 |
| C++ M1 host: emscripten side-module loader, `xcpp.js` imports in Go, sysroot `fs.FS`, pinned CppInterOp build in CI | go-embedded-cpp/cpp | gate passed | L | 3 |
| C++ M2 notebook bridge (` ```cpp ` cells) | go-embedded-cpp/cpp + go-notebook | C++ M1 | M | 3 |

Total for C: on the order of one rbgo milestone wave; M2 dominates. Repository names follow
`bib/julia.md` §9 (`go-embedded-julia/julia`): `go-embedded-c/c`, `go-embedded-cpp/cpp`.

## 6. Risks and open questions

- **Security of user C in the notebook.** A wasm module cannot address host memory: every load
  and store is checked against its linear memory, so a buffer overflow in a cell corrupts the
  session, not the kernel. Bounds the runtime must set, all available in wazero:
  `WithMemoryLimitPages` (`config.go:65`; a cell's `malloc(1<<31)` must fail, not grow the
  host), `WithCloseOnContextDone` with a context deadline — "especially useful when one wants to
  run untrusted Wasm binaries since otherwise, any invocation … can potentially block the
  corresponding Goroutine forever" (`config.go:149-170`); recursion is bounded by the wasm
  stack. Host imports are the attack surface: WASI calls get a read-only in-memory FS, no
  sockets, no `exec`, and every pointer/length pair a guest passes is validated before
  `Memory.Read`. The kernel itself parses untrusted text (cc/v4) and untrusted modules (the
  loader): both are fuzzed, the corpus committed, with size bounds (source bytes, macro
  expansion depth, `#include` depth, module size) asserting a refusal under `GOMEMLIMIT`.
  An interpreter (option 1) would have given UB *diagnostics* but not more host safety than this.
- **A crashed session.** UB inside a cell can corrupt the session's memory silently, as native
  C does; the notebook offers "restart C session" and replays cells, like a Jupyter kernel
  restart. Optional bounds-check instrumentation in the code generator is a later item.
- **cc/v4 is one maintainer's front end**, tested by its author against ccgo's needs; GNU
  extensions are its priority (`UnsignedEnums`, `IgnoreNegativeShiftAmounts` knobs). Its libc
  README warns that minor versions break the ccgo/libc contract; we do not use ccgo or modernc
  libc, only cc, but pin it exactly and run M0's corpus on every bump.
- **Code-generator correctness** is the main risk of option 3: structured control flow from
  arbitrary `goto`, `switch` into loops (Duff's device), and the C conversion rules. The oracle
  and csmith exist for exactly this.
- **Exception handling**: our `setjmp` uses wazero's experimental exnref support; the browser
  engines' support of exnref must be checked per browser before the playground relies on it.
- **C++ size**: 71 256 623 bytes of clang module plus 26 074 947 bytes of sysroot data, and the
  legacy-EH refusal. Unmeasured: wazero's compile time and memory for it, and whether
  `--translate-to-exnref` yields a module wazero accepts. These decide §3.4.
- **Unverified in this file**: every speed statement (nothing was timed); the pass rate of
  cc/v4 itself on gcc torture (its tests were read, not run); ccgo's pass rate is derived from its
  known-failure map, not run; csmith not built; c-testsuite not run; wasi-libc not built; the
  browser-engine path for cells (`WebAssembly.instantiate` from Go) not prototyped; whether
  wazero lets a host grow a table shared by several modules (to check in M3); the papers behind
  the LLVM structuring passes were not opened.

Verified on 2026-11-01: all file and line counts by `tools/srccount` over the revisions of §0;
the two wasm refusals by the `wazero` CLI v1.12.0 and `tools/wasmimports`; artefact sizes by
HTTP downloads from compiler-research.org; release asset sizes and org state by `gh api`; the
cc/v4 probe natively and under `wazero run` as `GOOS=wasip1`; licence files read in each clone.
