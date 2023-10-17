package main

import (
	"log"

	"github.com/ffiat/nostr"
)

type Repository struct {
	db  map[string]*nostr.Event
	cfg *Config
}

func (s *Repository) Store(e *nostr.Event) error {
	s.db[e.Id] = e
	return nil
}

func (s *Repository) All() []*nostr.Event {

	var events []*nostr.Event
	for _, v := range s.db {
		events = append(events, v)
	}

	return events
}

// TODO: Cache the pulled events.
func (s *Repository) FindByPubKey(key string) []*nostr.Event {

	log.Printf("Finding by Pubjey with id: %s\n", key)

	var events []*nostr.Event

	pk, err := nostr.DecodeBech32(key)
	if err != nil {
		log.Fatalf("\nunable to decode npub: %#v", err)
	}

	// List only the latest 3 event from the author.
	f := nostr.Filter{
		Authors: []string{pk.(string)},
		Kinds:   []uint32{nostr.KindTextNote},
		Limit:   10,
	}

	for _, v := range s.cfg.Relays {

		cc := NewConnection(v)
		err := cc.Listen()
		if err != nil {
			log.Fatalf("unable to listen to relay: %v", err)
		}

		sub, err := cc.Subscribe(nostr.Filters{f})
		if err != nil {
			log.Fatalf("\nunable to subscribe: %#v", err)
		}

		orDone := func(done <-chan struct{}, c <-chan *nostr.Event) <-chan *nostr.Event {
			valStream := make(chan *nostr.Event)
			go func() {
				defer close(valStream)
				for {
					select {
					case <-done:
						return
					case v, ok := <-c:
						if ok == false {
							return
						}
						valStream <- v
					}
				}
			}()
			return valStream
		}

		for e := range orDone(sub.Done, sub.EventStream) {
			events = append(events, e)
		}

		//cc.Close()
	}

	log.Println("event added to local cache")
	log.Println(events)

	return events
}

// TODO: Cache the pulled events.
func (s *Repository) FindByEventId(id string) []*nostr.Event {

	log.Printf("Finding event with id: %s\n", id)

	var events []*nostr.Event

	f := nostr.Filter{
		Ids:   []string{id},
		Kinds: []uint32{nostr.KindTextNote},
		Limit: 10,
	}

	for _, v := range s.cfg.Relays {

		cc := NewConnection(v)
		err := cc.Listen()
		if err != nil {
			log.Fatalf("unable to listen to relay: %v", err)
		}

		sub, err := cc.Subscribe(nostr.Filters{f})
		if err != nil {
			log.Fatalf("\nunable to subscribe: %#v", err)
		}

		orDone := func(done <-chan struct{}, c <-chan *nostr.Event) <-chan *nostr.Event {
			valStream := make(chan *nostr.Event)
			go func() {
				defer close(valStream)
				for {
					select {
					case <-done:
						return
					case v, ok := <-c:
						if ok == false {
							return
						}
						valStream <- v
					}
				}
			}()
			return valStream
		}

		for e := range orDone(sub.Done, sub.EventStream) {
			events = append(events, e)
		}

		//cc.Close()
	}

	log.Println("event added to local cache")
	log.Println(events)

	return events
}
