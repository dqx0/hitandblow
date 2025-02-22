# Hit and Blow API

The Hit and Blow API is a server application that provides a RESTful API for playing Hit and Blow games.  
The project is implemented in the Go language and Gin framework and uses Swagger (swaggo/swag) to automatically generate API documentation.  
A CI/CD pipeline using Docker and Kubernetes (GKE) is also in place.

## Features

- **Swagger Documentation**: You can find the API documentation at `/swagger/index.html`.
- **CI/CD pipeline**: GitHub Actions automatically build, push, and deploy Docker images to GKE.


## Directory structure
```
hitandblow/
├── cmd/
│ └── main.go
├── internal/ 
│ ├── server/ 
│ │ └── server.go
│ └── game/
│   └── game.go
├── docs/
├── .github/ 
│ └── workflows/ 
│   └── workflow.yaml
├── Dockerfile
├── go.mod
└── README.md
```
