# Inventory Management API

This project is the starting point for building Inventory Management API using [Go](https://go.dev/), the [Gin](https://github.com/gin-gonic/gin) framework, and [GORM](https://gorm.io/index.html) for database operations.

## Setting Up

1. Clone this project to your local machine, then navigate to the **/starter** directory.

2. Set up a PostgreSQL database and Redis on your local machine.

   - For PostgreSQL, create a database named `postgres` and update the connection details in the `.env` file.
   ```
   psql postgres postgres -h localhost
   ```
   - For Redis, ensure the Redis server is running on the default port.

    ```
    brew install redis
    brew services start redis
    ```

3. Before running the server, ensure all dependencies are properly managed.

    ```
    go mod tidy
    ```

4. Run the server.

    ```
    go run main.go
    ```

5. Run the tests

    ```
    go test ./tests/...
    ```

The server will connect to your PostgreSQL database and seed it with 20 sample items. Be sure to check the database to confirm that the items have been seeded successfully.
## API Endpoints. 
1. Get All Items with Pagination, Sorting, and Filtering


- Get all items
    ```
    curl http://localhost:8080/api/v1/inventory
    ```
- Get a JWT token
    ```
    curl -X POST http://localhost:8080/api/v1/login \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "password"}'
    ```
- Create new item
    ```
    curl -X POST http://localhost:8080/api/v1/inventory \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_TOKEN_HERE" \
    -d '{"name": "New Item", "stock": 10, "price": 29.99}'
    ```
- Get item by ID
    ```
    curl http://localhost:8080/api/v1/inventory/{id}
    ```

- Update item
    ```
    curl -X PUT http://localhost:8080/api/v1/inventory/{id} \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_TOKEN_HERE" \
    -d '{"name": "Updated Item", "stock": 15, "price": 39.99}'
    ```
- Delete item
    ```
    curl -X DELETE http://localhost:8080/api/v1/inventory/{id} \
    -H "Authorization: Bearer YOUR_TOKEN_HERE"
    ```
2. Rate Limiting Test

    ```
    for i in {1..10}; do curl "http://localhost:8080/api/v1/inventory"; echo; done
    ```
## Project Structure

```
Inventory-Management-API/
├── README.md
├── starter/
│   ├── .env
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── database/
│   │   └── cache.go
│   │   └── database.go
│   ├── handlers/
│   │   ├── auth.go
│   │   └── item_handler.go
│   ├── middleware/
│   │   ├── jwt.go
│   │   └── rate_limiter.go
│   ├── models/
│   │   └── item.go
│   ├── routes/
│   │   └── routes.go
│   ├── tests/
│   │   └── api_test.go
```
## Key Features Implemented. 
`Rate Limiting`: Token bucket algorithm with 1 request/second refill rate and burst capacity of 5.  
`Pagination`: Page-based pagination with metadata.  
`Sorting`: Sort by name, stock, or price in ascending/descending order.  
`Filtering`: Filter by minimum stock and name (case-insensitive partial match).  
`Error Handling`: Improved error responses and validation.  
`Performance`: Optimized database queries and connection pooling.  
`JWT Authentication`: Secure API endpoints with JSON Web Tokens.  
`In-memory Caching`: Improve performance with Redis in-memory caching for frequently accessed data.
## Query Parameters
`page`: Page number (default: 1).  
`page_size`: Items per page (default: 10, max: 100).  
`sort_by`: Field to sort by (name, stock, price).  
`sort_order`: Sort direction (asc, desc).  
`min_stock`: Minimum stock filter (default: 0).  
`name`: Name filter (partial match, case-insensitive).  
The response includes pagination metadata to help clients build proper pagination controls.
