# Tickers API

[![Build](https://github.com/falmar/sply-richerage/actions/workflows/build.yaml/badge.svg)](https://github.com/falmar/sply-richerage/actions/workflows/build.yaml)
[![Publish](https://github.com/falmar/sply-richerage/actions/workflows/publish.yaml/badge.svg)](https://github.com/falmar/sply-richerage/actions/workflows/publish.yaml)

> "Why did they call it 'richerage' instead of 'brokerage'?
> Because after you start investing, you're no longer broke, you're getting rich!"
>
> — AI

## Run Docker

```bash
$ docker run --rm -it -p 8080:80 docker.io/falmar/sply:http -p 80 -d
```

## API

### POST /login
```
POST /login HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

```json
{
  "username": "anonymous",
  "password": "anonymous"
}
```

To retrieve the auth token, you must first login with ANY username and password. The auth token is valid for 1 week.

```bash
$ curl -X POST -H "Host: localhost:8080" -H "Content-Type: application/json" -d '{"username": "anonymous", "password": "anonymous"}' http://localhost:8080/login
```

Long lived test token:
> `DbjHLiUBzrLBrYKqnC8HOHgnNhOGgl+iZKakolvRgEM=.YW5vdGhlcg==.1708069778`


### GET /tickers
```
GET /tickers HTTP/1.1
Host: localhost:8080
Authorization: Basic xxx
```

Will return a list of all tickers for a given user.

```bash
$ curl -X GET -H "Host: localhost:8080" -H "Authorization: Basic xxx" http://localhost:8080/tickers 
```

### GET /tickers/{ticker}/history
```
GET /tickers/AAPL/history HTTP/1.1
Host: localhost:8080
Authorization: Basic xxx
```

Will return a list of all historical prices for a given ticker.

```bash
$ curl -X GET -H "Host: localhost:8080" -H "Authorization: Basic xxx" http://localhost:8080/tickers/AAPL/history
```


## Summary

#### dependencies:
- For the project structure used https://github.com/golang-standards/project-layout
- Used [go-kit/kit](https://gokit.io/) for the writing the tickers/login services it lays out a good foundation for writing microservices in go. with the use of 3 layers (transport, endpoint, service) it makes it easy to add new transports (grpc, http, etc).
- Used [chi](https://github.com/go-chi/chi) as the http router, it's a lightweight router that has a lot of features.
- Used [Zap Logger](https://github.com/uber-go/zap) for logging, almost no logging is done but it's there.
- Used [Cobra](https://github.com/spf13/cobra)/[Viper](https://github.com/spf13/viper) for CLI/Configuration

#### code structure:
- The main logic for tickers is in `./internal/tickers`
- The main logic for login is in `./internal/login`
- Additional helper/shared code is in `./internal/pkg`
- The cli entrypoint is in `./cmd/main.go`
- Http command is in `./cmd/http/http.go`


#### data

About the code in `./internal/storage` the "seeded storage" its a deterministic `*rand.Rand` that uses the string hashes of the usernames/ticker/symbol to generate the data on the fly and provide consistent data across server restarts.


The tickers are variable for each user given its username is the seed. so users may see "different" *amount* of tickers but the data of the ticker itself is the same across all users.  

The time series of the tickers is seeded by the symbol, so it is also consistent across server restarts, the *amount* of series and the dates+price are also deterministic but the data is different for each symbol independently of the user. 

This storage is an interface like everything else in the codebase to allow easy plug-and-play of different storage backends. May it be database, object storage, filesystem etc. 

I do apologize in advance if I misunderstood the data generation, once I receive confirmation about it, it will be updated, given it is just an interface and simply plug-in the new storage implementation, it should have little to none impact on the rest of the codebase
