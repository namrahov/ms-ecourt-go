# ms-ecourt-go

A microservice for managing electronic court (e-court) applications and documents. It handles application lifecycle management, status tracking, document generation (PDF/Excel), and reporting, with permission-based access control via an external admin service.

## Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Building and Running](#building-and-running)
- [API Reference](#api-reference)
- [Database](#database)
- [External Services](#external-services)

---

## Architecture

```
ms-ecourt-go/
├── main.go                    # Entry point, server setup
├── config/                    # Environment variable configuration
├── handler/                   # HTTP handlers (controllers)
├── service/                   # Business logic
│   └── permission/            # Permission validation
├── repo/                      # Data access layer (PostgreSQL)
├── model/                     # Entities and DTOs
├── client/                    # External HTTP clients
├── mapper/                    # Entity ↔ DTO mapping
├── middleware/                # HTTP middleware (logging, context)
├── util/                      # Shared utilities
├── profiles/                  # Environment profiles (.env files)
└── infra/                     # Helm/K8s values, Docker entrypoint
```

**Stack:**
- **Language:** Go
- **Router:** `gorilla/mux`
- **Database:** PostgreSQL via `go-pg`
- **Logging:** `logrus`
- **Metrics:** Prometheus
- **Document generation:** `go-wkhtmltopdf` (PDF), `excelize` (Excel)
- **Secrets:** HashiCorp Vault (injected at container startup)

---

## Prerequisites

- Go 1.18+
- PostgreSQL
- `wkhtmltopdf` installed and available on `$PATH` (required for PDF generation)
- Access to the admin microservice

---

## Configuration

Configuration is loaded from environment variables. Profiles are stored in the `profiles/` directory.

`profiles/default.env` is always loaded first. An optional override profile can be specified with the `-p` flag (e.g., `-p prod` loads `profiles/prod.env`).

| Variable                  | Description                          | Required |
|---------------------------|--------------------------------------|----------|
| `PORT`                    | HTTP server port                     | Yes      |
| `HOSTNAME`                | Server hostname                      | Yes      |
| `DB_CARD_DELIVERY_HOST`   | PostgreSQL host                      | Yes      |
| `DB_CARD_DELIVERY_PORT`   | PostgreSQL port                      | Yes      |
| `DB_CARD_DELIVERY_NAME`   | PostgreSQL database name             | Yes      |
| `DB_CARD_DELIVERY_USER`   | PostgreSQL username                  | Yes      |
| `DB_CARD_DELIVERY_PASS`   | PostgreSQL password                  | Yes      |
| `MS_ADMIN`                | Admin microservice base URL          | Yes      |
| `LOG_LEVEL`               | Log level (`debug`, `info`, `warn`)  | No       |

---

## Building and Running

**Build:**
```bash
go build -o ms-ecourt .
```

**Run with default profile:**
```bash
./ms-ecourt
```

**Run with a named profile:**
```bash
./ms-ecourt -p dev
./ms-ecourt -p prod
```

**Run with Go directly:**
```bash
go run main.go -p dev
```

**Docker / Kubernetes:**

The entrypoint script (`infra/entrypoint-parent.sh`) sources secrets injected by Vault at `/vault/secrets/secrets` before starting the binary.

---

## API Reference

All application and document endpoints are under `/v1/ecourt` and require the `User-Id` header. Permission is validated against the `lscView` key via the admin service.

### Applications

| Method | Path                                        | Description                                  |
|--------|---------------------------------------------|----------------------------------------------|
| GET    | `/v1/ecourt/applications`                   | List applications (paginated, filterable)    |
| GET    | `/v1/ecourt/applications/{id}`              | Get a single application with comments       |
| GET    | `/v1/ecourt/applications/get/filter-info`   | Get distinct court names and judges          |
| GET    | `/v1/ecourt/applications/{id}/change-status`| Change the status of an application          |

**Query parameters for `GET /applications`:**

| Parameter      | Type   | Description                    |
|----------------|--------|--------------------------------|
| `page`         | int    | Page number (default: 1)       |
| `count`        | int    | Items per page (default: 10)   |
| `courtName`    | string | Filter by court name           |
| `judgeName`    | string | Filter by judge name           |
| `person`       | string | Filter by person name          |
| `createDateFrom` | string | Filter from date             |
| `createDateTo`  | string | Filter to date                |

**Status change body:**
```json
{
  "status": "IN_PROGRESS",
  "description": "Processing started"
}
```

**Application statuses:** `RECEIVED`, `IN_PROGRESS`, `SENT`, `HOLD`

**Application priorities:** `STANDARD`, `HIGH`

---

### Documents

| Method | Path                                    | Description                              |
|--------|-----------------------------------------|------------------------------------------|
| POST   | `/v1/ecourt/documents/generate-act`     | Generate a PDF act document              |
| GET    | `/v1/ecourt/documents/get-report`       | Generate an Excel report of applications |
| POST   | `/v1/ecourt/documents/upload`           | Upload an Excel file (max 25 KB)         |
| POST   | `/v1/ecourt/documents/upload-file`      | Upload a file with SHA1 hashing          |

- `upload` expects a `multipart/form-data` request with a `file` field.
- `upload-file` expects a `multipart/form-data` request with an `nf` field. Files are saved to `./public/`.
- `generate-act` returns a binary PDF.
- `get-report` returns an `.xlsx` file.

---

### Health & Metrics

| Method | Path        | Description           |
|--------|-------------|-----------------------|
| GET    | `/readiness` | Readiness probe      |
| GET    | `/health`    | Health check         |
| GET    | `/metrics`   | Prometheus metrics   |

---

### Request Headers

| Header           | Description                                      |
|------------------|--------------------------------------------------|
| `User-Id`        | Required. User identifier for permission checks. |
| `DP-User-ID`     | Optional. Logged for tracing.                    |
| `DP-Customer-ID` | Optional. Logged for tracing.                    |
| `x-request-id`   | Optional. Distributed tracing request ID.        |

---

## Database

PostgreSQL with the following tables:

| Table            | Description                                      |
|------------------|--------------------------------------------------|
| `application`    | Core court application entity                    |
| `document`       | Documents linked to applications                 |
| `comment`        | Internal and external comments on applications   |
| `status_history` | Audit trail of application status changes        |

**Connection pool settings:**
- Pool size: 5
- Dial timeout: 1 minute
- Max connection age: 15 minutes
- Max retries: 2

---

## External Services

### Admin Service (`MS_ADMIN`)

Used for user lookup and permission validation on every protected endpoint.

| Method | Endpoint                                   | Purpose                        |
|--------|--------------------------------------------|--------------------------------|
| GET    | `/users/{id}`                              | Fetch user details             |
| POST   | `/checkUserPermissionByUserIdAndKey`       | Validate user permission key   |

The required permission key for all e-court endpoints is `lscView`.
