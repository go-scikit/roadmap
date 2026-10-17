# Ruby / go-scikit overlap census: what the pure-Go Ruby duplicates, and who should own it

Date: 2026-10-17. Read-only census. Everything below was read from GitHub `main` through the API
(trees, `contents/<path>`, `contents/go.mod`); no local clone was consulted.

Scope. Ruby side: `go-embedded-ruby/ruby` ("rbgo", `internal/vm`), and the numeric `go-ruby-*`
orgs: `bigdecimal`, `prime`, `cmath`, `rational`, `complex`, `matrix`, `arrow`, `parquet`,
`images`, `money`, `securerandom`. go-scikit side: `go-ndarray/ndarray` (+ `linalg`, `random`,
`fft`), `go-fft/fft`, `go-numerics/{special,optimize,integrate}`, `go-statistics/stats`,
`go-dataframe/dataframe`, `go-plot/plot`, `go-images/*`, and the planned orgs of `PLAN.md`.

## 1. What exists today (facts the table rests on)

**Which orgs exist on GitHub.** `go-bignum`, `go-sparse`, `go-approx`, `go-dsp`, `go-units`,
`go-sciio`, `go-symbolic` have **no repository yet** (org listing is empty); so every "planned"
counterpart below is a plan, not code. `go-scikit/ruby` exists but holds only `LICENSE` and a
one-line `README.md` ("Ruby bindings to the go-scikit stack, for the pure-Go Ruby").

**rbgo `internal/vm` numeric files** (sizes in lines): `math.go` 270, `cmath.go` 97,
`complex.go` 1089, `rational.go` 656, `integer_bignum.go` 114, `numeric_edges.go` 963,
`random.go` 493, `prime_bind.go` 161, `matrix.go` 703, `ndarray.go` 615, `fft.go` 220,
`image.go` 186, `images.go`/`images_bind.go`/`images_canvas.go`, `bigdecimal_bind.go`,
`securerandom.go`. There is no BigMath, no statistics, no dataframe, no plot, no sparse binding.

