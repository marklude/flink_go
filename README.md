## FLINK GO

### Pre-requisit

[Install redis server](https://redis.io/download, "Redis")

Set environment variable `HISTORY_SERVER_LISTEN_ADDR` [default=8080]

Set environment variable `LOCATION_HISTORY_TTL_SECONDS` (Optional)

### Time To Live feature

I used redis to implement the `TTL` feature because there is no need for me to reinvent the wheel looking at the allocated time. Moreover, redis is a proven technology for caching which is tested in production by many big companies.
Redis is also modular and decoupled, hence can be easily scaled and maintained.

### Build and Run

> `go build .`

> `./flink_go`
