package cmdregistry

import "github.com/bwmarrin/discordgo"

type commandHandler func(*discordgo.Session, *discordgo.InteractionCreate)

var (
	commands          []*discordgo.ApplicationCommand
	commandHandlerMap = make(map[string]commandHandler)
)

func RegisterCommand(ac *discordgo.ApplicationCommand, h commandHandler) {
	commands = append(commands, ac)
	commandHandlerMap[ac.Name] = h
}

func GetCommandsToBeRegistered() []*discordgo.ApplicationCommand {
	return commands
}

func GetCommandHandlers() map[string]commandHandler {
	return commandHandlerMap
}

func CallCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandHandlerMap[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