**rbgo `go.mod` imports from the numeric world:** `go-ndarray/ndarray v0.11.0`,
`go-fft/fft v0.23.0`, `go-images/images` (pseudo-version of 2026-09-27), `go-ruby-bigdecimal`,
`go-ruby-cmath`, `go-ruby-matrix`, `go-ruby-prime`, `go-ruby-arrow`, `go-ruby-parquet`,
`go-ruby-images`, `go-ruby-money`, `go-ruby-securerandom`. It does **not** import
`go-ruby-rational` nor `go-ruby-complex` (both repositories exist, 100% coverage, "byte-exact
against MRI 4.0.5"): rbgo carries its own `rational.go` and `complex.go`. It does **not** import
`go-numerics/*`, `go-statistics/stats`, `go-dataframe` or `go-plot`. (ndarray is at v0.11.0 in
rbgo while go-scikit modules require v0.12.0; fft v0.23.0 against v0.27.1 required by ndarray.)

**`go.mod` of the `go-ruby-*` numeric libraries:** `bigdecimal`, `prime`, `cmath`, `rational`,
`complex`, `matrix`, `money` have **no dependency at all** (stdlib `math`, `math/big` only).
`arrow` requires `apache/arrow-go/v18 v18.8.0`; `parquet` requires `arrow-go/v18 v18.8.0` and
`go-ruby-arrow/arrow`; `images` requires `go-images/images v0.1.0` and `golang.org/x/image`.
None of them imports go-ndarray, go-dataframe or go-bignum.

**go-scikit module dependencies (relevant edges):** `special`, `optimize`, `integrate` and
`stats` require `go-ndarray/ndarray v0.12.0`; `stats` also `optimize` + `special`;
`go-dataframe/dataframe` requires `arrow-go/v18 v18.8.0`, `go-ndarray v0.12.0`, `go-plot v0.1.1`
(and through the plot sub-package `dfplot`, `go-gfx`, `go-opentype`, `go-tex`, `go-widgets`);
`go-ndarray` requires `go-fft v0.27.1`. Note `ndarray/random` already contains numpy's
`mt19937.go`, `pcg64.go`, `philox.go`, a ziggurat, and `random/libm.go`, a correctly rounded
double-double `log/log1p/exp/expm1/pow`.

**What the Ruby code actually does for Math.** `Math.sin/cos/exp/cbrt/erf/erfc/gamma/lgamma/...`
call Go's `math` package directly (`math.Gamma`, `math.Lgamma`, `math.Erf`), with Ruby's domain
checks and bignum-aware `log` around it. Go's `math` is not the C libm MRI calls, so rbgo is
**already not libm-exact in the last ulp**; the oracle that makes rbgo's tests pass is whatever the
author's machine produced.

## 2. Overlap table

Legend: duplicated? yes = the same algorithm exists twice in our code; partly = overlapping scope or
a missed reuse; no = looked at, nothing to dedupe. Effort S (days), M (weeks), L (months).

| # | Ruby component | What it implements | go-scikit counterpart | Dup? | Recommended owner and direction | Constraint | Effort |
|---|---|---|---|---|---|---|---|
| 1 | rbgo `Math.gamma`, `Math.lgamma`, `Math.erf`, `Math.erfc` (`math.go`) | Go `math.Gamma/Lgamma/Erf/Erfc` plus Ruby edge cases (gamma(-0.0) = -Inf, lgamma sign) | `go-numerics/special`: `Gamma`, `GammaSign`, `LogGamma`, `LogGammaSign`, `Erf`, `Erfc`, `Erfcx` (about 0.5 ulp against mpmath, integers from a correctly rounded factorial table) | partly (a missed reuse, not copied code) | rbgo may call special for these four, behind the existing Ruby edge-case wrappers. Optional, and decided by an oracle run, not by taste. The rest of `Math` (`sin`, `cos`, `exp`, `sqrt`, `hypot`, `atan2`, ...) has nothing in go-scikit worth calling: keep stdlib `math` | MRI's results come from the platform C libm (glibc `tgamma` is several ulp off the true value in places; macOS libm differs from glibc). A near-correctly-rounded value from special can differ from MRI in the last ulp where Go's `math.Gamma` already also differs, in different places. MRI computes `gamma` of an integer 1..23 from an exact table: special's correctly rounded factorial table agrees there. Change only if a measured diff against real MRI on the CI platforms shows no new mismatches | S | 
| 2 | rbgo `Math.sin/cos/tan/exp/log/log2/log10/sqrt/cbrt/hypot/atan2/...` | stdlib `math` | `ndarray/random/libm.go` (correctly rounded `log/log1p/exp/expm1/pow` in double-double); `go-bignum` plan item 1 (correctly rounded `big.Float` elementary functions) | partly, inside go-scikit: `libm.go` and go-bignum will compute the same correctly rounded functions twice | go-bignum should own correct rounding. When go-bignum ships, `random/libm.go` should become a thin wrapper or be deleted (the double-double fast path may stay in ndarray if it is the faster one, but go-bignum is its test oracle). Ruby: keep stdlib `math` | Ruby must stay MRI-exact where MRI is libm: a correctly rounded function is not what glibc/Apple libm return in 0.06% to 0.3% of arguments (figures from the `libm.go` header), and Go's `math` uses assembly on some architectures, so Ruby results already vary per architecture. A shared CR libm would give Ruby cross-architecture determinism at the price of differing from MRI in rare ulps. Decide explicitly; default: leave | M (scikit-internal) / none (Ruby) |
| 3 | rbgo `Random` (`random.go`): `RandomObj` with its own MT19937, `init_genrand`, `init_by_array` over the seed's 32-bit words, 53-bit float, byte stream for `securerandom` formatting | Mersenne Twister seeded as MRI does | `ndarray/random/mt19937.go`: `MT19937`, `NewMT19937Legacy`, `SeedLegacy`, `SeedLegacyArray` (the reference `init_by_array`), `Uint32`, `Uint64`, `Float64` ((a>>5, b>>6) 53-bit form) | **yes** | The MT19937 core (state, twist, tempering, init_genrand, init_by_array, 32-bit output, 53-bit float) belongs to `ndarray/random`, which has the numpy-bit-exact, jump-ahead (`zz_mtjump.go`) version. rbgo keeps the MRI layer: seed to key words, `rand(n)` via MRI's `limited_rand` rejection mask, `shuffle`/`sample`/`Random::Formatter`. Propose that rbgo import `ndarray/random` MT19937 for the core and keep its own tests as the check | Seeded output must match MRI bit for bit. MRI and numpy agree on the core generator (same Knuth init, same `init_by_array`, same 53-bit float), but **not** on bounded integers, shuffle or sampling: only the core is shareable. A import of `ndarray/random` pulls the ndarray module into the dependency (rbgo already imports it) | S |
| 4 | rbgo `rational.go` (656 lines) + the Rational parts of `numeric_edges.go` | Ruby `Rational` on `math/big.Rat` | none in go-scikit (stdlib `big.Rat`); the Ruby-side twin is `go-ruby-rational/rational` (arith, compare, convert, fastint, 100% coverage, MRI 4.0.5 byte-exact) | **yes, inside Ruby** | Not a go-scikit matter, but go-scikit/ruby must not build on a third copy: the rbgo and go-ruby-rational sessions should settle which one is canonical (go-ruby-rational by the org's own design) and make rbgo import it | The in-tree one may have VM-coupled behaviour (coercion, `bytecode`) that the library does not model; unverified here. Do not touch from the scikit side | M |
| 5 | rbgo `complex.go` (1089 lines), `cmath.go` binding; `go-ruby-complex`; `go-ruby-cmath` | Ruby `Complex`; CMath functions as a line-for-line port of the cmath gem (real branch stays real, otherwise complex) | go-scikit has complex128 dtypes (`ndarray`, `fft`) but no Ruby-visible complex math library; go-bignum Wave 4 plans complex ball arithmetic | complex vs scikit: no. Complex vs go-ruby-complex: **yes, inside Ruby** (rbgo does not import go-ruby-complex; it does import go-ruby-cmath) | Ruby owners as in row 4. go-scikit side: nothing to absorb. When go-bignum gets complex functions, they are a new capability, not a replacement for CMath | `go-ruby-cmath` is deliberately a port of Ruby source (Float results must agree with `ruby -rcmath`); do not replace its formulas with `math/cmplx` or go-bignum | S (Ruby-internal) |
| 6 | rbgo Integer/Bignum (`integer_bignum.go`, `numeric_edges.go`): `gcd/lcm/gcdlcm`, `bit_length`, `fdiv`, floor/ceil/round with digits, `Integer.sqrt`, `pow(a, mod)`, `digits`, `rationalize` | `math/big` plus Ruby rounding rules | go-bignum number theory (plan): gcd/xgcd, modular inverse, powmod, integer roots, binomial, factorial, partitions, Bernoulli (exact big-integer versions) | partly (the first four are stdlib `big.Int` calls already; the rest not inspected) | go-bignum should expose exact integer primitives with plain `*big.Int` signatures (`ISqrt`, `PowMod`, `ModInverse`, `Binomial`, `Factorial`, extended GCD) so rbgo can call them instead of keeping local copies. Rounding-with-digits stays in rbgo (it is MRI semantics, not mathematics) | Exact integers are deterministic: sharing costs nothing in MRI fidelity. Only error classes and argument coercion are Ruby's. I did not read every rbgo Integer method, so "which local copies exist" must be confirmed by that session | S |
| 7 | `go-ruby-prime/prime` (24.5 KB) + rbgo `prime_bind.go` | segmented sieve with a growing cache, deterministic Miller-Rabin below 2^64 with Montgomery arithmetic, Baillie-PSW (strong base 2 + strong Lucas, Jacobi), Pollard rho, `PrimeDivision` (with ZeroError), `Each`, `Take`, `Next`, `Prev` | go-bignum number theory (plan, bib `precision-numbertheory-symbolic.md`): `isprime` BPSW, `NextPrime`/`PrevPrime`, segmented sieve, `PrimePi`, factorisation (rho, ECM, quadratic sieve) | **yes** (by plan: go-bignum would write the same first four) | go-bignum should **start from** go-ruby-prime's primality, sieve and rho code (same licence family, same fleet) as `go-bignum/numtheory`, then extend (ECM, MPQS, `PrimePi`, ECPP). go-ruby-prime then shrinks to the Ruby API (enumerator semantics, `prime_division` ordering and the `[-1,1]` convention for negatives, error classes) over go-bignum. Direction is "absorb", not "rewrite" | Results are mathematical facts, so sharing is exact. Keep the Ruby wrapper's behaviour byte-identical (its tests are the guard). go-bignum must not regress the Montgomery fast path for < 2^64 (that is where Ruby loops spend time) | M |
| 8 | `go-ruby-matrix/matrix` (Matrix, Vector; `Num` = Integer / Rational / Float) + rbgo `matrix.go` | exact-arithmetic Ruby `Matrix`: `determinant`, `inverse` (exact Rationals), `rank`, `pow`, `trace`, minors, `Vector` inner/cross/normalize/angle | `ndarray/linalg` (`Solve`, `Inv`, `Det`, `Slogdet`, LU, QR, Cholesky, SVD, Eig, Eigh, lstsq, norm, power) in float64 / complex | partly (same operations, different number domain) | Do **not** merge. Exact Integer/Rational Matrix is a different algorithm family (fraction-free elimination) from LAPACK-style float LU. go-ruby-matrix stays the owner for Ruby `Matrix`/`Vector`. The float-only methods MRI lacks (SVD, QR, Cholesky, least squares) belong in the future go-scikit/ruby `NDArray`/`Linalg`, not in `Matrix`. Optional: a `Matrix#to_ndarray`/`NDArray#to_matrix` bridge | MRI-exact: `Matrix#determinant` and `#inverse` on Float entries follow MRI's own elimination order, so their last bits differ from LU with partial pivoting. A call into `linalg.Det` would change Ruby outputs. Also `eigen`/`lup` in MRI are ports of JAMA/own code with their own rounding | none (bridge S) |
| 9 | rbgo `ndarray.go` (615 lines, hand-written binding of `go-ndarray/ndarray`) vs the planned `go-scikit/ruby` | float64 NumPy-style arrays in Ruby: creation, views, ufuncs, reductions, `matmul`/`dot` | `go-ndarray/ndarray` (all dtypes, `linalg`, `random`, `fft`) | partly: a second binding of the same library would be written if go-scikit/ruby starts from zero | One binding, one owner. Recommended: go-scikit/ruby owns the binding code (all of ndarray including `linalg`, `random`, complex dtypes, and the new libraries), and rbgo drops `ndarray.go`/`fft.go` in its favour. This needs an **exported extension API in rbgo** (see section 3): `internal/vm` cannot be imported from another module | Go's `internal` rule: the binding cannot live outside rbgo today. The current Ruby surface was generated by the go-ndarray playground code generator (comment in `ndarray.go`): keep that surface as the compatibility baseline | L |
| 10 | rbgo `fft.go` (220 lines) | binds `go-fft/fft`: FFT, IFFT, RFFT, IRFFT, FFTFreq, FFTN/IFFTN, FFT2, windows (hann, hamming, blackman, blackman_harris, bartlett), PSD, Spectrogram | `go-fft/fft` itself has `window.go`, PSD and Spectrogram; planned `go-dsp` lists Welch, STFT, windows, resampling | Ruby vs scikit: no (direct reuse, already best case). Inside go-scikit: partly (go-dsp would duplicate go-fft windows/PSD/spectrogram) | go-fft keeps FFT, DCT/DST, windows and the simple PSD/spectrogram it has; go-dsp imports them and must not re-implement windows, adding the SciPy-compatible `get_window`, Welch, STFT, filters on top. rbgo/go-scikit/ruby bind whichever module the function lives in | None for Ruby (no MRI counterpart; the surface is a go-fft one). Bump rbgo's `go-fft v0.23.0` pin at the same time as the ndarray bump | S |
| 11 | rbgo `images.go`/`images_bind.go`/`image.go` vs `go-ruby-images/images` | Ruby image API; `go-ruby-images` wraps `go-images/images` v0.1.0 (canvas, color, format, process, transform); rbgo imports both `go-images/images` directly and `go-ruby-images`, and `image.go` is a stdlib `image`/`color` binding | `go-images/*` (8-bit RGBA today; the PLAN adds float/N-D images, morphology, labelling, registration) | partly: three binding layers over one library | `go-images` owns the algorithms (incl. the planned float image type). `go-ruby-images` owns the Ruby-gem-shaped API (it is the red-image/ruby-vips style facade). rbgo should keep only registration glue. When the float image type lands, expose it through go-scikit/ruby as an NDArray-backed image rather than a fourth binding | Ruby has no MRI numeric oracle for images; compatibility target is the gem API shape. The float image type in bib `approximation-signal-image.md` is "the pivotal design decision" and will change the Ruby binding | M |
| 12 | `go-ruby-bigdecimal/bigdecimal` + rbgo `bigdecimal_bind.go`; BigMath | `Decimal` = sign x big.Int coefficient x 10^exp, NaN/Inf, `Add/Sub/Mul/Div(prec)/IDiv/Mod/Pow(int)`, rounding modes, `Floor/Ceil/Truncate`, `FromFloat`, MRI-byte-exact `to_s` grammar. The tree shows **no `sqrt`, no `exp/log/sin/cos/atan`, no `BigMath`** | go-bignum: correctly rounded **binary** `big.Float` elementary functions and constants (item 1), later ball arithmetic and N-digit special functions | partly: the algorithms (argument reduction, Taylor/AGM, Newton for roots, pi by Machin/Chudnovsky) coincide, the number representation cannot | Share **algorithms and tests, not types**. (a) go-bignum is the **test oracle** for BigMath: `BigMath.PI(n)`, `E`, `sqrt`, `exp`, `log`, `sin`, `cos`, `atan` leading digits checked against correctly rounded values at n+guard digits. (b) If go-bignum writes its series kernels over `*big.Int` fixed point with the scale as a parameter (2^k for `big.Float`, 10^k for decimal), go-ruby-bigdecimal can reuse the kernels with a decimal scale. Ask go-bignum to expose that kernel only once its API is stable. (c) `BigDecimal#sqrt` needs only `big.Int.Sqrt` on a scaled coefficient: no go-bignum needed | BigDecimal is **decimal**: `0.1` is exact; a `big.Float` of any precision is not. A binary value cannot become the Ruby result without a decimal rounding step that MRI does not perform. MRI's BigMath (Ruby code in `bigdecimal/math.rb`, with `BigDecimal.double_fig` guard digits) returns specific trailing digits, which a correctly rounded result will not reproduce; byte-exactness must come from porting MRI's own algorithm, using go-bignum only to verify the leading digits | M (oracle: S; shared kernels: M) |
| 13 | `go-ruby-arrow/arrow` (array, datatype, ipc, record_batch, table) and `go-ruby-parquet/parquet` (reader, writer, compression, options) vs `go-dataframe/dataframe` `arrowio/` and `parquetio/` | Arrow IPC and Parquet read/write over `apache/arrow-go/v18 v18.8.0` (the same version in all three) | `go-dataframe/dataframe/arrowio` (IPC, pandas conventions), `parquetio` (with a read plan); `go-sciio` (planned) does **not** list Arrow or Parquet | **yes** at the Go level (two wrapper layers over one library) | `go-dataframe` owns the table-level logic (type mapping, compression options, pandas interchange). The Ruby gems (red-arrow, red-parquet shapes) own only their API. Preferred direction: go-ruby-arrow/parquet delegate type conversion to a lean module. **Caution**: `go-dataframe`'s `go.mod` also requires `go-plot`, `go-gfx`, `go-opentype`, `go-tex`, `go-widgets` (for `dfplot`), so importing it would bloat the Ruby dependency graph; first ask the dataframe owner to split `arrowio`/`parquetio` into their own module (or remove the plot edge) | Red-Arrow's type names, error classes and compression names are Ruby-API facts. Behaviour of Parquet round trips is arrow-go's in both, so identical by construction | M |
| 14 | `Array#sum`/`Float#sum` Kahan-Babuska summation, `Enumerable` min/max/minmax, `Integer#digits`, `Numeric#step` (not inspected) | MRI's compensated summation | `go-statistics/stats` (`accurate.go`, moments, `Describe`), `ndarray` pairwise `Sum` | not verified | Not recommended to share: MRI specifies the exact compensated algorithm; a different accurate sum changes last-bit results. `stats` is the owner of statistics APIs in go-scikit/ruby, where no MRI oracle constrains it | MRI-exact summation order | none |

No overlap found (looked at, nothing to dedupe): `go-ruby-money` (currency amounts, no physical units, so
`go-units` is unrelated), `go-ruby-securerandom` (cryptographic bytes and formatting), Ruby `Float#to_s`
and `round`, `go-ruby-set`/`ostruct`-style non-numeric libraries. There is no Ruby counterpart for
`go-sparse`, `go-approx`, `go-dsp`, `go-symbolic`, `go-units`, `go-sciio`, `go-plot`, `go-dataframe`,
`go-statistics`, `go-numerics/optimize|integrate`: those are pure go-scikit capabilities that the
planned `go-scikit/ruby` repository would expose, with nothing in Ruby to merge.

Count: 13 overlaps with substance (rows 1 to 13). "Yes" rows: 3, 4, 5 (Ruby-internal part), 7, 13;
the others are "partly". Row 14 is unverified.

## 3. Recommendations

**Order to deduplicate (cheapest and safest first).**

1. **Row 3, MT19937 (S).** `ndarray/random` already has the numpy-exact core. Propose to the rbgo
   session to import it for the generator and keep MRI's bounded-integer/shuffle layer. Check by
   keeping the existing seeded-output tests unchanged.
2. **Rows 4 and 5, Rational and Complex (Ruby-internal, M).** Before go-scikit/ruby exists, the Ruby
   sessions should choose `go-ruby-rational` and `go-ruby-complex` as canonical and make rbgo import
   them. Otherwise the binding repo ends up with a third copy of numeric tower code.
3. **Row 7, Prime (M), together with go-bignum creation.** Seed `go-bignum/numtheory` from
   `go-ruby-prime` (primality, sieve, rho) rather than writing it from the bib, then thin
   `go-ruby-prime` down to the Ruby API. This is the single largest piece of code that would
   otherwise be written twice.
4. **Row 12, BigMath (S for the oracle, M for kernels).** Use go-bignum as the validation oracle for
   BigMath digits first; share fixed-point kernels only if go-bignum's design allows a radix
   parameter.
5. **Row 9 and 10, binding consolidation (L).** Needs the exported rbgo extension API first (below).
   After it, move `ndarray.go`/`fft.go` into go-scikit/ruby, and let go-dsp import go-fft windows.
6. **Row 13, Arrow/Parquet (M).** Needs the dataframe owner to separate the plot dependency from
   `arrowio`/`parquetio`.
7. **Rows 1, 2 (optional).** Only after a measured comparison against real MRI on the CI platforms.

**go-scikit APIs that must exist or stabilise first.**

- `go-bignum` (nothing exists yet): exact-integer primitives with `*big.Int` signatures (`ISqrt`,
  `PowMod`, `ModInverse`, extended GCD, `Binomial`, `Factorial`), the number-theory package seeded from
  go-ruby-prime (`IsPrime`, `NextPrime`, `PrevPrime`, `Sieve`/`Each`, `Factor` returning ordered
  prime/exponent pairs), `big.Float` constants and `exp/log/sin/cos/atan/sqrt` correctly rounded at a
  requested number of **bits**, and a documented way to ask for **decimal** digits (so BigMath checks
  do not need hand conversion). Keep all of them free of Ruby types.
- `go-ndarray/ndarray/random`: a stable exported MT19937 surface (`NewMT19937Legacy`,
  `SeedLegacyArray`, `Uint32`, `Float64`, state get/set) that rbgo can depend on without pulling
  numpy-only parts; confirm it does not use `go-fft`-heavy imports that bloat the VM.
- `go-ndarray`: a v0.12+/v1 stability statement, because both rbgo (v0.11.0) and the planned binding
  repo would pin it; complex dtype results must be returnable as Ruby `Complex`.
- `go-fft`: keep `window.go` and PSD public and stable; go-dsp depends on them.
- `go-numerics/special`: nothing to add for Ruby; if row 1 proceeds, `Gamma` of integers 1..23 must
  stay exactly the factorial (it is, by its documented table).
- `go-dataframe`: split `arrowio`/`parquetio` so a Ruby module can import them without `go-plot`.

**Changes that belong to the sessions owning rbgo / go-ruby-*** (to be proposed, not done here):

- rbgo: import `go-ruby-rational` and `go-ruby-complex` instead of the in-tree copies (rows 4, 5);
  replace its hand-written MT19937 with `ndarray/random` (row 3); expose a **public extension API**
  (a registration package outside `internal/`) so that `go-scikit/ruby` can register `NDArray`, `FFT`,
  `Stats`, ... without importing `internal/vm`. Without it row 9 cannot start. Optionally bump the
  `go-ndarray`/`go-fft` pins.
- go-ruby-prime: after go-bignum numtheory exists, depend on it and keep the Ruby API only (row 7).
- go-ruby-bigdecimal: add `sqrt` and BigMath as ports of MRI's Ruby algorithms, using go-bignum only as
  an oracle (row 12). Its own session decides whether to adopt shared kernels later.
- go-ruby-matrix: no change; consider a conversion bridge to NDArray once the binding exists (row 8).
- go-ruby-arrow / go-ruby-parquet: consider delegating type mapping to go-dataframe's `arrowio`/
  `parquetio` after the module split (row 13).
- go-ruby-images / rbgo: agree that rbgo keeps registration glue only (row 11).

**What not to do.** Do not route `Matrix#determinant`, `Array#sum`, BigDecimal rounding or `CMath`
through a go-scikit routine because it is "more accurate": each is defined by an MRI algorithm, and a
better answer is a different answer.

## 4. Limits of this census

- Read from GitHub `main` on 2026-10-17; function lists were taken from `^func` declarations and
  `go.mod`, not from running anything. Unmeasured: how many Ruby-observable results change if rows 1
  or 2 are applied.
- rbgo's Integer/Rational methods in `numeric_edges.go` (963 lines) were skimmed by name, not read in
  full; row 6's "local copies" and row 14 need confirmation by the owning session.
- The planned orgs have no code; their rows rest on `PLAN.md` and `bib/precision-numbertheory-
  symbolic.md` (which also uses the older working names `go-numtheory` / `go-computeralgebra`,
  folded into `go-bignum` / `go-symbolic` by `PLAN.md`).
