# Survey App

## Requirements

Refer [here](./requirements.pdf)

## Prerequisites
1. [Go 1.16](https://golang.org/dl/)

## How To Run

From project root directory run:

```sh
$ go build -o survey-platform cmd/main.go
$ ./survey-platform
```

Alternatively, you can run the application after running tests with a single command 
if [GNU Make](https://www.gnu.org/software/make/) is installed
```sh
$ make all
```

## API docs
The API docs for this application can be found [here](./postman_collection.json) 
as a postman collection.
The environment required for the collection is added [here](./postman_environment.json)

The default value for `APP_PORT` is 8000.
It can be overridden by setting environment variable `APP_PORT` to the required port.

## How to Test
From project root directory run:
```sh
$ go test -cover -race ./...
```

## Running with Docker & docker-compose

### Prerequisites
1. [docker](https://docs.docker.com/engine/install/)
2. [docker-compose](https://docs.docker.com/compose/install/)

### Instructions

Run
```sh
$ docker-compose up
```

This will build docker and run application in a docker container.
Port mapping is done from `8000:8000`
In case need to change the port on host, change the first argument to the required port, like `9000:8000`


### Swagger Docs
 ``` 
 swag init -g cmd/main.go
 ```
Run server and go to http://localhost:8080/swagger/index.html from browser

### NOTE
No particular database is used for this application. 
On application exit, the current data is dumped to a json file 
to persist the data so that it will be loaded in the subsequent app run.

As an extension, if a db needs to be added it can be added easily into the repo layer with minimal changes as another implementation

While running in docker mode, volume mapping needs to considered for persisting the file to host machine