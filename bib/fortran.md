# Fortran — bibliography and feasibility plan for go-embedded-fortran

Written 2026-11-13. Facts only, each with a source; counts come from the sources, not from
estimates. Goal: ` ```fortran ` cells in go-notebook and its WebAssembly playground, backed by Go
code built with `CGO_ENABLED=0`, embeddable as a library and runnable as wasm, sharing as much code
as possible with go-ndarray, go-scikit and go-fft. The organisation
[go-embedded-fortran](https://github.com/go-embedded-fortran) was created on 2026-11-13
(`gh api orgs/go-embedded-fortran` `created_at` 12:04:12Z) and holds no repository
(`gh repo list` returned `[]`).

**Decisions.**

- **Front end: ours, in Go.** Lexer (free and fixed source form), parser, semantic analysis and the
  cell splitter are written in Go. No existing front end can be reused as is: LFortran's parser is a
  Bison GLR grammar with 189 shift/reduce and 189 reduce/reduce conflicts (`parser.yy:16-18`), not a
  goyacc (LALR) grammar; flang's front end is 114 477 lines of C++ (§3.2); the only Go Fortran
  parser with activity is F77-F90 and transpiler-oriented (§3.4). LFortran (BSD-3) is the design
  reference for the semantic tree and the cell semantics, flang-rt (Apache-2.0 WITH LLVM-exception)
  for the runtime. The front end is the largest item of the plan and is common to both back ends.
- **Back end: not decided here.** Two candidates, (A) a pure-Go evaluator in the manner of
  go-embedded-julia, where Fortran arrays *are* `*ndarray.Array` with column-major strides, and (B)
  compiling to wasm32 on go-embedded-c's back end and cell runtime. §4 states the measured gate,
  run at M1 on four kernels, and its decision rule, written before any measurement.
- **Libraries once, in Go.** Whichever back end wins, BLAS/LAPACK calls, the array intrinsics
  (`matmul`, `dot_product`, `transpose`, `sum`, …) and FFTs go to go-ndarray, go-ndarray/linalg and
  go-fft; formatted and list-directed I/O is one Go implementation. Nothing is compiled from
  reference LAPACK or from a Fortran runtime library into the kernel (§6).

## 0. Sources read, and how the counts were made

| source | revision read | what was used |
|---|---|---|
| J3 [18-007r1](https://j3-fortran.org/doc/year/18/18-007r1.pdf) | "WD 1539-1, J3/18-007r1 (F2018 Interpretation Document)", cover dated 9 October 2018, page headers 2018-08-28 | cover and table of contents (clauses 1-13 to sub-clause level): the F2018 text, base of the v1 subset |
| J3 [23-007r1](https://j3-fortran.org/doc/year/23/23-007r1.pdf) | "J3/23-007r1 (Draft Fortran 2023)", 13 June 2023 | cover and top-level clause list only: what v1 leaves out |
| [lfortran/lfortran](https://github.com/lfortran/lfortran) | `4f2cd049c87f`, 2026-11-13 (latest release v0.69.0, 2026-11-10) | `LICENSE`, `README.md`, `src/lfortran/parser/{parser.yy,tokenizer.re,fixedform_tokenizer.cpp}`, `src/lfortran/semantics`, `src/libasr/{ASR.asdl,pass,codegen,runtime}`, `src/runtime`, `doc/src/{interactive_cells,interactive_evaluation,jupyterlite}.md`, `integration_tests/CMakeLists.txt` |
| LFortran JupyterLite deployment, [lfortran.github.io/lfortran](https://lfortran.github.io/lfortran/) | fetched 2026-11-13 (`last-modified` 10:01:52 GMT) | `xeus/kernels.json`; sizes of `xeus/wasm-host/bin/xlfortran.{wasm,js,data}`; `xlfortran.wasm` downloaded and validated |
| [llvm/llvm-project](https://github.com/llvm/llvm-project) | tag `llvmorg-23.1.4`, `7a941dfb2877`, 2026-10-19 (the C oracle's LLVM, `bib/c-cpp.md` §0) | `flang/lib/{Parser,Semantics,Evaluate,Lower,Optimizer}`, `flang/include`, `flang/lib/Optimizer/CodeGen/Target.cpp`, `flang/docs/MeetingNotes/2025/2025-12-03.md`, `flang-rt/{lib,include,LICENSE.TXT}` (sparse checkout); `clang/cmake/caches/Release.cmake` at the tag; release asset list (`gh api …/releases/tags/llvmorg-23.1.4`) |
| gcc ([gcc-mirror/gcc](https://github.com/gcc-mirror/gcc)) | tag `releases/gcc-15.2.0`, `5115c7e447fc`, 2025-08-08 | `gcc/testsuite/{gfortran.dg,gfortran.fortran-torture}`, `gcc/fortran`, `libgfortran` (licence header of `io/write.c`) |
| [Reference-LAPACK/lapack](https://github.com/Reference-LAPACK/lapack) | `0c51225`, 2026-11-11 (latest release v3.12.1) | `LICENSE`, `SRC`, `BLAS/SRC`, `TESTING`, `BLAS/TESTING` |
| [llvm/llvm-test-suite](https://github.com/llvm/llvm-test-suite) | `19b349c`, 2026-11-12 | `Fortran/{gfortran,UnitTests,SNAP,cray}`, `Fortran/gfortran/README.md`, `Fortran/UnitTests/fcvs21_f95/CMakeLists.txt` |
| [soypat/go-fortran](https://github.com/soypat/go-fortran) | `04860ef`, 2026-07-10 | `README.md`, `LICENSE` (via `gh api`), `lexer.go`, `parser.go`, `repl.go` (declarations) |
| [Konstantin8105/f4go](https://github.com/Konstantin8105/f4go) | `568775f`, 2023-08-17 | `README.md` |
| GitHub search `fortran language:Go` | 2026-11-13 | the 17 repositories returned |
| [go-ndarray/ndarray](https://github.com/go-ndarray/ndarray) | v0.40.3, `d8ed921`, 2026-11-12 | `ndarray.go`, `storage.go`, `linalg.go`, `linalg/{linalg,blas}.go`, `fft/fft.go`; issue #76 |
| [go-embedded-julia/julia](https://github.com/go-embedded-julia/julia) | `71691db`, 2026-11-10 (latest release v0.12.0) | `README.md`, `interp/{eval,prepare}.go` (declarations) |
| [go-embedded-c/c](https://github.com/go-embedded-c/c) | `23b6d30`, 2026-11-11 (latest release v0.2.1) | `README.md`, `docs/design-m2.md`, `wasmgen/{func,structure}.go`, `internal/wasm` |
| [go-notebook/notebook](https://github.com/go-notebook/notebook) | v0.14.0, `14320ea`, 2026-11-10 | `analysis/lang.go` |
| [go-rouge/rouge](https://github.com/go-rouge/rouge) | `1f4653a`, 2026-11-08 (latest release v0.9.0) | `README.md`, `lexers.go`, `lexers_fortran.go` |
| [go-fft/fft](https://github.com/go-fft/fft) | v0.30.1 (`gh api`: BSD-3-Clause) | description only; its use by `go-ndarray/fft` read in `fft/fft.go:1-17` |
| this repository | `bib/c-cpp.md`, `bib/julia.md` (§ LinearAlgebra l. 218-228, array crossing l. 329-338), `bib/julia-runtime.md` §6.2, `bib/r.md` §1.2, PR #14 | the sibling plans this one builds on |

Counting method. Clones live under `_bibliography/fortran/` of the workspace (LLVM and gcc as added
paths of the sparse clones in `_bibliography/c-cpp/`); nothing was built from them. File and line
counts are `go run ./tools/srccount` (this repository), with the extensions stated per row. The
wasm verdict is `go run ./tools/wasmimports`. Label counts in LFortran's test registry are
`srccount --contains`, i.e. **lines** containing the text, stated as such.

## 1. What has to be executed

The notebook contract is the one of `bib/julia.md` §5.1 and `bib/c-cpp.md` §1: one session per
language, cells define the names they declare and use the names they read, variables cross to Go,
errors are `file:line: message` at the notebook line, a worker with time and memory limits.
Fortran adds three things:

- **Program units persist.** A `module`, `subroutine` or `function` defined in cell 3 is used in
  cell 7 (by `use` for a module, by name for an external procedure).
- **Top-level statements are not Fortran.** A program is made of program units (F2018 5.2); a cell
  like `x = 2.0; print *, x**2` needs the REPL treatment LFortran gives it (§7).
- **Arrays are column-major and 1-based by default**, and the scientific reason to write Fortran in a
  notebook is array code and LAPACK-style loops, so the crossing of arrays (§6) and the speed of
  loops (§4) decide the design.

### 1.1 The v1 language subset

Text: F2018 (18-007r1). Proposed v1 = F2018 **without**:

- coarrays and images: 5.4.7-5.4.8 (coarray, established coarrays), 9.4.3 (coindexed objects), 9.6
  (image selectors), 11.6 (image execution control: `SYNC *`, `EVENT`, `FORM TEAM`, `LOCK`), 11.5
  `FAIL IMAGE`;
- teams: 11.1.5 `CHANGE TEAM`, 11.6.6, 11.6.9;
- parallel semantics of `DO CONCURRENT` (11.1.7): accepted and run sequentially;
- asynchronous I/O (12.7 `WAIT`, 8.5.4 `ASYNCHRONOUS`), submodules (5.2.5), derived-type parameters
  (7.5.3), user-defined derived-type I/O (13.7.6) — later milestones, not v1;
- interoperability with C (clause 18) beyond accepting `bind(c)` on names: there is no C to call in a
  cell session unless option B shares go-embedded-c's session (§4.2).

**With** fixed source form (6.3.3): reference LAPACK and BLAS are fixed form (`SRC` holds 2 188 `.f`
files, 922 119 lines, and 18 `.f90` files, 10 354 lines), and so is the NIST F77 suite. The deleted
features of Annex B that FCVS uses are accepted as gfortran accepts them: its llvm-test-suite copy
notes "Tests make use of arithmetic IF which is now a deleted feature in Fortran 2008" and omits
nine tests for features deleted in F95 (`fcvs21_f95/CMakeLists.txt:21-31`). F2023 (23-007r1) is out
of v1: only its cover and clause list were read.

## 2. Reference

- **Language**: F2018, J3/18-007r1, clauses 6 (lexical tokens, free 6.3.2 and fixed 6.3.3 source
  form), 7-11 (types, declarations, data objects, expressions, execution control), 12-13 (I/O
  statements and editing), 14-16 (program units, procedures, intrinsics), 17 (IEEE), Annex B
  (deleted and obsolescent features).
- **Oracle compilers**:
  - **gfortran 15.2.0** natively, the primary output oracle: it is what LFortran's and
    llvm-test-suite's reference outputs are made with (6 232 lines of LFortran's test registry carry
    the `gfortran` label; FCVS's references were "generated with gfortran", `CMakeLists.txt:4-13`).
    GPL; run as a program, never ported.
  - **flang 23.1.4**, the LLVM of the C oracle: in the Linux release tarballs, not the macOS one —
    `Release.cmake` at the tag appends `flang` to the projects and `flang-rt` to the runtimes only
    `if (NOT ${CMAKE_HOST_SYSTEM_NAME} MATCHES "Darwin")`. It **cannot be the wasm32 oracle**: its
    target table (`flang/lib/Optimizer/CodeGen/Target.cpp:1902-1964`) has cases for x86, x86_64,
    aarch64, ppc, ppc64(le), sparc(v9), riscv64, amdgpu, nvptx64 and loongarch64, and anything else
    reaches `TODO(…, "target not implemented")`. Flang+wasm exists only as outside work (a FOSDEM
    2025 talk, cited in `flang/docs/MeetingNotes/2025/2025-12-03.md:49`).
  - **LFortran** as a second front-end opinion (its ASR is a typed semantic tree) and for its wasm
    path (§3.1), not as the equality oracle: it is a young compiler (its own CI treats third-party
    applications as "bug generators", `AGENTS.md` CI policy).

## 3. Existing implementations

### 3.1 LFortran

BSD-3 (`LICENSE`), "a modern open-source (BSD licensed) interactive Fortran compiler built on top
of LLVM" (`README.md:5-6`). Measured at `4f2cd049c87f` (`.cpp`/`.h`/`.yy`/`.re` unless stated):

| part | files | lines |
|---|---|---|
| `src/lfortran/parser` (Bison GLR grammar `parser.yy` 2 912, re2c `tokenizer.re` 1 060, `fixedform_tokenizer.cpp` 2 436) | 13 | 13 221 |
| `src/lfortran/semantics` (AST → ASR) | 12 | 56 018 |
| `grammar/AST.asdl` / `src/libasr/ASR.asdl` | 1 / 1 | 482 / 478 |
| `src/libasr/pass` (ASR → ASR passes) | 203 | 148 583 |
| `src/libasr/codegen/asr_to_llvm.cpp` | 1 | 31 666 |
| `src/libasr/codegen/asr_to_wasm.cpp` (+ `wasm_assembler.h` 420, `wasm_utils.cpp` 209) | 1 | 4 277 |
| `src/libasr/runtime/lfortran_intrinsics.c` (the C runtime: I/O, strings, math) | 1 | 21 050 |
| whole `src` (`.c .h .f90 .cpp .py .asdl .yy .re`) | 528 | 459 255 |
| `integration_tests` (`.f90 .f .f77`) | 7 153 | 303 283 |

- **Cell semantics (prior art).** "Every cell is compiled into its own `TranslationUnit`", symbol
  tables chained oldest first; "Cells are additive"; "Redeclaring shadows … Code compiled earlier
  keeps using what it resolved to then, and later cells resolve to the new declaration … The same
  holds for variables and modules" (`interactive_cells.md:11-53`). Cell-scope symbols are prefixed
  `__cell<N>_`, the first cell unprefixed; each cell runs through a wrapper
  `__lfortran_evaluate_<N>`; a `program` unit in a cell is called, not only compiled (`:86-136`).
  A cell containing `implicit none` disables implicit typing for later cells (`:55-68`). A bare
  expression cell returns its typed value, a statement cell `statement`, a declaration or `use`
  `none`; "A failed cell is not rolled back" (`interactive_evaluation.md:12`,
  `jupyterlite.md:390-391`).
- **wasm.** Two paths. (i) The direct back end `asr_to_wasm.cpp` (4 277 lines against 31 666 for
  LLVM) imports `wasi_snapshot_preview1.proc_exit`/`fd_write` and formats `print` output inline
  (`asr_to_wasm.cpp:99-107`, `174-195`, `238-352`); in the test registry, 246 lines carry the
  label `wasm` (158 + 88, `" wasm "` and `" wasm)"`) against 6 525 with `llvm`, out of 6 651
  `RUN(NAME` lines — the direct back end covers a small part of the language. (ii) LLVM's
  wasm target (`llvm_wasm`, 2 578 lines). The JupyterLite kernel is LFortran *itself* compiled to
  wasm with Emscripten, with LLVM in the host environment (`jupyterlite.md:34-35`, `143`).
- **The deployed kernel**, measured 2026-11-13: `xlfortran.wasm` **78 103 511** bytes,
  `xlfortran.js` 628 699, `xlfortran.data` 194 907. `tools/wasmimports` refuses it as wazero
  refuses the clang module of `bib/c-cpp.md` §3.3: `invalid function[0]
  export["__wasm_call_ctors"]: legacy exception handling instruction 0x6 not supported; recompile
  with wasm-opt --translate-to-exnref`.
- **Verdict.** Not embeddable as is (C++, LLVM, 78 MB, legacy EH). Portable by licence: its ASR
  design, its pass list and `lfortran_intrinsics.c` may be ported to Go. Its grammar cannot be fed
  to goyacc (GLR with 378 declared conflicts): a Go port is a hand-written parser.

### 3.2 LLVM flang and flang-rt (llvmorg-23.1.4)

| part | files | lines |
|---|---|---|
| `flang/lib/Parser` | 34 | 20 174 |
| `flang/lib/Semantics` | 90 | 67 875 |
| `flang/lib/Evaluate` (constant folding, intrinsic table) | 35 | 26 428 |
| front end total (Parser + Semantics + Evaluate) | 159 | **114 477** |
| `flang/lib/Lower` + `flang/lib/Optimizer` (FIR/MLIR) | 211 | 157 702 |
| `flang/include` | 347 | 73 273 |
| `flang-rt/lib` + `flang-rt/include` (`.cpp .h .c`) | 201 | 39 790 |
| flang-rt I/O core: `edit-input`, `edit-output`, `format`, `format-implementation.h`, `io-api`, `io-stmt`, `descriptor-io`, `namelist`, `unit`, `internal-unit`, `external-unit`, `connection`, `buffer`, `io-error` | 14 | 9 219 |

flang and flang-rt are Apache-2.0 WITH LLVM-exception (`LICENSE.TXT`). flang-rt has separate
`matmul-*`, `dot-product-*`, `sum`, `product`, `extrema-*`, `findloc-*`, `transformational`,
`random` files (`lib/runtime` listing): the intrinsic families a runtime must cover. Verdict:
**the reference to port the runtime from** (licence compatible, C++ readable, one file per
family), and the native second oracle on Linux; not the wasm oracle (§2).

The flang community records the same goal as ours: "Proposing an Interactive Fortran Workflow with
Flang using Jupyter Notebooks … Something is available using LFortran, but it's not perfect … There
is clang-repl for this with C++ … RFC coming" (`MeetingNotes/2025/2025-12-03.md:39-45`).

### 3.3 gfortran 15.2.0

GPL; libgfortran is GPL with the GCC Runtime Library Exception (`libgfortran/io/write.c:9-19`).
`gcc/fortran` 65 files, 246 002 lines; `libgfortran` 914 files (`.c .h .F90 .f90`), 335 092 lines
(`io` 22 files, 23 908 lines). Oracle only.

### 3.4 Fortran in Go

The GitHub search `fortran language:Go` returned 17 repositories; three are Fortran tools:

| project | licence | what it is | state | verdict |
|---|---|---|---|---|
| soypat/go-fortran | BSD-3 | "Fortran AST parsing and tokenization support in Go. FOTRAN77 through Fortran 90 supported"; a transpiler to Go tested byte-for-byte against gfortran on `testdata/golden.f90` (README); `REPL.Eval` evaluates expressions (`repl.go:397`) | 50 `.go` files, 35 776 lines; pushed 2026-07-10; lexer `Lexer90` (free form, `&` continuation, `lexer.go:18`, `702-765`); no fixed-form code found in `lexer.go` | read in M0 as a candidate base for the lexer/parser; no semantic tree, no F2003+ |
| Konstantin8105/f4go | MIT | F77 → Go transpiler, LAPACK 3.8.0 as test input (README) | 30 files, 10 706 lines; last commit 2023-08-17 | no |
| gharib85/go-fortran | BSD-3 | same description and licence as soypat/go-fortran (not opened) | pushed 2025-12-12 | no |

No pure-Go Fortran interpreter or compiler exists.

### 3.5 The fleet pieces

- **go-embedded-julia/julia** executes lowered IR with an interpreter: `interp` is "a port of
  julia's `src/interpreter.c`"; `prepare` checks each SSA value, slot and label once, then
  `evalBody` "walks the statements" (`README.md:116`; `interp/eval.go:39-52` switches on the
  statement type, `evalValue` returns `any`, `:176-178`). Its speed against julia's own interpreter
  is "Not measured yet" (`README.md:124`): **the fleet has no measured interpreter speed** to
  extrapolate option A from.
- **go-ndarray v0.40.3.** `Array` is (data, shape, strides, offset), "stored row-major (C-order)"
  (`ndarray.go:16-18`, `34-43`). A column-major array is a `Transpose` view of the reversed shape,
  no copy (`ndarray.go:458`, `bib/julia-runtime.md` §6.2). But `FromData` and `FromSlice` copy
  (`ndarray.go:134-151`, `storage.go:131-143`), and `Data` returns the array's own storage only when
  it is contiguous row-major, a copy otherwise (`storage.go:145-155`): **a column-major array's
  buffer cannot be reached without a copy, and a caller's buffer cannot be adopted.** The missing
  constructor is issue [#76](https://github.com/go-ndarray/ndarray/issues/76) "Wrap: an array over a
  caller's buffer with given strides (zero-copy, for Julia column-major arrays)", open since
  2026-10-27. `MatMul` reads any 2-D strided view in place ("a transposed or sliced operand … is
  never materialised", `linalg.go:14-29`) and returns a row-major result; `linalg` functions take
  "a contiguous copy" of their operands (`linalg/linalg.go:127-135`). Its internal BLAS layer
  (`linalg/blas.go`: `gemm`, `trsm*`, `herk`, `larfg`, block reflectors, generic over
  `float64 | complex128`, LAPACK's `dlamch` constants at `:15-21`) is unexported.
- **go-embedded-c/c v0.2.1.** M0 and M1 done, M2 "in progress" (1 440 gcc torture files run as
  under clang), **M3 cell runtime "planned"** (`README.md:18-29`). The M2 design says "Each
  instruction works on typed virtual values" (`design-m2.md:37-38`); the code as it stands holds
  **already-encoded wasm bytes** in each basic block (`wasmgen/func.go:8-24`, `code []byte`), the
  function generator is keyed by cc/v4 nodes (`func.go:45-57`, `structure.go:54`), and 13 of the
  14 files of `wasmgen` reference `cc.`. The front-end-neutral parts are the Ramsey structurer
  (over blocks of bytes, `structure.go:10-17`) and the object encoder `internal/wasm`; both live in
  the C module, the encoder under `internal/`, which another module cannot import.
- **go-notebook v0.14.0.** A language is one registered `Language` (`analysis/lang.go:61-78`), with
  optional `CellAnalyzer`, `NameChecker`, `Highlighter`, `Runtime`, `TypeLearner`, `Isolated` —
  and **`Sequential`** (`:88-99`), made for C: "cells depend on the earlier cells of the same
  language, in notebook order … The graph calls Begin once … then AnalyzeNext for each of the
  language's cells in order". Fortran's `use` of an earlier cell's module is that case.
- **go-rouge v0.9.0** has a `fortran` lexer: "Fortran 2008 free form: kind parameters, `B`/`O`/`Z`
  literals, `.EQ.`-style operator words, two-word keywords, `FORMAT` state" (`README.md:97`,
  registered `lexers.go:35`). Fixed form is not listed.

## 4. The back end: the central decision

### 4.1 The two options

| | (A) pure-Go evaluator | (B) compile to wasm32 on go-embedded-c |
|---|---|---|
| model | go-embedded-julia: the front end's typed tree → closures over typed frames | go-embedded-c: typed tree → CFG → Ramsey structuring → relocatable object → cell linker |
| arrays | **are** `*ndarray.Array` (column-major strides over a `[]float64` the evaluator owns, wrapped by #76's `Wrap` per crossing, O(1)) | live in the session's linear memory; cross by copy, or by a view through wazero's `Memory.Read` that "disconnects on any capacity change" (`bib/c-cpp.md` §2.5) and still needs #76 |
| BLAS/LAPACK, `matmul`, FFT | direct Go calls into go-ndarray/linalg and go-fft | host imports: copy out of linear memory, call the same Go code, copy back |
| scalar loops | one closure call per operation; unmeasured (§3.5) | wazero compiler mode on amd64/arm64 (`bib/c-cpp.md` §2.3); the browser engine's tiers in the playground; interpreter elsewhere |
| playground | the evaluator is Go compiled to wasm: an interpreter inside wasm | the cell module is handed to the browser engine: native wasm speed |
| host safety | safe Go; Fortran bounds checks are free to add (Go slices are checked) | bounds-checked linear memory; out-of-bounds *within* the session's memory is silent, as in native Fortran |
| what we write | the evaluator (+ closure compiler) | Fortran lowering to the shared IR; and first, a front-end-neutral IR + structurer + encoder extracted from go-embedded-c (§3.5) |
| depends on | go-ndarray #76 | go-embedded-c M2 finished **and** M3 (cell runtime, today "planned"), plus the extraction |
| shared with | go-embedded-julia (the arrays, the BLAS mapping of `bib/julia.md` l. 218-228) | go-embedded-c (and go-embedded-cpp through it) |

Both use the same front end, the same Go runtime library (I/O, intrinsics, §5) and the same Go
numerical libraries; they differ in who executes scalar code and where arrays live.

### 4.2 The gate (M1), decided in advance

**When.** As soon as the front end parses and checks the four kernels below (M1). Option A's
contender is a throw-away evaluator written for the gate (≤ S); B's contenders are made by the
oracles, so the gate does not wait for go-embedded-c M3.

**Kernels** (double precision, each with a checksum printed, inputs from a fixed seed):

| | kernel | size |
|---|---|---|
| K1 | scalar loop: `do i = 1, n; s = s + 1d0/real(i,8)**2; end do` | n = 10⁸ |
| K2 | 2-D five-point Jacobi stencil, column-major order (`j` outer, `i` inner), 100 sweeps | 1024 × 1024 |
| K3 | matmul as a triple loop `c(i,j) = c(i,j) + a(i,k)*b(k,j)` (loop order `j,k,i`); and separately `c = matmul(a, b)` | 512 × 512 |
| K4 | whole-array expression `a = b*c + d` | n = 10⁷ |

**Contenders.**

1. **A**: the Go evaluator, closures over typed slots (no `any` boxing in arithmetic), whole-array
   expressions and `matmul` delegated to go-ndarray.
2. **B-upper**: the same kernel compiled to wasm32 by an oracle and run in wazero (desktop) and in
   the browser engine (playground). The oracle is clang 23.1.4 `-O2 --target=wasm32-wasip1` on a
   line-for-line C transliteration of the kernel (flang has no wasm32 target, §2); LFortran's `wasm`
   back end on the Fortran source where it accepts it, as a cross-check. **B-near**: the same at
   clang `-O0`, because go-embedded-c's first performance target is "parity with clang `-O0`"
   (`design-m2.md:99-100`): what B delivers at its first release.
3. **Reference**: gfortran 15.2.0 `-O2` natively.

For K3-`matmul` and K4 under B, the time includes copying operands into linear memory and the result
out, since a cell's arrays cross to Go (§7); the copy cost is also reported alone.

**Protocol.** One core (`GOMAXPROCS=1`; `taskset` where available), interleaved rounds (A, B-near,
B-upper, gfortran in each round, 7 rounds, median and spread), never at 1-minute load above 8;
desktop on darwin/arm64 and linux/amd64; playground in headless Chromium, timed inside the page.
Each result's checksum within a stated error bound against a `math/big` evaluation of the kernel
(the fleet's correctness contract); a kernel out of bound is a failure, not a time.

**Decision rule.** Let *rA*, *rB0*, *rB2* be the times of A, B-near and B-upper divided by
gfortran's, *g* the geometric mean over K1-K3 (loop kernels).

- **B** if *g(rA) / g(rB0) ≥ 5* on the desktop **or** in the playground: A would be five times slower
  than B's first release on loop code, a permanent open defect under the fleet's performance rule.
- **A** if *g(rA) / g(rB0) ≤ 2* in **both** homes **and** A's K4 and K3-`matmul` are within 1.2× of
  calling go-ndarray directly from Go (the evaluator adds nothing to array code).
- **Between**: **A first**, because it ships without go-embedded-c M3 and shares the array model with
  Julia; the same gate is re-run when go-embedded-c M3 is released, and B replaces A's scalar
  execution if the rule then says B. The front end's typed tree is the contract both consume, so the
  switch does not touch the front end.
- In every case the gap to gfortran (*rA* or *rB2*) is published and is an open defect.

## 5. Runtime library: I/O and intrinsics

| piece | size of the references | options | proposal |
|---|---|---|---|
| formatted, list-directed, namelist I/O (F2018 clauses 12-13) | flang-rt I/O core 9 219 lines (Apache-2.0 WITH LLVM-exception); libgfortran `io` 23 908 (GPL, oracle); LFortran `lfortran_intrinsics.c` 21 050 (BSD-3, I/O and more) | (i) a Go package, called directly by A and as host functions by B; (ii) compile `lfortran_intrinsics.c` to wasm with go-embedded-c — serves B only, and ties formatting to a C file | **(i)**, ported from flang-rt's `edit-output`/`edit-input`/`format-implementation`; units are in-memory (`print`/unit 6 to the cell output, files on a read-only `fs.FS`) |
| elemental math intrinsics | flang-rt `numeric.cpp`; gfortran calls the system libm | Go `math` / go-simd kernels | Go, each with a stated ulp bound against `math/big`; differences with gfortran's libm are judged by the bound, not by equality |
| array intrinsics (`matmul`, `dot_product`, `transpose`, `sum`, `product`, `maxval`, `minloc`, `findloc`, `reshape`, `spread`, …) | flang-rt one file per family (§3.2) | go-ndarray kernels | go-ndarray; a family it lacks is added **to go-ndarray**, not to the Fortran runtime |
| `random_number`, `random_seed` | flang-rt `random.cpp` | go-ndarray/random | go-ndarray/random; gfortran's stream is not owed (to confirm in Annex A) |
| character intrinsics, `ieee_arithmetic`, `iso_fortran_env`, `iso_c_binding` | LFortran `src/runtime/pure` (4 `.f90` modules) | Go | Go; the intrinsic modules are generated from a table |

## 6. Sharing with the fleet

| what | owner | Fortran uses it for | note |
|---|---|---|---|
| the BLAS/LAPACK routines as Go functions over column-major buffers | **go-ndarray/linalg** (exported layer to add; today `linalg/blas.go` is unexported) | `call dgemm(…)`, `call dgesv(…)` from a cell resolve to Go, never to compiled reference LAPACK | the same mapping `bib/julia.md` l. 223-228 prescribes for Julia's `@blasfunc` sites: one implementation for both languages; families go-ndarray lacks are added there |
| `Wrap(data, shape, strides, offset)` | go-ndarray #76 | arrays crossing to Go without copy, in both options | already needed by Julia (`bib/julia-runtime.md` §6.2) |
| `matmul`, `dot_product`, `transpose`, reductions | go-ndarray | intrinsics | `MatMul` already reads strided views in place (`linalg.go:24-29`) |
| FFT | go-fft (through go-ndarray/fft) | FFTPACK-style calls, if a cell names them | `go-ndarray/fft` already wraps go-fft (`fft/fft.go:1-17`) |
| front end | go-embedded-fortran | also R's `.Fortran` sites | base R calls `.Fortran` at 7 sites and stats at 16 (14 names), backed by 18 206 lines of Fortran in stats (`bib/r.md` l. 68-69): a Fortran front end is a route for go-embedded-r to run them |
| CFG IR, Ramsey structurer, object encoder, cell linker, WASI host | go-embedded-c (to extract into an exported package) | option B only | §3.5: today keyed by cc/v4 and under `internal/` |
| syntax highlighting | go-rouge `fortran` lexer | go-notebook `Highlighter` | free form only; fixed form is a go-rouge item |

## 7. Notebook semantics

- **Session module.** The session owns one compiler scope chain, as LFortran's chained
  `TranslationUnit`s (§3.1). A cell is split into program units (`module`, `subroutine`,
  `function`, `program`, `block data`) and top-level material; top-level declarations go into an
  implicit session module `__cell<N>`; top-level executable statements are wrapped in a subroutine
  `__cell<N>_run` called after the cell is compiled; every earlier session module and every module
  a cell defined is in scope (an implicit `use`) for later cells.
- **Names are case-insensitive**: the session normalises to lower case; the name that
  crosses to Go is the lower-case one, and `NameChecker` refuses a Fortran name that is a Go keyword
  only on the Go side.
- **Redefinition** shadows, LFortran's rule: earlier code keeps its binding, later cells see the new
  one (`interactive_cells.md:24-53`). A changed module re-run in place is a new module version.
- **Implicit typing**: off by default in cells (an `implicit none` session default), on by a session
  option, because a typo silently creating a `real` is the classic Fortran error and a notebook
  re-runs cells out of order. LFortran's behaviour (implicit typing until a cell says `implicit
  none`, `interactive_cells.md:55-68`) is the documented alternative.
- **go-notebook**: `Sequential` (`lang.go:88-99`) carries the scope chain through `AnalyzeNext`; the
  defined names are the session-module variables and the program-unit names, the used names are the
  free references; `Crosses` lets variables cross and keeps procedures and modules Fortran's, like
  `localsLang` (`lang.go:286-298`).
- **Values to Go**: `integer(k)` → `int8…int64`; `real(4|8)` → `float32|float64`; `complex(8|16)` →
  `complex64|complex128`; `logical` → `bool`; `character(len=*)` → `string` (trailing blanks kept);
  arrays of those → `*ndarray.Array` with column-major strides and **1-based bounds recorded** as
  metadata (go-ndarray is 0-based); derived types → opaque handles. From Go into a cell: an ndarray
  becomes an assumed-shape array; a row-major one is accepted as a `Transpose` view, not copied.

## 8. Test corpora

| corpus | revision | size | licence | use |
|---|---|---|---|---|
| gcc `gfortran.dg` | gcc 15.2.0 | 8 306 files (`.f90 .f .F90 .F .f03 .f08 .f95 .f77 .F03 .F08`), 372 689 lines; 3 153 lines contain `dg-do run` | GPL — run as data | run tests: exit status and output; the rest: must compile or must be refused |
| gcc `gfortran.fortran-torture` | gcc 15.2.0 | 314 files, 9 551 lines | GPL | compile and execute |
| llvm-test-suite `Fortran/gfortran` | `19b349c` | 8 492 files, 371 744 lines | imported from GCC (README) | its `regression` and `torture` dirs, with flang's pass/fail lists |
| llvm-test-suite `Fortran/UnitTests` (incl. `fcvs21_f95`), `SNAP`, `cray` | `19b349c` | 202 files, 107 766 lines; 27 / 9 624; 25 / 2 222 | Apache-2.0 WITH LLVM-exception; FCVS: NIST, "allows redistribution but requires explicit acknowledgement of what has been changed" (`CMakeLists.txt:32-34`) | FCVS = the F77 conformance suite (180 files with `.f .f90 .FOR .for`, 102 840 lines), references from gfortran 7.5 |
| LFortran `integration_tests` | `4f2cd049c87f` | 7 153 files, 303 283 lines; 6 651 registered `RUN` lines | BSD-3 | programs that check themselves (`error stop`); labels say which ones LFortran runs |
| Reference LAPACK `TESTING` + `BLAS/TESTING` | `0c51225` (release v3.12.1) | 1 250 + 12 files, 579 218 + 36 218 lines | modified BSD (`LICENSE:15-27`) | the end-to-end exit test: reference BLAS (`BLAS/SRC` 170 files, 48 235 lines) and LAPACK (`SRC` 2 212 files, 935 245 lines) compiled by our front end, its tests passing — the dogfooding gate, as Base is for Julia |

"Correct" = same exit status and same stdout bytes as gfortran 15.2.0, except where the output is
processor-dependent (list-directed real formatting, `random_number`; Annex A lists these, not read)
— there,
values parsed back and compared within the stated bound. Report per corpus, by file: pass /
fail-only-ours / fail-both (`bib/julia.md` §8).

## 9. Validation and performance

- A Go harness compiles each corpus file with gfortran (and flang on Linux), runs it, and compares
  with our kernel; `go test -fuzz` on the lexer (both source forms), the parser, the cell splitter,
  the format-string parser and (B) the module loader, corpora committed.
- Performance per the fleet rule: §4.2's kernels stay as the standing benchmark after the gate;
  LAPACK routines through go-ndarray are measured against reference LAPACK compiled by gfortran
  `-O2` and against OpenBLAS where numpy uses it. Nothing was timed for this file.

## 10. Proposed home and plan

Effort scale of `bib/`: S days, M weeks, L months, XL several months.

| work item | org/repo | depends on | size | priority |
|---|---|---|---|---|
| M0 oracle harness: pinned gfortran 15.2.0 (and flang 23.1.4 on Linux), corpus runner (gfortran.dg run tests, fortran-torture, FCVS, LFortran integration tests), golden outputs; read soypat/go-fortran in depth | go-embedded-fortran/fortran (`internal/oracle`, `cmd/foracle`) | — | S | 1 |
| M1 front end, part 1: lexer (free and fixed form, continuation, `include`), parser to an AST for the gate's kernels; **the §4.2 gate**, report published | go-embedded-fortran/fortran | M0 | M | 1 |
| M2 front end, part 2: semantics to a typed tree (scopes, modules, generic resolution, implicit typing, kinds, array shapes, constant folding), the v1 subset of §1.1; cell splitter and session module (§7); go-notebook `Sequential` analysis | go-embedded-fortran/fortran + go-notebook | M1 | L | 1 |
| `Wrap` (issue #76); an exported BLAS/LAPACK-named layer over column-major buffers (shared with Julia) | go-ndarray/ndarray | — | M | 1 |
| M3 back end as the gate decided: A = closure compiler over the typed tree; B = Fortran lowering to an IR extracted from go-embedded-c (exported package), on go-embedded-c M3's session | go-embedded-fortran/fortran (A); + go-embedded-c (B) | M2; B also go-embedded-c M3 | L | 1 |
| M4 runtime: I/O (§5, ported from flang-rt), intrinsics on go-ndarray/go-fft, intrinsic modules | go-embedded-fortran/fortran (`rt`) | M2 | M | 1 |
| M5 notebook bridge: values table §7, `file:line` errors, ` ```fortran ` cells on desktop and playground | go-embedded-fortran/fortran + go-notebook | M3, M4 | M | 1 |
| M6 conformance: FCVS, gfortran.dg run tests, reference LAPACK and its `TESTING` built by our front end; performance table | go-embedded-fortran/fortran | M5 | L | 2 |
| fixed-form lexer for go-rouge | go-rouge/rouge | — | S | 3 |
| site, landing, brand | go-embedded-fortran/{docs, go-embedded-fortran.github.io, .github, brand} via `ghnew` | M5 | S | 2 |

Repository name follows `bib/julia.md` §9 and `bib/c-cpp.md` §5: `go-embedded-fortran/fortran`.
The front end (M1 + M2) dominates, as Sema does for C++: LFortran spends 56 018 lines on semantics
alone, flang 94 303 on Semantics + Evaluate.

## 11. Risks and open questions

- **Security of user Fortran.** Under A, a cell is data for a Go evaluator: array accesses are Go
  slice accesses and can never leave the array (Fortran's `-fcheck=bounds` for free, against
  silent corruption in native Fortran); loops and recursion are bounded by a context deadline
  checked on back edges and calls, allocations by a session memory budget checked in `allocate`
  and in automatic arrays (a cell's `allocate(a(huge(0_8)))` must fail, not grow the host). Under
  B, the limits of `bib/c-cpp.md` §6 apply (wazero memory pages, `WithCloseOnContextDone`). In both,
  the kernel parses untrusted text: lexer, parser, `FORMAT` strings and namelist input are fuzzed,
  with bounds asserting a refusal under `GOMEMLIMIT`: source bytes, statement continuation lines
  nesting depth, `include` depth, array rank,
  repeat counts in formats, `character` lengths. (The standard's own continuation limit, 6.3.2, was
  not read.)
- **Front-end size.** 114 477 lines in flang, 69 239 in LFortran's parser + semantics: the plan's
  dominant risk. v1's subset (§1.1) and the order "what LAPACK and FCVS need first" bound it.
- **Fixed form** has no reserved words and ignores blanks (`DO 10 I = 1.10` is an assignment):
  LFortran needed a separate 2 436-line tokenizer for it.
- **B depends on another plan's unfinished milestone** (go-embedded-c M3, "planned") and on an
  extraction from its code that the C plan did not schedule.
- **A's speed is unknown in the fleet** (julia's interpreter speed is unmeasured, §3.5); the gate
  exists for that.
- **LFortran's kernel module** uses legacy exception handling and is refused by wazero (§3.1): even
  as a fallback C++-style host (go-embedded-cpp's path) it would need `--translate-to-exnref`.
- **Unverified in this file**: every speed statement (nothing was timed); the contents of the
  23.1.4 Linux tarballs (flang's presence is read from `Release.cmake`, not from a listing); how
  webR builds R's Fortran for wasm with flang (`bib/r.md` l. 251) given that `Target.cpp` has no
  wasm32 case at 23.1.4 (patched flang, or a path that does not reach `CodeGenSpecifics`?); whether
  DO CONCURRENT may always be run sequentially (11.1.7 not read beyond its title); the F2018 clause
  numbers are from the table of contents only, the clause bodies were not read (so the
  case-insensitivity of names, processor dependencies and source-form limits are stated from
  common knowledge of the language, not from the text); soypat/go-fortran's
  coverage beyond its README and declarations (its tests not run; "no fixed form" is the absence of
  the word in `lexer.go`); LFortran's wasm coverage is a count of registry labels, not a run; the
  LFortran test-registry and `dg-do run` figures count lines, not tests; FCVS's file count includes
  every file with the four extensions in its directory; licence of each `gfortran.dg` file not read
  individually; the NIST FCVS licence text itself was not opened (only llvm-test-suite's summary
  of it); go-rouge's fixed-form behaviour not tested; no Fortran kernel, gfortran or flang was
  built or run.

Verified on 2026-11-13: all file and line counts by `tools/srccount` over the revisions of §0; the
wasm refusal by `tools/wasmimports`; artefact sizes by HTTP from lfortran.github.io; release assets,
org state, repository metadata and issue #76 by `gh api`; licence files read in each clone.
