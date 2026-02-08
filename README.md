# Survey Data Validator

A quick sanity check for your survey coordinates.

[![Live Demo](https://img.shields.io/badge/demo-live-brightgreen)](https://surveyvalidator.vercel.app)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![Deploy](https://img.shields.io/badge/deploy-Vercel-black?logo=vercel)](https://vercel.com)

> We've all been there: a typo in the coordinates, a duplicated point ID, a traverse that doesn't close. This tool catches those before you submit.

![Survey Validator Screenshot](https://via.placeholder.com/800x400/f8fafc/1e293b?text=Survey+Data+Validator)

---

## What It Does

Paste your points, click validate, get a report. It runs 8+ checks in parallel and tells you what's wrong (and what's fine):

| Check | What It Catches |
|-------|------------------|
| **Duplicates** | Same point ID used twice, or two shots at basically the same spot |
| **Traverse Closure** | Does your loop close? What's the precision ratio? |
| **Outliers** | That one point way off from everything else (probably a typo) |
| **Bad Input** | Missing coords, zeroes, empty point IDs |
| **Geometry** | Weird leg lengths, sudden direction changes, suspicious patterns |

It also:
- Computes Bowditch adjustment so you can see corrected coordinates
- Draws an interactive plot you can zoom, pan, and measure on
- Exports to JSON, CSV, or PNG with one click
- Lets you click any issue to jump straight to that row in your data

---

## Try It

**[surveyvalidator.vercel.app](https://surveyvalidator.vercel.app)**

No install, no signup. Just paste your data and go.

---

## Getting Started

### Use the web app
Just open [surveyvalidator.vercel.app](https://surveyvalidator.vercel.app). Done.

### Or run it locally

```bash
# Clone the repository
git clone https://github.com/yourusername/survey-validator.git
cd survey-validator

# Run the server
go run ./cmd/server

# Open in browser
open http://localhost:8080
```

---

## How It Works

```
  Your Data  ──▶  Validation Engine  ──▶  Report + Visualization
     │                  │                        │
  paste/upload     8 checks in            pass/fail, issue list,
  or type it       parallel (<100ms)      interactive plot
```

### 1. Get your data in

- **Paste** — Copy from Excel, your data collector, wherever
- **Upload** — Drag a .csv file onto the page
- **Type** — Just fill in the table manually (Tab and arrows work)

### 2. Pick a tolerance

| Preset | Dup Threshold | Near-Dup | When to use |
|--------|---------------|----------|-------------|
| Survey Grade | 1mm | 1cm | Precise control, boundaries |
| Engineering | 1cm | 10cm | Construction stakeout |
| Mapping | 10cm | 1m | GIS, topo, recon |

### 3. Hit validate

Click the button. In under 100ms you get:

- **PASS / WARNING / FAIL** — overall verdict
- **Confidence score** — rough quality rating (0-100%)
- **Issue list** — click any issue to jump to that row
- **Visualization** — see your points on a plot

---

## The Plot

The interactive canvas does more than just show dots:

| Feature | What it does |
|---------|---------------|
| **Zoom** | +/− buttons or scroll wheel |
| **Pan** | Click and drag when zoomed in |
| **Measure** | Click two points, get distance and bearing |
| **Labels** | Toggle point IDs on/off |
| **Adjusted** | Overlay the Bowditch-corrected positions |
| **Export** | Download as PNG |
| **Minimap** | Shows where you are when zoomed in |
| **Tooltips** | Hover for details |

**What the symbols mean:**
- ▲ Green triangle = Control point
- ● Blue circle = Traverse station  
- ◆ Gray diamond = Detail/topo shot
- Red ring around anything = Problem

---

## Traverse Closure

If you've got traverse points, we calculate how well your loop closes:

| Metric | What it means |
|--------|---------------|
| **Misclosure ΔE/ΔN** | How far off you are in easting and northing |
| **Linear Misclosure** | Total closure error (meters) |
| **Closure Ratio** | Like "1:12,500" — that's 1mm error per 12.5m traveled |
| **Rating** | Excellent / Good / Acceptable / Poor |

Rule of thumb:
- **1:10,000+** — Great. Control-quality work.
- **1:5,000+** — Good. Typical boundary survey.
- **1:3,000+** — Okay for topo.
- **Below that** — Probably need to re-run something.

If closure is acceptable, we run Bowditch adjustment automatically and you can see the corrected coordinates on the plot.

---

## API

If you want to integrate this into your own workflow, here's the API.

### Health check

```http
GET /health
```

```json
{ "status": "healthy", "service": "survey-validator" }
```

### Validate Survey

```http
POST /api/v1/validate
Content-Type: application/json
```

**Request:**
```json
{
  "project_id": "SITE-2026-001",
  "points": [
    { "point_id": "CP1", "easting": 500000.000, "northing": 600000.000, "height": 100.0, "survey_type": "control" },
    { "point_id": "T1", "easting": 500050.125, "northing": 600030.250, "height": 100.15, "survey_type": "traverse" },
    { "point_id": "T2", "easting": 500100.380, "northing": 600055.620, "height": 100.32, "survey_type": "traverse" }
  ]
}
```

**Response:**
```json
{
  "status": "PASS",
  "confidence_score": 95,
  "issues": [],
  "summary": { "total_points": 3, "traverse_points": 2, "control_points": 1 },
  "traverse_adjustment": { "closure_ratio": "1:12500", "status": "PASS" },
  "checks_performed": ["input_validation", "duplicate_detection", "traverse_closure", ...]
}
```

---

## Code Layout

If you want to dig into the code:

```
survey-validator/
├── api/                    # HTTP handlers
│   ├── validate/index.go   # Vercel serverless function
│   ├── health/index.go     # Health check endpoint
│   └── server.go           # Local dev server
├── domain/                 # Business logic
│   ├── validators.go       # Core validation checks
│   ├── traverse.go         # Traverse closure & adjustment
│   ├── spatial.go          # Geometric calculations
│   └── leveling.go         # Height validation
├── engine/                 # Orchestration
│   └── engine.go           # Concurrent check runner
├── models/                 # Data structures
│   ├── point.go            # Survey point model
│   ├── report.go           # Validation report
│   └── traverse.go         # Traverse adjustment model
├── public/                 # Frontend
│   └── index.html          # Single-file app (~3000 lines)
├── testdata/               # Sample datasets
│   ├── sample_survey.json
│   └── synthetic_survey.csv
└── vercel.json             # Deployment config
```

---

## Under the Hood

### Concurrent validation

All checks run in parallel (Go goroutines). A 500-point file takes about 50-100ms.

### Coordinate system

This works with **projected coordinates in meters**. If you're in UTM, State Plane, or a local grid, you're good. Lat/long won't work—project first.

### Outlier detection

We use a simple 3-sigma test: find the centroid of all points, compute the standard deviation of distances from it, and flag anything beyond 3 standard deviations. Works well for clustered data; linear traverses might trigger false positives.

### Bowditch adjustment

The classic compass rule: distribute the misclosure proportionally based on leg distances. Longer legs get more of the correction. After adjustment, your traverse closes perfectly.

---

## Running It Yourself

You'll need Go 1.21 or later.

```bash
# Run locally
go run ./cmd/server

# Then open http://localhost:8080

# Run tests
go test ./...
```

To deploy your own copy on Vercel:

```bash
npm i -g vercel
vercel --prod
```

---

## What It Doesn't Do (Yet)

- **Lat/long** — You need projected coordinates. Convert first.
- **Out-of-order traverses** — Points need to be in the order you walked them.
- **Leveling runs** — Vertical-only validation is on the roadmap.
- **Raw angles** — We work with coordinates, not field observations.
- **Huge files** — Keep it under ~1000 points or the browser gets sluggish.

---

## Maybe Someday

- Leveling run validation
- Angular misclosure from raw observations
- Coordinate transformation between systems
- PDF export
- Save/load projects
- Batch processing

---

## License

MIT. Do whatever you want with it.

---

## Thanks

Survey math from Ghilani & Wolf's *Elementary Surveying* and *Adjustment Computations*. Built with Go and plain JavaScript—no frameworks, no build step, no npm install.

---

<p align="center">
  <a href="https://surveyvalidator.vercel.app">Give it a try →</a>
</p>
