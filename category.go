package cli

// CommandCategory is a category containing commands.
type CommandCategory struct {
	Name     string
	Commands []*Command
}

// VisibleCommands returns a slice of the Commands with Hidden=false
func (category *CommandCategory) VisibleCommands() []*Command {
	cmds := []*Command{}
	for _, command := range category.Commands {
		if !command.Hidden {
			cmds = append(cmds, command)
		}
	}
	return cmds
}
