# Arbitrary precision, number theory, computer algebra — bibliography

One file per domain. Facts only, each with a source (URL or DOI). Terse: tables, not prose.
Items still marked (verify) were not confirmed against the source on 2026-10-17.

## 1. Feature inventory

What MATLAB (base + Symbolic Math Toolbox) and Mathematica offer, by function name, grouped by capability.
"go-scikit today" = what exists on GitHub main (checked via API 2026-10-17) plus Go's stdlib.
Docs: Mathematica https://reference.wolfram.com/language/ ; MATLAB https://www.mathworks.com/help/symbolic/ .

### (a) Arbitrary-precision numerics

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| Multiprecision float, N digits | `vpa`, `digits` | `N[expr, n]`, `SetPrecision`, `Precision`, `Accuracy` | stdlib `math/big.Float` (prec bits, rounding modes ToNearestEven/ToZero/AwayFromZero/ToNegativeInf/ToPositiveInf); `Sqrt` only | no exp/log/trig/pow/constants in stdlib |
| Correctly rounded elementary functions | none (`vpa` guard digits, no guarantee) | `N` raises working precision until the result is reliable (significance arithmetic) | `ALTree/bigfloat` (MIT): Exp, Log, Pow, Sqrt on `big.Float`; no guarantee of correct rounding stated | MPFR-style correct rounding of exp, log, sin, cos, atan, pow, hyperbolics |
| Constants to N digits | `vpa(pi,n)`, `vpa(exp(sym(1)))` | `N[Pi,n]`, `N[EulerGamma,n]`, `N[Catalan,n]`, `N[Glaisher,n]` | none | pi (Chudnovsky), e, ln2, gamma, Catalan, zeta(3) |
| Ball / interval arithmetic | none in base; Intlab is third party | `Interval`, `Around` (v12.0+) | none | Arb-style midpoint-radius real/complex balls |
| Special functions to N digits | `vpa` of `gamma`, `besselj`, `hypergeom`, `zeta`, ... | `Gamma`, `BesselJ`, `HypergeometricPFQ`, `Zeta`, `EllipticK`, ... with arbitrary precision | `go-numerics/special`: float64 API, `internal/dd` double-double internally | hypergeometric core (0F1,1F1,2F1,pFq), gamma/lgamma, zeta, Bessel, erf, Airy, elliptic at N digits |
| Multiprecision linear algebra | `vpa` matrices, `det`, `inv`, `eig` on sym/vpa | `LinearSolve`, `Inverse`, `Eigenvalues` on arbitrary-precision matrices | `go-ndarray/ndarray` (dtypes float64, float32, complex128, int*, bool; no big types) | LU/QR/eig on `big.Float`; ball matrices (Arb `arb_mat`) |
| Multiprecision FFT / fast mult | internal | internal (GMP) | `remyoudompheng/bigfft` (BSD-3): FFT multiplication for `big.Int`, last push 2023-01-29 | not wired into stdlib; Go's `big.Int` stays Karatsuba |
| Root finding / integration to N digits | `vpasolve` | `NIntegrate` with `WorkingPrecision`, `FindRoot`, `NSolve` | `go-numerics/integrate`, `optimize` (double) | tanh-sinh / Gauss-Legendre / Petras rigorous quadrature at N digits |

