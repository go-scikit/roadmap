# Coverage

Where the go-scikit stack stands against its references, and the order in which the gaps are closed. Goal: the whole mathematical surface of SciPy, NumPy, statsmodels, mpmath and the matching MATLAB and Mathematica functions, not a usable subset.

Each repository keeps its own `docs/coverage.md` (reference name, MATLAB name, Go name, status done / partial / missing, dependency). This page sums them up as of 2026-11-02. The MATLAB, Mathematica, FFTW and Arb columns of those files were written from knowledge, because the vendor function lists are client-rendered; each file marks them for re-check before a MATLAB-driven slice starts. The counts below are copied from the files and are not yet recomputed by a committed tool (open item).

## Counts

| repository | reference | done | partial | missing |
|---|---|---:|---:|---:|
| [go-numerics/optimize](https://github.com/go-numerics/optimize/blob/main/docs/coverage.md) | scipy.optimize 1.18.1 | 27 | 14 | 44 |
| [go-numerics/special](https://github.com/go-numerics/special/blob/main/docs/coverage.md) | scipy.special 1.18.1 | 33 | 1 | 89 |
| [go-numerics/integrate](https://github.com/go-numerics/integrate/blob/main/docs/coverage.md) | scipy.integrate 1.18.1 | 22 | 5 | 11 |
| [go-statistics/stats](https://github.com/go-statistics/stats/blob/main/docs/coverage.md) | scipy.stats 1.18.1 | 68 | 8 | 172 |
| [go-sparse/sparse](https://github.com/go-sparse/sparse/blob/main/docs/coverage.md) | scipy.sparse (+ linalg, csgraph) 1.18.1 | 19 | 21 | 65 |
| [go-approx/interp](https://github.com/go-approx/interp/blob/main/docs/coverage.md) | scipy.interpolate 1.18.1 | 11 | 14 | 31 |
| [go-approx/poly](https://github.com/go-approx/poly/blob/main/docs/coverage.md) | numpy.polynomial 2.5.3 | 12 | 12 | 28 |
| [go-dsp/signal](https://github.com/go-dsp/signal/blob/main/docs/coverage.md) | scipy.signal 1.18.1 | 65 | 18 | 49 |
| [go-fft/fft](https://github.com/go-fft/fft/blob/main/docs/coverage.md) | scipy.fft 1.18.1 | 10 | 17 | 14 |
| [go-ndarray/ndarray](https://github.com/go-ndarray/ndarray/blob/main/docs/coverage.md) | numpy 2.5.3 top level | 84 | 76 | 321 |
| | numpy.linalg | 18 | 7 | 6 |
| | scipy.linalg 1.18.1 | 7 | 13 | 85 |
| [go-bignum/bigfloat](https://github.com/go-bignum/bigfloat/blob/main/docs/coverage.md) | mpmath 1.4.1 functions | 42 | 1 | 238 |
| go-statistics/models | statsmodels 0.15.0 regression | design only ([design-models.md](https://github.com/go-statistics/stats/blob/main/docs/design-models.md)) | | |

## Order

### 1. Shared foundations

Most blocked slices wait on the same few building blocks. They come first, each in the lowest repository that can hold it.

| foundation | home | unblocks |
|---|---|---|
| complex number layer (elementary functions, then `loggamma`) | to decide: go-numerics/special or a new package | special: hypergeometric, Bessel extras, parabolic cylinder; bigfloat: almost all of mpmath; integrate: LSODA/VODE; fft: `fht` |
| QMC engines (Sobol, Halton, scrambling) | go-statistics/stats | stats.qmc, QMC integration, optimize `shgo` and DE `sobol`/`halton` init |
| `schur`, `hessenberg`, `lu`, `cho_*` exported, then `expm`, `logm`, `sqrtm` | go-ndarray | signal `lsim`, sparse `expm`, item 10 |
| tridiagonal (and banded) symmetric eigensolver | go-ndarray | orthogonal polynomials and Gauss rules (special, poly), integrate quadrature |
| triangular solve, sparse LU, `spsolve` | go-sparse | optimize MILP and conic, sparse Jacobians in integrate and optimize |
| Delaunay triangulation and k-d tree | go-geometry (item 7) | interp scattered data and RBF, optimize `shgo` |
| numerical Hessians; pivoted QR and public triangular solve | go-numerics/optimize; go-ndarray | go-statistics/models |

### 2. Slices with no dependency

Can run in parallel, one agent per repository.

- **special**: about 60 thin wrappers over existing functions; Cephes ports (expint, Fresnel, sici, spence, dawsn); Carlson and Jacobi elliptic.
- **stats**: about 60 closed-form continuous distributions; closed-form discrete and multivariate ones; about 50 tests and statistics; resampling.
- **optimize**: roots II (`toms748`, `fixed_point`, `root` lm/broyden/anderson/df-sane); trust-region family (dogleg, trust-ncg, trust-exact, SR1, `keep_feasible`, constraint Hessians); derivative-free (TNC, COBYLA, COBYQA); linear least squares (`nnls`, `lsq_linear`, `isotonic_regression`); global III (`dual_annealing`, MultiStart, GlobalSearch, `ga`, `particleswarm`, `patternsearch`).
- **sparse**: containers (LIL, DIA, BSR, `kron`, stacks, element-wise ops, axis reductions); csgraph; Krylov (BiCG, CGS, QMR, TFQMR, LGMRES, GCROT).
- **interp / poly**: `roots` and knot insertion; barycentric, Krogh, Floater–Hormann, Bernstein; grid cubic, pchip, makima; power basis and series arithmetic.
- **signal**: waveforms, missing windows; Savitzky–Golay, median and rank filters; residue and design leftovers; streaming filters.
- **fft**: hfft2/hfftn and `s`; `prev_fast_len`, `nextpow2`; DCT/DST `n`, batching, `orthogonalize=False`.
- **integrate**: `newton_cotes`, `cumulative_simpson`; `quad_vec`, `cubature`; `tanhsinh`, `nsum`; `solve_bvp`.
- **ndarray**: missing ufuncs; views and shape functions; reductions (`std`, `var`, `nan*`, multi-axis); missing random distributions.
- **bigfloat**: elementary completions (sec family, `cospi`/`sinpi`, constants); utilities (`floor`, `fmod`, `fsum`, `ulp`).

### 3. Slices behind a foundation

Started as soon as their foundation is released: special hypergeometric, Airy/Struve, orthogonal polynomials, Mathieu/spheroidal; bigfloat gamma, zeta, hypergeometric engine, elliptic, ball type, number theory; stats generic distribution layer, non-central and special-function distributions, Wishart, QMC, UNU.RAN; integrate LSODA/VODE, sparse Jacobians, DDE/PDE; optimize `shgo`, MILP and conic, Krylov roots, assignment, MATLAB SQP derivatives; sparse direct solvers, eigensolvers (`eigs`, `eigsh`, `lobpcg`, `svds`), factorizations; interp scattered data, RBF, FITPACK 2-D, AAA; signal state space, `ShortTimeFFT`/`CZT`, wavelets, N-D; ndarray matrix functions, banded, QZ, sort and sets, einsum, dtypes, `.npy`.

### 4. Plan items with their own design

go-statistics/models (item 6, waves W0–W12), then items 7–22 of [PLAN.md](PLAN.md). Left out of the item 6 design and to be planned as their own items: time series (statsmodels `tsa`), survival analysis (Cox, duration), imputation, empirical likelihood, `sqrt_lasso` (needs a conic solver), multivariate robust covariance, GMM and instrumental variables.

## Rules for every slice

Port the reference exactly and validate by equality where the reference is deterministic; measure against the reference on one core and treat any slower row as an open defect; use go-simd kernels without FMA and in the reference's reduction order where bit-parity is claimed; 100 % coverage, fuzzing, README, landing and docs for each release.
