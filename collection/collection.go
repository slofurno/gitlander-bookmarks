package collection

import (
	"errors"
	"sync"
)

type Tuple struct {
	Time  int64
	Key   []byte
	Value []byte
}

//stored by time ascending
type Collection interface {
	Insert(int64, []byte, []byte)
	Get([]byte) []byte
	Slice(int64, int) []*Tuple
	Delete([]byte)
}

func NewCollection() *OrderedCollection {
	return &OrderedCollection{
		lookup: make(map[string]*Tuple),
	}
}

type OrderedCollection struct {
	lookup map[string]*Tuple
	store  []*Tuple
	lock   sync.Mutex
}

func (s *OrderedCollection) Insert(time int64, key, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.lookup[string(key)] != nil {
		return errors.New("key already exists")
	}

	var i int
	store := s.store
	for i = len(store); i > 0; i-- {
		if store[i-1].Time < time {
			break
		}
	}

	n := &Tuple{time, key, value}
	next := make([]*Tuple, len(store)+1)
	copy(next, store[:i])
	copy(next[i+1:], store[i:])
	next[i] = n

	s.store = next
	s.lookup[string(key)] = n
	return nil
}

func (s *OrderedCollection) Get(key []byte) *Tuple {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.lookup[string(key)]
}

func (s *OrderedCollection) Update(key, value []byte) {

}

func (s *OrderedCollection) Slice(time int64, count int) []*Tuple {
	s.lock.Lock()
	defer s.lock.Unlock()

	store := s.store
	var i int
	for i = len(store); i > 0; i-- {
		if store[i-1].Time <= time {
			break
		}
	}

	start := i - count
	if start < 0 {
		start = 0
	}

	res := make([]*Tuple, count)
	copy(res, store[start:i])

	return res
}

func (s *OrderedCollection) Delete(key []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	item := s.lookup[string(key)]

	if item == nil {
		return
	}

	var i int

	for i = 0; i < len(s.store); i++ {
		if s.store[i] == item {
			break
		}
	}

	delete(s.lookup, string(key))
	s.store = append(s.store[:i], s.store[i+1:]...)
}
