### Run via Docker

**docker-compose.yml**
```bash
  db_backup:
    build: .
    env_file:
      - .env
    environment:
      - SSH_USER=${SSH_USER}
      - SSH_HOST=${SSH_HOST}
      - SSH_PORT=${SSH_PORT}
      - SSH_PASSWORD=${SSH_PASSWORD}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_HOST=${DB_HOST}
      - DB_NAME=${DB_NAME}
    volumes:
      - ./backups:/app/backups
```

Build
```bash
docker compose up -d
```

Run
```bash
docker compose run --rm db_backup
```