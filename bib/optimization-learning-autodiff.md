# Optimisation, statistical learning and automatic differentiation — bibliography

One file per domain. Facts only, each with a source (URL or DOI). Terse: tables, not prose.

Provenance: sections 1 (go-scikit column), 3 licences and the Go-library rows were checked on
2026-10-17 (`gh api repos/...` on GitHub main, `gh search repos`). MATLAB/Mathematica
function names and all DOIs were re-checked on 2026-10-17 (vendor doc pages, Crossref).

## 1. Feature inventory

Sources: MATLAB https://www.mathworks.com/help/optim/ug/fmincon.html [doc: algorithms interior-point (default), trust-region-reflective, sqp, sqp-legacy, active-set];
https://www.mathworks.com/help/optim/ , https://www.mathworks.com/help/gads/ , https://www.mathworks.com/help/stats/ , https://www.mathworks.com/help/econ/ , https://www.mathworks.com/help/deeplearning/ ;
Mathematica https://reference.wolfram.com/language/tutorial/ConstrainedOptimizationOverview.html [doc: LinearOptimization, FindMinimum, NMinimize, ConicOptimization, IPOPTLink; methods InteriorPoint, DifferentialEvolution, NelderMead, SimulatedAnnealing, RandomSearch].

go-scikit today (verified on GitHub main): `go-numerics/optimize` = Minimize (Nelder-Mead, Powell, CG, BFGS, Newton-CG, L-BFGS-B; `Bounds`), MinimizeScalar, Fminbound, Brent*/Bisect/Newton/Ridder, Root/Fsolve, Leastsq, LeastSquares (TRF, LM), CurveFit (Pcov/Perr), Brute, Linprog (revised simplex and interior-point/HSD, presolve, LPDual), Bracket, ApproxFprime, CheckGrad. `go-numerics/integrate` = Quad/NQuad, Romberg, SolveIVP (RK45, Radau, BDF, LSODA…), Odeint, takes `JacFunc` (user-supplied). `go-statistics/stats` = distributions + Fit*, tests, descriptive stats, KDE, Bootstrap, BinnedStatistic, FitMultivariateNormal; regression = Linregress only. `go-dataframe/dataframe`, `go-ndarray` exist. All BSD-3-Clause.

### (a) Optimisation

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| unconstrained local | fminunc, fminsearch | FindMinimum, FindMaximum | optimize.Minimize | none (check quasi-Newton/trust-region set vs SciPy) |
| constrained nonlinear (SQP, IP, trust-region) | fmincon (5 algorithms) | FindMinimum/NMinimize with constraints, Method InteriorPoint, IPOPTLink | none (Minimize has bounds only: L-BFGS-B, Nelder-Mead, Powell) | SLSQP, trust-constr, interior point, augmented Lagrangian |
| QP | quadprog | QuadraticOptimization | none | active set, ADMM, interior point |
| LP | linprog (dual-simplex-highs, interior-point) | LinearOptimization | Linprog: revised simplex + interior-point (dense), presolve | dual simplex, sparse |
| MILP / MIQP | intlinprog, (ga with integers) | LinearOptimization/NMinimize with Integers | none | branch and cut |
| global | ga, patternsearch, particleswarm, simulannealbnd, surrogateopt, GlobalSearch, MultiStart | NMinimize (DifferentialEvolution, SimulatedAnnealing, RandomSearch, NelderMead) | Brute only | DE, basin hopping, annealing, SHGO, DIRECT, CMA-ES |
| multi-objective | gamultiobj, paretosearch, fgoalattain, fminimax | (via NMinimize weighting; no dedicated) | none | NSGA-II, minimax |
| SDP / conic | (Optimization Toolbox: no SDP; via CVX/YALMIP) | SemidefiniteOptimization, ConicOptimization | none | SCS-class ADMM, ECOS-class IP |
| nonlinear least squares | lsqnonlin, lsqcurvefit, lsqlin, lsqnonneg | FindFit, NonlinearModelFit | Leastsq, LeastSquares, CurveFit | lsqlin, nnls (absent) |

