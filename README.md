# Order Service

This is a microservice for managing customer orders. It provides APIs for creating, retrieving, updating, and deleting orders. The service also includes a background job that synchronizes order data with Meilisearch on startup.

## Table of Contents
- [Meilisearch Sync Job](#meilisearch-sync-job)
- [Running the Service](#running-the-service)
- [Environment Variables](#environment-variables)
- [Database Initialization](#database-initialization)
- [API](#api)
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
- [Wire](https://github.com/google/wire) (for dependency injection)

### **Generate Wire Dependencies**
The wire dependency injection files can be found in the deps folder.

Install wire
```sh
go install github.com/google/wire/cmd/wire@latest
```

To create the dependency injection file launch:
```sh
./wiregen.sh
```

### **Start the Service**
```sh
docker-compose up --build
```
This starts the following services:
- Redis
- MariaDB
- Meilisearch
- Order Service

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

## Api

Post Create Order

```sh 
curl -X POST "http://localhost:8080/api/v1.0/order/" \
     -H "Content-Type: application/json" \
     -d '{
        "user_id": 1,
        "items": [
            {"product_id": 4, "quantity": 20},
            {"product_id": 5, "quantity": 30}
        ]
     }'
```

Get Order
```sh
curl -X GET "http://localhost:8080/api/v1.0/order/1"
```

Delete Order

```sh
curl -X DELETE "http://localhost:8080/api/v1.0/order/1"
```

Update Order
```sh
curl -X PUT "http://localhost:8080/api/v1.0/order/1" \
    -H "Content-Type: application/json" \
    -d '{
      "items": [
        {"product_id": 4, "quantity": 10},
        {"product_id": 5, "quantity": 15}]
    }'
```

List Order
```sh
curl -X GET "http://localhost:8080/api/v1.0/order?input=laptop&start_date=2025-03-29T12:30:00Z&end_date=2025-05-29T14:30:00Z&limit=10&offset=0" \
```

## Swagger

To generate the swagger use this to convert the postman apis to swagger
and run the swaggen.sh

https://www.npmjs.com/package/postman-to-openapi
