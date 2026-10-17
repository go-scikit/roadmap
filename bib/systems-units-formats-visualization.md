# Control and dynamical systems, units, scientific file formats, 3D visualisation — bibliography

One file per domain. Facts only, each with a source (URL or DOI). Terse: tables, not prose.

Inventory date 2026-10-17. Local state was read from the clones under `~/Documents/VCS/GIT/github.com/` (`GOWORK=off go doc`, `ls`, `grep`); remote facts from `gh api repos/<o>/<r>` (licence, last push). Items still marked (unverified) were not re-fetched.

## 1. Feature inventory

### 1a. Control and dynamical systems

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| model objects (TF, ZPK, SS, FRD, descriptor, delay) | `tf zpk ss frd dss`, Control System Toolbox https://www.mathworks.com/help/control/ | `TransferFunctionModel`, `StateSpaceModel`, `AffineStateSpaceModel` https://reference.wolfram.com/language/guide/ControlSystems.html | none | whole package |
| model algebra (series, parallel, feedback, connect) | `series parallel feedback connect lft` | `SystemsModelSeriesConnect`, `SystemsModelParallelConnect`, `SystemsModelFeedbackConnect` | none | whole |
| zeros, poles, DC gain, minimal realisation | `zero pole dcgain minreal ss2tf tf2ss balreal` | `TransferFunctionPoles`, `TransferFunctionZeros`, `MinimalStateSpaceModel` | none; needs non-symmetric eig + QZ | eig/QZ missing in go-ndarray |
| time response | `step impulse lsim initial stepinfo` | `OutputResponse`, `StateResponse` | integrate: RK45, DOP853, Radau, BDF, LSODA, `Event` (go-numerics/integrate `ivp.go`) | LTI-exact response via matrix exponential (`expm`) absent |
| frequency response | `bode nyquist nichols freqresp sigma margin bandwidth` | `BodePlot`, `NyquistPlot`, `NicholsPlot`, `GainMargins`, `PhaseMargins`, `GainPhaseMargins` | go-plot/plot 2D lines only | response evaluation, Nyquist encirclement, margins |
| root locus | `rlocus` | `RootLocusPlot` | none | whole |
| state feedback, pole placement | `place acker` | `StateFeedbackGains` | none | whole |
| LQR / LQG / Kalman | `lqr dlqr lqe kalman lqg lqi` | `LQRegulatorGains`, `LQEstimatorGains`, `LQOutputRegulatorGains` | none | CARE/DARE solvers |
| Riccati, Lyapunov, Sylvester | `icare idare lyap dlyap sylvester` | `LyapunovSolve`; Riccati via `LQRegulatorGains` | none | Schur ordering needed |
| discretisation | `c2d d2c d2d` (zoh, foh, tustin, matched) | `ToDiscreteTimeModel`, `ToContinuousTimeModel` | none | `expm` needed |
| system identification | System Identification Toolbox `ssest n4sid arx oe armax tfest` https://www.mathworks.com/help/ident/ | no identification function (no `EstimatedModel`); `FindFit`/`LinearModelFit` only | go-numerics/optimize (least squares) | subspace methods need SVD/QR |
| model reduction | `balred balreal modred` | `InternallyBalancedDecomposition` | none | whole |
| robust control | Robust Control Toolbox `hinfsyn musyn` | none | none | out of scope for first wave |
| block-diagram simulation (Simulink) | Simulink, `sim`, solvers ode45, ode23t, ode15s; algebraic-loop handling; zero-crossing events | System Modeler (Modelica), `SystemModel`, `SystemModelSimulate` https://reference.wolfram.com/language/ref/SystemModelSimulate.html | none | graph compiler + DAE solver |
| DAE solving | `ode15i`, `decic`, `ode15s` with mass matrix | `NDSolve` with `"IDA"` method, `DAE` index handling | integrate has stiff ODE (Radau, BDF) and events; no DAE, no mass matrix (grep `ivp.go`, `README.md`) | IDA/DASSL-type BDF, consistent initialisation, index reduction |
| symbolic/Modelica | Simscape language | Modelica import/export | none | out of scope |

### 1b. Units and physical constants

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| quantities with units | `symunit` (Symbolic Math Toolbox), `u = symunit`, `unitConvert`, `checkUnits`, `separateUnits`, `baseUnits` https://www.mathworks.com/help/symbolic/units-of-measurement.html | `Quantity`, `UnitConvert`, `QuantityMagnitude`, `CompatibleUnitQ`, `UnitDimensions`, `UnitSimplify` https://reference.wolfram.com/language/guide/UnitsAndQuantities.html | none (inventory of fleet: no unit library; go-datetime covers dates only) | whole |
| dimensional analysis | `checkUnits`; Buckingham pi by hand | `DimensionalCombinations` | none | whole |
| physical constants | `physconst` (Phased Array System Toolbox) | `Quantity["SpeedOfLight"]`, `Entity["PhysicalConstant"]` | none | CODATA table |
| temperature offsets (degC, degF) | `unitConvert` | `Quantity[..., "DegreesCelsius"]` | none | affine units |

### 1c. Scientific file formats

