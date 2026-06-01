Based on your project structure, here is the Go project layout written in a clean documentation style:

```text
project-root/
│
├── config/
│   └── config.go          # Application configuration (env, settings)
│
├── database/
│   └── database.go        # Database connection and initialization
│
├── handlers/
│   ├── auth.go            # Authentication handlers
│   ├── handler.go         # Common/base handlers
│   ├── upload.go          # File upload handlers
│   └── user.go            # User-related handlers
│
├── jobs/
│   ├── cleanup.go         # Cleanup cron jobs
│   └── email.go           # Email cron jobs
│
├── logger/
│   └── logger.go          # Logger configuration and setup
│
├── logs/
│   ├── demo-rest-api.log  # Application log file
│   └── wtf.json           # Log output / debug file
│
├── middlewares/
│   ├── auth.go            # Authentication middleware
│   └── logger.go          # Request logging middleware
│
├── models/
│   └── user.go            # Database models / structs
│
├── routes/
│   └── routes.go          # Route registration
│
├── scheduler/
│   └── scheduler.go       # Cron scheduler setup and job registration
│
├── uploads/               # Uploaded files storage
│
├── utils/
│   ├── functions.go       # Helper functions
│   ├── jwt.go             # JWT utilities
│   └── password.go        # Password hashing / verification
│
├── .env                   # Environment variables
├── .env.example           # Example environment config
├── .gitignore             # Git ignored files
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── gorm_raw_ref.md        # GORM raw query reference
└── main.go                # Application entry point
```

### Layer Responsibility

A typical request flow in this structure is:

```text
main.go
   ↓
config → logger → database
   ↓
routes
   ↓
middlewares
   ↓
handlers
   ↓
models / utils
   ↓
database
```

Background jobs flow:

```text
main.go
   ↓
scheduler
   ↓
jobs
   ↓
database / utils / logger
```

This structure follows a fairly standard **REST API + clean modular architecture** in Go:

* **config** → app settings
* **database** → DB setup
* **models** → schema/entities
* **handlers** → HTTP business logic
* **middlewares** → request pipeline
* **routes** → endpoint mapping
* **utils** → reusable helpers
* **scheduler + jobs** → cron/background processing
* **logger/logs** → logging system
* **uploads** → file storage
* **main.go** → bootstrap and graceful shutdown orchestration.