### (b) Number theory and combinatorics

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| Primality | `isprime` (to 2^64 for ints), `sym` variant | `PrimeQ`, `ProvablePrimeQ` (Wolfram Function Repository, not built in) | stdlib `big.Int.ProbablyPrime(n)`: Miller-Rabin + Baillie-PSW | proof certificates (ECPP), `NextPrime`, `PrevPrime` |
| Factorisation | `factor` (ints, to 2^64 and sym) | `FactorInteger` (methods undocumented) | none | Pollard rho/p-1, ECM, SIQS, trial division, perfect-power |
| Primes, sieve | `primes`, `nextprime`, `primes(n)` | `Prime`, `PrimePi`, `PrimeQ`, `Prime[n]` | none | segmented sieve, `PrimePi` (Lucy / Meissel-Lehmer) |
| Modular arithmetic | `mod`, `gcd`, `lcm`, `powermod` | `PowerMod`, `ModularInverse`, `ChineseRemainder`, `JacobiSymbol`, `MultiplicativeOrder`, `PrimitiveRoot` | stdlib `big.Int`: Exp(mod), ModInverse, ModSqrt, Jacobi, GCD (extended) | CRT, discrete log, primitive roots, Tonelli-Shanks general, Euler phi |
| Arithmetic functions | `eulerPhi`, `divisors` | `EulerPhi`, `DivisorSigma`, `MoebiusMu`, `Divisors` | none | all of them (from factorisation) |
| Continued fractions | `rat`, `cumprod` idioms | `ContinuedFraction`, `FromContinuedFraction`, `Convergents` | none | CF expansion of rationals/quadratic irrationals/reals at N digits; Stern-Brocot |
| Partitions, combinatorics | `nchoosek`, `factorial`, `perms`, `primes` | `PartitionsP`, `PartitionsQ`, `Binomial`, `Multinomial`, `StirlingS1/S2`, `BellB`, `CatalanNumber`, `Fibonacci`, `LucasL`, `BernoulliB`, `EulerE`, `HarmonicNumber` | `go-numerics/special` `comb.go` (double); stdlib `big.Int.Binomial` | exact big-integer versions; Hardy-Ramanujan-Rademacher partitions; Bernoulli (fast) |
| Diophantine, lattices | none | `FindInstance`, `Reduce` over integers, `LatticeReduce`, `FindIntegerNullVector`, `FindLinearRecurrence` | none | LLL, PSLQ / integer-relation (`identify`-style) |

### (c) Computer algebra

| capability | MATLAB Symbolic Math Toolbox | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| Expression trees, canonical forms | `sym`, `syms`, `symfun`, `children`, `subs`, `symvar` | `FullForm`, `Head`, `Replace`, `ReplaceAll`, `Orderless/Flat` attributes | none (`go-embedded-ruby/ruby` has Bignum/Rational/Complex only) | immutable hash-consed expression DAG, canonical ordering |
| Simplification | `simplify`, `expand`, `collect`, `factor`, `rewrite`, `combine`, `simplifyFraction` | `Simplify`, `FullSimplify`, `Expand`, `Together`, `Apart`, `Cancel`, `FunctionExpand`, `TrigReduce` | none | rational-function normal forms first; heuristic `Simplify` last |
| Univariate / multivariate polynomials | `poly2sym`, `coeffs`, `polynomialReduce`, `gcd`, `lcm`, `resultant`, `polynomialDegree` | `PolynomialGCD`, `PolynomialQuotientRemainder`, `Resultant`, `Discriminant`, `CoefficientList`, `Factor`, `FactorList`, `Cyclotomic` | none | dense/sparse polys over Z, Q, Z/p, GF(p^k); GCD, factorisation |
| Factorisation over Z/Q/F_p | `factor` | `Factor`, `FactorSquareFree`, `Factor[..., Modulus->p]`, `Extension->` | none | Zassenhaus + van Hoeij LLL; Cantor-Zassenhaus; multivariate via Hensel |
| Gröbner bases | no `groebner`; `polynomialReduce` (multivariate division, monomial orders), `resultant` | `GroebnerBasis`, `PolynomialReduce`, `Eliminate`, `Solve` of polynomial systems | none | Buchberger, F4, F5 |
| Rational functions | `partfrac`, `simplifyFraction`, `numden` | `Together`, `Apart`, `Cancel`, `Numerator`, `Denominator` | none | with poly GCD |
| Differentiation | `diff`, `jacobian`, `hessian`, `gradient`, `divergence`, `curl`, `laplacian` | `D`, `Dt`, `Derivative` | none (AD may live elsewhere) | rule table for ~60 functions |
| Indefinite / definite integration | `int` | `Integrate` | none | Risch-Norman, rational integration (Hermite/Rothstein-Trager), table/rule-based (Rubi) |
| Series, limits | `taylor`, `series`, `limit`, `symsum`, `symprod` | `Series`, `Limit`, `Sum`, `Product`, `AsymptoticSolve` | none | Gruntz, power series arithmetic |
| Equation solving | `solve`, `vpasolve`, `dsolve`, `linsolve` | `Solve`, `Reduce`, `Roots`, `NSolve`, `DSolve`, `Eliminate`, `FindInstance` | none | polynomial systems (Gröbner/triangular/resultants), inequalities (CAD, later) |
| Exact linear algebra | `det`, `inv`, `rref`, `eig`, `jordan`, `svd`, `null`, `colspace`, `charpoly`, `smithForm`, `hermiteForm` on `sym` | `Det`, `Inverse`, `RowReduce`, `NullSpace`, `CharacteristicPolynomial`, `JordanDecomposition`, `SmithDecomposition`, `HermiteDecomposition` | none (`go-ndarray` is floating point) | fraction-free (Bareiss), modular + CRT / p-adic lifting (Dixon), LLL |
| Special-function simplification | `simplify` with `Steps`, `rewrite(...,'gamma')` | `FunctionExpand`, `FullSimplify`, `Simplify` + `Assumptions` | none | identity tables, assumptions system |
| Pattern matching / rewriting | `isAlways`, `assume`, `assumeAlso`, `children` ; no general rules | `Rule`, `RuleDelayed`, `Replace`, `Condition`, `Blank` patterns, `PatternTest` | none | an AC-aware matcher; needed by rule-based integration |
| Output | `latex`, `mathml`, `pretty`, `ccode`, `fortran`, `matlabFunction` | `TeXForm`, `MathMLForm`, `CForm`, `FortranForm`, `Compile` | `go-tex` for TeX typesetting (engine, math layout) | expr -> LaTeX/MathML emitter; render via go-tex; `Lambdify`-style compile to Go closure |

