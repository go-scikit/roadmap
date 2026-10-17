# Interpolation/approximation, signal processing, image processing — bibliography

One file per domain. Facts only, each with a source (URL or DOI). Terse: tables, not prose.

Source status: "go-scikit today" re-checked against GitHub main (API, not clones) on 2026-10-17;
licences re-checked with the GitHub licence API / LICENSE files the same day; every DOI resolved
via Crossref (section 6).

## 1. Feature inventory

Function names: MATLAB (base + Curve Fitting, Signal Processing, Wavelet, Image Processing
toolboxes) — https://www.mathworks.com/help/matlab/ , /curvefit/ , /signal/ , /wavelet/ , /images/ .
Mathematica — https://reference.wolfram.com/language/guide/ (Interpolation, SignalProcessing,
ImageProcessing, Wavelets guides).

Local state: `go-images/images` (exported funcs listed), `go-fft/fft` (FFT/RFFT/FFTN/FFT2/RFFT2,
Hann/Hamming/Blackman/BlackmanHarris/Bartlett, FFTFreq/RFFTFreq, PSD, Spectrogram),
`go-numerics/special` (has `poly.go`; not a polynomial-fit package), `go-ndarray/ndarray`
(creation/indexing/manipulation/linalg). `go-fft/fft` has the DCT and DST (types I–IV, N-D:
`dct.go`, `dct32.go`, `dctn.go`, `dst.go` — checked by the coordinator on 2026-10-17; an
earlier draft of this file wrongly said they were missing).

