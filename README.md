# livetest
#### livetest is a continous testing service/tool for testing backend endpoints/features in development or production environments. Periodically request a collection of endpoints and verify the response with the option to store results in a database and send a notification when a request fails

# Get Started

## Binary Installation
```
$ go get -u github.com/rknizzle/livetest/cmd/livetest
```
NOTE: Must include a config file to run with the binary


## Test locally with Docker Compose
```
docker-compose up
```

# Config

### concurrency
#### number of requests that can run at the same time

### datastore
#### Store the result of each request to track responses over time  
Supported datastores:
- postgres

### notification
#### Send out a notification if a request fails or returns a bad response  
Supported notifications:
- HTTP request

### Example config:
```
{
    "concurrency": 2,
    "datastore": {
        "db": "postgres",
        "dbname": "postgres",
        "host": "localhost",
        "password": "password",
        "port": 5432,
        "user": "postgres"
    },
    "jobs": [
        {
            "expectedResponse": {
                "statusCode": 200
            },
            "frequency": 5000,
            "headers": {},
            "httpMethod": "GET",
            "requestBody": {},
            "title": "example GET",
            "url": "http://postman-echo.com/get?foo1=bar1&foo2=bar2"
        },
        {
            "expectedResponse": {
                "statusCode": 200
            },
            "frequency": 8000,
            "headers": {
                "Content-Type": "application/json"
            },
            "httpMethod": "POST",
            "requestBody": {
                "data": "value"
            },
            "title": "example POST",
            "url": "http://postman-echo.com/post"
        }
    ],
    "notification": {
        "msg": {
            "headers": {
                "Content-Type": "application/json"
            },
            "httpMethod": "POST",
            "requestBody": {
                "notification": "true"
            },
            "url": "http://postman-echo.com/post"
        },
        "type": "http"
    }
}
```
