
# Repono

<p align="center">
  <strong>Your Files. Your World.</strong>
  <br />
  A cloud-based file storage and sharing platform built with Go and React.
  <br /><br />
  <a href="https://repono-frontend.onrender.com">Live Application</a>
  ·
  <a href="https://repono-backend-3uq9.onrender.com/health">API Health</a>
  ·
  <a href="https://github.com/her0dx-56/Repono">Source Code</a>
</p>

---

## Overview

**Repono** is a cloud-based file storage and sharing platform that enables users to upload, manage, and securely share files through a modern web interface.

Built with a Go backend and React frontend, Repono uses PostgreSQL for persistent data storage, Redis for rate limiting and caching, and AWS S3 for file storage.

The application provides JWT-based authentication, public and private file sharing, Redis-backed API rate limiting, and caching for public share data. Its layered backend architecture separates HTTP handling, business logic, data access, and storage.

**Live Demo:** [https://repono-frontend.onrender.com](https://repono-frontend.onrender.com)

---

## Features

- **User Authentication** — Register and log in using JWT-based authentication.
- **File Upload** — Upload files to cloud-based object storage.
- **File Management** — View and manage uploaded files through a web dashboard.
- **Public File Sharing** — Generate unique shareable links for accessing files.
- **Private File Sharing** — Share files with specific users.
- **AWS S3 Integration** — Store uploaded file contents in cloud object storage.
- **Redis-Based Rate Limiting** — Control excessive API requests using Redis-backed rate limiting.
- **Redis Caching** — Cache public file-sharing data to reduce repeated database lookups.
- **PostgreSQL** — Persist user accounts, file metadata, and sharing information.
- **REST API** — Go-based backend for application operations.
- **Responsive Frontend** — Modern React interface for file management and sharing.
- **Docker Support** — Containerized application setup using Docker and Docker Compose.

---

## Tech Stack

| Component | Technology |
|---|---|
| Frontend | React, JavaScript, Vite |
| Routing | React Router |
| Backend | Go |
| API | REST |
| Database | PostgreSQL |
| Database Hosting | Neon |
| Caching | Redis |
| Rate Limiting | Redis |
| Redis Hosting | Upstash |
| Object Storage | AWS S3 |
| Authentication | JWT |
| Containerization | Docker, Docker Compose |
| Deployment | Render |
| Database Migrations | Goose |
| Version Control | Git, GitHub |

---

## Architecture

Repono follows a layered backend architecture that separates request handling, business logic, and data access.

```text
                         ┌─────────────────────┐
                         │        User         │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │    React + Vite     │
                         │      Frontend       │
                         └──────────┬──────────┘
                                    │
                               HTTPS / REST
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │      Go Server      │
                         │                     │
                         │   CORS Middleware   │
                         │   Auth Middleware   │
                         │   Rate Limiting     │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │      Handlers       │
                         │   HTTP Requests     │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │      Services       │
                         │   Business Logic    │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │    Repositories     │
                         │    Data Access      │
                         └──────┬──────┬──────┘
                                │      │
                      ┌─────────┘      └──────────┐
                      ▼                           ▼
             ┌─────────────────┐         ┌─────────────────┐
             │   PostgreSQL    │         │     Redis       │
             │                 │         │                 │
             │ User Metadata  │         │ Rate Limiting   │
             │ File Metadata  │         │ Public Share    │
             │ Sharing Data   │         │ Cache           │
             └─────────────────┘         └─────────────────┘
                               
                         ┌─────────────────┐
                         │      AWS S3     │
                         │                 │
                         │   File Storage  │
                         └─────────────────┘
```

### Design Principles

- **Separation of Concerns:** HTTP handlers, business logic, repositories, and storage have distinct responsibilities.
- **Layered Architecture:** Handlers communicate with services, while repositories manage data access.
- **Externalized File Storage:** File contents are stored separately from application metadata.
- **Redis Integration:** Redis supports API rate limiting and caching of public share data.
- **Environment-Based Configuration:** Database connections, storage credentials, and other deployment settings are configured through environment variables.

---

## How It Works

### Authentication

Users can register and log in to access their accounts. JWT-based authentication is used to authorize protected API requests.

### File Upload and Management

Authenticated users can upload files through the frontend. The backend stores file contents in AWS S3 and maintains associated metadata in PostgreSQL.

Users can then view and manage their uploaded files through the dashboard.

### Public and Private Sharing

Repono supports two file-sharing approaches:

- **Public Sharing:** Generate a unique link that recipients can use to access a shared file.
- **Private Sharing:** Share files with specific users through the application's sharing functionality.

### Redis-Based Rate Limiting

Redis is used to implement rate limiting for API requests. This helps control excessive requests and reduce the risk of endpoint abuse.

### Redis Caching

Repono uses Redis to cache public file-sharing data. Caching frequently requested share information can reduce repeated database queries and improve response efficiency.

---

## Project Structure

```text
Repono/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   └── jwt.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── database.go
│   │
│   ├── handler/
│   │   ├── auth_test_handler.go
│   │   ├── file_handler.go
│   │   ├── file_share_handler.go
│   │   ├── health.go
│   │   ├── public_file_share_handler.go
│   │   └── user_handler.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   └── rate_limit.go
│   │
│   ├── model/
│   │   ├── file.go
│   │   ├── file_response.go
│   │   ├── file_share.go
│   │   ├── file_share_request.go
│   │   ├── file_share_response.go
│   │   ├── login_response.go
│   │   ├── public_file_share.go
│   │   ├── user.go
│   │   ├── user_login_request.go
│   │   ├── user_request.go
│   │   └── user_response.go
│   │
│   ├── redis/
│   │   ├── client.go
│   │   ├── client_test.go
│   │   ├── public_share_cache.go
│   │   ├── ratelimiter.go
│   │   ├── ratelimiter_test.go
│   │   └──
│   │
│   ├── repository/
│   │   ├── file_repository.go
│   │   ├── file_share_repository.go
│   │   ├── public_file_share_repository.go
│   │   └── user_repository.go
│   │
│   ├── server/
│   │   └── server.go
│   │
│   ├── service/
│   │   ├── file_service.go
│   │   ├── file_service_test.go
│   │   ├── file_share_service.go
│   │   ├── health.go
│   │   ├── public_file_share_service.go
│   │   └── user_service.go
│   │
│   └── storage/
│       ├── local.go
│       ├── local_storage_test.go
│       ├── s3.go
│       ├── s3_client.go
│       └── storage.go
│
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_files.sql
│   ├── 003_add_indexes.sql
│   ├── 004_add_files_owner_unique.sql
│   ├── 005_create_file_shares.sql
│   └── 006_create_public_file_share.sql
│
├── frontend/
│   ├── public/
│   │   └── favicon.svg
│   │
│   ├── src/
│   │   ├── components/
│   │   │   ├── common/
│   │   │   │   ├── Brand.jsx
│   │   │   │   ├── Button.jsx
│   │   │   │   ├── ConfirmDialog.jsx
│   │   │   │   ├── EmptyState.jsx
│   │   │   │   ├── Modal.jsx
│   │   │   │   └── Toast.jsx
│   │   │   │
│   │   │   ├── files/
│   │   │   │   ├── FileActions.jsx
│   │   │   │   ├── FileCard.jsx
│   │   │   │   ├── FileTable.jsx
│   │   │   │   └── UploadZone.jsx
│   │   │   │
│   │   │   ├── layout/
│   │   │   │   ├── AppLayout.jsx
│   │   │   │   ├── Header.jsx
│   │   │   │   └── Sidebar.jsx
│   │   │   │
│   │   │   ├── mascot/
│   │   │   │   ├── CourierMascot.jsx
│   │   │   │   ├── MascotReaction.jsx
│   │   │   │   └── UploadAnimation.jsx
│   │   │   │
│   │   │   └── sharing/
│   │   │       ├── ShareDialog.jsx
│   │   │       └── SharedLinksTable.jsx
│   │   │
│   │   ├── context/
│   │   │   ├── AuthContext.jsx
│   │   │   └── FileContext.jsx
│   │   │
│   │   ├── hooks/
│   │   │   ├── useAuth.js
│   │   │   ├── useFiles.js
│   │   │   └── useToast.js
│   │   │
│   │   ├── pages/
│   │   │   ├── AuthPage.jsx
│   │   │   ├── DashboardPage.jsx
│   │   │   ├── LandingPage.jsx
│   │   │   └── SharedFilePage.jsx
│   │   │
│   │   ├── services/
│   │   │   ├── apiClient.js
│   │   │   ├── authService.js
│   │   │   ├── fileService.js
│   │   │   ├── mockService.js
│   │   │   └── shareService.js
│   │   │
│   │   ├── styles/
│   │   │   ├── animations.css
│   │   │   ├── dashboard.css
│   │   │   ├── global.css
│   │   │   ├── layout.css
│   │   │   └── tokens.css
│   │   │
│   │   ├── utils/
│   │   │   ├── fileHelpers.js
│   │   │   ├── formatFileSize.js
│   │   │   └── validators.js
│   │   │
│   │   ├── App.jsx
│   │   └── main.jsx
│   │
│   ├── .dockerignore
│   ├── .env.example
│   ├── .gitignore
│   ├── Dockerfile
│   ├── README.md
│   ├── index.html
│   ├── nginx.conf
│   ├── package-lock.json
│   ├── package.json
│   └── vite.config.js
│
├── .dockerignore
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/)
- [Node.js and npm](https://nodejs.org/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- PostgreSQL database
- Redis
- AWS account with an S3 bucket

### 1. Clone the Repository

```bash
git clone https://github.com/her0dx-56/Repono.git
cd Repono
```

### 2. Configure Environment Variables

Create a `.env` file in the project root:

```env
HOST=0.0.0.0
PORT=9000

DATABASE_URL=your_postgresql_connection_string
JWT_SECRET=your_random_secret

REDIS_ADDR=redis://localhost:6379

AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=your_aws_region
AWS_S3_BUCKET=your_s3_bucket_name

FRONTEND_ORIGIN=http://localhost:5173
```

Replace the placeholder values with your own configuration. Never commit real secrets to GitHub.

### 3. Set Up PostgreSQL

Create a PostgreSQL database and set its connection string in `DATABASE_URL`.

Install [Goose](https://github.com/pressly/goose), then run the migrations from the project root:

```bash
goose postgres "$DATABASE_URL" -dir ./migrations up
```

### 4. Start Redis

Start a local Redis instance, or use Docker Compose to run the configured Redis service.

For a local Redis instance, configure:

```env
REDIS_ADDR=redis://localhost:6379
```

### 5. Run the Backend

From the project root:

```bash
go mod download
go run ./cmd/server
```

The backend should be available at:

```text
http://localhost:9000
```

Health endpoint:

```text
http://localhost:9000/health
```

### 6. Configure the Frontend

Open a new terminal:

```bash
cd frontend
npm install
```

Create `frontend/.env`:

```env
VITE_API_BASE_URL=http://localhost:9000
VITE_USE_MOCKS=false
```

### 7. Run the Frontend

```bash
npm run dev
```

Open the local URL provided by Vite, typically:

```text
http://localhost:5173
```

---

## Docker Deployment

Repono includes Docker and Docker Compose configuration to run the application locally.

Ensure the root `.env` file has the required configuration, then run:

```bash
docker compose up --build
```

The configured frontend and backend are available at:

| Service | URL |
|---|---|
| Frontend | http://localhost:8081 |
| Backend | http://localhost:9000 |
| Redis | Internal Docker network |

Stop the services with:

```bash
docker compose down
```

---

## Environment Variables

### Backend

| Variable | Description |
|---|---|
| `HOST` | Backend bind address |
| `PORT` | Backend listening port |
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | JWT signing secret |
| `REDIS_ADDR` | Redis connection URL |
| `AWS_ACCESS_KEY_ID` | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |
| `AWS_REGION` | AWS region |
| `AWS_S3_BUCKET` | S3 bucket name |
| `FRONTEND_ORIGIN` | Allowed frontend origin |

### Frontend

| Variable | Description |
|---|---|
| `VITE_API_BASE_URL` | Backend API base URL |
| `VITE_USE_MOCKS` | Controls frontend mock data |

**Security:** Vite exposes `VITE_` variables in the client-side bundle. Never store passwords, API secrets, or private credentials in frontend environment variables.

---

## Deployment

Repono is deployed using Render and managed cloud services.

| Component | Provider |
|---|---|
| Frontend | Render Static Site |
| Backend | Render Web Service |
| PostgreSQL | Neon |
| Redis | Upstash |
| Object Storage | AWS S3 |

### Live Application

- **Frontend:** [repono-frontend.onrender.com](https://repono-frontend.onrender.com)
- **Backend:** [repono-backend-3uq9.onrender.com](https://repono-backend-3uq9.onrender.com)
- **Health Check:** [API Health](https://repono-backend-3uq9.onrender.com/health)

Production environment variables are configured in the hosting provider's dashboard and should not be committed to the repository.

---

## Security

Repono includes security-focused mechanisms such as:

- JWT-based authentication for protected endpoints.
- Redis-backed rate limiting to restrict excessive requests.
- CORS configuration for controlling browser-based cross-origin access.
- Environment-based configuration for sensitive credentials.
- Public sharing through unique share links.
- Private sharing with specific users.

Security requires continuous review, including authorization checks, upload validation, secret management, and protection of public share links.

---

## Testing

The repository includes Go unit tests for selected backend components, including:

- Redis client
- Redis rate limiter
- File service
- Local storage

Run the Go test suite from the project root:

```bash
go test ./...
```

---

## Future Improvements

- Broader automated unit and integration test coverage.
- File previews for common document and media formats.
- Upload progress tracking and improved large-file handling.
- Structured logging, metrics, and observability.
- CI/CD automation.
- Load testing and performance optimization.
- Additional security hardening.

---

## Contributing

Contributions, suggestions, and bug reports are welcome.

1. Fork the repository.
2. Create a feature branch:

   ```bash
   git checkout -b feature/your-feature
   ```

3. Implement and test your changes.
4. Commit your changes:

   ```bash
   git commit -m "Add: your feature"
   ```

5. Push your branch and open a pull request.

---

## Author

**Jatin Nayak**  
Electronics and Telecommunication Engineering Student | Software Developer

- **GitHub:** [@her0dx-56](https://github.com/her0dx-56)
- **Email:** [her0dev@protonmail.com](mailto:her0dev@protonmail.com)

---

## License

No license has been specified yet. Add a `LICENSE` file and update this section if you choose to distribute Repono under an open-source license.

---

<p align="center">
  Built with Go, React, and ☕
  <br />
  <strong>Repono — Your Files. Your World.</strong>
</p>
