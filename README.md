# Multi-Service Fullstack Demo

This project is a multi-service fullstack demo application with a **Go backend**, **React frontend**, **PostgreSQL database**, and observability stack (**Prometheus**, **Grafana**, **OpenTelemetry Collector**), all orchestrated via Docker Compose.

---

## 🏗 Tech Stack

- **Backend:** Go (Golang) + HTTP/GraphQL/REST API
- **Frontend:** React + Vite
- **Database:** PostgreSQL 15
- **Observability:** 
  - OpenTelemetry Collector
  - Prometheus
  - Grafana
- **Containerization:** Docker & Docker Compose 3.9

---

## 📁 Project Structure

.
├── backend/ # Go backend with Dockerfile
├── frontend/ # React frontend with Dockerfile
├── monitoring/ # Prometheus, Grafana, OpenTelemetry configs
├── sql/ # SQL scripts for DB initialization
├── docker-compose.yml # Docker Compose file
└── README.md

yaml
Copy code

---

## 🚀 Services

| Service     | Port   | Description |
|------------|--------|-------------|
| **backend** | 8080  | Go API service |
| **frontend** | 5173 | React frontend |
| **postgres** | 5432 | PostgreSQL database |
| **collector** | 4318 | OpenTelemetry Collector |
| **prometheus** | 9090 | Prometheus metrics server |
| **grafana** | 3001 | Grafana dashboards (host:container port mapping 3001:3000) |

---

## ⚡ Getting Started

### Prerequisites

- Docker
- Docker Compose
- Go (for local development if not using Docker)

### Start All Services

```bash
docker-compose up --build
This will build the backend and frontend images and start all services with the correct dependencies.

Stop Services
bash
Copy code
docker-compose down
🛠 Backend Healthcheck
The Go backend exposes a health endpoint:

bash
Copy code
GET http://localhost:8080/health
Docker healthcheck ensures the backend is ready before dependent services (frontend) start.

🗄 Database
PostgreSQL is initialized with scripts from ./sql

Default credentials:

ini
Copy code
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin123
POSTGRES_DB=multi_demo
Healthcheck ensures the database is ready before the backend starts.

📊 Observability
OpenTelemetry Collector runs on port 4318 and collects traces/metrics from the backend.

Prometheus runs on port 9090 and scrapes metrics from the collector.

Grafana runs on port 3001 and loads dashboards from ./monitoring/grafana/dashboards.

⚙️ Environment Variables
Backend:

ini
Copy code
DB_HOST=postgres
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=multi_demo
DB_PORT=5432
OTEL_EXPORTER_OTLP_ENDPOINT=http://collector:4318
Frontend:

ini
Copy code
REACT_APP_API_URL=http://backend:8080
🧰 Docker Volumes
db_data → PostgreSQL persistent data

grafana_data → Grafana dashboards and settings

