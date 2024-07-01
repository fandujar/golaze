package golaze

import (
	"context"
	"sync"
)

type State struct {
	Data map[string]interface{}
	lock sync.Mutex
}

func (s *State) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Data[key] = value
}

func (s *State) Get(key string) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Data[key]
}

func (s *State) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.Data, key)
}

func (s *State) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Data = make(map[string]interface{})
}

func (s *State) Context() context.Context {
	return context.WithValue(context.Background(), "state", s)
}
