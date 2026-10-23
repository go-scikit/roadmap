# Status

As of 2026-10-23. "Released" means a tagged version; "merged" means on main without a tag yet. Every child organisation has a landing, a documentation site, a brand and a profile; [go-scikit/docs](https://go-scikit.github.io/docs/latest/) links them all.

| # | item | home | state | evidence |
|---|---|---|---|---|
| 1 | correctly rounded big.Float functions, Go oracle harness | [go-bignum/bigfloat](https://github.com/go-bignum/bigfloat) ([docs](https://go-bignum.github.io/docs/latest/bigfloat/)) | v0.1.0 released; every function at or below mpmath's time ([docs/perf.md](https://github.com/go-bignum/bigfloat/blob/v0.1.0/docs/perf.md)) | correctly rounded against mpmath on 28 840 cases in six rounding modes |
| 2 | sparse containers, orderings, iterative solvers | [go-sparse/sparse](https://github.com/go-sparse/sparse) ([docs](https://go-sparse.github.io/docs/latest/sparse/)) | v0.1.0 released; every measured case at or below SciPy ([docs/perf.md](https://github.com/go-sparse/sparse/blob/main/docs/perf.md)) | iteration counts equal to SciPy's for CG, GMRES, MINRES and LSQR |
| 3 | interpolation, splines, polynomials | [go-approx/interp](https://github.com/go-approx/interp), [go-approx/poly](https://github.com/go-approx/poly) ([docs](https://go-approx.github.io/docs/latest/)) | interp v0.2.1 and poly v0.1.0 released, every measured case at or below NumPy/SciPy ([interp docs/perf.md](https://github.com/go-approx/interp/blob/main/docs/perf.md), [poly docs/perf.md](https://github.com/go-approx/poly/blob/main/docs/perf.md)) | numpy.interp, interp1d, PCHIP, Akima and poly1d bit-identical to NumPy/SciPy; splines within 1e-15–1e-10 |
| 4 | filter design, filtering, spectral estimation | [go-dsp/signal](https://github.com/go-dsp/signal) ([docs](https://go-dsp.github.io/docs/latest/signal/)) | v0.1.0 released; all 28 measured cases at or below SciPy ([docs/perf.md](https://github.com/go-dsp/signal/blob/main/docs/perf.md)) | 40 Remez designs within 3e-8 of SciPy; accuracy table over 225 groups in [docs/accuracy.md](https://github.com/go-dsp/signal/blob/main/docs/accuracy.md) |
| 5–22 | see [PLAN.md](PLAN.md) | | planned | |
| — | vector-update kernels `Axpy`, `Xpby`, `AxpyDot`, `Scale` (shared by the solvers) | [go-simd/floats](https://github.com/go-simd/floats) | v0.4.0 released: SIMD on amd64 and arm64, FMA loop elsewhere | — |