### (a) Interpolation and approximation

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| 1-D interpolation (linear, nearest, pchip, cubic, spline, makima) | `interp1`, `pchip`, `spline`, `makima` | `Interpolation`, `InterpolatingFunction` | none | all |
| N-D regular grid | `interp2`, `interp3`, `interpn`, `griddedInterpolant` | `Interpolation` (tensor) | `images.Warp` interpolates pixels only (`Interp` modes) | n-D array interpolator |
| Scattered data | `griddata`, `scatteredInterpolant` (Delaunay, natural neighbour) | `Interpolation` on scattered, `NearestFunction` | none | Delaunay + natural neighbour |
| RBF | none built in | none built in (`FunctionInterpolation` is not RBF) | none | RBF solver |
| Splines | `spline`, `ppval`, `mkpp`, `spapi`, `spaps`, `csaps` (Curve Fitting) | `BSplineFunction`, `BSplineCurve` | none | B-spline basis, smoothing spline |
| Polynomials | `polyfit`, `polyval`, `roots`, `conv`, `deconv`, `polyder`, `polyint`, `poly` | `Fit`, `InterpolatingPolynomial`, `NSolve`, `Roots` | none (`special/poly.go` internal evaluators) | polyfit/roots/eval |
| Chebyshev approximation | Chebfun (add-on), `chebfun` | `ChebyshevT` | none | Chebfun-style adaptive |
| Rational approximation | `rationalfit` (RF Toolbox), Chebfun `aaa` | `PadeApproximant` (`MiniMaxApproximation` is only in the legacy `FunctionApproximations`` package) | none | AAA, Padé, minimax |
| Curve/surface fitting | `fit`, `lsqcurvefit` (Optimization Toolbox), `curveFitter` | `FindFit`, `NonlinearModelFit` | `go-numerics/optimize` (`CurveFit`, `CurveFitModel` in `curvefit.go`; `leastsquares.go`, `lmdif.go`, `trf.go`) | fit front-end for splines/models |

### (b) Signal processing

| capability | MATLAB | Mathematica | go-scikit today | gap |
|---|---|---|---|---|
| IIR design | `butter`, `cheby1`, `cheby2`, `ellip`, `besself`, `buttord`, `zp2sos` | `ButterworthFilterModel`, `Chebyshev1FilterModel`, `Chebyshev2FilterModel`, `EllipticFilterModel`, `BesselFilterModel` | none | all, plus ZPK/SOS types |
| FIR design | `fir1`, `fir2`, `firls`, `firpm`, `kaiserord`, `designfilt` | `LeastSquaresFilterKernel`, `LowpassFilter` | windows only (`fft.Hann` etc.) | window/LS/Parks–McClellan |
| Filtering | `filter`, `filtfilt`, `sosfilt`, `conv`, `conv2`, `xcorr`, `fftfilt` | `LowpassFilter`, `ListConvolve`, `ListCorrelate` | `images.Convolve` (2-D, `Kernel`) | 1-D lfilter/filtfilt/sosfilt; conv/corr via FFT |
| Spectral estimation | `periodogram`, `pwelch`, `pmtm`, `spectrogram`, `stft`, `pburg`, `pcov` | `Periodogram`, `Spectrogram`, `SpectrogramArray` | `fft.PSD`, `fft.Spectrogram`, windows | Welch, multitaper (DPSS), STFT/ISTFT, more windows (Kaiser, Tukey, Gaussian, DPSS) |
| Resampling | `resample`, `decimate`, `interp`, `upfirdn`, `rat` | `TimeSeriesResample` | `images.Resize` (2-D only) | resample_poly, decimate, upfirdn |
| Analytic signal | `hilbert`, `envelope` | `HilbertTransform` | FFT exists | hilbert |
| Peak finding | `findpeaks`, `islocalmax` | `FindPeaks` | none | all |
| Wavelets | `dwt`, `idwt`, `wavedec`, `cwt`, `icwt`, `modwt`, `wpdec`, `wfilters` (Wavelet Toolbox) | `DiscreteWaveletTransform`, `ContinuousWaveletTransform`, `StationaryWaveletTransform`, wavelet families | none | DWT/SWT/CWT/WPT + families |
| Audio I/O | `audioread`, `audiowrite`, `sound` | `Import`/`Export` WAV, `Audio` | none | WAV/AIFF reader/writer (stdlib-only) |

### (c) Image processing

| capability | MATLAB (Image Processing) | Mathematica | go-scikit today (`go-images/images`) | gap |
|---|---|---|---|---|
| Linear filtering | `imfilter`, `imgaussfilt`, `fspecial`, `conv2` | `GaussianFilter`, `ImageConvolve`, `MeanFilter` | `Convolve`, `BoxBlur`, `GaussianBlur`, `Sharpen`, `UnsharpMask` | separable/FFT convolution, nD, float images (current API is `*image.RGBA` 8-bit) |
| Rank filters | `medfilt2`, `ordfilt2`, `stdfilt` | `MedianFilter`, `RankedMax`/`RankedMin` (list-level) | `Median` | general rank filter |
| Edges | `edge` (Sobel, Prewitt, Roberts, LoG, Canny) | `EdgeDetect`, `GradientFilter` | `Sobel`, `SobelX/Y`, `SobelMag`, `Prewitt`, `Scharr`, `Laplacian`, `Canny` | Roberts, LoG, hysteresis on float |
| Morphology | `imerode`, `imdilate`, `imopen`, `imclose`, `bwmorph`, `imreconstruct`, `bwdist` | `Erosion`, `Dilation`, `Opening`, `Closing`, `DistanceTransform`, `MorphologicalComponents` | `Erode`, `Dilate`, `Open`, `Close` (integer radius) | structuring elements, binary/grey, reconstruction, distance transform, skeleton, labelling |
| Geometric transforms | `imresize`, `imrotate`, `imtransform`, `imwarp`, `affine2d`, `projective2d` | `ImageResize`, `ImageRotate`, `ImageTransformation`, `GeometricTransformation` | `Resize` (modes), `Warp`, `Affine` (Identity/Translation/Scaling/Rotation/Then/Invert), `Rotate`, flips/transposes, `Crop` | projective/homography, TPS, polar/log-polar |
| Radon / Hough | `radon`, `iradon`, `hough`, `houghpeaks`, `houghlines`, `imfindcircles` | `Radon`, `InverseRadon`, `ImageLines` | none | all |
| Segmentation | `watershed`, `grayconnected`, `activecontour`, `imsegkmeans`, `graythresh`, `superpixels` | `WatershedComponents`, `Binarize`, `FindThreshold` | `OtsuThreshold`, `Otsu`, `Threshold` | watershed, region growing, active contours, k-means, adaptive threshold |
| Feature detection | `detectHarrisFeatures`, `detectSURFFeatures`, `detectORBFeatures`, `corner` (Computer Vision Toolbox) | `ImageKeypoints`, `ImageCorners`, `FindGeometricTransform` | none | Harris/Shi–Tomasi, FAST, blob (LoG/DoG), descriptors |
| Colour spaces | `rgb2hsv`, `rgb2lab`, `rgb2ycbcr`, `rgb2xyz` | `ColorConvert` (Lab, Luv, XYZ, HSB, ...) | `RGBToHSV`, `HSVToRGB`, `Grayscale` | Lab/Luv/XYZ/YCbCr, ICC |
| Registration | `imregister`, `imregcorr`, `imregtform`, `imregdemons` | `ImageAlign`, `FindGeometricTransform` | none | phase correlation, intensity-based, demons |
| Histogram ops | `imhist`, `histeq`, `adapthisteq`, `imadjust` | `HistogramTransform`, `ImageAdjust` | `AdjustBrightness`, `AdjustContrast`, `Invert` | histogram, equalisation, CLAHE |

Ruby binding: `go-embedded-ruby/ruby/internal/vm/images.go` (+ `images_bind.go`, `image.go`)
exposes `blur`, `gaussian_blur`, `median`, `canny`, `erode`, `dilate`, `sobel`; the file comment
(line 166) says these are "the scikit-image (go-images) surface" — it wraps `go-images/images`,
it does not hold separate algorithms.

## 2. Reference algorithms

| capability | algorithm(s) | reference | why |
|---|---|---|---|
| 1-D splines | B-spline basis via Cox–de Boor recurrence; not-a-knot cubic | de Boor, *A Practical Guide to Splines*, Springer 2001, ISBN 978-0-387-95366-3 | stable basis, banded solve O(n) |
| Shape-preserving | PCHIP | Fritsch & Carlson, SIAM J. Numer. Anal. 17(2), 1980, doi:10.1137/0717021 | monotone, local |
| Shape-preserving | Akima | Akima, J. ACM 17(4), 1970, doi:10.1145/321607.321609 | local, low overshoot |
| Smoothing spline | Reinsch | Reinsch, Numer. Math. 10, 1967, doi:10.1007/BF02162161 | O(n) banded |
| FITPACK | Dierckx curve/surface fitting with knot placement | Dierckx, *Curve and Surface Fitting with Splines*, OUP 1993, ISBN 978-0-19-853440-2 | what SciPy `splrep`/`splev` call |
| Scattered | Delaunay by Bowyer–Watson | Bowyer, Comput. J. 24(2), 1981, doi:10.1093/comjnl/24.2.162; Watson, same issue, doi:10.1093/comjnl/24.2.167 | incremental, simple; needs robust predicates |
| Scattered | natural neighbour | Sibson 1981, "A brief description of natural neighbour interpolation" (Interpreting Multivariate Data, Wiley) — no DOI | C1 except at data |
| Barycentric | barycentric Lagrange | Berrut & Trefethen, SIAM Review 46(3), 2004, doi:10.1137/S0036144502417715 | O(n) evaluation, stable |
| Chebyshev | adaptive Chebyshev expansions, FFT-based coefficients, colleague-matrix roots | Trefethen, *Approximation Theory and Approximation Practice*, SIAM 2013 (ext. ed. 2019), ISBN 978-1-61197-593-2; Chebfun https://www.chebfun.org | exponential convergence for analytic f |
| Rational | AAA | Nakatsukasa, Sète & Trefethen, SIAM J. Sci. Comput. 40(3) A1494, 2018, doi:10.1137/16M1106122, arXiv:1612.00337 (https://arxiv.org/abs/1612.00337) | barycentric form, greedy support points, no user parameters |
| Rational | Padé via SVD | Gonnet, Güttel & Trefethen, SIAM Review 55(1), 2013, doi:10.1137/110853236 | robust to spurious poles |
| Polynomial roots | companion-matrix eigenvalues | Edelman & Murakami, Math. Comp. 64, 1995, doi:10.1090/S0025-5718-1995-1262279-2 | backward stability analysis |
| IIR design | bilinear transform of analog prototypes, SOS cascade | Oppenheim & Schafer, *Discrete-Time Signal Processing*, 3rd ed., Pearson 2010, ISBN 978-0-13-198842-2 | SOS avoids high-order coefficient loss |
| FIR | Parks–McClellan (Remez exchange) | Parks & McClellan, IEEE Trans. Circuit Theory 19(2), 1972, doi:10.1109/TCT.1972.1083419 | equiripple optimum |
| Zero-phase | forward-backward with edge padding (filtfilt) | Gustafsson, IEEE Trans. Signal Process. 44(4), 1996, doi:10.1109/78.492552 | initial-condition handling |
| PSD | Welch averaged modified periodogram | Welch, IEEE Trans. Audio Electroacoust. 15(2), 1967, doi:10.1109/TAU.1967.1161901 | variance reduction |
| PSD | multitaper with DPSS | Thomson, Proc. IEEE 70(9), 1982, doi:10.1109/PROC.1982.12433; Slepian, Bell Syst. Tech. J. 57(5), 1978, doi:10.1002/j.1538-7305.1978.tb02104.x | low leakage bias |
| Wavelets | Mallat pyramid; Daubechies compact support | Mallat, IEEE PAMI 11(7), 1989, doi:10.1109/34.192463; Daubechies, CPAM 41(7), 1988, doi:10.1002/cpa.3160410705 | O(n) DWT; filters from tables |
| CWT | FFT-domain scaling with analytic wavelets | Torrence & Compo, BAMS 79(1), 1998, doi:10.1175/1520-0477(1998)079<0061:APGTWA>2.0.CO;2 | practical guide, significance tests |
| Gaussian filter | truncated separable kernel (4σ, SciPy default) or recursive IIR | Young & van Vliet, Signal Process. 44(2), 1995, doi:10.1016/0165-1684(95)00020-E | IIR: O(1) per pixel in σ |
| Morphology | structuring-element theory; grey-scale; fast van Herk/Gil–Werman | Haralick, Sternberg & Zhuang, IEEE PAMI 9(4), 1987, doi:10.1109/TPAMI.1987.4767941; Soille, *Morphological Image Analysis*, Springer 2004, ISBN 978-3-540-42988-3 | rectangle in O(1) per pixel |
| Edges | Canny | Canny, IEEE PAMI 8(6), 1986, doi:10.1109/TPAMI.1986.4767851 | hysteresis, NMS |
| Hough | (ρ,θ) accumulator | Duda & Hart, CACM 15(1), 1972, doi:10.1145/361237.361242 | standard |
| Segmentation | watershed by immersion | Vincent & Soille, IEEE PAMI 13(6), 1991, doi:10.1109/34.87344 | priority-flood, linear in pixels |
| Segmentation | Otsu | Otsu, IEEE Trans. SMC 9(1), 1979, doi:10.1109/TSMC.1979.4310076 | already present (`images.Otsu`) |
| Segmentation | snakes | Kass, Witkin & Terzopoulos, IJCV 1(4), 1988, doi:10.1007/BF00133570 | active contour baseline |
| Features | Shi–Tomasi; SIFT | Shi & Tomasi, CVPR 1994, doi:10.1109/CVPR.1994.323794; Lowe, IJCV 60(2), 2004, doi:10.1023/B:VISI.0000029664.99615.94 | SIFT patent status: check before shipping |
| Registration | phase correlation; sub-pixel upsampled DFT | Reddy & Chatterjee, IEEE TIP 5(8), 1996, doi:10.1109/83.506761; Guizar-Sicairos, Thurman & Fienup, Opt. Lett. 33(2), 2008, doi:10.1364/OL.33.000156 | translation to 1/100 px in O(N log N) |
| Textbooks | Gonzalez & Woods, *Digital Image Processing*, 4th ed., Pearson 2018, ISBN 978-0-13-335672-4; Mallat, *A Wavelet Tour of Signal Processing*, 3rd ed., Academic Press 2009, ISBN 978-0-12-374370-1 | | general |

## 3. Reference implementations

| implementation | language | licence | portable to Go? | notes |
|---|---|---|---|---|
| SciPy `signal`, `interpolate`, `ndimage` | Python/C/Fortran | BSD-3 (https://github.com/scipy/scipy/blob/main/LICENSE.txt) | port allowed | `interpolate` wraps FITPACK/Qhull; check bundled-code licences per file |
| FITPACK (Dierckx) | Fortran | netlib `dierckx` carries no explicit licence | oracle only unless a licence is found | http://www.netlib.org/dierckx/ |
| Qhull (Delaunay in SciPy) | C | Qhull licence (BSD-like, custom) | oracle; reimplement Bowyer–Watson | http://www.qhull.org |
| Chebfun | MATLAB | BSD-3 (LICENSE.txt, Oxford 2017) | port allowed | https://github.com/chebfun/chebfun ; AAA in `aaa.m` |
| baryrat (c-f-h/baryrat) / scipy.interpolate.AAA | Python | BSD-2 (baryrat); BSD-3 (SciPy) | port allowed | SciPy has `scipy/interpolate/_bary_rational.py` on main |
| PyWavelets | Python/C | MIT | port allowed | https://github.com/PyWavelets/pywt ; wavelet filter tables |
| scikit-image | Python/Cython | BSD-3 | port allowed (main files) | https://scikit-image.org/docs/stable/license.html — some files carry other licences |
| OpenCV | C++ | Apache-2.0 (>=4.5.0); BSD-3 earlier | port allowed (watch patented modules in contrib) | https://opencv.org/license/ |
| Octave `signal` package | Octave/C++ | GPL-3 | oracle only | https://gnu-octave.github.io/packages/signal/ |
| Octave core (interp1, polyfit, ...) | C++/Octave | GPL-3 | oracle only | https://octave.org |
| MATLAB toolboxes / Mathematica | proprietary | proprietary | oracle only if available; docs for names only | |
| Go: `gonum.org/v1/gonum` (`dsp/fourier`, `interp`, `stat`) | Go | BSD-3 | direct dependency or port | https://github.com/gonum/gonum — inventory with `go doc` before duplicating |
| Go: `golang.org/x/image` | Go | BSD-3 | dependency | resampling kernels (draw) |
| Go: `github.com/disintegration/imaging`, `github.com/anthonynsimon/bild` | Go | MIT | port allowed | reference for API shape |
| Go: `github.com/mjibson/go-dsp` | Go | ISC | port allowed | FFT/windows/spectral |

## 4. Validation

- Oracles (no root): SciPy via `uv run --with scipy==<pin> --with numpy python harness.py` (macOS
  arm64 and Linux wheels exist); scikit-image and PyWavelets via `uv run --with scikit-image
  --with pywavelets`; Octave 9 + `signal` package via pkgx/Homebrew-less release tarball or
  Docker on Linux (verify availability per platform before relying on it). Drive each from Go with
  `exec.Command` + JSON on stdin/stdout (inputs and outputs as base64 float64 or hex bits, to avoid
  decimal round-trip).
- Corpora: USC-SIPI Image Database (https://sipi.usc.edu/database/ — images free for research,
  check terms before committing any file; keep out of the repo, download in test with checksum);
  scikit-image `data` module images (BSD-3-compatible set, https://github.com/scikit-image/scikit-image/tree/main/skimage/data);
  standard test signals: chirp, Doppler, bumps, blocks (Donoho & Johnstone 1994, doi:10.1093/biomet/81.3.425);
  Runge function and Chebfun examples for approximation; MATLAB `zp2sos`/`butter` documented
  examples as fixtures; Rosetta/FIR designs from Parks–McClellan paper tables.
- "Correct" means:
  - interpolation/splines/polyfit: max relative residual vs oracle within 4 ULP on exact-node
    reproduction; knot/coefficients equal to 1e-12 relative; derivative continuity tested;
  - IIR/FIR design: max |Δ| of frequency response vs oracle, in dB and ULP; pole radius < 1;
  - filtering/convolution: relative L2 error vs oracle (<= 1e-13 for float64 direct form);
  - spectra: relative error per bin and total power (Parseval);
  - wavelets: perfect reconstruction error (<= 1e-12) and coefficients vs PyWavelets;
  - images: exact equality for integer morphology/threshold/labelling; max abs diff <= 1 LSB for
    filtered 8-bit; SSIM only as supplementary;
  - Hough/Radon: peak location equal, accumulator within 1e-9 relative.
  Report max and mean error per case, never a pass count alone.

## 5. Proposed home and plan

Org names checked with `gh api users/<name>` on 2026-10-17: `go-signal` EXISTS but is a
third-party organisation (created 2016, 1 public repo `signal`) — not usable; `go-interp` EXISTS (0
public repos, not ours); `go-dsp`, `go-interpolate`, `go-splines`, `go-imageproc`,
`go-wavelets`, `go-approx` returned 404 (free). Proposed NEW orgs: `go-approx` (interpolation +
approximation), `go-dsp` (filters, spectral, wavelets, audio I/O). Images: extend
`go-images/images` rather than a new org, except heavier items in `go-imageproc`.

| work item | org/repo | depends on | size S/M/L/XL | priority 1–3 |
|---|---|---|---|---|
| 1-D/N-D interpolation: linear, PCHIP, Akima, cubic, B-spline (de Boor), smoothing spline, `interpn` | NEW `go-approx/interp` | go-ndarray | L | 1 |
| Polynomials: polyfit (QR), polyval, roots (companion), conv/deconv, Chebyshev basis | NEW `go-approx/poly` | go-ndarray, linalg | M | 1 |
| IIR/FIR design + lfilter/filtfilt/sosfilt + ZPK/SOS types, windows (Kaiser, Tukey, DPSS), Parks–McClellan | NEW `go-dsp/filter` | go-fft | L | 1 |
| Spectral: Welch, multitaper, STFT/ISTFT, resample_poly, decimate, Hilbert, find_peaks | NEW `go-dsp/spectral` | go-fft | L | 1 |
| Image: float/nD image type, separable Gaussian, structuring elements, binary/grey morphology, labelling, distance transform, watershed, Hough/Radon, Lab/XYZ colour | `go-images/images` (+ `go-imageproc` if split) | go-ndarray | XL | 2 |
| Scattered data: robust Delaunay (Bowyer–Watson), natural neighbour, RBF | NEW `go-approx/scattered` | go-approx/interp, robust predicates | L | 2 |
| Chebfun-style adaptive approximation + AAA + Padé | NEW `go-approx/cheb` | go-fft, go-approx/poly | L | 2 |
| Wavelets: DWT/SWT/CWT/WPT, filter tables (PyWavelets, MIT) | NEW `go-dsp/wavelets` | go-fft | L | 2 |
| Image registration (phase correlation, demons), features (Harris/FAST/DoG), active contours | `go-images` / `go-imageproc` | image float type | L | 3 |
| Audio I/O (WAV/AIFF) | NEW `go-dsp/audio` | none | S | 3 |

## 6. Risks and open questions

- MATLAB names: MathWorks rate-limited (HTTP 429) the page checks, so MATLAB names are not page-verified (Curve Fitting `csaps`/`spaps`/`spapi`/`fit` and Wolfram names were).
- `go-signal` and `go-interp` are taken by third parties; the plan uses `go-dsp` and `go-approx`.
- `go-images/images` is built on 8-bit `*image.RGBA` / `go-gfx/gfx/raster.Image` (RGBA bytes; `*Raster` facades exist); scientific image processing needs float/nD arrays. Adding a float image type is the pivotal design decision and affects the Ruby binding (`images.go`).
- Chebyshev approximation and wavelets build on `go-fft/fft`'s DCT/DST, which exist.
- Delaunay robustness (degenerate/cocircular points) needs exact predicates; SciPy relies on Qhull (custom licence) — treat as oracle only.
- Bit-equality with SciPy is impossible for FITPACK-based splines (knot selection heuristics); define equality by tolerance and by fixed explicit knots.
- Licence hygiene: scikit-image and SciPy contain files with other licences; port per file. Octave signal is GPL: oracle only, no code reading-and-porting.
- SIFT-family descriptors: patent expiry and OpenCV contrib status to verify before any port.
- USC-SIPI image terms limit redistribution; do not commit images.

Verified on 2026-10-17: all 36 DOIs resolved via Crossref (title/year/first author match; Gonnet–Güttel–Trefethen corrected 110853996 -> 110853236; Sibson has no DOI); Wolfram names checked against reference.wolfram.com (nonexistent ones removed); MATLAB names only partly (MathWorks 429); licences via GitHub licence API/LICENSE files (SciPy BSD-3, PyWavelets MIT, OpenCV Apache-2.0, scikit-image BSD-3 per file, Chebfun BSD-3, baryrat BSD-2, go-dsp ISC, gonum BSD-3, imaging/bild MIT); go-images, go-fft, go-numerics/{special,optimize}, go-ndarray, go-gfx and org names checked on GitHub main.
