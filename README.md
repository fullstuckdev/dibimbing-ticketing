# Ticketing System API

A RESTful API for a ticketing system built with Go, Gin framework, GORM ORM, and MySQL database.

## Features

- User Management (Registration/Login)
- Event Management (CRUD operations)
- Ticket Management (Purchase/View/Cancel)
- Role-Based Access Control (Admin/User)
- Report Generation
- Pagination and Filtering
- Input Validation

## Tech Stack

- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL
- **Authentication**: JWT

## Project Structure

```
/project-root
  /controller       - HTTP request handlers
  /service          - Business logic
  /repository       - Data access layer
  /entity           - Database models
  /config           - Configuration
  /middleware       - HTTP middleware
  /reports          - Reporting functionality
  /docs             - Swagger documentation
  /postman          - Postman collection
  /tests            - Unit tests
    /service        - Service layer tests
    /repository     - Repository layer tests
    /controller     - Controller layer tests
  main.go           - Application entry point
```

## API Endpoints

### User Management

- `POST /register` - Register a new user
- `POST /login` - User login

### Event Management

- `GET /events` - List all events
- `GET /events/:id` - Get event details
- `POST /events` - Create a new event (Admin only)
- `PUT /events/:id` - Update event (Admin only)
- `DELETE /events/:id` - Delete event (Admin only)

### Ticket Management

- `GET /tickets` - List tickets (Admin sees all, users see their own)
- `POST /tickets` - Purchase a ticket
- `GET /tickets/:id` - View ticket details
- `PATCH /tickets/:id` - Cancel a ticket

### Reports (Admin only)

- `GET /reports/summary` - Get overall sales report
- `GET /reports/event/:id` - Get event-specific report

## Authentication

The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <token>
```

## Role-Based Access

- **Admin**: Full access to all endpoints
- **User**: Can view events, purchase/view/cancel their own tickets

## Setup Instructions

1. Clone the repository
2. Create `.env` file with the following variables:
   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=password
   DB_NAME=ticketing_system
   JWT_SECRET=your_jwt_secret
   PORT=8080
   ```
3. Create the MySQL database
   ```sql
   CREATE DATABASE ticketing_system;
   ```
4. Run the application
   ```
   go run main.go
   ```

## Running in Development

```bash
go run main.go
```

The server will start at http://localhost:8080 (or the port specified in your .env file).

## API Documentation

### Swagger UI

The API is documented using Swagger. After starting the application, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

This provides an interactive documentation where you can view all endpoints and even test them.

### Postman Collection

A Postman collection is included in the `postman` directory. To use it:

1. Open Postman
2. Click on "Import" button
3. Select the file `postman/ticketing_system_api.json`
4. After importing, set the necessary environment variables:
   - `base_url`: Your API base URL (default: http://localhost:8080)
   - `token`: JWT token obtained after login as a regular user
   - `admin_token`: JWT token obtained after login as an admin user

## Development with Hot Reload

This project supports hot reloading using [Air](https://github.com/air-verse/air), which automatically rebuilds and restarts the application when file changes are detected.

### Installing Air

```bash
# Install Air globally
go install github.com/air-verse/air@latest

# Add Go bin directory to your PATH if it's not already added
# For bash/zsh (add to ~/.bashrc or ~/.zshrc):
# export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
$(go env GOPATH)/bin/air -v
# Or just 'air -v' if Go bin is in your PATH
```

### Running with Air

```bash
# If Go bin is in your PATH:
air

# Or using the full path:
$(go env GOPATH)/bin/air
```

Air will watch for any changes in your Go files and automatically rebuild and restart the server, making development much faster.
