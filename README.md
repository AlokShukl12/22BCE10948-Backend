# File Sharing Platform

A Go-based file sharing platform that allows users to upload, manage, and share files. The system uses PostgreSQL for metadata storage, Redis for caching, and AWS S3 for file storage.

## Features

- User authentication with JWT
- File upload and management
- File sharing with expiration dates
- File search functionality
- Redis caching for file metadata
- Background job for cleaning up expired files
- Rate limiting
- Docker support

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- Redis
- AWS S3 bucket (or local storage)
- Docker (optional)

## Configuration

1. Copy the `.env.example` file to `.env` and update the values:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=filesharing

REDIS_HOST=localhost
REDIS_PORT=6379

AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_BUCKET_NAME=your-bucket-name

JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=24h

SERVER_PORT=8080
```

2. Create the PostgreSQL database:
```sql
CREATE DATABASE filesharing;
```

## Running the Application

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run main.go
```

### Using Docker

1. Build the Docker image:
```bash
docker build -t filesharing .
```

2. Run the container:
```bash
docker run -p 8080:8080 filesharing
```

## API Endpoints

### Authentication
- `POST /register` - Register a new user
- `POST /login` - Login and get JWT token

### Files
- `POST /upload` - Upload a new file
- `GET /files` - List all files
- `GET /files/search` - Search files
- `GET /files/:id` - Get file details
- `POST /files/:id/share` - Share a file

## Testing

Run the tests:
```bash
go test ./...
```

## Project Structure

```
.
├── auth/           # Authentication related code
├── database/       # Database initialization
├── handlers/       # HTTP request handlers
├── middleware/     # HTTP middleware
├── models/         # Data models
├── repositories/   # Database repositories
├── routes/         # Route configuration
├── services/       # Business logic
├── storage/        # Storage related code
├── main.go         # Application entry point
├── go.mod          # Go module file
├── Dockerfile      # Docker configuration
└── README.md       # Project documentation
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 