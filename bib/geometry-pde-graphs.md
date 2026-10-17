# Computational geometry, meshing/PDEs, graphs — bibliography

One file per domain. Facts only, each with a source (URL or DOI). Terse: tables, not prose.

Provenance: items marked [fetched] were read on 2026-10-17 from the URL given. All DOIs in section 2
were resolved via Crossref on 2026-10-17 (title, year, first author match).

## 1. Feature inventory

Existing fleet code (GitHub main, checked 2026-10-17): go-ndarray/ndarray, go-numerics/{optimize,integrate,special},
go-plot/plot, go-gfx/{gfx,qr}. Tree searches (sparse, csr, graph, delaun, mesh, hull, geom, voronoi, pde, fem,
polygon) find no geometry, mesh, PDE, sparse or graph code. go-graphdrawing (GitHub org)
holds only `.github`, `go-graphdrawing.github.io`, `brand` [`gh repo list go-graphdrawing`]: no code.
go-tex repos: engine, pdfrender, tex, texmf, math (no graph/tikz/pgf repo; `texmf` contents not audited).
go-numerics/integrate has `bdf.go`, `radau.go` (stiff ODE). So: nothing today for geometry/PDE/graphs.

### (a) Computational geometry

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| Delaunay 2D/3D/nD | `delaunay`, `delaunayTriangulation`, `triangulation` | `DelaunayMesh` | none | all |
| Voronoi | `voronoi`, `voronoiDiagram` | `VoronoiMesh` | none | all |
| convex hull | `convhull`, `convhulln` | `ConvexHullMesh` | none | all |
| alpha shapes / boundary | `alphaShape`, `boundary` | `BoundaryMesh`, `RegionBoundary` | none | all |
| nearest neighbour | `dsearchn`, `knnsearch` (Stats TB) | `Nearest`, `NearestFunction` | none | all (gonum has kd/vp trees, see 3) |
| point in polygon | `inpolygon`, `isinterior` | `RegionMember` | none | all |
| polygon booleans | `polyshape` (`union`, `intersect`, `subtract`) | `RegionUnion`, `RegionIntersection`, `RegionDifference` | none | all |
| measures | `area`, `perimeter`, `volume` | `RegionMeasure`, `Area`, `Volume` | none | all |
| scattered interpolation | `griddata`, `scatteredInterpolant` | `Interpolation` on scattered data | none | all |

Sources: https://www.mathworks.com/help/matlab/computational-geometry.html [fetched: delaunay,
delaunayTriangulation, triangulation, voronoi, convhull, alphaShape, boundary, dsearchn, griddata,
inpolygon, polyshape]; https://reference.wolfram.com/language/guide/ComputationalGeometry.html
[fetched: Nearest, NearestFunction, region tools]; every Wolfram name in (a) resolves at reference.wolfram.com/language/ref/<Name>.html (HTTP 200, 2026-10-17).

### (b) Meshing and PDEs

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| model/geometry object | `createpde`, `geometryFromEdges`, `multicuboid` | `Region`, `DiscretizeRegion` | none | all |
| mesh generation | `generateMesh` (Hmax, Hgrad) | `ToElementMesh`, `ToBoundaryMesh`, `ToGradedMesh` | none | all |
| FEM solve (elliptic/parabolic/hyperbolic) | `solvepde`, `specifyCoefficients`, `applyBoundaryCondition` | `NDSolve` (FEM spatial discretisation) | none | all |
| eigenproblems | `solvepdeeig` | `NDEigensystem` | none | all |
| structural/thermal | `structuralModel`, `thermalModel` | `NDSolve` with `DirichletCondition`, `NeumannValue`, `PeriodicBoundaryCondition` | none | all |
| 1D method of lines | `pdepe` | `NDSolve` (MOL) | none | all |