| format | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| CSV | `readtable readmatrix csvread writetable` | `Import[*, "CSV"]` | go-dataframe: `ReadCSV`, `ReadCSVFile`, `csv_write.go` | none |
| Parquet / Arrow IPC / Feather | `parquetread parquetwrite` | `Import[*, "Parquet"]`, `"ArrowIPC"` | go-dataframe sub-packages `parquetio`, `arrowio` | none |
| JSON | `jsondecode jsonencode` | `ImportString[*, "JSON"]` | no JSON reader in go-dataframe (grep finds only test fixtures); stdlib `encoding/json` | thin adapter |
| Excel xlsx | `readtable`, `xlsread` | `Import[*, "XLSX"]` | go-ooxml/ooxml `workbook.go`, `writeworkbook.go` over go-ooxml/opc | bridge to dataframe |
| MAT-file v4 / v5 / v7 / v7.3 | `load save matfile` https://www.mathworks.com/help/matlab/import_export/mat-file-versions.html | `Import[*, "MAT"]` (v5) | none | whole |
| HDF5 | `h5read h5write h5info h5create` | `Import[*, "HDF5"]` | none | whole |
| NetCDF classic, 64-bit offset, NetCDF-4 | `ncread ncwrite ncinfo` | `Import[*, "NetCDF"]` | none | whole |
| FITS | `fitsread fitswrite fitsinfo` | `Import[*, "FITS"]` | none | whole |
| WAV / FLAC / MP3 | `audioread audiowrite` | `Import[*, "WAV"]`, `"FLAC"` | none | WAV + FLAC |
| Matrix Market | none native (File Exchange `mmread`) | `Import[*, "MTX"]` | none | whole (small) |
| NumPy `.npy/.npz` | none native (File Exchange `npy-matlab`) | no `Import` format "NPY" documented | none | small |

### 1d. Visualisation beyond 2D

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| 3D lines, scatter | `plot3 scatter3` | `ListPointPlot3D ListLinePlot3D ParametricPlot3D` | go-plot/plot is 2D (axes, lines, scatter, contour, colorbar, patches, SVG/raster renderers; grep for surface/zbuf finds none) | whole |
| surfaces, meshes | `surf mesh surfc waterfall ribbon trisurf` | `Plot3D ListPlot3D ParametricPlot3D RegionPlot3D` | none | whole |
| contours in 3D, quiver | `contour3 quiver3 streamline` | `ContourPlot3D VectorPlot3D StreamPlot3D` | 2D contour exists (`go-plot/plot/contour`) | extend |
| isosurface, slices, volume | `isosurface isonormals isocaps slice volshow` | `ContourPlot3D`, `Image3D`, `ListSliceContourPlot3D` | none | marching cubes + volume |
| colormaps | `parula turbo viridis colormap` | `ColorData` gradients | go-plot/plot `colors/` (viridis, plasma, inferno, magma, cividis; no turbo, no Crameri) | add turbo, Crameri |
| camera, lighting | `view camlight lighting material camproj` | `ViewPoint`, `Lighting`, `ViewVertical` | go-gfx/gfx `iso` (isometric projection, depth sort, scene composition) | perspective + general camera |
| interactive rotate | `rotate3d` | direct drag in front end | none | browser side |
| animation | `getframe VideoWriter`, `animatedline` | `Animate`, `ListAnimate` | none | frame loop + encoder |
| interactive widgets | App Designer, `uicontrol`, Live Editor controls https://www.mathworks.com/help/matlab/live-scripts-and-functions.html | `Manipulate`, `Dynamic`, `Slider`, `ControlActive` https://reference.wolfram.com/language/ref/Manipulate.html | go-notebook: `notebook serve` re-runs only changed cells; cell API `Cell.Show`, `Cell.Result` (`runtime/runtime.go`); dependency analysis in `analysis/graph.go` | widget cell + reactive re-run |

## 2. Reference algorithms

### 2a. Control

| capability | algorithm | reference | why |
|---|---|---|---|
| CARE / DARE | ordered real Schur form of the Hamiltonian / symplectic pencil | Laub, "A Schur method for solving algebraic Riccati equations", IEEE TAC 24(6) 1979, doi:10.1109/TAC.1979.1102178 | the reference method; backward-stable up to conditioning; QZ variant for singular pencils |
| CARE refinement | Newton / Kleinman iteration, with Schur result as start | Kleinman, IEEE TAC 13(1), 1968, doi:10.1109/TAC.1968.1098829 | quadratic convergence polishes the residual |
| Sylvester / Lyapunov | Bartels-Stewart (real Schur reduction, back-substitution) | Bartels, Stewart, "Solution of AX+XB=C", CACM 15(9) 1972, doi:10.1145/361573.361582 | O(n^3), standard in `lyap`, `sylvester`, SLICOT SB04MD |
| Lyapunov, Cholesky factor directly | Hammarling | Hammarling, IMA J Numer Anal 2(3) 1982, doi:10.1093/imanum/2.3.303 | positive-definite factor without forming the product |
| matrix exponential (c2d zoh, lsim) | scaling and squaring with Pade, optimal degree | Al-Mohy, Higham, SIAM J Matrix Anal Appl 31(3) 2009, doi:10.1137/09074721X; Higham, "The scaling and squaring method for the matrix exponential revisited", SIAM Review 51(4) 2009, doi:10.1137/090768539 | backward error bound, used by SciPy `expm` |
| exact ZOH with input matrix | augmented-matrix exponential | Van Loan, IEEE TAC 23(3) 1978, doi:10.1109/TAC.1978.1101743 | one `expm` call gives Ad and Bd |
| pole placement | robust method (KNV) | Kautsky, Nichols, Van Dooren, Int J Control 41(5) 1985, doi:10.1080/0020718508961188 | numerically stable for multi-input; `acker` is unstable above order ~10 |
| Kalman filter (discrete, square-root) | Joseph form / UD factorisation | Kalman, J Basic Eng 82(1) 1960, doi:10.1115/1.3662552; Bierman, Factorization Methods for Discrete Sequential Estimation, 1977 (ISBN 978-0-486-44981-4) | covariance stays positive-definite |
| frequency response | Hessenberg reduction, then solve (jwI-A) x = B per frequency | Laub, "Efficient multivariable frequency response computations", IEEE TAC 26(2) 1981, doi:10.1109/TAC.1981.1102602 | O(n^2) per point instead of O(n^3) |
| stability margins | root-find on |L(jw)|=1 and arg L = -pi, from the dense grid | Skogestad, Postlethwaite, Multivariable Feedback Control, 2nd ed., Wiley 2005, ISBN 978-0-470-01168-3 | textbook definition; edge cases for delays |
| balanced truncation | square-root method | Laub, Heath, Paige, Ward, IEEE TAC 32(2) 1987, doi:10.1109/TAC.1987.1104549 | numerically safe |
| subspace identification | N4SID, MOESP | Van Overschee, De Moor, Subspace Identification for Linear Systems, Kluwer 1996, doi:10.1007/978-1-4613-0465-4 | non-iterative, needs SVD |
| prediction-error identification | Gauss-Newton / Levenberg-Marquardt on innovations | Ljung, System Identification: Theory for the User, 2nd ed., Prentice Hall 1999, ISBN 978-0-13-656695-3 | the MATLAB toolbox's reference |
| textbook basis | | Franklin, Powell, Emami-Naeini, Feedback Control of Dynamic Systems, 8th ed., Pearson 2019, ISBN 978-0-13-468571-7 | worked examples usable as tests |

