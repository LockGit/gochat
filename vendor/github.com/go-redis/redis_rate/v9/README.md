# Rate limiting for go-redis

[![Build Status](https://travis-ci.org/go-redis/redis_rate.svg?branch=master)](https://travis-ci.org/go-redis/redis_rate)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-redis/redis/v8)](https://pkg.go.dev/github.com/go-redis/redis_rate/v9)

> :heart: [**Uptrace.dev** - distributed traces, logs, and errors in one place](https://uptrace.dev)

This package is based on [rwz/redis-gcra](https://github.com/rwz/redis-gcra) and implements
[GCRA](https://en.wikipedia.org/wiki/Generic_cell_rate_algorithm) (aka leaky bucket) for rate
limiting based on Redis. The code requires Redis version 3.2 or newer since it relies on
[replicate_commands](https://redis.io/commands/eval#replicating-commands-instead-of-scripts)
feature.

## Installation

redis_rate supports 2 last Go versions and requires a Go version with
[modules](https://github.com/golang/go/wiki/Modules) support. So make sure to initialize a Go
module:

```shell
go mod init github.com/my/repo
```

And then install redis\_rate/v9 (note **_v9_** in the import; omitting it is a popular mistake):

```shell
go get github.com/go-redis/redis_rate/v9
```

## Example

```go
package redis_rate_test

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

func ExampleNewLimiter() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	limiter := redis_rate.NewLimiter(rdb)
	res, err := limiter.Allow(ctx, "project:123", redis_rate.PerSecond(10))
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
	// Output: allowed 1 remaining 9
}
```
