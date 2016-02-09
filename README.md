
# ratelimiter

A simple rate limiter.  Warning: it's not perfectly accurate — you will still have to handle overages.  With that said, it works pretty well.

```go
import (
    "github.com/brynbellomy/go-ratelimiter"
)

func main() {
    rl := ratelimiter.New(1000, 1 * time.Second)

    for {
        doWork(rl)
    }
}

func doWork(rl *ratelimiter.RateLimiter) {
    rl.GetCapacity(1) // blocks
    defer rl.ReleaseCapacity(1)

    // ... do some work ...
}
```