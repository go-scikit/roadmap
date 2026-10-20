# Status

As of 2026-10-20. "Released" means a tagged version; "merged" means on main without a tag yet.

| # | item | home | state | evidence |
|---|---|---|---|---|
| 1 | correctly rounded big.Float functions, Go oracle harness | go-bignum/bigfloat | in progress, not yet on GitHub: error model, Ziv loop, roots, constants, exp/log done; trigonometric, hyperbolic and pow in progress | — |
| 2 | sparse containers, orderings, iterative solvers | go-sparse/sparse | in progress; performance pull request in review and being re-measured: [go-sparse/sparse#2](https://github.com/go-sparse/sparse/pull/2) | iteration counts equal to SciPy's for CG, GMRES, MINRES and LSQR |
| 3 | interpolation, splines, polynomials | [go-approx/interp](https://github.com/go-approx/interp), [go-approx/poly](https://github.com/go-approx/poly) | interp v0.1.0 released (speed comparison in [docs/perf.md](https://github.com/go-approx/interp/blob/main/docs/perf.md)); remaining interp gaps in progress; poly merged, performance pull request in review: [go-approx/poly#2](https://github.com/go-approx/poly/pull/2) | numpy.interp, interp1d, PCHIP, Akima and poly1d bit-identical to NumPy/SciPy; splines within 1e-15–1e-10 |
| 4 | filter design, filtering, spectral estimation | go-dsp/signal | in progress; performance pull request in review: [go-dsp/signal#2](https://github.com/go-dsp/signal/pull/2) | — |
| 5–22 | see [PLAN.md](PLAN.md) | | planned | |
| — | vector-update kernels `Axpy`, `Xpby`, `AxpyDot`, `Scale` (shared by the solvers) | [go-simd/floats](https://github.com/go-simd/floats) | v0.4.0 released: SIMD on amd64 and arm64, FMA loop elsewhere | — |
