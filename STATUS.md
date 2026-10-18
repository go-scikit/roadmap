# Status

As of 2026-10-18. "Released" means a tagged version; "merged" means on main without a tag yet.

| # | item | home | state | evidence |
|---|---|---|---|---|
| 1 | correctly rounded big.Float functions, Go oracle harness | go-bignum/bigfloat | in progress, not yet on GitHub: error model, Ziv loop, roots, constants, exp/log done; trigonometric, hyperbolic and pow in progress | — |
| 2 | sparse containers, orderings, iterative solvers | go-sparse/sparse | in progress, not yet on GitHub: formats, Matrix Market, CG/GMRES/BiCGSTAB/MINRES/LSQR with preconditioners, RCM/AMD/COLAMD | iteration counts equal to SciPy's for CG, GMRES, MINRES and LSQR |
| 3 | interpolation, splines, polynomials | [go-approx/interp](https://github.com/go-approx/interp), [go-approx/poly](https://github.com/go-approx/poly) | merged; not released until the speed comparison with SciPy is measured | numpy.interp, interp1d, PCHIP, Akima and poly1d bit-identical to NumPy/SciPy; splines within 1e-15–1e-10 |
| 4 | filter design, filtering, spectral estimation | go-dsp/signal | in progress, not yet on GitHub: IIR design (Butterworth, Chebyshev, elliptic, Bessel), conversions, frequency responses | — |
| 5–22 | see [PLAN.md](PLAN.md) | | planned | |
