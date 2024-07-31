# User Management Microservice

This microservice provides an HTTP API to manage user data. It supports adding, updating, deleting, and listing users with pagination. The service is implemented in Go and uses an in-memory storage mechanism.

## Requirements

- Go 1.16 or later
- Docker

## Clone the Repository

```git clone https://github.com/John-AG/user-service```

```cd user-service```

## Build the Docker image

```docker build -t user-service .```

## Run the application in Docker

```docker run -p 8080:8080 user-service```

## Features

- Add a new User
- Modify an existing User
- Remove a User
- Return a paginated list of Users
- Health check endpoint
- Logging

## User Schema

The user data is stored using the following schema:

json
{
    "id": "d2a7924e-765f-4949-bc4c-219c956d0f8b",
    "first_name": "Alice",
    "last_name": "Bob",
    "nickname": "AB123",
    "password": "supersecurepassword",
    "email": "alice@bob.com",
    "country": "UK",
    "created_at": "2019-10-12T07:20:50.52Z",
    "updated_at": "2019-10-12T07:20:50.52Z"
}

## API Endpoints

## Health Check

```curl -X GET http://localhost:8080/health```

## User Management

- Add a New User.

```
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{
    "first_name": "Alice",
    "last_name": "Smith",
    "nickname": "ASmith",
    "password": "supersecurepassword",
    "email": "alice@smith.com",
    "country": "UK"
}'
```

- List Users with Pagination.

```curl -X GET "http://localhost:8080/users?page=1&pageSize=10"```

- List Users with Filtering.

```curl -X GET "http://localhost:8080/users?page=1&pageSize=10&country={country}"```

- Update an existing user.

```
curl -X PUT http://localhost:8080/users/{id} -H "Content-Type: application/json" -d '{
    "first_name": "AliceUpdated",
    "last_name": "Smith",
    "nickname": "AUpdated",
    "password": "newpassword",
    "email": "alice.updated@smith.com",
    "country": "US"
}'
```

- Delete a user.

```curl -X DELETE http://localhost:8080/users/{id}```

## Testing

Unit tests are included in the main_test.go file. To run the tests, use the following command:

```go test```

## Explanation of choices

HTTP vs gRPC API

I chose HTTP due to the simplicity of the application, the improved performance that gRPC can provide wasn’t necessary. I was also much more familiar with HTTP so in terms of required performance and time, it was what made the most sense for this task.

Docker vs Direct Execution

There are many reasons why using Docker is superior to direct execution, such as advantages in terms of environment consistency, isolation, portability, scalability, and simplified dependency management. The key one for this task being environment consistency, ensuring that there are no issues for the user when running this API on their machine. Additionally it was mentioned that this was the preferred method on the task. 

## Possible extensions

- Persistent Storage: Replace in-memory storage with a persistent database to ensure data persistence.
- Advanced Filtering: Implement advanced filtering options such as partial matches and ranges (e.g., date ranges).
- Full-Text Search: Integrate a full-text search engine like Elasticsearch for more powerful search capabilities.
- Sorting: Add support for sorting the results by different fields (e.g., first name, last name, created date). This was not implemented to avoid the solution becoming unnecessarily complex.
- Throttling: Implement throttling to control the number of requests a client can make in a given time period.
