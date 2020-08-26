# Suru - slack bot 

## Description

This slack bot is designed for saving question and multiple answers to an Elasticsearch datastore. 

This service is used for interfacing with the Slack API with an api built using Golang, and an Elasticsearch datastore.

## Local Development
#### Running Go API from project location
Copy .env.example and configure ENV variables 

```
$ cp .env.example .env
```

Compile and run the go api
```
go run cmd/chatservice.go
```
#### Setup slack bot
Go to [https://api.slack.com](https://api.slack.com) for instructions on setting up your own slack app.

#### Expose local server to the public
Use [https://ngrok.com](https://ngrok.com) for quick setup to connect your local go api to the slack api.

#### Spin up ElasticSearch datastore
To spin up just the Elasticsearch docker container
```
docker-compose up es -d
```

#### Run API & ES in containers
```
docker-compose up -d
```

## Tests
### A health point test has been implemented.
```
go test ./...
```