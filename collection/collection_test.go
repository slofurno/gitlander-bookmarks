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
		c.Insert(int64(i), b, []byte("this is value number "+strconv.Itoa(i)))
	}

}

func TestAddDelete(t *testing.T) {

	c.Insert(87, []byte("HEYKEY"), []byte("some kidna value"))

	before := c.Get([]byte("HEYKEY"))

	if before == nil {
		t.Error("key should exist")
	}

	c.Delete([]byte("HEYKEY"))

	after := c.Get([]byte("HEYKEY"))

	if after != nil {
		t.Error("key should no longer exist", after)
	}

}
