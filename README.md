# How to run

1. Copy `example-config.yaml` to a new file called `config.yaml`.
2. The default port is 5000, but can be changed if you like.
3. Set up a postgreSQL database and add its url to the config file.
4. Migration files are by default in the `migrations` folder but can be changed in the config.
4. Run `go run *.go`


# Endpoints

### These endpoints return info about services that a user has installed. Each service has a type code, name, description, and the versions of the service the user has currently installed. If the user has not installed any versions of a specific service type, it will not be included in the results.

### Included service types:
- Code 1: Database service
- Code 2: Reporting service
- Code 3: Currency conversion service
- Code 4: Translation service
- Code 5: Notifications service

## `GET /services`: Returns a list of all the services installed.
If no query parameters are provided, all installed services are returned.

Query parameters
- `search` will search for services of a type matching the search string.
- `sort` will sort the result. The options are ascending `name`, `description`, or `typecode`. Default is `typecode`.
- `pagesize` and `page` are pagination options for the returned results. They will group the results into groups of size `pagesize` and return the specific group based on `page`. Default `pagesize` is all results and default `page` is 1. Any values less than 1 will use the defaults.

## `GET /services/:type_code`: Returns info about a specific installed service.
## `GET /services/:type_code/versions`: Returns info about the currently in use (installed) versions of a specific service.
## `GET /services/:type_code/versions/:version_number`: Returns info about a specific in use (installed) version of a service.

# Design Decisions

My general approach was to keep the project pretty simple.
- The whole project is one package. Given more time, I would have split the API layer, the service layer, and the repository layer into separate packages.
- For errors and logging, I stuck with Go's built in error handling and logging.
- There is no concept of a user or user accounts to keep things simple. The assumption is that anyone using the api is part of the same account.
- I used the Gin web framework to handle http requests because I am familiar with it.
- For the data layer, I used postgreSQL since a relation database fits the project well.


# Config

While configuration values are often set as environment variables in production settings, reading from a YAML file is quick and simple for local development.

Validation of the config is done to prevent deploying the server with a missing configuration value. A user of an API usually doesn't have access to the server to fix missing configuration values, so we prevent the server running in the first place if a config is missing.

# API layer

Gin was chosen as the http framework since I am familiar with it and using it will speed things up.

I used pretty traditional http status codes when returning responses to the user. The errors returned to the user are basic, but given more time, they could be expanded to provide a bit more context info to the user.

# Service (business logic) layer

This layer holds all of the business logic and deals with 3 different concepts.

## Service types
There are the different types of services a user can install. I made the assumption that these types cannot be changed by the API user and are not changed very often by the developers. I defined these different types as unsigned integer constants of the type `ServiceTypeCode`. For showing the different types to the end user, the types are turned into structs of the type `ServiceType`. These structs hold other human readable info about the type.

With these service types defined in memory instead of a database, they are very fast to fetch and search. The changes to the types are also very easy to track with version control software.

The obvious downside to not storing them in a database is that new types cannot be created by the API user. New code must be pushed to get a new service type.

## Installed service versions
These are the versions of a service that a user currently has installed. These records are stored in the database since a user can deploy different versions of the service and we want to keep track of those new deployments.

## Service
This data type is for representing the API data in a pretty format to the user. It combines the information about a service's type along with the versions of that service a user currently has installed.

# Repository (data) layer

The repository layer is represented by an interface. I did this to make it easy to swap out the technology used to do persistence. During development, I used a mock persistence layer that worked in memory. Once I had business logic figured out, I swapped this mock persistence layer for a postgres implementation.

# Testing (for the future)

Given more time, I would write a good set of unit tests for the business layer. It could use a good set of unit tests especially for the code that handles the searching, ordering, and pagination of the API results.

The API layer could also be a good place to write a set of unit or integration tests. The API layer is essentially the user interface so we can't assume that the user will do only exactly what we want them to. It would be good to test for a wide variety of user inputs to make sure we are properly handling the inputs or we are providing useful errors to the user.

I would also want to build a set of integration tests that test the interaction between the layers. This project is basic enough that an integration test would not be too cumbersome to set up.