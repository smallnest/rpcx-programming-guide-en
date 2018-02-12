# Circuit Breaker

Circuit Breaker is use to prevent a network or service failure from cascading to other services, see the post: [Pattern: Circuit Breaker](http://microservices.io/patterns/reliability/circuit-breaker.html).

A service client should invoke a remote service via a proxy that functions in a similar fashion to an electrical circuit breaker. When the number of consecutive failures crosses a threshold, the circuit breaker trips, and for the duration of a timeout period all attempts to invoke the remote service will fail immediately. After the timeout expires the circuit breaker allows a limited number of test requests to pass through. If those requests succeed the circuit breaker resumes normal operation. Otherwise, if there is a failure the timeout period begins again.


Rpcx defines a `Breaker` interface and you can implement your customized circuit creaker.

```go
type Breaker interface {
	Call(func() error, time.Duration) error
}
```

Rpcx implements a simple `ConsecCircuitBreaker` that is a window sliding CircuitBreaker with failure threshold.

You can set your CircuitBreaker instance to `Option.Breaker`.