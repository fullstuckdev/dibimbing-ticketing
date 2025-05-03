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

API documentation is available using Swagger annotations and can be accessed at `/swagger/index.html` when running in development mode.