## 2. Reference algorithms

### (a) Arbitrary precision

| capability | algorithm | reference | why |
|---|---|---|---|
| Correct rounding | MPFR model: every function returns the correctly rounded result for the given precision and mode; internal Ziv loop (raise working precision until rounding is decided) | Fousse, Hanrot, Lefèvre, Pélissier, Zimmermann, "MPFR: A multiple-precision binary floating-point library with correct rounding", ACM TOMS 33(2) 2007, doi:10.1145/1236463.1236468 ; Ziv, "Fast evaluation of elementary mathematical functions with correctly rounded last bit", ACM TOMS 17(3) 1991, doi:10.1145/114697.116813 | strongest contract: results are reproducible, testable bit for bit |
| Ball arithmetic | midpoint-radius real/complex balls, rigorous error bounds propagated through every operation | Johansson, "Arb: efficient arbitrary-precision midpoint-radius interval arithmetic", IEEE Trans. Comput. 66(8) 2017, doi:10.1109/TC.2017.2690633, arXiv:1611.02831 | certified results; ball width is the accuracy measure, so tests need no external oracle |
| Hypergeometric / special functions | rigorous hypergeometric series with bounds, asymptotic series, argument reduction; gamma, zeta, Bessel, erf, Airy all derive from it | Johansson, "Computing hypergeometric functions rigorously", ACM TOMS 45(3) 2019, doi:10.1145/3328732 | one core serves ~all special functions |
| Fast integer multiplication | Karatsuba, Toom-Cook, Schönhage-Strassen / FFT; asymptotically O(n log n) | Brent & Zimmermann, "Modern Computer Arithmetic", CUP 2010, https://members.loria.fr/PZimmermann/mca/pub226.html (free PDF); Harvey & van der Hoeven, Annals of Math 193(2) 2021, doi:10.4007/annals.2021.193.2.4 | constants (pi, e) and elementary functions are bound by multiplication speed; `bigfft` already provides FFT mult for Go |
| Constants | pi: Chudnovsky with binary splitting; e, ln 2, Euler gamma (Brent-McMillan) | MCA ch. 4 and 5 (same URL) | binary splitting makes the cost O(M(n) log^2 n) |
| Elementary functions | argument reduction + Taylor with binary splitting (small precision); AGM for log/exp at very high precision | MCA ch. 4 | switch point by precision |
| Linear algebra | LU/QR in `big.Float` with extra guard bits; ball matrices (Arb `arb_mat`) for verified results | Arb paper above | |