Sources: https://www.mathworks.com/help/pde/index.html [fetched: structural, heat, electromagnetics,
general PDEs, geometry/mesh workflow]; all MATLAB names in (b) resolve under https://www.mathworks.com/help/pde/ug/
(HTTP 200: createpde, geometryFromEdges, multicuboid, generateMesh, solvepde, specifyCoefficients,
applyBoundaryCondition, solvepdeeig, structuralModel, thermalModel) and `pdepe` under matlab/ref/pdepe.html;
https://reference.wolfram.com/language/FEMDocumentation/guide/FiniteElementMethodGuide.html [fetched:
NDSolve, NDEigensystem, ToElementMesh, ToBoundaryMesh, ToGradedMesh, DirichletCondition, NeumannValue,
PeriodicBoundaryCondition]. ToElementMesh/ToBoundaryMesh/ToGradedMesh live in the `NDSolve`FEM`` context
(docs under /language/FEMDocumentation/ref/).

### (c) Graphs and networks

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| construction, adjacency/incidence/Laplacian | `graph`, `digraph`, `adjacency`, `incidence`, `laplacian` | `Graph`, `AdjacencyMatrix`, `KirchhoffMatrix` | none | all (gonum/graph covers structures) |
| shortest paths | `shortestpath`, `shortestpathtree`, `distances` | `FindShortestPath`, `GraphDistanceMatrix` | none | all |
| max flow / min cut | `maxflow` | `FindMaximumFlow`, `FindMinimumCut` | none | all |
| spanning tree | `minspantree` | `FindSpanningTree` | none | all |
| connectivity | `conncomp`, `biconncomp`, `condensation` | `ConnectedComponents`, `ConnectedGraphComponents`, `FindVertexCut` (no biconnected-components function) | none | all |
| centrality | `centrality` | `BetweennessCentrality`, `PageRankCentrality`, ... | none | all |
| traversal / cycles / toposort | `bfsearch`, `dfsearch`, `allcycles`, `cyclebasis`, `toposort` | `FindCycle`, `TopologicalSort` | none | all |
| isomorphism | `isomorphism` | `IsomorphicGraphQ`, `FindGraphIsomorphism` | none | all |
| matching | (none in base list) | `FindIndependentEdgeSet` | none | all |
| community detection | (none in base list) | `FindGraphCommunities` | none | all |
| layout / drawing | `plot` + `layout` (force, layered, circle, subspace) | `GraphPlot`, `GraphPlot3D`, `GraphLayout` | none | all |

Sources: https://www.mathworks.com/help/matlab/graph-and-network-algorithms.html [fetched: graph,
digraph, shortestpath, shortestpathtree, distances, maxflow, minspantree, conncomp, biconncomp,
centrality, bfsearch, dfsearch, allcycles, cyclebasis, hascycles, condensation, toposort, isomorphism,
plot, highlight, layout, adjacency, incidence, laplacian]. Wolfram names in this table resolve at
/language/ref/<Name>.html (HTTP 200, 2026-10-17); `BiconnectedComponents` returned 404 and was replaced.

## 2. Reference algorithms

