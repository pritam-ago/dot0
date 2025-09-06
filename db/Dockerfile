# Use the official PostgreSQL image
FROM postgres:16-alpine

# Add custom initialization scripts
COPY ./init.sql /docker-entrypoint-initdb.d/

# Set default environment variables
ENV POSTGRES_DB=cloudstore
ENV POSTGRES_USER=cloudstore
ENV POSTGRES_PASSWORD=development_password

# Expose PostgreSQL port
EXPOSE 5432
