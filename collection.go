package main

type Callback struct {
	added   func(string, interface{})
	changed func(string, interface{})
	removed func(string, interface{})
}

type Cursor interface {
	fetch2() []interface{}
}

type Collection struct {
	store     map[string]interface{}
	callbacks []*Callback
	events    chan collectionEvent
}

func newCollection() *Collection {
	c := &Collection{
		store:     make(map[string]interface{}),
		callbacks: make([]*Callback, 0),
		events:    make(chan collectionEvent, 256),
	}

	go c.eventLoop()

	return c
}

func (c *Collection) eventLoop() {

	for {
		select {
		case e := <-c.events:
			e()
		}
	}

}

type collectionEvent func()

func (c *Collection) Fetch() []interface{} {

	result := make(chan []interface{}, 1)

	fetch := func() {
		dump := []interface{}{}
		for _, value := range c.store {
			dump = append(dump, value)
		}
		result <- dump
	}

	c.events <- fetch
	return <-result

}

func (c *Collection) Add(key string, value interface{}) {

	add := func() {

		if _, ok := c.store[key]; !ok {
			c.store[key] = value
			for _, f := range c.callbacks {
				f.added(key, value)
			}
		} else {
			c.store[key] = value
			for _, f := range c.callbacks {
				f.changed(key, value)
			}
		}
		//TODO: is this how i want to handle repeat adds?
	}

	c.events <- add
}

func (c *Collection) Update(key string, value interface{}) {

	update := func() {
		c.store[key] = value

		for _, f := range c.callbacks {
			f.changed(key, value)
		}
	}

	c.events <- update
}

func (c *Collection) ObserveChanges(callback *Callback) func() {

	//TODO:instead of n calls on subscribe, maybe 1 call with n elements
	add := func() {
		for key, el := range c.store {
			callback.added(key, el)
		}
		c.callbacks = append(c.callbacks, callback)
	}

	c.events <- add

	rem := func() {
		for i := 0; i < len(c.callbacks); i++ {
			if c.callbacks[i] == callback {
				c.callbacks = append(c.callbacks[:i], c.callbacks[i+1:]...)
				return
			}
		}
	}

	onstop := func() {
		c.events <- rem
	}

	return onstop
}
