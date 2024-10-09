package commands

import (
	"godiscards/internal/cmdregistry"

	"github.com/bwmarrin/discordgo"
)

func init() {
	cmdregistry.RegisterCommand(
		&discordgo.ApplicationCommand{
			Name:        "drop",
			Description: "Get a random card",
		},
		dropHandler,
	)
}
func dropHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hi",
		},
	})
	return
}
