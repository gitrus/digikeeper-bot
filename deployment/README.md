# Digikeeper Bot Deployment

Production-ready Docker Compose setup with TimescaleDB and database migrations.

## Usage

```bash
#!/bin/bash
# Script to setup secrets directory and files

# Create secrets directory with secure permissions
mkdir -p ./secrets
chmod 700 ./secrets

echo "postgres" > ./secrets/db_user
echo "$(openssl rand -base64 32 | tr -d '\n')" > ./secrets/db_password
echo "your_telegram_token_here" > ./secrets/telegram_token

chmod 600 ./secrets/*
```

## Structure

- `docker-compose.yml` - Container configuration
- `migrations/` - Database migration scripts golang-migrate
