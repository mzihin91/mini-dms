# Device Management System

A dockerized microservices-based device management system that allows users to register devices (Access Controllers, Face Recognition Readers, ANPR cameras), activate them to spawn independent simulation workers, and monitor transaction generation in real-time.

## Features

✅ **Device Registration** - Register devices with unique names and IP addresses  
✅ **Multi-Device Support** - Access Controllers, Face Recognition Readers, ANPR cameras  
✅ **Device Activation** - Spawn simulation workers that generate transactions  
✅ **Real-Time Monitoring** - View transactions as they are generated  
✅ **Single-Command Deployment** - Docker Compose orchestration

## Architecture

- **Frontend**: React 18+ with TailwindCSS
- **Backend**: Go 1.21+ REST API with gorilla/mux
- **Database**: PostgreSQL 15+ with JSONB support
- **Deployment**: Docker Compose with bridge networking

## Quick Start

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Git 2.30+

### Installation

1. Clone the repository:
```bash
git clone <repository-url> mini-dms
cd mini-dms
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Build and start all services:
```bash
docker compose up --build
```

### Access the Application

- **Frontend UI**: http://localhost:3000
- **Backend API**: http://localhost:8080/api/v1
- **Health Check**: http://localhost:8080/api/v1/health

## Usage

### Registering a Device

1. Navigate to the Devices page
2. Fill in the device form:
   - Name: Unique device identifier
   - Type: Select from Access Controller, Face Reader, or ANPR
   - IP Address: Valid IPv4/IPv6 address
3. Click "Create Device"

### Activating a Device

1. Find the device in the device list
2. Click the "Activate" button
3. The device status changes to "Active" and begins generating transactions

### Monitoring Transactions

1. Navigate to the Transactions page
2. View real-time transaction feed (auto-refreshes every 3 seconds)
3. Filter by device or event type as needed

## API Endpoints

### Devices

- `GET /api/v1/devices` - List all devices
- `POST /api/v1/devices` - Create a new device
- `GET /api/v1/devices/{id}` - Get device by ID
- `DELETE /api/v1/devices/{id}` - Delete a device
- `POST /api/v1/devices/{id}/activate` - Activate a device
- `POST /api/v1/devices/{id}/deactivate` - Deactivate a device

### Transactions

- `GET /api/v1/transactions` - List all transactions
- `GET /api/v1/transactions?device_id={id}` - Filter by device
- `GET /api/v1/transactions?event_type={type}` - Filter by event type

### Health

- `GET /api/v1/health` - Service health check

## API Examples

### Create an Access Controller

```bash
curl -X POST http://localhost:8080/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Front Door AC",
    "device_type": "access_controller",
    "ip_address": "192.168.1.10"
  }'
```

### Activate a Device

```bash
curl -X POST http://localhost:8080/api/v1/devices/1/activate
```

### List Transactions

```bash
curl http://localhost:8080/api/v1/transactions
```

## Development

### Project Structure

```
backend/
├── cmd/api/            # Application entrypoint
├── internal/
│   ├── models/         # Domain models
│   ├── handlers/       # HTTP handlers
│   ├── services/       # Business logic
│   ├── repository/     # Database access
│   ├── deviceworker/   # Device simulation
│   ├── database/       # DB connection
│   └── logging/        # Structured logging
└── migrations/         # SQL migrations

frontend/
├── src/
│   ├── components/     # Reusable UI components
│   ├── pages/          # Page-level components
│   └── services/       # API client
└── public/             # Static assets
```

## Troubleshooting

### Port Already in Use

If ports 3000, 5432, or 8080 are already in use:
```bash
# Stop conflicting services or change ports in docker-compose.yml
docker compose down
```

### Database Connection Issues

Check database container is running:
```bash
docker compose ps
docker compose logs db
```

### Backend Not Starting

View backend logs:
```bash
docker compose logs backend
```

## License

MIT

## Support

For issues and questions, please open an issue in the repository.
