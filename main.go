package main

import (
	"Discord-bot/controller"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var cid string = "container-id"

func main() {
	Token := "token-here"
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	message := m.Content
	if message == "!h" {
		s.ChannelMessageSend(m.ChannelID, "!h - Help Menu\n"+
			"!s {commands} - Server Commands\n"+
			"{commands}\n"+
			"    update   - Update Image\n"+
			"    start       - Start Container\n"+
			"    stop       - Stop Container\n"+
			"    restart   - Restart Container\n"+
			"    status     - Container Status\n"+
			"    stats      - Server Stats\n\n"+
			"!m {minecraft-commands} - Minecraft Commands\n")
	}
	if strings.HasPrefix(message, "!s ") {
		message = message[3:]
		controller.Server(message, cid, s, m)
	}
}
