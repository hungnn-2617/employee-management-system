# Employee Management System

An employee management system built with Go and MySQL.

## Requirements

- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0+

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd employee-management-system
```

2. Configure environment:

```bash
cp .env.example .env
# Edit the .env file with your database credentials
```

3. Install dependencies:

```bash
go mod download
```

4. Start the database with Docker:

```bash
docker-compose up -d
```

5. Run the application:

```bash
go run internal/cmd/app/main.go
```
