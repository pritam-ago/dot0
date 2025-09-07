-- Drop tables in reverse order to handle dependencies
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS pairing_sessions;
DROP TABLE IF EXISTS user_plans;
DROP TABLE IF EXISTS users;

-- Drop the UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
