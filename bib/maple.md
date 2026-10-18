# Maple — delta bibliography

One file per domain. Facts only, each with a source (URL or DOI). Terse: tables, not prose.
Written 2026-10-18. **Delta only**: this file lists what Maple brings that the other `bib/` files do
not already cover, or where Maple is the reference implementation and changes an algorithm choice
or the validation oracle. For everything else see
[precision-numbertheory-symbolic.md](precision-numbertheory-symbolic.md) (computer algebra),
[systems-units-formats-visualization.md](systems-units-formats-visualization.md) (DAE, units,
notebook), [optimization-learning-autodiff.md](optimization-learning-autodiff.md) (statistics, AD).

Maple help pages are cited as `help:<path>` = `https://www.maplesoft.com/support/help/Maple/view.aspx?path=<path>`
(all fetched 2026-10-18). Pages that returned HTTP 410 on that date and were therefore NOT read:
`odeadvisor` (the working path is `DEtools/odeadvisor`), `DifferentialGeometry/references`,
`DAEtools`, `worksheet/documenttools`, `CodeGeneration/LanguageDefinition`. Nothing is claimed
from them. Maple is proprietary: no Maple output is copied here; its pages are the only thing quoted.

## 1. Feature inventory (Maple delta)

"In plan?" = already listed by name in the existing bib files. "go-scikit today" = GitHub main
checked with `gh api` on 2026-10-18: `go-symbolic`, `go-control`, `go-units`, `go-bignum`,
`go-autodiff`, `go-fem` exist as organisations created 2026-10-17 with **no repository**;
`go-numerics/integrate` has RK/Radau/BDF/LSODA-like and events but no DAE or mass matrix (file list
and README grepped); `go-notebook/notebook` exists (Markdown file, Go cells, dependency order,
kept results); `go-tex/math` exists (pure-Go TeX math to SVG, parses a TeX subset, wasm); `go-statistics`
has only `stats`.

