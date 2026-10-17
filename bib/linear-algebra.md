# Linear algebra (sparse, matrix functions, structured, tensor, GPU) — bibliography

One file per domain. Facts only, each with a source (URL or DOI). Terse: tables, not prose.

Scope: linear algebra beyond dense BLAS-level products. Inventory date 2026-10-17. DOIs below were
resolved against the Crossref API (title, year, first author matched). Unmarked = verified.

## 1. Feature inventory

Sources: MATLAB sparse index https://www.mathworks.com/help/matlab/sparse-matrices.html ;
expm page https://www.mathworks.com/help/matlab/ref/expm.html ;
Mathematica https://reference.wolfram.com/language/guide/SparseArrays.html ,
https://reference.wolfram.com/language/ref/LinearSolve.html ,
https://reference.wolfram.com/language/ref/Eigenvalues.html ,
https://reference.wolfram.com/language/ref/MatrixFunction.html ;
SciPy https://docs.scipy.org/doc/scipy/reference/sparse.linalg.html .

go-scikit today: `go-ndarray/ndarray` on GitHub (v0.12.0, checked by the coordinator through the
GitHub API on 2026-10-17) has `(*Array).MatMul`, `Dot`, `Inner`, `Outer` over a packed GEMM kernel
and a `linalg` package: LU, QR, Cholesky, SVD (divide and conquer), eig, eigh, lstsq, norm,
matrix power (`lu.go`, `qr.go`, `cholesky.go`, `svd.go`, `bdc.go`, `dc.go`, `eig.go`, `eigh.go`,
`lstsq.go`, `norm.go`, `power.go`) and the public functions `Solve`, `Inv`, `Det`, `Slogdet`, `Eig`,
`Eigvals`, `Eigh`, `Eigvalsh`, `Lstsq`, `Pinv`, `MatrixRank`, `Cond` (norm-based, not an estimator),
`Norm`, `MatrixPower`, `MultiDot`, `QR`, `Cholesky`, `SVD`, `SVDVals`. An earlier draft of this file
read a STALE local clone (1a4cfae) and wrongly concluded those were missing. Still missing (no file
or symbol in the GitHub tree): Schur form with reordering, QZ, kron, einsum, 1-norm
condition estimators (`rcond`/`condest`). `go-numerics/optimize` contains Cholesky/least-squares
code inside lbfgsb, linprog, curvefit (private helpers, not a public API). No sparse type exists
on GitHub main (tree paths of go-ndarray/ndarray, go-statistics/stats, go-dataframe/dataframe,
go-numerics/{special,optimize,integrate} searched for sparse/csr/csc: only `dcsrch.go`, a line search).

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| sparse containers | `sparse`, `spalloc`, `spdiags`, `speye`, `sprand`, `sprandn`, `sprandsym`, `spconvert`, `find`, `nnz`, `nonzeros`, `spfun`, `spones` | `SparseArray`, `Normal`, `ArrayRules`, `CoefficientArrays` | none | CSR, CSC, COO, DOK + conversions, Matrix Market I/O |
| sparse I/O | `spconvert`, `.mat` | import MTX, HarwellBoeing, MAT | none | Matrix Market, Rutherford-Boeing readers |
| fill-reducing ordering | `amd`, `colamd`, `symamd`, `symrcm`, `dissect`, `colperm`, `dmperm` | internal to `LinearSolve` (not documented as functions) | none | AMD, COLAMD, RCM, nested dissection |
| structural analysis | `sprank`, `etree`, `symbfact`, `spaugment` | none found | none | etree, symbolic factorisation, Dulmage-Mendelsohn |
| sparse direct solve | `mldivide`/`decomposition`, `lu`, `chol`, `qr` on sparse (SuiteSparse) | `LinearSolve` Method `"Multifrontal"`, `"Pardiso"`, `"Mumps"`, `"Banded"`, `"Cholesky"` | none (dense LU/Cholesky/QR exist in ndarray/linalg) | sparse LU, Cholesky/LDL', QR |
| iterative solvers | `pcg`, `gmres`, `bicg`, `bicgstab`, `bicgstabl`, `cgs`, `qmr`, `tfqmr`, `minres`, `symmlq`, `lsqr` | `LinearSolve` Method `"Krylov"` with `"ConjugateGradient"`, `"GMRES"`, `"BiCGSTAB"` | none | CG, GMRES, BiCGSTAB, MINRES, LSQR/LSMR over a LinearOperator |
| preconditioners | `ichol`, `ilu`, `equilibrate` | `"ILU0"`, `"ILUT"`, `"ILUTP"` (suboptions FillIn, Tolerance) | none | IC0/ICT, ILU0/ILUT, Jacobi/SSOR, AMG (not in MATLAB base; in PETSc/hypre/pyamg) |
| sparse eigen | `eigs` (ARPACK) | `Eigenvalues` Method `"Arnoldi"` (Criteria, Shift), `"FEAST"` (Hermitian), `"Banded"` | none | Lanczos/IRLM, Krylov-Schur, shift-invert, LOBPCG |
| sparse SVD | `svds` | `SingularValueList` (sparse supported) | none | Lanczos bidiagonalisation (IRLBA/PROPACK style) |
| norm / condition estimate | `normest`, `condest`, `cond`, `rcond` | `Norm`, `LUDecomposition` (third result: condition number estimate); no `MatrixConditionNumber` function | `Cond` (norm-based, via SVD), `Norm` | Hager/Higham 1-norm estimator, `onenormest`, `rcond` from LU |
| matrix exponential | `expm` (scaling and squaring, Higham 2005 / Al-Mohy-Higham 2009) | `MatrixExp` | none | `expm`, `expm_multiply` (action), Frechet derivative |
| matrix log / sqrt / power | `logm`, `sqrtm`, `mpower` | `MatrixLog`, `MatrixPower` | none | `logm`, `sqrtm`, fractional `matrix_power` |
| general matrix function | `funm` | `MatrixFunction` Method `"Schur"` (block Parlett), `"Jordan"`, `"SpectralBasis"`; accepts `SparseArray`, dense result | none | Schur-Parlett `funm` |
| Kronecker | `kron` | `KroneckerProduct` | none | `Kron`, structured `kron` operator (matrix-free) |
| Toeplitz / banded / circulant | `toeplitz`, `hankel`, `gallery`, `spdiags` | `ToeplitzMatrix`, `HankelMatrix`, `Banded` method (Eigenvalues, LinearSolve) | none | Toeplitz solvers (Levinson, Gohberg-Semencul), banded LU/Cholesky, circulant via FFT |
| tensor ops | `pagemtimes`, `tensorprod` | `TensorContract`, `TensorProduct`, `ArrayFlatten` | `Outer`, `Dot`, `Inner` only | `einsum`, `tensordot`, contraction-path optimiser, Tucker/CP (Toolbox-level, low priority) |
| GPU | `gpuArray` (Parallel Computing Toolbox, NVIDIA CUDA only); sparse gpuArray supported with limits: no logical sparse, contiguous-submatrix indexing, no index assignment, element and nnz counts below intmax('int32'); single-precision sparse since R2025a. Per-function support list not extracted — https://www.mathworks.com/help/parallel-computing/gpuarray.html | `CUDALink`: `CUDADot`, `CUDATranspose`, `CUDAFourier`, `CUDASparseMatrix`, `CUDASparseVector`; `OpenCLLink` alternative — https://reference.wolfram.com/language/CUDALink/guide/CUDALink.html | none | see section 6 (no realistic pure-Go path) |

