package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"strings"

	"github.com/Emyrk/progbot/playground"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// TODO: Maybe use regex?
	if len(m.Content) > 2 && m.Content[:2] == ">c" {
		command := m.Content[3:]
		res := strings.SplitN(command, " ", 2)
		if len(res) < 2 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Invalid command. Try `|compile LANG CODE`\nCom: %#v, Res: %#v", command, res))
			return
		}

		switch strings.ToLower(res[0]) {
		case "go", "golang":
			handleGo(s, m, res[1])
		case "e", "elixir":
			handleElixir(s, m, res[1])
		default:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Invalid language, found %s", res[0]))
			return
		}
	} else if m.Content == "hello" {
		s.ChannelMessageSend(m.ChannelID, "Hi!")
	}
}

func handleElixir(s *discordgo.Session, m *discordgo.MessageCreate, code string) {
	res, err := ExecuteElixirWithTimeout(wrapelixircode(code), 10)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Compiled and ran:\n``` %s ```\n", res))
}

func wrapelixircode(code string) string {
	code = strings.Replace(code, "`", "", -1)
	return code
}
func wrapgocode(code string) string {
	code = strings.Replace(code, "`", "", -1)
	code = `
package main
func main() {

}


			` //+ code
	return code
}

func handleGo(s *discordgo.Session, m *discordgo.MessageCreate, code string) {
	resp, err := playground.CompileAndRun(&playground.Request{
		Body: wrapgocode(code),
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%#v\n", resp))
}
