package main

import (
	"bufio"
	"godiscards/internal/cmdregistry"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"

	_ "godiscards/internal/commands"
)

var (
	removeCommands = false
)

func LoadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		log.Print(key, value)
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := LoadEnv(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	botToken, exists := os.LookupEnv("BOT_TOKEN")
	if !exists {
		log.Fatalf("Bot Token not set")
	}

	isDev, exists := os.LookupEnv("IS_DEV")
	if !exists {
		log.Fatalf("Is Dev not set")
	}

	var guildIds []string

	if isDev == "true" {
		devGuildId, exists := os.LookupEnv("DEV_GUILDID")
		if !exists {
			log.Fatalf("Dev GuildID not set")
		}

		guildIds = append(guildIds, devGuildId)
	}

	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(cmdregistry.CallCommandHandler)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Connot open session: %v", err)
	}

	commandsToBeRegistered := cmdregistry.GetCommandsToBeRegistered()
	var registeredCommands []*discordgo.ApplicationCommand

	for _, v := range commandsToBeRegistered {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildIds[0], v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v, err)
		}

		registeredCommands = append(registeredCommands, cmd)
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if removeCommands {
		log.Println("Removing commands...")
		for _, v := range commandsToBeRegistered {
			if err := s.ApplicationCommandDelete(s.State.User.ID, guildIds[0], v.ID); err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
