package memroute

import (
	"errors"
	"github.com/micro-go/lock"
	"sync"
)

type Router interface {
	Connect() (Client, error)
	Send(topic string, data interface{}) error
}

func NewRouter() Router {
	mutex := &sync.RWMutex{}
	subscriptions := newSubsciptions(mutex)
	clients := make(map[*client]interface{})
	return &router{mutex, subscriptions, clients}
}

const (
	separator = "/"
)

var (
	alreadyExistsErr = errors.New("Subscription already exists")
	badRequestErr    = errors.New("Bad request")
	mismatchErr      = errors.New("Mismatched subscriptions")
	noSubErr         = errors.New("No subscription")
)

type router struct {
	// Store a cache of every client that a topic results in.
	// The primary use case currently favours performance over memory,
	// but this is also an easy way to prevent a client from receiving duplicates.
	// NOTE: Only access clients with the subscription lock set.
	mutex         *sync.RWMutex
	subscriptions subscriptions
	clients       map[*client]interface{}
}

func (r *router) Connect() (Client, error) {
	// Do the simple thing of forcing all routes to be recached
	// whenever a subscription changes. This is a good match for
	// the current use case of creating all subscriptions on startup,
	// but we would want to optimize if adding and removing subscriptions
	// happened frequently.
	defer lock.Write(r.mutex).Unlock()
	r.subscriptions.lockedClear()
	c := newClient(r)
	r.clients[c] = nil
	return c, nil
}

func (r *router) disconnect(c *client) error {
	if c == nil {
		return badRequestErr
	}
	defer lock.Write(r.mutex).Unlock()
	r.subscriptions.lockedClear()
	close(c.receiver)
	delete(r.clients, c)
	return nil
}

func (r *router) Send(topic string, payload interface{}) error {
	return r.sendFrom(nil, topic, payload)
}

func (r *router) sendFrom(fromC *client, topic string, payload interface{}) error {
	sub := r.subscriptions.find(topic)
	if sub == nil {
		sub = r.subscriptions.make(topic, r)
	}
	if sub != nil && len(sub.clients) > 0 {
		m := &Message{topic, payload}
		for toC, _ := range sub.clients {
			if fromC != toC {
				toC.Receiver() <- m
			}
		}
	}
	return nil
}