### (b) Statistical models and machine learning

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| linear model + inference | fitlm, regress, anova, stepwiselm | LinearModelFit (SE, ANOVA table, CI) | stats.Linregress (1 regressor) | multiple regression, SE, t/F, ANOVA, diagnostics, robust, ridge/lasso |
| GLM | fitglm, glmfit | GeneralizedLinearModelFit | none | IRLS, families, links, deviance |
| nonlinear regression | fitnlm | NonlinearModelFit | CurveFit (Pcov/Perr = parameter SE) | CI, t-tests, residual diagnostics |
| mixed models | fitlme, fitlmematrix | none native | none | REML/ML LMM |
| nonparametric regression | fitrgp, fitrkernel, smoothing (lowess via Curve Fitting) | SmoothKernelDistribution, Predict | KDE only | LOWESS, splines, GP |
| time series | arima (estimate, forecast methods), ssm (filter, smooth), exponential smoothing (Econometrics) | TimeSeriesModelFit, ARMAProcess, StateSpaceModel | none | ARIMA, SARIMAX, state space/Kalman, ETS |
| PCA / FA | pca, factoran, rotatefactors | PrincipalComponents, no FA | none | PCA, FA, rotations |
| clustering | kmeans, linkage/cluster, dbscan, fitgmdist | FindClusters, ClusterClassify | none | k-means++, hierarchical, DBSCAN, GMM-EM |
| classification | fitclinear, fitcsvm, fitctree, TreeBagger, fitcensemble, fitcknn | Classify, Predict | none | logistic, SVM, trees/forests, boosting, k-NN |
| model selection | crossval, cvpartition, hyperparameter opt | ClassifierMeasurements, PredictorMeasurements (no CrossValidate) | Bootstrap | CV splitters, grid search, metrics |
| neural networks | Deep Learning Toolbox (trainnet, dlnetwork) | NetChain/NetGraph/NetTrain | none | out of first scope; needs AD |

### (c) Automatic differentiation

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| symbolic/AD derivatives | Symbolic Math Toolbox diff/jacobian; dlgradient/dlarray AD (Deep Learning Toolbox); Optimization problem-based workflow uses AD | D, Dt (symbolic); FindMinimum auto-derives | ApproxFprime (finite differences), CheckGrad | forward and reverse AD |
| Jacobian for stiff ODE | odeset Jacobian / JPattern | NDSolve symbolic Jacobian | integrate.JacFunc user-supplied | AD-built Jacobian feeding JacFunc |
| Hessian for optimise | fmincon HessianFcn, quasi-Newton | FindMinimum Hessian | none | forward-over-reverse |

## 2. Reference algorithms

