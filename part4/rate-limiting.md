# rate-limiting

This plugin uses [juju/ratelimit](https://github.com/juju/ratelimit) to limit requests.

Use `func NewRateLimitingPlugin(fillInterval time.Duration, capacity int64) *RateLimitingPlugin` to create a new plugin.

It contains a new token bucket that fills at the rate of one token every fillInterval, up to the given maximum capacity. Both arguments must be positive. The bucket is initially full.