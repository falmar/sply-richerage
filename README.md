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



> Test Token: `DbjHLiUBzrLBrYKqnC8HOHgnNhOGgl+iZKakolvRgEM=.YW5vdGhlcg==.1708069778`



