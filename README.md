# roadmap

What the [go-scikit](https://go-scikit.github.io/) stack must add to reach feature parity with
MATLAB, Mathematica and Maple, where each piece lives, and how far it has got.

- [PLAN.md](PLAN.md) — the rules every work item follows, the home of each domain, and the
  order of the work in four waves.
- [STATUS.md](STATUS.md) — each item with its repository, pull requests, release and measured
  parity.
- [bib/](bib/) — one file per domain: the MATLAB and Mathematica inventory against what
  go-scikit has, the reference algorithms and papers, the reference implementations and their
  licences, and how results are validated. DOIs, function names and licences were checked on
  2026-10-17; each file ends with what was verified.
  - [linear algebra](bib/linear-algebra.md)
  - [approximation, signal and image processing](bib/approximation-signal-image.md)
  - [arbitrary precision, number theory, computer algebra](bib/precision-numbertheory-symbolic.md)
  - [optimisation, statistical learning, automatic differentiation](bib/optimization-learning-autodiff.md)
  - [geometry, PDEs and graphs](bib/geometry-pde-graphs.md)
  - [systems, units, file formats and visualisation](bib/systems-units-formats-visualization.md)
  - [code shared with the pure-Go Ruby](bib/ruby-reuse.md) — where go-embedded-ruby and the
    go-ruby-* gems duplicate go-scikit code, and which side should own each piece
  - [Maple](bib/maple.md) — only what Maple adds to the two inventories above: ODE/PDE solving,
    algebraic numbers, triangular and differential algebra, code generation to Go, tensors,
    groups, and MapleSim-style modelling
  - [Julia](bib/julia.md) — bibliography and feasibility plan for go-embedded-julia, a pure-Go
    Julia interpreter over lowered IR for go-notebook cells and the WebAssembly playground:
    parser, lowering, runtime boundary (94 intrinsics, 896 `ccall` sites in Base), milestones
  - [template](bib/TEMPLATE.md) for a new domain

`go run ./tools/checklinks` checks that every relative link in the tree resolves; CI runs it.