| capability | algorithm | paper | why |
|---|---|---|---|
| exact predicates | adaptive-precision orient2d/orient3d/incircle/insphere | Shewchuk 1997, "Adaptive Precision Floating-Point Arithmetic and Fast Robust Geometric Predicates", Discrete Comput Geom 18:305-363, doi:10.1007/PL00009321 | exact sign at near-float cost; prerequisite of every item below |
| convex hull nD | Quickhull | Barber, Dobkin, Huhdanpaa 1996, ACM TOMS 22(4):469-483, doi:10.1145/235815.235821 | used by Qhull, SciPy, MATLAB; handles nD, merged facets for degeneracy |
| 2D Delaunay | divide and conquer / incremental flip | Guibas, Stolfi 1985, ACM TOG 4(2):74-123, doi:10.1145/282918.282923 | O(n log n); quad-edge structure |
| Voronoi | sweep line, or dual of Delaunay | Fortune 1987, Algorithmica 2:153-174, doi:10.1007/BF01840357 | O(n log n) |
| alpha shapes | filtration of Delaunay | Edelsbrunner, Mücke 1994, ACM TOG 13(1):43-72, doi:10.1145/174462.156635 | MATLAB `alphaShape` basis |
| kd-tree / NN | k-d tree | Bentley 1975, CACM 18(9):509-517, doi:10.1145/361002.361007 | exact NN in low dimension |
| polygon booleans | scan-line clipping | Vatti 1992, CACM 35(7):56-63, doi:10.1145/129902.129906 | handles holes, self-touching input |
| constrained Delaunay + refinement | Ruppert refinement; Chew | Ruppert 1995, J. Algorithms 18(3):548-585, doi:10.1006/jagm.1995.1021; Shewchuk 2002, Comput. Geom. 22:21-74, doi:10.1016/S0925-7721(01)00047-5 | min-angle guarantee (about 20.7 degrees); basis of Triangle |
| 3D tetrahedral meshing | constrained Delaunay + refinement | Si 2015, ACM TOMS 41(2):11, doi:10.1145/2629697 | TetGen method |
| general mesh generation | Delaunay/frontal, many | Geuzaine, Remacle 2009, IJNME 79(11):1309-1331, doi:10.1002/nme.2579 [citation fetched from https://gmsh.info/] | Gmsh |
| FEM theory | Lagrange P1/P2, variational form | Brenner, Scott, "Mathematical Theory of Finite Element Methods", Springer, doi:10.1007/978-0-387-75934-0; Ern, Guermond, "Theory and Practice of Finite Elements", Springer 2004, doi:10.1007/978-1-4757-4355-5 | error estimates, a priori bounds |
| FEM software design | automated assembly from forms | Logg, Mardal, Wells (eds.), "Automated Solution of Differential Equations by the Finite Element Method" (FEniCS book), Springer 2012, doi:10.1007/978-3-642-23099-8 | design of element/dofmap/assembler split |
| triangle quadrature | symmetric Gauss rules | Dunavant 1985, IJNME 21(6):1129-1148, doi:10.1002/nme.1620210612; Xiao, Gimbutas 2010, Comput. Math. Appl. 59(2):663-676, doi:10.1016/j.camwa.2009.10.027 | tabulated, high degree, positive weights |
| finite volumes | conservative flux, Godunov | LeVeque 2002, "Finite Volume Methods for Hyperbolic Problems", CUP, doi:10.1017/CBO9780511791253 | conservation laws |
| method of lines | spatial FD + stiff ODE | Skeel, Berzins 1990, SIAM J. Sci. Stat. Comput. 11(1):1-32, doi:10.1137/0911001 (the `pdepe` basis per MathWorks docs) | reuses go-numerics/integrate |
| shortest paths | Dijkstra; Bellman-Ford; Floyd-Warshall; Johnson | Dijkstra 1959, Numer. Math. 1:269-271, doi:10.1007/BF01386390; Cormen et al., "Introduction to Algorithms", 4th ed., MIT Press 2022, ISBN 9780262046305 | textbook |
| max flow | push-relabel; Dinic | Goldberg, Tarjan 1988, JACM 35(4):921-940, doi:10.1145/48014.61051; Ahuja, Magnanti, Orlin, "Network Flows", Prentice Hall 1993, ISBN 013617549X | practical speed |
| matching | Hopcroft-Karp (bipartite); blossom (general) | Hopcroft, Karp 1973, SIAM J. Comput. 2(4):225-231, doi:10.1137/0202019; Edmonds 1965, Can. J. Math. 17:449-467, doi:10.4153/CJM-1965-045-4 | |
| strong components | Tarjan | Tarjan 1972, SIAM J. Comput. 1(2):146-160, doi:10.1137/0201010 | linear |
| betweenness | Brandes | Brandes 2001, J. Math. Sociol. 25(2):163-177, doi:10.1080/0022250X.2001.9990249 | O(nm) not O(n^3) |
| communities | Louvain; Leiden; Girvan-Newman | Blondel et al. 2008, J. Stat. Mech. P10008, doi:10.1088/1742-5468/2008/10/P10008; Traag, Waltman, van Eck 2019, Sci. Rep. 9:5233, doi:10.1038/s41598-019-41695-z; Girvan, Newman 2002, PNAS 99:7821, doi:10.1073/pnas.122653799 | Leiden fixes Louvain's disconnected communities |
| isomorphism | VF2; canonical labelling | Cordella et al. 2004, IEEE TPAMI 26(10):1367-1372, doi:10.1109/TPAMI.2004.75; McKay, Piperno 2014, J. Symb. Comput. 60:94-112, doi:10.1016/j.jsc.2013.09.003 | VF2 for matching, nauty for canonical form |
| force-directed layout | spring-electrical; stress | Fruchterman, Reingold 1991, Softw. Pract. Exper. 21(11):1129-1164, doi:10.1002/spe.4380211102; Kamada, Kawai 1989, Inf. Process. Lett. 31(1):7-15, doi:10.1016/0020-0190(89)90102-6 | |
| layered layout | Sugiyama; dot | Sugiyama, Tagawa, Toda 1981, IEEE Trans. SMC 11(2):109-125, doi:10.1109/TSMC.1981.4308636; Gansner, Koutsofios, North, Vo 1993, IEEE TSE 19(3):214-230, doi:10.1109/32.221135 | de facto DAG drawing |
| orthogonal layout | topology-shape-metrics | Tamassia 1987, SIAM J. Comput. 16(3):421-444, doi:10.1137/0216030 | |

## 3. Reference implementations

Licences marked [gh] come from the GitHub API on 2026-10-17; [file] = read from the project's licence file the same day.

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| gonum/gonum (`graph/...`, `spatial/{kdtree,vptree,r2,r3,barneshut}`) | Go | BSD-3-Clause [gh] | yes, import directly | graph: path (Dijkstra, A*, Bellman-Ford, Floyd-Warshall, Johnson, Yen, spanning_tree), topo (Tarjan, cliques, cycles), community (Louvain, Leiden), layout (Eades, Isomap), network (`MaxFlowDinic` in maximum_flow.go, Betweenness/EdgeBetweenness (Brandes), PageRank, HITS), `flow` = dominators only [files read via gh api on main]. No matching, no push-relabel, no Delaunay/hull, no Sugiyama, no sparse matrices (gonum/mat) |
| paulmach/orb | Go | MIT [gh] | yes | 2D geometry types, clipping |
| twpayne/go-geom | Go | BSD-2-Clause [gh] | yes | geometry types, WKB/GeoJSON |
| golang/geo | Go | Apache-2.0 [gh] | yes | S2 spherical geometry |
| fogleman/delaunay | Go | MIT [gh] | yes | 2D Delaunay (port of delaunator), last push 2023, no robust predicates claimed |
| Shewchuk predicates.c | C | public domain ("This code is in the public domain" [fetched https://www.cs.cmu.edu/~quake/robust.html]) | yes | https://www.cs.cmu.edu/~quake/robust.html |
| Triangle 1.6 | C | "freely available, copyrighted, may not be sold or included in commercial products without a license" [fetched https://www.cs.cmu.edu/~quake/triangle.html] | NO (non-free): oracle only | CDT, Ruppert/Chew refinement, Voronoi; Wilkinson Prize 2003 |
| Qhull | C/C++ | custom "Qhull, Copyright (c) 1993-2020" with five redistribution conditions, not a standard SPDX [fetched https://github.com/qhull/qhull/blob/master/COPYING.txt] | read the five conditions before any port; oracle meanwhile | used by SciPy, MATLAB lineage |
| TetGen 1.5.0 (Nov 2013, H. Si, WIAS) | C++ | dual: AGPL-3.0-or-later or commercial [fetched https://wias-berlin.de/software/tetgen/1.5/FAQ-license.html] | NO: oracle only | |
| Gmsh 4.15.2 (2026-03-24) | C++ | GPL-2.0-or-later with linking exception [fetched https://gmsh.info/] | oracle only; call as external binary | |
| CGAL | C++ | per-file GPL-3.0+ or LGPL-3.0+, plus commercial from GeometryFactory [file Installation/LICENSE] | oracle only | |
| SciPy `spatial` | Python/C | BSD-3-Clause [gh] | algorithms free to port; code wraps Qhull | oracle for hull/Delaunay/Voronoi/KDTree |
| NetworkX | Python | BSD-3-Clause [file LICENSE.txt "3-clause BSD"] | port allowed | oracle for all graph algorithms |
| igraph C core | C | GPL-2.0 [gh] | NO: oracle only | |
| Boost.Graph | C++ | BSL-1.0 [gh] | yes (permissive) | |
| Graphviz | C | EPL-2.0 [file, gitlab.com/graphviz/graphviz LICENSE; the GitHub repo is an empty archived mirror] | EPL: oracle only |  `dot`, `neato`, `fdp`, `sfdp`, `-Tjson` output |
| deal.II | C++ | dual: Apache-2.0 with LLVM exception, and separately LGPL-2.1-or-later [file LICENSE.md]; some files under bundled/, examples/, tests/ differ | Apache-2.0 side permits a port; check per-file SPDX header | |
| MFEM | C++ | BSD-3-Clause [gh] | yes | clean design reference for FE elements and assembly |
| FEniCS / DOLFINx | C++/Python | LGPL-3.0 [gh] | oracle only | MMS tests, UFL forms |

## 4. Validation

Oracles (drive from a Go harness by exec, JSON on a pipe; no root):
- SciPy (`scipy.spatial.{ConvexHull,Delaunay,Voronoi,cKDTree}`): `uv run --with scipy python harness.py`; works on macOS arm64 and Linux wheels.
- NetworkX: `uv run --with networkx python ...`; node-link JSON in and out.
- Octave (`delaunay`, `convhull`, `inpolygon`, `polyshape` partial): `pkgx octave` or distro package; mirrors MATLAB semantics. MATLAB itself is unavailable.
- Graphviz `dot -Tjson` for layered layout: crossing counts, ranks and bounding boxes can be compared; coordinates cannot be required to match.
- Triangle and Qhull binaries as oracles only (licence above): `qconvex`, `qdelaunay` via release binaries or pkgx; `triangle` built from Netlib source locally and never vendored.
- FEniCS/DOLFINx (conda-forge or docker; no pure uv wheel, so Linux-only lane) for P1/P2 Poisson solutions on a shared mesh.

Corpora (licences not verified this session unless stated):
- SNAP, https://snap.stanford.edu/data/ [fetched]: ego-Facebook 4,039 nodes up to com-Friendster over 65 M nodes; ground-truth communities; road networks. Check per-dataset terms.
- DIMACS Implementation Challenges, http://archive.dimacs.rutgers.edu/Challenges/ [fetched]: 1st network flows and matching, 5th priority queues/point sets, 6th near-neighbour search, 9th shortest paths, 10th partitioning/clustering.
- Method of manufactured solutions (MMS): pick u, derive f and boundary data analytically, no external corpus. Report observed convergence order on a refinement sequence (expect h^(p+1) in L2, h^p in H1 for Pp).
- Shewchuk predicates: degenerate input families (cocircular grids, collinear sets) generated in-repo.

What "correct" means:
- Predicates: exact sign, compared against big.Rat arithmetic on random and near-degenerate input; zero mismatches.
- Hull/Delaunay: combinatorial validity (Euler relation, empty-circumsphere test with exact predicates), vertex set equal to oracle up to degeneracies; report count of differing simplices only for cocircular input where non-uniqueness is legitimate.
- kd-tree/NN: exact equality of index sets with brute force; ties broken by index.
- Polygon booleans: area of result vs oracle within 1e-12 relative, plus topological validity.
- Mesh quality: min angle, max aspect ratio, area conservation; report distribution, not only the minimum.
- PDE: observed convergence order within 0.1 of theory, relative L2 error against exact solution, linear-solve relative residual; never compare against an oracle's mesh node by node.
- Graphs: exact equality for integer-weight shortest paths, flow value, matching size, component partitions; modularity within tolerance for Louvain/Leiden (non-deterministic: fix seed, compare modularity, not labels); layouts: no coordinate equality, only edge crossings, overlap count and stress.

## 5. Proposed home and plan

Org names below were checked with `gh api users/<name>`: 404 means free at 2026-10-17. Taken (exist): go-mesh, go-networks. Not derived from Python projects.

| work item | org/repo | depends on | size S/M/L/XL | priority 1–3 |
|---|---|---|---|---|
| exact predicates (port of Shewchuk, with big.Rat oracle test) | NEW org `go-polytope` / `predicates` | none | S | 1 |
| convex hull 2D/3D/nD (Quickhull), Delaunay 2D/3D, Voronoi, alpha shapes | `go-polytope/hull`, `delaunay` | predicates, go-ndarray | M-L | 1 |
| point-in-polygon, polygon booleans, distances, measures | `go-polytope/polygon` | predicates | M | 2 |
| adopt gonum kd-tree/vp-tree behind an ndarray facade, add MATLAB-style `knnsearch` API | `go-polytope/nearest` | gonum (BSD-3), go-ndarray | S | 2 |
| graph gap-fill over gonum/graph (already has Dinic max flow, Brandes betweenness, PageRank, HITS): min cut extraction, push-relabel if needed, matching (Hopcroft-Karp, blossom), isomorphism (VF2), remaining centralities, MATLAB-style `graph`/`digraph` facade | NEW org `go-graphalgo` | gonum/graph | M-L | 1 |
| graph layout: force-directed first, then Sugiyama (dot-like), then orthogonal; render via go-gfx and go-plot | existing `go-graphdrawing` (empty) | go-graphalgo, go-gfx | XL | 2 |
| constrained Delaunay + Ruppert refinement 2D, then 3D | NEW org `go-meshgen` | go-polytope | L (2D), XL (3D) | 2 |
| Lagrange P1/P2 elements, quadrature tables, assembly, BCs; Poisson then heat/wave/elasticity/Stokes | NEW org `go-finiteelement` | go-meshgen, a sparse solver (none in go-ndarray/go-numerics/gonum mat on main, 2026-10-17) | XL | 2 |
| MMS convergence harness (shared test kit) | `go-finiteelement/mms` | none | S | 1 |
| 1D method of lines (`pdepe` analogue) | `go-numerics/integrate` or `go-finiteelement/mol` | stiff ODE solver in go-numerics/integrate (`bdf.go`, `radau.go` on main) | M | 3 |
| finite differences / volumes on structured grids | `go-finiteelement/fdm` | go-ndarray | M | 3 |

## 6. Risks and open questions

- Licences: Triangle is non-free; TetGen is AGPL-3.0+, CGAL GPL/LGPL-3+, Graphviz EPL-2.0 (oracle only). Qhull's custom terms: read the five conditions before any port. deal.II is dual Apache-2.0 (LLVM exception)/LGPL-2.1+, so its Apache side may be ported (check per-file SPDX). Safe sources: Shewchuk's predicates (public domain) and Guibas–Stolfi's papers, MFEM (BSD-3), gonum (BSD-3), Boost (BSL), deal.II (Apache side).
- Robustness is the main technical risk: Delaunay, hull and booleans silently produce invalid topology without exact predicates; build and test predicates first and never ship hull before them.
- Scope: 3D meshing and general FEM (Stokes, elasticity, 3D) are multi-year; the MATLAB PDE Toolbox surface is far larger than the Poisson-first plan covers.
- Sparse linear algebra is a hidden dependency for FEM; not found in go-ndarray, go-numerics or gonum/mat on main.
- gonum overlap: graph work should extend or wrap gonum, not duplicate it; gonum already has Dinic max flow and Brandes betweenness; confirm the maintainers' stance on adding matching upstream.
- Layout output is not comparable across tools by coordinates; the acceptance test must be metric-based.
- Open: does go-tex/texmf bundle pgf graphdrawing Lua algorithms (would be reusable specs or tests)? Not audited.

Verified on 2026-10-17: 34 DOIs resolved via Crossref (title/year/first author all match, none changed); MATLAB and Mathematica names checked by HTTP status of their doc pages; licences via gh API and licence files; fleet/gonum contents via GitHub main trees; org names go-polytope, go-graphalgo, go-meshgen, go-finiteelement free, go-mesh/go-networks taken.
