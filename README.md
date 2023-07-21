# Tickers API

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
