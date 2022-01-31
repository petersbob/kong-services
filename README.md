# How to run

1. Copy `example-config.yaml` to a new file called `config.yaml`.
2. The default port is 5000, but can be changed if you like.
3. Set up a postgreSQL database and add its url to the config file.
4. Migration files are by default in the `migrations` folder but can be changed in the config.
4. Run `go run *.go`


# Endpoints

### These endpoints return info about services that a user has installed. Each service has a type code, name, description, and the versions of the service the user has currently installed. If the user has not installed any versions of a specific service type, it will not be included in the results.

### Invluded service types:
- Code 1: Database service
- Code 2: Reporting service
- Code 3: Currency conversion service
- Code 4: Translation service
- Code 5: Notifications service

## `GET /services`: Returns a list of all the services installed.
If no query parameter are provided, all installed services are returned.

Query parameters
- `search` will search for services of a type matching the search string.
- `sort` will sort the result. The options are ascending `name`, `description`, or `typecode`. Default is `typecode`.
- `pagesize` and `page` are pagination options for the returned results. They will group the results into groups of size `pagesize` and return the specific group based on `page`. Default `pagesize` is all results and default `page` is 1. Any values less than 1 will use the defaults.

## `GET /services/:type_code`: Returns info about a specific installed service.
## `GET /services/:type_code/versions`: Returns info about the currently in use (installed) versions of a specific service.
## `GET /services/:type_code/versions/:version_number`: Returns info about a specific in use (installed) version of a service.

# Design Decisions

Kept the whole project in one package to keep things simple. I may change this as it gets more complicated.

There is no concept of a user or user accounts to keep things simple. The assumption is that anyone using the api is part of the same account.

# Config

I used yaml since I am familiar with it.

Validation of the config was done to prevent deploying the server with a missing configuration value. A user of an API usually doens't have access to server to set configuration values, so we prevent the server running in the first place if a config is missing. Also, from personal experience, errors related to missing config values can be really frustrating to solve.

# API layer

Gin was chosen as the http framework since I am familiar with it and using it will speed things up.

Did not let the server start if configs do not exist. The user of the server's api can't do anything about fixing missing configs, so we prevent the server from running in the the first place.

# Service (business logic) layer

I set up the service types as constants in code. The benefit is that any logic related to different service types will be really fast and changes to the service types can be tracked via version control. The downside of this is that it does not allow the end user of the product to create new service types since code access would be needed to create a new service type. I made the assumption that new service types are added very rarely and done by developers.

# Repository (data) layer

The repository layer is mean to interface with the database. The repository layer is abstracted out to be an interface to allow changing of the underlying implementation of the persistence.