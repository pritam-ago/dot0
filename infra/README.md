# CloudStore Infrastructure

Infrastructure configuration for CloudStore using Docker Compose and Traefik.

## 🛠️ Technologies

- Docker Compose
- Traefik (Reverse Proxy)
- PostgreSQL
- Redis

## 📁 Structure

```
infra/
├── docker-compose.yml  # Main Docker Compose configuration
└── traefik.yml        # Traefik reverse proxy configuration
```

## 🚀 Getting Started

1. **Start the infrastructure**

   ```bash
   docker-compose up -d
   ```

2. **Access services**

   - Backend API: http://localhost:8080
   - Traefik Dashboard: http://localhost:8081
   - PostgreSQL: localhost:5432
   - Redis: localhost:6379

3. **Stop services**
   ```bash
   docker-compose down
   ```

## 🔧 Configuration

### Environment Variables

Backend:

- `POSTGRES_*` - PostgreSQL connection settings
- `REDIS_*` - Redis connection settings

Database:

- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`

Redis:

- `REDIS_PASSWORD`

### Volumes

- `postgres-data`: PostgreSQL data persistence
- `redis-data`: Redis data persistence

### Networks

- `cloudstore-network`: Internal Docker network for service communication

## 📝 Notes

- Development environment uses insecure defaults
- Production deployment should:
  - Use secure passwords
  - Enable SSL/TLS
  - Configure proper access controls
  - Use production-ready volume configurations
