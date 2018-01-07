package memroute

import (
	"github.com/micro-go/lock"
	"github.com/micro-go/parse"
	"sync"
)

// --------------------
// SUBSCRIPTIONS
// --------------------

type subscriptions struct {
	mutex *sync.RWMutex
	all   map[string]*subscription
}

func newSubsciptions(mutex *sync.RWMutex) subscriptions {
	all := make(map[string]*subscription)
	return subscriptions{mutex, all}
}

// func clear() removes all subscriptions.
func (s *subscriptions) clear() {
	defer lock.Write(s.mutex).Unlock()
	s.lockedClear()
}

// func clear() removes all subscriptions. Must already be locked.
func (s *subscriptions) lockedClear() {
	s.all = make(map[string]*subscription)
}

// func find() finds the current subscription for a given
// topic. Answer nil if it hasn't been constructed.
func (s *subscriptions) find(topic string) *subscription {
	defer lock.Read(s.mutex).Unlock()
	return s.all[topic]
}

// func make() creates a cache of every client the route
// resolves to.
func (s *subscriptions) make(topic string, r *router) *subscription {
	defer lock.Write(s.mutex).Unlock()
	found := s.all[topic]
	if found != nil {
		return found
	}
	sub := newSubscription()
	s.all[topic] = sub
	for c, _ := range r.clients {
		if routeMatchesAny(topic, c.topics) {
			sub.clients[c] = nil
		}
	}
	return sub
}

func routeMatchesAny(topic string, patterns map[string]interface{}) bool {
	if topic == "" || patterns == nil {
		return false
	}
	for pattern, _ := range patterns {
		m := parse.NewMqttStringMatch(pattern)
		if m.Matches(topic) {
			return true
		}
	}
	return false
}

// --------------------
// SUBSCRIPTION
// --------------------

// struct subscription provides all clients in a subscription.
type subscription struct {
	clients map[*client]interface{}
}

func newSubscription() *subscription {
	clients := make(map[*client]interface{})
	return &subscription{clients}
}
