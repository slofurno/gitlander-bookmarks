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
	Insert(*Tuple) error
	Update(*Tuple) error
	Get() []*Tuple
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

func (s *OrderedCollection) Insert(item *Tuple) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.lookup[string(item.Key)] != nil {
		return errors.New("key already exists")
	}

	var i int
	store := s.store
	for i = len(store); i > 0; i-- {
		if store[i-1].Time < item.Time {
			break
		}
	}

	next := make([]*Tuple, len(store)+1)
	copy(next, store[:i])
	copy(next[i+1:], store[i:])
	next[i] = item

	s.store = next
	s.lookup[string(item.Key)] = item
	return nil
}

func (s *OrderedCollection) Get() []*Tuple {
	s.lock.Lock()
	defer s.lock.Unlock()

	res := make([]*Tuple, len(s.store))
	copy(res, s.store)
	return res
}

func (s *OrderedCollection) Update(item *Tuple) error {
	s.Delete(item.Key)
	return s.Insert(item)
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

	res := make([]*Tuple, len(s.store)-1)
	copy(res, s.store[:i])
	copy(res[i:], s.store[i+1:])

	delete(s.lookup, string(key))
	//s.store = append(s.store[:i], s.store[i+1:]...)

	s.store = res
}
