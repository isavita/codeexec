# Code Execution API

This is a Go-based API that allows you to execute code in various programming languages, including Python and JavaScript. The API provides a simple endpoint for code submission and returns the output of the executed code.

## Features

- Execute code in multiple programming languages
- Containerized execution environment for security and isolation
- API authentication using an API key
- Configurable port for running the API server
- Detailed error messages and execution output

## Prerequisites

- Docker: Make sure you have Docker installed on your machine. You can download it from the official Docker website: https://www.docker.com/get-started

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/yourusername/code-execution-api.git
```

2. Build the Docker images:

```bash
./build.sh
```

3. Run the API server:

```bash
docker run -p 8080:8080 codeexec
```

The API server will start running on `http://localhost:8080`.

## Configuration

### Port

By default, the API server runs on port `8080`. You can change the port by setting the `PORT` environment variable when running the container:

```bash
docker run -p 8000:8000 -e PORT=8000 codeexec
```

### API Authentication

The API uses an API key for authentication. When making requests to the API endpoint, include the `X-Api-Key` header with your API key.

To enable API authentication, set the `API_KEY_CHECK_ENABLED` environment variable to `"true"` and provide your API key using the `API_KEY` environment variable when running the container:

```bash
docker run -p 8080:8080 -e API_KEY_CHECK_ENABLED=true -e API_KEY=your-api-key codeexec
```

To disable API authentication, set the `API_KEY_CHECK_ENABLED` environment variable to `"false"` or omit it entirely.

## API Endpoint

### Execute Code

- URL: `/api/execute`
- Method: `POST`
- Headers:
  - `Content-Type: application/json`
  - `X-Api-Key: your-api-key` (if API authentication is enabled)
- Request Body:
  ```json
  {
    "code": "your-code-here",
    "language": "python"
  }
  ```

## Examples

### Execute Python Code

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "code": "print(\"Hello, World!\")",
  "language": "python"
}' http://localhost:8080/api/execute
```

Response:
```json
{
  "output": "Hello, World!"
}
```

### Execute JavaScript Code

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "code": "console.log(\"Hello, World!\");",
  "language": "javascript"
}' http://localhost:8080/api/execute
```

Response:
```json
{
  "output": "Hello, World!"
}
```

### Authentication Error

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "code": "print(\"Hello, World!\")",
  "language": "python"
}' http://localhost:8080/api/execute
```

Response:
```json
{
  "error": "unauthorized"
}
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
