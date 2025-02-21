# SWIFT Codes API

This is a Go-based REST API for managing SWIFT codes (Bank Identifier Codes). The application allows you to retrieve, add, and delete SWIFT codes, as well as fetch all SWIFT codes for a specific country. The project is containerized using Docker Compose, with a PostgreSQL database for data storage. The database has been denormalized to optimize query performance. This reduces the number of joins and improves read speed at the cost of redundancy.

## Table of Contents

- [Setup](#setup)
- [Running Tests](#running-tests)


## Setup

### Prerequisites

- Docker and Docker Compose installed on your machine.
- Git for cloning the repository.

### Steps

1. **Clone the repository:**

   ```bash
   git clone https://github.com/czyz-bartosz/swift-codes-api.git
   cd swift-codes-api
   ```

2. **Set up environment variables:**

    - Copy the `.env.example` file to `.env`:
      ```bash
      cp .env.example .env
      ```
    - Update the `.env` file with your database credentials:
      ```bash
      DB_USER=your_db_username
      DB_PASSWORD=your_db_password
      DB_NAME=mydb
      ```

3. **Build and run the application using Docker Compose:**

> ⚠️ **Warning: Port Conflicts**  
> Before running the application, ensure that the following ports are not already in use on your machine:
>
> - **Port 8080**: Used by the Go application.
> - **Port 5432**: Used by the PostgreSQL database.
>
> If these ports are occupied, the application will fail to start.


   ```bash
   docker compose up --build
   ```

   This will start the PostgreSQL database and the Go application. The API will be accessible at `http://localhost:8080`.


## Running Tests

> ⚠️ **Warning: Running Integration Tests Will Clear the Database**  
> Please note that running the integration tests will clear the database.  
> To restore the initial state, simply restart the application.

To run the unit and integration tests on your local computer, execute the following command.:

```bash
go test ./...
```

or on the container:

```bash
docker compose exec app go test ./...
```

Make sure the Docker containers are running before running the integration tests.
