# Example

This folder contains an example of implementation of the package.

## Environment Variables

The database credentials will be read from the following environment variables:

- DB_USER: database user
- DB_PASSWORD: database password
- DB_HOSTNAME: database hostname
- DB_PORT: database port
- DB_SCHEMA: database schema

## Running the example

`DB_USER=<user> DB_HOSTNAME=<host> DB_PORT=<port> DB_SCHEMA=<schema> DB_PASSWORD=<pwd> go run ./examples/complete`
