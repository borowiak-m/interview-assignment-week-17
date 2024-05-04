# MongoDB Records Management API
### Overview
The MongoDB Records Management API is a backend system designed to manage MongoDB documents and handle in-memory data through RESTful endpoints. Developed using Go standard lib, it demonstrates interaction with MongoDB for record retrieval based on specific query parameters and manages temporary in-memory data. The system is containerized using Docker.
The project is part of an assignment in "Total Coder" skool community: https://www.skool.com/totalcoder

## Features
- MongoDB Operations: Facilitates querying MongoDB to retrieve records by date and count criteria.
- In-Memory Data Handling: Manages temporary in-memory data storage and retrieval.
- Custom Response Structure: Constructs API responses with status codes and detailed messages.
- Docker Integration: Utilizes Docker and Docker Compose for environment setup and dependency management.

## How to run it
### Prerequisites
To run this project, you will need:

- Docker and Docker Compose to manage the application and MongoDB instance.
- Go installed on your machine (Version 1.15 or later).
- An HTTP client (e.g., curl, Postman) for interacting with the API.

### Environment setup
Clone the Repository: Clone the project to your local machine.
Docker Compose: Navigate to the project directory and use docker-compose up to start the services defined in docker-compose.yml

### API Endpoints
#### POST /fetchMongoRecords
Request payload: JSON with startDate, endDate, minCount, and maxCount.
Response: JSON with code, msg, and an array of records each containing key, createdAt, and totalCount.

Example request payload: 

```
{
    "startDate": "2024-01-01",
    "endDate": "2024-12-31",
    "minCount": 100,
    "maxCount": 800
}
```

Example request:

```
curl -X POST http://localhost:3000/fetchMongoRecords \
-H "Content-Type: application/json" \
-d '{
    "startDate": "2024-01-01",
    "endDate": "2024-12-31",
    "minCount": 100,
    "maxCount": 800
}'
```

Example response:

```
{
  "code": 0,
  "msg": "Success",
  "records": [
    {
      "key": "TAKwGc6Jr4i8Z487",
      "createdAt": "2024-01-05T14:15:22Z",
      "totalCount": 2050
    },
    {
      "key": "NAeQ8eX7e5TEg70H",
      "createdAt": "2024-02-15T14:15:22Z",
      "totalCount": 1850
    }
  ]
}
```

#### GET /inmemory
Response: JSON array of all in-memory records with key and value.
Example request: 

```
curl http://localhost:3000/inmemory
```

Example response:

```
[
  {
    "key": "active-tabs",
    "value": "getir"
  },
  {
    "key": "inactive-tabs",
    "value": "getout!"
  }
]
```

#### POST /inmemory
Request payload: JSON with key and value
Response: JSON with key and value from request

Example request payload:

```
{
  "key": "session-data",
  "value": "user123"
}
```

Example request: 

```
curl -X POST http://localhost:3000/fetchMongoRecords -d '{"startDate":"2024-01-01", "endDate":"2024-12-31", "minCount":100, "maxCount":800}'
```

### Running the project 
Start the Server: After setting up Docker, run the Docker Compose command which builds the Go application and starts the MongoDB service.
Interact with the API: Use the endpoints to interact with the database and in-memory data.

### Docker configuration
- Dockerfile: Builds a Go environment, installs dependencies, and compiles the application.
- docker-compose.yml: Defines services for the app and MongoDB, including volume mounts for persistent MongoDB data and initialization scripts.




