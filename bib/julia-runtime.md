# Julia M3 — bibliography and design of the Core runtime

Written 2026-10-27. The design of milestone M3 of [`julia.md`](julia.md) §8 (types, subtyping,
method tables, dispatch, builtins, intrinsics, the interpreter over lowered IR), against the
release the port is pinned to. Facts only, each with the file and lines it was read from.
No M3 code exists yet; M2 (lowering) is in progress in
[go-embedded-julia/julia](https://github.com/go-embedded-julia/julia), whose `expr` package
already holds the lowered-IR values this document builds on.

## 0. Sources read

| source | revision | what was used |
|---|---|---|
| [JuliaLang/julia](https://github.com/JuliaLang/julia) tag `v1.12.6` | `15346901f003`, `VERSION` = `1.12.6`: the commit go-embedded-julia pins | `src/interpreter.c`, `src/builtin_proto.h`, `src/intrinsics.h`, `src/julia.h`, `src/gf.c`, `src/subtype.c`, `src/toplevel.c`, `doc/src/devdocs/ast.md`, sizes of `base/boot.jl` and `test/*.jl` |
| go-embedded-julia/julia, branch of M2 (1/n), commit `da6eab0` | 2026-10-27 | `expr/ir.go`: the Go types of the lowered IR |
| [go-ndarray/ndarray](https://github.com/go-ndarray/ndarray) `main` | 2026-10-27 | `ndarray.go` (`FromData`, `Transpose`), `storage.go` (`FromSlice`, `Data`) |
| [go-notebook/notebook](https://github.com/go-notebook/notebook) `main` | 2026-10-27 | `analysis/ruby.go` (the per-language cell analysis a Julia cell will mirror) |
| `gh repo list go-bignum` | 2026-10-27 | `bigfloat` is the only library repository |

Counts come from `go run ./tools/jlcount <julia-tree>` (this repository, with a test that plants
each table form and refuses an empty table) and from `wc -l`, which agreed on the two files
checked by hand (`subtype.c` 5 556, `interpreter.c` 952).

## 1. Corrections to `julia.md`

`julia.md` §0 read the julia **master** of 2026-10-26 (`1.15.0-DEV`), not the release the port
was then pinned to. On v1.12.6 the runtime surface is smaller, and some figures change:

| figure in `julia.md` | 1.15.0-DEV (as written) | **v1.12.6** (pinned) |
|---|---|---|
| builtins in `src/builtin_proto.h` | 83 | **63** (`builtin_proto.h:11-74`) |
| intrinsics in `src/intrinsics.h` | 94 (81 `ADD_I` + 12 `ALIAS` + 1 hidden) | **93** (80 `ADD_I` + 12 `ALIAS` + 1 `ADD_HIDDEN`, `intrinsics.h:3-110`); `cpu_supports` does not exist yet, `have_fma` is the only CPU test |
| semantic core, 15 C files | 36 477 lines | **29 963** (`subtype.c` 5 556, `gf.c` 5 009, `jltypes.c` 3 993, `builtins.c` 2 598, `datatype.c` 2 355, `module.c` 2 099, `runtime_intrinsics.c` 1 756, `method.c` 1 607, `typemap.c` 1 521, `toplevel.c` 1 269, `interpreter.c` 952, `genericmemory.c` 630, `array.c` 339, `symbol.c` 168, `simplevector.c` 111) |
| `base/boot.jl` | 1 338 | **1 086** |
| oracle test files | `subtype.jl` 3 666, `core.jl` 9 501 | `subtype.jl` **2 788**, `specificity.jl` 329, `intrinsics.jl` 566, `core.jl` **8 686** |
| M2 ports JuliaLowering (§2.3) | recommendation | **superseded**: v1.12.6 lowers with `src/julia-syntax.scm` by default, so M2 ports it (go-embedded-julia README, "Lowering parity") |
| go-ndarray accepts Fortran strides (§10, unverified) | open | **answered** in §6.2 below: as a view yes, without a copy no |

The prose of `julia.md` is otherwise unchanged; every M3 count below is of v1.12.6.

## 2. What the interpreter must execute

### 2.1 The reference: `src/interpreter.c` (952 lines)

julia's own interpreter is two functions. `eval_value` (`interpreter.c:198-364`) evaluates an
argument; `eval_body` (`:470-712`) walks the statements with a program counter. Its state is
`interpreter_state` (`:17-27`): the `CodeInfo`, the module, one `locals` array holding the slots
**then** the SSA values (`s->locals[nslots + ip]`, `:195`), the static parameters, `ip`, and
`continue_at` for `:leave`.

| statement / value | where | semantics to reproduce |
|---|---|---|
| `SSAValue`, `SlotNumber`, `Argument` | `:201-216` | out of range is `ErrorException("access to invalid SSAValue")` / `"…slot number"`; an unassigned slot is `UndefVarError(name, :local)` |
| `QuoteNode`, `GlobalRef`, bare `Symbol` | `:217-225` | a global read in the **task's world age** (`jl_current_task->world_age`) |
| `GotoNode`, `GotoIfNot` | `:484-495` | 1-based labels; a non-`Bool` condition is `TypeError(:if, Bool, cond)` |
| `ReturnNode` | `:496-498` | |
| `EnterNode` | `:508-575` | pushes a handler (`setjmp`); optional scope (`ScopedValue`); the statement's value is the exception-stack depth, consumed by `pop_exception` |
| `:leave` | `:595-623` | pops one handler per argument that is an `SSAValue` of an `EnterNode` not replaced by `nothing` |
| `:pop_exception` | `:624-627` | restores the exception stack |
| `NewvarNode` | `:696-702` | resets a slot to undefined |
| `:(=)` | `:579-594` | the left side **must** be a slot: `"Invalid IR: Assignment LHS not a Slot"` (globals are assigned through the `setglobal!` builtin) |
| `:call`, `:invoke`, `:invoke_modify` | `:242-250`, `:115-162` | `jl_apply` on the evaluated arguments; `invoke` targets a `MethodInstance`/`CodeInstance` |
| `:new`, `:splatnew` | `:293-310` | allocation; `ast.md:418-420`: "does no checking. Evaluating arbitrary `new` expressions can easily segfault" |
| `:isdefined`, `:throw_undef_if_not`, `:static_parameter` | `:251-292`, `:322-333` | |
| `:copyast`, `:the_exception`, `:boundscheck` (always `true`) | `:334-342` | |
| `:meta`, `:inbounds`, `:loopinfo`, `:aliasscope`, `:popaliasscope`, `:inline`, `:noinline`, `:coverageeffect`, `:gc_preserve_begin/_end` | `:343-352` | `nothing` |
| `:new_opaque_closure` | `:311-321` | |
| `:method` | `:80-111`, `:353-354`, `:629-631` | 1 argument: declare the generic function; 3 arguments (top level only): `jl_method_def(atypes, mt, codeinfo, module)` |
| top level only: `:toplevel`, `:global`, `:const`, `:latestworld`, `:module`, `:import`, `:using`, `:export`, `:public`, `:thunk` | `:628-690`, `toplevel.c:572-584` | `:latestworld` moves the task to the newest world (`:651-653`) |
| `:foreigncall`, `:cfunction` | `:356-361` | **an error in julia's interpreter**: "`ccall` requires the compiler" |
| `PiNode`, `PhiNode`, `PhiCNode`, `UpsilonNode` | `:226-235`, `:367-468`, `:499-507` | only in **optimised** IR; `Meta.lower` does not produce them |

Two consequences for the Go port:

1. **The lowered IR M2 produces needs no φ-nodes.** The interpreter runs unoptimised IR, like
   `jl_interpret_toplevel_thunk` (`:884`), so `eval_phi` (`:367-468`, the most intricate part of
   the file) is not ported in M3. It is needed only if optimised IR is ever executed.
2. **`foreigncall` must be implemented**, where julia's interpreter gives up. julia never meets
   the error because Base is compiled; a Go interpreter runs Base itself, so every `ccall` of
   Base reaches it. The shim table of `julia.md` §3.2 is part of the interpreter, keyed by the
   literal callee name of `args[1]` (`ast.md:481-512`).

### 2.2 The contract with M2 (what `lower` must hand over)

The M2 types (`expr/ir.go`, commit `da6eab0`) mirror julia's serialisation: `SSAValue`,
`SlotNumber`, `Argument`, `GotoNode`, `GotoIfNot{Cond, Dest}`, `ReturnNode{Val}`, `Unreachable`,
`NewvarNode{Slot}`, `EnterNode{Catch, Scope, HasScope}`, `CodeInfo{Code, DebugInfo,
SSAValueTypes, SSAFlags, SlotNames, SlotFlags, …}` with most fields typed `any`. That is right
for an oracle; M3 asks four things of it and nothing else:

| requirement | why |
|---|---|
| `lower.Lower` returns the same `CodeInfo` as `Meta.lower` (M2's own exit test) | the interpreter is written against julia's IR, statement for statement (§2.1) |
| the statement heads are among those of §2.1 | an unknown head is julia's `"unsupported or misplaced expression"` (`interpreter.c:362`), which M3 reproduces — but a head M2 invents would be a silent divergence |
| `SlotNames`/`SlotFlags` have one entry per slot; `SSAValueTypes` is the statement count | `jl_source_nslots` = `length(slotflags)`, `jl_source_nssavalues` = the count (`interpreter.c:182-190`) |
| `DebugInfo` is kept | backtraces and the notebook's `file:line` (`julia.md` §5.1) |

M3 does **not** execute `expr.CodeInfo` directly. A `prepare` step converts it once per method
into an executable form: statements as a Go sum type (a `switch` on concrete types, not on
`any`), slot and SSA counts as `int`, every `SSAValue`/`SlotNumber`/label **checked in range
once**, heads interned. julia checks slot bounds in `=` only with `assert` (`interpreter.c:584`),
so a malformed `CodeInfo` built by user code (`Core.eval` of an `Expr(:thunk)`) is unchecked in a
release build; in Go the check is the cost of one pass at preparation and leaves the hot loop
free of it.

## 3. The runtime boundary

### 3.1 Builtins: 63 (`builtin_proto.h:11-74`)

Each is a Go function `func(vm *VM, args []Value) (Value, error)` in a table indexed like
`jl_builtin_ids` (`:83-88`), so that `Core.Builtin` instances keep julia's order.

| group | builtins |
|---|---|
| type construction (called by lowered `struct`/`abstract type` code, `ast.md:404-406`) | `_abstracttype`, `_primitivetype`, `_structtype`, `_setsuper!`, `_typebody!`, `_equiv_typedef`, `_defaultctors`, `_typevar`, `apply_type`, `_compute_sparams` |
| type queries | `typeof`, `isa`, `<:`, `typeassert`, `fieldtype`, `nfields`, `sizeof`, `===` |
| fields | `getfield`, `setfield!`, `setfieldonce!`, `swapfield!`, `modifyfield!`, `replacefield!`, `isdefined` |
| globals | `getglobal`, `setglobal!`, `setglobalonce!`, `swapglobal!`, `modifyglobal!`, `replaceglobal!`, `isdefinedglobal`, `get_binding_type` |
| memory | `memorynew`, `memoryrefnew`, `memoryrefget`, `memoryrefset!`, `memoryrefsetonce!`, `memoryrefswap!`, `memoryrefmodify!`, `memoryrefreplace!`, `memoryref_isassigned`, `memoryrefoffset` |
| calls | `_apply_iterate`, `invoke`, `invoke_in_world`, `invokelatest`, `_call_in_world_total`, `applicable`, `intrinsic_call`, `opaque_closure_call`, `ifelse` |
| values | `tuple`, `svec`, `_svec_len`, `_svec_ref`, `_expr`, `throw`, `throw_methoderror`, `current_scope`, `finalizer`, `donotdelete`, `compilerbarrier` |

### 3.2 Intrinsics: 93 (`intrinsics.h:3-110`)

The `julia.md` §3.1 mapping stands (`math/bits`, `math`, explicit widths up to 128 bits,
software `Float16`/`BFloat16`), minus `cpu_supports`. Conformance oracle: `test/intrinsics.jl`
(566 lines). The fleet rule "one implementation per kernel" applies: checked 128-bit arithmetic
and `Float16` conversion go in **go-simd** (or the lowest library that needs them) if they are
not there already, and the VM imports them.

### 3.3 Types and subtyping (`jltypes.c`, `datatype.c`, `subtype.c`)

`subtype.c:3-21` names its algorithm: "section 4.2.2 of" Bezanson's PhD thesis
(<https://github.com/JeffBezanson/phdthesis/>), plus concrete-only variables, the **diagonal
rule** ("a type variable is concrete if it occurs more than once in covariant position, and never
in invariant position … it depends on which occurrences of a type variable the algorithm
actually uses"), non-type values matching unconstrained variables, `Vararg{T,N}`, `Type{T}<:S if
isa(T,S)`, and tolerated free variables. Its state is `jl_stenv_t` (`:102-`): a linked list of
`jl_varbinding_t` (`:67-91`, bounds `lb`/`ub`, `occurs_inv`, `occurs_cov`, `concrete`,
`intvalued`, …) and two `jl_unionstate_t` bit stacks (`:43-53`) through which
`forall_exists_subtype` (`:1807`) enumerates the `Union` components of each side.

Design: **a statement-by-statement port of `subtype.c`** (as `julia.md` §10 already says), with
the bit-stack enumeration kept as is, because the order in which unions are explored decides
which occurrences set `occurs_cov`, and therefore the diagonal rule. Entry points to expose:
`jl_subtype_env` (`:2246`), `jl_type_intersection_env_s`, `jl_type_morespecific` (`:5516`).
The formal model (Zappa Nardelli et al., OOPSLA 2018, doi:10.1145/3276483) is a reading aid,
not the specification. Oracle: `test/subtype.jl` (2 788 lines) and `test/specificity.jl` (329),
plus a differential fuzz: random type pairs from a small grammar, `<:` and `typeintersect`
asked of the reference julia through the M0 oracle and of Go.

### 3.4 Dispatch (`gf.c`, `typemap.c`, `method.c`)

`jl_apply_generic` (`gf.c:4207-4215`) is: world age of the current task, `jl_lookup_generic_`,
then `_jl_invoke`. The lookup (`:4130-4205`) tries a 4-way associative call-site cache, then a
leaf cache keyed by the argument type tuple, then the method cache typemap
(`jl_typemap_assoc_exact`), and only on a miss `jl_mt_assoc_by_type`; no method is
`jl_method_error` (`:4184`).

Design: the same three levels, in Go, keyed by the concrete argument types — a per-call-site
inline cache in the prepared statement (§2.2), a leaf cache (`map` from an interned type-tuple
id), and the full typemap. World ages are part of every key, as in julia, because `latestworld`
(§2.1) makes a notebook cell that redefines a method visible to the next statement, and stale
caches would be a wrong answer, not a slow one.

### 3.5 Values and memory

- A `Value` is a Go interface; isbits values of primitive types are stored unboxed in typed
  fields (`int64`, `float64` …) of a small set of concrete Go types, and `Memory{T}` of an isbits
  `T` is a Go `[]byte` with the layout `datatype.c` computes, so that `sizeof`, `fieldoffset`,
  `reinterpret` and `pointerref` agree with julia.
- `Ptr{T}` is (buffer, offset); every `pointerref`/`pointerset`/`unsafe_load` is bounds-checked
  and raises, never reads outside the buffer (`julia.md` §10).
- Exceptions: a Go `panic` carrying a `*JuliaException` would cross `defer` frames unchecked;
  instead `eval_body` returns `(Value, *Exception)` and the `EnterNode` handler stack is a slice
  in the frame, which reproduces `setjmp`/`longjmp` (`interpreter.c:548-574`) without
  unwinding Go frames by panic.

## 4. Package layout (in go-embedded-julia/julia)

All public, per `ruby-reuse.md` §3 (go-scikit/julia must be able to register methods from its
own module).

| package | contents | ports |
|---|---|---|
| `core` | `Value`, `DataType`, `TypeName`, `Union`, `UnionAll`, `TypeVar`, `Vararg`, `Module`, `Binding`, `Method`, `MethodInstance`, `SimpleVector`, `Memory`, `Symbol` interning | `julia.h`, `jltypes.c`, `datatype.c`, `module.c`, `symbol.c`, `simplevector.c`, `genericmemory.c` |
| `core/subtype` | `Subtype`, `Intersect`, `MoreSpecific` | `subtype.c` |
| `core/dispatch` | method tables, caches, world ages, ambiguity | `gf.c`, `typemap.c`, `method.c` |
| `core/builtins` | the 63 builtins | `builtins.c` |
| `core/intrinsics` | the 93 intrinsics | `runtime_intrinsics.c` |
| `interp` | `prepare`, `Frame`, `evalBody`, `evalValue`, toplevel eval, the `foreigncall` shim table | `interpreter.c`, `toplevel.c` |
| `julia` (module root) | `VM`, `Eval(cell)`, `Call`, `Register`, Go ↔ Julia values (`julia.md` §5.2) | — |

## 5. M3 exit, measured as the fleet measures

| check | reference | pass condition |
|---|---|---|
| IR execution | `testdata/programs/*.jl` with julia's printed values (M0 already records "parse, `Meta.lower` and evaluation", go-embedded-julia README) | identical values, Core only |
| `boot.jl` (1 086 lines) | julia 1.12.6 | evaluates to a `Core` whose `names(Core; all=true)` and the `typeof`/`supertype` of each name equal julia's |
| subtyping | `test/subtype.jl`, `test/specificity.jl` + differential fuzz | per-`@test` records equal (`julia.md` §8, conformance) |
| intrinsics | `test/intrinsics.jl` | idem |
| `core.jl` | the parts not needing Base I/O | idem, with the excluded testsets listed and justified |
| speed | julia 1.12.6 with `--compile=min` on the same programs, one core, interleaved rounds, load below 8. Under `min`, a method with no compiled code that does not need the compiler runs in `jl_fptr_interpret_call` (`gf.c:3471-3482`), but a method that already has native code, such as Base's in the system image, keeps it (`gf.c:3441-3469`): the reference is julia's interpreter for user code over a compiled Base, a stricter bar than interpreter against interpreter | no slower on Core-only programs; a slower row is an open defect (`julia.md` §10 explains why the default JIT itself is not the bar) |
| security | — | `prepare` fuzzed on mutated `CodeInfo` (never a Go panic, always a julia exception); recursion bound (`StackOverflowError`); allocation bound in `memorynew` (`OutOfMemoryError`); a cell cannot reach a `foreigncall` outside the shim table |

## 6. The boundary with the rest of the fleet

### 6.1 go-bignum

`go-bignum` holds one library, `bigfloat` ("correctly rounded elementary functions and constants
on `math/big.Float`, any precision, all rounding modes"). It is the MPFR target of `julia.md`
§3.2 (M5, not M3). `BigInt` stays on `math/big.Int`; `fac_ui`/`bin_ui` have no fleet home yet.

### 6.2 go-ndarray: Fortran order is a view, but not without a copy

- `Array` is (data, shape, strides, offset) (`ndarray.go:17`, `:39`); `Transpose` "returns a view
  with the axes reversed. No data is copied" (`ndarray.go:455-466`). So a Julia `m×n` column-major
  buffer is exactly `Transpose` of an `n×m` row-major array: go-ndarray **can represent** it.
- But every public constructor copies: `FromData` ("The data slice is copied",
  `ndarray.go:133-150`) and `FromSlice` ("wraps a copy of data", `storage.go:131-143`). `Data`
  returns the array's own storage only when it is contiguous row-major (`storage.go:145-149`).
- So `julia.md` §5.2's "`*ndarray.Array` sharing the same buffer, no copy" is **not possible with
  today's API**. Missing: a constructor that adopts a buffer with given strides
  (`Wrap[T](data []T, shape, strides []int, offset int)`, validated against `len(data)`). It
  belongs in **go-ndarray** (the lowest shared library), not in the VM; it is an M6 dependency,
  to open as an issue there when M6 starts.

### 6.3 go-notebook

`analysis/ruby.go` analyses a Ruby cell into a `*Cell` of defined and used names
(`AnalyzeRubyCell`, `:26-43`); the Julia analysis (parse and names, being added now by another
session) follows the same shape. M3 adds nothing to the notebook; execution of ` ```julia `
cells starts when `julia.VM.Eval` exists, and the notebook's "not run yet (M3)" message names this
milestone.

## 7. Order of work

| step | contents | depends on | size |
|---|---|---|---|
| 1 | `core` values and type objects; `boot.jl`'s first forms (`struct`, `abstract type` through the type-construction builtins) | M2 lowering `boot.jl` | M |
| 2 | `interp.prepare` + `evalBody` for §2.1 without `foreigncall`; fuzz `prepare` | step 1 | M |
| 3 | `core/subtype`, statement-by-statement; differential fuzz | step 1 | L |
| 4 | `core/dispatch` with world ages | 3 | M |
| 5 | the 63 builtins, the 93 intrinsics | 1, 2 | M |
| 6 | all of `boot.jl`; the M3 exit table | 1-5 | M |

## 8. Unverified here

- Nothing of M3 was run: no speed figure exists, and §5's bars are the protocol, not results.
- `jl_type_intersection_env_s` and `typemap.c` were located, not read line by line.
- That M2's goldens contain only the heads of §2.1 was not checked against
  `testdata/golden/lowering.jsonl`; it is the first test step 2 should write.
- Whether go-simd already has checked 128-bit arithmetic and `Float16` conversion (§3.2).
