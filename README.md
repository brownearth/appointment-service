# Appointment Service: Future's Take Home Assignment

Appointment management service, RESTful HTTP API using JSON, built with Go.

This is a solution for Future's take home project.



## 📣 Jeffrey's Feedback on The Assignment

If I was an interviewer or a reviewer of an assignment's submission, I would have a handful of questions.  So, I will put those questions here and answer them and I hope that makes reviewing a little easier.


## 😏 How did I feel about the project?
It was a wonderful experience. Regardless of how the interviewing process goes, this was great.

I don't get to work with Go much at my current position, and I only have a small amount Go experience on personal projects.  So, I THOROUGHLY enjoyed working on this. 

I wanted the submission to inlcude more than just "Get it Done" code. This created a need for thinking through design choices, trade-offs, what to include and what to note as potential future improvements. So, regardless of how the interview goes, thank you for this.

I really enjoyed working on how the `makefile` is grabbing the values from an `.environment/*` and exporting them such that they can be loaded as environment vars to the service.  The env vars in QA, staging or Production environment likey are set via secrets and never committed.  But having these simple files makes dev'ing relatively simple. Also the setting of VERSION, BUILD_TIME and COMMIT_SHA vi `-ldflags` is something I just learned and it was fun to use that to set some standard attributes in the structured logger.

I also loved setting up air for hot restarts and setting up the migrate tool for DB DDL.  These made iterating fast and enjoyable.

Did I go a little overboard: ✅

Would I do things slightly different next time: ✅

## How long did you work on this?
Longer than I would like to admit.  Probably around 20 hours.  And maybe 2 of thosee 20 on this readme :) 

But it was enjoyable and I learned a bunch.

## What were the hardest parts?
*Timezones* I wanted the API to only deal with UTC, and there service should convert to pacific for validation reasons. Some of the auto-magic marshalling was converting to local TZ.

*Scaffolding, Tooling, etc* I wanted this to feel a little more "real world" and so setting up reasonable make targets, environment handling, project layout, etc was a tough one.  And time consuming.

## Did you use AI?
Absolutely.  I use copilot embedded in VS-Code, and I use Claude.  They are wonderful, and wonderfully misleading and plain wrong at times.  But I do think any good engineer needs to lean on AI to become even more proficient. I personally learn more interacting with AI than I would googling and using StackOverflow.

I do think engineers need to be cautious though, on screen-shares I have seen API keys, PII, etc pasted into ChatGPT. 


## 🏗️↔️🏗️ How did I address Separation of Concerns (SoC)?
The email that came with the assignment spoke to separation of concerns.  That was handled a few different ways.

### 🥞 SoC: Layer Isolation
Layer isolation basically means you have clear boundaries between the "layers" of the software. In this project those layers can be found via directory:

```
internal/
├── api/           # HTTP handlers and server setup
├── service/       # Business logic layer
├── repository/    # Data access layer
├── model/         # Domain entities
└── dto/           # Data transfer objects
```

### 📝 SoC: Code to Interfaces (not implementations)
Coding to an interface allows the underlying implementation to change without affecting the caller.  This allows for dependency injection.  For example, a service typically depends upon a repository (persistance later).  When the service is consctructed, an implementation of the repository can be provided.  This can be done dynamically also, as shown in this project

#### Repository Pattern
The repository layer is interface-driven with three implementations:
1. **Memory Repository**: 
   - A simple, home-grown, in-memory storage layer.
   - Enables rapid development and testing
   - Perfect for local development
   - Fast unit tests
2. **SQLite3**:
   - File-based, serverless database solution
   - Zero-configuration deployment
   - Ideal for small to medium workloads
   - Great for prototyping and testing environments
3. **PostgreSQL Repository**:
   - **NOTE:** Not implemented yet, just stubbed out
   - **NOTE:** Some production grade RDBMS when running service in production
   - Production-ready implementation
   - Docker support for local testing
   - Suitable for QA and production environments

