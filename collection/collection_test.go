package collection

import (
	"crypto/rand"
	"strconv"
	"testing"
)

var c *OrderedCollection

func init() {
	c = NewCollection()

	for i := 0; i < 100; i++ {
		b := make([]byte, 8)
		rand.Read(b)

		item := &Tuple{
			Key:   b,
			Value: []byte("this is value number " + strconv.Itoa(i)),
			Time:  int64(i),
		}

		c.Insert(item)
	}

}

func contains(src []*Tuple, key []byte) bool {

	match := string(key)

	for _, c := range src {
		if string(c.Key) == match {
			return true
		}
	}

	return false
}

func find(src []*Tuple, key []byte) *Tuple {

	match := string(key)

	for _, c := range src {
		if string(c.Key) == match {
			return c
		}
	}

	return nil
}

func compare(n []byte, m []byte) bool {
	return string(n) == string(m)
}

func TestAddDelete(t *testing.T) {

	key := []byte("HEYKEY")

	item := &Tuple{
		Key:   key,
		Value: []byte("some kidna value"),
		Time:  87,
	}

	c.Insert(item)

	before := c.Get()

	if !contains(before, key) {
		t.Error("key should exist")
	}

	c.Delete(key)

	after := c.Get()

	if contains(after, key) {
		t.Error("key should no longer exist", after)
	}
}

func TestMutate(t *testing.T) {

	key := []byte("someotherkey")
	val1 := []byte("myvalue")
	val2 := []byte("somenewvalue")

	item := &Tuple{
		Key:   key,
		Value: val1,
		Time:  50,
	}

	item2 := &Tuple{
		Key:   key,
		Value: val2,
		Time:  60,
	}

	c.Insert(item)
	before := c.Get()

	c.Update(item2)
	after := c.Get()

	t1 := find(before, key)
	t2 := find(after, key)

	if compare(t1.Value, t2.Value) {
		t.Error("original result should not be changed")
	}
}
