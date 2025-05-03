# Ticketing System API Testing with Postman

This guide explains how to test the Ticketing System API using Postman.

## Prerequisites

- [Postman](https://www.postman.com/downloads/) installed on your system
- Ticketing system backend running on `http://localhost:8080`

## Setup

1. Import the collection file `ticketing_system_api_tests.postman_collection.json` into Postman:

   - Click "Import" in Postman
   - Select the collection file
   - Click "Import"

2. Import the environment file `ticketing_system_environment.postman_environment.json`:

   - Click "Import" in Postman
   - Select the environment file
   - Click "Import"

3. Select the "Ticketing System" environment from the environment dropdown in the top-right corner of Postman

## Running Tests

The collection is organized into folders by feature. You can run tests in sequence by following the order below:

### Test Flow

1. **Authentication**

   - Register User: Creates a regular user account
   - Register Admin: Creates an admin user account
   - Login (User): Logs in as a regular user and automatically saves the token
   - Login (Admin): Logs in as an admin and automatically saves the token
   - Get User Profile: Confirms you can access user data with the token
   - (Leave Logout for the end of testing)

2. **Events**

   - Create Event (Admin): Creates a test event and saves the event ID for later use
   - Get All Events: Lists all events
   - Get Event by ID: Gets details of the specific event
   - Update Event (Admin): Modifies the test event
   - (Leave Delete Event for the end of testing)

3. **Tickets**

   - Purchase Ticket: Buys a ticket for the test event and saves the ticket ID
   - Get All Tickets: Lists all tickets for the current user
   - Get Ticket by ID: Gets details of the specific ticket
   - (Leave Cancel Ticket for the end of testing)

4. **Audit Logs**

   - Get My Audit Logs: View the current user's activity logs
   - Get All Audit Logs (Admin): View all users' activity logs
   - Get Entity Audit Logs (Admin): View logs for a specific entity

5. **Reports**
   - Get Sales Report (Admin): View overall sales information
   - Get Event Sales Report (Admin): View sales for a specific event

### Cleanup Testing

After the main tests, you can test the following endpoints:

1. Cancel Ticket: Tests ticket cancellation
2. Delete Event (Admin): Tests event deletion
3. Logout: Tests user logout

## Automatic Variable Setting

When you run certain requests, variables are automatically set for later use:

- After Login: `user_token` and `admin_token`
- After Create Event: `event_id`
- After Purchase Ticket: `ticket_id`

## Troubleshooting

- If you receive authentication errors, ensure you've run the Login requests and that tokens were properly saved
- If you get "Not Found" errors for specific IDs, make sure you've created the resources first
- Check that your server is running on the correct port (8080)

## Additional Notes

- The event creation uses `{{$isoTimestamp}}` for start and end dates, which automatically sets them to the current time. For production testing, you may want to modify these to future dates.
- The audit logging endpoints include optional query parameters that are disabled by default. Enable them to test filtering capabilities.
