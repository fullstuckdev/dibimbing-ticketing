-- Create database
CREATE DATABASE IF NOT EXISTS ticketing_system;

-- Use database
USE ticketing_system;

-- Create admin user after running the application
-- INSERT INTO users (name, email, password, role, created_at, updated_at)
-- VALUES ('Admin User', 'admin@example.com', '$2a$10$HASHED_PASSWORD', 'admin', NOW(), NOW());
-- Note: You'll need to generate a proper bcrypt hash for the password using the application 