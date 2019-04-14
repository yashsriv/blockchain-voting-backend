# Blockchain Voting Backend

In order to run this, ensure your go version is above `1.11`:

```
$ go version
go version go1.12.3 linux/amd64
```

Also redis is needs to be running:
```
$ docker run --rm -p 6379:6379 -e ALLOW_EMPTY_PASSWORD=yes bitnami/redis 
```

Then:
```
$ go run ./cmd/server
```

## To run the command line client

```
$ go run ./cmd/contractor
```
