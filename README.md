# impulse

A proof-of-concept FAAS runtime, built on runc. Exposes an API for starting and monitoring "chambers", which are
lightweight containers that contain client functions and a platform function wrapper that binds to a port.

Users of impulse can make API calls to schedule functions, which are at runtime bound to a baseimage,
given a port allocation, started, and monitored. Users can list all chambers running and query
their status and allocated ports.

This is a component that would be used at the core of a larger FAAS platform.

## Future projects:
- An edge proxy for all API requests to this platform, responsible for directing requests to live chambers.
- A scheduler, responsible for making requests to impulse instances to start and stop containers in response
  to demand.
- A bakery service for creating images based on baseimages for client functions

## Roadmap
- Implement destroy endpoint
- Tests!
- Implement baseimages and support for languages beyond python 3.9
- Implement support for pulling images from a repository, not relying on them being present on disk at startup time.
- Deeper audit of required container capabilities

## Terminology
- Chamber: An instance of a function. The only chamber runtime currently is minimal runc containers.

## Related projects
Python baseimage: https://github.com/izaaklauer/baseimage-python3.9

Sample guest function: https://github.com/izaaklauer/guestimage-python3.9

## API

### Create a chamber

```
POST /chambers

{
    "app": "sampleapp",
    "version": "1.0.0",
    "runtime": "python3.9"
}
```

#### Response
```
Status: 201 Created
```


### List chambers

```
GET /chambers
```

#### Response
```
Status: 200

[
    {
        "id": "sampleapp-1.0.0-bkqek",
        "status", "STARTING",
        "createdTimeMillis": 1617809436251,
        "port": 6000,
        "spec": {
            "app": "sampleapp",
            "version": "1.0.0",
            "runtime": "python3.9"
        }
    },
    {
        "id": "sampleapp-1.0.0-k1hbd",
        "status", "RUNNING",
        "createdTimeMillis": 1617809436100,
        "port": 6001,
        "spec": {
            "app": "sampleapp",
            "version": "1.0.0",
            "runtime": "python3.9"
        }
    }
]

```

## Running the impulse server

`go run main.go`



