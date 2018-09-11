package bot

import (
	"fmt"
)

type commandFunc func(from string, params []string) string

type Command struct {
	Name        string
	Description string
	Commands    []*Command
	Action      commandFunc
}

func (c Command) Run(from string, args []string) string {
	if len(args) > 0 {
		cmdName := args[0]
		if cmdName == "help" {
			return c.Help()
		}
		if len(c.Commands) > 0 {
			for _, cmd := range c.Commands {
				if cmd.Name == cmdName {
					return cmd.Run(from, args[1:])
				}
			}
			return fmt.Sprintf("command %s not found\n%s", cmdName, c.Help())
		}
	}
	if c.Action != nil {
		return c.Action(from, args)
	}
	return c.Help()
}

func (c Command) Help() string {
	if len(c.Commands) > 0 {
		msg := fmt.Sprintf("%s\n-------------------", c.Description)
		for _, cmd := range c.Commands {
			msg = fmt.Sprintf("%s\n%s: %s", msg, cmd.Name, cmd.Description)
		}
		return msg
	}
	return c.Description
}
