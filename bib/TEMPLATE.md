# <Domain> — bibliography

One file per domain. Facts only, each with a source (URL or DOI). Terse: tables, not prose.

## 1. Feature inventory

What MATLAB (base + toolboxes) and Mathematica offer in this domain, by function name, grouped
by capability. Mark what go-scikit already has (repo and symbol) and what it lacks.

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|

## 2. Reference algorithms

Per capability: the algorithm(s) the reference tools use or that the literature recommends,
with the paper (authors, year, title, DOI) and why it is the right choice (accuracy, cost).

## 3. Reference implementations

| implementation | language | licence | portable to Go? (BSD/MIT/Apache = port; GPL/LGPL/proprietary = oracle only) | notes |
|---|---|---|---|---|

## 4. Validation

- Oracle tools to compare against (name, version, how to obtain it on macOS arm64 and Linux
  without root: uv, pkgx or a release binary), and how to drive each from a Go harness
  (exec + JSON on a pipe).
- Public test corpora and suites (name, URL, size, licence).
- What "correct" means in this domain (ULPs, relative residual, backward error, exact
  equality, certified bounds) and the measure to report.

## 5. Proposed home and plan

| work item | org/repo (existing, or NEW org to create) | depends on | size S/M/L/XL | priority 1–3 |
|---|---|---|---|---|

## 6. Risks and open questions
