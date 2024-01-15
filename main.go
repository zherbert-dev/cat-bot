package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

const CatAPIURL = "https://api.thecatapi.com/v1/images/search?category_ids=3&format=src"

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

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type Gopher struct {
	Name string `json:"name"`
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!cat" {

		//Call the KuteGo API and retrieve our cute Dr Who Gopher
		response, err := http.Get(CatAPIURL)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "cat.png", response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get Cat! :-(")
		}
	}

	if strings.Contains(m.Content, "!catsay") {
		cat, err := os.ReadFile("ascii.txt")
		if err != nil {
			fmt.Println(err)
		}

		message := m.Content[8:]
		content := message + "\n" + string(cat)

		s.ChannelMessageSend(m.ChannelID, content)
	}

	if m.Content == "!hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello there!")
	}
}
