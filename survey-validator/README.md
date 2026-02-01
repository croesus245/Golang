# Survey Data Validator

A simple tool for surveyors to check their coordinate data for common errors.

## What This Tool Does

Before you submit your survey or use it in the field, run your points through this validator. It catches mistakes that are easy to miss:

- **Duplicate points** - Same coordinates entered twice (copy/paste errors)
- **Zero coordinates** - Points with 0,0 (forgot to enter data)
- **Outliers** - A point way off from the others (typo in coordinates)
- **Traverse closure** - Does your loop close? What's the precision?
- **Distance issues** - Unusual jumps between traverse points

## How To Use

1. Start the server (see below)
2. Open http://localhost:8080 in your browser
3. Enter your points in the table - Point ID, Easting, Northing, Height, Type
4. Click **Validate Survey Data**
5. Review results - problems show in red/orange, good data shows green

## Quick Start

```bash
cd survey-validator
go build -o survey-validator.exe ./cmd/server
.\survey-validator.exe
```

Open http://localhost:8080 and start entering points.

## Features

- Web interface - no JSON knowledge needed, just fill in the table
- Instant validation - results in under a second
- Confidence score - overall quality rating for your data
- Detailed reports - shows exactly which points have problems

## Point Types

- **traverse** - Points in your traverse loop
- **control** - Known control points
- **detail** - Detail/topo points

## Project Structure

```
survey-validator/
├── api/                    # HTTP handlers and server
│   ├── server.go
│   └── request.go
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── domain/                 # Survey logic and validation rules
│   ├── spatial.go          # Geometric calculations
│   ├── traverse.go         # Traverse-specific checks
│   └── validators.go       # Core validation functions
├── engine/
│   └── engine.go           # Concurrent validation orchestrator
├── models/                 # Data structures
│   ├── point.go            # Survey point model
│   └── report.go           # Validation report model
├── testdata/               # Sample data for testing
│   ├── sample_survey.json
│   └── sample_with_errors.json
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or later

### Installation

```bash
# Clone or navigate to the project
cd survey-validator

# Download dependencies (none required - uses standard library only)
go mod tidy

# Build the application
go build -o survey-validator.exe ./cmd/server

# Or run directly
go run ./cmd/server
```

### Running the Server

```bash
# Default port 8080
go run ./cmd/server

# Custom port
go run ./cmd/server -port 3000

# Using environment variable
PORT=3000 go run ./cmd/server
```

## API Endpoints

### Health Check

```
GET /health
```

Response:
```json
{
  "status": "healthy",
  "service": "survey-validator"
}
```

### Validate Survey Data

```
POST /api/v1/validate
Content-Type: application/json
```

#### Request Body

```json
{
  "project_id": "SURVEY-001",
  "coordinate_system": "UTM Zone 36N",
  "points": [
    {
      "point_id": "T1",
      "easting": 500050.123,
      "northing": 6000025.456,
      "height": 101.200,
      "survey_type": "traverse"
    },
    {
      "point_id": "T2",
      "easting": 500100.789,
      "northing": 6000050.321,
      "height": 102.100,
      "survey_type": "traverse"
    }
  ]
}
```

#### Survey Types

- `traverse` - Traverse survey points
- `control` - Control points
- `detail` - Detail/topographic points

#### Response

```json
{
  "project_id": "SURVEY-001",
  "timestamp": "2026-02-01T10:30:00Z",
  "status": "PASS",
  "confidence_score": 100,
  "summary": {
    "total_points": 2,
    "traverse_points": 2,
    "control_points": 0,
    "detail_points": 0,
    "points_with_height": 2,
    "bounding_box": {
      "min_easting": 500050.123,
      "max_easting": 500100.789,
      "min_northing": 6000025.456,
      "max_northing": 6000050.321
    },
    "centroid_easting": 500075.456,
    "centroid_northing": 6000037.889
  },
  "issues": [],
  "checks_performed": [
    "input_validation",
    "duplicate_detection",
    "distance_bearing_check",
    "outlier_detection",
    "traverse_closure"
  ],
  "processing_time": "1.234ms"
}
```

## Validation Checks

### 1. Input Validation
- Checks for empty point IDs
- Flags zero coordinates
- Validates survey types

### 2. Duplicate Detection
- **Duplicates**: Points within 0.001m (error)
- **Near-duplicates**: Points within 0.01m (warning)

### 3. Distance & Bearing Consistency
- Flags very short distances (< 0.1m)
- Detects large bearing changes (> 170°)
- Identifies unusual distance ratios

### 4. Outlier Detection
- Calculates centroid of all points
- Flags points beyond 3 standard deviations

### 5. Traverse Closure
- Calculates linear misclosure
- Computes relative precision
- Quality ratings:
  - Good: better than 1:10000
  - Acceptable: 1:5000 to 1:10000
  - Poor: 1:1000 to 1:5000
  - Unacceptable: worse than 1:1000

## Testing with Sample Data

```bash
# Start the server
go run ./cmd/server

# In another terminal, test with sample data
curl -X POST http://localhost:8080/api/v1/validate ^
  -H "Content-Type: application/json" ^
  -d @testdata/sample_survey.json

# Test with error-containing data
curl -X POST http://localhost:8080/api/v1/validate ^
  -H "Content-Type: application/json" ^
  -d @testdata/sample_with_errors.json
```

## Validation Status

| Status | Description |
|--------|-------------|
| `PASS` | No errors or warnings detected |
| `WARNING` | Warnings detected but no critical errors |
| `FAIL` | Critical errors detected |

## Confidence Score

The confidence score (0-100) is calculated based on issues found:
- Error: -15 points
- Warning: -5 points
- Info: -1 point

## Architecture

The system follows clean architecture principles:

- **API Layer** (`api/`) - HTTP request/response handling
- **Engine** (`engine/`) - Orchestrates concurrent validation
- **Domain** (`domain/`) - Survey logic and spatial calculations
- **Models** (`models/`) - Data structures

Each validation check runs independently in its own goroutine, allowing:
- Fast parallel processing
- Easy addition of new checks
- Clean separation of concerns

## Limitations

This system intentionally focuses on validation only:

- ❌ No coordinate transformations
- ❌ No spatial database integration
- ❌ No graphical visualization
- ❌ No adjustment algorithms

## Future Extensions

Possible enhancements:
- Coordinate system transformations
- GIS platform integration
- Advanced statistical analysis
- Survey adjustment modules
- Land registry integration

## License

MIT License
