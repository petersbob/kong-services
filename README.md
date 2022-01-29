## How to run

1. Copy `example-config.yaml` to a new file called `config.yaml`
2. The default port is 5000, but can be changed if you like
3. run `go run main.go`


# Design Decisions

Kept the whole project in one package to keep things simple. I may change this as it gets more complicated.

# Config

I used yaml since I am familiar with it.

Validation of the config was done to prevent deploying the server with a missing configuration value. A user of an API usually doens't have access to server to set configuration values, so we prevent the server running in the first place if a config is missing. Also, from personal experience, errors related to missing config values can be really frustrating to solve.

# API layer

Gin was chosen as the http framework since I am familiar with it and using it will speed things up.

Did not let the server start if configs do not exist. The user of the server's api can't do anything about fixing missing configs, so we prevent the server from running in the the first place.