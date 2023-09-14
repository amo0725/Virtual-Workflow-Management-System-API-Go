# Virtual Workflow Management System

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technologies](#technologies)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)

## Overview

Virtual Workflow Management System is a RESTful API backend built in Go. It's designed to manage workflows in a virtual environment. With user authentication and Attribute-Based Access Control, this project is a one-stop solution for managing workflows efficiently and securely.

## Features

- User Authentication
- Session Management
- CRUD operations for workflows
- Attribute-Based Access Control (ABAC)
- Transfer Workflow Ownership

## Technologies

- Go
- Gin
- MongoDB
- Redis
- JWT for authentication
- Swagger for API documents
- godotenv for environment variables
- crypto for password hashing

## Prerequisites

- Go v1.21.1
- MongoDB
- Redis
- `git` (Optional)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/amo0725/Virtual-Workflow-Management-System-API-Go.git
   ```

2. Change to project directory:

   ```bash
   cd Virtual-Workflow-Management-System-API-Go
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

## Configuration

All configurations are specified in the `.env` file. Here are some important configurations:

- `HOST`: Host address for the server
- `PORT`: Port number (example: 8080)
- `BASE_PATH`: Base path for API endpoints
- `MONGO_HOST`: MongoDB connection string
- `MONGO_DB_NAME`: MongoDB database name
- `REDIS_USERNAME`: Redis username
- `REDIS_PASSWORD`: Redis password
- `REDIS_HOST`: Redis connection string
- `JWT_SECRET`: JWT secret key

## API Endpoints

- `/api/login`: Authenticate user and create session
- `/api/logout`: Logout and invalidate session
- `/api/register`: Register new user
- `/api/workflows`: CRUD operations for workflows

Refer to the [API Documentation](http://localhost:8080/swagger/index.html) for more details. (Make sure the server is running, and the `PORT` in `.env` file is same as the url.)

## Testing

To run the tests:

```bash
go test ./...
```