| capability | algorithm / reference | why |
|---|---|---|
| constrained NLP, large | IPOPT primal-dual filter line-search IP: Wächter & Biegler 2006, Math. Program. 106:25–57, doi:10.1007/s10107-004-0559-y | default of fmincon-class; handles sparse, inequality |
| constrained NLP, small dense | SLSQP: Kraft 1988, DFVLR-FB 88-28 (no DOI; report); SQP theory Nocedal & Wright 2006 *Numerical Optimization* 2nd ed., doi:10.1007/978-0-387-40065-5 | SciPy `SLSQP`, cheap per iteration |
| trust-region constrained | Byrd, Hribar, Nocedal 1999, SIAM J. Optim. 9:877, doi:10.1137/S1052623497325107 (basis of SciPy `trust-constr` / KNITRO) | robust without exact Hessian |
| augmented Lagrangian | Conn, Gould, Toint 1991 (LANCELOT), SIAM J. Numer. Anal. 28:545, doi:10.1137/0728030; Andreani, Birgin, Martínez, Schuverdt 2008 (ALGENCAN), doi:10.1137/060654797 | simple, scales, tolerates degenerate constraints |
| QP active set (dense, convex) | Goldfarb & Idnani 1983, Math. Program. 27:1, doi:10.1007/BF02591962 | dual method, stable, as `quadprog` R/Python packages |
| QP first-order | OSQP: Stellato, Banjac, Goulart, Bemporad, Boyd 2020, Math. Prog. Comp. 12:637, doi:10.1007/s12532-020-00179-2 | ADMM, one factorisation, warm start; Apache-2.0 |
| LP/QP interior point | Mehrotra predictor-corrector 1992, SIAM J. Optim. 2:575, doi:10.1137/0802028 | standard IPM for LP/QP |
| LP at scale | HiGHS dual revised simplex: Huangfu & Hall 2018 (Parallelizing the dual revised simplex method), Math. Prog. Comp. 10:119, doi:10.1007/s12532-017-0130-5 | MATLAB `linprog` algorithm `dual-simplex-highs`; SciPy default |
| MILP | branch and cut: Achterberg 2007 PhD thesis (SCIP), doi:10.14279/depositonce-1634; Wolsey *Integer Programming* 2nd ed. 2020, doi:10.1002/9781119606475 | needed for intlinprog parity |
| differential evolution | Storn & Price 1997, J. Glob. Optim. 11:341, doi:10.1023/A:1008202821328 | Mathematica NMinimize default family, SciPy |
| basin hopping | Wales & Doye 1997, J. Phys. Chem. A 101:5111, doi:10.1021/jp970984n | SciPy `basinhopping` |
| simulated annealing | Kirkpatrick, Gelatt, Vecchi 1983, Science 220:671, doi:10.1126/science.220.4598.671 | `simulannealbnd` |
| SHGO | Endres, Sandrock, Focke 2018, J. Glob. Optim. 72:181, doi:10.1007/s10898-018-0645-y | SciPy `shgo` |
| DIRECT | Jones, Perttunen, Stuckman 1993, JOTA 79:157, doi:10.1007/BF00941892 | deterministic global, bounds only |
| CMA-ES | Hansen & Ostermeier 2001, Evol. Comput. 9:159, doi:10.1162/106365601750190398 | gonum already has `cmaes.go` (BSD-3) |
| multi-objective | NSGA-II: Deb et al. 2002, IEEE TEVC 6:182, doi:10.1109/4235.996017 | `gamultiobj` basis |
| conic ADMM | SCS: O'Donoghue, Chu, Parikh, Boyd 2016, JOTA 169:1042, doi:10.1007/s10957-016-0892-3 | LP/QP/SOCP/SDP/exp cone |
| conic IPM | ECOS: Domahidi, Chu, Boyd 2013, ECC, doi:10.23919/ECC.2013.6669541 | SOCP, small, accurate |
| linear/GLM | IRLS: McCullagh & Nelder 1989 *Generalized Linear Models* 2nd ed., ISBN 0-412-31760-5; QR/SVD least squares: Golub & Van Loan 4th ed. | exact SE from (X'WX)^-1 |
| statistical learning | Hastie, Tibshirani, Friedman 2009 ESL 2nd ed., doi:10.1007/978-0-387-84858-7 (free PDF https://hastie.su.domains/ElemStatLearn/) | trees, boosting, SVM, CV, regularisation |
| mixed models | Bates, Mächler, Bolker, Walker 2015, J. Stat. Softw. 67(1), doi:10.18637/jss.v067.i01 | profiled REML, sparse Cholesky |
| ARIMA | Box, Jenkins, Reinsel, Ljung *Time Series Analysis* 5th ed. 2015, ISBN 978-1-118-67502-1; exact likelihood via Kalman (Gardner, Harvey, Phillips 1980, Appl. Stat. 29:311, doi:10.2307/2346910) | exact ML, as statsmodels |
| state space | Durbin & Koopman 2012 *Time Series Analysis by State Space Methods* 2nd ed., doi:10.1093/acprof:oso/9780199641178.001.0001 | Kalman, smoother, diffuse init |
| exponential smoothing | Hyndman, Koehler, Snyder, Grose 2002, Int. J. Forecasting 18:439, doi:10.1016/S0169-2070(01)00110-8 | ETS state-space form |
| k-means | Lloyd 1982, IEEE Trans. IT 28:129, doi:10.1109/TIT.1982.1056489; k-means++ Arthur & Vassilvitskii 2007, SODA, ISBN 978-0-89871-624-5 | seeding bounds |
| DBSCAN | Ester, Kriegel, Sander, Xu 1996, KDD, https://www.aaai.org/Papers/KDD/1996/KDD96-037.pdf | density, no k |
| SVM | Cortes & Vapnik 1995, doi:10.1007/BF00994018; LIBSVM Chang & Lin 2011, doi:10.1145/1961189.1961199; LIBLINEAR Fan et al. 2008, JMLR 9:1871 | SMO + dual coordinate descent |
| forests / boosting | Breiman 2001, doi:10.1023/A:1010933404324; Friedman 2001, Ann. Stat. 29:1189, doi:10.1214/aos/1013203451; XGBoost Chen & Guestrin 2016, doi:10.1145/2939672.2939785 | |
| AD | Griewank & Walther 2008 *Evaluating Derivatives* 2nd ed., doi:10.1137/1.9780898717761; Baydin, Pearlmutter, Radul, Siskind 2018, JMLR 18(153), https://jmlr.org/papers/v18/17-468.html (arXiv:1502.05767) | forward = n sweeps/Jacobian-vector, reverse = 1 sweep/gradient; forward-over-reverse Hessian-vector |

AD in Go (design facts): the Go spec has no operator overloading (https://go.dev/ref/spec), so AD is one of: (1) dual-number / hyper-dual types with explicit methods (`Add`, `Mul`) — forward mode, easy, verbose; (2) tape (Wengert list) built by calling tape methods — reverse mode, runtime graph; (3) source transformation with `go/parser`, `go/ast`, `go/types` (https://pkg.go.dev/go/ast) rewriting a function into value+derivative code — no runtime overhead, needs restricted subset; (4) generics (https://go.dev/ref/spec#Type_parameters) can let one numeric kernel run on `float64` and on a dual type only if operations are method-constrained, not `+`, so kernels need an interface of ops. Reference designs: Tapenade (source transformation, MIT per https://tapenade.gitlabpages.inria.fr/tapenade/distrib/LICENSE.html), Enzyme (LLVM, Apache-2.0 with LLVM exception), JAX (Apache-2.0).

## 3. Reference implementations

Licences read from the GitHub API on 2026-10-17 unless noted.

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| HiGHS (ERGO-Code/HiGHS) | C++ | MIT | port (large) or oracle | LP simplex/IPM, QP, MIP; https://highs.dev |
| IPOPT (coin-or/Ipopt) | C++ | EPL-2.0 | weak copyleft: oracle; clean-room port from paper | needs sparse linear solver (MUMPS/MA27) |
| OSQP (osqp/osqp) | C | Apache-2.0 | port | small core, QDLDL factorisation (Apache-2.0) |
| SCS (cvxgrp/scs) | C | MIT | port | cone projections incl. PSD need eigen-decomposition |
| ECOS (embotech/ecos) | C | GPL-3.0 | oracle only | |
| SciPy optimize (scipy/scipy) | Python/C/Fortran | BSD-3-Clause | port; reference for SLSQP, trust-constr, DE, basinhopping, shgo, dual_annealing, DIRECT | go-numerics/optimize already follows it |
| NLopt (stevengj/nlopt) | C | mixed (API: NOASSERTION); COPYING: conjunction of sub-licences, `luksan` is LGPL-2.1+, so the combined library is LGPL | oracle | SLSQP, COBYLA, ISRES |
| OR-Tools (google/or-tools) | C++ | Apache-2.0 | oracle (CP-SAT, GLOP) | |
| gonum/optimize (+ convex/lp) | Go | BSD-3-Clause | reuse or port | files: bfgs, lbfgs, cg, cmaes, neldermead, newton, gradientdescent; `convex/lp` only: no SQP, no QP, no global |
| gonum/stat | Go | BSD-3-Clause | reuse or port | pca_cca.go, roc.go, distmv, sampleuv; no GLM, no regression-with-inference |
| statsmodels (statsmodels/statsmodels) | Python | BSD-3-Clause | port; reference for OLS/GLM/ARIMA/state space output | https://www.statsmodels.org |
| scikit-learn (scikit-learn/scikit-learn) | Python/Cython | BSD-3-Clause | port | k-means, DBSCAN, trees, forests, GMM |
| liblinear, libsvm (cjlin1) | C/C++ | BSD-3-Clause | port | SVM, logistic |
| lme4 | R/C++ | GPL-2+ | oracle only | |
| R stats (arima, glm, lm) | R | GPL-2+ | oracle only | `Rscript`; reference numbers for ARIMA |
| golearn (sjwhitworth/golearn) | Go | MIT | inspect only, last push 2024-01 | k-NN, trees, linear |
| gorgonia (gorgonia/gorgonia) | Go | Apache-2.0 | reuse ideas (graph-based AD on tensors) | last push 2024-08 |
| itsubaki/autograd | Go | MIT | reference for tape AD | active 2026-10 |
| pointlander/gradient | Go | BSD-3-Clause | reference | |
| pbenner/autodiff | Go | GPL-3.0 | oracle/read-only | licence blocks reuse |
| tamnd/dual, meirizarrygelpi/dual | Go | MIT | inspect | stale (2016, 2018); dual-number algebras |
| malaschitz/randomForest | Go | Apache-2.0 | reuse candidate | 2025-11 |
| goarima (albertyw, sartorproj) | Go | MIT | inspect | quality unverified |
| CppAD / Tapenade / Enzyme / JAX | C++/Java/LLVM/Python | EPL-2.0, secondary GPL-2.0+ (CppAD COPYING; API: NOASSERTION) / MIT (Tapenade LICENSE page) / Apache-2.0 with LLVM exceptions (Enzyme LICENSE) / Apache-2.0 | CppAD: oracle; Tapenade, Enzyme, JAX: port or oracle | derivative oracles for tests |
| MATLAB, Mathematica | — | proprietary | oracle only if the user owns a licence | not required in CI |

## 4. Validation

Oracles (drive from a Go harness by exec, JSON on stdin/stdout, as for the other go-scikit domains):

| oracle | obtain without root | drive |
|---|---|---|
| SciPy (optimize, stats) | `uv run --with scipy python …` (macOS arm64 and Linux wheels) | JSON in, JSON out |
| statsmodels, scikit-learn | `uv run --with statsmodels --with scikit-learn …` | same |
| HiGHS CLI / highspy | release binary from https://github.com/ERGO-Code/HiGHS/releases or `uv run --with highspy` | MPS file in, solution file out |
| IPOPT | `uv run --with cyipopt` needs system IPOPT; use conda-forge/pkgx if present | oracle only |
| OSQP / SCS / ECOS | `uv run --with osqp --with scs --with ecos` | JSON |
| R | pkgx `r` or system; `Rscript -e` | CSV/JSON |
| JAX / autograd | `uv run --with jax` | Jacobians/Hessians as reference values |

Public corpora:

| corpus | URL | use | licence / note |
|---|---|---|---|
| CUTEst | doi:10.1007/s10589-014-9687-3, https://github.com/ralna/CUTEst | NLP | LGPL-3.0; SIF files, needs decoder: oracle side only |
| Hock–Schittkowski | Hock & Schittkowski 1981, Lecture Notes in Econ. & Math. Systems 187, doi:10.1007/978-3-642-48320-2 | small NLP with known optimum | problem statements are facts |
| Netlib LP | https://www.netlib.org/lp/data/ | LP, MPS format, known objectives | public, some files restricted notes in README |
| Maros–Meszaros QP | Maros & Mészáros 1999, Optim. Methods Softw. 11:671, doi:10.1080/10556789908805768; https://github.com/qpsolvers/maros_meszaros_qpbenchmark | 138 convex QPs | |
| MIPLIB 2017 | https://miplib.zib.de, doi:10.1007/s12532-020-00194-3 | MILP, "benchmark" set 240 | per-instance licences vary: check |
| Mittelmann benchmarks | https://plato.asu.edu/bench.html | LP/QP/SDP solver comparison | reference timings |
| SDPLIB | http://euler.nmt.edu/~brian/sdplib/sdplib.html | SDP | Borchers 1999 (SDPLIB 1.2), doi:10.1080/10556789908805769 |
| Global test functions | Surjanovic & Bingham https://www.sfu.ca/~ssurjano/optimization.html | DE/annealing/DIRECT | |
| NIST StRD | https://www.itl.nist.gov/div898/strd/ (linear: Norris, Longley, Filip, Wampler1–5; nonlinear: higher/average/lower difficulty) | certified values to 11 digits for OLS and NLS | US government work |
| UCI ML Repository | https://archive.ics.uci.edu | iris, wine, digits… for classification/clustering | per-dataset CC BY 4.0 mostly; check |
| M3/M4 forecasting | https://github.com/Mcompetitions/M4-methods | time series | repo declares no licence |
| Box–Jenkins series (airline, sunspots) | bundled in statsmodels/R datasets | ARIMA | |
| AD identities | analytic derivatives, plus finite-difference `CheckGrad` and JAX oracle | AD | |

What "correct" means:

| domain | criterion | reported measure |
|---|---|---|
| LP/QP/NLP | KKT residuals (primal, dual, complementarity) below tol; objective vs reference | relative objective gap, max KKT residual; solved count over corpus |
| MILP | status + objective within 1e-6 relative; proven optimal gap | solved/timeouts, objective gap |
| global heuristics | statistical: success rate over N seeds to known optimum | success rate, median evaluations |
| OLS/NLS vs NIST | certified digits | LRE = −log10(rel. error) per coefficient and SE; target ≥ 9 on easy, report Filip/Wampler5 honestly |
| GLM/ARIMA | coefficients, SE, log-likelihood, AIC vs statsmodels/R | relative error, ULP for deterministic parts |
| clustering/classification | deterministic given seed: partitions exact match vs sklearn with same init; metrics (ARI, accuracy) | ARI = 1 on shared init; accuracy within CV noise |
| AD | exact to rounding for forward/reverse vs analytic; vs FD to O(h) | max relative error, ULPs |

## 5. Proposed home and plan

Existing orgs: `go-numerics` (optimize, integrate), `go-statistics` (stats). New orgs below were checked with `gh api users/<name>` (404 = free) on 2026-10-17: go-optim, go-mlkit, go-autodiff, go-timeseries, go-gradient, go-solvers, go-modeling, go-mathopt, go-fitkit, go-regress, go-gradmath, go-tapes are free; `go-learn` is TAKEN.

| work item | org/repo | depends on | size | priority |
|---|---|---|---|---|
| Forward + reverse AD library (dual/hyper-dual, tape, Jacobian/Hessian helpers); feeds `optimize` gradients and `integrate.JacFunc` | NEW org `go-autodiff`, repo `ad` | go-ndarray | M | 1 |
| Multiple linear regression + GLM (IRLS), SE/t/F/ANOVA, diagnostics; certify on NIST StRD | go-statistics (new repo `models`) | stats, go-ndarray, a QR/SVD (gonum port or own) | M | 1 |
| Constrained NLP: SLSQP, then interior point/augmented Lagrangian; QP (active set, then OSQP port) | go-numerics/optimize extension or NEW org `go-optim` repo `constrained` | optimize, linear algebra, AD | L | 1 |
| LP at scale: sparse dual simplex + presolve + IPM; then branch and cut MILP | NEW org `go-optim` repo `lp`/`mip` (HiGHS MIT as design reference) | sparse LU | XL | 2 |
| Global optimisation: DE, basin hopping, dual annealing, SHGO, DIRECT, CMA-ES (gonum reference), NSGA-II | go-numerics/optimize or `go-optim/global` | optimize | M | 2 |
| Clustering, PCA/FA, SVM (liblinear/libsvm ports), trees/forests/boosting, CV and metrics | NEW org `go-mlkit` | stats, dataframe, ndarray | XL (split per algorithm) | 2 |
| Time series: ARIMA/SARIMAX, state space/Kalman, ETS, mixed models | NEW org `go-timeseries` (mixed models under go-statistics) | models, optimize, AD | L | 2 |
| Conic: SCS/ECOS-class SOCP then SDP | `go-optim/conic` | eigen-solver, QP | XL | 3 |
| Neural networks (NetTrain-class) | defer; `go-mlkit/nn` on top of go-autodiff reverse mode | AD, ndarray | XL | 3 |
| Validation harness (exec + JSON oracles, NIST/Netlib/Maros–Meszaros loaders) | each repo's `internal/oracle` | none | S per repo | 1 |

## 6. Risks and open questions

- Linear algebra is the bottleneck: constrained NLP, QP, GLM, SDP all need QR, Cholesky, sparse LU/LDL and symmetric eigen in pure Go, tested to LAPACK accuracy. Check what go-ndarray/gonum provide before sizing anything; none of this was inventoried here.
- IPOPT is EPL-2.0 and needs MA27/MUMPS-class sparse solvers: a clean-room IPM from the paper is a large effort; start with SLSQP/OSQP.
- HiGHS-class dual simplex and MILP are research-grade work (years in the reference). Scope "MIPLIB benchmark parity" must not be promised.
- AD in Go: no operator overloading means user code must be written against an AD type or be source-transformed; retrofitting existing `func([]float64) float64` objectives is only possible by go/ast rewriting of a restricted subset.
- Evidence gaps: Wolfram `IPOPTLink` appears in the optimization tutorial but has no `ref/` page; Hock–Schittkowski and Achterberg DOIs resolve (Springer / TU Berlin) but Crossref has no metadata for Achterberg's thesis; Tapenade MIT is from its distribution LICENSE page, source availability not checked.
- Corpus licences (MIPLIB instances, UCI datasets, CUTEst LGPL) must be checked per file before vendoring; prefer download-on-demand into a directory outside the repositories.
- Name `go-learn` is taken; avoid Python-derived names (e.g. scikit-style) per instruction.

Verified on 2026-10-17: all DOIs via Crossref (title/year/first author; 2 fixed: Gardner-Harvey-Phillips, SDPLIB; Achterberg resolves via doi.org); MATLAB and Mathematica names via vendor pages (CrossValidate and `estimate`/`forecast` as standalone names corrected); licences via `gh api` and LICENSE/COPYING files; go-numerics/optimize, integrate, go-statistics/stats, gonum contents via GitHub API; free org names via `gh api users/`.
