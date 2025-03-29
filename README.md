# Order Service

This is a microservice for managing customer orders. It provides APIs for creating, retrieving, updating, and deleting orders. The service also includes a background job that synchronizes order data with Meilisearch on startup.

## Table of Contents
- [Meilisearch Sync Job](#meilisearch-sync-job)
- [Running the Service](#running-the-service)
- [Environment Variables](#environment-variables)
- [Database Initialization](#database-initialization)
- [Swagger](#swagger)


## Meilisearch Sync Job

On startup, the service synchronizes existing order data with Meilisearch to ensure search accuracy. The sync job:
1. Fetches all orders from the database.
2. Indexes them into Meilisearch.
3. Runs in the background to ensure search data is always up to date.

## Running the Service

### **Prerequisites**
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

### **Start the Service**
```sh
docker-compose up --build
```
This starts the following services:
- Redis
- MariaDB
- Meilisearch
- Order Service

The Order Service will start  once the database is healthy and the other 2 services have started

### **Stop the Service**
```sh
docker-compose down
```

## Environment Variables
| Variable                 | Description                         | Default Value  |
|--------------------------|---------------------------------|----------------|
| `DATABASE_HOST`          | Database host                   | `mariadb`      |
| `DATABASE_PORT`          | Database port                   | `3306`         |
| `DATABASE_USER`          | Database username               | `root`         |
| `DATABASE_PASSWORD`      | Database password               | `password`     |
| `REDIS_HOST`             | Redis host                      | `redis`        |
| `REDIS_PORT`             | Redis port                      | `6379`         |
| `MEILISEARCH_HOST`       | Meilisearch host               | `meilisearch`  |
| `MEILISEARCH_PORT`       | Meilisearch port               | `7700`         |
| `MEILISEARCH_MASTER_KEY` | Meilisearch port               | MASTER_API_KEY |

## Database Initialization
The database is initialized using an `init.sql` file, which is automatically executed when MariaDB starts.

To manually connect to the database:
```sh
docker exec -it <mariadb-container-id> mysql -u root -p
```

## Swagger

To generate the swagger use this to convert the postman apis to swagger
and run the swaggen.sh

https://www.npmjs.com/package/postman-to-openapi