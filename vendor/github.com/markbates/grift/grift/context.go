package grift

// Context used to pass information between grifts
type Context struct {
	Name    string
	Args    []string
	Verbose bool
	data    map[string]interface{}
}

// Get a piece of data from the Context.
func (c *Context) Get(key string) interface{} {
	return c.data[key]
}

// Set a piece of data onto the Context.
func (c *Context) Set(key string, val interface{}) {
	c.data[key] = val
}

// NewContext builds and returns a new default Context.
func NewContext(name string) *Context {
	return &Context{
		Name: name,
		Args: []string{},
		data: map[string]interface{}{},
	}
}
