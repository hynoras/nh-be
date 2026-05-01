# Noheir Backend

Backend service for **Noheir** — a chemistry-focused platform for managing experiments, procedures, and lab workflows.

The system is designed as a **modular, production-ready API** with built-in observability, scalable infrastructure, and secure session-based authentication.

🌐 Live: [https://noheir-api.duckdns.org](https://noheir-api.duckdns.org)

---

## Overview

Noheir backend is built with a focus on:

* **Modular architecture** — feature-based organization inside a structured monolith
* **Observability-first** — tracing, metrics, and health checks are first-class
* **Scalable infrastructure** — Redis, RabbitMQ, and containerized deployment
* **Production mindset** — CI/CD, health-based rollout, and rollback mechanisms

---

## Features

### Core Modules

* **Authentication**

  * Session-based authentication using cookies (`auth_session`)
  * Redis-backed session store
  * Email verification via Resend
  * Password change support

* **Users**

  * User registration and login
  * Account management

* **Experiments**

  * Experiment lifecycle management
  * Experiment result tracking

* **Procedures**

  * Define and manage experimental procedures

* **Permissions**

  * Granular access control system

---

### Infrastructure Capabilities

* Redis for session storage and rate limiting
* RabbitMQ for asynchronous messaging
* **Full Grafana Stack:**
    * **Prometheus:** Metrics collection and alerting
    * **Grafana Tempo:** Distributed tracing
    * **Grafana Loki:** Log aggregation
    * **OpenTelemetry Collector:** Observability pipeline
    * **Grafana Alloy:** Agent for observability data
    * **Grafana Dashboards:** Unified visualization
* Optional Nginx reverse proxy for routing, SSL termination, and service exposure
* Optional Secure secret management via HashiCorp Vault (Raft persistent storage)
* Swagger API documentation

---

## Tech Stack

| Category         | Technology                           |
| ---------------- | ------------------------------------ |
| Language         | Go 1.25.3                            |
| Framework        | Gin                                  |
| ORM              | GORM                                 |
| Database         | PostgreSQL                           |
| Cache            | Redis                                |
| Messaging        | RabbitMQ                             |
| Observability    | Prometheus + Grafana + Loki + Tempo  |
| API Docs         | Swagger                              |
| Email Service    | Resend                               |
| Secrets          | HashiCorp Vault                      |
| Reverse Proxy    | Nginx                                |
| Containerization | Docker + Docker Compose              |

---

## Architecture

The system follows a modular monolith design:

Client → Nginx → API (Gin)
                ├── PostgreSQL (primary data store)
                ├── Redis (sessions, caching)
                ├── RabbitMQ (async messaging)
                └── Observability pipeline
                     ├── OpenTelemetry Collector
                     ├── Prometheus (metrics)
                     ├── Tempo (tracing)
                     └── Loki (logs)

---

## Diagram

```mermaid
graph LR
    Client[Client Browser] --> Nginx[Nginx Proxy]
    
    subgraph "Noheir Backend"
        Nginx --> API[API (Gin)]
        
        subgraph "Core Services"
            API -->|Primary Storage| DB[(PostgreSQL)]
            API -->|Sessions/Cache| Redis[(Redis)]
            API -->|Async Tasks| MQ[(RabbitMQ)]
        end
        
        subgraph "Observability Stack"
            API -.->|Traces| OTel[OTel Collector]
            API -.->|Metrics| OTel
            API -.->|Logs| OTel
            
            OTel -->|Metrics| Prometheus[(Prometheus)]
            OTel -->|Traces| Tempo[(Grafana Tempo)]
            OTel -->|Logs| Loki[(Grafana Loki)]
            
            Prometheus --> Grafana[Grafana]
            Tempo --> Grafana
            Loki --> Grafana
        end
    end
    
    subgraph "External"
        Nginx -->|Forwarding| Vault[(HashiCorp Vault)]
        Nginx -->|HTTPS/SSL| Internet[(Internet)]
    end
    
    style DB fill:#9f9,stroke:#333,stroke-width:2px
    style Redis fill:#f9f,stroke:#333,stroke-width:2px
    style MQ fill:#f99,stroke:#333,stroke-width:2px
    style Vault fill:#ff9,stroke:#333,stroke-width:2px
```

---

## Project Structure

```text
/build         # Build artifacts or scripts
/cmd           # Application entrypoint (main.go)
/config        # Configuration loading + Vault integration
/docs          # Swagger documentation
/infra         # Infrastructure setup (database, redis, observability, secrets)
/internal
  /features    # Feature modules (auth, experiment, permission, procedure, user)
  /middleware  # Gin middlewares (auth, tracing, metrics, request ID)
/mq            # RabbitMQ configuration
/pkg           # Shared utilities
/router        # Route definitions
/templates     # Email templates
/tests         # Test suite
```

---

## Getting Started

### Prerequisites

* Docker & Docker Compose
* Go 1.25+ (Required if not running the API via Docker)
* (Optional) Air for hot reload

> **Note on Managed Infrastructure:** The application is designed to be flexible. You can use your own managed infrastructure (such as **Supabase** for PostgreSQL or **Redis Cloud**) by updating the `.env` variables, or spin up local instances using the provided `docker-compose.yml`. By default, the local PostgreSQL and Redis services in the compose file are commented out, assuming the use of a managed database like Supabase.

> **Note on HashiCorp Vault:** To mirror production behavior, Vault is configured to use the **Raft consensus algorithm** for persistent storage (mapped to `infra/secret/data`). This setup requires a manual initialization and unseal process, ensuring a secure, production-grade workflow even in development.

---

## Run in Development Mode

You can run the backend in development mode by starting only the required dependencies and running the API locally.

### Minimal Dev (Fast Startup)

Start only core dependencies:

```bash
docker compose up -d rabbitmq
```

Then run the API:

```bash
go run cmd/main.go
```

Or with hot reload:

```bash
air
```

> Use this mode for day-to-day development. It has the fastest startup time and lowest resource usage.

---

### Full Dev (With Observability & Infra)

Start all supporting services, including observability and Vault:

```bash
docker compose up -d --build rabbitmq vault prometheus tempo loki otel-collector alloy grafana
```

Then run the API:

```bash
go run cmd/main.go
```

> This mode mirrors a production-like environment, including tracing, metrics, logging, and secret management.
> It is more resource-intensive and recommended when debugging, testing observability, or validating infrastructure behavior.

---

### Notes

* PostgreSQL and Redis may be external (e.g., Supabase, Redis Cloud) depending on your `.env` configuration
* Vault requires manual initialization and unsealing before use
* Ensure `.env` is correctly configured before starting

---

## Environment Variables

Create a `.env` file in the root directory:

```env
APP_ENV=dev
PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=noheir

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD= # Leave blank for local development

RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest

USE_VAULT=false

FRONTEND_URL=http://localhost:3000
RESEND_API_KEY=your_key_here
```

---

## API

### Base Path

```
/api/v1
```

### Example Endpoints

**Authentication**

```http
POST /api/v1/auth/signup
POST /api/v1/auth/login
GET  /api/v1/auth/verify/{token}
```

**Health Checks**

```http
GET /health/live
GET /health/ready
```

**Observability**

```http
GET /metrics        # Prometheus metrics
GET /swagger/*any   # Swagger documentation
```

---

## Observability

Noheir features a comprehensive observability stack:

* **Tracing:** OpenTelemetry + Grafana Tempo
* **Metrics:** Prometheus (`/metrics`)
* **Logging:** Grafana Loki
* **Dashboards:** Pre-configured Grafana instance
* **Health Checks:**
  * `/health/live` — container is running
  * `/health/ready` — dependencies are ready

This enables:
* End-to-end request visibility
* Real-time performance monitoring
* Log-trace correlation
* Health-based automated deployments

---

## Deployment

Deployment is handled via **GitLab CI/CD**.

### Pipeline Stages

* `lint`
* `test`
* `build`
* `deploy`

### Strategy

* Docker image is built per commit
* Tagged with commit SHA (not `latest` in production)
* Pushed to GitLab Container Registry
* Deployed via SSH to VPS
* Updated using:

```bash
docker compose up -d
```

### Reliability

* Health-check-based rollout
* Automatic rollback on failure

---

## Design Philosophy

Noheir backend is built with the following principles:

* **Feature modularity**
  Avoid tightly coupled layers by organizing logic per domain

* **Operational visibility**
  Metrics, tracing, and health checks are non-optional

* **Production-first mindset**
  Deployment, reliability, and observability are considered early

* **Scalable evolution**
  Designed to evolve from monolith → modular monolith → microservices

---

## Roadmap

* Compound management module
* File upload system (experiment assets)
* Experiment branching and versioning
* Microservices decomposition (Experiment domain)
* Enhanced observability dashboards (e.g. Grafana)

---

## Notes

* Redis is required for session management
* RabbitMQ enables async workflows and future scaling
* Vault integration is optional but recommended for production

---

## License

Proprietary - All Rights Reserved