### (b) Number theory

| capability | algorithm | reference | why |
|---|---|---|---|
| Probable prime | Baillie-PSW (strong Fermat base 2 + strong Lucas) ; no counterexample known | Baillie & Wagstaff, Math. Comp. 35 (1980) 1391-1417, doi:10.1090/S0025-5718-1980-0583518-6 ; Baillie, Fiori, Wagstaff, "Strengthening the Baillie-PSW primality test", Math. Comp. 90 (2021), doi:10.1090/mcom/3616, arXiv:2006.14425 | what Go `ProbablyPrime` already does; deterministic below 2^64 with fixed bases |
| Proof of primality | ECPP | Atkin & Morain, "Elliptic curves and primality proving", Math. Comp. 61 (1993) 29-68, doi:10.1090/S0025-5718-1993-1199989-X | certificates verifiable in polynomial time; practical to ~10^4 digits |
| Deterministic poly-time | AKS | Agrawal, Kayal, Saxena, Annals of Math 160 (2004) 781-793, doi:10.4007/annals.2004.160.781 | theoretical only, far slower than ECPP |
| Small factors | Pollard rho, Brent's cycle detection | Pollard, BIT 15 (1975) 331-334, doi:10.1007/BF01933667 ; Brent, BIT 20 (1980) 176-184, doi:10.1007/BF01933190 | O(p^{1/2}), cheap |
| Medium factors | ECM | Lenstra, "Factoring integers with elliptic curves", Annals of Math 126 (1987) 649-673, doi:10.2307/1971363 | cost depends on the smallest prime factor, not n |
| Large semiprimes | Quadratic sieve (multiple polynomial, self-initialising) | Pomerance, "The quadratic sieve factoring algorithm", EUROCRYPT 1984, doi:10.1007/3-540-39757-4_17 | best to ~100 digits; NFS beyond, out of scope |
| Partitions p(n) | Hardy-Ramanujan-Rademacher series, ball arithmetic bounds | Johansson, "Efficient implementation of the Hardy-Ramanujan-Rademacher formula", LMS J. Comput. Math. 15 (2012) 341-359, doi:10.1112/S1461157012001088 | p(10^20) in seconds |
| Integer relations / lattices | LLL, PSLQ | Lenstra, Lenstra, Lovász, Math. Ann. 261 (1982) 515-534, doi:10.1007/BF01457454 | also needed inside polynomial factorisation |

### (c) Computer algebra

Textbooks: von zur Gathen & Gerhard, "Modern Computer Algebra", 3rd ed., CUP 2013, doi:10.1017/CBO9781139856065 ; Geddes, Czapor, Labahn, "Algorithms for Computer Algebra", Kluwer 1992, doi:10.1007/b102438 ; Bronstein, "Symbolic Integration I: Transcendental Functions", 2nd ed., Springer 2005, doi:10.1007/b138171.

