# R — bibliography and feasibility plan for a pure-Go, embeddable R

Written 2026-11-01. Facts only, each with a source; counts come from the sources, not from estimates.

> **Lesson carried over from [`julia.md`](julia.md).** Its counts were first taken on julia
> master instead of the release the port is pinned to, and had to be corrected in
> [`julia-runtime.md`](julia-runtime.md) §1. Here every count is taken on the **R 4.6.2 release
> tarball** and nothing else. R has no release tags in git: the
> [wch/r-source](https://github.com/wch/r-source) mirror has one tag, `trunk`
> (`gh api repos/wch/r-source/tags`, 2026-11-01), so the pin is the CRAN tarball and its digest,
> not a commit.

Goal: an R interpreter written in Go (working name **go-embedded-r**), `CGO_ENABLED=0`,
embeddable and `GOOS=js`/`wasip1`-capable, that backs ` ```r ` cells in go-notebook and its
WebAssembly playground, the way [go-embedded-julia/julia](https://github.com/go-embedded-julia/julia)
backs ` ```julia ` cells and [go-embedded-ruby/ruby](https://github.com/go-embedded-ruby/ruby)
(rbgo) backs ` ```ruby ` cells. No code exists; the organisation
[go-embedded-r](https://github.com/go-embedded-r) exists and holds no repository
(`gh repo list go-embedded-r`, 2026-11-01).

## 0. Sources read, and how the counts were made

| source | revision read | what was used |
|---|---|---|
| R source, [`R-4.6.2.tar.gz`](https://cran.r-project.org/src/base/R-4/R-4.6.2.tar.gz) | SHA-256 `c297040f6392b1dac6e5c0eb398422cfdaedf540ae2205edb9ad0d10cb10ee00`; `VERSION` = `4.6.2`; `SVN-REVISION` 90784, 2026-10-29 | `src/main/{names.c,gram.y,eval.c,summary.c,RNG.c}`, `src/include/Rinternals.h`, `src/appl/dqrdc2.f`, `src/modules/lapack/{README,LICENSE,Lapack.c}`, `src/library/*/R`, `tests/`, `doc/COPYRIGHTS`, `src/library/base/R/apply.R` (licence header) |
| CRAN [`VERSION-INFO.dcf`](https://cran.r-project.org/src/base/VERSION-INFO.dcf) | 2026-11-01 | release 4.6.2 "Diversified Portfolio" 2026-10-29; old release 4.5.3; devel 4.7.0 |
| *R Language Definition*, `doc/manual/R-lang.texi` | in the 4.6.2 tarball | section list; §Evaluation of expressions (l. 860), §Indexing (1501), §Scope of variables (1917), §Functions / Evaluation (2171), §Object-oriented programming (2460: `UseMethod` 2728, `NextMethod` 2815, group methods 2891), §Computing on the language (2975), §.Internal and .Primitive (3707), §Exception handling (3718), §Parser (3977: modes of parsing 3988-4024, tokens 4074, expressions 4345) |
| *R Internals*, `doc/manual/R-ints.texi` | in the 4.6.2 tarball | §SEXPs (73), §Environments and variable lookup (499), §Contexts (738), §Argument evaluation (838-877), §Serialization Formats (1092), §S4 objects (1386), §Lazy loading (1852), §.Internal vs .Primitive (1926-1985), §Graphics Devices / Device structures (2795-2839), §Testing R code (4878-4890) |
| CRAN [`src/contrib/PACKAGES`](https://cran.r-project.org/src/contrib/PACKAGES) | fetched 2026-11-01 11:15 CET, SHA-256 `eaca4e078582ed1d688522e9e76bf2525a0f5c462bb682e4c1985964f57ffeac` | `NeedsCompilation`, `Depends`, `Imports`, `LinkingTo` of every package (§1.4) |
| Prior art READMEs via `gh api repos/<r>/readme` | 2026-11-01 | [bedatadriven/renjin](https://github.com/bedatadriven/renjin) (+ `tools/gcc-bridge/README.md`), [oracle/fastr](https://github.com/oracle/fastr), [radfordneal/pqR](https://github.com/radfordneal/pqR), [jtalbot/riposte](https://github.com/jtalbot/riposte), [r-wasm/webr](https://github.com/r-wasm/webr) (+ `LICENSE.md`) |
| Oracle binaries | 2026-11-01 | [conda-forge `r-base`](https://api.anaconda.org/package/conda-forge/r-base) file list; CRAN [macOS arm64 directory](https://cran.r-project.org/bin/macosx/big-sur-arm64/base/); Posit [`versions.json`](https://cdn.posit.co/r/versions.json); `pkgx +r-project.org` |
| Fleet | `main`, 2026-11-01 | `gh repo list` of go-scikit, go-ndarray, go-statistics, go-numerics, go-plot, go-embedded-julia, go-embedded-r; `go-scikit/ruby` `go.mod`; `go-ndarray/ndarray` `linalg/` file list; `go-statistics/stats` file list; `go-notebook/notebook` `analysis/` and `go-notebook/docs` `content/en/` file lists; [`julia.md`](julia.md) §5 for the cell contract |

Counting method. The tarball was read only (no `configure`, no build). Every count in §1-§2 is
the output of [`tools/rcount`](../tools/rcount/main.go) (this repository; a test plants each
table form, a commented-out entry, an entry outside the table, a symbolic offset, and refuses an
empty table or a missing file; fuzzed):

```
go run ./tools/rcount --cran PACKAGES <R-4.6.2 root>
```

Verification: the `R_FunTab` count (773) equals the number of lines of `names.c` that begin with
`{"` (`grep -c`). A first scratch version that required a numeric offset column read 670 entries
and 182/488 primitives/internals: 28 entries have a symbolic offset (`PLUSOP`, …); the line count
exposed it. Line counts are physical lines (`wc -l`, plus one for a final unterminated line).

## 1. What R is made of (sizes, R 4.6.2)

### 1.1 Interpreter and numerical code

| part | path | files | lines | role for a Go port |
|---|---|---|---|---|
| interpreter, base built-ins, I/O, graphics engine | `src/main/*.c` | 128 | 156 069 C (incl. generated `gram.c` 6 858) + 117 Fortran | spec only (GPL, §3.2) |
| grammar + hand-written lexer | `src/main/gram.y` | 1 | 4 549 | spec only; §2.2 |
| evaluator | `src/main/eval.c` | 1 | 9 426 | AST `eval` at `eval.c:1082`; byte-code engine `bcEval` at `:7500`, `bcEval_loop` at `:7520` |
| key runtime files | `envir.c` 5 084, `memory.c` 5 130, `connections.c` 7 489, `engine.c` 4 017, `errors.c` 3 011, `coerce.c` 3 132, `grep.c` 3 524, `serialize.c` 3 383, `plotmath.c` 3 237, `unique.c` 2 690, `arithmetic.c` 2 466, `subassign.c` 2 340, `attrib.c` 2 122, `deparse.c` 1 962, `subset.c` 1 423, `names.c` 1 465 | | | |
| Rmath (d/p/q/r of distributions, special functions) | `src/nmath/` | 127 | 18 847 C | oracle only (GPL-2+, §3.2) |
| LINPACK + optimisers + `pretty` | `src/appl/` | 21 | 9 829 C + 2 337 Fortran | `dqrdc2.f` (QR used by `lm`), `optim.c`, `lbfgsb.c`, `uncmin.c`, `integrate.c` |
| reference BLAS | `src/extra/blas/` | 9 | 22 593 Fortran + 225 C | BSD (from LAPACK 3.12.0); replaced by go-ndarray |
| LAPACK subset | `src/modules/lapack/` | 17 | 250 684 Fortran + 2 739 C | BSD; base R calls **26** routines (§3.3) |
| public headers | `src/include/` | 58 | 16 758 C | `Rinternals.h` etc. are LGPL-2.1+ (§3.2) |

### 1.2 The base packages: lines per language

| package | R | C | Fortran | Rd (help) | native call sites in its R code |
|---|---|---|---|---|---|
| base | 23 469 (169 files) | (in `src/main`) | (in `src/appl`) | 47 029 | `.Internal` 514 sites / 409 names; `.Fortran` 7; `.Call` 3 |
| stats | 36 490 (181) | 23 087 | 18 206 | 30 775 | `.Call` 230/184; `.Fortran` 16/14; `.External` 15/8; `.C` 12/10; `.External2` 10/7 |
| utils | 23 589 (103) | 5 766 | — | 16 121 | `.Call` 32/30; `.External2` 26/21; `.External` 17/15; `.Internal` 4/3 |
| methods (S4) | 18 015 (45) | 1 553 | — | 11 928 | `.Call` 23/17; `.Internal` 5/3 |
| graphics | 8 341 (57) | 12 167 | — | 8 790 | `.External2` 4/4; `.Call` 3/3; `.External` 2/2 |
| grDevices | 6 837 (46) | 26 446 (`devPS.c` 10 117, `cairo/cairoFns.c` 2 586) | — | 8 059 | `.External` 48/25; `.Call` 25/22; `.External2` 14/9 |
| grid | 17 744 (45) | 11 483 (`grid.c` 5 470) | — | 7 183 | `.Call` 27/24 |
| tools | 53 678 (67) | 10 252 (+ 3 366 yacc) | — | 5 781 | `.Call` 34/19; `.Internal` 14/7 |
| compiler (byte-code compiler) | 4 303 (11) | — | — | 197 | — |

Base R is R code over a C boundary: `base` alone has 409 distinct `.Internal` names in 514 call
sites and almost no `.Call`. `stats` is the opposite: 41 293 lines of C and Fortran behind 230
`.Call` sites, i.e. the numerical half of `stats` is not R code at all.

### 1.3 The interpreter's surface

| item | count | source |
|---|---|---|
| `R_FunTab` entries | **773**, 498 distinct C functions (`do_*`) | `names.c:83-1028`; eval column `XYZ`, Y=1 internal, Z=1 builtin (`names.c:62-71`) |
| `.Primitive` | **210** (168 builtin, 42 special) | same |
| `.Internal` | **563** (555 builtin, 8 special) | same |
| `.Internal` names reached from `base`'s R code | 409 of 563 | §1.2 |
| SEXPTYPE codes | 24 (0-10, 13-25; `ANYSXP` 18 is a pseudo-type; `S4SXP` = `OBJSXP` 25) | `Rinternals.h:111-144` |
| grammar | 44 tokens, 22 precedence lines, **12 nonterminals**, 124 rule lines; the lexer is the 3 980-line C epilogue | `gram.y:399-418`, `:421-442`, `:444-569` |
| byte-code opcodes | 129 (`BCMISMATCH_OP` … `DECLNKSTK_OP`) | `eval.c:4602-4733` |
| JIT | on by default at level 3 | `eval.c:1460` ("turn JIT on by default"), overridable by `R_ENABLE_JIT` |
| test suite | `tests/`: 95 R files, 33 316 lines; 32 `.Rout.save` reference outputs | listing of `tests/` |

R-ints §.Internal vs .Primitive (1926-1975): primitives "make no use of R code", `formals` and
`body` return `NULL` for them, and argument matching "can be handled differently"; `.Internal`
functions are preferred "wrapped in a closure" because that "ensures standard handling of named
and default arguments". So the port's surface is 210 primitives with bespoke matching plus 563
internals called positionally from R closures.

### 1.4 CRAN: the native-code wall, measured

From the PACKAGES index (§0), counting each package name once:

| measure | packages | share |
|---|---|---|
| packages on CRAN | 25 608 | 100 % |
| `NeedsCompilation: yes` | 5 620 | 22 % |
| `LinkingTo: Rcpp` (C++) | 3 402 | 13 % |
| **pure R and every hard dependency (Depends, Imports, LinkingTo), transitively, pure R** | **3 851** | **15 %** |

A dependency missing from the index (Bioconductor, archived) counts as not runnable; base
packages count as runnable. So an interpreter with no native-code story runs 15 % of CRAN, not
the 78 % that need no compiler themselves: most pure-R packages import a compiled one.

Ranking compiled packages by how many packages depend on them (hard dependencies, transitively),
then providing the top *k* natively (a Go implementation of the package):

| top *k* provided natively | runnable packages |
|---|---|
| 0 | 3 851 |
| 1 (`rlang`) | 3 890 |
| 5 (+ `cli`, `glue`, `vctrs`, `Rcpp`) | 4 011 |
| 10 (+ `magrittr`, `lattice`, `utf8`, `tibble`, `Matrix`) | 4 235 |
| 20 | 5 484 |
| 50 | 7 440 |
| 100 | 10 099 (39 %) |

The ranking is greedy by reach, not an optimal cover; it shows the shape: the wall is a few
dozen foundational packages (tidyverse substrate, Rcpp, Matrix, lattice), then a long tail.

## 2. Reference semantics and architecture

### 2.1 Pipeline

```
source ──lexer+parser (Go, clean-room)──▶ R language objects (LANGSXP/SYMSXP/constants, srcrefs)
   ──▶ eval (AST walker): promises, argument matching, closures, specials/builtins, contexts
   ──▶ 210 primitives (Go) ──▶ .Internal dispatch ──▶ 563 internals (Go)
   ──▶ base/stats/utils R code (fetched, §2.6) ──▶ numerical kernels (go-ndarray, go-numerics,
       go-statistics) ──▶ graphics engine ──▶ SVG device
```

### 2.2 Parser

- R-lang §Parser (3977-4600) defines the three modes (REPL, `parse(file)`, `parse(text=)`,
  3994-4024), comments, tokens and expressions; `gram.y` fixes precedence (`:421-442`, 22 lines
  from `?` down to `(`/`[`/`[[`) and has only 12 nonterminals: the grammar is small, the lexer
  (3 980 lines of C) carries the detail. Its token list (`gram.y:401-418`) shows the
  context-dependent part: `=` is `EQ_ASSIGN`, `EQ_SUB` or `EQ_FORMALS` by position, a symbol is
  `SYMBOL`, `SYMBOL_FUNCTION_CALL`, `SYMBOL_SUB`, `SYMBOL_FORMALS` or `SYMBOL_PACKAGE`, plus
  `PIPE` (`|>`), `PIPEBIND`, `PLACEHOLDER`, `COMMENT`, `LINE_DIRECTIVE`, `INCOMPLETE_STRING`
  (the REPL's continuation state).
- `gram.y` is GPL and "based on the parser package, Copyright (C) 2009--2011 Romain Francois"
  (`doc/COPYRIGHTS:611-616`): not translatable. Write a recursive-descent or Pratt parser from
  R-lang plus the precedence table, and check it differentially.
- **Oracle**: `utils::getParseData()` (`src/library/utils/R/sourceutils.R:147`) returns the token
  table (line/column ranges, token kind, parent id, text) that `parse(keep.source=TRUE)`
  records; `deparse()` of the parsed expressions gives the tree. Compare both on every R file of
  the 4.6.2 tree (724 files in the nine packages of §1.2, plus 95 in `tests/`) and on CRAN
  sources.

### 2.3 Evaluator core

| topic | reference text | Go component |
|---|---|---|
| objects | R-ints §SEXPs (73); `Rinternals.h:111-144` | one Go struct per SEXPTYPE; attributes as a pairlist; `NA_integer_` = `INT_MIN`, `NA_real_` a NaN payload (to confirm on the oracle) |
| environments, lookup | R-ints §Environments (499); R-lang §Scope (1917) | frames as Go maps with insertion order where R exposes it (`ls(sorted=FALSE)`) |
| promises, lazy evaluation | R-ints §Argument evaluation (838-877): "replaces the actual argument list by a list of promises", defaults are "promises to be evaluated in the environment just created" | `Promise{expr, env, value, forced}`; `substitute`, `missing`, `force` |
| argument matching | R-lang §Functions / Evaluation (2171) | exact, partial, positional, `...`, `..1` |
| contexts, non-local exits | R-ints §Contexts (738); R-lang §Exception handling (3718) | Go `panic`/`recover` per context, never across a wasm callback; `on.exit`, `sys.call`, `sys.function`, `parent.frame` |
| S3 | R-lang §OOP (2460-2975); R-ints §Argument evaluation (864-877) on how a method sees the generic's promises | `UseMethod`, `NextMethod`, group generics `Math`/`Ops`/`Summary`/`Complex` |
| S4 | R-ints §S4 objects (1386); `methods` is 18 015 lines of R + 1 553 C | late milestone |
| copy semantics | R-ints §Argument evaluation (862): "The `NAMED` mechanism has been replaced by reference counting" | reference counts on vectors for copy-on-modify; observable through `tracemem` only |
| visibility | R-ints §Visibility (1773); `names.c:62-66` (eval column X) | per primitive, from the table |
| computing on the language | R-lang §Computing on the language (2975-3597) | `quote`, `bquote`, `eval`, `body<-`, `formals<-`, `match.call`, `sys.call` |

No byte-code in v1: the 129-opcode engine (`eval.c:4602-4733`, `bcEval` `:7500`) and the
`compiler` package (4 303 lines of R) are a speed device; the AST evaluator is the semantics.
R itself runs with the JIT on (`eval.c:1460`), so the performance reference (fleet rule) is
GNU R **with** its JIT; a Go byte-code tier is a later, measured decision (§6).

### 2.4 Primitives and internals

The 210 primitives are the language (`if`, `for`, `function`, `<-`, `[`, `[[`, `$`, `@`,
arithmetic and comparison operators, `c`, `list`, `length`, `attr`, `class`, `UseMethod`,
`standardGeneric`, `.Internal` itself …) and are written first. The 563 internals are written in
the order base R's own code reaches them (409 names from `base/R`), measured by the oracle
harness: the first target is the set needed to load `base` without error. Riposte's README (l.
54-55) put that set at "about 15 primitive and external functions" for its own design; measure it
for ours instead of assuming.

### 2.5 Numerical boundary inside base R

| R facility | in R 4.6.2 | Go route |
|---|---|---|
| `sum` of doubles, other reductions | `rsum` accumulates in `LDOUBLE` (`summary.c:152-154`; also `:128`, `:176`, `:336`) | on arm64 macOS `long double` is `double`, on x86-64 Linux it is 80-bit extended: results are platform-dependent. Bit parity per platform needs a software 64-bit-mantissa accumulator; decide which platform the oracle is (§6) |
| `qr`, `lm`, `glm` | `qr(LAPACK=FALSE)` and `lm.fit` go through `dqrdc2`, "a *modification* of Linpack's dqrdc … for R" with "a limited column pivoting strategy" (`src/appl/dqrdc2.f:1-15`) | port LINPACK `dqrdc` (public domain status "unclear", `doc/COPYRIGHTS:619-627`) into go-ndarray and reproduce dqrdc2's pivoting rule from its documented behaviour, checked on the oracle; `qr(LAPACK=TRUE)` maps to go-ndarray QR |
| `solve`, `chol`, `svd`, `eigen`, `det`, `qr(LAPACK=TRUE)`, `rcond`, `norm` | 26 LAPACK routines (§3.3) | go-ndarray `linalg` (`lu.go`, `qr.go`, `cholesky.go`, `svd.go`, `eig.go`, `eigh.go`, `lstsq.go`, `norm.go`); missing pivoted Cholesky `dpstrf`, pivoted QR `dgeqp3`, condition estimators `dgecon`/`dtrcon` go **into go-ndarray** |
| d/p/q/r distribution functions, `gamma`, `beta`, `choose`, Bessel | Rmath, `src/nmath/` 18 847 lines, GPL-2+ | go-statistics/stats (distributions, "accuracy measured against mpmath") and go-numerics/special (gamma, beta, erf, incomplete gamma/beta, Bessel, zeta); parity with R is a **tolerance**, not bit equality (§4) |
| RNG | Mersenne-Twister default, normal by inversion, `sample` by rejection (`RNG.c:31-32`, `:50`) | MT19937 from its authors' reference (the MT part of `RNG.c` is "Copyright (C) 1997, 1999 Makoto Matsumoto and Takuji Nishimura", `doc/COPYRIGHTS:264-267`); seed scrambling, inversion and rejection sampling reproduced from `?RNG`/`?set.seed` and checked stream-for-stream on the oracle. go-statistics draws match **NumPy**, not R: the R streams are new code |
| regex `perl=FALSE` | TRE 0.8.0, BSD (`doc/COPYRIGHTS:446-452`), `src/extra/tre` 8 939 lines | **portable**: translate TRE to Go (POSIX ERE semantics, approximate matching for `agrep`) |
| regex `perl=TRUE` | PCRE2 (binary distributions, `doc/COPYRIGHTS:688-763`) | PCRE2 syntax front end on [go-regexp/engine](https://github.com/go-regexp/engine), as `julia.md` §3.3 plans for Julia; share it |
| `serialize`, `saveRDS`, `load` | R-ints §Serialization Formats (1092); XDR from `src/extra/xdr` (BSD) | Go reader/writer from the manual; also the oracle's wire format (§4) |
| `radixsort` (`order`, `sort(method="radix")`) | `src/main/radixsort.c` "largely based on code from the data.table package" (`doc/COPYRIGHTS:634-636`) | reimplement; results are defined (stable order) so equality is exact |

### 2.6 Base library strategy: the licence decides it

Julia's plan runs Base unchanged through the interpreter because Base is MIT (`julia.md` §2.6).
R's base packages are GPL-2+: `src/library/base/R/apply.R:4-9` is "free software … GNU General
Public License … version 2 … or (at your option) any later version", and `doc/COPYRIGHTS:8-14`
says "R as a whole is distributed under GPL version 2 or 3". The fleet rule (TEMPLATE §3) is
GPL = oracle only. Three options:

| option | interpreter (Go) | base/stats/utils R code (≈ 100 000 lines) | wasm/notebook | cost |
|---|---|---|---|---|
| **A. clean-room everything** | BSD-3 | rewritten in Go or new R, from manuals + Rd + oracle | self-contained | rewriting 23 469 + 36 490 + 23 589 + 18 015 lines of R whose behaviour *is* the definition of R; Riposte planned exactly this ("reimplementing R internal functions in a limited subset of R", README l. 48) and stopped |
| **B. GPL port** | GPL-2+ (C may then be translated) | shipped | GPL-3 binaries, like webR (`LICENSE.md:15`) | fastest and most faithful; makes go-notebook binaries that embed it GPL; breaks the fleet's BSD-3 rule |
| **C. BSD interpreter, R code fetched** | BSD-3, clean-room (primitives, internals, parser, evaluator written from R-lang, R-ints, Rd pages and oracle behaviour; no `src/main` file open while its Go counterpart is written) | **not redistributed**: fetched from the pinned CRAN tarball on first use, checked against its SHA-256, cached — the go-tex/texmf pattern ("it redistributes nothing … a download is not a redistribution", fleet capability map) | the browser fetches the R files from a CORS-enabled mirror, then evaluates them (or a cached snapshot built on that machine) | the interpreter is the large part; base R code runs unchanged, so upstream fixes come for free |

Recommendation: **C**, with B kept as the fallback if C's legal reading does not hold. This is
a licence question to put to the user before M2 (§6): this file does not settle it. Under C the
start-up image of `julia.md` §6 must be built on the user's machine (it contains GPL code),
never shipped.

### 2.7 Graphics: an SVG device

R-ints §Graphics Devices (2795-2839): a device is a `DevDesc` holding "a set of callback
functions that will be used by the graphics engine to obtain information about the device and
initiate actions (e.g. a new page, plotting a line or some text)", several optional (`activate`, `cap`, `deactivate`, `locator`,
`holdflush`, `mode`, `newFrameConfirm`, `path`, `raster`, `size`); `GEDevDesc` adds the display
list. The engine (`engine.c` 4 017, `plotmath.c` 3 237) sits between the device and two
systems: base graphics (`graphics` 8 341 R + 12 167 C) and grid (17 744 R + 11 483 C; ggplot2
and lattice draw through grid). Plan: a Go graphics engine with one device writing SVG (text
metrics from go-opentype fonts, so `strwidth` agrees with what the SVG shows), and
`png`/`pdf` later through go-gfx and the fleet PDF writer. R's own `svg()` is a Cairo device
(`grDevices/src/cairo/cairoFns.c`), not a model to port. Exit test: base and grid example plots
compared to the oracle's `svg()` output structurally (element count, coordinates within a
tolerance), not byte for byte.

## 3. Reference implementations and prior art

### 3.1 Other R implementations

| project | approach | native code (the CRAN wall) | licence | state (as read) |
|---|---|---|---|---|
| [Renjin](https://github.com/bedatadriven/renjin) | JVM interpreter; "a drop-in replacement for GNU R" (README 8-11) | **GCC-Bridge**: GCC front end → Gimple dumped as JSON → JVM class files, with a partial C library runtime (`tools/gcc-bridge/README.md:4-16`) | GPL-2+, "includes code from R which is redistributed here" (README 73-79) | last push 2023-09-26 |
| [FastR](https://github.com/oracle/fastr) | Truffle/GraalVM; "based on the GNU-R runtime" (README 74) | implements the R extensions C API; packages using it in hot paths "may exhibit slower performance … due to the high cost of transitions between the native and managed code", mitigated by the GraalVM LLVM runtime (README 48-54); runs `ggplot2`, `dplyr`, `Rcpp`, `Shiny` examples (38) | GPLv3 (README 87) | pushed 2026-07-30 |
| [pqR](https://github.com/radfordneal/pqR) | "a modification of R-2.15.0" for speed (README 5-11) | inherits GNU R's C API | GPL-2 or 3 (README 17-19) | last push 2020-07-23 |
| [Riposte](https://github.com/jtalbot/riposte) | new VM + tracing JIT for vector code (README 48) | "Initial support for R's C-level API … remains one of the big unknowns" (README 69-70) | BSD-2 (README 34) | last push 2015-09-24; roadmap: ~200 primitives (50 done), ~580 internals (30 done), July 2013 (README 57-62) |
| [webR](https://github.com/r-wasm/webr) | GNU R compiled to wasm with Emscripten, patches to R's source (README 4); Fortran through LLVM flang (21) | packages compiled to wasm too | R parts GPL; "webR binaries are distributed under the GNU GPL Version 3" (`LICENSE.md:5,15`) | pushed 2026-10-27 |

What they learned, as read: every implementation that reached real-world R code either reused
GNU R's own base library and C (Renjin, FastR, pqR, webR — and therefore is GPL) or did not get
past the internals (Riposte). Every one that took on CRAN had to answer native code with a
compiler (GCC-Bridge, LLVM/Sulong, Emscripten). §1.4 puts a number on why: 85 % of CRAN is
behind a compiled package. A CGO-free Go interpreter has no such compiler; its answers are §5's
M8 (Go ports of the top blockers; a C-to-Go translation route was not studied here).

### 3.2 Licences: what is portable

| component | licence (source) | verdict |
|---|---|---|
| interpreter `src/main`, Rmath `src/nmath`, base package R code, `gram.y` | GPL-2+ (`doc/COPYRIGHTS:8-14`; file headers) | oracle only |
| `R.h`, `Rinternals.h`, `Rmath.h`, `R_ext/*.h` | LGPL-2.1+ (`doc/COPYRIGHTS:27-37`) | facts (type codes, API names) usable; not needed beyond SEXPTYPE codes and the serialization spec |
| LAPACK, reference BLAS | modified BSD (`src/modules/lapack/LICENSE`; `doc/COPYRIGHTS:289-360`); R's copy is LAPACK 3.12.1 and R "announced the intention never to update this collection" (`src/modules/lapack/README:1-9`) | portable, but go-ndarray owns linear algebra: add missing routines there |
| LINPACK (`src/appl/d*.f`) | "status … unclear: they are said to be in the public domain in the USA" (`doc/COPYRIGHTS:22-23`, `619-627`); `dqrdc2.f` is R Core's modification | port `dqrdc` from netlib; reproduce `dqrdc2`'s rule from behaviour |
| TRE regex | BSD-2 (`doc/COPYRIGHTS:446-481`) | portable |
| XDR | BSD (`doc/COPYRIGHTS:362-395`) | trivial in Go |
| tzone `strftime` | UCB BSD-style (`doc/COPYRIGHTS:483-502`) | portable; Go `time` covers most |
| Mersenne-Twister | Matsumoto & Nishimura (`doc/COPYRIGHTS:264-267`) | implement from the authors' reference |
| AS algorithms in nmath/stats (`qnorm`, `qbeta`, `pnt`, `swilk`, …), Loader's `dbinom`/`dpois` | "based in whole or in part on Applied Statistics algorithms (C) Royal Statistical Society"; Loader (C) Lucent (`doc/COPYRIGHTS:240-261`) | implement from the published papers (not opened here), not from nmath |
| R manuals (R-lang, R-ints) and Rd pages | documentation | the specification for a clean-room implementation |

### 3.3 What the fleet already has

| need | fleet library (observed 2026-11-01) | gap |
|---|---|---|
| dense linear algebra | go-ndarray/ndarray `linalg` (LU, QR, Cholesky, SVD, Eig, Eigh, lstsq, norm, BLAS kernels) | of the 26 routines `Lapack.c` calls (`dgecon dgeev dgeqp3 dgesdd dgesv dgetrf dlange dormqr dpotrf dpotri dpstrf dsyevr dtrcon dtrtrs ilaver` + 11 complex `z*`): pivoted QR, pivoted Cholesky, condition estimators, complex variants; LINPACK `dqrdc`/`dqrdc2` |
| distributions, tests | go-statistics/stats (normal … zipf, KS, t-tests, rank tests, KDE, bootstrap) | R's parameterisations and argument conventions (`lower.tail`, `log.p`); R-identical RNG streams |
| special functions | go-numerics/special | `choose`/`lchoose`, `psigamma` orders, `besselI(expon.scaled)` conventions, checked on Rmath as oracle |
| optimisation, integration | go-numerics/optimize, go-numerics/integrate | R's `optim`, `nlm`, `integrate` have their own algorithms (`src/appl/optim.c`, `lbfgsb.c`, `uncmin.c`, `integrate.c`, GPL) whose iterates users see: re-derive from the published methods, tolerance parity |
| data frames | go-dataframe/dataframe | the bridge to R `data.frame` (§5) |
| FFT | go-fft/fft | `fft`, `mvfft`: tolerance parity |
| plots | go-plot/plot (SVG/PNG figures), go-gfx/gfx | an R graphics *engine* and device (§2.7) is new |
| regex | go-regexp/engine (Onigmo) | TRE port; PCRE2 front end shared with go-embedded-julia |

## 4. Validation

### 4.1 Oracle: GNU R 4.6.2

| route | macOS arm64 | Linux x86-64 | root? | status |
|---|---|---|---|---|
| conda-forge `r-base` via micromamba | 4.6.2 published for `osx-arm64` | 4.6.2 for `linux-64` | no | **recommended**: one pinned version on both (api.anaconda.org file list, 2026-11-01) |
| pkgx | `r-project.org` recipe at 4.6.1, "not published" for darwin/aarch64 | 4.6.1 | no | not usable on the dev machine today; a go-pkgx recipe bump is the sovereign route later |
| CRAN macOS `.pkg` | `big-sur-arm64/base/` lists up to `R-4.5.3-arm64.pkg` | — | installer needs admin | no |
| Posit builds | — | `versions.json` lists 4.6.2 | no (tarball) | Linux only; not opened beyond the list |
| build from the tarball | gfortran + toolchain | same | no | the fallback; slow |

Harness (M0): a Go program that starts **one** `R --vanilla --slave` process per worker and
speaks a framed protocol on stdin/stdout. Requests are R source text; replies are the value
`serialize(x, NULL, xdr = TRUE)` (R-ints §Serialization Formats) plus captured stdout, warnings
and the condition class/message. Using R's own serialization as the wire format avoids a JSON
package on the oracle side (jsonlite is a compiled CRAN package) and makes the Go RDS reader an
M0 deliverable that the interpreter needs anyway. Timeouts and memory limits per worker; golden
files regenerated by one command. The harness was **not run** in this session (fetching
micromamba was refused by the sandbox): every oracle statement here is a plan.

### 4.2 Corpora

- R's own `tests/`: 95 R files (33 316 lines), 32 with `.Rout.save` reference transcripts
  (`arith`, `eval-etc`, `method-dispatch`, `reg-tests-*`, `print-tests`, `d-p-q-r-tests`,
  `lapack`, `lm-tests`, …); R-ints §Testing R code (4878-4890) names `make test-Specific` and
  `tests/no-segfault.Rin`'s stop list. These are the conformance suite, compared output line by
  output line.
- Examples in the Rd files of the base packages (47 029 lines of Rd in `base` alone), run under
  both.
- Parser corpus: the 724 R files of §1.2's nine packages, `tests/`, and the R sources of the
  3 851 runnable CRAN packages (§1.4).
- Fuzzing: parser, RDS reader (gzip-compressed by default: decompression bombs must be
  refused), regex front ends.

### 4.3 What "correct" means

| layer | measure |
|---|---|
| parser | token table (`getParseData`) and `deparse` output identical |
| evaluator, base | printed output identical to `.Rout.save` / the oracle, warnings and error messages included |
| integer, logical, character, sorting, set operations, RNG streams | exact equality |
| double arithmetic in R code | exact where R computes in `double`; `sum`/`mean` per the `LDOUBLE` decision (§2.5) |
| Rmath via go-statistics / go-numerics | ULP distance to R reported per function, and to mpmath; the better of the two is not a defect, a worse one is |
| linear algebra via go-ndarray | backward error and residual; equality of `lm` coefficients to a stated tolerance; signs of QR/eigenvectors normalised before comparing |
| graphics | SVG structure and coordinates within a tolerance |

## 5. Proposed home and plan

Effort in the `bib/` scale: S days, M weeks, L months, XL several months. No repository was
created by this file; names follow the go-embedded-julia layout.

| # | work item | org/repo | depends on | size | priority | exit test |
|---|---|---|---|---|---|---|
| M0 | oracle harness + RDS (serialize v3) reader/writer | go-embedded-r/r (`rds/`, `internal/oracle`) | conda-forge r-base 4.6.2 | S | 1 | golden values round-trip R → Go → R; RDS fuzzed |
| M1 | lexer + parser, srcrefs, `deparse` | go-embedded-r/r (`syntax/`) | M0 | M | 1 | `getParseData` + `deparse` identical on the 724 + 95 files and on CRAN sources; fuzzed with depth and size bounds |
| M2 | evaluator core: SEXPs, environments, promises, matching, closures, contexts/conditions, S3, the 210 primitives | go-embedded-r/r (`eval/`, `prim/`) | M1; **licence decision** (§2.6) | XL | 1 | `tests/simple-true.R`, `arith-true.R`, `eval-etc.R`, `any-all.R`, `method-dispatch.R` against `.Rout.save` |
| M3 | base: the 563 internals in reach order, base R code loaded per §2.6, lazy-load/start-up snapshot | go-embedded-r/r (`internal/…`, `base/`) | M2 | XL | 1 | `base` loads; `reg-tests-1*.R`, `print-tests.R`, `datetime*.R`, `complex.R`, `iec60559.R`, `utf8-regex.R`; base Rd examples |
| M4 | stats, utils, methods: Rmath on go-statistics/go-numerics, R RNG streams, `lm`/`glm` (dqrdc2 rule), `optim`, `nls`, formulas and model frames; S4 | go-embedded-r/r (`stats/`), kernels into go-ndarray / go-statistics / go-numerics | M3 | XL | 2 | `d-p-q-r-tests.R`, `p-r-random-tests.R` (streams exact), `lm-tests.R`, `lapack.R`, `reg-S4.R`; stats Rd examples |
| M4b | regex: TRE port (`perl=FALSE`, `agrep`) and PCRE2 front end on go-regexp/engine | go-regexp (new `tre` repo) + go-regexp/engine | — | M | 2 | `tests/utf8-regex.R`, `PCRE.R`, base `grep`/`regexpr`/`gsub` examples |
| M5 | graphics engine + SVG device; base graphics, grid | go-embedded-r/r (`graphics/`, `grdevices/svg`) | M3, go-opentype | L | 2 | base and grid example plots vs the oracle's `svg()` (§2.7); ggplot2 later (M8) |
| M6 | notebook bridge + go-scikit bindings | go-notebook/notebook (`analysis/r.go`, next to `julia.go`, `ruby.go`), go-scikit/r (new, next to go-scikit/ruby) | M3 | M | 2 | a mixed Go/R notebook transcript in go-notebook docs (`content/en/r.md`, next to `julia.md`, `ruby.md`) |
| M7 | wasm playground: lazy R core, base fetched and cached (§2.6), SVG through the existing sanitizer | go-notebook/playground | M5, M6 | M | 3 | first result time and size recorded against the Ruby core |
| M8 | CRAN: Go ports of the top blockers (`rlang`, `cli`, `glue`, `vctrs`, `magrittr`, `utf8`, `tibble`, …; `Rcpp` cannot be "ported", only its users) | `go-r-<package>` orgs, as `go-ruby-<gem>` | M3 | XL (open-ended) | 3 | each package's own `tests/testthat` under the Go interpreter; the runnable count of §1.4 re-measured |
| — | site repos | go-embedded-r/{docs, go-embedded-r.github.io, .github, brand} via `ghnew` | — | S | 1 | org conformance |

Value bridge (M6), following the Ruby/Julia contract (`julia.md` §5.1): one shared interpreter
per notebook; names a cell assigns (`<-`, `=`, `->`, `<<-` at top level, from the M1 parser) are
its definitions; `NULL` ↔ nil; length-1 logical/integer/double/complex/character ↔ `bool`,
`int32`/`int`, `float64`, `complex128`, `string` with NA as a distinct Go value (never a silent
zero); numeric vectors ↔ `[]float64` or `*ndarray.Array`; matrices and arrays (`dim` attribute)
↔ `*ndarray.Array` with **column-major** strides (the question `julia-runtime.md` §6.2 answered
for Julia: a view yes, without a copy no); `data.frame` ↔ go-dataframe; named lists ↔ ordered
maps; factors ↔ codes + levels; anything else an opaque `r.Value`. The VM exposes a public
extension API outside `internal/` from day one (`bib/ruby-reuse.md` §3), so go-scikit/r can
register Go functions as R closures in their own namespace.

## 6. Risks and open questions

- **Biggest risk: the licence of base R** (§2.6). Option C's premise, that a BSD interpreter may
  fetch and run unmodified GPL R code it does not redistribute, is a legal reading, not a
  verified opinion. If it fails, the choice is between A (rewrite ≈ 100 000 lines of R whose
  behaviour defines the language — Riposte's unfinished road) and B (GPL, like every
  implementation that reached real code). Decide before M2.
- **The CRAN wall** (§1.4): 15 % of CRAN runs on a pure interpreter; 100 Go package ports reach
  39 %. "R in the notebook" must be promised as base + recommended-like numerics + the ported
  set, not as CRAN.
- **Clean-room discipline**: implementers must work from manuals, Rd pages and oracle behaviour,
  never with `src/main` or `src/nmath` open; this file read them only to count. The rule must be
  in the repository's CONTRIBUTING before the first commit.
- **Platform-dependent R results**: `LDOUBLE` accumulation (`summary.c:152-154`) differs between
  the arm64 macOS and x86-64 Linux oracles. Pick the oracle platform per CI job and record it;
  bit parity on x86-64 needs a software extended accumulator.
- **Interpretive speed**: GNU R runs with its byte-code JIT on (`eval.c:1460`); an AST walker
  will be slower on scalar loops. Vector operations (most R code) route to go-ndarray/go-simd
  kernels and can be at par; report the ratio per milestone against R single-core, and open a
  defect where slower, per the fleet rule.
- **Rmath parity is by tolerance**: go-statistics is accurate against mpmath and draws like NumPy;
  R users compare printed digits. Where R and mpmath disagree, document it; never copy nmath.
- **Start-up**: R loads base from a lazy-load database (R-ints §Lazy loading, 1852); under
  option C the Go equivalent is a snapshot built on the user's machine, so the first run in a
  browser pays parsing and evaluating base R. Unknown until M3; measure against the Ruby core.
- **Security** (untrusted cells): fuzz the parser and RDS reader; bound recursion depth,
  vector allocation (refuse beyond a configured limit with R's own "cannot allocate vector"
  error class), regex steps; RDS bombs
  refused under `GOMEMLIMIT`; `system()`, `file()`, `url()` absent on wasm and policy-gated on
  the desktop.
- **Unverified in this file**: everything about running the oracle (not run); the Posit Linux
  builds beyond their version list; the published AS / TOMS / Loader papers (cited through
  `doc/COPYRIGHTS` only); NA bit patterns; a C-to-Go translation route for CRAN native code
  (not studied); the claim that `dqrdc2`'s rule can be reproduced from behaviour alone.

Verified on 2026-11-01: every count in §1 and §2 by `tools/rcount` on the 4.6.2 tarball of §0;
licence texts read in `doc/COPYRIGHTS`, `src/modules/lapack/LICENSE` and file headers; manual
sections read at the cited lines; prior-art READMEs read through the GitHub API; fleet facts via
`gh repo list` and the cited files.