### 2b. DAE and simulation

| capability | algorithm | reference | why |
|---|---|---|---|
| DAE of index 1, fully implicit F(t,y,y')=0 | variable-order, variable-step BDF (DASSL / IDA) | Brenan, Campbell, Petzold, Numerical Solution of Initial-Value Problems in DAEs, SIAM 1996, doi:10.1137/1.9781611971224; Petzold, SIAM J Sci Stat Comput 4(1) 1983, doi:10.1137/0904010 | the algorithm behind MATLAB `ode15i`-class tools and SUNDIALS IDA |
| semi-explicit / mass matrix M y' = f | Radau IIA, RODAS; existing Radau extends | Hairer, Wanner, Solving ODEs II, 2nd ed., Springer 1996, doi:10.1007/978-3-642-05221-7 (chapters VI, VII) | integrate already implements Radau IIA order 5 (`radau.go` header cites IV.8) |
| consistent initial conditions | algebraic Newton on the constraint block (IDACalcIC, `decic`) | Brown, Hindmarsh, Petzold, "Consistent initial condition calculation for differential-algebraic systems", SIAM J Sci Comput 19(5) 1998, doi:10.1137/S1064827595289996 | solver fails otherwise |
| index reduction | Pantelides algorithm; dummy-derivative method | Pantelides, SIAM J Sci Stat Comput 9(2) 1988, doi:10.1137/0909014; Mattsson, Soderlind, SIAM J Sci Comput 14(3) 1993, doi:10.1137/0914043 | what Modelica tools apply before the solver |
| structural analysis, block ordering | Tarjan strongly-connected components, BLT | Tarjan, SIAM J Comput 1(2) 1972, doi:10.1137/0201010 | sorted evaluation order, algebraic loops |
| events (zero crossing, discontinuities) | root-find on dense output, restart | already in integrate (`Event`, `TEvents`) | state-event handling for hybrid models |
| language | Modelica Language Specification 3.6 https://specification.modelica.org/ | defines equations, connectors, `when`, synchronous semantics | only as target vocabulary; a first wave can be a Go block-diagram DSL |

### 2c. Units

| capability | source | note |
|---|---|---|
| SI definition (2019 redefinition, 7 base units, prefixes) | BIPM, The International System of Units (SI Brochure), 9th ed., 2019 (v3.01 2025) https://www.bipm.org/en/publications/si-brochure | normative |
| constants values | CODATA 2022 recommended values, Mohr, Newell, Taylor, Tiesinga, Rev. Mod. Phys. 97 (2025) https://physics.nist.gov/cuu/Constants/ | exact for h, e, k, N_A; others with uncertainty and correlations |
| unit grammar and database | UDUNITS-2 https://www.unidata.ucar.edu/software/udunits/ (grammar, XML database of ~600 names) | NetCDF CF conventions use it; licence: BSD-3-Clause text plus a patent-termination clause (COPYRIGHT file; `gh api` reports NOASSERTION): permissive, may port |
| dimensional analysis | Buckingham, Phys. Rev. 4 (1914) 345, doi:10.1103/PhysRev.4.345 | pi theorem via null space of exponent matrix |
| representation | vector of 7 exponents (rational) plus scale factor and offset; offset only for degC, degF | Hart, "Multidimensional Analysis", Springer 1995, doi:10.1007/978-1-4612-4208-6 |
| unified codes | UCUM https://ucum.org/ucum | case-sensitive, used in health data |
| dates and time scales | out of scope (go-datetime) | |

### 2d. File formats

| format | specification | key points |
|---|---|---|
| HDF5 | HDF Group, HDF5 File Format Specification v3.0 https://support.hdfgroup.org/documentation/hdf5/latest/_f_m_t3.html | superblock v0-v3, object headers v1/v2, B-trees v1/v2, fractal heap, local heap, chunked layout, filters (deflate, shuffle, fletcher32), global heap for vlen |
| MAT v4 | MathWorks, MAT-File Format (PDF) https://www.mathworks.com/help/pdf_doc/matlab/matfile_format.pdf (chapter v4 legacy) | 20-byte header per matrix |
| MAT v5 / v7 | same PDF | 128-byte header, tagged data elements, miCOMPRESSED (zlib) in v7 |
| MAT v7.3 | same PDF | HDF5 file with 512-byte user block; `MATLAB_class` attribute; complex as compound; sparse as group |
| NetCDF classic (CDF-1,2,5) | Unidata, NetCDF File Format Specification https://docs.unidata.ucar.edu/netcdf-c/current/file_format_specifications.html | big-endian, record dimension |
| NetCDF-4 | same page | HDF5 with dimension scales; "classic model" subset of HDF5 |
| CF conventions | https://cfconventions.org/ | units attribute is UDUNITS syntax |
| FITS | FITS Standard 4.0 (IAU FITS Working Group, 2018) https://fits.gsfc.nasa.gov/fits_standard.html | 2880-byte blocks, 80-char cards, big-endian, BINTABLE, compressed images (Rice, GZIP, HCOMPRESS) |
| WAV | Microsoft RIFF WAVE; EBML not needed | PCM 8/16/24/32, float, extensible header |
| FLAC | RFC 9639 https://www.rfc-editor.org/rfc/rfc9639 | frames, Rice coding, LPC, MD5 of samples |
| Matrix Market | NIST https://math.nist.gov/MatrixMarket/formats.html | coordinate / array, real/complex/integer/pattern, general/symmetric/skew/Hermitian |
| NumPy NPY | NEP 1 https://numpy.org/neps/nep-0001-npy-format.html | header dict, little/big endian, Fortran order |

### 2e. 3D visualisation

| capability | algorithm | reference | why |
|---|---|---|---|
| isosurface | marching cubes | Lorensen, Cline, "Marching cubes: a high resolution 3D surface construction algorithm", SIGGRAPH 1987, doi:10.1145/37402.37422 | baseline; ambiguous cases in 15-case table |
| isosurface, topologically correct | asymptotic decider; marching cubes 33 | Nielson, Hamann, IEEE Visualization 1991, doi:10.1109/VISUAL.1991.175782; Chernyaev 1995 (unverified); Lewiner et al., J Graphics Tools 8(2) 2003, doi:10.1080/10867651.2003.10487582 | resolves ambiguities |
| survey | Newman, Yi, "A survey of the marching cubes algorithm", Computers and Graphics 30(5) 2006, doi:10.1016/j.cag.2006.07.021 | choose between MC, marching tetrahedra, dual methods |
| volume rendering | ray casting with transfer function, front-to-back compositing | Levoy, "Display of surfaces from volume data", IEEE CG&A 8(3) 1988, doi:10.1109/38.511 | basis of `volshow` |
| hidden surface, SVG vector output | painter's algorithm, back-to-front depth sort by polygon centroid (cannot resolve cycles; split polygons or BSP) | Newell, Newell, Sancha, ACM Annual Conf. 1972, doi:10.1145/800193.569954; Fuchs, Kedem, Naylor, SIGGRAPH 1980, doi:10.1145/965105.807481 | vector output stays small; go-gfx `iso` already depth-sorts |
| hidden surface, raster | z-buffer | Catmull, PhD thesis, Univ. Utah 1974 | exact per pixel, needed for intersecting surfaces |
| perspective transform, camera | view and projection matrices, look-at | Hughes et al., Computer Graphics: Principles and Practice, 3rd ed., 2013 (ISBN 978-0-321-39952-6) | |
| lighting | Lambert diffuse plus Blinn-Phong specular, vertex normals by area-weighted average | Blinn, SIGGRAPH 1977, doi:10.1145/563858.563893 | `lighting gouraud` equivalent |
| colour maps | perceptually uniform: viridis (CC0), CIE Lab-based | van der Walt, Smith (viridis), https://bids.github.io/colormap/ ; Crameri et al., Nat. Commun. 11 (2020) 5444, doi:10.1038/s41467-020-19160-7 | MATLAB `parula` is proprietary: do not copy; use `turbo` (Apache 2.0, Google) https://research.google/blog/turbo-an-improved-rainbow-colormap-for-visualization/ |
| browser rendering | WebGL 2 or WebGPU from Go compiled to WASM (`GOOS=js GOARCH=wasm`) calling `syscall/js` | WebGPU spec https://www.w3.org/TR/webgpu/ ; WebGL 2.0 https://registry.khronos.org/webgl/specs/latest/2.0/ | interactive rotation at 60 fps; fall back to SVG/raster server-side |

### 2f. Manipulate-style interactivity in go-notebook

Facts about the current implementation (read in `go-notebook/notebook/README.md`, lines 66-70, and `analysis/graph.go`): `notebook serve analysis.md` re-runs only changed cells; the `.md` file is the single source of truth; cells receive a `*runtime.Cell` with `Show(vs ...any)` and `Result(v any, err error)`.

Design reference (not yet implemented): reactive dependency re-execution as in Observable (https://observablehq.com/@observablehq/how-observable-runs), Marimo (https://docs.marimo.io/guides/reactivity/) and Jupyter `ipywidgets.interact` (https://ipywidgets.readthedocs.io/). Mathematica `Manipulate` evaluates its body on every control change, with `ControlActive` selecting a cheap rendering while the control is dragged (https://reference.wolfram.com/language/ref/ControlActive.html). The needed pieces:

| piece | description |
|---|---|
| widget value type | `Slider`, `Checkbox`, `Dropdown`, `Text` returned by a cell; value stored per notebook session, not in the `.md` (or in a `param` fence if persistence is wanted) |
| dependency graph | `analysis/graph.go` already derives cell dependencies; a widget change marks its dependents dirty and re-runs them only |
| transport | WebSocket or SSE from `serve` to the page: widget value in, replacement outputs (SVG, PNG, WebGL scene JSON) out |
| cheap mode | during drag, re-run with a `Dragging()` flag so cells can lower resolution (the `ControlActive` idea) |
| debounce and cancel | cancel the running cell when a newer value arrives; kernel cache (`kernel/cache.go`) keys on inputs |

## 3. Reference implementations

### 3a. Control and DAE

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| python-control https://github.com/python-control/python-control | Python | BSD-3-Clause (`gh api`) | port | API template: `TransferFunction`, `StateSpace`, `bode_plot`, `lqr`, `place`; uses SLICOT for Riccati (via `slycot`) |
| SLICOT-Reference https://github.com/SLICOT/SLICOT-Reference | Fortran 77 | BSD-3-Clause (`gh api`, pushed 2026-10-12) | port (routines: SB02OD CARE, SB02MD, SB03MD Lyapunov, SB04MD Sylvester, TB05AD frequency response, AB01ND controllability, TB01*, SB10* H-infinity) | relicensed from GPL to BSD-3 by the SLICOT team (the repository carries the BSD licence; check each routine's header before copying) |
| slycot https://github.com/python-control/Slycot | Python/Fortran wrapper | GPL-2.0 (LICENSE.txt: GPL v2; `gh api` reports NOASSERTION) | oracle only | |
| Octave `control` package https://gnu-octave.github.io/packages/control/ | Octave/C++/Fortran | GPL-3.0+ | oracle only | MATLAB-compatible function names |
| SUNDIALS (IDA, CVODE, ARKODE) https://github.com/LLNL/sundials | C | BSD-3-Clause (`gh api`) | port IDA (DAE BDF, variable order 1-5, Newton, `IDACalcIC`, root finding) | |
| DASSL / DASPK (Petzold) https://www.netlib.org/ode/ddassl.f | Fortran | SLATEC library member (netlib header `C***LIBRARY SLATEC`; SLATEC is public domain) | port possible | ancestor of IDA |
| Hairer codes RADAU5, RODAS https://www.unige.ch/~hairer/software.html | Fortran | BSD-2-style (https://www.unige.ch/~hairer/prog/licence.txt; RADAR5 and SUMEXP codes have separate licences) | port Radau5 mass-matrix path | Radau already ported in integrate; mass matrix not |
| SciPy `scipy.signal` (`lti`, `place_poles`, `solve_continuous_are`) https://github.com/scipy/scipy | Python/C/Fortran | BSD-3-Clause | port / oracle | oracle for CARE and `expm` |
| OpenModelica https://github.com/OpenModelica/OpenModelica | C/C++ | OSMC-PL / GPL-3 (OSMC-License.txt; `gh api` NOASSERTION) | oracle only | Modelica frontend, index reduction reference |
| JModelica / Modelica Standard Library | Modelica | BSD-3 (MSL) | library of models usable as tests | |
| Julia ControlSystems.jl https://github.com/JuliaControl/ControlSystems.jl | Julia | MIT (LICENSE.md; `gh api` NOASSERTION) | port | clean API, tests against MATLAB values |
| MATLAB Control System Toolbox, Simulink, System Identification Toolbox | MATLAB | proprietary | docs only | |
| Mathematica control systems, System Modeler | Wolfram | proprietary | docs only | |
| Existing local code | Go | go-numerics/integrate (licence file present, check header), go-ndarray (`linalg`: LU, QR, Cholesky, SVD, eig, eigh, lstsq; no Schur reordering, QZ or `expm`) | n/a | linear algebra is the blocking dependency |
| gonum/mat https://github.com/gonum/gonum | Go | BSD-3-Clause (`gh api`, pushed 2026-10-08) | port/import; has `Eigen`, `SVD`, `QR`, `LU`, `Cholesky`; no real Schur form or QZ (check `go doc`) | candidate dependency or source of ports |

### 3b. Units

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| UDUNITS-2 https://github.com/Unidata/UDUNITS-2 | C | BSD-3-Clause text plus patent-termination clause (COPYRIGHT file, 2025; `gh api`: NOASSERTION) | may port the XML database and grammar (keep notice) | CF reference |
| pint https://github.com/hgrecco/pint | Python | BSD-3-Clause (LICENSE; `gh api` NOASSERTION) | port the definitions file `default_en.txt` (data) | oracle |
| astropy.units https://github.com/astropy/astropy | Python | BSD-3-Clause | oracle | includes CODATA 2018/2022 constants (`astropy.constants`) |
| Boost.Units, mp-units https://github.com/mpusz/mp-units | C++ | Boost / MIT | design reference | compile-time dimension checking |
| Go: google/go-units, docker/go-units, alecthomas/units | Go | Apache-2.0 (google), Apache-2.0 (docker), MIT (alecthomas) (`gh api`); `dbrun3/go-units` does not exist | none is a dimensional-analysis engine (typed floats or byte-size parsers) | do not depend |
| Go: bcicen/go-units | Go | MIT (`gh api`) | | conversion tables only |
| NIST CODATA ASCII table https://physics.nist.gov/cuu/Constants/Table/allascii.txt | text | public domain (US government) | embed as generated Go table | the data source |
| SciPy `scipy.constants` https://github.com/scipy/scipy | Python | BSD-3-Clause | oracle | CODATA 2022 in recent releases |

### 3c. File formats

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| scigolib/hdf5 https://github.com/scigolib/hdf5 | Go, pure | MIT (`gh api`, pushed 2026-09-02, 32 stars) | import candidate; audit | read complete, writes partial per its README |
| scigolib/matlab https://github.com/scigolib/matlab | Go, pure | MIT (pushed 2026-09-02, 12 stars) | import candidate; audit | MAT-file reader/writer, v7.3 via scigolib/hdf5 |
| flywave/hdf5 https://github.com/flywave/hdf5 | Go, pure | MIT (pushed 2026-07-10) | audit | young, 0 stars |
| go-hdf5/hdf5 https://github.com/go-hdf5/hdf5 | Go, pure | BSD-3-Clause (pushed 2018-10-19) | stale, org owned by a third party | the name `go-hdf5` is therefore NOT available to us |
| gonum/hdf5 | Go + cgo | BSD-3-Clause | no (cgo) | out |
| batchatco/go-native-netcdf https://github.com/batchatco/go-native-netcdf | Go, pure | MIT (pushed 2026-03-14, 62 stars) | import candidate; reads CDF and NetCDF-4 via own HDF5 subset | writes CDF only |
| fhs/go-netcdf | Go + cgo | MIT | no (cgo) | |
| astrogo/fitsio https://github.com/astrogo/fitsio (moved to codeberg.org/astrogo/fitsio) | Go, pure | BSD-3-Clause (archived on GitHub, last push 2025-03-04) | port / fork | image + table; read/write |
| dmarkham/fits https://github.com/dmarkham/fits | Go, pure | MIT (pushed 2026-04-14) | audit | stdlib only |
| mewkiz/flac https://github.com/mewkiz/flac | Go, pure | Unlicense (pushed 2026-08-19) | import | decoder and encoder |
| go-audio/wav https://github.com/go-audio/wav | Go | Apache-2.0 (archived 2024-10-28) | port or write fresh (WAV is simple) | the `go-audio` GitHub name is taken by a third party |
| h5py / h5dump / HDF5 C library https://github.com/HDFGroup/hdf5 | C | BSD-3-Clause-like | oracle | reference reader |
| scipy.io (`loadmat`, `savemat`, `mmread`, `netcdf_file`, `wavfile`) | Python | BSD-3-Clause | oracle for v4/v5/v7, Matrix Market, classic NetCDF, WAV | v7.3 not supported by `loadmat`; use h5py |
| astropy.io.fits https://github.com/astropy/astropy | Python | BSD-3-Clause | oracle | |
| Unidata netcdf-c https://github.com/Unidata/netcdf-c | C | BSD-3-Clause-like | oracle (`ncdump`) | |
| Local go-dataframe | Go | see repo LICENSE | n/a | I/O inventory: CSV (`csv.go`, `csv_parse.go`, `csv_write.go`), Parquet (`parquetio`), Arrow IPC and pandas interop (`arrowio`). No JSON, no xlsx bridge, no HDF5/MAT/NetCDF/FITS |
| Local go-ooxml/ooxml | Go | see repo LICENSE | n/a | `workbook.go`, `writeworkbook.go` (xlsx); depends on go-ooxml/opc |

### 3d. 3D visualisation

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| Local go-gfx/gfx (`svg`, `raster`, `vector`, `geometry`, `iso`, `color`) | Go | see repo LICENSE | extend | `iso` already has projection, depth sort and scene compositing (README line 22) |
| Local go-plot/plot | Go | see repo LICENSE | extend with `plot3d` axes | 2D only today |
| scikit-image `measure.marching_cubes` https://github.com/scikit-image/scikit-image | Python/Cython | BSD-3-Clause | oracle | Lewiner variant; compare mesh area and Euler characteristic |
| VTK (marching cubes, volume rendering) https://github.com/Kitware/VTK | C++ | BSD-3-Clause | oracle / reference | |
| three.js https://github.com/mrdoob/three.js | JavaScript | MIT | JS side of the browser viewer | could be embedded as the front end of notebook output |
| matplotlib `mplot3d` https://github.com/matplotlib/matplotlib | Python | PSF-based (BSD-compatible) | design reference | uses painter's algorithm per polygon (its documented limitation) |
| Plotly.js 3D https://github.com/plotly/plotly.js | JavaScript | MIT | candidate front end | WebGL scatter3d/surface/isosurface/volume |
| Observable Plot / d3 | JavaScript | ISC | 2D only | |
| Go: fogleman/fauxgl https://github.com/fogleman/fauxgl | Go | MIT (`gh api`) | design reference for software z-buffer rasteriser | |
| Go: g3n/engine, ebitengine | Go | BSD-2 / Apache-2.0, cgo/GL | no (cgo) | |
| Mathematica Manipulate, MATLAB Live Editor | proprietary | | docs only | |

## 4. Validation

- Oracles (macOS arm64 and Linux, no root):
  - python-control, scipy, scikit-image, h5py, astropy, pint: `uv run --with python-control --with scipy --with h5py --with astropy --with scikit-image --with pint python oracle.py` (uv from https://docs.astral.sh/uv/). Go harness: `exec.Command("uv", "run", ...)`, JSON on stdout, `float64` serialised with `repr` (17 significant digits) so values round-trip.
  - Octave control package (GPL, oracle only): `pkgx octave` or Homebrew formula; `octave --eval "pkg load control; ..." | jq`; slower start (seconds), batch the cases.
  - HDF5 tools: `h5dump -j`/`h5ls` from the HDF Group binary releases https://www.hdfgroup.org/downloads/hdf5/ (tarball, no root) or `uv run --with h5py`. NetCDF: `ncdump` from conda-forge via `pkgx` or `uv run --with netcdf4`. FITS: `fitsverify`, `astropy.io.fits`.
  - UDUNITS: `udunits2 -H` binary (build from source) or `pint`/`astropy` as secondary.
- Public corpora:
  - SLICOT benchmark collections https://www.slicot.org/20-site/126-benchmark-collections (that page lists model-reduction examples, Chahlaoui, Van Dooren, SLICOT Working Note 2002-2). Riccati: CAREX/DAREX (Abels, Benner, SLICOT Working Notes 1999-14/15) (unverified: not found on the pages fetched; see https://www.slicot.org/working-notes).
  - COMPleib: control design benchmark (Leibfritz) http://www.complib.de/ (reachable over http only).
  - DAE test sets: Test Set for IVP Solvers (Bari), https://archimede.uniba.it/~testset/testsetivpsolvers/ (stiff ODE and index-1 to index-3 DAEs with reference solutions to ~16 digits; licence unspecified, cite and generate, do not redistribute).
  - Modelica Standard Library examples https://github.com/modelica/ModelicaStandardLibrary (BSD-3-Clause).
  - HDF Group example files https://support.hdfgroup.org/documentation/hdf5/latest/_examples.html ; h5py test files; scipy test data `scipy/io/matlab/tests/data` (v4, v5, v7 MAT files, BSD-3); astropy test data (`astropy/io/fits/tests/data`, BSD-3); netcdf-c test files https://github.com/Unidata/netcdf-c ; NASA FITS sample files https://fits.gsfc.nasa.gov/fits_samples.html ; Matrix Market and SuiteSparse collections https://sparse.tamu.edu/ ; FLAC test suite https://github.com/ietf-wg-cellar/flac-test-files (CC0-1.0 per `gh api`).
  - Units: NIST SP 811 (Thompson, Taylor 2008) conversion factors https://www.nist.gov/pml/special-publication-811 ; pint test suite.
  - Marching cubes: Stanford volumetric data (bunny CT) https://graphics.stanford.edu/data/voldata/ (licence: research use, do not commit); analytic level sets (sphere, torus, Genus-2) generated.
- What "correct" means, by area:

| area | criterion | measure to report |
|---|---|---|
| Riccati / Lyapunov | relative residual ||A'X+XA-XBR^-1B'X+Q||_F / ||X||_F below 1e3 * n * eps; X symmetric to 1 ulp; closed-loop poles in left half-plane | worst residual over the benchmark set, median ratio versus SLICOT/scipy |
| poles/zeros | relative error vs reference after sorting (conjugate pairs matched), conditioning accounted | max error in ULPs of the eigenvalue condition number times eps |
| frequency response | complex relative error per point vs python-control | max over grid, report in ULPs of |H| |
| `expm`, c2d | relative error vs high-precision (mpmath) reference | norm-wise, 1e1 * eps |
| time response (LTI) | match to `expm`-based exact solution; ODE-based path judged with tolerance | max relative error versus exact |
| DAE | error vs reference solution of the Bari test set at stated tolerance; index-1 constraint residual | final-time error, constraint drift, steps, Newton iterations, event times to 1e-10 |
| file formats | exact equality of bytes read (ints, strings, raw floats); round-trip: Go write then h5py/astropy/scipy read gives identical arrays; Go read of their files identical | count of files passing; fuzz with `go test -fuzz` for panics; limit on allocations for hostile headers (see go-gfx decompression-bomb tests) |
| units | exact rational dimension vectors; conversion factors within 1 ulp of the CODATA decimal | table diff against NIST allascii.txt |
| mesh (marching cubes) | closed manifold for smooth level sets; Euler characteristic matches (sphere 2, torus 0); area within discretisation error vs analytic | vertices, faces, Euler characteristic, area error |
| rendering | golden-image SVG (byte equality) and raster (per-pixel tolerance, SSIM) | pixels differing above threshold; depth-order inversions on a known intersecting scene |
| Manipulate | latency from slider event to updated output; only dependent cells re-run | ms p50/p95; count of cells executed |

## 5. Proposed home and plan

Names checked with `gh api users/<name>` on 2026-10-17: free (404) = `go-controls`, `go-units`, `go-sciio`, `go-dae`, `go-sysid`, `go-codata`, `go-isosurface`, `go-scene3d`, `go-fitsio`, `go-astrofits`, `go-ncdf`. Taken (do NOT use): `go-hdf5` (third party, 2018 BSD-3), `go-ctrl`, `go-signal`, `go-mesh`, `go-audio`.

| work item | org/repo | depends on | size S/M/L/XL | priority 1-3 |
|---|---|---|---|---|
| Dense linear algebra prerequisites: real Schur form with eigenvalue reordering, QZ, `expm`, Sylvester/Lyapunov (Bartels-Stewart) | existing go-ndarray (or the linalg home proposed in the linear-algebra bib file) | none | XL | 1 |
| LTI model types (TF/ZPK/SS/FRD), poles, zeros, series/parallel/feedback | NEW org `go-controls`, repo `lti` | linalg item | L | 1 |
| CARE/DARE (Schur), LQR, Kalman gain, pole placement (KNV), c2d | `go-controls/design` (port SLICOT SB02OD, SB03MD, SB04MD under BSD-3) | linalg, `lti` | L | 1 |
| Time and frequency response plus Bode/Nyquist/Nichols/root-locus plots, margins | `go-controls/lti` and `go-plot/plot` (new `bode`, `nyquist`, `nichols` sub-packages) | `lti`, plot | M | 2 |
| System identification (ARX/ARMAX least squares, N4SID) | `go-controls/ident` | `lti`, SVD | L | 3 |
| DAE solver: IDA-type variable-order BDF, mass-matrix option for Radau, consistent initial conditions, events | existing go-numerics/integrate (or NEW `go-dae/ida` if kept separate from ODE) | integrate events, sparse/dense LU | XL | 2 |
| Index reduction (Pantelides, dummy derivatives) and block-diagram compiler (BLT sorting, algebraic loops) | `go-dae/blocks` | DAE solver | XL | 3 |
| Units: dimension vector, `Quantity`, UDUNITS-grammar parser, SI/US/derived database, offset units, CODATA 2022 generated table | NEW org `go-units`, repo `units` (+ `constants`) | none | M | 2 |
| NumPy `.npy/.npz`, Matrix Market, WAV, JSON adapters, xlsx to dataframe bridge | NEW org `go-sciio` (repos `npy`, `mtx`, `wav`) ; bridge in go-dataframe / go-ooxml | none | M | 2 |
| HDF5 reader (superblock v0-3, chunked, deflate/shuffle) and writer; decide: audit and wrap scigolib/hdf5 (MIT) or write own | `go-sciio/h5` | none | XL | 1 |
| MAT v4/v5/v7 and v7.3 (on top of h5) | `go-sciio/mat` | `h5` | M | 1 |
| NetCDF classic reader/writer, NetCDF-4 on top of h5 | `go-sciio/netcdf` | `h5`, `units` (CF) | L | 2 |
| FITS (images, BINTABLE, compressed images) | `go-sciio/fits` (fork or audit astrogo/fitsio BSD-3, archived) | none | L | 3 |
| FLAC decode (import mewkiz/flac Unlicense or port) | `go-sciio/flac` | none | S | 3 |
| 3D scene core: camera, perspective, lighting, depth-sorted SVG, z-buffer raster | `go-gfx/gfx` new package `scene3d` (or NEW `go-scene3d` if kept out of gfx) | gfx raster, svg | L | 1 |
| 3D plot API: `plot3`, `scatter3`, `surf`, `mesh`, `contour3`, `quiver3`, axes, colorbar, `view` | `go-plot/plot` new sub-package `plot3d` | scene core, colormaps | L | 1 |
| Marching cubes (Lewiner-corrected), marching tetrahedra, volume slices, isocaps | NEW org `go-isosurface` or `go-plot/plot3d/iso` | ndarray | M | 2 |
| Volume rendering (ray cast on CPU, WebGL later) | `go-plot/plot3d/volume` | scene core | M | 3 |
| Browser interactive viewer: WASM + WebGL2 scene JSON; rotate/zoom | `go-notebook` output renderer and `go-webengine`/`wasmdesk` reuse (check overlap with go-xrkit before starting) | scene core | L | 2 |
| Animation: frame loop, GIF/APNG, MP4 via go-avkit if present | `go-plot/plot` `anim` | raster, go-avkit | M | 3 |
| Manipulate/Dynamic: widget values, dependent-cell re-run over WebSocket, cheap-while-dragging mode | `go-notebook/notebook` (`serve`, `analysis/graph.go`) | notebook serve | L | 1 |
| Colormaps: viridis family, turbo, Crameri; verify what `plot/colors` has | `go-plot/plot/colors` | none | S | 2 |

Top order of attack: (1) linalg Schur/QZ/expm, (2) HDF5 + MAT, (3) 3D scene + plot3d, (4) Manipulate in notebook, (5) LTI/Riccati.

## 6. Risks and open questions

- Linear algebra is the blocker: go-ndarray's `linalg` package (v0.12.0 on GitHub) has LU, QR, Cholesky, SVD, eig, eigh and lstsq, but no real Schur with reordering, QZ or `expm` — which everything in 1a needs (an earlier draft read a stale local clone and said linalg had only MatMul/Dot/Inner/Outer). gonum/mat gives Eigen/SVD/QR but not Schur reordering (verify with `go doc`); LAPACK `dtrsen`, `dgges`, `dgees` would have to be ported (LAPACK is BSD-3-Clause-modified).
- SLICOT licence: the repository reports BSD-3-Clause, historically it was GPL/LGPL. Confirm per-routine headers in the pinned commit before porting; keep a record of the commit SHA.
- HDF5 is a large specification (object header versions, fractal heap, filters). Wrapping scigolib/hdf5 (MIT, 32 stars, young, partial write) is cheaper but couples to one maintainer; the audit must include hostile-input limits (allocation caps, recursion depth) in line with the fleet's decompression-bomb standard. MAT v7.3 and NetCDF-4 both inherit its bugs.
- `go-hdf5` is owned by a third party; do not reuse that name. `astrogo/fitsio` is archived on GitHub (moved to Codeberg): fork under our own licence terms (BSD-3 requires keeping the notice).
- DAE: index > 1 needs structural analysis and symbolic derivative of equations; without an expression system (no go-symbolic seen in the fleet) only index-1 and user-supplied reduced systems are feasible. Modelica compatibility is out of reach in a first wave; say so in the README.
- Units: UDUNITS licence field reads NOASSERTION; do not copy its database before reading its LICENSE. `parula` is proprietary; avoid. CODATA tables change every four years; generate from the NIST ASCII file with a reproducible tool and record the edition.
- Depth sorting by centroid gives wrong results for interpenetrating or large overlapping polygons; vector output needs either polygon splitting or a documented limitation. Raster z-buffer is exact but loses vectors. Golden SVG tests are fragile to float formatting: round to fixed decimals.
- Browser 3D requires WASM or JS shipped with the notebook; check the go-webengine/wasmdesk/go-xrkit memory entries to avoid duplicating an existing GPU or capture path. WebGPU availability is still uneven (Safari, Firefox versions); keep WebGL2 as baseline.
- Manipulate: reactive re-run changes the semantics of cells with side effects; need an explicit "pure" marker or re-run only widget dependents. Open question: persist widget state in the `.md` (a fenced block) or keep it session-local.
- Open: are the numerical oracles (python-control via slycot, Octave control) acceptable as test-time dependencies on CI runners, or must expected values be frozen into `testdata` JSON generated once (preferred: frozen with the generating script and versions recorded).
- Items still marked (unverified) (Chernyaev 1995, Bierman ISBN, CAREX/DAREX numbers) need a fetch before being relied on.

Verified on 2026-10-17: all DOIs resolved on Crossref (title, year, first author; 5 fixed); MATLAB/Mathematica names (Wolfram ref pages by HTTP status, MathWorks via WebFetch for the doubtful ones); licences of every listed repository by `gh api` and LICENSE text; go-scikit/fleet statements against GitHub main (org repo lists, trees, file contents); org-name availability; reference URLs (HTTP status).
