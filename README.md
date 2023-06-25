# ais

Simple web service for manage articles

## Setup

1. To setup first you need to run docker compose on project root directory to build container for database (MariaDB) and cache (Redis).

```
docker compose up -d
```

2. Then create new database table or you can run following SQL query.

```
CREATE DATABASE `ais_db` DEFAULT CHARACTER SET = `utf8mb4` DEFAULT COLLATE = `utf8mb4_general_ci`;
```

3. Make sure your redis container is running and can be access using the credential in docker-compose.yaml

4. Then copy .env.example file as .env file and change some of the variables according to your local setup

```
cp .env.example .env
```

5. Then run the Go project

```
go run main.go
```
