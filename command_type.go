package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	handlerFunctions map[string]func(*State, command) error
}

func (c *commands) run(s *State, cmd command) error {
	handler, ok := c.handlerFunctions[cmd.Name]
	if !ok {
		return errors.New("Unknown command")
	}
	err := handler(s, cmd)
	return err
}

func (c *commands) register(name string, f func(*State, command) error) {
	c.handlerFunctions[name] = f
}