The determination of which concrete implementation to construct via factory is decided base on configuration. See environment files:
```
.environments/development.memory
.environments/development.sqlite3
.environments/development.postgres
```
Within each is a `STORAGE_TYPE` value that is the switch to determine which repository is injected.  There are `make` targets to execute with a particular repository implementation.

### 📊 Data Representation Layers: DTOs, Domain Models, and DB Models
This project implements clear separation between data representations across different layers:

1. **DTOs (Data Transfer Objects)**
   - Define external API contract
   - Handle serialization/deserialization
   - Isolate API changes from internal implementations

2. **Domain Models**
   - Encapsulate business logic and rules
   - Independent of persistence and API concerns
   - Core objects used by application services

3. **DB Models**
   - Optimized for database persistence
   - Handle table/column mappings and relationships
   - Separate storage concerns from business logic

**Benefits:**
- API changes don't cascade through the system
- Database schema changes are isolated from business logic
- Each layer can evolve independently
- Clear boundaries enable focused testing
- Reduced coupling between external interfaces, business logic, and data persistence
- Easier to maintain and modify individual layers

### ✂️ Cross-Cutting Concerns
This project isolates infrastructure and operational concerns into dedicated packages:
```
internal/
├── config/     # Application settings & environment configuration
├── errors/     # Structured error types with appropriate HTTP mappings
├── logger/     # Structured logging with slog
└── middleware/ # Request/response processing pipeline
```

**Current Middleware:**
- Request logging

**Common Middleware Use Cases To Consider For Future:**
- Authentication/Authorization
- CORS handling
- Distributed tracing / OTEL
- Rate limiting
- Request validation

## 🛠️ Technology Choices
When building this project, I evaluated various libraries and frameworks. Here's why I chose these specific ones:

### Gin Web Framework
- Built-in middleware support
- Strong request validation
- Active community
- Production-ready
- Clean API design

### SQLite3
- Zero-config database
- Single file deployment
- ACID compliant
- Great for development
- SQL compatibility

### sqlx
- Type-safe SQL
- Simple but powerful
- No heavy ORM overhead
- Native SQL when needed
- Great connection pooling

### golang-migrate
- SQL migration management
- Support for multiple databases
- Version control for schemas
- Rollback capability
- CLI and library support

### slog
- Structured logging
- Built into stdlib
- High performance
- Flexible handlers
- Good error integration

### Testify
- Rich assertions
- Solid mocking support
- Test suite organization 
- Works with Go testing
- Clear error messages


## 📋 Project Overview

This service provides an API for scheduling appointments between trainers and clients.

### API Endpoints Required
- Get available appointment times for trainer between 2 dates
- Schedule new appointments
- View scheduled appointments for a trainer

### Business Requirements
- Appointments are 30 minutes long
- Scheduled at :00 or :30 past the hour
- Business hours: M-F 8am-5pm Pacific Time
- No overlapping appointments allowed

## 🛠 Build and Development

### Make Targets

```bash
# Build and Run
make build            # Build the service
make run              # Build and run the service
make clean            # Clean build artifacts

# Database
make migrate-up       # Apply pending migrations
make migrate-down     # Rollback migrations
make clean-db        # Remove database file
make rebuild-db      # Full database rebuild

# Development
make run-dev-mem      # Run with hot reload (memory storage)
make run-dev-postgres # Run with hot reload (postgres storage)
make fmt              # Format code
make vet              # Run Go vet
make lint             # Run linter
make tidy             # Tidy up modules

# Testing
make test                      # Run tests
make test-verbose              # Run tests with verbose output
make coverage-report-functions # Display test coverage at function granularity
make coverage-report-packages  # Display test coverage at package granularity
make coverage-report-html
make coverage-all     # Generate and show all coverage reports
make coverage-browser # Open coverage report in browser

# Tools
make install-tools    # Install required development tools
```

### Prerequisites
- Go 1.21 or higher
- Air (for hot reload during development)
- Migrate (for handling db migrations)

### Environment Setup
1. Install required tools:
   ```bash
   make install-tools
   ```
