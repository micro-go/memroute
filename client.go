package memroute

// interface Client manages a single connection to the router.
// You can create a single client that has multiple subscriptions
// if you want to deal with resolving the route, or you could
// create a client per subscription if you want to listen to
// multiple channels. Note that subscribing and unsubscribing
// on a client is not thread safe.
type Client interface {
	Disconnect()
	Subscribe(topic string) error
	Unsubscribe(topic string)
	Receiver() chan *Message
	Send(topic string, data interface{}) error
}

func newClient(r *router) *client {
	topics := make(map[string]interface{})
	receiver := make(chan *Message, 100)
	return &client{r, topics, receiver}
}

type client struct {
	r        *router
	topics   map[string]interface{}
	receiver chan *Message
}

func (c *client) Disconnect() {
	c.r.disconnect(c)
}

func (c *client) Subscribe(topic string) error {
	// See router.Connect() for an explanation.
	c.r.subscriptions.clear()
	c.topics[topic] = nil
	return nil
}

func (c *client) Unsubscribe(topic string) {
	// See router.Connect() for an explanation.
	c.r.subscriptions.clear()
	delete(c.topics, topic)
}

func (c *client) Receiver() chan *Message {
	return c.receiver
}

func (c *client) Send(topic string, data interface{}) error {
	return c.r.sendFrom(c, topic, data)
}