| area | Maple capability (command, page) | in plan? | go-scikit today | delta |
|---|---|---|---|---|
| ODE classification | `DEtools[odeadvisor]` returns the matching types: first order Abel, Bernoulli, Chini, Clairaut, dAlembert, exact, homogeneous (several), linear, quadrature, Riccati, separable, ...; second order Bessel, Duffing, Emden, Gegenbauer, Hermite, Jacobi, Laguerre, Lienard, Liouville, Painleve, Van_der_Pol, ...; higher order exact_linear, missing, reducible, linear_ODEs (`help:DEtools/odeadvisor`) | no (`dsolve`/`DSolve` appear by name only) | none | a typed catalogue of ODE classes with a recogniser |
| Symbolic ODE solving | `dsolve` combines "classification, integrating factor and symmetry methods" (`help:dsolve/details`); option `Lie` = symmetry-only, `class` = avoid symmetries, `way=` selects a `DEtools[symgen]` algorithm, `HINT=` suggests infinitesimals; integration schemes `fat`, `can`, `can2`, `gon`, `gon2`, `dif` (`help:dsolve/Lie`); `DEtools[intfactor]`; Riccati/Abel solutions in Liouvillian or hypergeometric terms detected automatically; `HeunB`, `HeunC` outputs (`help:dsolve/details`); `odetest` checks a solution | no | none | the whole ODE stack |
| Linear ODE decision procedures | "decision procedures like Kovacic's algorithm", hypergeometric solutions (`help:dsolve/details`, `help:dsolve/hypergeometric`) | no | none | Kovacic, hypergeometric-type solutions |
| PDE solving | `pdsolve`: separation of variables (`HINT = '+'` / `'*'`), characteristic strips (`HINT = strip`, `charstrip`, `splitstrip`), travelling waves (`HINT = TWS`), `build`, `INTEGRATE`, options of `PDEtools[casesplit]` for systems (`help:pdsolve`) | no | none | PDE heuristics and system splitting |
| Algebraic numbers | `RootOf` is "the standard representation for Maple algebraic numbers, algebraic functions (see evala), and finite fields" (`help:RootOf`); `evala` computes in the smallest algebraic number/function field and offers `Gcd`, `Factor`, `Factors`, `Sqrfree`, `Normal`, `Quo`, `Rem`, `Resultant`, `Norm`, `Trace`, `Minpoly`, `Primfield`, `Subfields`, `AFactor`; for reducible `RootOf`s it does "a case discussion and recombination"; checks independence of the `RootOf`s unless told `'independent'` (`help:evala`) | partly: Mathematica `Factor[..., Extension->]` and `Root` objects exist, the plan lists `Extension->` only inside factorisation | none | algebraic number fields as a first-class domain |
| Triangular decomposition | `RegularChains`: `Triangularize`, `RealTriangularize`, `RegularGcd`, `ExtendedRegularGcd`, `Intersect`, `NormalForm`, `Inverse`, `SamplePoints`, `ComprehensiveTriangularize`; subpackages `MatrixTools`, `ConstructibleSetTools`, `SemiAlgebraicSetTools`, `ParametricSystemTools`, `ChainTools`, `FastArithmeticTools`, `AlgebraicGeometryTools` (`help:RegularChains`). `SemiAlgebraicSetTools[CylindricalAlgebraicDecompose]` defaults to the incremental algorithm CM (Chen, Moreno Maza), alternative `method=recursive` = CMXY, with options EC and RC-TTICAD (`help:RegularChains/SemiAlgebraicSetTools/CylindricalAlgebraicDecompose`) | plan names "triangular sets" and "CAD (later)" once, no algorithm | none | regular chains and a triangular-decomposition CAD |
| Differential algebra | `DifferentialAlgebra` (based on Boulier's BLAD): `RosenfeldGroebner` (cases = regular differential chains, representation of the radical of the differential ideal, singular cases), `DifferentialRing`, `BelongsTo`, `ReducedForm`, `PowerSeriesSolution`, `Tools[DifferentialPrem]` = Ritt reduction (`help:DifferentialAlgebra`, `help:DifferentialAlgebra/RosenfeldGroebner`) | no | none | differential rings, rankings, elimination |
| Code generation | `CodeGeneration` translates Maple code to C, CSharp, Fortran, Java, JavaScript, Julia, Matlab, Perl, Python, R, Swift, VisualBasic; `LanguageDefinition` extends it with new target languages (named on the page; its own page returned 410) (`help:CodeGeneration`); `codegen[optimize]` = common-subexpression elimination, temporaries `t0, t1, ...`, option `tryhard` (`help:codegen/optimize`) | partly: MATLAB `ccode`, `fortran`, `matlabFunction`, Mathematica `CForm`, `FortranForm`, `Compile` are listed with the target "Lambdify-style compile to Go closure" | none | Go as a target language, CSE, and a pluggable target definition |
| Tensors and mathematical physics | `Physics`: noncommutative variables, `Bra`, `Ket`, `Dagger`, `Commutator`, `AntiCommutator`, tensors with spacetime/space/tetrad/spinor indices, `Christoffel`, `Einstein`, `Ricci`, `Riemann`, `Weyl`, `KillingVectors`, `Psigma`, `Dgamma`, `Setup`, `Vectors` (`help:Physics`) | no (tensor contraction in linear-algebra.md is numeric `einsum` only) | none | abstract-index tensor algebra |
| Differential geometry, Lie algebras | `DifferentialGeometry` with `Tensor`, `LieAlgebras`, `JetCalculus`, `GroupActions`, `Tools`; `DGsetup`; `Library` of Lie algebras, vector-field systems, differential equations and exact solutions in general relativity (`help:DifferentialGeometry`) | no | none | Lie algebra computations, jet calculus |
| Group theory | `GroupTheory`: permutation, finitely presented, Cayley-table, custom and symbolic groups; `PermutationGroup`, `GroupOrder`, `ConjugacyClasses`, `CayleyTable`, `IsSolvable`, `SylowSubgroup`, `SylowBasis`, `SylowTower`, `FPGroup`, `Relators`, `SmallGroup`, `AllSmallGroups`, `SearchSmallGroups`, `IdentifySmallGroup` (`help:GroupTheory`) | no (the plan lists none; whether Mathematica's permutation-group functions need a separate inventory entry was not checked) | none (stdlib has no group type) | permutation and finitely presented groups |
| Exact and modular linear algebra | `LinearAlgebra[Modular]`: `Mod`, `Multiply`, `Inverse`, `LUDecomposition`, `RowReduce`, `Rank`, `RankProfile`, `Determinant`, `CharacteristicPolynomial`, `LinearSolve`, `ChineseRemainder`, `IntegerDeterminant`, `IntegerLinearSolve`; datatypes `integer[4]` (moduli up to 2^16-1), `integer[8]` (up to 2^32-1), `float[8]` (up to 2^25-1, uses BLAS), cites Dumas and Pernet 2001 (`help:LinearAlgebra/Modular`) | partly: "Bareiss, Dixon, Smith/Hermite" in precision file, no word-size-prime layer | none | word-size prime field kernels, rank profile |
| Symbolic statistics | `Statistics`: 35+ distributions plus `Distribution`, `RandomVariable`; symbolic `PDF`, `CDF`, `MGF`, `CGF`, `CharacteristicFunction`, `ExpectedValue`, `Moment`, `Cumulant`, `Mean`, `Variance`, `Skewness`, `Kurtosis`; estimation `MaximumLikelihoodEstimate`, `FisherInformation`, `LogLikelihood`; resampling `Bootstrap`; tests `OneSampleTTest`, `ShapiroWilkWTest`, ... (`help:Statistics`; the page does not state which results are symbolic) | numeric side yes (go-statistics `models`, `timeseries`); symbolic side no | `go-statistics/stats` (numeric) | random variables as symbolic objects |
| Units | `Units` package with environments Standard, Natural, Simple (default); `AddUnit`, `AddSystem`, `UsingSystem`, `GetDimensions`; systems listed by `GetSystems`: Atomic, CGS, EMU, ESU, FPS, IPS, MKS, MTS, SI; conversion by `convert/units` and `convert/temperature` (`help:Units`; the page cites neither NIST SP 811 nor CODATA) | yes (`symunit`, `Quantity`; go-units) | none | user-defined unit systems (EMU, ESU, Atomic) only |
| Document model | native `.mw` = XML (introduced in Maple 9), Document and Worksheet modes, 2-D math as the default input style, 1-D Maple notation optional, typesetting levels standard/extended, `DocumentTools` (search excerpt of `worksheet/documenting/...` and `Formats/MW` pages, 2026-10-18; pages not read in full); 22 embedded components in `DocumentTools:-Components`: Button, CheckBox, CodeEditRegion, ComboBox, DataTable, Dial, Label, ListBox, MathContainer, Meter, Microphone, Plot, RadioButton, RotaryGauge, Shortcut, Slider, Speaker, State, TextArea, ToggleButton, VideoPlayer, VolumeGauge (`help:DocumentTools/Components`) | `Manipulate` yes (systems file 2f); 2-D math input no | go-notebook: Markdown + Go cells; go-tex/math renders a TeX subset | structured math input and editing; widgets beyond sliders |
| Physical modelling (MapleSim) | MapleSim API: `CreateModel` from a Modelica file, `GetModelica`, `GetEquations`, `GetCompiledProc` (compiles the subsystem "into a C Function"), `Tools` generates Modelica from Maple (`help:MapleSim`). Vendor claims: "fully symbolic math engine", symbolic routines "can actually reduce the index of the DAEs" (https://www.maplesoft.com/products/maplesim/symbolic_computation.aspx); numeric DAE solvers in `dsolve`: `rkf45_dae`, `ck45_dae`, `rosenbrock_dae`, `mebdfi` (`help:dsolve/numeric/DAE`) | solvers yes (DAE, Pantelides, BLT in systems file 2b); the Modelica-style acausal compiler is marked "out of scope" there | `integrate` has no DAE | symbolic equation assembly feeding Pantelides, then code generation |

### Not a delta (checked, already covered)

`Statistics` descriptive and regression commands (`Fit`, `LinearFit`, `NonlinearFit`, `Lowess`,
`PCA`, `KernelDensity`, `Bootstrap`, ...) match MATLAB/Mathematica rows in the optimisation file.
Maple's `Units` conversion core duplicates `symunit`/`Quantity`. `Manipulate`-style sliders are
covered by the `go-notebook` item. Maple's plotting, MapleSim component libraries and Maple
Learn/Cloud content are not inventoried (proprietary content).

## 2. Reference algorithms

Maple's help pages name few algorithms: `evala`, `Triangularize` (one reference), `RosenfeldGroebner`
(short labels) and `codegen[optimize]` (no paper) give none or little. The algorithms below come
from the literature those pages and the Maple authors point to. DOIs were resolved through
Crossref on 2026-10-18 (title, year, first author).

### (a) Symbolic ODE and PDE

| capability | algorithm | reference | why |
|---|---|---|---|
| Second-order linear ODE, Liouvillian solutions | Kovacic | Kovacic, "An algorithm for solving second order linear homogeneous differential equations", J. Symbolic Comput. 2(1) 1986, doi:10.1016/S0747-7171(86)80010-4 (cited on `help:dsolve/references`) | decision procedure: finds a Liouvillian solution or proves none exists |
| n-th order linear, Liouvillian | Singer | Singer, "Liouvillian solutions of n-th order homogeneous linear differential equations", Amer. J. Math. 103(4) 1981, doi:10.2307/2374045 ; van der Put, Singer, "Galois Theory of Linear Differential Equations", Springer 2003, doi:10.1007/978-3-642-55750-7 | the theory behind Kovacic |
| Beyond Liouvillian | solutions through reduction to hypergeometric/Heun type | Chan, Cheb-Terrab, "Non-Liouvillian solutions for second order linear ODEs", ISSAC 2004, doi:10.1145/1005285.1005299 ; Cheb-Terrab, "Solutions for the general, confluent and bi-confluent Heun equations and their connection with Abel equations", J. Phys. A 37(42) 2004, doi:10.1088/0305-4470/37/42/007 (both on `help:dsolve/references`) | what lets `dsolve` return `HeunC`-type answers |
| First- and second-order nonlinear ODE | Lie symmetries: ansatz for infinitesimals (`symgen`), canonical coordinates, integrating factors | Cheb-Terrab, Duarte, da Mota, Comput. Phys. Commun. 101 (1997) 254, doi:10.1016/S0010-4655(97)00018-0 (first order) and 108 (1998) 90, doi:10.1016/S0010-4655(97)00132-X (second order); Cheb-Terrab, Roche, "Symmetries and first order ODE patterns", CPC 113 (1998) 239, doi:10.1016/S0010-4655(98)00071-X ; Olver, "Applications of Lie Groups to Differential Equations", 2nd ed., GTM 107, Springer 1993, doi:10.1007/978-1-4612-4350-2 ; Bluman, Kumei, "Symmetries and Differential Equations", Springer 1989, doi:10.1007/978-1-4757-4307-4 (the last two credited with Stephani on `help:dsolve/Lie`; the first two are the page's references) | heuristic, not complete: effective on equations that no classification pattern matches |
| ODE classification | pattern catalogue keyed on Kamke, Zwillinger | Kamke, "Differentialgleichungen: Loesungsmethoden und Loesungen", Chelsea 1959 ; Zwillinger, "Handbook of Differential Equations", 2nd ed., Academic Press 1992 (both on `help:DEtools/odeadvisor`) ; Cheb-Terrab, Roche, "Abel ODEs: equivalence and integrable classes", CPC 130 (2000) 204, doi:10.1016/S0010-4655(00)00042-4 | each class has a closed-form solution recipe; cheaper than symmetries |
| PDE | heuristic separation, characteristic strips, travelling waves, splitting into cases | Cheb-Terrab, von Buelow, "A computational approach for the analytical solving of partial differential equations", CPC 90 (1995) 102, doi:10.1016/0010-4655(95)00083-R (cited on `help:pdsolve`) | |

### (b) Algebraic numbers and exact linear algebra

| capability | algorithm | reference | why |
|---|---|---|---|
| Computing with `RootOf` when the minimal polynomial may be reducible | dynamic evaluation: split into cases when a non-invertible element appears, recombine | Della Dora, Dicrescenzo, Duval, "About a new method for computing in algebraic number fields", EUROCAL '85, LNCS 204, doi:10.1007/3-540-15984-3_279 | this is the behaviour `help:evala` describes ("case discussion and recombination"); the attribution to the D5 principle is mine, not the page's |
| Factoring over number fields | norm method | Trager, "Algebraic factoring and rational function integration", SYMSAC '76, doi:10.1145/800205.806338 ; Cohen, "A Course in Computational Algebraic Number Theory", GTM 138, Springer 1993, doi:10.1007/978-3-662-02945-9 | |
| GCD over number fields | modular GCD with rational reconstruction, several extensions | van Hoeij, Monagan, "A modular GCD algorithm over number fields presented with multiple extensions", ISSAC 2002, doi:10.1145/780506.780520 | avoids coefficient swell |
| Linear algebra over F_p, word-size p | blocking onto floating-point BLAS, rank profile, LU | Dumas, Gautier, Pernet, "Finite field linear algebra subroutines", ISSAC 2002, doi:10.1145/780506.780515 ; Dumas, Giorgi, Pernet, "Dense linear algebra over word-size prime fields: the FFLAS and FFPACK packages", ACM TOMS 35(3) 2008, doi:10.1145/1391989.1391992 | the Maple page's `float[8]` datatype (moduli below 2^25) is this idea |

### (c) Triangular decomposition, differential algebra

| capability | algorithm | reference | why |
|---|---|---|---|
| RegularChains library | regular chains, regular GCD | Lemaire, Moreno Maza, Xie, "The RegularChains library in MAPLE", ACM SIGSAM Bulletin 39(3) 2005, doi:10.1145/1113439.1113456 | |
| Triangular decomposition of polynomial systems | incremental, Kalkbrener and Lazard senses | Chen, Moreno Maza, "Algorithms for computing triangular decomposition of polynomial systems", J. Symbolic Comput. 47(6) 2012, doi:10.1016/j.jsc.2011.12.023 ; Moreno Maza, "On triangular decompositions of algebraic varieties", MEGA 2000 (the only reference on `help:RegularChains/Triangularize`; no DOI found) | triangular sets avoid Gröbner-basis blow-up in lex order |
| Probabilistic option of `Triangularize` | lifting | Dahan, Moreno Maza, Schost, Wu, Xie, "Lifting techniques for triangular decompositions", ISSAC 2005, doi:10.1145/1073884.1073901 | |
| Real solving, CAD | CAD by incremental triangular decomposition (CM), recursive (CMXY), equational constraints, TTICAD | Collins, "Quantifier elimination for real closed fields by cylindrical algebraic decomposition", LNCS 33, 1975, doi:10.1007/3-540-07407-4_17 (Crossref returns the title with the misspelling "decompostion") ; Chen, Moreno Maza, Xia, Yang, "Computing cylindrical algebraic decomposition via triangular decomposition", ISSAC 2009, doi:10.1145/1576702.1576718 ; Chen, Moreno Maza, "An incremental algorithm for computing cylindrical algebraic decompositions", Springer 2014, doi:10.1007/978-3-662-43799-5_17 ; Chen, Davenport, May, Moreno Maza, Xia, Xiao, "Triangular decomposition of semi-algebraic systems", ISSAC 2010, doi:10.1145/1837934.1837972 | CAD is doubly exponential in the number of variables; the triangular route is Maple's default (`help:RegularChains/SemiAlgebraicSetTools/CylindricalAlgebraicDecompose`) |
| Differential elimination | Rosenfeld-Groebner, regular differential chains, Ritt reduction | Ritt, "Differential Algebra", AMS Colloquium 33, 1950, doi:10.1090/coll/033 ; Kolchin, "Differential Algebra and Algebraic Groups", Academic Press 1973 (no DOI) ; Boulier, Lazard, Ollivier, Petitot, "Representation for the radical of a finitely generated differential ideal", ISSAC 1995, doi:10.1145/220346.220367 ; same authors, "Computing representations for radicals of finitely generated differential ideals", AAECC 20 (2009) 73-121, doi:10.1007/s00200-009-0091-7 ; Hubert, "Notes on triangular sets and triangulation-decomposition algorithms II: differential systems", LNCS 2630 (2003), doi:10.1007/3-540-45084-X_2 | Ritt and Kolchin are the labels [R50, K73] on the Maple pages |

### (d) Code generation

| capability | algorithm | reference | why |
|---|---|---|---|
| Common subexpression elimination | hash-consing of the expression DAG (linear time); `codegen[optimize]` finds only syntactically identical subexpressions, "x+y" is not shared with "x+y+z"; `tryhard` is a costlier second algorithm (`help:codegen/optimize`) | Downey, Sethi, Tarjan, "Variations on the common subexpression problem", J. ACM 27(4) 1980, doi:10.1145/322217.322228 | a hash-consed DAG in go-symbolic gives the Maple behaviour for free; going further (shared sums) is a choice |
| Evaluation cost | multivariate Horner schemes | Ceberio, Kreinovich, "Greedy algorithms for optimizing multivariate Horner schemes", ACM SIGSAM Bulletin 38(1) 2004, doi:10.1145/980175.980179 | fewer operations for polynomials |
| Derivative code | reverse-mode sharing | Griewank, Walther, "Evaluating Derivatives", 2nd ed., SIAM 2008, doi:10.1137/1.9780898717761 | the generated code for gradients should match go-autodiff's tape |
| Floating-point semantics of the emitted Go | Go may fuse `x*y+z` across statements; an explicit `float64(x*y)` conversion forbids fusion (Go spec, section Floating-point operators, https://go.dev/ref/spec, read 2026-10-18) | | generated code is bit-reproducible across amd64/arm64/ppc64le/s390x only if every product is wrapped; golden tests must pick a policy |

### (e) Tensors, Lie algebras, groups

| capability | algorithm | reference | why |
|---|---|---|---|
| Tensor index canonicalisation | Butler-Portugal, permutation-group based | Portugal, "Algorithmic simplification of tensor expressions", J. Phys. A 32(44) 1999, doi:10.1088/0305-4470/32/44/313 ; Martin-Garcia, "xPerm: fast index canonicalization for tensor computer algebra", Comput. Phys. Commun. 2008, doi:10.1016/j.cpc.2008.05.009 ; Peeters, "Cadabra: a field-theory motivated symbolic computer algebra system", CPC 2007, doi:10.1016/j.cpc.2007.01.003 | needs the group layer below |
| Permutation groups | Schreier-Sims, base and strong generating set | Sims, "Computational methods in the study of permutation groups", in Computational Problems in Abstract Algebra, Pergamon 1970, doi:10.1016/B978-0-08-012975-4.50020-5 ; Butler, "Fundamental Algorithms for Permutation Groups", LNCS 559, 1991, doi:10.1007/3-540-54955-2 ; Seress, "Permutation Group Algorithms", CUP 2003, doi:10.1017/CBO9780511546549 ; Holt, Eick, O'Brien, "Handbook of Computational Group Theory", CRC 2005, doi:10.1201/9781420035216 (the last three appear in `help:GroupTheory/references`: Butler 1991, Seress 2003, Holt-Eick-O'Brien 2005) | |
| Small groups | table lookup, not computation | Besche, Eick, O'Brien, The Small Groups library (cited by `help:GroupTheory/references`) | data licence to read before use |

## 3. Reference implementations

Maple itself: proprietary, doc reference only, not runnable in CI.

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| SymPy `solvers/ode` (`ode.py`, `single.py` classification hints, `lie_group.py` 39 kB, `riccati.py`, `hypergeometric.py`, `systems.py`, `nonhomogeneous.py`, `subscheck.py`) and `solvers/pde.py` | Python | BSD-3 text (LICENSE read; GitHub NOASSERTION, see precision file) | port permitted file by file after reading appended notices | closest free analogue of `dsolve` + `odeadvisor`; `classify_ode` and `checkodesol` are the oracle APIs |
| FriCAS `kovacic.spad` (5.8 kB), `odeef.spad`, `oderf.spad`, `numode.spad`, `groebsol.spad`, `primelt.spad`, `algfact.spad` | Spad | BSD-3-Clause (GitHub API 2026-10-18) | port permitted | Kovacic's algorithm exists here, only free implementation found |
| Reduce packages `odesolve`, `crack`, `cde`, `cdiff`, `eds`, `excalc`, `symmetry`, `geometry`, `tri`, `wu`, `arnum`, `ranum`, `cali`, `gentran`, `redlog`, `hephys`, `atensor` | Lisp/RLISP | BSD (top-level LICENSE read: "Reduce License"; per-package headers not read) | port permitted once each package header is read | `excalc` = exterior calculus, `gentran` = Fortran/C code generator (Gates, "Gentran", ACM SIGSAM Bulletin 1985, doi:10.1145/1089411.1089415), `crack` = differential elimination |
| Maxima `ode2`, `contrib_ode`, `desolve`, `itensor`, `ctensor`, `atensor` | Lisp | GPL-2 (not re-verified; precision file marks it verify) | oracle only | doc page for ODEs not fetched (404 on `maxima_105.html`) |
| SymPy `combinatorics` (`perm_groups.py`, `fp_groups.py`, `coset_table.py`, `named_groups.py`), `diffgeom`, `liealgebras`, `stats`, `physics/units`, `physics/quantum`, `codegen`, `simplify/cse_main.py`, `utilities/codegen.py`, `printing/{c,fortran,julia,rust,octave}.py` | Python | BSD-3 text | port permitted | free counterparts of GroupTheory, DifferentialGeometry (partial), Statistics (symbolic RVs), Units, Physics (partial), CodeGeneration |
| GAP | C/GAP | GPL-2.0 (GitHub API) | oracle only | reference for permutation groups; `smallgrp` package licence field reads NOASSERTION: read before using its data |
| msolve | C | GPL-2.0 (GitHub API) | oracle only | Gröbner bases and real solving over Q, to compare triangular decompositions |
| Singular | C++ | NOASSERTION on GitHub (GPL-2/3 historically; read COPYING before use) | oracle only | `triang.lib` for triangular sets |
| LinBox, FFLAS-FFPACK | C++ | LGPL-2.1 (GitHub API) | oracle only | the free word-size-prime kernels; algorithm papers in 2(b) are portable knowledge |
| Cadabra2 | C++/Python | GPL-3.0 (GitHub API) | oracle only | tensor canonicalisation |
| SageManifolds | Python (Sage) | licence field empty on GitHub (Sage is GPL) | oracle only | differential geometry on manifolds |
| BLAD / Boulier's DifferentialAlgebra back end | C | not checked | unknown | verify licence before any reuse; algorithms in 2(c) are in published papers |
| OpenModelica, Modelica Standard Library | Modelica | OSMC-PL/GPL-3 (NOASSERTION on GitHub); MSL BSD-3-Clause (GitHub API) | MSL models usable as tests; OpenModelica oracle only | see systems file 3a |
| Go side | | | | no Go library for any row of section 1 was found (searched with the same method as the precision file; absence not proven). `go-codegen` is a third-party organisation (created 2023-08-09, one repo): do not use that name |

## 4. Validation

Maple cannot be an oracle on CI (proprietary, licence-locked, not installable without root and a
licence). Free oracles replace it and expected values are frozen into `testdata` by a Go generator,
recording tool and version.

### Oracles

| domain | oracle (free) | drive from Go |
|---|---|---|
| ODE solutions | SymPy `dsolve`/`classify_ode`/`checkodesol` (`uv run --with sympy`); FriCAS for Kovacic; Maxima `ode2` as a third opinion (verify availability) | exec + JSON on a pipe, as in the precision file |
| ODE correctness | no oracle needed for acceptance: substitute the solution into the ODE and check the residual is zero symbolically (the `odetest` idea, `help:dsolve/Lie`) or numerically at random points with big-float arithmetic | in-process |
| Algebraic numbers | PARI/GP (`nffactor`, `polredabs`), python-flint | as in the precision file |
| Triangular sets | msolve, Singular, Sage (all GPL: oracle only); check by radical membership and by solving numerically | exec |
| Differential algebra | no free complete oracle found; check Rosenfeld-Groebner output by Ritt reduction of the input equations and by comparing Taylor-series solutions | in-process |
| Groups | GAP, SymPy `combinatorics`; exact equality of group order and conjugacy-class sizes | exec |
| Code generation | SymPy `cse`, `ccode`, `fortran_code`; plus compile and run | `go build` on the output, compare with the interpreter |
| Modelling | Modelica Standard Library examples and OpenModelica results; Bari IVP test set (see systems file section 4) | exec |

### Public corpora

| corpus | URL | size | licence / caveat |
|---|---|---|---|
| Kamke's collection | book, Chelsea 1959 (see 2a); equations referred to by number | 1429 ODEs according to a Maplesoft comparison PDF (https://aws.webstore.maplesoft.com/compare/mathematica_analysis/Comparison_Maple_Mathmatica_DEs_Kamke.pdf; host not resolvable from here on 2026-10-18, figure taken from a search excerpt, **verify**) | book is copyrighted: store equation numbers and our own re-typed forms, derive our own reference solutions, copy nothing from Maple output |
| Polyanin, Zaitsev, Handbook of Exact Solutions for ODEs, 2nd ed., CRC 2003, doi:10.1201/9781420035339 ; EqWorld https://eqworld.ipmnet.ru/en/solutions/ode.htm | | | EqWorld reads "Copyright 2004-2017 Andrei D. Polyanin", no reuse terms: cite, do not copy |
| Maplesoft's Kamke comparison | link above | reports Maple 17.02 against Mathematica 9.0.1, 3 minutes per ODE, Maple about 97.5% of solvable equations, Mathematica about 80.3% (search excerpt, **verify**) | vendor-run, old versions: use as a sanity range for a solved-fraction target, not as evidence |
| SymPy ODE and PDE tests | https://github.com/sympy/sympy/tree/master/sympy/solvers/ode/tests | hundreds of cases | BSD-3 text |
| Rubi suite | see precision file | 72,000+ integrals | MIT: also validates `dsolve` outputs that reduce to quadratures |
| Modelica Standard Library examples | https://github.com/modelica/ModelicaStandardLibrary | | BSD-3-Clause |
| OEIS A000001 (number of groups of order n) | https://oeis.org/A000001 | | reuse terms: verify |

### What "correct" means

| area | criterion | measure |
|---|---|---|
| ODE/PDE solve | residual after substitution is 0 (symbolic) or below 1e-12 relative at 100 random points (numeric, big-float); the solution is general (right number of constants) | solved fraction per class; wrong answers must be 0; classification agreement with SymPy `classify_ode` and with `odeadvisor` classes named on the page |
| Kovacic | for each equation: Liouvillian solution found iff the oracle finds one; non-existence compared with FriCAS | exact agreement on the test set |
| Algebraic numbers | exact equality after canonical form; norm and minimal polynomial equal PARI's | pass/fail; independence check never skipped silently |
| Triangular decomposition | union of zero sets equals the input variety: check by radical membership of the input in the saturated ideal, and the reverse inclusion by sampling | pass/fail; component counts against msolve |
| CAD | sample points satisfy the formula in every cell; cells partition the space; sign-invariance of each input polynomial | cell count vs Maple's published examples is not usable (no oracle); compare with a QE oracle |
| Differential elimination | each returned chain reduces the input to 0 under Ritt reduction | pass/fail |
| Code generation | generated Go compiles with `go vet` clean; equals the expression's interpreter on 1000 random points within a stated ULP bound; operation count after CSE not above the SymPy `cse` count | ULP, op counts, compile success |
| Groups | exact: order, orbit sizes, class counts | pass/fail |
| DAE/MapleSim-like | index detected, constraint drift, final-time error against Modelica reference runs | as systems file section 4 |

## 5. Proposed home and plan

No new organisation. Every item fits an organisation that already exists on GitHub (checked with
`gh api users/<name>` on 2026-10-18; the thirteen organisations of PLAN.md are present and empty).
`go-codegen` is taken by a third party, hence code generation is a repository of `go-symbolic`.

| work item | org/repo | depends on | size S/M/L/XL | priority 1-3 |
|---|---|---|---|---|
| Expression-to-Go code generation: IR, CSE on the hash-consed DAG, Horner option, Go emitter with a fusion policy; C and Fortran emitters used as test oracles only; pluggable target interface (Maple's `LanguageDefinition` idea) | `go-symbolic/codegen` | expression DAG, `diff` | M | 1 within go-symbolic |
| Algebraic numbers: minimal-polynomial representation, dynamic evaluation for reducible cases, norm/trace/GCD/factor over number fields, primitive element | `go-symbolic/algnum` | poly GCD and factorisation | L | 2 |
| Word-size prime-field linear algebra with rank profile (FFLAS-style blocking) feeding the exact linear algebra item | `go-symbolic/linalg` (same repo as the precision file's exact linear algebra) | poly, go-ndarray BLAS-like kernels | M | 3 |
| ODE: class catalogue and recogniser, integrating factors, Lie symmetry heuristics, Kovacic, hypergeometric reductions; `checkodesol`-style verifier first | `go-symbolic/ode` | expr, poly, integrate, series | XL | 3 |
| PDE heuristics and system splitting | `go-symbolic/pde` | ode, diff | L | 3 |
| Regular chains, triangular decomposition, then CAD by triangular decomposition | `go-symbolic/triangular` | poly, algnum, Gröbner | XL | 3 |
| Differential algebra: rankings, Ritt reduction, Rosenfeld-Groebner | `go-symbolic/diffalg` | poly, triangular | L | 3 |
| Symbolic random variables (PDF/CDF/MGF/moments from distribution definitions) | `go-symbolic/stats`, numerics in `go-statistics/stats` | expr, integrate, series, special | M | 3 |
| Permutation groups (Schreier-Sims), conjugacy classes, Cayley tables, finitely presented groups; small-group data only after a licence read | `go-symbolic/group` | none | L | 3 |
| Abstract-index tensors with canonicalisation, general-relativity curvature tensors, noncommutative algebra, Lie algebra library | `go-symbolic/tensor` | group, expr | XL | 3 |
| User-defined unit systems (EMU, ESU, Atomic) | `go-units/units` | the units item | S | 3 |
| 2-D math input: typed or LaTeX input to a go-symbolic tree and back; embedded widgets beyond sliders in the Manipulate design | `go-notebook/notebook` with `go-tex/math` for rendering | expr, LaTeX printer | L | 3 |
| Symbolic acausal modelling: connectors, equation assembly, Pantelides/dummy derivatives on go-symbolic equations, BLT, emission of Go through `codegen`; Modelica-subset reader last | `go-control/model` | DAE solver in integrate or go-control, go-symbolic expr, codegen | XL | 3 |

## 6. Risks and open questions

- **The Maple help pages describe interfaces, rarely algorithms.** `evala`, `Triangularize` (one reference), `RosenfeldGroebner` (labels only), `codegen[optimize]` (no paper) and `dsolve` (pointer to a references page) give little. Attribution of an algorithm to Maple (e.g. dynamic evaluation in `evala`) is an inference from the page's wording and is flagged as such.
- **Maple is a black box with no public test list.** The only solved-fraction figures found are vendor-run (Maplesoft against Mathematica, 2013-era versions) and unverified here. A target such as "match Maple on Kamke" has no auditable denominator; define the target on our own corpus.
- **ODE solving is heuristic and open-ended.** Lie-symmetry ansatz search and class catalogues grow without bound; Kovacic is the only complete piece. Cost control (time limit per ODE, as the vendor comparison used) and a verifier-first design are needed so a wrong "solution" is impossible.
- **Licences.** SymPy, FriCAS, Reduce are portable (read per-file notices; Reduce per-package headers not read); GAP, msolve, LinBox, FFLAS-FFPACK, Cadabra, Maxima, Singular are oracle-only. Do not copy Maple help text beyond fair-use quotation or any Maple output into tests. Kamke, Polyanin and Murphy are books: equation numbers and re-derived solutions only.
- **CAD and triangular decomposition are doubly exponential in the worst case**; a Go port needs the same optimisations Maple names (equational constraints, TTICAD) to be usable, and polynomial arithmetic over Z/p with word-size fast paths first.
- **Generated Go floating point.** Fusion across statements is allowed by the spec; the codegen item must choose between wrapping every product (reproducible, slower) and leaving fusion (faster, platform-dependent results) and test on arm64, ppc64le and s390x.
- **MapleSim parity is a product, not a library.** MapleSim's component libraries are proprietary content; the vendor claims about symbolic index reduction come from marketing pages (the whitepaper on high-index DAEs in FMI was seen only as a search excerpt: https://www.maplesoft.com/whitepapers/Japan2018C.PDF). Realistic scope: symbolic assembly plus Pantelides on go-symbolic, validated on Modelica Standard Library examples.
- **Document model.** `.mw` round-tripping is not planned; the facts about 2-D math input come from search excerpts, not full pages (pages `worksheet/documenttools` returned 410). Whether a structured math editor belongs in go-notebook or in a separate front end is open.
- **Dependency chain.** Nearly everything here waits for the go-symbolic core (expression DAG, polynomials, GCD, factorisation), which is itself an XL item with no repository yet. Only `codegen` and the ODE verifier are useful soon after it.
- **Unverified in this file:** the Maxima licence and ODE documentation, BLAD licence, the Maple Kamke figures, OEIS reuse terms, Reduce per-package licences, MapleSim internals beyond vendor statements.

Verified on 2026-10-18: Maple command names listed in section 1 against the help pages cited there (names the pages did not list were dropped: `symtest`, `Classify`, `ParametricSystemSolve`, `SchreierSims`, `NumberOfSmallGroups` are not claimed); every DOI resolved through Crossref against title, year and first author (volume and page numbers not shown by Crossref are as printed on the Maple reference pages or marked by omission) (several first-guess DOIs for Singer 1981, Hubert, Lemaire, Dahan, van Hoeij-Monagan, Sims, Downey and Collins were wrong and were replaced by the Crossref results); licences and file lists of SymPy, FriCAS, Reduce, GAP, msolve, LinBox, FFLAS-FFPACK, Cadabra2, MSL via GitHub API; go-scikit state via GitHub main.