2. Start development server:
   ```bash
   # For memory storage:
   make run-dev-mem
   
   # For SQLite3:
    make migrate-up       # ensure the db is up to date
    make run-dev-sqlite3
   ```

## 🧪 Testing

The project includes some unit testing *BUT NEEDS MUCH MORE*:
- Unit tests for core components *<-- need more coverage*
- Repository implementation tests *<-- need more coverage*
- Coverage reporting
- Race condition detection

Generate coverage reports:
```bash
make coverage-all
```

## 📐 API Endpoints

### Get Available Appointments
```
GET /appointments/available
Parameters:
  - trainer_id: int
  - starts_at: string
  - ends_at: string
```

### Schedule Appointment
```
POST /appointments
Body:
{
  "trainer_id": int,
  "user_id": int,
  "starts_at": string,
  "ends_at": string
}
```

### List Trainer's Appointments
```
GET /appointments/trainer/{trainer_id}
```


## 🔧 Future Improvements

Potential areas for enhancement:
- Add production grade RDBMS repository
- Targets for working with local Postgres container
- Metrics and monitoring / OTEL
- Improve logging
- Ensure graceful shutdown is functioning
  - I scarmbled to get this in, but it is not critical to the submission, I hope
- Add metrics and monitoring
- Expand test coverage
- Add API documentation

## Proof of Functionality
There is a bash script with a handful of `curl`s that can be exercised.  That script is `scripts/run_scenario_1.sh`


I opened 2 terminal sessions
* Terminal 1
   * Clean, rebuild, tear down DB, mirgate up DB, run the server
     * all `make` targets
* Terminal 2
   * Run the testing script, review results

### Terminal 1: Run Server
```bash
$ make clean
rm -rf build
go clean -testcache
rm -f coverage.out
rm -f coverage.html
rm -f .env.tmp 

$ make build
go mod tidy
go fmt ./...
go vet ./...
Building the Go service...
go build -ldflags "\
        -X appointment-service/internal/version.Version=0.0.1 \
        -X appointment-service/internal/version.Commit=unknown \
        -X appointment-service/internal/version.BuildTime=2025-02-11T08:58:32Z" \
        -o build/appointment-service ./cmd/api
Build complete. Output: build/appointment-service

$ make migrate-down
echo "y" | migrate -database "sqlite3://data/appointments.db" -path migrations down
Are you sure you want to apply all down migrations? [y/N]
Applying all down migrations
1/d create_appointments_table (1.445958ms)

$ make migrate-up
migrate -database "sqlite3://data/appointments.db" -path migrations up
1/u create_appointments_table (1.6975ms)

$ make run-dev-sqlite3
go mod tidy
go fmt ./...
go vet ./...
Building the Go service...
go build -ldflags "\
        -X appointment-service/internal/version.Version=0.0.1 \
        -X appointment-service/internal/version.Commit=unknown \
        -X appointment-service/internal/version.BuildTime=2025-02-11T08:59:20Z" \
        -o build/appointment-service ./cmd/api
Build complete. Output: build/appointment-service
Loading environment from .environments/development.sqlite3:
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=debug
LOG_SOURCE=true
STORAGE_TYPE=sqlite3
DB_FILE=data/appointments.db
Starting the Go service in development (sqlite3 storage) mode with live reload...
env APP_ENV=development APP_PORT=8080 LOG_LEVEL=debug LOG_SOURCE=true STORAGE_TYPE=sqlite3 DB_FILE=data/appointments.db air -c .air.toml

  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.61.7, built with Go go1.23.5

[03:59:22] mkdir /Users/jbrown/personal/future/try-again/appointment-service/tmp
[03:59:22] watching .
[03:59:22] !exclude build
[03:59:22] watching cmd
[03:59:22] watching cmd/api
[03:59:22] watching data
[03:59:22] watching internal
[03:59:22] watching internal/api
[03:59:22] watching internal/app
[03:59:22] watching internal/config
[03:59:22] watching internal/dto
[03:59:22] watching internal/errors
[03:59:22] watching internal/logger
[03:59:22] watching internal/middleware
[03:59:22] watching internal/model
[03:59:22] watching internal/repository
[03:59:22] watching internal/repository/factory
[03:59:22] watching internal/repository/memory
[03:59:22] watching internal/repository/postgres
[03:59:22] watching internal/repository/sqlite3
[03:59:22] watching internal/service
[03:59:22] watching internal/service/factory
[03:59:22] watching internal/version
[03:59:22] watching migrations
[03:59:22] watching scripts
[03:59:22] !exclude tmp
[03:59:22] building...
go mod tidy
go fmt ./...
go vet ./...
Building the Go service...
go build -ldflags "\
        -X appointment-service/internal/version.Version=0.0.1 \
        -X appointment-service/internal/version.Commit=unknown \
        -X appointment-service/internal/version.BuildTime=2025-02-11T08:59:22Z" \
        -o build/appointment-service ./cmd/api
Build complete. Output: build/appointment-service
[03:59:24] running...
Version: 0.0.1
Commit: unknown
BuildTime: 2025-02-11T08:59:22Z
2025/02/11 03:59:24 Connected to SQLite DB at: data/appointments.db
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/appointments/trainers/:trainer_id --> appointment-service/internal/api.(*Server).ListAppointments-fm (3 handlers)
[GIN-debug] POST   /api/v1/appointments      --> appointment-service/internal/api.(*Server).CreateAppointment-fm (3 handlers)
[GIN-debug] GET    /api/v1/appointments/trainers/:trainer_id/availability --> appointment-service/internal/api.(*Server).GetAvailability-fm (3 handlers)
time=2025-02-11T03:59:24.690-05:00 level=INFO source=/Users/jbrown/personal/future/try-again/appointment-service/cmd/api/main.go:47 msg="Starting server" service=appointment-service version=0.0.1 commit_sha=unknown build_time=2025-02-11T08:59:22Z port=8080
```

