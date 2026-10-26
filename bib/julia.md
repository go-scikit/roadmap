# Julia — bibliography and feasibility plan for a pure-Go, embeddable Julia

Written 2026-10-26. Facts only, each with a source; counts come from the sources, not from estimates.
Goal: a Julia interpreter written in Go (working name **go-embedded-julia**, "jlgo") that backs
` ```julia ` cells in go-notebook and its WebAssembly playground, the way
[rbgo](https://github.com/go-embedded-ruby/ruby) backs ` ```ruby ` cells. No code exists; the
organisation does not exist (`gh api users/go-embedded-julia` returned 404 on 2026-10-26).

## 0. Sources read, and how the counts were made

| source | revision read | what was used |
|---|---|---|
| [JuliaLang/julia](https://github.com/JuliaLang/julia) | `5cd977df8812`, 2026-10-26, `VERSION` = `1.15.0-DEV` | `src/`, `base/`, `test/`, `doc/src/{manual,devdocs}`, `THIRDPARTY.md`, `Make.inc`, `stdlib/*.version` |
| `JuliaSyntax/` (in the julia tree) | same commit | parser source; the standalone [JuliaSyntax.jl](https://github.com/JuliaLang/JuliaSyntax.jl) README says that since 2.0 "active development has moved to `JuliaSyntax/` in the `JuliaLang/julia` repository (for Julia 1.14+)" |
| `JuliaLowering/` (in the julia tree) | same commit | lowering source; the standalone [c42f/JuliaLowering.jl](https://github.com/c42f/JuliaLowering.jl) README read for design and status |
| [JuliaDebug/JuliaInterpreter.jl](https://github.com/JuliaDebug/JuliaInterpreter.jl) | `0964ddb7c00b`, 2026-10-24 | `src/`, `docs/src/index.md` |
| [JuliaLang/LinearAlgebra.jl](https://github.com/JuliaLang/LinearAlgebra.jl) | `80cce3f44348`, 2026-10-25 (HEAD; julia pins `66e8f07c734b` in `stdlib/LinearAlgebra.version`, not the one read) | `src/` ccall survey |
| [JuliaStats/Statistics.jl](https://github.com/JuliaStats/Statistics.jl) | `59c5fece0f64` = the SHA pinned in `stdlib/Statistics.version` | ccall survey, size |
| Zappa Nardelli, Belyakova, Pelenitsyn, Chung, Bezanson, Vitek, *Julia subtyping: a rational reconstruction*, PACMPL 2 (OOPSLA), 2018, [doi:10.1145/3276483](https://doi.org/10.1145/3276483) | Crossref record | the formal model of `subtype.c` (title, authors, venue checked via Crossref; paper body not read) |
| go-notebook [`docs/content/en/ruby.md`](https://github.com/go-notebook/docs/blob/main/content/en/ruby.md), [go-scikit/ruby](https://github.com/go-scikit/ruby) `doc.go` and file list, [go-notebook/playground](https://github.com/go-notebook/playground) README, [go-ruby-regexp/regexp](https://github.com/go-ruby-regexp/regexp) README, `gh repo list go-bignum` | `main`, 2026-10-26 | fleet side |

Counting method. The Julia clones were read only (no build script run). Every count below was
produced by a small Go program (scratch tool `jlcount`, with a test that plants each `ccall`
form and each non-form: comment, docstring, string literal, `@ccallable`, export list), not by
`grep`. It strips comments and triple-quoted strings, then parses the first argument of each
`ccall(` and the callee of each `@ccall`. Names built at run time with `$(...)` inside `@eval`
loops cannot be resolved statically; they are counted as "generated" and resolved by reading the
loops (cited by line).

## 1. What Julia is made of (sizes)

| part | path | files | lines | role for a Go port |
|---|---|---|---|---|
| C/C++ runtime + codegen | `src/*.{c,cpp,h}` | 229 | 169 453 | spec; the LLVM half is dropped |
| semantic core of it | `subtype.c` 7 737, `gf.c` 6 567, `jltypes.c` 4 728, `datatype.c` 3 200, `builtins.c` 3 161, `module.c` 2 504, `runtime_intrinsics.c` 1 862, `typemap.c` 1 795, `method.c` 1 538, `interpreter.c` 1 118, `toplevel.c` 1 020, `genericmemory.c` 631, `array.c` 324, `symbol.c` 169, `simplevector.c` 123 | 15 | 36 477 | what must be re-expressed in Go |
| dropped with LLVM | `codegen.cpp` 11 790, `cgutils.cpp` 5 351, `jitlayers.cpp` 2 874, `aotcompile.cpp` 2 824, `ccall.cpp` 2 459, `intrinsics.cpp` 1 930, `llvm-*.cpp` … | | | not needed by an interpreter |
| Base, in Julia | `base/**/*.jl` | 202 | 140 705 | run as is through the interpreter |
| type inference + optimiser | `Compiler/src/*.jl` | 41 | 36 867 | not needed for speed; see 2.5 |
| parser, in Julia | `JuliaSyntax/src` | 18 | 11 326 (`parser.jl` 3 856, `tokenize.jl` 1 317) | port to Go |
| lowering, in Julia | `JuliaLowering/src` | 17 | 13 940 (`desugaring.jl` 4 693, `linear_ir.jl` 1 406, `scope_analysis.jl` 961, `closure_conversion.jl` 733) | port to Go |
| legacy front end, femtolisp | `src/*.scm` | 11 | 10 622 (`julia-syntax.scm` 5 779, `julia-parser.scm` 2 778, `macroexpand.scm` 733, `ast.scm` 568) | oracle and fallback |
| femtolisp itself | `src/flisp/flisp.c` | | 2 490 | not ported |
| JuliaInterpreter | `src/*.jl` | 12 | 6 943 (`construct.jl` 1 114, `interpret.jl` 1 103, `builtins.jl` 708) | design reference |
| test suite | `test/**/*.jl` | 297 | 130 413 (`core.jl` 9 501, `syntax.jl` 4 797, `subtype.jl` 3 666, `arrayops.jl` 3 554) | oracle, section 6 |
| LinearAlgebra | `src/*.jl` | 35 | 37 491 | run as Julia over go-ndarray kernels |
| Random + Statistics + Printf + Dates | `src/*.jl` | 24 | 11 219 | run as Julia |

The lowered IR (`doc/src/devdocs/ast.md:244-603`) has **12 node types** (`Expr`, `SlotNumber`,
`Argument`, `CodeInfo`, `GotoNode`, `GotoIfNot`, `ReturnNode`, `QuoteNode`, `GlobalRef`,
`Core.BindingPartition`, `SSAValue`, `NewvarNode`; lines 253-318) and **30 `Expr` heads**
(`call`, `invoke`, `invoke_modify`, `static_parameter`, `=`, `method`, `struct_type`,
`abstract_type`, `primitive_type`, `global`, `const`, `new`, `splatnew`, `isdefined`,
`the_exception`, `enter`, `leave`, `pop_exception`, `inbounds`, `boundscheck`, `loopinfo`,
`gc_preserve_begin`, `gc_preserve_end`, `aliasscope`, `popaliasscope`, `copyast`, `meta`,
`foreigncall`, `new_opaque_closure`; lines 327-567). That is the whole instruction set an
interpreter must execute; Julia's own C interpreter over it, `src/interpreter.c`, is 1 118 lines.

## 2. Architecture

### 2.1 Pipeline

```
source ──JuliaSyntax (Go port)──▶ Expr ──macroexpand──▶ Expr ──JuliaLowering (Go port)──▶ CodeInfo
   ──▶ toplevel eval (toplevel.c semantics) ──▶ method tables ──call──▶ dispatch (gf.c, typemap.c,
       subtype.c) ──▶ interpreter over lowered IR (interpreter.c semantics) ──▶ builtins (83),
       intrinsics (94), foreigncall shim table (section 3)
```

This is exactly the pipeline `doc/src/devdocs/eval.md:32-71` describes, minus step 12 (codegen):
`jl_toplevel_eval_flex` already "uses some simple heuristics to decide whether to JIT compile the
AST or to interpret it directly" and "the bulk of the work to interpret code is handled by `eval`
in `interpreter.c`" (eval.md:59-61). The port makes the interpreted branch the only branch.

### 2.2 Parser: port JuliaSyntax

- Default parser since 1.10 (`HISTORY.md:1228`: "JuliaSyntax.jl is now used as the default
  parser"); `doc/src/devdocs/eval.md:75` says the same. Base itself loads it into the system image
  (`base/Base.jl:325-327`, `set_syntax_version(Base, v"1.14")`).
- **Bootstrapping forces a native parser.** Until `Base.jl:325` runs, Julia parses `base/*.jl`
  with the femtolisp parser in C. A Go VM has no femtolisp, so it needs a Go parser before any
  Julia runs: running JuliaSyntax through the interpreter is not an option for the first files.
- JuliaSyntax is a hand-written recursive-descent parser over a token stream with a separate
  green tree (`core/parse_stream.jl` 1 003 lines) and `integration/expr.jl` (717) to build `Expr`.
  It ports line by line to Go (rbgo's parser was built the same way). 11 326 lines of Julia is
  about the size of rbgo's parser work.
- Oracle: the real `JuliaSyntax.parseall(Expr, …)` on every file of `base/`, `stdlib/`, `test/`
  and the `JuliaSyntax/test` corpus, compared as serialised `Expr` trees.

### 2.3 Lowering: port JuliaLowering, keep julia-syntax.scm as the oracle

| | `src/julia-syntax.scm` (+ `macroexpand.scm`, `ast.scm`) | `JuliaLowering/src` |
|---|---|---|
| language | femtolisp Scheme, pattern-matching macros | Julia, plain functions over a `SyntaxGraph` |
| size | 5 779 + 733 + 568 | 13 940 |
| status | **default** | in the system image but off: `base/client.jl:644-650` activates it only when `JULIA_USE_FLISP_LOWERING` is false; `base/Base.jl:329-330` loads it with `activate!(false)` |
| provenance | line numbers only | byte-precise (README "Goals") |
| hygiene | `macroexpand.scm` | "automatic hygiene" prototype (README) |
| test | — | `JuliaLowering`, `JuliaLowering_stdlibs` are entries of `test/choosetests.jl:35` |

Recommendation: **port JuliaLowering**, because it is Julia with explicit data structures
(closer to Go than Scheme pattern rewriting), it is MIT, it carries provenance the notebook needs
for `file:line` errors, and upstream is moving to it (it is now in the tree and in the test list).
Risk: its README still says "many types of syntax are not yet handled" (standalone repo); before
the port starts, measure `JuliaLowering_stdlibs` on the reference julia to know what fraction of
Base it lowers identically. Where it is short, port the matching `julia-syntax.scm` rule. The
oracle in all cases is `Meta.lower` on the reference binary (flisp, the default), compared on
`CodeInfo` after renumbering slots and SSA values.

### 2.4 Runtime in Go

| C file (lines) | Go component | notes |
|---|---|---|
| `jltypes.c` 4 728, `datatype.c` 3 200 | type objects: `DataType`, `UnionAll`, `Union`, `TypeVar`, `Vararg`, layout, type cache | field layout must match for `sizeof`, `fieldoffset`, `reinterpret`; isbits values stored unboxed in Go byte slices |
| `subtype.c` 7 737 | `issubtype`, `type_intersection`, `morespecific` | the hardest part; formal model in Zappa Nardelli et al. 2018; specificity rules in `devdocs/types.md:513-530`; oracle `test/subtype.jl` (3 666 lines) and `test/specificity.jl` |
| `gf.c` 6 567, `typemap.c` 1 795, `method.c` 1 538 | method tables, dispatch cache, world ages, ambiguity | `devdocs/functions.md:13-22`: call = form `Tuple{typeof(f), typeof(x)…}`, look it up in the method table, cache |
| `builtins.c` 3 161 | the **83 builtins** of `src/builtin_proto.h:12-94` (`getfield`, `setfield!`, `isa`, `===`, `<:`, `apply_type`, `_apply_iterate`, `invoke`, `memorynew`, `memoryref*` ×12, `*global*` ×14, `_structtype`, `_typebody!`, …) | each one a Go function |
| `runtime_intrinsics.c` 1 862 | the **94 intrinsics** of `src/intrinsics.h:3-111` | section 3.1 |
| `interpreter.c` 1 118, `toplevel.c` 1 020 | interpreter over `CodeInfo` and toplevel eval | |
| `module.c` 2 504 | modules, bindings, binding partitions | |
| `genericmemory.c` 631, `array.c` 324 | `Memory{T}` / `MemoryRef`; `Array` is defined in Julia on top | |
| `task.c`, `partr.c`, `gc-*.c`, `signals-*.c` | goroutine per Task with a hand-off channel (one runs at a time, as in Julia's single-threaded scheduler); Go GC; no signals | `Threads.@spawn` runs on one OS thread in v1 |
| `staticdata.c` 5 103 | own image format (section 5) | not a port of `sys.ji` |

The interpreter shape follows JuliaInterpreter: a `Frame` with `locals` (slots) and `ssavalues`,
a program counter over `CodeInfo.code`, `GotoNode`/`GotoIfNot`/`ReturnNode` as control flow,
`:enter`/`:leave` as an exception-handler stack. JuliaInterpreter is the evidence that lowered IR
is directly executable for the whole of Base: `@interpret sum(list)` (its `docs/src/index.md:25`)
runs Base's own `sum` statement by statement. Its limits are also the port's: it hands
`foreigncall` and `llvmcall` to compiled wrappers (`src/interpret.jl:203-225`,
`src/optimize.jl:1`, `src/packagedef.jl:18`), which is exactly where the Go port needs its shim
table instead; and a short list of Base methods are always run compiled
(`src/packagedef.jl:58-69`: `vcat(::Vector)`, `_getindex_ra`, `decompose(::BigFloat)`,
`dsfmt_jump`), a hint of which paths are too slow when interpreted.

### 2.5 Type inference: only where Base's results depend on it

Base calls inference for **results**, not just speed: `Base._return_type = Compiler.return_type`
(`base/Base_compiler.jl:418`), used by `promote_op` (`base/promotion.jl:644`) and through it by
`map`/`broadcast`/`collect` on empty inputs and by `accumulate` (the files that name
`_return_type`/`promote_op`/`infer_return_type`: `accumulate.jl` 8 lines, `promotion.jl` 15,
`reflection.jl` 10, `lock.jl` 6, `broadcast.jl` 2, `iterators.jl` 2, `reducedim.jl` 2, five
others 1 each). `map(x -> x + 1, Int[])` returning `Vector{Int}` and not `Vector{Any}` is an
observable result. Options:

1. interpret `Compiler/src` (36 867 lines) itself, lazily, only when `_return_type` is called;
   correct by construction, slow (a cold inference of a broadcast kernel through an interpreter);
2. a small Go abstract interpreter over the same lowered IR (types only, no constant propagation,
   widen to `Any` on recursion), fast, and wrong in the corners.

Recommendation: (1) first, because it is the reference, with a per-signature cache that is part of
the image (section 5); (2) only if measurement shows the notebook waiting on it, and then checked
against (1) on every call in the test suite.

### 2.6 Base bootstrap: interpret base/*.jl, port nothing of Base natively

Base is bootstrapped by evaluating `base/boot.jl` (Core, 1 338 lines) and then the file list of
`base/Base_compiler.jl` (`include` lines 139-423: `essentials.jl`, `int.jl`, `float.jl`,
`array.jl`, …, then `Compiler/src/Compiler.jl` at 415) and `base/Base.jl`, as
`devdocs/eval.md:182-187` and `devdocs/sysimg.md:3-9` describe. Recommendation: **run Base
unchanged through the interpreter, at build time, and ship the resulting heap as an image**. Port
to Go only (a) builtins, intrinsics and `foreigncall` targets (sections 2.4 and 3), and (b) hot
leaf methods after profiling, replaced by method overrides registered in Go (same signature, same
results, a test comparing both). Reasons: Base is 140 705 lines whose semantics are the
definition of Julia; a native rewrite would diverge (the lesson of `ruby-reuse.md` row 14: "a
better answer is a different answer"), and every Base fix upstream would have to be ported twice.

## 3. The runtime boundary

### 3.1 Intrinsics: 94 (`src/intrinsics.h:3-111`)

81 `ADD_I` + 12 `ALIAS` + 1 `ADD_HIDDEN`. By group: wrap 1 (`bitcast`); arithmetic 17; fast
arithmetic 7 (aliases); comparisons 15 (incl. 4 aliases, `fpiseq`); bitwise 11 (incl.
`bswap_int`, `ctpop_int`, `ctlz_int`, `cttz_int`); conversion 9; checked arithmetic 10; functions
9 (`abs_float`, `copysign_float`, `flipsign_int`, `ceil/floor/trunc/rint/sqrt_llvm`,
`sqrt_llvm_fast`); pointer arithmetic 2; pointer access 2; pointer atomics 6; C interface 2
(`cglobal`, `llvmcall`); CPU tests 2 (`have_fma`, `cpu_supports`); hidden 1 (`cglobal_auto`).

All except `llvmcall` and the pointer group map to Go in a few lines each: `math/bits`
(`OnesCount`, `LeadingZeros`, `TrailingZeros`, `ReverseBytes`, `Add64`/`Mul64` for checked
arithmetic), `math` (`FMA`, `Sqrt`, `Floor`, `Copysign`, `RoundToEven`), explicit widths for
`trunc_int`/`sext_int`/`zext_int` on 8 to 128 bits (`Int128` needs a two-word type), Float16 and
BFloat16 by software conversion. `muladd_float` may be computed as `FMA` or not: choose FMA and
document it. `pointerref`/`pointerset` and the atomics operate on the Go-side byte storage of
`Memory`/`Ref`; there is no raw address space (see risks). `llvmcall` is not supported: Base has
**0** `llvmcall(` sites; `base/fastmath.jl:304-306` uses `ccall("llvm.powi", llvmcall, …)` three
times (map to `math.Pow`); `stdlib/Random/src/XoshiroSimd.jl` has 10 `llvmcall` and 2 generated
`ccall($fshl)` — its SIMD bulk fill must be replaced by the scalar Xoshiro path.

### 3.2 `ccall` / `@ccall` in base/: 896 call sites in 202 files

856 `ccall` + 40 `@ccall`; 24 `cglobal(`. By target:

| target | sites | unique names | files (sites) | Go implementation |
|---|---|---|---|---|
| Julia runtime `jl_*`, `ios_*` | 524 | 387 | 72 files; most in `runtime_internals.jl` 43, `iostream.jl` 38, `loading.jl` 34, `boot.jl` 25 | **is** the Go runtime: reflection (`jl_method_*`, `jl_type_*`, `jl_module_*` 10, `jl_get_*` 29, `jl_set_*` 18), IdDict (`jl_eqtable_*` 4, `jl_idset_*` 4), strings (`jl_alloc_string`, `jl_string_ptr`, `jl_pchar_to_string`), files (`jl_fs_*` 12, `jl_stat_*` 12, `ios_*` 22, `jl_uv_*` 19), GC (`jl_gc_*` 14). No-ops or errors: compiler/precompile/cpu (`jl_compile_*`, `jl_precompile_*` 5, `jl_cpufeatures_*` 8, `jl_cpu_*` 6, `jl_cumulative_*` 4, coverage, `jl_restore_*`) |
| libc / OS by process symbol | 88 | 70 | `libc.jl` 27, `env.jl` 12, … | Go `os`, `syscall`, `time` (`getenv`/`setenv`, `strerror`, `localtime_r`, `mktime`, `strptime`, `memcpy`/`memcmp`/`memchr`/`memset` on Go slices, `malloc`/`free` on a Go arena); part of the names are Windows or macOS only (`GetEnvironmentStringsW`, `_NSGetEnviron`, `mach_timebase_info`, …) |
| **MPFR** | 81 static + 29 generated | 72 static | `base/mpfr.jl` (1 334 lines); generated by `$(string(:mpfr_, fC, …))` at `mpfr.jl:324, 579-665, 753-769, 821, 879, 949, 969, 976, 1123` | **go-bignum/bigfloat** ("correctly rounded elementary functions and constants on `math/big.Float`, any precision, all rounding modes") + `math/big.Float` for arithmetic; `mpfr_asprintf`/`mpfr_snprintf` (Printf) need a Go formatter |
| libuv `uv_*` | 78 | 67 | `stream.jl` 24, `file.jl` 10, `libuv.jl` 10, `sysinfo.jl` 10 | Go `os`, `os/exec`, `net`, `time.Timer`; the event loop becomes goroutines that wake the scheduler (section 2.4) |
| **GMP** | 40 with `libgmp` named (37 static names + `libgmp.$fname`) + 3 interpolated `$(gmpz(op))`/`$(gmpq(op))`; 32 names generated by the loops | 37 static | `base/gmp.jl` (1 127 lines); generated by `Symbol(:__gmpz_, op)` for 14 binary ops (`gmp.jl:171-175`: add, sub, mul, fdiv_q, tdiv_q, cdiv_q, fdiv_r, tdiv_r, cdiv_r, gcd, lcm, and, ior, xor), 7 `_ui` ops (191), neg/com/sqrt/set (215), scan0/scan1 (205), set_ui/si/d (230), mpq set_ui/si (1034) | `math/big.Int`/`big.Rat` (all of these exist there: `Quo`/`Rem` vs `Div`/`Mod` must follow GMP's t/f/c rounding), `__gmpz_export`, `__gmpz_limbs_*` emulated on `big.Int.Bits`; go-bignum for `fac_ui`, `bin_ui` |
| other dynamic (pointer / `@static` / interpolated) | 7 | | `libc.jl:46,48,127` (`dup`, `dup2`, `fdopen` chosen by OS), `c.jl:537`, `stat.jl:194`, `threadcall.jl:41`, `irrationals.jl:276` (`mpfr_const_*`) | resolved per case |
| **PCRE2** | 18 | 18 | `base/pcre.jl` (284 lines), used by `regex.jl` (929) | see 3.3 |
| utf8proc | 15 | 14 | `base/strings/unicode.jl` | Go `unicode` + `golang.org/x/text/unicode/norm` for NFC/NFD/NFKC/NFKD; `utf8proc_charwidth` and grapheme breaks need tables generated from the same Unicode version as the pinned utf8proc |
| Windows-only (`path.jl`, `util.jl` credui/ole32), macOS `libproc` | 13 | 10 | | not on wasm; `GOOS=windows` later |
| **libm** | **0** | | `base/math.jl:306` only defines `const libm = Base.libm_name` | nothing: Base's elementary functions are written in Julia (`base/special/*`) over intrinsics, so results are bit-identical by construction |

Standard libraries read (Random, Dates, Printf, Statistics): 14 sites. `Random/src/DSFMT.jl` 7
(`:libdSFMT`, the `MersenneTwister` generator; the default `Xoshiro`/`TaskLocalRNG` is pure Julia
in `Xoshiro.jl`), MPFR 3 (`mpfr_snprintf` in Printf, `mpfr_sub_d`/`mpfr_sub_ui` in Random),
`uv_clock_gettime` 1 (Dates), `jl_gc_init_finalizer_rng_state` 1. Statistics: 0.

LinearAlgebra (HEAD read): 194 sites, of which 182 call BLAS/LAPACK through
`(@blasfunc($name), libblastrampoline)` inside `@eval` loops over element types (`blas.jl` 38,
`lapack.jl` 144; 105 distinct callee expressions, nearly all one routine family each, such as
`$getrf`, `$geev`, `$gesdd`, `$syevr`, `$potrf`, `$trtrs`, `$gges3`, plus generic `$fname`
loops for BLAS levels 1-3), 9 sites over 8 `lbt_*` configuration names (`lbt.jl`), 3
`sysctlbyname`. Mapping: implement the Fortran BLAS/LAPACK ABI each family uses as Go functions
over **go-ndarray/linalg** (LU, QR, Cholesky, SVD, Eig, Eigh, lstsq exist there per
`bib/ruby-reuse.md` §1, row 8) on the column-major `Memory{Float64}` buffer the Julia `Matrix`
owns, without copying. Families go-ndarray lacks (generalized Schur `gges3`, `tgsen`,
`ggsvd`, banded `gbtrf`, tridiagonal `gttrf`/`pttrf`, Bunch-Kaufman `sytrf`/`hetrf`) are missing
kernels to add **to go-ndarray**, not to the VM.

### 3.3 Regex

`base/regex.jl` is PCRE2 semantics (`r"…"`, `match`, `eachmatch`, named groups, `\K`, Unicode
classes, JIT). The fleet engine is [go-ruby-regexp/regexp](https://github.com/go-ruby-regexp/regexp),
a pure-Go "faithful backtracking VM" reimplementation of **Onigmo** (named groups, backrefs,
`\g<…>` calls, fixed-width lookbehind, atomic groups, possessive quantifiers, step budget and
timeout against ReDoS; BSD-3-Clause). PCRE2 and Onigmo differ in syntax corners (`\K`, `(*VERB)`s,
callouts, variable-length lookbehind accepted by PCRE2 10.43+, `\h`, option letters). Plan: a
PCRE2 syntax front end that compiles to the same VM, with `test/regex.jl` as the oracle, the 18
`pcre2_*` names implemented on it. If the VM cannot take a construct, raise Julia's
`ErrorException("PCRE compilation error: …")` rather than match differently.

### 3.4 FFT

Julia's Base and stdlib have **no FFT** (no `fftw` in `base/`, `stdlib/` or `deps/`; FFTW.jl is a
registry package). A `go-julia-fftw` package presenting the FFTW.jl / AbstractFFTs.jl API over
**go-fft** is the port; it is not part of the VM.

### 3.5 Licences

- Julia is MIT (`LICENSE.md`), and `THIRDPARTY.md:1-4` states the language "consists of the
  compiler (the contents of `src/`), most of the standard library (`base/` and `stdlib/`)". So
  `src/` (C), `base/`, `stdlib/`, `JuliaSyntax/` (MIT), `JuliaLowering/` (MIT),
  JuliaInterpreter.jl (MIT), LinearAlgebra.jl (MIT), Statistics.jl (MIT) may all be translated to
  Go with the notice kept. `src/flisp` is BSD-3 (THIRDPARTY.md:12).
- Exceptions inside Base with their own licence: `base/ryu/*` (Boost) and
  `base/special/{rem_pio2,hyperbolic}.jl` (FDLIBM notice) (THIRDPARTY.md:14-17). Both run as
  Julia through the VM; the notices travel with the image.
- External libraries Base calls: **GMP** "LGPL3+ or GPL2+", **MPFR** LGPL3+, PCRE BSD-3,
  utf8proc MIT, libuv MIT, OpenBLAS/LAPACK BSD-3, dSFMT BSD-3 (THIRDPARTY.md:22-45). None is
  linked or translated: GMP and MPFR are replaced by `math/big` and go-bignum, written from the
  documented semantics, not from their source. dSFMT (BSD) may be translated.
- **GPL to avoid**: SuiteSparse's `libumfpack`, `libspqr`, `librbio` are GPL-2.0+ and
  `libcholmod` LGPL-2.1+ and GPL-2.0+ (THIRDPARTY.md:45-57); they back `SparseArrays` (an external
  stdlib, `stdlib/SparseArrays.version`), whose `\` on sparse matrices goes to UMFPACK/CHOLMOD.
  A Go port of SparseArrays must use fleet sparse solvers (`go-sparse`, planned in
  `PLAN.md`) and must not read SuiteSparse source. `Make.inc:296-297` has `USE_GPL_LIBS ?= 1`,
  i.e. the stock build includes them. FFTW (GPL) is not in Julia's tree (3.4). ITTAPI (BSD-3 and
  GPL2, optional profiling) is irrelevant.

## 4. A minimal first version for the notebook

### 4.1 Scope of v1

- **Language**: all syntax JuliaSyntax parses; functions, multiple dispatch with parametric
  methods and `where`, keyword arguments, closures, `struct`/`mutable struct`, abstract and
  parametric types, `Union`, `Nothing`/`Missing`, tuples, named tuples, `try`/`catch`/`finally`,
  macros (hygienic) and `@generated` functions (generators run in the interpreter), string
  interpolation, broadcasting (`.+`, `@.`), comprehensions and generators, `do` blocks, modules,
  `using`/`import` of in-image modules. Tasks and `Channel` single-threaded.
- **Out of v1**: `ccall` to arbitrary user libraries (an error naming the symbol), `llvmcall`,
  `@cfunction`, `Threads.@threads` parallelism (runs serially), `Distributed`, `Pkg` (no package
  installation from the playground; packages ship in the image), precompilation of user code.
- **Base**: everything that loads, with BigInt/BigFloat, Regex, Unicode, `Printf`, `Dates`,
  `Random` (Xoshiro default, MersenneTwister via ported dSFMT), `show`/`display` with MIME
  `text/plain` and `image/svg+xml`.
- **stdlib**: `LinearAlgebra` (dense, over go-ndarray), `Statistics` (pure Julia, 0 ccall),
  `Random`, `Printf`, `Dates`, `Test` (for the oracle). `SparseArrays` later (3.5).

### 4.2 What the playground can show, milestone by milestone

| milestone | the playground shows |
|---|---|
| M1 parser | a cell's `Expr` tree (`dump`) and parse errors with JuliaSyntax's diagnostics |
| M2 lowering | `Meta.@lower` of a cell: the `CodeInfo` listing |
| M3 Core | arithmetic, `isa`, `typeof`, `<:` queries, user-defined types and methods with dispatch, using Core only (no `println`: values shown by the host through `jl_static_show` semantics) |
| M4 Base | `println`, strings, `Vector`/`Dict`/`Set`, comprehensions, broadcasting, `sort`, `sum`, `BigInt` factorials, `r"…"` matches, `@show`, `@time` (wall clock only) |
| M5 numeric stdlib | `LinearAlgebra` (`A \ b`, `eigen`, `svd`), `Statistics` (`mean`, `std`, `quantile`, `cor`), `Random` with seeds identical to julia's output, `Printf`, `Dates` |
| M6 bridge | mixed Go/Julia notebook; go-scikit bindings (`Stats`, `Special`, `Plot` SVG) |

## 5. Bridging with go-notebook and go-scikit

### 5.1 The Ruby page as the template

The Ruby cells (`content/en/ruby.md`) define the contract to copy: one shared interpreter for all
cells of the language; only **variables** cross; a cell defines the top-level names it assigns
and uses the names it reads first (taken from the language's own parser); Go → Ruby converts
numbers, strings, booleans, slices to Arrays, maps to Hashes, `*ndarray.Array` **by identity**,
anything else as an opaque object that comes back as itself; Ruby → Go types each local by its
value after the cell runs, and a Go cell reading a not-yet-typed variable waits for it; errors are
`file:line: message (Class)` at the notebook line with at most twelve frames; worker process with
timeout and memory limits.

### 5.2 Julia ↔ Go value table (proposed)

| Julia value | Go type | Go → Julia |
|---|---|---|
| `nothing` | `any` (nil) | `nil` → `nothing` |
| `Bool` | `bool` | same |
| `Int8…Int64`, `UInt8…UInt64` | the same-width Go integer (`int64` for `Int`) | Go `int` → `Int` |
| `Int128`/`UInt128` | `jl.Int128` (two words) | |
| `BigInt` | `*big.Int` | `*big.Int` → `BigInt` |
| `Float16/32/64` | `jl.Float16`, `float32`, `float64` | |
| `BigFloat` | `*big.Float` (precision and rounding kept) | |
| `Rational{T}` | `*big.Rat` (or `jl.Rational[T]`) | |
| `Complex{Float64}` | `complex128` | |
| `String` | `string` | |
| `Symbol` | `julia.Symbol` | |
| `Char` | `rune` (invalid UTF-8 kept as `julia.Char` raw bits) | |
| `Array{Float64,N}` (any isbits eltype) | `*ndarray.Array` sharing the same buffer, **column-major** strides kept | `*ndarray.Array` → `Array` view, no copy |
| `Vector{String}`, `Vector{Any}` | `[]string`, `[]any` | slices |
| `Dict{String,V}` | `map[string]any` | maps |
| `Tuple`/`NamedTuple` | `[]any` / `map[string]any` (order kept in `julia.NamedTuple`) | |
| a go-scikit object | the wrapped Go value | itself |
| any other Julia object | `julia.Value`, callable from Go through `Call` (dispatch happens in Julia) | an opaque Go value arrives as a `GoObject{T}` whose methods are the Go methods |

Array memory order is the one real difference from Ruby: Julia is column-major, go-ndarray is
NumPy-style strided, so a Julia matrix crosses as an ndarray with Fortran strides, which
go-ndarray must accept without copying (to confirm against its stride rules).

Cell names: the Go port of JuliaSyntax gives the parse tree; the cell's defined names are the
top-level assignment targets, `function`/`struct`/`const`/`module` names; used names are free
identifiers read before assignment (JuliaLowering's `scope_analysis.jl`, 961 lines, computes
exactly this and is reused). A Julia function defined in a cell is a name but, as in Ruby, does
not cross to Go.

### 5.3 go-scikit bindings: reuse the go-scikit/ruby shape

go-scikit/ruby is one Go package whose `Register(in *rb.Interp)` defines modules `Stats`,
`Special`, `Optimize`, `Integrate`, `DataFrame`, `Plot` (its `doc.go`); "every Ruby method calls
the Go function it is named after, on the same inputs; the results are the Go library's, bit for
bit", with shared `conv.go`/`reflect.go` conversion code and generated tables (`enums_gen.go`,
`special_gen.go`). A `go-scikit/julia` should copy that: `Register(vm *julia.VM)` defining Julia
modules `Stats`, `Special`, … whose functions are Go builtins, plus — because Julia has multiple
dispatch — each Go function registered as a **method** with a Julia signature
(`Special.gamma(x::Float64)`), so that user code can add methods for its own types. The
`bib/ruby-reuse.md` §3 lesson applies from day one: the VM must expose a **public** extension
API (method registration, value conversion, object wrapping) outside `internal/`, or
go-scikit/julia cannot be written in its own module.

## 6. WebAssembly, image and start-up

- Reference points: the playground's Go core is 7 561 584 bytes on the wire (29.8 MB raw),
  first result 0.73-0.85 s cold, 0.22-0.43 s warm (go-notebook/playground README, measured
  2026-10-17). The Ruby core is reported at about 12.7 MB gzipped with a first result in 2.3 s
  (figures from the request for this file; not re-measured here). The Ruby core is fetched only
  when a notebook has a Ruby cell; a Julia core must be a third lazily loaded `.wasm`.
- **The image is the start-up.** Julia without a system image evaluates `boot.jl` and all of
  Base at start (`devdocs/eval.md:182-187`); through an interpreter that is 140 705 lines plus
  the 36 867-line Compiler if inference is interpreted (2.5) — minutes, not seconds. So the build
  runs the bootstrap once, natively, and writes the resulting heap (modules, bindings, types,
  method tables, lowered `CodeInfo`, interned symbols, the `_return_type` cache, constant data
  such as Unicode tables) to a versioned binary image, embedded with `go:embed` or fetched next to
  the wasm. This is Julia's own design (`sys.ji`, "a frozen version of the `Main`, `Core`, and
  `Base` modules", `devdocs/eval.md:175-180`, written by `staticdata.c`) without native code.
- Size levers: store `CodeInfo` in a compact encoding (Julia itself compresses IR:
  `jl_uncompress_ir` is in the ccall list); drop docstrings from the wasm image
  (`base/docs/basedocs.jl` alone is 4 327 lines) and serve them on demand; **trim** as
  `devdocs/sysimg.md:150-169` describes for `juliac`, with the same caveat that reflection
  (`getglobal(Base, Symbol(readchomp(stdin)))`) defeats it — a notebook is such a program, so
  trim only stdlibs the notebook does not import, never Base methods.
- Start-up budget to hold: image deserialisation must stay under the Ruby core's 2.3 s on the
  same machine. Unknown until M4: the image size; measure it and decide (lazy per-module images:
  `LinearAlgebra` only when `using LinearAlgebra` runs).
- The interpreter must not use goroutine-per-Task in a way that blocks the wasm main thread;
  Go's wasm port runs goroutines cooperatively on one thread, which matches Julia's
  single-threaded scheduler.

## 7. Prior art: Julia in the browser or without LLVM

| project | approach | state (as read) |
|---|---|---|
| [Keno/julia-wasm](https://github.com/Keno/julia-wasm) | the C runtime built with Emscripten, a web REPL | README: "an extremely early alpha", "many things are likely (and known) to be broken"; no size or timing figures |
| [tshort/WebAssemblyCompiler.jl](https://tshort.github.io/WebAssemblyCompiler.jl/) | compiles Julia's typed IR to wasm with Binaryen and wasm-GC | docs: "Code must be type stable (no dynamic dispatches)"; unsupported: multi-dimensional arrays, `Ptr`, Union types, exception handling, errors, some integer types, "BLAS and all other C dependencies"; MIT |
| [T. Short, "Executing Julia in the browser"](https://nextjournal.com/tshort/executing-julia-in-the-browser) (Nextjournal; search excerpt only, page not opened) | Emscripten build of the runtime | per the search excerpt: booted but failed in `jl_init` reading `boot.jl`, "the point where Julia needs code generation and LLVM"; an interpreter able to run without codegen was named as the missing piece |
| R. M. Huffman, *Julia in WebAssembly*, MIT thesis, 2023, [dspace 1721.1/150151](https://dspace.mit.edu/handle/1721.1/150151) | survey of strategies | abstract: several strategies are feasible; running the Julia compiler entirely in wasm is future work |
| [MikeInnes/Charlotte.jl](https://github.com/MikeInnes/Charlotte.jl) | Julia IR to wasm | listed in search results only; not opened |
| JuliaInterpreter.jl | interpreter over lowered IR, **inside** a running julia | works for all of Base (2.4), but relies on the native runtime for builtins, ccall and llvmcall |

No pure-reimplementation interpreter of Julia (outside the C runtime) was found. The two
browser routes both stop at the C runtime: one needs LLVM at `boot.jl`, the other refuses
dynamic dispatch. A Go VM that interprets lowered IR and ships a pre-built image is the
combination none of them tried: JuliaInterpreter's execution model plus a runtime that does not
need LLVM.

## 8. Milestones, effort, conformance

Effort in the `bib/` scale: S days, M weeks, L months, XL several months.

| # | milestone | contents | size | exit test |
|---|---|---|---|---|
| M0 | oracle harness | a Go tool that runs the reference `julia` (pinned release) on a corpus and stores its outputs: parse trees, `Meta.lower`, `Test` summaries, printed values | S | outputs checked in as golden files, regenerated by one command |
| M1 | parser | Go port of `JuliaSyntax/src` (11 326 lines) → Go `Expr`; fuzzed (`go test -fuzz`) with a depth and size bound | L | `Expr`-identical on all 202 `base/` files, `stdlib/`, all 297 `test/` files, JuliaSyntax's own tests |
| M2 | lowering + macro expansion | Go port of `JuliaLowering/src` (13 940), missing rules from `julia-syntax.scm` | L | `CodeInfo`-identical (modulo numbering) to `Meta.lower` on Base's top-level forms |
| M3 | Core runtime | types, subtyping, method tables, dispatch, 83 builtins, 94 intrinsics, interpreter, `boot.jl` | XL | `test/subtype.jl`, `test/specificity.jl`, `test/intrinsics.jl`, `test/core.jl` parts not needing Base I/O |
| M4 | Base bootstrap + image | `Base_compiler.jl` and `Base.jl` evaluated; the `jl_*`/libc/libuv shim (≈ 524 + 88 + 78 sites); Compiler interpreted for `_return_type`; image writer/reader | XL | Base loads; `test/` files `arrayops`, `strings`, `dict`, `tuple`, `reduce`, `broadcast`, `iterators`, `sorting`, `numbers`, `int`, `math`, `float16`, `ryu`, `show` |
| M5 | numeric boundary | GMP/MPFR on `math/big` + go-bignum, PCRE2 front end on go-ruby-regexp, utf8proc tables, dSFMT port, Printf, Dates, Random, Statistics, LinearAlgebra on go-ndarray | L | `test/gmp.jl`, `mpfr.jl`, `regex.jl`, `unicode.jl`, `rational.jl`, `complex.jl`; stdlib tests of Random (seeded streams byte-identical), Dates, Printf, Statistics, LinearAlgebra (dense) |
| M6 | notebook bridge + go-scikit/julia | value table 5.2, cell names, errors, worker limits | M | mixed Go/Julia example notebook transcript in go-notebook docs |
| M7 | wasm playground | lazy Julia core, image loading, SVG display through the existing sanitizer | M | first result under the Ruby core's time on the same machine; size recorded |

Total: on the order of the rbgo effort plus a parser and a lowering pass rbgo did not need in
this form; M3 and M4 dominate.

**Conformance, as the fleet does with reference tools.** Julia's suite is written with the `Test`
stdlib (`@test`, `@testset`), and `test/choosetests.jl:8-36` lists the test names
(`"core"`, `"subtype"`, `"arrayops"`, … `"JuliaSyntax"`, `"JuliaLowering"`). The harness:

1. runs each listed file under the reference `julia` and under jlgo, one file per worker process
   with a timeout, each with a `Test` reporter that emits one JSON record per `@test`
   (file, line, expression, pass/fail/error/broken);
2. compares records, not totals: a test passing in both, failing only in jlgo (defect), or
   failing in both (upstream), so that an equal count with different members is not read as
   parity;
3. keeps a committed allowlist of files out of scope (`precompile`, `loading`, `llvmcall*`,
   `threads`, `ccall` to user libraries, `spawn`, `cmdlineargs`, `compiler`, `jit`), each with a
   reason, and reports the pass rate per remaining file in STATUS.md;
4. stores every jlgo-only failure as a regression test once fixed.

The differential oracles per stage (M1 `Expr`, M2 `CodeInfo`, Random streams, printed `show`
output) are stronger than the suite for the front end and must be built first.

## 9. Naming in the fleet

Observed on 2026-10-26: `go-embedded-ruby` holds `ruby` (rbgo), `docs`,
`go-embedded-ruby.github.io`, `.github`, `brand` (`gh repo list go-embedded-ruby`); a GitHub
search for organisations whose login contains `go-ruby-` returns 206 (e.g. `go-ruby-regexp`,
`go-ruby-bigdecimal`, `go-ruby-rational`, `go-ruby-complex`, `go-ruby-csv`, `go-ruby-rspec`):
one organisation per gem or stdlib library, each holding the library repository and its own site.
The search for `go-julia` logins finds only unrelated accounts (`go-julian`, `GO-JuliaM`);
`go-embedded-julia` is free.

Recommendation: **mirror it exactly.**

- `go-embedded-julia/julia` — the VM (module `github.com/go-embedded-julia/julia`), with `docs`,
  `go-embedded-julia.github.io`, `.github`, `brand`, created with `ghnew`.
- `go-julia-<package>` organisations for package ports that are not part of the VM image's
  stdlib, named after the Julia package without `.jl`: `go-julia-fftw` (FFTW.jl/AbstractFFTs.jl
  over go-fft), `go-julia-sparsearrays` (over go-sparse, no SuiteSparse), `go-julia-dataframes`,
  `go-julia-distributions`, `go-julia-specialfunctions`. Stdlibs that ship in the image
  (`LinearAlgebra`, `Statistics`, `Random`, `Dates`, `Printf`, `Test`) stay in the VM repository,
  as rbgo carries Ruby's default libraries.
- `go-scikit/julia` for the go-scikit bindings, next to `go-scikit/ruby`.

One correction to the Ruby precedent: some `go-ruby-*` numeric libraries duplicate rbgo's own
copies (`bib/ruby-reuse.md` rows 4-5). For Julia, a `go-julia-*` package must be Julia code run
by the VM plus Go kernels **imported from the go-scikit libraries**, never a second Go
implementation of the same algorithm.

## 10. Risks and open questions

- **Subtyping and specificity** (`subtype.c`, 7 737 lines) are where a reimplementation diverges
  silently: diagonal rule, `Union` distribution, `Vararg`, bounds of `TypeVar`s
  (`devdocs/types.md:320-480`). Port `subtype.c` statement by statement rather than from the paper.
- **JuliaLowering's coverage** is not measured here; if it lowers much less than Base, M2 becomes
  a port of `julia-syntax.scm` (5 779 lines of Scheme pattern rewriting).
- **Interpretive speed.** Base's generic code relies on inlining and specialisation; interpreted,
  `sum` of a `Vector{Float64}` dispatches `+` per element. Mitigation: dispatch caches keyed by
  concrete argument types, specialised fast paths for intrinsic-only methods, and Go overrides
  for measured hot leaves (2.6) checked against the Julia method. Against the fleet rule "no slower
  than the reference": an interpreter cannot match Julia's JIT on loops; the comparison that can
  hold is against the reference for **start-up and first result** (Julia's time-to-first-plot)
  and for kernels routed to go-ndarray/go-bignum. State the measured ratio per milestone.
- **Pointers and `unsafe_*`.** Base uses `Ptr` arithmetic, `pointer(::Array)`, `unsafe_load`,
  `unsafe_wrap` and 24 `cglobal` sites. A Go VM must model a pointer as (buffer, offset) into Go
  memory, check bounds on every access (a malformed `unsafe_load` must raise, not corrupt the
  VM), and refuse `Ptr` values that point nowhere known.
- **Security** (notebook and playground run untrusted cells): fuzz the parser and the image
  reader; bound recursion depth (Julia's `StackOverflowError`), allocation sizes (`Vector{UInt8}(undef, 2^62)` must raise `OutOfMemoryError`), and regex steps; no `ccall` escape on wasm; `run`/`open`/`ENV` as Ruby cells have them on the desktop and absent on wasm.
- **Unicode version** must match utf8proc's for `textwidth`, `isletter`, normalisation, else
  `test/unicode.jl` diverges.
- **Float formatting** is in Julia (`base/ryu`), so `show(1.0)` is exact for free; BigFloat
  printing goes through `mpfr_asprintf` and needs care.
- **Upstream moves fast**: 1.15-DEV is read here; pin a **release** for the port (and its
  LinearAlgebra/Statistics SHAs from `stdlib/*.version`) and track the next one deliberately.
- **Unverified in this file**: how much of Base JuliaLowering lowers identically to femtolisp;
  every interpreter speed statement (nothing was run); the Ruby core figures (not re-measured); the body of the OOPSLA paper; T. Short's
  Nextjournal post and Charlotte.jl (search excerpts only); go-ndarray's acceptance of
  Fortran-order strides; the LinearAlgebra survey is of HEAD, not of the SHA julia pins.

Verified on 2026-10-26: all counts in sections 1-3 by the Go scratch tool over the clones at the
commits of section 0; licence headers of julia, JuliaSyntax, JuliaLowering, JuliaInterpreter,
LinearAlgebra, Statistics, femtolisp read in the clones; THIRDPARTY.md and Make.inc read in the
clone; fleet facts via `gh api` / `gh repo list` and the cited READMEs; prior-art pages opened
except where marked.