| capability | algorithm | reference | why |
|---|---|---|---|
| Polynomial GCD | modular GCD (Brown) + CRT; heuristic GCD (GCDHEU) for small cases; sparse Zippel for multivariate | Geddes et al. ch. 7; von zur Gathen & Gerhard ch. 6 and 11 | avoids coefficient explosion of Euclid |
| Factoring over F_p | square-free, distinct-degree, equal-degree (Cantor-Zassenhaus) | Cantor & Zassenhaus, Math. Comp. 36 (1981) 587-592, doi:10.1090/S0025-5718-1981-0606517-5 | standard |
| Factoring over Z | factor mod p, Hensel lifting, recombination by LLL | van Hoeij, J. Number Theory 95 (2002) 167-189, doi:10.1006/jnth.2001.2763 | polynomial-time recombination, avoids Zassenhaus exponential case |
| Gröbner bases | Buchberger (criteria, sugar strategy); F4 (linear-algebra reductions); F5 (signatures) | Buchberger, Aequationes Math. 4 (1970) 374-383, doi:10.1007/BF01844169 ; Faugère F4, J. Pure Appl. Algebra 139 (1999) 61-88, doi:10.1016/S0022-4049(99)00005-5 ; Faugère F5, ISSAC 2002, doi:10.1145/780506.780516 | F4 is the practical default; compute mod p then lift |
| Rational integration | Hermite reduction + Lazard-Rioboo-Trager | Bronstein, Symbolic Integration I (above) | no factorisation over extension fields |
| Transcendental integration | Risch algorithm | Risch, Trans. AMS 139 (1969) 167-189, doi:10.1090/S0002-9947-1969-0237477-8 ; Bronstein (above) | complete decision procedure for elementary functions |
| Practical integration | Risch-Norman heuristic (ansatz + linear solve) ; rule-based (Rubi) | Rich, Scheibe, Abbasi, "Rule-based integration: an extensive system of symbolic integration rules", JOSS 3(32) 2018, doi:10.21105/joss.01073 ; https://rulebasedintegration.org/ | Rubi: 6,700+ rules in a decision tree, 72,000+ test problems with optimal antiderivatives (rulebasedintegration.org, fetched 2026-10-17) |
| Limits, asymptotics | Gruntz algorithm (comparability classes, MRV set) | Gruntz, "On computing limits in a symbolic manipulation system", PhD thesis, ETH Zurich 1996, ETH Diss. 11432 (no DOI found) | what SymPy `limit` implements |
| Series | power series arithmetic with lazy/truncated coefficients; Newton iteration | von zur Gathen & Gerhard ch. 9; MCA ch. 4 | |
| Linear algebra over Q | Bareiss fraction-free; multi-modular + CRT; Dixon p-adic lifting | von zur Gathen & Gerhard ch. 5; Geddes et al. ch. 9 | |
| Solving polynomial systems | Gröbner (lex via FGLM), triangular sets, resultants | von zur Gathen & Gerhard ch. 21 ; Geddes et al. ch. 10 | |
| Pattern matching / rewriting | AC (associative-commutative) matching, discrimination nets | Rubi's implementation as test; Mathematica docs https://reference.wolfram.com/language/tutorial/Patterns.html | rule-based integration needs it |
| Output | LaTeX / MathML printers | SymPy `printing.latex` (oracle); `go-tex` renders | |

