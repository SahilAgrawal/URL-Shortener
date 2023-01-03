# URL-Shortener
I created a REST API using Golang and Gin, which can shorten the URL.

## Run the application in localhost:8080
```bash
go run ./main.go
```

## Endpoints
| Method  | Path        | Description               |
|---------|-------------|---------------------------|
| POST    | /url        | Shorten the url             |
| GET     | /url/{code} | Find shorten url by code  |
| GET     | /url        | List of shorten urls      |

## Written By
[Sahil Agrawal](https://github.com/SahilAgrawal)
