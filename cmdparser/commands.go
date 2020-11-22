package cmdparser

type command struct {
	input  string
	action string
	body   string
}

func (c command) Input() string {
	return c.input
}

func (c command) Action() string {
	return c.action
}

func (c command) Body() string {
	return c.body
}

func newCommand(input, action, body string) command {
	return command{input, action, body}
}