## 2. Reference algorithms

DOIs checked via Crossref. SIAM book DOIs are the book-level DOI.

| capability | algorithm | reference | why this one |
|---|---|---|---|
| sparse direct, textbook | CSparse design: elimination tree, symbolic analysis, left-looking LU, up-looking Cholesky | Davis 2006, *Direct Methods for Sparse Linear Systems*, doi:10.1137/1.9780898718881 | 2,000-line reference; right size for a Go port and for tests |
| sparse LU, partial pivoting | left-looking, time proportional to flops | Gilbert & Peierls 1988, doi:10.1137/0909058 | basis of CSparse `cs_lu` and KLU |
| sparse LU, unsymmetric multifrontal | UMFPACK | Davis 2004, "Algorithm 832", doi:10.1145/992200.992206 | MATLAB `\` for general sparse (GPL: oracle only) |
| supernodal / BLAS-3 LU | SuperLU | Li 2005, doi:10.1145/1089014.1089017 | BSD-style; SciPy `splu` uses it |
| multifrontal theory | multifrontal method | Liu 1992, doi:10.1137/1034004 | explains the elimination tree to frontal-matrix mapping |
| circuit matrices LU | KLU: BTF + AMD/COLAMD + Gilbert-Peierls | Davis & Palamadai Natarajan 2010, "Algorithm 907", doi:10.1145/1824801.1824814 | very sparse, low fill; LGPL |
| sparse Cholesky | CHOLMOD | Chen, Davis, Hager, Rajamanickam 2008, "Algorithm 887", doi:10.1145/1391989.1391995 | supernodal; MATLAB `chol` on sparse |
| sparse QR | SuiteSparseQR (multifrontal) | Davis 2011, "Algorithm 915", doi:10.1145/2049662.2049670 | rank-revealing least squares; GPL |
| ordering, symmetric | AMD | Amestoy, Davis, Duff 1996, doi:10.1137/S0895479894278952 | de facto fill reducer; BSD-3 original |
| ordering, unsymmetric / QR | COLAMD | Davis, Gilbert, Larimore, Ng 2004, doi:10.1145/1024074.1024079 | column ordering for LU and QR; BSD-3 |
| ordering, nested dissection | multilevel graph partitioning (METIS) | Karypis & Kumar 1998, doi:10.1137/S1064827595287997 | better than AMD on 2-D/3-D meshes and for parallel factorisation |
| CG | conjugate gradients | Hestenes & Stiefel 1952, doi:10.6028/jres.049.044 ; Saad 2003 ch. 6, doi:10.1137/1.9780898718003 | SPD systems |
| GMRES | Arnoldi + least squares, restarted | Saad & Schultz 1986, doi:10.1137/0907058 | minimal residual, nonsymmetric |
| BiCGSTAB | stabilised Bi-CG | van der Vorst 1992, doi:10.1137/0913035 | short recurrence, smooth convergence |
| MINRES / SYMMLQ | Lanczos-based, indefinite symmetric | Paige & Saunders 1975, doi:10.1137/0712047 | symmetric indefinite systems |
| incomplete factorisation | IC(0)/ILU(0) | Meijerink & van der Vorst 1977, doi:10.1090/S0025-5718-1977-0438681-4 | cheapest preconditioner for M-matrices |
| ILUT | dual threshold ILU | Saad 1994, doi:10.1002/nla.1680010405 | the ILU in MATLAB/Mathematica `"ILUT"` family; SciPy `spilu` is SuperLU's ILU |
| algebraic multigrid | smoothed aggregation | Vanek, Mandel, Brezina 1996, doi:10.1007/BF02238511 | robust AMG for elliptic problems; classical Ruge-Stueben not verified here |
| Arnoldi eigensolver | implicitly restarted Arnoldi/Lanczos (ARPACK) | Sorensen 1992, doi:10.1137/0613025 ; Lehoucq & Sorensen 1996 (deflation), doi:10.1137/S0895479895281484 ; *ARPACK Users' Guide* 1998, doi:10.1137/1.9780898719628 | what MATLAB `eigs` and SciPy `eigs/eigsh` use |
| Arnoldi, modern | Krylov-Schur | Stewart 2002, doi:10.1137/S0895479800371529 | simpler restart than IRAM, same cost; basis of SLEPc default |
| symmetric block eigensolver | LOBPCG | Knyazev 2001, doi:10.1137/S1064827500366124 | needs only a preconditioner and mat-vec; SciPy `lobpcg` |
| sparse SVD | augmented implicitly restarted Lanczos bidiagonalisation (IRLBA) | Baglama & Reichel 2005, doi:10.1137/04060593X | what R `irlba` uses; PROPACK alternative (not verified here) |
| expm | scaling and squaring, Pade order chosen by backward error | Higham 2005, doi:10.1137/04061101X ; Al-Mohy & Higham 2009, doi:10.1137/09074721X | MATLAB `expm` cites exactly these two |
| expm action | truncated Taylor with norm-based parameter selection | Al-Mohy & Higham 2011, doi:10.1137/100788860 | `expm(A)b` without forming `expm(A)`; sparse-friendly; SciPy `expm_multiply` |
| survey of expm | 19 methods | Moler & Van Loan 2003, doi:10.1137/S00361445024180 | why not Taylor/eigendecomposition |
| logm | inverse scaling and squaring, Pade | Al-Mohy & Higham 2012, doi:10.1137/110852553 | improved version of Higham's `logm` |
| sqrtm | real/complex Schur method | Bjorck & Hammarling 1983, doi:10.1016/0024-3795(83)80010-X ; Higham 1987 (real roots), doi:10.1016/0024-3795(87)90118-2 | O(n^3), stable |
| fractional power | Schur-Pade | Higham & Lin 2011, doi:10.1137/10081232X | `matrix_power(A, p)` for real p |
| general f(A) | Schur-Parlett (blocked) | Davies & Higham 2003, doi:10.1137/S0895479802410815 | what `funm` does; Mathematica `"Schur"` is block Parlett too |
| f(A) theory and catalogue | all of the above | Higham 2008, *Functions of Matrices*, doi:10.1137/1.9780898717778 | the reference |
| norm estimation | block 1-norm estimator | Higham & Tisseur 2000, doi:10.1137/S0895479899356080 ; Hager 1984, doi:10.1137/0905023 ; Higham 1988 codes, doi:10.1145/50063.214386 | `rcond`, `condest`, `onenormest`; needed by expm parameter choice |
| Kronecker | `(B (x) A) vec(X) = vec(A X B^T)` identities | Van Loan 2000, doi:10.1016/S0377-0427(00)00393-9 | never form the Kronecker product for mat-vec |
| tensors | CP/Tucker decompositions | Kolda & Bader 2009, doi:10.1137/07070111X | survey; low priority |
| einsum paths | contraction-order optimisation | Smith & Gray 2018 (opt_einsum), doi:10.21105/joss.00753 | `einsum` performance |
| Toeplitz / banded | Levinson-Durbin, Gohberg-Semencul, banded LU/Cholesky, circulant embedding via FFT | no single verified source; textbook (Golub & Van Loan) (unverified here) | dedicated O(n^2)/O(n log n) paths |
| test collection paper | SuiteSparse Matrix Collection | Davis & Hu 2011, doi:10.1145/2049662.2049663 | citing the corpus |

## 3. Reference implementations

SuiteSparse licences are per package, from
https://raw.githubusercontent.com/DrTimothyAldenDavis/SuiteSparse/dev/LICENSE.txt (read 2026-10-17).
Port rule from the template: BSD/MIT/Apache = port; GPL/LGPL/proprietary = oracle only.

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| SuiteSparse AMD, COLAMD, CAMD, CCOLAMD | C | BSD-3-Clause | yes (port) | the orderings; COLAMD+AMD are the first work items |
| SuiteSparse CSparse / CXSparse | C | LGPL-2.1-or-later | oracle only (clean-room from Davis 2006 text; book algorithms are not the code) | LGPL; do not translate the source |
| SuiteSparse KLU, BTF, LDL | C | LGPL-2.1-or-later | oracle only | |
| SuiteSparse CHOLMOD | C | mixed: Cholesky, Check, Utility, Partition LGPL-2.1+; Supernodal, Modify, MatrixOps GPL-2.0+; Include Apache-2.0 | oracle only | |
| SuiteSparse UMFPACK, SPQR | C/C++ | GPL-2.0-or-later | oracle only | |
| SuiteSparse ParU | C++ | GPL-3.0-or-later | oracle only | |
| SuiteSparse Mongoose | C++ | GPL-3.0 | oracle only | |
| SuiteSparse SPEX | C | LGPL-3.0-or-later OR GPL-2.0-or-later | oracle only | exact arithmetic solves |
| SuiteSparse GraphBLAS, LAGraph | C | Apache-2.0 (with third-party parts under other licences, per README) | graph-algebra only; out of scope | https://github.com/DrTimothyAldenDavis/SuiteSparse |
| METIS 5.x | C | Apache-2.0 (LICENSE on KarypisLab/METIS main, read 2026-10-17) | yes (port) | SuiteSparse bundles a copy |
| SuperLU | C | BSD-3-Clause style (License.txt on xiaoyeli/superlu main, read 2026-10-17) | yes (port) | |
| ARPACK-ng | Fortran 77 | BSD-3-Clause (COPYING read: https://raw.githubusercontent.com/opencollab/arpack-ng/master/COPYING) | yes (port; Fortran to Go by hand, test against it) | reverse-communication interface suits a Go callback design |
| SciPy `scipy.sparse`, `scipy.sparse.linalg` | Python/C/C++ | BSD-3-Clause (SciPy 1.18.1 on PyPI https://pypi.org/pypi/scipy/json; bundled deps differ) | yes for SciPy's own code; its SuperLU/ARPACK/UMFPACK glue follows the component licences | API model to copy: `LinearOperator`, `aslinearoperator`, `eigs`, `eigsh`, `lobpcg`, `svds`, `expm_multiply`, `funm_multiply_krylov`, `onenormest`, `spilu` |
| Eigen | C++ | MPL-2.0 (COPYING.MPL2 on gitlab.com/libeigen/eigen master, read 2026-10-17) | weak copyleft: treat as oracle | https://gitlab.com/libeigen/eigen |
| PETSc | C | BSD-2-Clause (https://raw.githubusercontent.com/petsc/petsc/main/LICENSE, BSD text confirmed 2026-10-17) | yes (port of ideas/algorithms) | KSP/PC reference design; SLEPc (eigen) licence not checked |
| gonum `mat` | Go | BSD-3-Clause (https://pkg.go.dev/gonum.org/v1/gonum/mat) | yes, usable as dependency or source | dense only: Cholesky, LU, QR, LQ, SVD, GSVD, Eigen, EigenSym, BandCholesky, PivotedCholesky; `Exp`, `Pow`, `Kronecker`, `Condition`, banded and tridiagonal types. No sparse types and no `logm`/`sqrtm` found on that page. Pure-Go BLAS/LAPACK in `gonum.org/v1/gonum/blas/gonum`, `lapack/gonum` (CGO=0 capable; not re-checked here) |
| james-bowman/sparse | Go | MIT (GitHub API licence, 2026-10-17) | yes | DOK, COO, CSR, CSC on top of gonum `mat.Matrix`; third-party, not gonum core |
| Octave (`expm`, `logm`, `sqrtm`, `funm`, `eigs`, `pcg`, `gmres`...) | C++/Octave | GPL-3.0 (GitHub API licence of gnu-octave/octave, 2026-10-17) | oracle only | |
| MATLAB, Mathematica | proprietary | proprietary | oracle only (not obtainable without licence) | documentation only |

## 4. Validation

Oracles (Go harness: `exec.Command`, write a Matrix Market or JSON on stdin, read JSON on stdout):

| oracle | version seen | obtain without root | drive |
|---|---|---|---|
| SciPy `scipy.sparse.linalg`, `scipy.linalg` (`expm`, `logm`, `sqrtm`, `funm`, `splu`, `eigsh`, `svds`) | 1.18.1, Python >= 3.12 (PyPI JSON) | `uv run --with scipy python harness.py` (macOS arm64 and Linux wheels) | JSON arrays on a pipe; send CSR as `data/indices/indptr` |
| GNU Octave | not checked | release tarball / pkgx; exact availability unverified | `octave --eval` with `jsonencode` |
| SuiteSparse tools (`umfpack`, `cholmod` demo) | 7.x per README header, date unverified | build from source (no root) | files in Matrix Market |
| ARPACK-ng | master | build from source | via SciPy `eigs` instead |

Corpora:

| corpus | URL | size | licence/terms |
|---|---|---|---|
| SuiteSparse Matrix Collection | https://sparse.tamu.edu/ (paper doi:10.1145/2049662.2049663) | 2,904 matrices per site on 2026-10-17; formats MATLAB, Rutherford-Boeing, Matrix Market; site banner says collection OFFLINE due to download volume, no ETA | no terms on the page read; check the GitHub repo before redistributing |
| Matrix Market | https://math.nist.gov/MatrixMarket/formats.html | exchange format: `%%MatrixMarket matrix coordinate real|complex|integer|pattern general|symmetric|skew-symmetric|hermitian`, 1-based indices | format only; matrices in the NIST repository have their own terms (not checked) |
| Higham's `matlab-toolbox` / Matrix Computation Toolbox, MATLAB `gallery` | not fetched | small | unverified; gallery matrices (Frank, Kahan, Hilbert, etc.) are the classic ill-conditioned set for expm/logm/sqrtm |

What "correct" means:

| area | criterion | measure to report |
|---|---|---|
| sparse direct solve | normwise backward error | `||b - Ax||_inf / (||A||_inf ||x||_inf + ||b||_inf)`, expect O(u) on well-posed cases; also fill ratio nnz(L+U)/nnz(A) vs AMD/COLAMD reference |
| orderings | exact permutation match is NOT required (tie-breaking differs) | nnz(L) and flop count within a stated tolerance of SuiteSparse; permutation is valid |
| iterative solvers | final true residual, not the recurrence residual | `||b - Ax|| / ||b||` recomputed explicitly; iteration count vs SciPy (same tolerance, same preconditioner) |
| eigen (sparse) | residual per pair | `||Av - lambda v|| / (|lambda| ||v||)` for returned pairs, plus eigenvalues vs dense `eigh` on matrices small enough |
| sparse SVD | residual | `||A v - sigma u||`, singular values vs dense SVD |
| matrix functions | relative forward error in the 1-norm vs a high-precision reference (e.g. mpmath or exact for diagonalisable test cases), normalised by `cond(f,A) u` | Higham's standard: error/(cond * u) should be O(1); plot as performance profile |
| norm and condition estimates | estimate is a lower bound, ratio to true 1-norm | ratio in (0,1], typically > 0.5 |
| structured (Toeplitz, banded, Kronecker) | agree with the dense path | max relative difference vs dense solve; operation count scaling checked by benchmark |
| exact integer results | equality | exact |

## 5. Proposed home and plan

Proposed NEW orgs (checked 2026-10-17 with `gh api users/<name>`): `go-sparse` -> 404 (free);
`go-linalg` -> 404 (free). `go-matrix` EXISTS (taken; its single repo `matrix` was last pushed 2014-08, no description, no usable tree). Dense work stays in
`go-ndarray/ndarray/linalg`, where the factorisations already are; `go-linalg` is not needed.

| work item | org/repo | depends on | size S/M/L/XL | priority 1–3 |
|---|---|---|---|---|
| Dense additions: real Schur with reordering, QZ, kron, 1-norm condition estimation (LU/QR/Cholesky/SVD/eig/eigh/lstsq/solve/inv/det/cond exist) | `go-ndarray/ndarray/linalg` | ndarray GEMM; gonum `mat`/`lapack` as BSD source or oracle | L | 1 |
| Sparse containers CSR/CSC/COO/DOK, conversions, SpMV, SpMM, transpose, Matrix Market + Rutherford-Boeing I/O, `LinearOperator` | NEW org `go-sparse`, repo `sparse` | ndarray (optional) | M | 1 |
| Orderings: AMD, COLAMD, RCM (BSD sources); nested dissection as separate item | `go-sparse/order` | sparse | M (AMD+COLAMD+RCM); L (ND/METIS-class) | 1 / 3 |
| Sparse direct: etree + symbolic + LU (Gilbert-Peierls, CSparse design) and Cholesky/LDL' (up-looking), then supernodal; QR later | `go-sparse/direct` | sparse, order, dense kernels | L (basic); XL (supernodal/multifrontal, QR) | 2 |
| Iterative solvers CG, GMRES, BiCGSTAB, MINRES, LSQR/LSMR + preconditioners Jacobi, SSOR, IC0, ILU0, ILUT | `go-sparse/iter` | sparse, `LinearOperator` | M | 1 |
| AMG (smoothed aggregation) | `go-sparse/amg` | iter | L | 3 |
| Sparse eigen: Lanczos/IRLM, Krylov-Schur, shift-invert; LOBPCG; sparse SVD (Lanczos bidiagonalisation) | `go-sparse/eigs` | iter, direct (shift-invert), dense eigh | L | 2 |
| Matrix functions: expm (Al-Mohy-Higham 2009), expm_multiply, logm, sqrtm, fractional power, Schur-Parlett funm, Frechet derivatives | `go-linalg/matfun` | dense Schur/eig (item 1), norm estimator | L | 1 |
| Norm/condition estimators (Higham-Tisseur block 1-norm, `rcond`, `onenormest`) | `go-linalg/condest` (needed by matfun) | dense LU | S | 1 |
| Structured: Kronecker (matrix-free), Toeplitz/circulant (Levinson, FFT), banded LU/Cholesky | `go-linalg/structured` | dense, `go-ndarray` fft branch | M | 2 |
| Tensor: einsum + tensordot + contraction path; CP/Tucker | `go-ndarray/ndarray` (einsum), decompositions later | ndarray | M (einsum); L (decompositions) | 2 / 3 |
| Conformance harness vs SciPy/Octave + SuiteSparse corpus runner | `go-scikit/roadmap` or `go-sparse/conformance` | all | M | 1 (build first) |
| GPU | none | none | n/a | not planned |

## 6. Risks and open questions

- Every sparse and matrix-function item depends on dense Schur/eig/LU/QR/SVD; the factorisations
  exist in go-ndarray/linalg, Schur reordering and QZ do not, so those are the critical path.
  go-ndarray is owned by another session: coordinate before changing it.
- Licence trap: the high-quality sparse direct solvers (CHOLMOD supernodal, UMFPACK, SPQR, ParU) are
  GPL/LGPL. Only AMD/COLAMD are BSD-3. LU/Cholesky must be written from the Davis 2006 text and the
  papers, not translated from CSparse (LGPL). Record clean-room provenance in each package.
- SuiteSparse Matrix Collection is reported OFFLINE (site banner, 2026-10-17); validation tests must
  not download at test time. Mirror a small named subset outside every repository (per the
  global rule on large artefacts) and pin checksums; check the redistribution terms first.
- Sparse SVD (IRLBA) and ARPACK shift-invert are numerically delicate (restart, deflation,
  loss of orthogonality); matching SciPy iteration counts is not a realistic criterion, matching
  residuals is.
- Matrix-function accuracy claims need a high-precision reference; Go has no mpmath. Use
  `math/big` based Taylor for small test cases, or SciPy-generated reference values with
  mpmath, stored as fixtures.
- Fill-reducing ordering results differ across implementations (tie-breaking); do not test for
  permutation equality.
- GPU: MATLAB `gpuArray` is NVIDIA-CUDA-only and needs Parallel Computing Toolbox (sparse supported
  with the limits in section 1); Mathematica has `CUDALink` with `CUDASparseMatrix`. A pure-Go
  (CGO=0) path is not realistic: every vendor GPU stack (CUDA, ROCm, Metal, Vulkan compute) is
  reached through C or Objective-C ABIs. Only a purego/dlopen-style binding could avoid cgo, and
  that is an engineering judgement of this document, not a sourced fact. Not worth planning
  before the CPU stack exists. SIMD via `go-asmgen` is the realistic accelerator.
- Still unverified: SLEPc licence, Ruge-Stueben AMG, PROPACK, Golub & Van Loan Toeplitz source,
  Octave/SuiteSparse-tools availability, SuiteSparse collection terms.

Verified on 2026-10-17: all 41 DOIs resolved on Crossref (title, year, first author match; Hestenes-Stiefel
DOI added); MATLAB sparse-index names from the MathWorks page, `funm`, `logm`, `expm`, `tensorprod`,
`pagemtimes` pages (other core names not re-fetched: MathWorks rate-limits curl); Mathematica names via
reference.wolfram.com (`MatrixConditionNumber` does not exist, removed; CUDALink names exist under
`/language/CUDALink/ref/`); LinearSolve/Eigenvalues/MatrixFunction method strings; licences from
LICENSE files (SuiteSparse per-package table matches, Mongoose GPL-3.0 confirmed; METIS, SuperLU, PETSc,
ARPACK-ng, Eigen, SciPy, gonum, james-bowman/sparse, Octave); go-ndarray and sparse claims against the
GitHub API tree and `linalg/*.go` symbols.
