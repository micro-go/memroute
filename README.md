# memroute
Go library that performs in-memory message routing.

## Description

This is a simple layer on top of channels to provide string-based topic routing and multiple destinations.

## Example

Construct a router:

```
import "github.com/micro-go/memroute"

r := memroute.NewRouter()
```

Create clients:

```c1 := r.Connect()```

Subscribe to topics:

```c1.Subscribe("path/to/something")```

Listen to the subscriptions:

```
for {
	select {
	...
	case m, more := <-c1.Receiver():
		if more {
			fmt.Println("received", m.Payload, "on topic", m.Topic)
		}
	}
}
```

And send messages, either via another client:

```c2.Send("path/to/something", payload)```

Or directly from the router:

```r.Send("path/to/something", payload)```

For greater convenience, a singleton is available if you only need a single global router:

```
import "github.com/micro-go/memroute"

c1 := memroute.Connect()

memroute.Send("path/to/something", payload)
```

## Topics

Topics follow the MQTT spec: Pattern matching where `/` is a separator, `+` is a wild card match on a single level, and `#` is a wild card match for multiple levels. For example, a client subscribed to `a/#` will receive topics `a/b`, `a/c`, `a/b/c` etc. A client subscribed to `a/+/c` will receive `a/b/c` but not `a/b/b`.

## Payloads

`payload` is an interface{} and can be anything; it's up to your application to define a data model.