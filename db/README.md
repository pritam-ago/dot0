# CloudStore Database

Database configuration and migrations for CloudStore.

## 🛠️ Technologies

- PostgreSQL 16
- Redis
- golang-migrate (for migrations)

## 📁 Structure

```
db/
├── migrations/        # Database migrations
└── init.sql         # Initial database setup
```

## 🗄️ Database Schema

### Tables

1. **users**

   - User account information
   - Authentication details

2. **storage_paths**

   - Registered storage locations
   - Links to user accounts

3. **files**

   - File metadata
   - Storage location mapping

4. **file_versions**

   - Version history
   - File change tracking

5. **shares**

   - File sharing information
   - Access control

6. **sync_events**
   - Synchronization status
   - Event tracking

## 🚀 Getting Started

1. **Setup PostgreSQL and Redis**

   Install and configure PostgreSQL and Redis according to your system's package manager.

2. **Run migrations**

   ```bash
   migrate -path db/migrations -database "postgresql://cloudstore:development_password@localhost:5432/cloudstore?sslmode=disable" up
   ```

3. **Reset database**
   ```bash
   migrate -path db/migrations -database "postgresql://cloudstore:development_password@localhost:5432/cloudstore?sslmode=disable" down -all
   ```

## 🔧 Configuration

Environment variables:

- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `REDIS_PASSWORD`
