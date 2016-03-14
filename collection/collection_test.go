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
			Key:   string(b),
			Value: "this is value number " + strconv.Itoa(i),
			Time:  int64(i),
		}

		c.Insert(item)
	}

}

func contains(src []*Tuple, match string) bool {

	for _, c := range src {
		if c.Key == match {
			return true
		}
	}

	return false
}

func find(src []*Tuple, match string) *Tuple {
	for _, c := range src {
		if c.Key == match {
			return c
		}
	}
	return nil
}

func compare(n, m string) bool {
	return n == m
}

func TestAddDelete(t *testing.T) {

	key := "HEYKEY"

	item := &Tuple{
		Key:   key,
		Value: "some kidna value",
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

	key := "someotherkey"
	val1 := "myvalue"
	val2 := "somenewvalue"

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

func TestCrc(t *testing.T) {

	b := make([]byte, 16)
	count := make(map[int]int)

	for i := 0; i < 9999; i++ {
		rand.Read(b)
		c := Crc16(b) % 4

		count[c] = count[c] + 1
	}

}
