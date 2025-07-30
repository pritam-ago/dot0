# Docker Setup for Relay Project

This document provides instructions for containerizing and running the entire relay project using Docker.

## Project Structure

The project consists of three main components:

- **go-relay**: Go backend server handling WebSocket connections and database operations
- **web-client**: React frontend for web users
- **pc-app**: Tauri-based desktop application for PC users

## Prerequisites

- Docker (version 20.10 or higher)
- Docker Compose (version 2.0 or higher)

## Quick Start

1. **Clone the repository and navigate to the project root:**

   ```bash
   cd /path/to/your/project
   ```

2. **Build and start all services:**

   ```bash
   docker-compose up --build
   ```

3. **Access the services:**
   - Go Relay Server: http://localhost:8080
   - Web Client: http://localhost:3001
   - PC App: http://localhost:3002
   - PostgreSQL Database: localhost:5432

## Service Details

### PostgreSQL Database

- **Port**: 5432
- **Database**: relay_db
- **Username**: relay_user
- **Password**: relay_password
- **Data Persistence**: Uses Docker volume `postgres_data`

### Go Relay Server

- **Port**: 8080
- **Health Check**: http://localhost:8080/health
- **Dependencies**: PostgreSQL
- **Features**: WebSocket connections, PIN management, session handling

### Web Client

- **Port**: 3001
- **Health Check**: http://localhost:3001/health
- **Features**: React SPA with nginx serving
- **Dependencies**: Go Relay Server

### PC App

- **Port**: 3002
- **Health Check**: http://localhost:3002/health
- **Features**: Tauri-based desktop app served via nginx
- **Dependencies**: Go Relay Server

## Docker Commands

### Build and Start

```bash
# Build and start all services
docker-compose up --build

# Start in detached mode
docker-compose up -d --build

# Start specific service
docker-compose up --build go-relay
```

### Stop and Cleanup

```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Stop and remove images
docker-compose down --rmi all
```

### Logs and Monitoring

```bash
# View logs for all services
docker-compose logs

# View logs for specific service
docker-compose logs go-relay

# Follow logs in real-time
docker-compose logs -f

# Check service status
docker-compose ps
```

### Individual Service Management

```bash
# Rebuild specific service
docker-compose build go-relay

# Restart specific service
docker-compose restart go-relay

# Execute commands in running container
docker-compose exec go-relay sh
docker-compose exec postgres psql -U relay_user -d relay_db
```

## Environment Variables

### Go Relay Server

The following environment variables are automatically set by Docker Compose:

- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_NAME=relay_db`
- `DB_USER=relay_user`
- `DB_PASSWORD=relay_password`

### PostgreSQL

- `POSTGRES_DB=relay_db`
- `POSTGRES_USER=relay_user`
- `POSTGRES_PASSWORD=relay_password`

## Development Workflow

### Making Changes

1. Make your code changes
2. Rebuild the affected service:
   ```bash
   docker-compose build <service-name>
   docker-compose up <service-name>
   ```

### Hot Reloading (Development)

For development with hot reloading, you can override the Dockerfile commands:

```bash
# For Go Relay (development mode)
docker-compose run --service-ports go-relay go run ./cmd/server

# For Web Client (development mode)
docker-compose run --service-ports web-client npm start

# For PC App (development mode)
docker-compose run --service-ports pc-app npm run dev
```

## Troubleshooting

### Common Issues

1. **Port Already in Use**

   ```bash
   # Check what's using the port
   netstat -tulpn | grep :8080

   # Kill the process or change ports in docker-compose.yml
   ```

2. **Database Connection Issues**

   ```bash
   # Check if PostgreSQL is running
   docker-compose ps postgres

   # Check PostgreSQL logs
   docker-compose logs postgres
   ```

3. **Build Failures**

   ```bash
   # Clean build cache
   docker-compose build --no-cache

   # Remove all images and rebuild
   docker-compose down --rmi all
   docker-compose up --build
   ```

4. **Permission Issues**
   ```bash
   # Fix file permissions
   sudo chown -R $USER:$USER .
   ```

### Health Checks

All services include health checks. You can monitor them:

```bash
# Check health status
docker-compose ps

# View health check logs
docker inspect <container-name> | grep -A 10 "Health"
```

## Production Deployment

### Environment Variables

For production, create a `.env` file:

```env
POSTGRES_PASSWORD=your_secure_password
DB_PASSWORD=your_secure_password
```

### Security Considerations

1. Change default passwords
2. Use secrets management
3. Enable SSL/TLS
4. Configure firewall rules
5. Regular security updates

### Scaling

```bash
# Scale specific service
docker-compose up --scale go-relay=3

# Use Docker Swarm for production orchestration
docker stack deploy -c docker-compose.yml relay-stack
```

## Monitoring and Logging

### Log Aggregation

```bash
# View all logs with timestamps
docker-compose logs -t

# Filter logs by service and time
docker-compose logs --since="2024-01-01T00:00:00" go-relay
```

### Resource Monitoring

```bash
# Monitor container resources
docker stats

# Check disk usage
docker system df
```

## Backup and Restore

### Database Backup

```bash
# Create backup
docker-compose exec postgres pg_dump -U relay_user relay_db > backup.sql

# Restore backup
docker-compose exec -T postgres psql -U relay_user relay_db < backup.sql
```

### Volume Backup

```bash
# Backup PostgreSQL volume
docker run --rm -v relay_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz -C /data .

# Restore PostgreSQL volume
docker run --rm -v relay_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres_backup.tar.gz -C /data
```

## Support

For issues related to:

- **Docker setup**: Check this README and Docker documentation
- **Application logic**: Check the individual service documentation
- **Database issues**: Check PostgreSQL logs and connection settings
