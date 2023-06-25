Simple web service to manage articles using **MariaDB** for database and **Redis** for cache

Import `postman_collection.json` to get following available endpoints in your Postman

| Method | Route | Description |
| -- | -- | -- |
| GET | /articles | List articles |
| GET | /articles/:id | View article |
| POST | /articles | Post new article |
| POST | /articles/update | Update article |
| POST | /articles/delete | Delete article |


## Setup

1. To setup first you need to run docker compose on project root directory to build container for **database** and **cache**.

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

## Testing

1. To do test, you need copy .env.test.example and change some of the variables according to your local setup. Make sure app port not conflicted with other services.

```
cp .env.test.example .env.test
```

2. Then go to `handler` directory and run Go test

```
cd handler
go test -v
```