## 3. Reference implementations

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| MPFR | C | LGPL-3.0+ (mpfr.org) | oracle only | the model for correct rounding; https://www.mpfr.org/ |
| GMP | C | LGPL-3 / GPL-2 dual (gmplib.org) | oracle only | |
| Arb | C | LGPL-2.1 (GitHub `flintlib/arb`, last push 2024-03-16; merged into FLINT 3) | oracle only | ball arithmetic; https://arblib.org/ |
| FLINT | C | LGPL-3.0 (GitHub `flintlib/flint`, active 2026-10-16) | oracle only | polynomials, factorisation, linear algebra, number theory incl. ECM and quadratic sieve; https://flintlib.org/ |
| python-flint | Python/Cython | MIT (GitHub `flintlib/python-flint`) | oracle (wrapper itself MIT but binds LGPL FLINT) | easiest way to drive FLINT/Arb |
| mpmath | Python | BSD-3-Clause (GitHub `mpmath/mpmath`, active 2026-10-04) | port permitted (algorithms and code BSD) | the oracle go-scikit uses today; clear pure-Python reference for special functions, quadrature, `identify`, `findroot` |
| PARI/GP | C | GPL-2+ (pari.math.u-bordeaux.fr FAQ) | oracle only | number theory reference (`isprime`, `factorint`, `ellmul`, `qfbclassno`); https://pari.math.u-bordeaux.fr/ |
| SymPy | Python | BSD-3 (LICENSE read; GitHub reports NOASSERTION because it appends third-party notices) | port permitted file by file after reading the appended notices | reference for Gruntz, Risch (`ratint`, `risch_integrate`), `solve`, `simplify`; https://www.sympy.org/ |
| SymEngine | C++ | MIT core plus BSD-3 for `utilities/teuchos` (LICENSE read; GitHub NOASSERTION) | port permitted | fast core, expression design reference |
| FriCAS | Common Lisp | BSD-3-Clause (GitHub `fricas/fricas`, API-verified 2026-10-17) | port permitted | most complete Risch implementation (Bronstein's algorithms in the Axiom lineage) |
| Axiom / OpenAxiom | Lisp | modified BSD (verify) | port permitted | |
| Reduce | Lisp | BSD-2 (`reduce-algebra/reduce-algebra` LICENSE read; GitHub NOASSERTION) | port permitted | https://reduce-algebra.sourceforge.io/ |
| Maxima | Common Lisp | GPL-2 (verify) | oracle only | `integrate`, `limit`, `solve`; https://maxima.sourceforge.io/ |
| Giac / Xcas | C++ | GPL-3 (verify) | oracle only | fast Gröbner (F4), https://www-fourier.univ-grenoble-alpes.fr/~parisse/giac.html |
| Sage | Python | GPL-3 distribution, original code GPL-2+ (`sagemath/sage` COPYING.txt) | oracle only | glue over all the above |
| Rubi (rules + tests) | Mathematica (translated to Maxima, SymPy, Mathematica) | MIT (GitHub `RuleBasedIntegration/Rubi`, API-verified; last push 2024-02-22) | port permitted | rules and 72,000-problem test suite are reusable data |
| Mathematica / MATLAB Symbolic | closed | proprietary | doc reference only | no output copying beyond fair use |
| `math/big` (Go stdlib) | Go | BSD-3 | in-tree | Int, Rat, Float; ProbablyPrime = MR + BPSW; Float has Sqrt only |
| `ALTree/bigfloat` | Go | MIT (API-verified; last push 2026-08-18) | port/depend | Exp, Log, Pow on `big.Float`; base for a multiprecision layer |
| `remyoudompheng/bigfft` | Go | BSD-3 (API-verified; last push 2023-01-29) | depend | FFT mult for `big.Int` |
| `modernc.org/mathutil` (ex cznic) | Go | BSD-3 (API-verified) | depend | small helpers, primality for uint64 |
| `modernc.org/nerdamer` ("nerdamer-go") | Go | BSD-3-Clause (pkg.go.dev, published 2026-08-29) | not a Go CAS | wrapper: back end is Nerdamer-prime (JavaScript); scope small |
| `Konstantin8105/sm` | Go | MIT (API-verified; last push 2023-10-14) | port/depend | symbolic simplify and diff, small |
| `rpopplewell/calculus` | Go | MIT (API-verified; last push 2024-04-04) | port | single-variable diff/integrate, numerics-flavoured |
| `ncw/gmp` | Go + cgo | BSD-3 (API-verified) | not usable (cgo) | violates CGO=0 |
| `go-cas` | GitHub organisation (2015), 1 repo `cas` | MIT | n/a | "Golang CAS Client package" (Central Authentication Service, unrelated to computer algebra); name taken |

No Go library covers: ball arithmetic, correctly-rounded special functions, polynomial factorisation, Gröbner bases, or Risch integration (searched 2026-10-17: pkg.go.dev and GitHub results above; absence is not proven).

## 4. Validation

### Oracle tools

| tool | version / obtain (no root) | drive from Go |
|---|---|---|
| mpmath | `uv run --with mpmath python script.py` (1.x, BSD) | exec python, JSON on stdin/stdout; values as decimal strings, `mp.dps` set explicitly |
| python-flint | `uv run --with python-flint python script.py` (wheels for macOS arm64 and Linux; verify per release) | same pipe; use `arb`/`acb` and read `.rad()` to get certified enclosures; `fmpz_poly.factor()`, `fmpz.is_prime()`, `fmpz_mpoly` |
| PARI/GP | release tarball built in `$HOME`, or pkgx/brew user prefix (verify availability); `gp -q` | exec `gp -q -f`, print with `print1`; JSON via `printf`; exact integers/rationals as strings |
| SymPy | `uv run --with sympy` | exec python; compare by `simplify(a-b)==0` or numeric spot-checks at random points, not by string equality |
| Maxima | release binary or `pkgx maxima` (verify) | `maxima --very-quiet --batch-string`; `tex()` output for LaTeX |
| FriCAS | build from source (BSD) | verification of Risch results |
| Rubi suite | https://rulebasedintegration.org/ (Test Problems menu); repo `RuleBasedIntegration/Rubi` (MIT) | convert to JSON of (integrand, variable, optimal antiderivative); check by differentiation |

### Public corpora

| corpus | URL | size | licence |
|---|---|---|---|
| Rubi integration test suite | https://rulebasedintegration.org/ | 72,000+ problems | MIT (repo) |
| mpmath test suite | https://github.com/mpmath/mpmath/tree/master/mpmath/tests | thousands of assertions (BSD-3) | BSD-3 |
| FLINT / Arb tests | https://github.com/flintlib/flint | `t-*.c` per module | LGPL (data only as oracle) |
| DLMF formulas and tables | https://dlmf.nist.gov/ | ~ 36 chapters; check reuse terms at https://dlmf.nist.gov/about/ (verify) | NIST terms |
| OEIS | https://oeis.org/ ; b-files per sequence | 380,000+ sequences (verify) | check https://oeis.org/wiki/Welcome#Copyright (verify) |
| Known factorisations | Cunningham project tables, https://homes.cerias.purdue.edu/~ssw/cun/ | n/a (verify) | n/a |

### What "correct" means and what to report

| area | criterion | measure |
|---|---|---|
| Multiprecision floats | correct rounding: result == round(exact) in the given mode, 0 ULP | count of mismatches vs MPFR-grade oracle at 53, 113, 256, 1024, 10000 bits and all five rounding modes; report 0 |
| Ball arithmetic | enclosure: oracle value lies inside [mid-rad, mid+rad]; and tightness | enclosure failures (must be 0); median log2(rad/ulp) |
| Special functions at N digits | relative error < 10^-(N-2) vs mpmath at N+20 digits (mpmath itself differs from Arb at hard points: take Arb as the arbiter when they disagree) | max relative error in ULP at target precision |
| Number theory | exact equality | pass/fail; factorisation: product == n and every factor prime (check by ProbablyPrime + BPSW pseudoprime-free ranges); primality vs PARI `isprime` on N random + known pseudoprime lists |
| Polynomials, linear algebra over exact fields | exact equality after canonical form | pass/fail; for factorisation: product equals input, factors irreducible per oracle |
| Gröbner | reduced Gröbner basis is unique for a given order: exact equality with Singular/Giac/Sage output | pass/fail |
| Integration | derivative of result minus integrand simplifies to 0; and a "grade" per Rubi (A optimal, B ok but larger, C, F fail) | grade histogram on the Rubi suite; report counts of A/B/C/F |
| Limits, series, solve | exact equality of normal forms, or numeric spot-checks at random points with N-digit arithmetic | pass rate; cases where `Solve` returns no solution when the oracle does |
| Output | LaTeX compiles and typesets via go-tex; round-trips to the same tree via a parser | compile success; round-trip equality |

Go harness shape: one `oracle` package with `Run(ctx, tool, inputJSON) (outputJSON, error)`, golden files committed, oracle used only to regenerate; the unit suite must pass with no oracle installed.

## 5. Proposed home and plan

Names free (`gh api users/<name>` 404) on 2026-10-17: go-numtheory, go-computeralgebra, go-multiprecision, go-exactmath, go-arith, go-algebra. TAKEN: `go-cas` (unrelated CAS-client org, 2015) and `go-bignum` (organisation created 2026-10-17 20:16 UTC, 0 repos: someone, presumably this campaign, already claimed it; confirm ownership before relying on it).

| work item | org/repo (existing, or NEW org to create) | depends on | size S/M/L/XL | priority 1–3 |
|---|---|---|---|---|
| Correctly-rounded `big.Float` elementary functions (exp, log, sin, cos, atan, pow, hyperbolics, constants pi/e/ln2) with Ziv loop and five rounding modes | NEW org `go-bignum`, repo `float` (reuse ideas from ALTree/bigfloat, MIT) | stdlib `math/big`, `bigfft` | L | 1 |
| Test harness vs mpmath/python-flint with golden JSON (self-hosting the accuracy checks), then flip `go-numerics/special` validation to it | `go-bignum/oracle` + `go-numerics/special` | item above | M | 1 |
| Ball arithmetic (real and complex), rigorous elementary and hypergeometric core, then gamma, zeta, Bessel, erf at N digits | `go-bignum/ball` | float | XL | 2 |
| Number theory: sieve, factor (rho, ECM, SIQS), CRT, discrete log, partitions, Bernoulli, continued fractions, ECPP later | NEW org `go-numtheory` | stdlib `big.Int` | L | 2 |
| Computer algebra core: expression DAG, canonical forms, Z/Q/F_p polynomials, GCD, factorisation, rational functions, `diff`, `expand`, LaTeX output | NEW org `go-computeralgebra`, repo `expr` + `poly` | go-numtheory (primes), go-tex (render) | XL | 2 |
| Gröbner (Buchberger then F4 mod p) and polynomial-system solve | `go-computeralgebra/groebner` | poly | L | 3 |
| Integration: rational (Hermite/LRT), Risch-Norman, rule engine and Rubi rule import with the 72,000-problem suite as the grade | `go-computeralgebra/integrate` | expr, poly, pattern matcher | XL | 3 |
| Series and limits (Gruntz) | `go-computeralgebra/series` | expr, poly | L | 3 |
| Exact linear algebra (Bareiss, Dixon, Smith/Hermite) | `go-computeralgebra/linalg` | poly, numtheory | M | 3 |
| Multiprecision matrices (LU/QR/eig on big.Float, ball matrices) | `go-bignum/matrix` | float, ball, go-ndarray | M | 3 |

## 6. Risks and open questions

- Scope: a CAS comparable to Mathematica is a decades-long target; Rubi, Risch, Gröbner and polynomial factorisation are each a project. Priority 1 items (correctly rounded floats, oracle harness) are the only ones with immediate payoff for go-scikit's accuracy validation.
- Correct rounding is expensive to prove: a Ziv loop is only correct if the error bounds on each function are rigorous. Without ball arithmetic underneath, bounds must be hand-derived per function (mpmath-style "guard bits" are not a proof). Decision needed: build floats on top of balls from the start.
- Oracle trust: mpmath is not certified; where mpmath and Arb/FLINT disagree, trust Arb/FLINT enclosures. Record which oracle version produced each golden file.
- Licence hygiene: MPFR/FLINT/Arb/PARI/Maxima/Giac/Sage are (L)GPL: use as oracles and read papers, never copy code. mpmath (BSD), FriCAS (BSD), Rubi (MIT) are portable. SymPy, SymEngine, Reduce need their LICENSE files read before any port (GitHub reports NOASSERTION).
- Integration results are not canonical: two correct antiderivatives can differ by a constant or a form; the oracle comparison must be by differentiation, and Rubi grades measure size of the antiderivative.
- Performance: Go `math/big` has no assembly tuned for FFT sizes and uses Karatsuba only; ten-million-digit constants need `bigfft`. GC pressure on many small `big.Int` is a risk for polynomial arithmetic; prefer word-sized (uint64) coefficient fast paths for F_p.
- Unverified in this file: Maxima, Giac, Axiom/OpenAxiom licence versions, OEIS and DLMF reuse terms, availability of PARI/Maxima through pkgx on macOS arm64.
- `go-cas` is an unrelated CAS-client (authentication) repository; no conflict with `go-computeralgebra`.

Verified on 2026-10-17: all 24 DOIs resolved via Crossref against title/year/first author (van Hoeij DOI corrected to 10.1006/jnth.2001.2763, BFW DOI added, Gruntz thesis has none); MATLAB names checked on mathworks.com (`degree` -> `polynomialDegree`; `discriminant`, `reduce` removed), Mathematica names on reference.wolfram.com; licences of MPFR, Arb, FLINT, python-flint, mpmath, SymPy, SymEngine, FriCAS, Reduce, Sage, Rubi, PARI and the Go libraries via GitHub API/LICENSE; go-scikit state via GitHub main.
