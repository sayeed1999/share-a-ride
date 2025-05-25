# Share-A-Ride

A ride-sharing application built with Go, following clean architecture principles and best practices.

## Project Structure

The project follows a standard Go project layout with clean architecture principles:

```unset
.
├── cmd/                    # Application entry points
│   └── main.go            # Main application entry point
├── internal/              # Private application code
│   ├── app/              # Application core logic
│   ├── config/           # Configuration management
│   ├── domain/           # Business domain models and interfaces
│   ├── mocks/            # Mock implementations for testing
│   ├── pkg/              # Internal packages
│   └── provider/         # External service providers and implementations
├── docs/                 # Documentation files
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── docker-compose.yml    # Docker compose configuration
```

### Directory Details

#### `cmd/`

Contains the main application entry points. This is where the application bootstrap happens and where the dependency injection is configured.

#### `internal/`

Contains the private application code that shouldn't be imported by other projects. It's organized into several subdirectories:

- `app/`: Contains the core application logic and use cases
- `config/`: Handles configuration management and environment variables
- `domain/`: Defines the core business logic and interfaces
- `mocks/`: Contains mock implementations used in testing
- `pkg/`: Houses internal packages that can be shared across the application
- `provider/`: Implements external service providers and adapters

## Architecture

The project follows Clean Architecture principles, with clear separation of concerns:

1. **Domain Layer** (`internal/domain/`)
   - Contains business logic and domain models
   - Defines interfaces for external dependencies
   - No dependencies on external packages

2. **Application Layer** (`internal/app/`)
   - Implements use cases
   - Orchestrates the flow of data
   - Depends only on the domain layer

3. **Infrastructure Layer** (`internal/provider/`)
   - Implements interfaces defined in the domain layer
   - Handles external concerns (database, external services)
   - Contains adapters for external dependencies

## Docker Support

The project includes Docker support for easy deployment and development:

- `Dockerfile`: Multi-stage build for creating optimized production images
- `.dockerignore`: Specifies which files should be excluded from Docker builds
- `docker-compose.yml`: Defines services, networks, and volumes for local development

### Running with Docker

To run the application using Docker:

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d

# Stop all services
docker-compose down
```

### Environment Variables

The application uses the following environment variables (configured in docker-compose.yml):

- `APP_ENV`: Application environment (development/production)
- `POSTGRES_USER`: Database user
- `POSTGRES_PASSWORD`: Database password
- `POSTGRES_DB`: Database name

## Development

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)

### Local Development

1. Clone the repository
2. Install dependencies:

   ```bash
   go mod download
   ```

3. Start the development environment:

   ```bash
   docker-compose up -d
   ```

4. Run the application:

   ```bash
   go run cmd/main.go
   ```

## Testing

To run tests:

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the terms of the LICENSE file included in the repository.
