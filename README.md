# File Parser

**File Parser** is a Go service for parallel parsing of `.tsv` files. It processes files in multiple goroutines, stores parsed data and parsing errors in PostgreSQL, and generates **RTF reports** containing both parsed data and errors. The service provides REST API endpoints and Swagger documentation.

---

## Features

- Multi-goroutine TSV file parsing.
- Stores parsed data and parsing errors in PostgreSQL.
- Generates **RTF output files** with parsed data and parsing errors.
- REST API endpoints with pagination (`limit`, `page`).
- Fetch error reports by filename.
- Health check endpoint.
- Swagger documentation available at `/swagger/*any`.
- Configurable log level and format.
- Configurable number of parser workers and input/output directories.

---

## Configuration

The service is configured via `config/config.yaml` and environment variables.

### Example `config.yaml`

```yaml
# Server config
server:
  read_timeout: 4s
  write_timeout: 10s
  max_header_bytes: 1048576
  time_for_graceful_shutdown: 10s

# Logger config
logger:
  level: info      # e.g., debug, info, warn, error
  format: json     # e.g., json, text

# Parser config
parser:
  count_of_workers: 3
  scan_frequency: 10s
  input_dir: ./input-data
  output_dir: ./output-data
  error_dir: ./error-data
```

### Environment Variables

```env
CONFIG_PATH=config/config.yaml

API_KEY=api_key
API_HOST=0.0.0.0
API_PORT=3505

DB_HOST=postgres
DB_PORT=5432
DB_USER=parser
DB_PASSWORD=12345
DB_NAME=parsed_files
SSL_MODE=disable
```

---

## API Endpoints

All endpoints under `/api/v1` require an API key (except `/health` and Swagger docs).

### Health Check

- **GET** `/api/v1/health`
- **Description:** Checks service availability.
- **Response:** `200 OK`

### Get Parsed Data

- **GET** `/api/v1/get-data`
- **Query Parameters:**
  - `unit_guid` (string, required) — Unit GUID
  - `limit` (int, optional, 1–1000, default 1000) — Number of records per page
  - `page` (int, optional, ≥1, default 1) — Page number
- **Response:** JSON with parsed data.

### Get Parsing Errors

- **GET** `/api/v1/get-errors`
- **Query Parameters:**
  - `filename` (string, required) — Name of the file
  - `limit` (int, optional, 1–1000, default 1000)
  - `page` (int, optional, ≥1, default 1)
- **Response:** JSON with parsing errors.

### Swagger Documentation

- Available at `/swagger/*any`
- Provides interactive API documentation for all endpoints.

### Authorization

Added at HTTP header new line
```http
Authorization: API-KEY
```

---

## Running the Service

1. Install dependencies:

```bash
go mod tidy
```

> If you launch API not inside the Docker container

2. Set environment variables:

```bash
export CONFIG_PATH=config/config.yaml

export API_KEY=your_api_key
export API_HOST=0.0.0.0
export API_PORT=3505

export DB_HOST=postgres
export DB_PORT=5432
export DB_USER=parser
export DB_PASSWORD=12345
export DB_NAME=parsed_files
export SSL_MODE=disable
```

3. Run/Stop the service:

```bash
make build

make up

make down
```

> If you launch API not inside the Docker container, you need to execute this command:

```bash
make launch
```

4. Launch tests

```bash
make test
```

5. Access:

- Health check: `http://localhost:3505/api/v1/health`
- Swagger docs: `http://localhost:3505/swagger/index.html`

---

## Output Files

- Parsed data is written to **RTF files** in the configured `output_dir`.
- Parsing errors are written to **RTF files** in the configured `error_dir`.

---

## Directory Structure

```
├── cmd/api/main.go     # Entry point
├── config/             # Configuration files
├── input-data/         # TSV files to parse
├── output-data/        # Generated RTF files with parsed data
├── error-data/         # RTF files with parsing errors
└── internal/           # Internal packages (parser, service, handler, db, etc.)
```