### Terminal 2: Run Test Script
```bash
$ scripts/run_scenario_1.sh                                                                                                                                                INT ✘  03:59:40 AM  ▓▒░
============================================================================
 FUTURE TAKE HOME ASSIGNMENT
----------------------------------------------------------------------------
  * Demostrating the API service using curl


~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
** TEST CASE: List appointments for trainer 1
** EXPECTED: Empty list
curl -s -w \nStatus code: %{http_code}\n http://localhost:8080/api/v1/appointments/trainers/1
[]
Status code: 200

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
** TEST CASE: Get Availability of trainer 1 on June 1
** EXPECTED: List of all time slots on June 1 between 10AM and 2PM in UTC
curl -s -w \nStatus code: %{http_code}\n http://localhost:8080/api/v1/appointments/trainers/1/availability?starts_at=2025-06-01T18:00:00Z&ends_at=2025-06-01T22:00:00Z
[
  {
    "start_time": "2025-06-01T18:00:00Z",
    "end_time": "2025-06-01T18:30:00Z"
  },
  {
    "start_time": "2025-06-01T18:30:00Z",
    "end_time": "2025-06-01T19:00:00Z"
  },
  {
    "start_time": "2025-06-01T19:00:00Z",
    "end_time": "2025-06-01T19:30:00Z"
  },
  {
    "start_time": "2025-06-01T19:30:00Z",
    "end_time": "2025-06-01T20:00:00Z"
  },
  {
    "start_time": "2025-06-01T20:00:00Z",
    "end_time": "2025-06-01T20:30:00Z"
  },
  {
    "start_time": "2025-06-01T20:30:00Z",
    "end_time": "2025-06-01T21:00:00Z"
  },
  {
    "start_time": "2025-06-01T21:00:00Z",
    "end_time": "2025-06-01T21:30:00Z"
  },
  {
    "start_time": "2025-06-01T21:30:00Z",
    "end_time": "2025-06-01T22:00:00Z"
  }
]
Status code: 200



~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
** TEST CASE: Book appointment outside business hours
** EXPECTED: 4xx ERROR - outside business hours
curl -s -w \nStatus code: %{http_code}\n -X POST http://localhost:8080/api/v1/appointments -H Content-Type: application/json -d {
        "start_time": "2025-06-01T14:00:00Z",
        "end_time": "2025-06-01T14:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }
{
  "error": "appointment must start between 8am and 5pm Pacific"
}
Status code: 400



~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
** TEST CASE: Book valid appointment at 11AM
** EXPECTED: 200 OK with new appointment
curl -s -w \nStatus code: %{http_code}\n -X POST http://localhost:8080/api/v1/appointments -H Content-Type: application/json -d {
        "start_time": "2025-06-01T19:00:00Z",
        "end_time": "2025-06-01T19:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }
{
  "id": 1,
  "trainer_id": 1,
  "start_time": "2025-06-01T19:00:00Z",
  "end_time": "2025-06-01T19:30:00Z",
  "user_id": 10
}
Status code: 201



~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
** TEST CASE: Book conflicting appointment at 11AM
** EXPECTED: 4xx ERROR - time slot taken
curl -s -w \nStatus code: %{http_code}\n -X POST http://localhost:8080/api/v1/appointments -H Content-Type: application/json -d {
        "start_time": "2025-06-01T19:00:00Z",
        "end_time": "2025-06-01T19:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }
{
  "error": "trainer 1 is not available between 2025-06-01 19:00:00 +0000 UTC and 2025-06-01 19:30:00 +0000 UTC"
}
Status code: 409



~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
** TEST CASE: Book valid appointment for June 1 at 1PM
** EXPECTED: 200 OK with new appointment
curl -s -w \nStatus code: %{http_code}\n -X POST http://localhost:8080/api/v1/appointments -H Content-Type: application/json -d {
        "start_time": "2025-06-01T21:00:00Z",
        "end_time": "2025-06-01T21:30:00Z",
        "trainer_id": 1,
        "user_id": 12
    }
{
  "id": 2,
  "trainer_id": 1,
  "start_time": "2025-06-01T21:00:00Z",
  "end_time": "2025-06-01T21:30:00Z",
  "user_id": 12
}
Status code: 201



~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
** TEST CASE: Get Availability of trainer 1 on June 1
** EXPECTED: List of all time slots on June 1 between 10AM and 2PM in UTC, should not show 11AM and 1PM
curl -s -w \nStatus code: %{http_code}\n http://localhost:8080/api/v1/appointments/trainers/1/availability?starts_at=2025-06-01T18:00:00Z&ends_at=2025-06-01T22:00:00Z
[
  {
    "start_time": "2025-06-01T18:00:00Z",
    "end_time": "2025-06-01T18:30:00Z"
  },
  {
    "start_time": "2025-06-01T18:30:00Z",
    "end_time": "2025-06-01T19:00:00Z"
  },
  {
    "start_time": "2025-06-01T19:30:00Z",
    "end_time": "2025-06-01T20:00:00Z"
  },
  {
    "start_time": "2025-06-01T20:00:00Z",
    "end_time": "2025-06-01T20:30:00Z"
  },
  {
    "start_time": "2025-06-01T20:30:00Z",
    "end_time": "2025-06-01T21:00:00Z"
  },
  {
    "start_time": "2025-06-01T21:30:00Z",
    "end_time": "2025-06-01T22:00:00Z"
  }
]
Status code: 200



~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
** TEST CASE: List appointments for trainer 1
** EXPECTED: appointments for 11am and 1pm
curl -s -w \nStatus code: %{http_code}\n http://localhost:8080/api/v1/appointments/trainers/1
[
  {
    "id": 1,
    "trainer_id": 1,
    "start_time": "2025-06-01T19:00:00Z",
    "end_time": "2025-06-01T19:30:00Z",
    "user_id": 10
  },
  {
    "id": 2,
    "trainer_id": 1,
    "start_time": "2025-06-01T21:00:00Z",
    "end_time": "2025-06-01T21:30:00Z",
    "user_id": 12
  }
]
Status code: 200


End of API Tests
```
